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
	id        uuid.UUID
	factory   app.Factory
	markup    app.Markup
	lastFocus time.Time

	onLoad func(compo app.Component)
}

func newMenu(d *Driver, c app.MenuConfig) (app.Menu, error) {
	menu := &Menu{
		id:        uuid.New(),
		factory:   d.factory,
		markup:    html.NewMarkup(d.factory),
		lastFocus: time.Now(),
	}
	d.elements.Add(menu)

	var err error
	if len(c.DefaultURL) != 0 {
		err = menu.Load(c.DefaultURL)
	}
	return menu, err
}

// ID satisfies the app.Element interface.
func (m *Menu) ID() uuid.UUID {
	return m.id
}

// Load satisfies the app.ElementWithComponent interface.
func (m *Menu) Load(rawurl string, v ...interface{}) error {
	rawurl = fmt.Sprintf(rawurl, v...)

	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}

	compo, err := m.factory.NewComponent(app.ComponentNameFromURL(u))
	if err != nil {
		return err
	}

	if _, err = m.markup.Mount(compo); err != nil {
		return errors.Wrapf(err, "loading %s in test menu %p failed", u, m)
	}

	if m.onLoad != nil {
		m.onLoad(compo)
	}
	return nil
}

// Contains satisfies the app.ElementWithComponent interface.
func (m *Menu) Contains(compo app.Component) bool {
	return m.markup.Contains(compo)
}

// Render satisfies the app.ElementWithComponent interface.
func (m *Menu) Render(compo app.Component) error {
	_, err := m.markup.Update(compo)
	return err
}

// LastFocus satisfies the app.ElementWithComponent interface.
func (m *Menu) LastFocus() time.Time {
	return m.lastFocus
}
