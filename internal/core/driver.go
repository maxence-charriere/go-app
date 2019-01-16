package core

import (
	"os"
	"path/filepath"

	"github.com/murlokswarm/app"
	"github.com/pkg/errors"
)

// Driver is a modular implementation of the app.Driver interface that provide
// the common driver logic. It used as a base for platform specific driver
// implementations.
type Driver struct {
	Elems                  *ElemDB
	Events                 *app.EventRegistry
	Factory                *app.Factory
	Go                     *Go
	JSToPlatform           string
	OpenDefaultBrowserFunc func(string) error
	NewContextMenuFunc     func(*Driver) *Menu
	NewMenuBarFunc         func(*Driver) *Menu
	NewWindowFunc          func(*Driver) *Window
	Platform               *Platform
	ResourcesFunc          func() string
	StorageFunc            func() string
	UIChan                 chan func()
}

// AppName satisfies the app.Driver interface.
func (d *Driver) AppName() string {
	wd, _ := os.Getwd()
	return filepath.Base(wd)
}

// DockTile satisfies the app.Driver interface.
func (d *Driver) DockTile() app.DockTile {
	dt := &DockTile{}
	dt.SetErr(app.ErrNotSupported)
	return dt
}

// ElemByCompo satisfies the app.Driver interface.
func (d *Driver) ElemByCompo(c app.Compo) app.Elem {
	return d.Elems.GetByCompo(c)
}

// HandleMenu returns a go RPC handler that handle menu requests.
func (d *Driver) HandleMenu(h func(m *Menu, in map[string]interface{})) GoHandler {
	return func(in map[string]interface{}) {
		e := d.Elems.GetByID(in["ID"].(string))
		if e.Err() != nil {
			return
		}

		switch m := e.(type) {
		case *Menu:
			h(m, in)

		default:
			app.Panic(errors.Errorf("%T is not a supported menu"))
		}
	}
}

// HandleWindow returns a go RPC handler that handle window requests.
func (d *Driver) HandleWindow(h func(w *Window, in map[string]interface{})) GoHandler {
	return func(in map[string]interface{}) {
		e := d.Elems.GetByID(in["ID"].(string))
		if e.Err() == app.ErrElemNotSet {
			return
		}

		h(e.(*Window), in)
	}
}

// MenuBar satisfies the app.Driver interface.
func (d *Driver) MenuBar() app.Menu {
	return &Menu{Elem: Elem{err: app.ErrNotSupported}}
}

// NewContextMenu satisfies the app.Driver interface.
func (d *Driver) NewContextMenu(c app.MenuConfig) app.Menu {
	if d.NewContextMenu == nil {
		return &Menu{Elem: Elem{err: app.ErrNotSupported}}
	}

	m := d.NewContextMenuFunc(d)
	m.kind = "context menu"
	m.Create(c)

	if m.err != nil {
		return m
	}

	m.err = d.Platform.Call("driver.SetContextMenu", nil, struct {
		ID string
	}{
		ID: m.id,
	})

	return m
}

// NewMenuBar creates a menu bar.
func (d *Driver) NewMenuBar(c app.MenuBarConfig) *Menu {
	if d.NewMenuBarFunc == nil {
		return &Menu{Elem: Elem{err: app.ErrNotSupported}}
	}

	m := d.NewMenuBarFunc(d)
	m.kind = "menu bar"
	m.Create(app.MenuConfig{URL: menuBarConfigToAddr(c)})

	if m.err != nil {
		return m
	}

	m.err = d.Platform.Call("driver.SetMenubar", nil, struct {
		ID string
	}{
		ID: m.id,
	})

	return m
}

// NewController satisfies the app.Driver interface.
func (d *Driver) NewController(c app.ControllerConfig) app.Controller {
	controller := &Controller{}
	controller.SetErr(app.ErrNotSupported)
	return controller
}

// NewFilePanel satisfies the app.Driver interface.
func (d *Driver) NewFilePanel(c app.FilePanelConfig) app.Elem {
	return &Elem{err: app.ErrNotSupported}
}

// NewNotification satisfies the app.Driver interface.
func (d *Driver) NewNotification(c app.NotificationConfig) app.Elem {
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

// NewStatusMenu satisfies the app.Driver interface.
func (d *Driver) NewStatusMenu(c app.StatusMenuConfig) app.StatusMenu {
	s := &StatusMenu{}
	s.SetErr(app.ErrNotSupported)
	return s
}

// NewWindow satisfies the app.Driver interface.
func (d *Driver) NewWindow(c app.WindowConfig) app.Window {
	if d.NewWindowFunc == nil {
		return &Window{Elem: Elem{err: app.ErrNotSupported}}
	}

	w := d.NewWindowFunc(d)
	w.Create(c)
	return w
}

// OpenDefaultBrowser satisfies the app.Driver interface.
func (d *Driver) OpenDefaultBrowser(url string) error {
	if d.OpenDefaultBrowserFunc == nil {
		return app.ErrNotSupported
	}

	return d.OpenDefaultBrowserFunc(url)
}

// Render satisfies the app.Driver interface.
func (d *Driver) Render(c app.Compo) {
	e := d.ElemByCompo(c)

	e.WhenView(func(v app.View) {
		v.Render(c)
	})
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

// Run satisfies the app.Driver interface.
func (d *Driver) Run(c app.DriverConfig) error {
	return app.ErrNotSupported
}

// Stop satisfies the app.Driver interface.
func (d *Driver) Stop() {
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

// Target satisfies the app.Driver interface.
func (d *Driver) Target() string {
	return "test"
}

// UI satisfies the app.Driver interface.
func (d *Driver) UI(f func()) {
	d.UIChan <- f
}
