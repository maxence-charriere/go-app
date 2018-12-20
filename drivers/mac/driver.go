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
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/bridge"
	"github.com/murlokswarm/app/internal/core"
	"github.com/pkg/errors"
)

var (
	driver     *Driver
	goappBuild = os.Getenv("GOAPP_BUILD")
	debug      = os.Getenv("GOAPP_DEBUG") == "true"
)

func init() {
	if len(goappBuild) != 0 {
		app.Logger = func(format string, a ...interface{}) {}
		return
	}

	logger := core.ToWriter(os.Stderr)
	app.Logger = core.WithColoredPrompt(logger)
	app.EnableDebug(debug)
}

// Driver is the app.Driver implementation for MacOS.
type Driver struct {
	core.Driver

	// The app package settings.
	// It is used by goapp to define how the app is built and packaged.
	Settings Settings

	// The URL of the component to load in the default window. A non empty value
	// triggers the creation of the default window when the app in openened. It
	// overrides the DefaultWindow.URL value.
	URL string

	// The default window configuration.
	DefaultWindow app.WindowConfig

	// Menubar configuration
	MenubarConfig MenuBarConfig

	// The URL of the component to load in the dock.
	DockURL string

	// The func called when the app is focused.
	OnFocus func()

	// The func called when the app loses focus.
	OnBlur func()

	// The func called when files associated with the app are opened.
	OnFilesOpen func(filenames []string)

	// The func called when the app URLScheme is invoked.
	OnURLOpen func(u *url.URL)

	// The func called when the app is about to quit.
	OnQuit func()

	ui           chan func()
	factory      *app.Factory
	events       *app.EventRegistry
	elems        *core.ElemDB
	devID        string
	macRPC       bridge.PlatformRPC
	goRPC        bridge.GoRPC
	stop         func()
	menubar      *Menu
	docktile     *DockTile
	droppedFiles []string
}

// Run satisfies the app.Driver interface.
func (d *Driver) Run(c app.DriverConfig) error {
	if len(goappBuild) != 0 {
		return d.runGoappBuild()
	}

	if driver != nil {
		return errors.New("running already")
	}

	d.ui = c.UI
	d.factory = c.Factory
	d.events = c.Events
	d.elems = core.NewElemDB()
	d.devID = generateDevID()
	d.macRPC.Handler = macCall

	d.goRPC.Handle("driver.OnRun", d.onRun)
	d.goRPC.Handle("driver.OnFocus", d.onFocus)
	d.goRPC.Handle("driver.OnBlur", d.onBlur)
	d.goRPC.Handle("driver.OnReopen", d.onReopen)
	d.goRPC.Handle("driver.OnFilesOpen", d.onFilesOpen)
	d.goRPC.Handle("driver.OnURLOpen", d.onURLOpen)
	d.goRPC.Handle("driver.OnFileDrop", d.onFileDrop)
	d.goRPC.Handle("driver.OnQuit", d.onQuit)

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
	d.goRPC.Handle("windows.OnAlert", handleWindow(onWindowAlert))

	d.goRPC.Handle("menus.OnClose", handleMenu(onMenuClose))
	d.goRPC.Handle("menus.OnCallback", handleMenu(onMenuCallback))

	d.goRPC.Handle("controller.OnDirectionChange", handleController(onControllerDirectionChange))
	d.goRPC.Handle("controller.OnButtonPressed", handleController(onControllerButtonPressed))
	d.goRPC.Handle("controller.OnConnected", handleController(onControllerConnected))
	d.goRPC.Handle("controller.OnDisconnected", handleController(onControllerDisconnected))
	d.goRPC.Handle("controller.OnPause", handleController(onControllerPause))
	d.goRPC.Handle("controller.OnClose", handleController(onControllerClose))

	d.goRPC.Handle("filePanels.OnSelect", handleFilePanel(onFilePanelSelect))
	d.goRPC.Handle("saveFilePanels.OnSelect", handleSaveFilePanel(onSaveFilePanelSelect))

	d.goRPC.Handle("notifications.OnReply", handleNotification(onNotificationReply))

	driver = d

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	d.stop = cancel

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for {
			select {
			case <-ctx.Done():
				wg.Done()
				return

			case fn := <-d.ui:
				fn()
			}
		}
	}()

	err := d.macRPC.Call("driver.Run", nil, nil)
	wg.Wait()
	return err
}

// AppName satisfies the app.Driver interface.
func (d *Driver) AppName() string {
	out := struct {
		AppName string
	}{}

	if err := d.macRPC.Call("driver.Bundle", &out, nil); err != nil {
		app.Panic(err)
	}

	if len(out.AppName) != 0 {
		return out.AppName
	}

	wd, err := os.Getwd()
	if err != nil {
		app.Panic(errors.Wrap(err, "app name unreachable"))
	}

	return filepath.Base(wd)
}

// Resources satisfies the app.Driver interface.
func (d *Driver) Resources(path ...string) string {
	out := struct {
		Resources string
	}{}

	if err := d.macRPC.Call("driver.Bundle", &out, nil); err != nil {
		app.Panic(err)
	}

	r := filepath.Join(path...)
	return filepath.Join(out.Resources, r)
}

