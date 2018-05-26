package app

import "fmt"

// StatusBarMenu is the interface that describes a status bar menu.
type StatusBarMenu interface {
	Menu

	// Set the menu icon in the status bar.
	// The icon should be a .png file.
	SetIcon(name string) error
}

type statusBarWithLogs struct {
	StatusBarMenu
}

func (s *statusBarWithLogs) Load(url string, v ...interface{}) error {
	parsedURL := fmt.Sprintf(url, v...)

	WhenDebug(func() {
		Debug("status bar is loading %s", parsedURL)
	})

	err := s.StatusBarMenu.Load(url, v...)
	if err != nil {
		Log("status bar failed to load %s: %s",
			parsedURL,
			err,
		)
	}
	return err
}

func (s *statusBarWithLogs) Render(c Component) error {
	WhenDebug(func() {
		Debug("status bar is rendering %T", c)
	})

	err := s.StatusBarMenu.Render(c)
	if err != nil {
		Log("status bar menu failed to render %T: %s",
			c,
			err,
		)
	}
	return err
}

func (s *statusBarWithLogs) SetIcon(name string) error {
	WhenDebug(func() {
		Debug("status bar is setting icon to %s", name)
	})

	err := s.StatusBarMenu.SetIcon(name)
	if err != nil {
		Log("status bar failed to set icon: %s", err)
	}
	return err
}
