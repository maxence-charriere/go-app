// +build darwin,amd64

package mac

import (
	"encoding/json"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
	"github.com/murlokswarm/app/internal/dom"
	"github.com/pkg/errors"
)

func newContextMenu(d *core.Driver) *core.Menu {
	return &core.Menu{
		DOM:    dom.Engine{Resources: d.Resources},
		Driver: d,
	}
}

func newMenuBar(d *core.Driver) *core.Menu {
	return &core.Menu{
		DOM:       dom.Engine{Resources: d.Resources},
		Driver:    d,
		NoDestroy: true,
	}
}

func newDockTile(d *core.Driver) *core.DockTile {
	return &core.DockTile{
		core.Menu{
			DOM:       dom.Engine{Resources: d.Resources},
			Driver:    d,
			NoDestroy: true,
		},
	}
}

func onMenuCallback(m *core.Menu, in map[string]interface{}) {
	mappingStr := in["Mapping"].(string)

	var mapping dom.Mapping
	if err := json.Unmarshal([]byte(mappingStr), &mapping); err != nil {
		app.Logf("menu callback failed: %s", err)
		return
	}

	c, err := m.DOM.CompoByID(mapping.CompoID)
	if err != nil {
		app.Logf("menu callback failed: %s", err)
		return
	}

	var f func()
	if f, err = mapping.Map(c); err != nil {
		app.Logf("menu callback failed: %s", err)
		return
	}

	if f != nil {
		f()
		return
	}

	app.Render(c)
}

func onMenuClose(m *core.Menu, in map[string]interface{}) {
	if m.NoDestroy {
		return
	}

	// menuDidClose: is called before clicked:.
	// We call CallOnUIGoroutine in order to defer the close operation
	// after the clicked one.
	driver.UI(func() {
		if err := driver.Platform.Call("menus.Delete", nil, struct {
			ID string
		}{
			ID: m.ID(),
		}); err != nil {
			app.Panic(errors.Wrap(err, "onMenuClose"))
		}

		driver.Elems.Delete(m)
	})
}
