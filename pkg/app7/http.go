// +build !wasm

package app

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/maxence-charriere/go-app/v6/pkg/errors"
)

const (
	defaultThemeColor = "#2d2c2c"
)

// Handler is an HTTP handler that serves an HTML page that loads a Go wasm app
// and its resources.
type Handler struct {
	// The page authors.
	Author string

	// A placeholder background color for the application page to display before
	// its stylesheets are loaded.
	//
	// DEFAULT: #2d2c2c.
	BackgroundColor string

	// The path of the static resources that the browser is caching in order to
	// provide offline mode.
	//
	// Note that Icon, Styles and Scripts are already cached by default.
	//
	// Paths are relative to the root directory.
	CacheableResources []string

	// The page description.
	Description string

	// The environment variables that are passed to the progressive web app.
	//
	// Reserved keys:
	// - GOAPP_VERSION
	// - GOAPP_REMOTE_ROOT_DIR
	Env Environment

	// The icon that is used for the PWA, favicon, loading and default not
	// found component.
	Icon Icon

	// The page keywords.
	Keywords []string

	// The text displayed while loading a page.
	LoadingLabel string

	// The name of the web application as it is usually displayed to the user.
	Name string

	// Additional headers to be added in head element.
	RawHeaders []string

	// The URL or path of the root directory. The root directory is the location
	// where app resources are located.
	//
	// DEFAULT: ".".
	RootDir string

	// The paths or urls of the JavaScript files to use with the page.
	//
	// Paths are relative to the root directory.
	Scripts []string

	// The name of the web application displayed to the user when there is not
	// enough space to display Name.
	ShortName string

	// The paths or urls of the CSS files to use with the page.
	//
	// Paths are relative to the root directory.
	Styles []string

	// The theme color for the application. This affects how the OS displays the
	// app (e.g., PWA title var or Android's task switcher).
	//
	// DEFAULT: #2d2c2c.
	ThemeColor string

	// The page title.
	Title string

	// UseMinimalDefaultStyles makes app.css to adopt a minimal implementation
	// that does not style the default HTML elements.
	UseMinimalDefaultStyles bool

	// The version number. This is used in order to update the PWA application
	// in the browser. It must be set when deployed on a live system in order to
	// prevent recurring updates.
	//
	// Default: Auto-generated in order to trigger pwa update on a local
	// development system.
	Version string

	once             sync.Once
	etag             string
	appWasmPath      string
	hasRemoteRootDir bool
	page             bytes.Buffer
	manifestJSON     bytes.Buffer
	appJS            bytes.Buffer
	appWorkerJS      bytes.Buffer
	wasmExecJS       []byte
	appCSS           []byte
}

func (h *Handler) init() {
	h.initVersion()
	h.initRootDir()
	h.initStyles()
	h.initScripts()
	h.initIcon()
	h.initPWA()
	h.initPage()
	h.initWasmJS()
	h.initAppJS()
	h.initWorkerJS()
	h.initManifestJSON()
	h.initScripts()
	h.initAppCSS()
}

func (h *Handler) initVersion() {
	if h.Version == "" {
		t := time.Now().UTC().String()
		h.Version = fmt.Sprintf(`%x`, sha1.Sum([]byte(t)))
	}
	h.etag = `"` + h.Version + `"`
}

