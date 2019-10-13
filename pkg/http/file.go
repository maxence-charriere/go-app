package http

import (
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Files returns a handler that serves files from the given directory.
func Files(dir string) RouteHandler {
	return files{
		handler:   http.FileServer(http.Dir(dir)),
		dir:       dir,
		windowsOS: runtime.GOOS == "windows",
	}
}

type files struct {
	handler   http.Handler
	dir       string
	windowsOS bool
}

func (f files) CanHandle(r *http.Request) bool {
	filename := r.URL.Path
	if f.windowsOS {
		filename = strings.ReplaceAll(filename, "/", `\`)
	}
	filename = filepath.Join(f.dir, filename)

	fi, err := os.Stat(filename)
	return err == nil && !fi.IsDir()
}

func (f files) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.handler.ServeHTTP(w, r)
}
