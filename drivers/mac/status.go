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

	menu := &StatusMenu{
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

// StatusMenu represents a menu that lives in the status bar.
type StatusMenu struct {
	Menu
	onClose func()
}

// Load loads the component targetted by the given url.
// It satisfies the app.StatusMenu interface.
func (s *StatusMenu) Load(url string, v ...interface{}) error {
	if err := s.Menu.Load(url, v...); err != nil {
		return err
	}

	return driver.macRPC.Call("statusMenus.SetMenu", nil, struct {
		ID string
	}{
		ID: s.ID().String(),
	})
}

// SetText set the status menu text.
// It satisfies the app.StatusMenu interface.
func (s *StatusMenu) SetText(text string) error {
	return driver.macRPC.Call("statusMenus.SetText", nil, struct {
		ID   string
		Text string
	}{
		ID:   s.ID().String(),
		Text: text,
	})
}

// SetIcon set the status menu icon.
// It satisfies the app.StatusMenu interface.
func (s *StatusMenu) SetIcon(name string) error {
	if _, err := os.Stat(name); err != nil && len(name) != 0 {
		return err
	}

	return driver.macRPC.Call("statusMenus.SetIcon", nil, struct {
		ID   string
		Icon string
	}{
		ID:   s.ID().String(),
		Icon: name,
	})
}

// Close closes the status menu, releasing allocated resources and removing
// it from the status bar.
// It satisfies the app.StatusMenu interface.
func (s *StatusMenu) Close() error {
	return driver.macRPC.Call("statusMenus.Close", nil, struct {
		ID string
	}{
		ID: s.ID().String(),
	})
}
