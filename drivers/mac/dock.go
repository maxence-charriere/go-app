package mac

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/bridge"
	"github.com/pkg/errors"
)

// DockTile implements the app.DockTile interface.
type DockTile struct {
	menu *Menu
}

func newDockTile(config app.MenuConfig) (dock *DockTile, err error) {
	dock = &DockTile{}

	if dock.menu, err = newMenu(app.MenuConfig{}); err != nil {
		err = errors.Wrap(err, "creating the dock failed")
	}

	if len(config.DefaultURL) != 0 {
		err = dock.Load(config.DefaultURL)
	}
	return
}

// ID satisfies the app.Element interface.
func (d *DockTile) ID() uuid.UUID {
	return d.menu.ID()
}

// Load satisfies the app.Menu interface.
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

// Contains satisfies the app.Menu interface.
func (d *DockTile) Contains(compo app.Component) bool {
	return d.menu.Contains(compo)
}

// Render satisfies the app.Menu interface.
func (d *DockTile) Render(compo app.Component) error {
	return d.menu.Render(compo)
}

// LastFocus satisfies the app.Menu interface.
func (d *DockTile) LastFocus() time.Time {
	return d.menu.LastFocus()
}

// SetIcon satisfies the app.DockTile interface.
func (d *DockTile) SetIcon(name string) error {
	icon := struct {
		Path string `json:"path"`
	}{
		Path: name,
	}

	_, err := driver.macos.Request(
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
