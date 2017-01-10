package app

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/murlokswarm/errors"
	"github.com/murlokswarm/log"
)

// ResourcePath represents the path of the app resource directory.
type ResourcePath string

// Path is a convenient method that cast a ResourcePath into a string.
func (r ResourcePath) Path() string {
	return string(r)
}

// Join joins any number of path elements into a single path, adding a
// Separator if necessary.
// It calls the Join function from package path/filepath with the resource
// location as the first element.
func (r ResourcePath) Join(elems ...string) string {
	elems = append([]string{r.Path()}, elems...)
	return filepath.Join(elems...)
}

// CSS returns a slice containing the css filenames of the css directory.
func (r ResourcePath) CSS() (css []string) {
	cssPath := r.Join("css")
	info, err := os.Stat(cssPath)
	if err != nil {
		log.Warnf("%v doesn't exists", cssPath)
		return
	}
	if !info.IsDir() {
		err := errors.Newf("%v is not a directory", cssPath)
		log.Error(err)
		return
	}

	files, _ := ioutil.ReadDir(cssPath)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if strings.HasSuffix(f.Name(), ".css") {
			css = append(css, filepath.Join("css", f.Name()))
		}
	}
	return
}

// JS returns a slice containing the js filenames of the js directory.
func (r ResourcePath) JS() (css []string) {
	cssPath := r.Join("js")
	info, err := os.Stat(cssPath)
	if err != nil {
		return
	}
	if !info.IsDir() {
		err := errors.Newf("%v is not a directory", cssPath)
		log.Error(err)
		return
	}

	files, _ := ioutil.ReadDir(cssPath)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if strings.HasSuffix(f.Name(), ".js") {
			css = append(css, filepath.Join("js", f.Name()))
		}
	}
	return
}

// Resources returns the path of the app resource directory.
func Resources() ResourcePath {
	return driver.Resources()
}

// IsSupportedImageExtension returns true if path has a supported image
// extensions.
// Supported extensions are jpg and png.
func IsSupportedImageExtension(path string) bool {
	ext := filepath.Ext(path)
	switch ext {
	case ".jpg", ".jpeg", ".png":
		return true

	default:
		return false
	}
}
