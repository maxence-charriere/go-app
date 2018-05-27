// +build darwin,amd64

package mac

import (
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/html"
)

func newStatusBar() (app.StatusBarMenu, error) {
	var markup app.Markup = html.NewMarkup(driver.factory)
	markup = app.ConcurrentMarkup(markup)

	statbar := &statusBar{
		Menu: Menu{
			id:        uuid.New(),
			markup:    markup,
			lastFocus: time.Now(),
		},
	}

	if err := driver.macRPC.Call("menus.New", nil, struct {
		ID string
	}{
		ID: statbar.ID().String(),
	}); err != nil {
		return nil, err
	}

	if err := driver.elements.Add(statbar); err != nil {
		return nil, err
	}

	if err := driver.macRPC.Call("driver.SetStatusBar", nil, statbar.ID()); err != nil {
		return nil, err
	}

	return statbar, nil
}

type statusBar struct {
	Menu
}

func (s *statusBar) SetIcon(name string) error {
	if _, err := os.Stat(name); err != nil && len(name) != 0 {
		return err
	}

	return driver.macRPC.Call("driver.SetStatusBarIcon", nil, name)
}
