// +build darwin,amd64

package mac

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/bridge"
	"github.com/murlokswarm/app/html"
)

// DockTile implements the app.DockTile interface.
type DockTile struct {
	menu Menu
}

func newDockTile(config app.MenuConfig) (app.DockTile, error) {
	var markup app.Markup = html.NewMarkup(driver.factory)
	markup = app.NewConcurrentMarkup(markup)

	rawDock := &DockTile{
		menu: Menu{
			id:        uuid.New(),
			markup:    markup,
			lastFocus: time.Now(),
		},
	}

	dock := app.NewDockTileWithLogs(rawDock)

	if _, err := driver.macos.Request(
		fmt.Sprintf("/menu/new?id=%s", rawDock.menu.id),
		nil,
	); err != nil {
		return nil, err
	}

	if err := driver.elements.Add(dock); err != nil {
		return nil, err
	}

	if len(config.DefaultURL) != 0 {
		return dock, dock.Load(config.DefaultURL)
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
	if err := d.menu.Load(url, v...); err != nil {
		return err
	}

	_, err := driver.macos.Request(
		fmt.Sprintf("/driver/dock/set?menu-id=%v", d.ID()),
		nil,
	)
	return err
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

	icon := struct {
		Path string `json:"path"`
	}{
		Path: name,
	}

	_, err := driver.macos.RequestWithAsyncResponse(
		"/driver/dock/icon",
		bridge.NewPayload(icon),
	)
	return err
}

// SetBadge satisfies the app.DockTile interface.
func (d *DockTile) SetBadge(v interface{}) error {
	badge := struct {
		Message string `json:"message"`
	}{}

	if v == nil {
		badge.Message = ""
	} else {
		badge.Message = fmt.Sprint(v)
	}

	_, err := driver.macos.Request(
		"/driver/dock/badge",
		bridge.NewPayload(badge),
	)
	return err
}
