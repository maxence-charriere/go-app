package app

import (
	"path/filepath"
)

// ResourcePath represents the path of the app resource directory.
type ResourcePath string

// Path is a convenient method that cast a ResourcePath into a string.
func (m ResourcePath) Path() string {
	return string(m)
}

// Join joins any number of path elements into a single path, adding a
// Separator if necessary.
// It calls the Join function from package path/filepath with the resource
// location as the first element.
func (m ResourcePath) Join(elems ...string) string {
	elems = append([]string{m.Path()}, elems...)
	return filepath.Join(elems...)
}

// Resources returns the path of the app resource directory.
func Resources() ResourcePath {
	return driver.Resources()
}
