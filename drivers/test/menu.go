package test

import (
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/html"
	"github.com/pkg/errors"
)

// A Menu implementation for tests.
type Menu struct {
	id          uuid.UUID
	typ         string
	factory     app.Factory
	markup      app.Markup
	lastFocus   time.Time
	component   app.Component
	simulateErr bool

	onClose func()
}

func newMenu(d *Driver, c app.MenuConfig) (app.Menu, error) {
	var markup app.Markup = html.NewMarkup(d.factory)
	markup = app.ConcurrentMarkup(markup)

	menu := &Menu{
		id:          uuid.New(),
		typ:         c.Type,
		factory:     d.factory,
		markup:      markup,
		lastFocus:   time.Now(),
		simulateErr: d.SimulateElemErr,
	}

	d.elements.Add(menu)
	menu.onClose = func() {
		d.elements.Remove(menu)
	}

	var err error
	if len(c.DefaultURL) != 0 {
		err = menu.Load(c.DefaultURL)
	}
	return menu, err
}

// ID satisfies the app.Menu interface.
func (m *Menu) ID() uuid.UUID {
	return m.id
}

// Load satisfies the app.ElementWithComponent interface.
func (m *Menu) Load(rawurl string, v ...interface{}) error {
	if m.simulateErr {
		return ErrSimulated
	}

	if m.component != nil {
		m.markup.Dismount(m.component)
	}

	rawurl = fmt.Sprintf(rawurl, v...)

	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}

	compo, err := m.factory.New(app.ComponentNameFromURL(u))
	if err != nil {
		return err
	}

	m.component = compo

	if _, err = m.markup.Mount(compo); err != nil {
		return errors.Wrapf(err, "loading %s in test menu %p failed", u, m)
	}
	return nil
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
	if m.simulateErr {
		return ErrSimulated
	}

	_, err := m.markup.Update(compo)
	return err
}

// LastFocus satisfies the app.Menu interface.
func (m *Menu) LastFocus() time.Time {
	return m.lastFocus
}

// Type satisfies the app.Menu interface.
func (m *Menu) Type() string {
	return m.typ
}
