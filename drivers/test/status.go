package test

import (
	"os"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/dom"
)

// StatusMenu is a teststatus menu that implements the app.StatusMenu interface.
type StatusMenu struct {
	Menu
}

func newStatusMenu(d *Driver, c app.StatusMenuConfig) *StatusMenu {
	s := &StatusMenu{
		Menu{
			driver: d,
			dom:    dom.NewDOM(d.factory),
			id:     uuid.New().String(),
		},
	}

	d.elems.Put(s)

	if len(c.URL) != 0 {
		s.Load(c.URL)
	}

	return s
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
	_, err := os.Stat(path)
	s.SetErr(err)
}

// SetText satisfies the app.StatusMenu interface.
func (s *StatusMenu) SetText(text string) {
	s.SetErr(nil)
	s.driver.setElemErr(s)
}

// Close satisfies the app.StatusMenu interface.
func (s *StatusMenu) Close() {
	s.driver.elems.Delete(s)
	s.SetErr(nil)
	s.driver.setElemErr(s)
}
