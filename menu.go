package app

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
	URL string

	// The function that is called when the menu is closed.
	OnClose func() `json:"-"`
}

// StatusMenu is the interface that describes a status menu menu.
type StatusMenu interface {
	Menu
	Closer

	// Set the menu button icon.
	// The icon should be a .png file.
	SetIcon(path string)

	// Set the menu button text.
	SetText(text string)
}

// StatusMenuConfig is a struct that describes a status menu.
// Accept only components that contain menu and menuitem tags.
type StatusMenuConfig struct {
	// The menu button text.
	Text string

	// The menu button icon.
	// Should be a .png file.
	Icon string

	// The URL of the component to load when the status menu is created.
	URL string

	// The function that is called when the status menu is closed.
	OnClose func() `json:"-"`
}

// DockTile is the interface that describes a dock tile.
// Accept only components that contain menu and menuitem tags.
type DockTile interface {
	Menu

	// SetIcon set the dock tile icon with the named file.
	// The icon should be a .png file.
	SetIcon(path string)

	// SetBadge set the dock tile badge with the string representation of the
	// value.
	SetBadge(v interface{})
}
