// +build darwin,amd64

// Package mac is the driver to be used for apps that run on MacOS.
// It is build on the top of Cocoa and Webkit.
package mac

/*
#include "driver.h"
#include "bridge.h"
*/
import "C"
import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/bridge"
	"github.com/pkg/errors"
)

var (
	driver *Driver
)

func init() {
	app.Loggers = []app.Logger{
		app.NewLogger(os.Stdout, os.Stderr, true, true),
	}
}

// Driver is the app.Driver implementation for MacOS.
type Driver struct {
	app.BaseDriver

	// Menubar configuration
	MenubarConfig MenuBarConfig

	// The URL of the component to load in the dock.
	DockURL string

	// The bundle configuration.
	// Only applied when the app is build with goapp mac build -bundle.
	Bundle Bundle

	// The handler called when the app is running.
	OnRun func()

	// The handler called when the app is focused.
	OnFocus func()

	// The handler called when the app loses focus.
	OnBlur func()

	// The handler called when the app is reopened.
	OnReopen func(hasVisibleWindows bool)

	// The handler called when a file associated with the app is opened.
	OnFilesOpen func(filenames []string)

	// The handler called when the app URI is invoked.
	OnURLOpen func(u *url.URL)

	// The handler called when the quit button is clicked.
	OnQuit func() bool

	// The handler called when the app is about to exit.
	OnExit func()

	factory      app.Factory
	elements     app.ElemDB
	uichan       chan func()
	macRPC       bridge.PlatformRPC
	goRPC        bridge.GoRPC
	menubar      app.Menu
	statusbar    app.StatusBarMenu
	dock         app.DockTile
	devID        string
	droppedFiles []string
}

// Name satisfies the app.Driver interface.
func (d *Driver) Name() string {
	return "MacOS"
}

// Run satisfies the app.Driver interface.
func (d *Driver) Run(f app.Factory) error {
	if bundle := os.Getenv("GOAPP_BUNDLE"); bundle == "true" {
		data, err := json.MarshalIndent(d.Bundle, "", "  ")
		if err != nil {
			os.Exit(1)
		}
		return ioutil.WriteFile("goapp-mac.json", data, 0777)
	}

	if driver != nil {
		return errors.Errorf("driver is already running")
	}

	d.devID = generateDevID()
	d.factory = f

	elements := app.NewElemDB()
	elements = app.ConcurrentElemDB(elements)
	d.elements = elements

	d.uichan = make(chan func(), 4096)
	defer close(d.uichan)

	d.macRPC.Handler = macCall

	d.goRPC.Handle("driver.OnRun", d.onRun)
	d.goRPC.Handle("driver.OnFocus", d.onFocus)
	d.goRPC.Handle("driver.OnBlur", d.onBlur)
	d.goRPC.Handle("driver.OnReopen", d.onReopen)
	d.goRPC.Handle("driver.OnFilesOpen", d.onFilesOpen)
	d.goRPC.Handle("driver.OnURLOpen", d.onURLOpen)
	d.goRPC.Handle("driver.OnFileDrop", d.onFileDrop)
	d.goRPC.Handle("driver.OnQuit", d.onQuit)
	d.goRPC.Handle("driver.OnExit", d.onExit)

	d.goRPC.Handle("windows.OnMove", handleWindow(onWindowMove))
	d.goRPC.Handle("windows.OnResize", handleWindow(onWindowResize))
	d.goRPC.Handle("windows.OnFocus", handleWindow(onWindowFocus))
	d.goRPC.Handle("windows.OnBlur", handleWindow(onWindowBlur))
	d.goRPC.Handle("windows.OnFullScreen", handleWindow(onWindowFullScreen))
	d.goRPC.Handle("windows.OnExitFullScreen", handleWindow(onWindowExitFullScreen))
	d.goRPC.Handle("windows.OnMinimize", handleWindow(onWindowMinimize))
	d.goRPC.Handle("windows.OnDeminimize", handleWindow(onWindowDeminimize))
	d.goRPC.Handle("windows.OnClose", handleWindow(onWindowClose))
	d.goRPC.Handle("windows.OnCallback", handleWindow(onWindowCallback))
	d.goRPC.Handle("windows.OnNavigate", handleWindow(onWindowNavigate))

	d.goRPC.Handle("menus.OnClose", handleMenu(onMenuClose))
	d.goRPC.Handle("menus.OnCallback", handleMenu(onMenuCallback))

	d.goRPC.Handle("filePanels.OnSelect", handleFilePanel(onFilePanelSelect))
	d.goRPC.Handle("saveFilePanels.OnSelect", handleSaveFilePanel(onSaveFilePanelSelect))

	d.goRPC.Handle("notifications.OnReply", handleNotification(onNotificationReply))

	driver = d
	errC := make(chan error)

	go func() {
		errC <- d.macRPC.Call("driver.Run", nil, nil)
	}()

	for {
		select {
		case f := <-d.uichan:
			f()

		case err := <-errC:
			return err
		}
	}
}

