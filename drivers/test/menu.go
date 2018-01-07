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
	config    app.MenuConfig
	id        uuid.UUID
	factory   app.Factory
	markup    app.Markup
	lastFocus time.Time
}

func newMenu(driver *Driver, config app.MenuConfig) app.Menu {
	menu := &Menu{
		id:        uuid.New(),
		factory:   driver.factory,
		markup:    app.NewConcurrentMarkup(html.NewMarkup(driver.factory)),
		lastFocus: time.Now(),
	}
	driver.elements.Add(menu)

	if len(config.DefaultURL) != 0 {
		if err := menu.Load(config.DefaultURL); err != nil {
			driver.Test.Log(err)
		}
	}
	return menu
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
