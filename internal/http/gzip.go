package http

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	webDir string
}

func (h *gzipHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		h.Handler.ServeHTTP(w, r)
		return
	}

	filename := strings.TrimPrefix(r.URL.Path, "/")
	filename = filepath.Join(h.webDir, filename)
	mimeType := mime.TypeByExtension(filepath.Ext(filename))
	gzipname := filename + ".gz"

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
