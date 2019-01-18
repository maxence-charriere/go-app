package core

import (
	"fmt"
	"os"

	"github.com/murlokswarm/app"
)

// DockTile is a modular implementation of the app.DockTile interface that can
// be configured to address the different drivers needs.
type DockTile struct {
	Menu
}

// WhenDockTile satisfies the app.DockTile interface.
func (d *DockTile) WhenDockTile(f func(app.DockTile)) {
	f(d)
}

// SetIcon satisfies the app.DockTile interface.
func (d *DockTile) SetIcon(path string) {
	if _, d.err = os.Stat(path); len(path) != 0 && d.err != nil {
		return
	}

	d.err = d.Driver.Platform.Call("docks.SetIcon", nil, struct {
		Icon string
	}{
		Icon: path,
	})
}

// SetBadge satisfies the app.DockTile interface.
func (d *DockTile) SetBadge(v interface{}) {
	badge := ""
	if v != nil {
		badge = fmt.Sprint(v)
	}

	d.err = d.Driver.Platform.Call("docks.SetBadge", nil, struct {
		Badge string
	}{
		Badge: badge,
	})
}
