package test

import (
	"context"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
)

// Driver is an app.Driver implementation for testing.
type Driver struct {
	core.Driver

	// A boolean that reports whether driver set element errors.
	Err bool

	ui       chan func()
	factory  *app.Factory
	events   *app.EventRegistry
	elems    *core.ElemDB
	stop     func()
	menubar  *Menu
	docktile *DockTile
}

// Target satisfies the app.Driver interface.
func (d *Driver) Target() string {
	return "web"
}

// Run satisfies the app.Driver interface.
func (d *Driver) Run(c app.DriverConfig) error {
	d.ui = c.UI
	d.factory = c.Factory
	d.events = c.Events
	d.elems = core.NewElemDB()
	d.menubar = newMenu(d, app.MenuConfig{})
	d.docktile = newDockTile(d)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	d.stop = cancel

	d.events.Emit(app.Running)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case fn := <-d.ui:
			fn()
		}
	}
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
	w := newWindow(d, c)
	d.setElemErr(w)
	return w
}

// NewPage satisfies the app.Driver interface.
func (d *Driver) NewPage(c app.PageConfig) app.Page {
	p := newPage(d, c)
	d.setElemErr(p)
	return p
}

// NewContextMenu satisfies the app.Driver interface.
func (d *Driver) NewContextMenu(c app.MenuConfig) app.Menu {
	m := newMenu(d, c)
	d.setElemErr(m)
	return m
}

// NewController satisfies the app.Driver interface.
func (d *Driver) NewController(c app.ControllerConfig) app.Controller {
	m := newController(d, c)
	d.setElemErr(m)
	return m
}

// MenuBar satisfies the app.Driver interface.
func (d *Driver) MenuBar() app.Menu {
	d.setElemErr(d.menubar)
	return d.menubar
}

// NewStatusMenu satisfies the app.Driver interface.
func (d *Driver) NewStatusMenu(c app.StatusMenuConfig) app.StatusMenu {
	m := newStatusMenu(d, c)
	d.setElemErr(m)
	return m
}

// DockTile satisfies the app.Driver interface.
func (d *Driver) DockTile() app.DockTile {
	d.setElemErr(d.docktile)
	return d.docktile
}

// UI satisfies the app.Driver interface.
func (d *Driver) UI(f func()) {
	d.ui <- f
}

// Stop satisfies the app.Driver interface.
func (d *Driver) Stop() {
	if d.stop != nil {
		d.stop()
	}
}

func (d *Driver) setElemErr(e errSetter) {
	if d.Err && e.Err() == nil {
		e.SetErr(app.ErrNotSupported)
	}
}

type errSetter interface {
	app.Elem
	SetErr(error)
}
