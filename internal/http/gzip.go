package http

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// GzipHandler returns a decorated version of the given handler that serves
// available gzipped static resources.
func GzipHandler(h http.Handler, webDir string) http.Handler {
	return &gzipHandler{
		Handler: h,
		webDir:  webDir,
	}
}

type gzipHandler struct {
	http.Handler
	once    sync.Once
	version string
	webDir  string
}

func (h *gzipHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.once.Do(h.init)

	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		h.Handler.ServeHTTP(w, r)
		return
	}

	filename := strings.TrimPrefix(r.URL.Path, "/")
	filename = filepath.Join(h.webDir, filename)
	mimeType := mime.TypeByExtension(filepath.Ext(filename))

	gzipname := filename
	if h.version != "" {
		gzipname += "." + h.version
	}
	gzipname += ".gz"

	fmt.Println(gzipname)

	if _, err := os.Stat(gzipname); err != nil {
		h.Handler.ServeHTTP(w, r)
		return
	}

	r = r.WithContext(r.Context())
	r.URL.Path += ".gz"
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", mimeType)
	h.Handler.ServeHTTP(w, r)
}

func (h *gzipHandler) init() {
	h.version = GetEtag(h.webDir)
}
