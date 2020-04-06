//go:generate go run gen/html.go
//go:generate go run gen/scripts.go
//go:generate go fmt

package app

import (
	"net/url"
	"strings"

	"github.com/maxence-charriere/go-app/v6/pkg/log"
)

var (
	// LocalStorage is a storage that uses the browser local storage associated
	// to the document origin. Data stored are encrypted and has no expiration
	// time.
	LocalStorage BrowserStorage

	// SessionStorage is a storage that uses the browser session storage
	// associated to the document origin. Data stored are encrypted and expire
	// when the page session ends.
	SessionStorage BrowserStorage

	// NotFound is the ui element that is displayed when a request is not
	// routed.
	NotFound UI = &notFound{}

	remoteRootDir string
	routes        router
	dispatcher    Dispatcher = Dispatch
	uiChan                   = make(chan func(), 256)
)

// Dispatcher is a function that executes the given function on the goroutine dedicated to UI.
type Dispatcher func(func())

// EventHandler represents a function that can handle HTML events.
type EventHandler func(src Value, e Event)

// Route binds the requested path to the given UI node.
func Route(path string, node UI) {
	routes.route(path, node)
}

// RouteWithRegexp binds the regular expression pattern to the given UI node.
// Patterns use the Go standard regexp format.
func RouteWithRegexp(pattern string, node UI) {
	routes.routeWithRegexp(pattern, node)
}

// Run starts the wasm app and displays the UI node associated with the
// requested URL path.
//
// It panics if Go architecture is not wasm.
func Run() {
	run()
}

// Navigate navigates to the given URL.
func Navigate(rawurl string) {
	dispatcher(func() {
		u, err := url.Parse(rawurl)
		if err != nil {
			log.Error("navigating to page failed").
				T("url", rawurl).
				T("error", err).
				Panic()
		}

		if u.String() == Window().URL().String() {
			return
		}

		if err = navigate(u, true); err != nil {
			log.Error("navigating to page failed").
				T("url", u).
				T("error", err).
				Panic()
		}
	})
}

// Reload reloads the current page.
func Reload() {
	dispatcher(func() {
		reload()
	})
}

// Window returns the JavaScript "window" object.
func Window() BrowserWindow {
	return window
}

// NewContextMenu displays a context menu filled with the given menu items.
func NewContextMenu(menuItems ...MenuItemNode) {
	dispatcher(func() {
		newContextMenu(menuItems...)
	})
}

// Dispatch executes the given function on the UI goroutine.
func Dispatch(f func()) {
	uiChan <- f
}

// ResolveStaticResourcePath makes a static resource path point to the right
// location whether the root directory is remote or not.
//
// Static resources are resources located in the web directory.
//
// This call is used internally to resolve paths within Cite, Data, Href, Src,
// and SrcSet. Paths already resolved are skipped.
func ResolveStaticResourcePath(path string) string {
	if !strings.HasPrefix(path, "/web/") &&
		!strings.HasPrefix(path, "web/") ||
		remoteRootDir == "" {
		return path
	}

	path = strings.TrimPrefix(path, "/")

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return remoteRootDir + path
}

// Getenv retrieves the value of the environment variable named by the key. It
// returns the value, which will be empty if the variable is not present.
func Getenv(k string) string {
	return getenv(k)
}
