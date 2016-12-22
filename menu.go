package app

// Menu represents a context menu.
type Menu struct{}

// ContextMenu represents a context menu.
type ContextMenu Menu

// NewContextMenu creates a new context menu.
func NewContextMenu() Contexter {
	return driver.NewContext(ContextMenu{})
}
