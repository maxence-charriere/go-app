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
			id:             uuid.New(),
			markup:         markup,
			lastFocus:      time.Now(),
			keepWhenClosed: true,
		},
		onClose: c.OnClose,
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

	if len(c.DefaultURL) != 0 {
		return menu, menu.Load(c.DefaultURL)
	}
	return menu, nil
}

type statusMenu struct {
	Menu
	onClose func()
}

func (s *statusMenu) Load(url string, v ...interface{}) error {
	if err := s.Menu.Load(url, v...); err != nil {
		return err
	}

	return driver.macRPC.Call("statusMenus.SetMenu", nil, struct {
		ID string
	}{
		ID: s.ID().String(),
	})
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
	return driver.macRPC.Call("statusMenus.Close", nil, struct {
		ID string
	}{
		ID: s.ID().String(),
	})
}
