// +build darwin,amd64

// Package mac is the driver to be used for applications that will run on MacOS.
package mac

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Cocoa
#cgo LDFLAGS: -framework WebKit
#cgo LDFLAGS: -framework CoreImage
#cgo LDFLAGS: -framework Security
#include "driver.h"
*/
import "C"
import (
	"net/url"
	"runtime"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/bridge"
	"github.com/murlokswarm/app/markup"
	"github.com/pkg/errors"
)

var (
	driver *Driver
)

func init() {
	runtime.LockOSThread()
}

// Driver is the app.Driver implementation for MacOS.
type Driver struct {
	MenubarURL string
	DockURL    string

	OnRun       func()
	OnFocus     func()
	OnBlur      func()
	OnReopen    func(hasVisibleWindows bool)
	OnFilesOpen func(filenames []string)
	OnURLOpen   func(u *url.URL)
	OnQuit      func() bool
	OnExit      func()

	components markup.CompoBuilder
	elements   app.ElementStore
	uichan     chan func()
	macos      bridge.PlatformBridge
	golang     bridge.GoBridge
	menubar    app.Menu
	dock       app.DockTile
}

// Run satisfies the app.Driver interface.
func (d *Driver) Run(b markup.CompoBuilder) error {
	d.components = b
	d.elements = app.NewElementStore()

	d.uichan = make(chan func(), 256)
	defer close(d.uichan)

	d.macos = bridge.NewPlatformBridge(handleMacOSRequest)
	d.golang = bridge.NewGoBridge(d.uichan)

	d.golang.Handle("/driver/run", d.onRun)
	d.golang.Handle("/driver/focus", d.onFocus)
	d.golang.Handle("/driver/blur", d.onBlur)
	d.golang.Handle("/driver/reopen", d.onReopen)
	d.golang.Handle("/driver/filesopen", d.onFilesOpen)
	d.golang.Handle("/driver/urlopen", d.onURLOpen)
	d.golang.Handle("/driver/quit", d.onQuit)
	d.golang.Handle("/driver/exit", d.onExit)

	go func() {
		for f := range d.uichan {
			f()
		}
	}()

	driver = d
	C.driver_run()
	return nil
}

func (d *Driver) onRun(u *url.URL, p bridge.Payload) (res bridge.Payload) {
	if d.OnRun != nil {
		return
	}

	d.OnRun()
	return
}

func (d *Driver) onFocus(u *url.URL, p bridge.Payload) (res bridge.Payload) {
	if d.OnFocus == nil {
		return
	}

	d.OnFocus()
	return
}

func (d *Driver) onBlur(u *url.URL, p bridge.Payload) (res bridge.Payload) {
	if d.OnBlur == nil {
		return
	}

	d.OnBlur()
	return
}

func (d *Driver) onReopen(u *url.URL, p bridge.Payload) (res bridge.Payload) {
	if d.OnReopen == nil {
		return
	}

	var hasVisibleWindows bool
	p.Unmarshal(&hasVisibleWindows)
	d.OnReopen(hasVisibleWindows)
	return
}

func (d *Driver) onFilesOpen(u *url.URL, p bridge.Payload) (res bridge.Payload) {
	if d.OnFilesOpen == nil {
		return
	}

	var filenames []string
	p.Unmarshal(&filenames)
	d.OnFilesOpen(filenames)
	return
}

func (d *Driver) onURLOpen(u *url.URL, p bridge.Payload) (res bridge.Payload) {
	if d.OnURLOpen == nil {
		return
	}

	var rawurl string
	p.Unmarshal(&rawurl)

	purl, err := url.Parse(rawurl)
	if err != nil {
		panic(errors.Wrap(err, "parsing url failed"))
	}

	d.OnURLOpen(purl)
	return
}

func (d *Driver) onQuit(u *url.URL, p bridge.Payload) (res bridge.Payload) {
	if d.OnQuit == nil {
		return
	}

	res = bridge.NewPayload(d.OnQuit())
	return
}

func (d *Driver) onExit(u *url.URL, p bridge.Payload) (res bridge.Payload) {
	if d.OnExit == nil {
		return
	}
	d.OnExit()
	return
}

// Render satisfies the app.Driver interface.
func (d *Driver) Render(c markup.Component) error {
	elem, err := d.elements.ElementByComponent(c)
	if err != nil {
		return err
	}
	return elem.Render(c)
}

// Context satisfies the app.Driver interface.
func (d *Driver) Context(c markup.Component) (elem app.ElementWithComponent, err error) {
	return d.elements.ElementByComponent(c)
}

// NewContextMenu satisfies the app.Driver interface.
func (d *Driver) NewContextMenu(c app.MenuConfig) app.Menu {
	panic("not implemented")
}

// Resources satisfies the app.Driver interface.
func (d *Driver) Resources() string {
	panic("not implemented")
}

// CallOnUIGoroutine satisfies the app.Driver interface.
func (d *Driver) CallOnUIGoroutine(f func()) {
	d.uichan <- f
}

// Storage satisfies the app.DriverWithStorage interface.
func (d *Driver) Storage() string {
	panic("not implemented")
}

// NewWindow satisfies the app.DriverWithWindows interface.
func (d *Driver) NewWindow(c app.WindowConfig) app.Window {
	panic("not implemented")
}

// MenuBar satisfies the app.DriverWithMenuBar interface.
func (d *Driver) MenuBar() app.Menu {
	panic("not implemented")
}

// Dock satisfies the app.DriverWithDock interface.
func (d *Driver) Dock() app.DockTile {
	panic("not implemented")
}

// Share satisfies the app.DriverWithShare interface.
func (d *Driver) Share(v interface{}) {
	panic("not implemented")
}

// NewFilePanel satisfies the app.DriverWithFilePanels interface.
func (d *Driver) NewFilePanel(c app.FilePanelConfig) app.Element {
	panic("not implemented")
}

// NewPopupNotification satisfies the app.DriverWithPopupNotifications
// interface.
func (d *Driver) NewPopupNotification(c app.PopupNotificationConfig) app.Element {
	panic("not implemented")
}
