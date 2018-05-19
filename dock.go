package app

// DockTile is the interface that describes a dock tile.
type DockTile interface {
	ElementWithComponent

	// SetIcon set the dock tile icon with the named file.
	// It returns an error if the file doesn't exist or if it is not a supported
	// image.
	SetIcon(name string) error

	// SetBadge set the dock tile badge with the string representation of the
	// value.
	SetBadge(v interface{}) error
}
