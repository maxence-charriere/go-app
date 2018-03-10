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
	"fmt"
	"io"
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

// Driver is the app.Driver implementation for MacOS.
type Driver struct {
	app.BaseDriver

	// Menubar configuration
	MenubarConfig MenuBarConfig

	// Component url to load in the dock.
	DockURL string

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
	macos        bridge.PlatformBridge
	golang       bridge.GoBridge
	macRPC       bridge.RPC
	menubar      app.Menu
	dock         app.DockTile
	devID        string
	droppedFiles []string
}

// Name satisfies the app.Driver interface.
func (d *Driver) Name() string {
	return "MacOS"
}

// Base satisfies the app.Driver interface.
func (d *Driver) Base() app.Driver {
	return d
}

// Run satisfies the app.Driver interface.
func (d *Driver) Run(f app.Factory) error {
	if driver != nil {
		return errors.Errorf("driver is already running")
	}

	d.devID = generateDevID()
	d.factory = f

	elements := app.NewElemDB()
	elements = app.NewConcurrentElemDB(elements)
	d.elements = elements

	d.uichan = make(chan func(), 4096)
	defer close(d.uichan)

	d.macos = bridge.NewPlatformBridge(handleMacOSRequest)
	d.golang = bridge.NewGoBridge(d.uichan)

	d.macRPC.Handler = macCall

	d.golang.Handle("/driver/run", d.onRun)
	d.golang.Handle("/driver/focus", d.onFocus)
	d.golang.Handle("/driver/blur", d.onBlur)
	d.golang.Handle("/driver/reopen", d.onReopen)
	d.golang.Handle("/driver/filesopen", d.onFilesOpen)
	d.golang.Handle("/driver/urlopen", d.onURLOpen)
	d.golang.Handle("/driver/filedrop", d.onFileDrop)
	d.golang.Handle("/driver/quit", d.onQuit)
	d.golang.Handle("/driver/exit", d.onExit)

	d.golang.Handle("/window/move", windowHandler(onWindowMove))
	d.golang.Handle("/window/resize", windowHandler(onWindowResize))
	d.golang.Handle("/window/focus", windowHandler(onWindowFocus))
	d.golang.Handle("/window/blur", windowHandler(onWindowBlur))
	d.golang.Handle("/window/fullscreen", windowHandler(onWindowFullScreen))
	d.golang.Handle("/window/fullscreen/exit", windowHandler(onWindowExitFullScreen))
	d.golang.Handle("/window/minimize", windowHandler(onWindowMinimize))
	d.golang.Handle("/window/deminimize", windowHandler(onWindowDeminimize))
	d.golang.Handle("/window/close", windowHandler(onWindowClose))
	d.golang.Handle("/window/callback", windowHandler(onWindowCallback))
	d.golang.Handle("/window/navigate", windowHandler(onWindowNavigate))

	d.golang.Handle("/menu/close", menuHandler(onMenuClose))
	d.golang.Handle("/menu/callback", menuHandler(onMenuCallback))

	d.golang.Handle("/file/panel/select", filePanelHandler(onFilePanelClose))
	d.golang.Handle("/file/savepanel/select", saveFilePanelHandler(onSaveFilePanelClose))

	d.golang.Handle("/notification/reply", notificationHandler(onNotificationReply))

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

func (d *Driver) onRun(u *url.URL, p bridge.Payload) (res bridge.Payload) {
	err := d.newMenuBar()
	if err != nil {
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

func (d *Driver) onFocus(u *url.URL, p bridge.Payload) (res bridge.Payload) {
	if d.OnFocus == nil {
		return nil
	}

	d.OnFocus()
	return nil
}

func (d *Driver) onBlur(u *url.URL, p bridge.Payload) (res bridge.Payload) {
	if d.OnBlur == nil {
		return nil
	}

	d.OnBlur()
	return nil
}

func (d *Driver) onReopen(u *url.URL, p bridge.Payload) (res bridge.Payload) {
	if d.OnReopen == nil {
		return nil
	}

	var hasVisibleWindows bool
	p.Unmarshal(&hasVisibleWindows)
	d.OnReopen(hasVisibleWindows)
	return nil
}

func (d *Driver) onFilesOpen(u *url.URL, p bridge.Payload) (res bridge.Payload) {
	if d.OnFilesOpen == nil {
		return nil
	}

	var filenames []string
	p.Unmarshal(&filenames)
	d.OnFilesOpen(filenames)
	return nil
}

func (d *Driver) onURLOpen(u *url.URL, p bridge.Payload) (res bridge.Payload) {
	if d.OnURLOpen == nil {
		return nil
	}

	var rawurl string
	p.Unmarshal(&rawurl)

	purl, err := url.Parse(rawurl)
	if err != nil {
		panic(errors.Wrap(err, "parsing url failed"))
	}

	d.OnURLOpen(purl)
	return nil
}

func (d *Driver) onFileDrop(u *url.URL, p bridge.Payload) (res bridge.Payload) {
	p.Unmarshal(&d.droppedFiles)
	return nil
}

// AppName satisfies the app.Driver interface.
func (d *Driver) AppName() string {
	var out struct {
		AppName string
	}

	if err := d.macRPC.Call("driver.Bundle", nil, &out); err != nil {
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

	if err := d.macRPC.Call("driver.Bundle", nil, &out); err != nil {
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

	if err := d.macRPC.Call("driver.Bundle", nil, &out); err != nil {
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
	m, err := newMenu(c, "context menu")
	if err != nil {
		return nil, err
	}

	in := struct {
		MenuID string
	}{
		MenuID: m.ID().String(),
	}

	err = d.macRPC.Call("driver.SetContextMenu", in, nil)
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

	return d.macRPC.Call("driver.Share", in, nil)
}

// NewNotification satisfies the app.DriverWithPopupNotifications
// interface.
func (d *Driver) NewNotification(config app.NotificationConfig) error {
	_, err := newNotification(config)
	return err
}

// MenuBar satisfies the app.DriverWithMenuBar interface.
func (d *Driver) MenuBar() (app.Menu, error) {
	return d.menubar, nil
}

func (d *Driver) newMenuBar() error {
	menubar, err := newMenu(app.MenuConfig{}, "menubar")
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

	in := struct {
		MenuID string
	}{
		MenuID: menubar.ID().String(),
	}

	if err = d.macRPC.Call("driver.SetMenubar", in, nil); err != nil {
		return errors.Wrap(err, "set menu bar")
	}
	return nil
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
func (d *Driver) Close() {
	d.macRPC.Call("driver.Quit", nil, nil)
}

func (d *Driver) onQuit(u *url.URL, p bridge.Payload) (res bridge.Payload) {
	quit := true

	if d.OnQuit != nil {
		quit = d.OnQuit()
	}

	res = bridge.NewPayload(quit)
	return
}

func (d *Driver) onExit(u *url.URL, p bridge.Payload) (res bridge.Payload) {
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
