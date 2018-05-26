package test

import (
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/html"
)

// A StatusBarMenu implementation for tests.
type StatusBarMenu struct {
	Menu
}

func newStatusBarMenu(d *Driver) app.StatusBarMenu {
	var markup app.Markup = html.NewMarkup(d.factory)
	markup = app.ConcurrentMarkup(markup)

	menu := &StatusBarMenu{
		Menu: Menu{
			id:          uuid.New(),
			typ:         "status bar",
			factory:     d.factory,
			markup:      markup,
			lastFocus:   time.Now(),
			simulateErr: d.SimulateElemErr,
		},
	}

	d.elements.Add(menu)
	return menu
}

// SetIcon satisfies the app.StatusBar interface.
func (s *StatusBarMenu) SetIcon(name string) error {
	if s.simulateErr {
		return ErrSimulated
	}
	_, err := os.Stat(name)
	return err
}
