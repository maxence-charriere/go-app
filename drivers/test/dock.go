package test

import (
	"os"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/dom"
)

// DockTile is a teststatus menu that implements the app.DockTile interface.
type DockTile struct {
	Menu
}

func newDockTile(d *Driver) *DockTile {
	dt := &DockTile{
		Menu{
			driver: d,
			dom:    dom.NewDOM(d.factory),
			id:     uuid.New().String(),
		},
	}

	d.elems.Put(dt)
	return dt
}

// WhenDockTile satisfies the app.DockTile interface.
func (d *DockTile) WhenDockTile(f func(app.DockTile)) {
	f(d)
}

// Type satisfies the app.Menu interface.
func (d *DockTile) Type() string {
	return "dock tile"
}

// SetIcon satisfies the app.DockTile interface.
func (d *DockTile) SetIcon(path string) {
	_, err := os.Stat(path)
	d.SetErr(err)
}

// SetBadge satisfies the app.DockTile interface.
func (d *DockTile) SetBadge(v interface{}) {
	d.SetErr(nil)
	d.driver.setElemErr(d)
}
