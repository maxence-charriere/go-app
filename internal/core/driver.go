package core

import (
	"os"
	"path/filepath"

	"github.com/murlokswarm/app"
)

// Driver is a modular implementation of the app.Driver interface that provide
// the common driver logic. It used as a base for platform specific driver
// implementations.
type Driver struct {
	Elems              *ElemDB
	Events             *app.EventRegistry
	Factory            *app.Factory
	Go                 *Go
	JSToPlatform       string
	OpenDefaultBrowser func(string) error
	NewWindowFunc      func(*Driver) *Window
	Platform           *Platform
	ResourcesFunc      func() string
	StorageFunc        func() string
	UIChan             chan func()
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
	if d.ResourcesFunc == nil {
		d.StorageFunc = func() string { return "resources" }
	}

	r := filepath.Join(p...)
	r = filepath.Join(d.ResourcesFunc(), r)
	return r
}

// Storage satisfies the app.Driver interface.
func (d *Driver) Storage(p ...string) string {
	if d.StorageFunc == nil {
		d.StorageFunc = func() string { return "storage" }
	}

	s := filepath.Join(p...)
	s = filepath.Join(d.StorageFunc(), s)
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
	if d.NewWindowFunc == nil {
		w := &Window{}
		w.err = app.ErrNotSupported
		return w
	}

	w := d.NewWindowFunc(d)
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
