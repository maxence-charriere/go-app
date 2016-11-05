package app

import (
	"path/filepath"
)

// ResourceLocation represents the path of the app resource directory.
type ResourceLocation string

// Path is a convenient method that cast a ResourceLocation into a string.
func (m ResourceLocation) Path() string {
	return string(m)
}

// Join joins any number of path elements into a single path, adding a
// Separator if necessary.
// It calls the Join function from package path/filepath with the resource
// location as the first element.
func (m ResourceLocation) Join(elems ...string) string {
	elems = append([]string{m.Path()}, elems...)
	return filepath.Join(elems...)
}

// Resources returns the path of the app resource directory.
func Resources() ResourceLocation {
	return driver.Resources()
}
