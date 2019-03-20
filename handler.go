// +build !wasm

package app

import (
	"net/http"
	"strings"
	"sync"

	apphttp "github.com/maxence-charriere/app/internal/http"
)

// ProgressiveAppConfig represents the configuration used to describe a
// progressive app.
type ProgressiveAppConfig struct {
	// Enforces landscape mode.
	LanscapeMode bool

	// Provides a short human-readable name for the application. This is
	// intended for when there is insufficient space to display the full name of
	// the web application, like device homescreens.
	//
	// Default is the app name where space are replaces by '-'.
	ShortName string

	// Defines the navigation scope of this website's context. This restricts
	// what web pages can be viewed while the manifest is applied. If the user
	// navigates outside the scope, it returns to a normal web page inside a
	// browser tab/window.
	//
	// Default is "/".
	Scope string

	// The URL that loads when a user launches the application (e.g. when added
	// to home screen), typically the index.
	// Default is "/".
	StartURL string

	// Defines the default theme color for an application. This sometimes
	// affects how the OS displays the site (e.g., on Android's task switcher,
	// the theme color surrounds the site).
	ThemeColor string
}

// Handler is a http handler that serves UI components created with this
// package.
type Handler struct {
	http.Handler

	// The app author.
	Author string

	// The app description.
	Description string

	// The app keywords.
	Keywords []string

	// The text displayed while loading a page.
	LoadingLabel string

	// The app name.
	Name string

	// The progressive app mode configuration.
	ProgressiveApp ProgressiveAppConfig

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

	themeColor := h.ProgressiveApp.ThemeColor
	if themeColor == "" {
		themeColor = "#21252b"
	}

	var pages http.Handler = &apphttp.PageHandler{
		Author:       h.Author,
		Description:  h.Description,
		Keywords:     h.Keywords,
		LoadingLabel: h.LoadingLabel,
		Name:         h.Name,
		ThemeColor:   themeColor,
		WebDir:       webDir,
	}
	pages = apphttp.CacheHandler(pages, webDir)

	var manifest http.Handler = &apphttp.ManifestHandler{
		BackgroundColor: themeColor,
		Name:            h.Name,
		Orientation:     orientation(h.ProgressiveApp.LanscapeMode),
		ShortName:       shortName(h.Name, h.ProgressiveApp.ShortName),
		Scope:           entryPoint(h.ProgressiveApp.Scope),
		StartURL:        entryPoint(h.ProgressiveApp.StartURL),
		ThemeColor:      themeColor,
	}
	manifest = apphttp.CacheHandler(manifest, webDir)

	h.Handler = &apphttp.RouteHandler{
		Files:    files,
		Manifest: manifest,
		Pages:    pages,
		WebDir:   webDir,
	}
}

func shortName(name, shortName string) string {
	if shortName != "" {
		return shortName
	}

	shortName = strings.Replace(name, " ", "", -1)
	shortName = strings.Replace(shortName, "\n", "", -1)
	shortName = strings.Replace(shortName, "\t", "", -1)
	return shortName
}

func orientation(landscapeMode bool) string {
	if landscapeMode {
		return "landscape"
	}
	return "any"
}

func entryPoint(entryPoint string) string {
	if entryPoint == "" {
		return "/"
	}
	return entryPoint
}