// Storage satisfies the app.Driver interface.
func (d *Driver) Storage(path ...string) string {
	s := filepath.Join(path...)
	return filepath.Join(d.support(), "storage", s)
}

// Render satisfies the app.Driver interface.
func (d *Driver) Render(c app.Compo) {
	e := d.ElemByCompo(c)
	e.(app.ElemWithCompo).Render(c)
}

// ElemByCompo satisfies the app.Driver interface.
func (d *Driver) ElemByCompo(c app.Compo) app.Elem {
	return d.elems.GetByCompo(c)
}

// NewWindow satisfies the app.Driver interface.
func (d *Driver) NewWindow(c app.WindowConfig) app.Window {
	return newWindow(c)
}

// NewContextMenu satisfies the app.Driver interface.
func (d *Driver) NewContextMenu(c app.MenuConfig) app.Menu {
	m := newMenu(c, "context menu")
	if m.Err() != nil {
		return m
	}

	err := d.macRPC.Call("driver.SetContextMenu", nil, m.ID())
	m.SetErr(err)
	return m
}

// NewController statisfies the app.Driver interface.
func (d *Driver) NewController(c app.ControllerConfig) app.Controller {
	return newController(c)
}

// NewFilePanel satisfies the app.Driver interface.
func (d *Driver) NewFilePanel(c app.FilePanelConfig) app.Elem {
	return newFilePanel(c)
}

// NewSaveFilePanel satisfies the app.Driver interface.
func (d *Driver) NewSaveFilePanel(c app.SaveFilePanelConfig) app.Elem {
	return newSaveFilePanel(c)
}

// NewShare satisfies the app.Driver interface.
func (d *Driver) NewShare(v interface{}) app.Elem {
	return newSharePanel(v)
}

// NewNotification satisfies the app.Driver interface.
func (d *Driver) NewNotification(c app.NotificationConfig) app.Elem {
	return newNotification(c)
}

// MenuBar satisfies the app.Driver interface.
func (d *Driver) MenuBar() app.Menu {
	return d.menubar
}

// NewStatusMenu satisfies the app.Driver interface.
func (d *Driver) NewStatusMenu(c app.StatusMenuConfig) app.StatusMenu {
	return newStatusMenu(c)
}

// DockTile satisfies the app.Driver interface.
func (d *Driver) DockTile() app.DockTile {
	return d.docktile
}

// CallOnUIGoroutine satisfies the app.Driver interface.
func (d *Driver) CallOnUIGoroutine(f func()) {
	d.ui <- f
}

// Stop satisfies the app.Driver interface.
func (d *Driver) Stop() {
	if err := d.macRPC.Call("driver.Quit", nil, nil); err != nil {
		app.Log("stop failed:", err)
		d.stop()
	}
}

func (d *Driver) runGoappBuild() error {
	b, err := json.MarshalIndent(d.Settings, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(goappBuild, b, 0777)
}

func (d *Driver) support() string {
	out := struct {
		Support string
	}{}

	if err := d.macRPC.Call("driver.Bundle", &out, nil); err != nil {
		app.Panic(err)
	}

	// Set up the support directory in case of the app is not bundled.
	if strings.HasSuffix(out.Support, "{appname}") {
		wd, err := os.Getwd()
		if err != nil {
			app.Panic(errors.Wrap(err, "support unreachable"))
		}

		appname := filepath.Base(wd) + "-" + d.devID
		out.Support = strings.Replace(out.Support, "{appname}", appname, 1)
	}

	return out.Support
}

func (d *Driver) onRun(in map[string]interface{}) interface{} {
	d.configureDefaultWindow()
	d.menubar = newMenuBar(d.MenubarConfig)
	d.docktile = newDockTile(app.MenuConfig{URL: d.DockURL})

	if len(d.URL) != 0 {
		app.NewWindow(d.DefaultWindow)
	}

	d.events.Emit(app.Running, d)
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
	hasVisibleWindow := in["HasVisibleWindows"].(bool)

	if !hasVisibleWindow && len(d.URL) != 0 {
		app.NewWindow(d.DefaultWindow)
	}

	d.events.Emit(app.Reopened, d)
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
			app.Panic(errors.Wrap(err, "onURLOpen"))
		}

		d.OnURLOpen(u)
	}
	return nil
}

func (d *Driver) onFileDrop(in map[string]interface{}) interface{} {
	d.droppedFiles = bridge.Strings(in["Filenames"])
	return nil
}

func (d *Driver) onQuit(in map[string]interface{}) interface{} {
	if d.OnQuit != nil {
		d.OnQuit()
	}

	d.stop()
	return nil
}

func (d *Driver) configureDefaultWindow() {
	if d.DefaultWindow == (app.WindowConfig{}) {
		d.DefaultWindow = app.WindowConfig{
			Title:     d.AppName(),
			MinWidth:  480,
			Width:     1280,
			MinHeight: 480,
			Height:    768,
			URL:       d.URL,
		}
	}

	if len(d.DefaultWindow.URL) == 0 {
		d.DefaultWindow.URL = d.URL
	}
}

func generateDevID() string {
	h := md5.New()
	wd, _ := os.Getwd()
	io.WriteString(h, wd)
	return fmt.Sprintf("%x", h.Sum(nil))
}
