// +build darwin,amd64

package mac

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/html"
)

// DockTile implements the app.DockTile interface.
type DockTile struct {
	menu Menu
}

func newDockTile(c app.MenuConfig) (app.DockTile, error) {
	var markup app.Markup = html.NewMarkup(driver.factory)
	markup = app.NewConcurrentMarkup(markup)

	rawDock := &DockTile{
		menu: Menu{
			id:        uuid.New(),
			markup:    markup,
			lastFocus: time.Now(),
		},
	}

	dock := app.DockTileWithLogs(rawDock)

	if err := driver.macRPC.Call("menus.New", nil, struct {
		ID string
	}{
		ID: dock.ID().String(),
	}); err != nil {
		return nil, err
	}

	if err := driver.elements.Add(dock); err != nil {
		return nil, err
	}

	if len(c.DefaultURL) != 0 {
		if err := dock.Load(c.DefaultURL); err != nil {
			return nil, err
		}
	}

	if err := driver.macRPC.Call("driver.SetDock", nil, dock.ID()); err != nil {
		return nil, err
	}

	return dock, nil
}

// ID satisfies the app.DockTile interface.
func (d *DockTile) ID() uuid.UUID {
	return d.menu.ID()
}

// Base satisfies the app.DockTile interface.
func (d *DockTile) Base() app.Menu {
	return d.menu.Base()
}

// Load satisfies the app.DockTile interface.
func (d *DockTile) Load(url string, v ...interface{}) error {
	return d.menu.Load(url, v...)
}

// Contains satisfies the app.DockTile interface.
func (d *DockTile) Contains(compo app.Component) bool {
	return d.menu.Contains(compo)
}

// Component satisfies the app.DockTile interface.
func (d *DockTile) Component() app.Component {
	return d.menu.component
}

// Render satisfies the app.DockTile interface.
func (d *DockTile) Render(compo app.Component) error {
	return d.menu.Render(compo)
}

// LastFocus satisfies the app.DockTile interface.
func (d *DockTile) LastFocus() time.Time {
	return d.menu.LastFocus()
}

// SetIcon satisfies the app.DockTile interface.
func (d *DockTile) SetIcon(name string) error {
	if _, err := os.Stat(name); err != nil && len(name) != 0 {
		return err
	}

	return driver.macRPC.Call("driver.SetDockIcon", nil, name)
}

// SetBadge satisfies the app.DockTile interface.
func (d *DockTile) SetBadge(v interface{}) error {
	var badge string
	if v != nil {
		badge = fmt.Sprint(v)
	}

	return driver.macRPC.Call("driver.SetDockBadge", nil, badge)
}