func (d *Driver) onRun(in map[string]interface{}) interface{} {
	err := d.newMenuBar()
	if err != nil {
		panic(err)
	}

	if d.statusbar, err = newStatusBar(); err != nil {
		panic(err)
	}

	if d.dock, err = newDockTile(app.MenuConfig{
		DefaultURL: d.DockURL,
	}); err != nil {
		panic(err)
	}

	if d.OnRun != nil {
		d.OnRun()
	}
	return nil
}

func (d *Driver) onFocus(in map[string]interface{}) interface{} {
	if d.OnFocus != nil {
		d.OnFocus()
	}
	return nil
}

func (d *Driver) onBlur(in map[string]interface{}) interface{} {
	if d.OnBlur != nil {
		d.OnBlur()
	}
	return nil
}

func (d *Driver) onReopen(in map[string]interface{}) interface{} {
	if d.OnReopen != nil {
		d.OnReopen(in["HasVisibleWindows"].(bool))
	}
	return nil
}

func (d *Driver) onFilesOpen(in map[string]interface{}) interface{} {
	if d.OnFilesOpen != nil {
		d.OnFilesOpen(bridge.Strings(in["Filenames"]))
	}
	return nil
}

func (d *Driver) onURLOpen(in map[string]interface{}) interface{} {
	if d.OnURLOpen != nil {
		u, err := url.Parse(in["URL"].(string))
		if err != nil {
			panic(errors.Wrap(err, "onURLOpen"))
		}
		d.OnURLOpen(u)
	}
	return nil
}

func (d *Driver) onFileDrop(in map[string]interface{}) interface{} {
	d.droppedFiles = bridge.Strings(in["Filenames"])
	return nil
}

// AppName satisfies the app.Driver interface.
func (d *Driver) AppName() string {
	var out struct {
		AppName string
	}

	if err := d.macRPC.Call("driver.Bundle", &out, nil); err != nil {
		panic(err)
	}

	if len(out.AppName) != 0 {
		return out.AppName
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(errors.Wrap(err, "app name unreachable"))
	}
	return filepath.Base(wd)
}

