package app

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Storer describes the directory locations to use during app lifecycle.
type Storer interface {
	// Resources returns resources directory filename.
	// Represents the root location where files related to the operation of
	// the app should be located.
	Resources() string

	// CSS returns the location where the .css files should be located.
	CSS() string

	// JS returns the location where the .js files should be located.
	JS() string

	// Storage returns the root location where common files should be stored.
	// eg db, cache, downloaded content.
	Default() string
}

// IsSupportedExtension returns a boolean indicating whether or not extensions
// contains the extension of name.
// extensions must contain the dot prefix. eg ".png".
func IsSupportedExtension(name string, extensions ...string) bool {
	ext := filepath.Ext(name)
	for _, e := range extensions {
		if ext == e {
			return true
		}
	}
	return false
}

// IsSupportedImageExtension returns a boolean indicating whether or not name
// is a .jpg, .jpeg or .png.
func IsSupportedImageExtension(name string) bool {
	return IsSupportedExtension(name, ".jpg", ".jpeg", ".png")
}

// GetFilenamesWithExtensionsFromDir returns the filenames of files within
// dirname. names are not prefixed with dirname.
func GetFilenamesWithExtensionsFromDir(dirname string, extension ...string) (names []string, err error) {
	info, err := os.Stat(dirname)
	if err != nil {
		return
	}
	if !info.IsDir() {
		err = errors.Errorf("%v is not a directory", dirname)
		return
	}

	files, _ := ioutil.ReadDir(dirname)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if IsSupportedExtension(f.Name(), extension...) {
			names = append(names, f.Name())
		}
	}
	return
}
