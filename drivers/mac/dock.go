// +build darwin,amd64

package mac

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/dom"
)

// DockTile implements the app.DockTile interface.
type DockTile struct {
	Menu
}

func newDockTile(c app.MenuConfig) *DockTile {
	d := &DockTile{
		Menu: Menu{
			dom:            dom.NewDOM(driver.factory),
			id:             uuid.New().String(),
			typ:            "dock tile",
			keepWhenClosed: true,
		},
	}

	if err := driver.macRPC.Call("menus.New", nil, struct {
		ID string
	}{
		ID: d.id,
	}); err != nil {
		d.SetErr(err)
		return d
	}

	driver.elems.Put(d)

	if len(c.URL) != 0 {
		d.Load(c.URL)
	}

	return d
}

// WhenDockTile satisfies the app.DockTile interface.
func (d *DockTile) WhenDockTile(f func(app.DockTile)) {
	f(d)
}

// Load the app.StatusMenu interface.
func (d *DockTile) Load(urlFmt string, v ...interface{}) {
	d.Menu.Load(urlFmt, v...)
	if d.Err() != nil {
		return
	}

	err := driver.macRPC.Call("docks.SetMenu", nil, struct {
		ID string
	}{
		ID: d.id,
	})

	d.SetErr(err)
}

// SetIcon satisfies the app.DockTile interface.
func (d *DockTile) SetIcon(path string) {
	if _, err := os.Stat(path); err != nil && len(path) != 0 {
		d.SetErr(err)
		return
	}

	err := driver.macRPC.Call("docks.SetIcon", nil, struct {
		Icon string
	}{
		Icon: path,
	})

	d.SetErr(err)
}

// SetBadge satisfies the app.DockTile interface.
func (d *DockTile) SetBadge(v interface{}) {
	var badge string
	if v != nil {
		badge = fmt.Sprint(v)
	}

	err := driver.macRPC.Call("docks.SetBadge", nil, struct {
		Badge string
	}{
		Badge: badge,
	})

	d.SetErr(err)
}
