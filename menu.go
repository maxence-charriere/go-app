package app

// ContextMenu is a struct that describes a context menu.
// It will be used by a driver to create a context on the top of a native
// context menu.
type ContextMenu struct{}

// NewContextMenu creates a new context menu.
func NewContextMenu() Contexter {
	return driver.NewElement(ContextMenu{}).(Contexter)
}
