// +build darwin,amd64

package mac

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/bridge"
	"github.com/murlokswarm/app/internal/core"
	"github.com/murlokswarm/app/internal/html"
	"github.com/pkg/errors"
)

// Menu implements the app.Menu interface.
type Menu struct {
	core.Menu

	markup         app.Markup
	id             string
	typ            string
	compo          app.Compo
	keepWhenClosed bool

	onClose func()
}

func newMenu(c app.MenuConfig, typ string) *Menu {
	m := &Menu{
		markup: app.ConcurrentMarkup(html.NewMarkup(driver.factory)),
		id:     uuid.New().String(),
		typ:    typ,

		onClose: c.OnClose,
	}

	if err := driver.macRPC.Call("menus.New", nil, struct {
		ID string
	}{
		ID: m.id,
	}); err != nil {
		m.SetErr(err)
		return m
	}

	driver.elems.Put(m)

	if len(c.URL) != 0 {
		m.Load(c.URL)
	}

	return m
}

// ID satisfies the app.Menu interface.
func (m *Menu) ID() string {
	return m.id
}

// Load satisfies the app.Menu interface.
func (m *Menu) Load(urlFmt string, v ...interface{}) {
	var err error
	defer func() {
		m.SetErr(err)
	}()

	if m.compo != nil {
		m.markup.Dismount(m.compo)
		m.compo = nil
	}

	u := fmt.Sprintf(urlFmt, v...)
	n := core.CompoNameFromURLString(u)

	var c app.Compo
	if c, err = driver.factory.NewCompo(n); err != nil {
		return
	}

	if _, err = m.markup.Mount(c); err != nil {
		return
	}

	if nav, ok := c.(app.Navigable); ok {
		navURL, _ := url.Parse(u)
		nav.OnNavigate(navURL)
	}

	var root app.Tag
	if root, err = m.markup.Root(c); err != nil {
		return
	}

	if root, err = m.markup.FullRoot(root); err != nil {
		return
	}

	m.compo = c

	err = driver.macRPC.Call("menus.Load", nil, struct {
		ID  string
		Tag app.Tag
	}{
		ID:  m.id,
		Tag: root,
	})
}

// Compo satisfies the app.Menu interface.
func (m *Menu) Compo() app.Compo {
	return m.compo
}

// Contains satisfies the app.Menu interface.
func (m *Menu) Contains(c app.Compo) bool {
	return m.markup.Contains(c)
}

// Render satisfies the app.Menu interface.
func (m *Menu) Render(c app.Compo) {
	var err error
	defer func() {
		m.SetErr(err)
	}()

	var syncs []app.TagSync
	if syncs, err = m.markup.Update(c); err != nil {
		return
	}

	for _, s := range syncs {
		if s.Replace {
			err = m.render(s)
		} else {
			err = m.renderAttributes(c, s)
		}

		if err != nil {
			return
		}
	}
}

// Type satisfies the app.Menu interface.
func (m *Menu) Type() string {
	return m.typ
}

func (m *Menu) render(sync app.TagSync) error {
	tag, err := m.markup.FullRoot(sync.Tag)
	if err != nil {
		return err
	}

	return driver.macRPC.Call("menus.Render", nil, struct {
		ID  string
		Tag app.Tag
	}{
		ID:  m.id,
		Tag: tag,
	})
}

func (m *Menu) renderAttributes(c app.Compo, sync app.TagSync) error {
	root, err := m.markup.Root(c)
	if err != nil {
		return err
	}

	tag := sync.Tag
	if root.ID != tag.ID {
		// Ensure that objc will not do extra initializations.
		tag.Children = nil
	}

	return driver.macRPC.Call("menus.RenderAttributes", nil, struct {
		ID  string
		Tag app.Tag
	}{
		ID:  m.id,
		Tag: tag,
	})
}

func onMenuCallback(m *Menu, in map[string]interface{}) interface{} {
	mappingString := in["Mapping"].(string)

	var mapping app.Mapping
	if err := json.Unmarshal([]byte(mappingString), &mapping); err != nil {
		app.Log("menu callback failed: %s", err)
		return nil
	}

	function, err := m.markup.Map(mapping)
	if err != nil {
		app.Log("menu callback failed: %s", err)
		return nil
	}

	if function != nil {
		function()
		return nil
	}

	var c app.Compo
	if c, err = m.markup.Compo(mapping.CompoID); err != nil {
		app.Log("menu callback failed: %s", err)
		return nil
	}

	m.Render(c)
	if m.Err() != nil {
		app.Log("menu callback failed: %s", err)
	}

	return nil
}

func onMenuClose(m *Menu, in map[string]interface{}) interface{} {
	if m.keepWhenClosed {
		return nil
	}

	// menuDidClose: is called before clicked:.
	// We call CallOnUIGoroutine in order to defer the close operation
	// after the clicked one.
	driver.CallOnUIGoroutine(func() {
		if m.onClose != nil {
			m.onClose()
		}

		if err := driver.macRPC.Call("menus.Delete", nil, struct {
			ID string
		}{
			ID: m.id,
		}); err != nil {
			panic(errors.Wrap(err, "onMenuClose"))
		}

		driver.elems.Delete(m)
	})

	return nil
}

func handleMenu(h func(m *Menu, in map[string]interface{}) interface{}) bridge.GoRPCHandler {
	return func(in map[string]interface{}) interface{} {
		id, _ := in["ID"].(string)
		e := driver.elems.GetByID(id)

		switch menu := e.(type) {
		case *Menu:
			return h(menu, in)

		// case *DockTile:
		// 	return h(&menu.Menu, in)

		// case *StatusMenu:
		// 	return h(&menu.Menu, in)

		default:
			panic("menu not supported")
		}
	}
}
