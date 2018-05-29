// +build darwin,amd64

package mac

import (
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/html"
)

func newStatusMenu(c app.StatusMenuConfig) (app.StatusMenu, error) {
	var markup app.Markup = html.NewMarkup(driver.factory)
	markup = app.ConcurrentMarkup(markup)

	menu := &statusMenu{
		Menu: Menu{
			id:        uuid.New(),
			markup:    markup,
			lastFocus: time.Now(),
		},
	}

	if err := driver.macRPC.Call("statusMenus.New", nil, struct {
		ID   string
		Text string
		Icon string
	}{
		ID:   menu.ID().String(),
		Text: c.Text,
		Icon: c.Icon,
	}); err != nil {
		return nil, err
	}

	if err := driver.elements.Add(menu); err != nil {
		return nil, err
	}

	// implement autoload.

	return menu, nil
}

type statusMenu struct {
	Menu
}

func (s *statusMenu) SetText(text string) error {
	panic("not implemented")
	// return driver.macRPC.Call("driver.SetStatusBarIcon", nil, name)
}

func (s *statusMenu) SetIcon(name string) error {
	if _, err := os.Stat(name); err != nil && len(name) != 0 {
		return err
	}

	panic("not implemented")
	// return driver.macRPC.Call("driver.SetStatusBarIcon", nil, name)
}

func (s *statusMenu) Close() error {
	// TO IMPLEMENT.
	return nil
}
