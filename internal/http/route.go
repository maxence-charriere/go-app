package http

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// RouteHandler is a handler that routes requests to the appropriate handler.
type RouteHandler struct {
	Files    http.Handler
	Pages    http.Handler
	Manifest http.Handler
	WebDir   string
}

func (h *RouteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/")
	filename = filepath.Join(h.WebDir, filename)

	if info, err := os.Stat(filename); err == nil && !info.IsDir() {
		h.Files.ServeHTTP(w, r)
		return
	}

	if r.URL.Path == "/manifest.json" {
		h.Manifest.ServeHTTP(w, r)
		return
	}

	h.Pages.ServeHTTP(w, r)
}
