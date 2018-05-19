package app

// Menu is the interface that describes a menu.
// Accept only components that contain menu and menuitem tags.
type Menu interface {
	ElementWithComponent
}

// MenuConfig is a struct that describes a menu.
type MenuConfig struct {
	DefaultURL string

	OnClose func() `json:"-"`
}
