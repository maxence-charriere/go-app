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
	// Defines the expected “background color” for the website. This value
	// repeats what is already available in the site’s CSS, but can be used by
	// browsers to draw the background color of a shortcut when the manifest is
	// available before the stylesheet has loaded. This creates a smooth
	// transition between launching the web application and loading the site's
	// content.
	BackgroundColor string

	// Defines the preferred display mode for the website.
	//
	// Default is Standalone.
	Display Display

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

// Display describes how progressive webapp is displayed.
type Display string

const (
	// FullScreen opens the web application without any browser UI and takes up
	// the entirety of the available display area.
	FullScreen Display = "fullscreen"

	// Standalone opens the web app to look and feel like a standalone native
	// app. The app runs in its own window, separate from the browser, and hides
	// standard browser UI elements like the URL bar, etc.
	Standalone = "standalone"

	// MinimalUI is similar to fullscreen, but provides the user with some means
	// to access a minimal set of UI elements for controlling navigation (i.e.,
	// back, forward, reload, etc).
	//
	// It is only supported by Chrome on mobile.
	MinimalUI = "minimal-ui"

	// Browser provides a standard browser experience.
	Browser = "browser"
)

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

	var pages http.Handler = &apphttp.PageHandler{
		Author:       h.Author,
		Description:  h.Description,
		Keywords:     h.Keywords,
		LoadingLabel: h.LoadingLabel,
		Name:         h.Name,
		WebDir:       webDir,
	}
	pages = apphttp.CacheHandler(pages, webDir)

	var manifest http.Handler = &apphttp.ManifestHandler{
		BackgroundColor: backgroundColor(h.ProgressiveApp.BackgroundColor),
		Display:         string(h.ProgressiveApp.Display),
		Name:            h.Name,
		Orientation:     orientation(h.ProgressiveApp.LanscapeMode),
		ShortName:       shortName(h.Name, h.ProgressiveApp.ShortName),
		Scope:           entryPoint(h.ProgressiveApp.Scope),
		StartURL:        entryPoint(h.ProgressiveApp.StartURL),
		ThemeColor:      backgroundColor(h.ProgressiveApp.ThemeColor),
	}
	manifest = apphttp.CacheHandler(manifest, webDir)

	h.Handler = &apphttp.RouteHandler{
		Files:    files,
		Manifest: manifest,
		Pages:    pages,
		WebDir:   webDir,
	}
}

func display(display string) string {
	if display == "" {
		return "standalone"
	}
	return display
}

func shortName(name, shortName string) string {
	if shortName != "" {
		return shortName
	}

	shortName = strings.ReplaceAll(name, " ", "")
	shortName = strings.ReplaceAll(shortName, "\n", "")
	shortName = strings.ReplaceAll(shortName, "\t", "")
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

func backgroundColor(color string) string {
	if color == "" {
		return "#21252b"
	}
	return color
}
