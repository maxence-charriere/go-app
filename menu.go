package app

import "fmt"

// Menu is the interface that describes a menu.
// Accept only components that contain menu and menuitem tags.
type Menu interface {
	ElemWithCompo

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

func (m *menuWithLogs) Load(url string, v ...interface{}) {
	parsedURL := fmt.Sprintf(url, v...)

	WhenDebug(func() {
		Debug("%s %s is loading %s",
			m.Type(),
			m.ID(),
			parsedURL,
		)
	})

	m.Menu.Load(url, v...)
	if m.Err() != nil {
		Log("%s %s failed to load %s: %s",
			m.Type(),
			m.ID(),
			parsedURL,
			m.Err(),
		)
	}
}

func (m *menuWithLogs) Render(c Compo) {
	WhenDebug(func() {
		Debug("%s %s is rendering %T",
			m.Type(),
			m.ID(),
			c,
		)
	})

	m.Menu.Render(c)
	if m.Err() != nil {
		Log("%s %s failed to render %T: %s",
			m.Type(),
			m.ID(),
			c,
			m.Err(),
		)
	}
}
