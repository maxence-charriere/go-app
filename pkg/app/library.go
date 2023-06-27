package app

// The interface that describes a library that contains custom components.
type Library interface {
	// Returns the script and its path. The script must be a standard
	// CSS file.
	Script() (path, script string)
}
