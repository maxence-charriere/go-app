package app

import "fmt"

// StatusMenuConfig is a struct that describes a status menu.
// Accept only components that contain menu and menuitem tags.
type StatusMenuConfig struct {
	// The menu button text.
	Text string

	// The menu button icon.
	// Should be a .png file.
	Icon string

	// The URL of the component to load when the status menu is created.
	DefaultURL string

	// The function that is called when the status menu is closed.
	OnClose func()
}

// StatusMenu is the interface that describes a status menu menu.
type StatusMenu interface {
	Menu
	Closer

	// Set the menu button icon.
	// The icon should be a .png file.
	SetIcon(name string) error

	// Set the menu button text.
	SetText(text string) error
}

type statusMenuWithLogs struct {
	StatusMenu
}

func (s *statusMenuWithLogs) Load(url string, v ...interface{}) {
	parsedURL := fmt.Sprintf(url, v...)

	WhenDebug(func() {
		Debug("status menu %s is loading %s",
			s.ID(),
			parsedURL,
		)
	})

	s.StatusMenu.Load(url, v...)
	if s.Err() != nil {
		Log("status menu %T failed to load %s: %s",
			s.ID(),
			parsedURL,
			s.Err(),
		)
	}
}

func (s *statusMenuWithLogs) Render(c Compo) {
	WhenDebug(func() {
		Debug("status menu %s is rendering %T",
			s.ID(),
			c,
		)
	})

	s.StatusMenu.Render(c)
	if s.Err() != nil {
		Log("status menu %s failed to render %T: %s",
			s.ID(),
			c,
			s.Err(),
		)
	}
}

func (s *statusMenuWithLogs) SetIcon(name string) error {
	WhenDebug(func() {
		Debug("status menu %s is setting icon to %s",
			s.ID(),
			name,
		)
	})

	err := s.StatusMenu.SetIcon(name)
	if err != nil {
		Log("status menu %s failed to set icon: %s",
			s.ID(),
			err,
		)
	}
	return err
}

func (s *statusMenuWithLogs) SetText(text string) error {
	WhenDebug(func() {
		Debug("status menu %s is setting text to %s",
			s.ID(),
			text,
		)
	})

	err := s.StatusMenu.SetText(text)
	if err != nil {
		Log("status menu %s failed to set text: %s",
			s.ID(),
			err,
		)
	}
	return err
}

func (s *statusMenuWithLogs) Close() {
	WhenDebug(func() {
		Debug("status menu %s is closing", s.ID())
	})

	s.StatusMenu.Close()
	if s.Err() != nil {
		Log("status menu %s failed to close: %s",
			s.ID(),
			s.Err(),
		)
	}
}
