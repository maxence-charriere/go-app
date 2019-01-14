package core

import (
	"os"
	"path/filepath"

	"github.com/murlokswarm/app"
)

// Driver is a base struct to embed in app.Driver implementations.
type Driver struct {
	// The database that contain the created elements.
	Elems *ElemDB

	// The event registry that emit events.
	Events *app.EventRegistry

	// The factory used to import and create components.
	Factory *app.Factory

	// The RPC object to deliver procedure calls to Go.
	Go *Go

	// The function name that javascript use to send data to the targetted
	// platform.
	JSToPlatform string

	// The function to open an URL in the targetted platform default browser.
	OpenDefaultBrowser func(string) error

	WindowFactory func(*Driver) *Window

	// The RPC object to call targetted platform procedures.
	Platform *Platform

	// A function that returns the targetted platform resources directory
	// path.
	ResourcesPath func() string

	// A function that returns the targetted platform storage directory path.
	StoragePath func() string

	// The channel used to execute function on the UI goroutine.
	UIChan chan func()
}

// Target satisfies the app.Driver interface.
func (d *Driver) Target() string {
	return "test"
}

// Run satisfies the app.Driver interface.
func (d *Driver) Run(c app.DriverConfig) error {
	return app.ErrNotSupported
}

// AppName satisfies the app.Driver interface.
func (d *Driver) AppName() string {
	wd, _ := os.Getwd()
	return filepath.Base(wd)
}

// Resources satisfies the app.Driver interface.
func (d *Driver) Resources(p ...string) string {
	r := filepath.Join(p...)
	r = filepath.Join(d.ResourcesPath(), r)
	return r
}

// Storage satisfies the app.Driver interface.
func (d *Driver) Storage(p ...string) string {
	s := filepath.Join(p...)
	s = filepath.Join(d.StoragePath(), s)
	return s
}

// Render satisfies the app.Driver interface.
func (d *Driver) Render(c app.Compo) {
	e := d.ElemByCompo(c)

	e.WhenView(func(v app.View) {
		v.Render(c)
	})
}

// ElemByCompo satisfies the app.Driver interface.
func (d *Driver) ElemByCompo(c app.Compo) app.Elem {
	return d.Elems.GetByCompo(c)
}

// NewWindow satisfies the app.Driver interface.
func (d *Driver) NewWindow(c app.WindowConfig) app.Window {
	if d.WindowFactory == nil {
		w := &Window{}
		w.err = app.ErrNotSupported
		return w
	}

	w := d.WindowFactory(d)
	w.Create(c)
	return w
}

// NewContextMenu satisfies the app.Driver interface.
func (d *Driver) NewContextMenu(c app.MenuConfig) app.Menu {
	m := &Menu{}
	m.SetErr(app.ErrNotSupported)
	return m
}

// NewStatusMenu satisfies the app.Driver interface.
func (d *Driver) NewStatusMenu(c app.StatusMenuConfig) app.StatusMenu {
	s := &StatusMenu{}
	s.SetErr(app.ErrNotSupported)
	return s
}

// NewFilePanel satisfies the app.Driver interface.
func (d *Driver) NewFilePanel(c app.FilePanelConfig) app.Elem {
	return &Elem{err: app.ErrNotSupported}
}

// NewSaveFilePanel satisfies the app.Driver interface.
func (d *Driver) NewSaveFilePanel(c app.SaveFilePanelConfig) app.Elem {
	return &Elem{err: app.ErrNotSupported}
}

// NewShare satisfies the app.Driver interface.
func (d *Driver) NewShare(v interface{}) app.Elem {
	return &Elem{err: app.ErrNotSupported}
}

// NewNotification satisfies the app.Driver interface.
func (d *Driver) NewNotification(c app.NotificationConfig) app.Elem {
	return &Elem{err: app.ErrNotSupported}
}

// NewController satisfies the app.Driver interface.
func (d *Driver) NewController(c app.ControllerConfig) app.Controller {
	controller := &Controller{}
	controller.SetErr(app.ErrNotSupported)
	return controller
}

// MenuBar satisfies the app.Driver interface.
func (d *Driver) MenuBar() app.Menu {
	m := &Menu{}
	m.SetErr(app.ErrNotSupported)
	return m
}

// DockTile satisfies the app.Driver interface.
func (d *Driver) DockTile() app.DockTile {
	dt := &DockTile{}
	dt.SetErr(app.ErrNotSupported)
	return dt
}

// UI satisfies the app.Driver interface.
func (d *Driver) UI(f func()) {
	d.UIChan <- f
}

// Stop satisfies the app.Driver interface.
func (d *Driver) Stop() {
}
