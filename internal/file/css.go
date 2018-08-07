package file

import (
	"os"
	"path/filepath"
)

// CSS returns a list that contains the path of the css files located
// in given directory path.
func CSS(dirname string) []string {
	var css []string

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if ext := filepath.Ext(path); ext != ".css" {
			return nil
		}

		css = append(css, path)
		return nil
	}

	filepath.Walk(dirname, walker)
	return css
}
