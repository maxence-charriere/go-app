package test

import (
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app/markup"
)

// A DockTile implementation for tests.
type DockTile struct {
	Menu
}

func newDockTile(d *Driver) *DockTile {
	dock := &DockTile{
		Menu: Menu{
			id:           uuid.New(),
			compoBuilder: d.compoBuilder,
			env:          markup.NewEnv(d.compoBuilder),
			lastFocus:    time.Now(),
		},
	}
	d.elements.Add(dock)
	return dock
}

// SetIcon satisfies the app.DockTile interface.
func (d *DockTile) SetIcon(name string) error {
	return nil
}

// SetBadge satisfies the app.DockTile interface.
func (d *DockTile) SetBadge(v interface{}) {
}