func (h *Handler) initRootDir() {
	rootDir := h.RootDir
	if rootDir == "" {
		rootDir = "."
	}
	rootDir = strings.TrimSuffix(rootDir, "/")
	rootDir = strings.TrimSuffix(rootDir, `\`)
	h.RootDir = rootDir

	h.hasRemoteRootDir = isRemoteLocation(rootDir)
	h.appWasmPath = filepath.Join(rootDir, "app.wasm")
}

func (h *Handler) initStyles() {
	for i, path := range h.Styles {
		h.Styles[i] = h.staticResource(path)
	}
}

func (h *Handler) initScripts() {
	for i, path := range h.Scripts {
		h.Scripts[i] = h.staticResource(path)
	}
}

func (h *Handler) initIcon() {
	if h.Icon.Default == "" {
		h.Icon.Default = "https://storage.googleapis.com/murlok-github/icon-192.png"
		h.Icon.Large = "https://storage.googleapis.com/murlok-github/icon-512.png"
	}

	if h.Icon.AppleTouch == "" {
		h.Icon.AppleTouch = h.Icon.Default
	}

	h.Icon.Default = h.staticResource(h.Icon.Default)
	h.Icon.Large = h.staticResource(h.Icon.Large)
	h.Icon.AppleTouch = h.staticResource(h.Icon.AppleTouch)
}

func (h *Handler) initPWA() {
	if h.Name == "" && h.ShortName == "" {
		h.Name = "App PWA"
	}
	if h.ShortName == "" {
		h.ShortName = h.Name
	}
	if h.Name == "" {
		h.Name = h.ShortName
	}

	if h.BackgroundColor == "" {
		h.BackgroundColor = defaultThemeColor
	}
	if h.ThemeColor == "" {
		h.ThemeColor = defaultThemeColor
	}

	if h.LoadingLabel == "" {
		h.LoadingLabel = "Loading"
	}
}

func (h *Handler) initPage() {
	h.page.WriteString("<!DOCTYPE html>\n")

	html := Html().Body(
		Head().Body(
			Meta().Charset("UTF-8"),
			Meta().
				HTTPEquiv("Content-Type").
				Content("text/html; charset=utf-8"),
			Meta().
				Name("author").
				Content(h.Author),
			Meta().
				Name("description").
				Content(h.Description),
			Meta().
				Name("keywords").
				Content(strings.Join(h.Keywords, ", ")),
			Meta().
				Name("viewport").
				Content("width=device-width, initial-scale=1, maximum-scale=1, user-scalable=0, viewport-fit=cover"),
			Title().
				Body(
					Text(h.Title),
				),
			Link().
				Rel("icon").
				Type("image/png").
				Href(h.Icon.Default),
			Link().
				Rel("apple-touch-icon").
				Href(h.Icon.AppleTouch),
			Link().
				Rel("manifest").
				Href("/manifest.json"),
			Link().
				Type("text/css").
				Rel("stylesheet").
				Href("/app.css"),
			Range(h.Styles).Slice(func(i int) UI {
				return Link().
					Type("text/css").
					Rel("stylesheet").
					Href(h.Styles[i])
			}),
			Range(h.RawHeaders).Slice(func(i int) UI {
				return Raw(h.RawHeaders[i])
			}),
		),
		Body().Body(
			Div().
				Class("app-wasm-layout").
				Body(
					Img().
						ID("app-wasm-loader-icon").
						Class("app-wasm-icon app-spin").
						Src(h.Icon.Default),
					P().
						ID("app-wasm-loader-label").
						Class("app-wasm-label").
						Body(Text(h.LoadingLabel)),
				),
			Div().ID("app-context-menu"),
			Script().Src("/wasm_exec.js"),
			Script().Src("/app.js"),
			Range(h.Scripts).Slice(func(i int) UI {
				return Script().
					Src(h.Scripts[i])
			}),
			Div().ID("app-end"),
		),
	)

	html.(writableNode).html(&h.page)
}

func (h *Handler) initWasmJS() {
	h.wasmExecJS = []byte(wasmExecJS)
}

func (h *Handler) initAppJS() {
	wasmURL := "/app.wasm"
	if h.hasRemoteRootDir {
		wasmURL = h.RootDir + wasmURL
	}

	if h.Env == nil {
		h.Env = make(map[string]string, 2)
	}
	h.Env["GOAPP_VERSION"] = h.Version

	if h.hasRemoteRootDir {
		h.Env["GOAPP_REMOTE_ROOT_DIR"] = h.RootDir
	}

	env, err := json.Marshal(h.Env)
	if err != nil {
		panic(errors.New("encoding pwa env failed").
			Tag("env", h.Env).
			Wrap(err),
		)
	}

	if err := template.
		Must(template.New("app.js").Parse(appJS)).
		Execute(&h.appJS, struct {
			Env  string
			Wasm string
		}{
			Env:  btos(env),
			Wasm: wasmURL,
		}); err != nil {
		panic(errors.New("initializing app.js failed").Wrap(err))
	}
}

func (h *Handler) initWorkerJS() {
	cacheableResources := make(map[string]struct{})
	cacheableResources["/wasm_exec.js"] = struct{}{}
	cacheableResources["/app.js"] = struct{}{}
	cacheableResources["/app.css"] = struct{}{}
	cacheableResources["/manifest.json"] = struct{}{}
	cacheableResources[h.Icon.Default] = struct{}{}
	cacheableResources[h.Icon.Large] = struct{}{}
	cacheableResources[h.Icon.AppleTouch] = struct{}{}
	cacheableResources["/"] = struct{}{}

	wasmPath := "/app.wasm"
	if h.hasRemoteRootDir {
		wasmPath = h.RootDir + wasmPath
	}
	cacheableResources[wasmPath] = struct{}{}

	cacheResources := func(res []string) {
		for _, r := range res {
			r = h.staticResource(r)
			cacheableResources[r] = struct{}{}
		}
	}
	cacheResources(h.Styles)
	cacheResources(h.Scripts)
	cacheResources(h.CacheableResources)

	if err := template.
		Must(template.New("app-worker.js").Parse(appWorkerJS)).
		Execute(&h.appWorkerJS, struct {
			Version          string
			ResourcesToCache map[string]struct{}
		}{
			Version:          h.Version,
			ResourcesToCache: cacheableResources,
		}); err != nil {
		panic(errors.New("initializing app-worker.js failed").Wrap(err))
	}
}

func (h *Handler) initManifestJSON() {
	if err := template.
		Must(template.New("manifest.json").Parse(manifestJSON)).
		Execute(&h.manifestJSON, struct {
			ShortName       string
			Name            string
			DefaultIcon     string
			LargeIcon       string
			BackgroundColor string
			ThemeColor      string
		}{
			ShortName:       h.ShortName,
			Name:            h.Name,
			DefaultIcon:     h.Icon.Default,
			LargeIcon:       h.Icon.Large,
			BackgroundColor: h.BackgroundColor,
			ThemeColor:      h.ThemeColor,
		}); err != nil {
		panic(errors.New("initializing manifest.json failed").Wrap(err))
	}
}

func (h *Handler) initAppCSS() {
	css := appCSS
	if h.UseMinimalDefaultStyles {
		css = appLightCSS
	}
	h.appCSS = []byte(css)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.once.Do(h.init)

	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("ETag", h.etag)

	etag := r.Header.Get("If-None-Match")
	if etag == h.etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	path := r.URL.Path

	switch path {
	case "/wasm_exec.js":
		h.serveWasmExecJS(w, r)
		return

	case "/app.js", "/goapp.js":
		h.serveAppJS(w, r)
		return

	case "/app-worker.js":
		h.serveAppWorkerJS(w, r)
		return

	case "/manifest.json":
		h.serveManifestJSON(w, r)
		return

	case "/app.css":
		h.serveAppCSS(w, r)
		return

	case "/app.wasm", "/goapp.wasm":
		http.ServeFile(w, r, h.appWasmPath)
		return
	}

	if strings.HasPrefix(path, "/web/") {
		filename := normalizeFilePath(path)
		filename = filepath.Join(h.RootDir, filename)
		if fi, err := os.Stat(filename); err == nil && !fi.IsDir() {
			http.ServeFile(w, r, filename)
			return
		}
	}

	h.servePage(w, r)
}

func (h *Handler) servePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Length", strconv.Itoa(h.page.Len()))
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(h.page.Bytes())
}

func (h *Handler) serveWasmExecJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Length", strconv.Itoa(len(h.wasmExecJS)))
	w.Header().Set("Content-Type", "application/javascript")
	w.WriteHeader(http.StatusOK)
	w.Write(h.wasmExecJS)
}

func (h *Handler) serveAppJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Length", strconv.Itoa(h.appJS.Len()))
	w.Header().Set("Content-Type", "application/javascript")
	w.WriteHeader(http.StatusOK)
	w.Write(h.appJS.Bytes())
}

func (h *Handler) serveAppWorkerJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Length", strconv.Itoa(h.appWorkerJS.Len()))
	w.Header().Set("Content-Type", "application/javascript")
	w.WriteHeader(http.StatusOK)
	w.Write(h.appWorkerJS.Bytes())
}

func (h *Handler) serveManifestJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Length", strconv.Itoa(h.manifestJSON.Len()))
	w.Header().Set("Content-Type", "application/manifest+json")
	w.WriteHeader(http.StatusOK)
	w.Write(h.manifestJSON.Bytes())
}

func (h *Handler) serveAppCSS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Length", strconv.Itoa(len(h.appCSS)))
	w.Header().Set("Content-Type", "text/css")
	w.WriteHeader(http.StatusOK)
	w.Write(h.appCSS)
}

func (h *Handler) staticResource(path string) string {
	if isRemoteLocation(path) {
		return path
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if h.hasRemoteRootDir {
		path = h.RootDir + path
	}
	return path
}

// Icon describes a square image that is used in various places such as
// application icon, favicon or loading icon.
type Icon struct {
	// The path or url to a square image/png file. It must have a side of 192px.
	//
	// Path is relative to the root directory.
	Default string

	// The path or url to larger square image/png file. It must have a side of
	// 512px.
	//
	// Path is relative to the root directory.
	Large string

	// The path or url to a square image/png file that is used for IOS/IPadOS
	// home screen icon. It must have a side of 192px.
	//
	// Path is relative to the root directory.
	//
	// DEFAULT: Icon.Default
	AppleTouch string
}

// Environment describes the environment variables to pass to the progressive
// web app.
type Environment map[string]string

func normalizeFilePath(path string) string {
	if runtime.GOOS == "windows" {
		return strings.ReplaceAll(path, "/", `\`)
	}
	return path
}

func isRemoteLocation(path string) bool {
	u, _ := url.Parse(path)
	return u.Scheme != ""
}
