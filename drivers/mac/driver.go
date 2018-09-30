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
	"github.com/murlokswarm/app/internal/bridge"
	"github.com/murlokswarm/app/internal/core"
	"github.com/murlokswarm/app/internal/logs"
	"github.com/pkg/errors"
)

var (
	driver       *Driver
	goappBundle  = os.Getenv("GOAPP_BUNDLE")
	debug        = os.Getenv("GOAPP_DEBUG") == "true"
	goappLogAddr = os.Getenv("GOAPP_LOGS_ADDR")
	goappLogs    *logs.GoappClient
)

func init() {
	if len(goappBundle) != 0 {
		app.Logger = func(format string, a ...interface{}) {}
		return
	}

	if len(goappLogAddr) != 0 {
		app.EnableDebug(debug)
		goappLogs = logs.NewGoappClient(goappLogAddr, logs.WithColoredPrompt)
		app.Logger = goappLogs.Logger()
		return
	}

	logger := logs.ToWriter(os.Stderr)
	app.Logger = logs.WithColoredPrompt(logger)
}

// Driver is the app.Driver implementation for MacOS.
type Driver struct {
	core.Driver

	// The bundle configuration.
	// Only applied when the app is build with goapp mac build -bundle.
	Bundle Bundle

	// Menubar configuration
	MenubarConfig MenuBarConfig

	// The URL of the component to load in the main window.
	// The main window is not created when OnRun or OnReopen are set.
	URL string

	// The URL of the component to load in the dock.
	DockURL string

	// The func called right after app.Run.
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

	factory      *app.Factory
	elems        *core.ElemDB
	devID        string
	macRPC       bridge.PlatformRPC
	goRPC        bridge.GoRPC
	uichan       chan func()
	stopchan     chan error
	menubar      *Menu
	docktile     *DockTile
	droppedFiles []string
}

// Run satisfies the app.Driver interface.
func (d *Driver) Run(f *app.Factory) error {
	if len(goappBundle) != 0 {
		return d.runGoappBundle()
	}

	if driver != nil {
		return errors.New("running already")
	}

	d.factory = f
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
	d.goRPC.Handle("windows.OnAlert", handleWindow(onWindowAlert))

	d.goRPC.Handle("menus.OnClose", handleMenu(onMenuClose))
	d.goRPC.Handle("menus.OnCallback", handleMenu(onMenuCallback))

	d.goRPC.Handle("filePanels.OnSelect", handleFilePanel(onFilePanelSelect))
	d.goRPC.Handle("saveFilePanels.OnSelect", handleSaveFilePanel(onSaveFilePanelSelect))

	d.goRPC.Handle("notifications.OnReply", handleNotification(onNotificationReply))

	d.uichan = make(chan func(), 256)
	d.stopchan = make(chan error)
	defer close(d.uichan)
	defer close(d.stopchan)

	driver = d

	if goappLogs != nil {
		go goappLogs.WaitForStop(d.Stop)
	}

	go func() {
		d.stopchan <- d.macRPC.Call("driver.Run", nil, nil)
	}()

	for {
		select {
		case err := <-d.stopchan:
			return err

		case fn := <-d.uichan:
			fn()
		}
	}
}

// AppName satisfies the app.Driver interface.
func (d *Driver) AppName() string {
	var out struct {
		AppName string
	}

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
	var out struct {
		Resources string
	}

	if err := d.macRPC.Call("driver.Bundle", &out, nil); err != nil {
		app.Panic(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		app.Panic(errors.Wrap(err, "resources unreachable"))
	}

	if filepath.Dir(out.Resources) == wd {
		out.Resources = filepath.Join(wd, "resources")
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
	if e.Err() == nil {
		e.(app.ElemWithCompo).Render(c)
	}
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
	d.uichan <- f
}

// Stop satisfies the app.Driver interface.
func (d *Driver) Stop() {
	if err := d.macRPC.Call("driver.Quit", nil, nil); err != nil {
		d.stopchan <- err
	}
}

func (d *Driver) runGoappBundle() error {
	b, err := json.MarshalIndent(d.Bundle, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(goappBundle, b, 0777)
}

func (d *Driver) support() string {
	var out struct {
		Support string
	}

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
	d.menubar = newMenuBar(d.MenubarConfig)
	d.docktile = newDockTile(app.MenuConfig{URL: d.DockURL})

	if d.OnRun == nil {
		d.OnRun = d.newMainWindow
	}

	d.OnRun()
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
	if d.OnReopen == nil {
		d.OnReopen = func(hasVisibleWindows bool) {
			if !hasVisibleWindows {
				d.newMainWindow()
			}
		}
	}

	d.OnReopen(in["HasVisibleWindows"].(bool))
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

	if goappLogs != nil {
		goappLogs.Close()
	}

	return nil
}

func (d *Driver) newMainWindow() {
	app.NewWindow(app.WindowConfig{
		Title:          d.AppName(),
		TitlebarHidden: true,
		MinWidth:       480,
		Width:          1280,
		MinHeight:      480,
		Height:         768,
		URL:            d.URL,
	})
}

func generateDevID() string {
	h := md5.New()
	wd, _ := os.Getwd()
	io.WriteString(h, wd)
	return fmt.Sprintf("%x", h.Sum(nil))
}
