package app

import "fmt"

// Menu is the interface that describes a menu.
// Accept only components that contain menu and menuitem tags.
type Menu interface {
	ElementWithComponent

	// The menu type.
	Type() string
}

// MenuConfig is a struct that describes a menu.
type MenuConfig struct {
	// The URL of the component to load when the menu is created.
	DefaultURL string

	Type string

	// The function that is called when the menu is closed.
	OnClose func() `json:"-"`
}

type menuWithLogs struct {
	Menu
}

func (m *menuWithLogs) Load(url string, v ...interface{}) error {
	parsedURL := fmt.Sprintf(url, v...)

	WhenDebug(func() {
		Debug("%s %s is loading %s",
			m.Type(),
			m.ID(),
			parsedURL,
		)
	})

	err := m.Menu.Load(url, v...)
	if err != nil {
		Log("%s %s failed to load %s: %s",
			m.Type(),
			m.ID(),
			parsedURL,
			err,
		)
	}
	return err
}

func (m *menuWithLogs) Render(c Component) error {
	WhenDebug(func() {
		Debug("%s %s is rendering %T",
			m.Type(),
			m.ID(),
			c,
		)
	})

	err := m.Menu.Render(c)
	if err != nil {
		Log("%s %s failed to render %T: %s",
			m.Type(),
			m.ID(),
			c,
			err,
		)
	}
	return err
}
