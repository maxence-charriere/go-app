package test

import (
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/html"
)

// A StatusMenu implementation for tests.
type StatusMenu struct {
	Menu
}

func newStatusMenu(d *Driver, c app.StatusMenuConfig) app.StatusMenu {
	var markup app.Markup = html.NewMarkup(d.factory)
	markup = app.ConcurrentMarkup(markup)

	menu := &StatusMenu{
		Menu: Menu{
			id:          uuid.New(),
			typ:         "status menu",
			factory:     d.factory,
			markup:      markup,
			lastFocus:   time.Now(),
			simulateErr: d.SimulateElemErr,
		},
	}

	d.elements.Add(menu)
	return menu
}

// SetText satisfies the app.StatusBar interface.
func (s *StatusMenu) SetText(text string) error {
	if s.simulateErr {
		return ErrSimulated
	}
	return nil
}

// SetIcon satisfies the app.StatusBar interface.
func (s *StatusMenu) SetIcon(name string) error {
	if s.simulateErr {
		return ErrSimulated
	}
	_, err := os.Stat(name)
	return err
}

// Close satisfies the app.StatusBar interface.
func (s *StatusMenu) Close() error {
	if s.simulateErr {
		return ErrSimulated
	}
	return nil
}
