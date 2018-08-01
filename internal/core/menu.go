package core

import (
	"github.com/murlokswarm/app"
)

// Menu is a base struct to embed in app.Window implementations.
type Menu struct {
	Elem
}

// WhenMenu satisfies the app.Menu interface.
func (m *Menu) WhenMenu(f func(app.Menu)) {
	f(m)
}

// Load satisfies the app.Menu interface.
func (m *Menu) Load(url string, v ...interface{}) {
	m.SetErr(app.ErrNotSupported)
}

// Compo satisfies the app.Menu interface.
func (m *Menu) Compo() app.Compo {
	return nil
}

// Contains satisfies the app.Menu interface.
func (m *Menu) Contains(c app.Compo) bool {
	return false
}

// Render satisfies the app.Menu interface.
func (m *Menu) Render(c app.Compo) {
	m.SetErr(app.ErrNotSupported)
}

// Type satisfies the app.Menu interface.
func (m *Menu) Type() string {
	return "menu"
}

// DockTile is a base struct to embed in app.DockTile implementations.
type DockTile struct {
	Menu
}

// SetIcon satisfies the app.DockTile interface.
func (d *DockTile) SetIcon(name string) {
	d.SetErr(app.ErrNotSupported)
}

// SetBadge satisfies the app.DockTile interface.
func (d *DockTile) SetBadge(v interface{}) {
	d.SetErr(app.ErrNotSupported)
}
