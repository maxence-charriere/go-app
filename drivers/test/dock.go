package test

import (
	"encoding/json"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/html"
)

// A DockTile implementation for tests.
type DockTile struct {
	Menu
}

func newDockTile(d *Driver) app.DockTile {
	var markup app.Markup = html.NewMarkup(d.factory)
	markup = app.ConcurrentMarkup(markup)

	dock := &DockTile{
		Menu: Menu{
			id:        uuid.New(),
			factory:   d.factory,
			markup:    markup,
			lastFocus: time.Now(),
		},
	}

	d.elements.Add(dock)
	return dock
}

// SetIcon satisfies the app.DockTile interface.
func (d *DockTile) SetIcon(name string) error {
	_, err := os.Stat(name)
	return err
}

// SetBadge satisfies the app.DockTile interface.
func (d *DockTile) SetBadge(v interface{}) error {
	_, err := json.Marshal(v)
	return err
}
