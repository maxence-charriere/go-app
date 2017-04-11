package app

// Share is a struct that describes a share.
// It will be used by a driver to create a native share panel.
type Share struct {
	Value interface{}
}

// NewShare creates a new sharing.
func NewShare(s Share) Elementer {
	return driver.NewElement(s)
}
