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

	// The function executed after a Run call.
	OnRun func()

	factory  *app.Factory
	elems    *core.ElemDB
	stop     func()
	uichan   chan func()
	menubar  *Menu
	docktile *DockTile
}

// Run satisfies the app.Driver interface.
func (d *Driver) Run(f *app.Factory) error {
	d.factory = f
	d.elems = core.NewElemDB()
	d.uichan = make(chan func(), 64)
	d.menubar = newMenu(d, app.MenuConfig{})
	d.docktile = newDockTile(d)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	d.stop = cancel

	if d.OnRun != nil {
		d.CallOnUIGoroutine(d.OnRun)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case fn := <-d.uichan:
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

// CallOnUIGoroutine satisfies the app.Driver interface.
func (d *Driver) CallOnUIGoroutine(f func()) {
	d.uichan <- f
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
