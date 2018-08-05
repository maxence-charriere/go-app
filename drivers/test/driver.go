package test

import (
	"context"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
)

// Driver is an app.Driver implementation for testing.
type Driver struct {
	core.Driver

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
	return newWindow(d, c)
}

// NewPage satisfies the app.Driver interface.
func (d *Driver) NewPage(c app.PageConfig) app.Page {
	return newPage(d, c)
}

// NewContextMenu satisfies the app.Driver interface.
func (d *Driver) NewContextMenu(c app.MenuConfig) app.Menu {
	return newMenu(d, c)
}

// MenuBar satisfies the app.Driver interface.
func (d *Driver) MenuBar() app.Menu {
	return d.menubar
}

// NewStatusMenu satisfies the app.Driver interface.
func (d *Driver) NewStatusMenu(c app.StatusMenuConfig) app.StatusMenu {
	return newStatusMenu(d, c)
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
	if d.stop != nil {
		d.stop()
	}
}
