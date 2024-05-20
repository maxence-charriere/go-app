package app

// The interface that describes a library that contains custom components.
type Library interface {
	// Returns the styles and its path. The styles must be a standard
	// CSS code.
	Styles() (path, styles string)
}
