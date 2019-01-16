package app

// Menu is the interface that describes a menu.
// Accept only components that contain menu and menuitem tags.
type Menu interface {
	View

	// The menu type.
	Kind() string
}

// MenuConfig is a struct that describes a menu.
type MenuConfig struct {
	// The URL of the component to load when the menu is created.
	URL string
}

// MenuBarConfig is a struct that describes a menu bar.
type MenuBarConfig struct {
	// The URL of the component that describes the app menu. A default app menu
	// is displayed if this URL is empty.
	AppURL string

	// The URLs of custom menus to display in the menu bar.
	CustomURLs []string

	// The URL of the component that describes the edit menu. A default edit
	// menu is displayed if this URL is empty.
	EditURL string

	// The URL of the component that describes the file menu. File menu is not
	// displayed if this URL is empty.
	FileURL string

	// The URL of the component that describes the help menu. A default help
	// menu is displayed if this URL is empty.
	HelpURL string

	// The URL of the component that describes the window menu. A default window
	// menu is displayed if this URL is empty.
	WindowURL string
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
