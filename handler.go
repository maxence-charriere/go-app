// +build !wasm

package app

import (
	"net/http"
	"sync"

	apphttp "github.com/maxence-charriere/app/internal/http"
)

// Handler is a http handler that serves UI components created with this
// package.
type Handler struct {
	http.Handler

	// The app author.
	Author string

	// The app description.
	Description string

	// The path of the icon relative to the web directory.
	Icon string

	// The app keywords.
	Keywords []string

	// The text displayed while loading a page.
	LoadingLabel string

	// The app name.
	Name string

	// The he path of the web directory. Default is web.
	WebDir string

	once sync.Once
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.once.Do(h.init)
	h.Handler.ServeHTTP(w, r)
}

func (h *Handler) init() {
	webDir := h.WebDir
	if webDir == "" {
		webDir = "web"
	}

	files := apphttp.FileHandler(webDir)
	files = apphttp.GzipHandler(files, webDir)
	files = apphttp.CacheHandler(files, webDir)

	var pages http.Handler = &apphttp.PageHandler{
		Author:       h.Author,
		Description:  h.Description,
		Icon:         h.Icon,
		Keywords:     h.Keywords,
		LoadingLabel: h.LoadingLabel,
		Name:         h.Name,
		WebDir:       webDir,
	}
	pages = apphttp.CacheHandler(pages, webDir)

	h.Handler = &apphttp.RouteHandler{
		Files:  files,
		Pages:  pages,
		WebDir: webDir,
	}
}
