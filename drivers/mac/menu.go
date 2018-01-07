// +build darwin,amd64

package mac

import (
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/html"
)

// Menu implements the app.Menu interface.
type Menu struct {
	id        uuid.UUID
	markup    app.Markup
	lastFocus time.Time
}

func newMenu(config app.MenuConfig) (m *Menu, err error) {
	var markup app.Markup = html.NewMarkup(driver.factory)
	markup = app.NewConcurrentMarkup(markup)

	m = &Menu{
		id:        uuid.New(),
		markup:    markup,
		lastFocus: time.Now(),
	}

	if err = driver.elements.Add(m); err != nil {
		return
	}

	if len(config.DefaultURL) != 0 {

	}

	return
}

// ID satisfies the app.Element interface.
func (m *Menu) ID() uuid.UUID {
	return m.id
}

// Load satisfies the app.Menu interface.
func (m *Menu) Load(url string, v ...interface{}) error {
	panic("not implemented")
}

// Contains satisfies the app.Menu interface.
func (m *Menu) Contains(c app.Component) bool {
	panic("not implemented")
}

// Render satisfies the app.Menu interface.
func (m *Menu) Render(c app.Component) error {
	panic("not implemented")
}

// LastFocus satisfies the app.Menu interface.
func (m *Menu) LastFocus() time.Time {
	panic("not implemented")
}
