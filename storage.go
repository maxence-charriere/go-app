package app

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/murlokswarm/log"
)

// FileHaveExtension returns a boolean indicating whether or not name have an
// extension defined in exts.
func FileHaveExtension(name string, exts ...string) bool {
	ext := filepath.Ext(name)
	for _, e := range exts {
		if ext == e {
			return true
		}
	}
	return false
}

// FileIsSupportedIcon returns a boolean indicating whether or not name is a
// supported icon.
func FileIsSupportedIcon(name string) bool {
	return FileHaveExtension(name, ".jpg", ".jpeg", ".png")
}

// GetFilenamesFromDir returns the filenames within dirname.
// names are not prefixed by dirname.
func GetFilenamesFromDir(dirname string, extension ...string) (names []string) {
	info, err := os.Stat(dirname)
	if err != nil {
		return
	}
	if !info.IsDir() {
		log.Errorf("%v is not a directory", dirname)
		return
	}

	files, _ := ioutil.ReadDir(dirname)
	for _, f := range files {
		if f.IsDir() {
			subdirname := filepath.Join(dirname, f.Name())
			subfilenames := GetFilenamesFromDir(subdirname, extension...)

			for _, n := range subfilenames {
				subfilename := filepath.Join(f.Name(), n)
				names = append(names, subfilename)
			}
			continue
		}
		if FileHaveExtension(f.Name(), extension...) {
			names = append(names, f.Name())
		}
	}
	return
}
