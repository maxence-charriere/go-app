package test

import (
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/markup"
	"github.com/pkg/errors"
)

// A Menu implementation for tests.
type Menu struct {
	config       app.MenuConfig
	id           uuid.UUID
	compoBuilder markup.CompoBuilder
	env          markup.Env
	lastFocus    time.Time
}

func newMenu(d *Driver, c app.MenuConfig) app.Menu {
	menu := &Menu{
		id:           uuid.New(),
		compoBuilder: d.compoBuilder,
		env:          markup.NewEnv(d.compoBuilder),
		lastFocus:    time.Now(),
	}
	d.elements.Add(menu)

	if len(c.DefaultURL) != 0 {
		if err := menu.Load(c.DefaultURL); err != nil {
			d.Test.Log(err)
		}
	}
	return menu
}

// ID satisfies the app.Element interface.
func (m *Menu) ID() uuid.UUID {
	return m.id
}

// Load satisfies the app.ElementWithComponent interface.
func (m *Menu) Load(rawurl string) error {
	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}

	componame, ok := markup.ComponentNameFromURL(u)
	if !ok {
		return nil
	}

	compo, err := m.compoBuilder.New(componame)
	if err != nil {
		return err
	}

	if _, err = m.env.Mount(compo); err != nil {
		return errors.Wrapf(err, "loading %s in test menu %p failed", u, m)
	}
	return nil
}

// Contains satisfies the app.ElementWithComponent interface.
func (m *Menu) Contains(c markup.Component) bool {
	return m.env.Contains(c)
}

// Render satisfies the app.ElementWithComponent interface.
func (m *Menu) Render(c markup.Component) error {
	_, err := m.env.Update(c)
	return err
}

// LastFocus satisfies the app.ElementWithComponent interface.
func (m *Menu) LastFocus() time.Time {
	return m.lastFocus
}
