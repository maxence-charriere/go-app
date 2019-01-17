package core

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/dom"
	"github.com/pkg/errors"
)

// Menu is a modular implementation of the app.Menu interface that can be
// configured to address the different drivers needs.
type Menu struct {
	Elem

	DOM       dom.Engine
	Driver    *Driver
	History   History
	NoDestroy bool

	compo app.Compo
	kind  string
}

// Create creates and display the menu.
func (m *Menu) Create(c app.MenuConfig) {
	m.id = uuid.New().String()
	m.DOM.AllowedNodes = []string{"menu", "menuitem"}
	m.DOM.Factory = m.Driver.Factory
	m.DOM.Sync = m.render
	m.DOM.UI = m.Driver.UI

	if m.err = m.Driver.Platform.Call("menus.New", nil, struct {
		ID string
	}{
		ID: m.id,
	}); m.err != nil {
		return
	}

	m.Driver.Elems.Put(m)

	if len(c.URL) != 0 {
		m.Load(c.URL)
	}
}

// Contains satisfies the app.Menu interface.
func (m *Menu) Contains(c app.Compo) bool {
	return m.DOM.Contains(c)
}

// WhenMenu satisfies the app.Menu interface.
func (m *Menu) WhenMenu(f func(app.Menu)) {
	f(m)
}

// Load satisfies the app.Menu interface.
func (m *Menu) Load(rawurl string) {
	compoName := CompoNameFromURLString(rawurl)

	if m.compo, m.err = m.Driver.Factory.NewCompo(compoName); m.err != nil {
		return
	}

	if rawurl != m.History.Current() {
		m.History.NewEntry(rawurl)
	}

	if m.err = m.Driver.Platform.Call("menus.Load", nil, struct {
		ID   string
		Kind string
	}{
		ID:   m.id,
		Kind: m.kind,
	}); m.err != nil {
		return
	}

	if m.err = m.DOM.New(m.compo); m.err != nil {
		return
	}

	if nav, ok := m.compo.(app.Navigable); ok {
		u, _ := url.Parse(rawurl)

		m.Driver.UI(func() {
			nav.OnNavigate(u)
		})
	}
}

// Reload satisfies the app.Menu interface.
func (m *Menu) Reload() {
	url := m.History.Current()
	if len(url) == 0 {
		m.err = errors.New("no component to reload")
		return
	}

	m.Load(url)
}

// CanPrevious satisfies the app.Menu interface.
func (m *Menu) CanPrevious() bool {
	return m.History.CanPrevious()
}

// Previous satisfies the app.Menu interface.
func (m *Menu) Previous() {
	url := m.History.Previous()
	if len(url) == 0 {
		m.err = errors.New("no previous component to load")
		return
	}

	m.Load(url)
}

// CanNext satisfies the app.Menu interface.
func (m *Menu) CanNext() bool {
	return m.History.CanNext()
}

// Next satisfies the app.Menu interface.
func (m *Menu) Next() {
	url := m.History.Next()
	if len(url) == 0 {
		m.err = errors.New("no next component to load")
		return
	}

	m.Load(url)
}

// Compo satisfies the app.Menu interface.
func (m *Menu) Compo() app.Compo {
	return m.compo
}

// Render satisfies the app.Menu interface.
func (m *Menu) Render(c app.Compo) {
	m.err = m.DOM.Render(c)
}

func (m *Menu) render(changes interface{}) error {
	b, err := json.Marshal(changes)
	if err != nil {
		return errors.Wrap(err, "encoding changes failed")
	}

	return m.Driver.Platform.Call("menus.Render", nil, struct {
		ID      string
		Changes string
	}{
		ID:      m.id,
		Changes: string(b),
	})
}

// Kind satisfies the app.Menu interface.
func (m *Menu) Kind() string {
	return m.kind
}

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
	if _, d.err = os.Stat(path); path != "" && d.err != nil {
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

// StatusMenu is a base struct to embed in app.StatusMenu implementations.
type StatusMenu struct {
	Menu
}

// WhenStatusMenu satisfies the app.StatusMenu interface.
func (s *StatusMenu) WhenStatusMenu(f func(app.StatusMenu)) {
	f(s)
}

// Type satisfies the app.Menu interface.
func (s *StatusMenu) Type() string {
	return "status menu"
}

// SetIcon satisfies the app.StatusMenu interface.
func (s *StatusMenu) SetIcon(path string) {
	s.SetErr(app.ErrNotSupported)
}

// SetText satisfies the app.StatusMenu interface.
func (s *StatusMenu) SetText(text string) {
	s.SetErr(app.ErrNotSupported)
}

// Close satisfies the app.StatusMenu interface.
func (s *StatusMenu) Close() {
	s.SetErr(app.ErrNotSupported)
}
