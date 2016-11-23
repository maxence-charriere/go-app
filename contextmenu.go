package app

type ContextMenu struct{}

// NewContextMenu creates a new context menu.
func NewContextMenu() Contexter {
	return driver.NewContext(ContextMenu{})
}