// Resources satisfies the app.Driver interface.
func (d *Driver) Resources(path ...string) string {
	var out struct {
		Resources string
	}

	if err := d.macRPC.Call("driver.Bundle", &out, nil); err != nil {
		panic(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(errors.Wrap(err, "resources unreachable"))
	}

	if filepath.Dir(out.Resources) == wd {
		out.Resources = filepath.Join(wd, "resources")
	}

	resources := filepath.Join(path...)
	return filepath.Join(out.Resources, resources)
}

// Storage satisfies the app.DriverWithStorage interface.
func (d *Driver) Storage(path ...string) string {
	storage := filepath.Join(path...)
	return filepath.Join(d.support(), "storage", storage)
}

func (d *Driver) support() string {
	var out struct {
		Support string
	}

	if err := d.macRPC.Call("driver.Bundle", &out, nil); err != nil {
		panic(err)
	}

	// Set up the support directory in case of the app is not bundled.
	if strings.HasSuffix(out.Support, "{appname}") {
		wd, err := os.Getwd()
		if err != nil {
			panic(errors.Wrap(err, "support unreachable"))
		}

		appname := filepath.Base(wd) + "-" + d.devID
		out.Support = strings.Replace(out.Support, "{appname}", appname, 1)
	}
	return out.Support
}

// NewWindow satisfies the app.DriverWithWindows interface.
func (d *Driver) NewWindow(c app.WindowConfig) (app.Window, error) {
	return newWindow(c)
}

// NewContextMenu satisfies the app.Driver interface.
func (d *Driver) NewContextMenu(c app.MenuConfig) (app.Menu, error) {
	c.Type = "context menu"

	m, err := newMenu(c)
	if err != nil {
		return nil, err
	}

	err = d.macRPC.Call("driver.SetContextMenu", nil, m.ID())
	return m, err
}

// Render satisfies the app.Driver interface.
func (d *Driver) Render(c app.Component) error {
	e, err := d.elements.ElementByComponent(c)
	if err != nil {
		return err
	}
	return e.Render(c)
}

// ElementByComponent satisfies the app.Driver interface.
func (d *Driver) ElementByComponent(c app.Component) (app.ElementWithComponent, error) {
	return d.elements.ElementByComponent(c)
}

// NewFilePanel satisfies the app.DriverWithFilePanels interface.
func (d *Driver) NewFilePanel(c app.FilePanelConfig) error {
	return newFilePanel(c)
}

// NewSaveFilePanel satisfies the app.DriverWithFilePanels interface.
func (d *Driver) NewSaveFilePanel(c app.SaveFilePanelConfig) error {
	return newSaveFilePanel(c)
}

// NewShare satisfies the app.DriverWithShare interface.
func (d *Driver) NewShare(v interface{}) error {
	in := struct {
		Share string
		Type  string
	}{
		Share: fmt.Sprint(v),
	}

	switch v.(type) {
	case url.URL, *url.URL:
		in.Type = "url"

	default:
		in.Type = "string"
	}

	return d.macRPC.Call("driver.Share", nil, in)
}

// NewNotification satisfies the app.DriverWithPopupNotifications
// interface.
func (d *Driver) NewNotification(config app.NotificationConfig) error {
	return newNotification(config)
}

// MenuBar satisfies the app.DriverWithMenuBar interface.
func (d *Driver) MenuBar() (app.Menu, error) {
	return d.menubar, nil
}

func (d *Driver) newMenuBar() error {
	menubar, err := newMenu(app.MenuConfig{Type: "menubar"})
	if err != nil {
		return errors.Wrap(err, "creating the menu bar failed")
	}
	d.menubar = menubar

	if len(d.MenubarConfig.URL) == 0 {
		format := "mac.menubar?appurl=%s&editurl=%s&windowurl=%s&helpurl=%s"
		for _, u := range d.MenubarConfig.CustomURLs {
			format += "&custom=" + u
		}

		err = menubar.Load(
			format,
			d.MenubarConfig.AppURL,
			d.MenubarConfig.EditURL,
			d.MenubarConfig.WindowURL,
			d.MenubarConfig.HelpURL,
		)
	} else {
		err = menubar.Load(d.MenubarConfig.URL)
	}
	if err != nil {
		return err
	}

	if err = d.macRPC.Call("driver.SetMenubar", nil, menubar.ID()); err != nil {
		return errors.Wrap(err, "set menu bar")
	}
	return nil
}

func (d *Driver) StatusBar() (app.StatusBarMenu, error) {
	panic("not implemented")
}

// Dock satisfies the app.DriverWithDock interface.
func (d *Driver) Dock() (app.DockTile, error) {
	return d.dock, nil
}

// CallOnUIGoroutine satisfies the app.Driver interface.
func (d *Driver) CallOnUIGoroutine(f func()) {
	d.uichan <- f
}

// Close quits the app.
func (d *Driver) Close() error {
	return d.macRPC.Call("driver.Quit", nil, nil)
}

func (d *Driver) onQuit(in map[string]interface{}) interface{} {
	out := struct {
		Quit bool
	}{
		Quit: true,
	}

	if d.OnQuit != nil {
		out.Quit = d.OnQuit()
	}
	return out
}

func (d *Driver) onExit(in map[string]interface{}) interface{} {
	if d.OnExit != nil {
		d.OnExit()
	}
	return nil
}

func generateDevID() string {
	h := md5.New()
	wd, _ := os.Getwd()
	io.WriteString(h, wd)
	return fmt.Sprintf("%x", h.Sum(nil))
}
