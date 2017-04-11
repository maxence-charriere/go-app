package app

// Menu is a struct that describes a menu.
// It will be used by a driver to create a menu on the top of a native menu.
type Menu struct{}

// ContextMenu is a struct that describes a context menu.
// It will be used by a driver to create a context on the top of a native
// context menu.
type ContextMenu Menu

// NewContextMenu creates a new context menu.
func NewContextMenu() Contexter {
	return driver.NewElement(ContextMenu{}).(Contexter)
}
