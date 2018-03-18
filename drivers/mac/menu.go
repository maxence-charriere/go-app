// +build darwin,amd64

package mac

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/bridge"
	"github.com/murlokswarm/app/html"
	"github.com/pkg/errors"
)

// Menu implements the app.Menu interface.
type Menu struct {
	id        uuid.UUID
	markup    app.Markup
	lastFocus time.Time
	component app.Component

	onClose func()
}

func newMenu(c app.MenuConfig, name string) (app.Menu, error) {
	var markup app.Markup = html.NewMarkup(driver.factory)
	markup = app.NewConcurrentMarkup(markup)

	rawMenu := &Menu{
		id:        uuid.New(),
		markup:    markup,
		lastFocus: time.Now(),

		onClose: c.OnClose,
	}

	menu := app.NewMenuWithLogs(rawMenu, name)

	if err := driver.macRPC.Call("menus.New", nil, struct {
		ID string
	}{
		ID: menu.ID().String(),
	}); err != nil {
		return nil, err
	}

	if err := driver.elements.Add(menu); err != nil {
		return nil, err
	}

	if len(c.DefaultURL) != 0 {
		if err := menu.Load(c.DefaultURL); err != nil {
			return nil, err
		}
	}
	return menu, nil
}

// ID satisfies the app.Menu interface.
func (m *Menu) ID() uuid.UUID {
	return m.id
}

// Base satisfies the app.Menu interface.
func (m *Menu) Base() app.Menu {
	return m
}

// Load satisfies the app.Menu interface.
func (m *Menu) Load(rawurl string, v ...interface{}) error {
	rawurl = fmt.Sprintf(rawurl, v...)
	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}

	var compo app.Component
	compo, err = driver.factory.New(app.ComponentNameFromURL(u))
	if err != nil {
		return err
	}

	if m.component != nil {
		m.markup.Dismount(m.component)
	}

	if _, err = m.markup.Mount(compo); err != nil {
		return err
	}
	m.component = compo

	if navigable, ok := compo.(app.Navigable); ok {
		navigable.OnNavigate(u)
	}

	var root app.Tag
	if root, err = m.markup.Root(compo); err != nil {
		return err
	}
	if root, err = m.markup.FullRoot(root); err != nil {
		return err
	}

	return driver.macRPC.Call("menus.Load", nil, struct {
		ID  string
		Tag app.Tag
	}{
		ID:  m.ID().String(),
		Tag: root,
	})
}

// Component satisfies the app.Menu interface.
func (m *Menu) Component() app.Component {
	return m.component
}

// Contains satisfies the app.Menu interface.
func (m *Menu) Contains(compo app.Component) bool {
	return m.markup.Contains(compo)
}

// Render satisfies the app.Menu interface.
func (m *Menu) Render(compo app.Component) error {
	syncs, err := m.markup.Update(compo)
	if err != nil {
		return err
	}

	for _, sync := range syncs {
		if sync.Replace {
			err = m.render(sync)
		} else {
			err = m.renderAttributes(compo, sync)
		}

		if err != nil {
			return err
		}
	}
	return nil
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
		ID:  m.ID().String(),
		Tag: tag,
	})
}

func (m *Menu) renderAttributes(compo app.Component, sync app.TagSync) error {
	root, err := m.markup.Root(compo)
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
		ID:  m.ID().String(),
		Tag: tag,
	})
}

// LastFocus satisfies the app.Menu interface.
func (m *Menu) LastFocus() time.Time {
	return m.lastFocus
}

func onMenuClose(m *Menu, in map[string]interface{}) interface{} {
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
			ID: m.ID().String(),
		}); err != nil {
			panic(errors.Wrap(err, "onMenuClose"))
		}

		driver.elements.Remove(m)
	})

	return nil
}

func onMenuCallback(m *Menu, in map[string]interface{}) interface{} {
	mappingString := in["Mapping"].(string)

	var mapping app.Mapping
	if err := json.Unmarshal([]byte(mappingString), &mapping); err != nil {
		app.Error(errors.Wrap(err, "onMenuCallback"))
		return nil
	}

	function, err := m.markup.Map(mapping)
	if err != nil {
		app.Error(errors.Wrap(err, "onMenuCallback"))
		return nil
	}

	if function != nil {
		function()
		return nil
	}

	var compo app.Component
	if compo, err = m.markup.Component(mapping.CompoID); err != nil {
		app.Error(errors.Wrap(err, "onMenuCallback"))
		return nil
	}

	if err = m.Render(compo); err != nil {
		app.Error(errors.Wrap(err, "onMenuCallback"))
	}
	return nil
}

func handleMenu(h func(m *Menu, in map[string]interface{}) interface{}) bridge.GoRPCHandler {
	return func(in map[string]interface{}) interface{} {
		id, _ := uuid.Parse(in["ID"].(string))

		elem, err := driver.elements.Element(id)
		if err != nil {
			return nil
		}

		menu := elem.(app.Menu).Base().(*Menu)
		return h(menu, in)
	}
}
