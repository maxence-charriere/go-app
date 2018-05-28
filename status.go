package app

import "fmt"

// StatusMenuConfig is a struct that describes a status menu menu.
type StatusMenuConfig struct {
	// The menu button text.
	Text string

	// The menu button icon.
	// Should be a .png file.
	Icon string
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

func (s *statusMenuWithLogs) Load(url string, v ...interface{}) error {
	parsedURL := fmt.Sprintf(url, v...)

	WhenDebug(func() {
		Debug("status menu %s is loading %s",
			s.ID(),
			parsedURL,
		)
	})

	err := s.StatusMenu.Load(url, v...)
	if err != nil {
		Log("status menu %T failed to load %s: %s",
			s.ID(),
			parsedURL,
			err,
		)
	}
	return err
}

func (s *statusMenuWithLogs) Render(c Component) error {
	WhenDebug(func() {
		Debug("status menu %s is rendering %T",
			s.ID(),
			c,
		)
	})

	err := s.StatusMenu.Render(c)
	if err != nil {
		Log("status menu %s failed to render %T: %s",
			s.ID(),
			c,
			err,
		)
	}
	return err
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

func (s *statusMenuWithLogs) Close() error {
	WhenDebug(func() {
		Debug("status menu %s is closing", s.ID())
	})

	err := s.StatusMenu.Close()
	if err != nil {
		Log("status menu %s failed to close: %s",
			s.ID(),
			err,
		)
	}
	return err
}
