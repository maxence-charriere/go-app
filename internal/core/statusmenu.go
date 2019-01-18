package core

import (
	"os"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
)

// StatusMenu is a modular implementation of the app.StatusMenu interface that
// can be configured to address the different drivers needs.
type StatusMenu struct {
	Menu
}

// Create creates and display the status menu.
func (s *StatusMenu) Create(c app.StatusMenuConfig) {
	s.id = uuid.New().String()
	s.DOM.AllowedNodes = []string{"menu", "menuitem"}
	s.DOM.Factory = s.Driver.Factory
	s.DOM.Sync = s.render
	s.DOM.UI = s.Driver.UI
	s.NoDestroy = true

	if s.err = s.Driver.Platform.Call("statusMenus.New", nil, struct {
		ID   string
		Text string
		Icon string
	}{
		ID:   s.id,
		Text: c.Text,
		Icon: c.Icon,
	}); s.err != nil {
		return
	}

	s.Driver.Elems.Put(s)

	if len(c.URL) != 0 {
		s.Load(c.URL)
	}
}

// WhenStatusMenu satisfies the app.StatusMenu interface.
func (s *StatusMenu) WhenStatusMenu(f func(app.StatusMenu)) {
	f(s)
}

// Load the app.StatusMenu interface.
func (s *StatusMenu) Load(rawurl string) {
	s.Menu.Load(rawurl)
	if s.err != nil {
		return
	}

	s.err = s.Driver.Platform.Call("statusMenus.SetMenu", nil, struct {
		ID string
	}{
		ID: s.id,
	})
}

// SetIcon satisfies the app.StatusMenu interface.
func (s *StatusMenu) SetIcon(path string) {
	if _, s.err = os.Stat(path); len(path) != 0 && s.err != nil {
		return
	}

	s.err = s.Driver.Platform.Call("statusMenus.SetIcon", nil, struct {
		ID   string
		Icon string
	}{
		ID:   s.id,
		Icon: path,
	})
}

// SetText satisfies the app.StatusMenu interface.
func (s *StatusMenu) SetText(text string) {
	s.err = s.Driver.Platform.Call("statusMenus.SetText", nil, struct {
		ID   string
		Text string
	}{
		ID:   s.id,
		Text: text,
	})
}

// Close satisfies the app.StatusMenu interface.
func (s *StatusMenu) Close() {
	s.err = s.Driver.Platform.Call("statusMenus.Close", nil, struct {
		ID string
	}{
		ID: s.id,
	})

	s.Driver.Elems.Delete(s)
}
