//go:generate go run gen/html.go
//go:generate go run gen/scripts.go
//go:generate go fmt

package app

import (
	"net/url"

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

	routes router
	uiChan = make(chan func(), 256)
)

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
	Dispatch(func() {
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
	Dispatch(func() {
		reload()
	})
}

// Window returns the JavaScript "window" object.
func Window() BrowserWindow {
	return window
}

// NewContextMenu displays a context menu filled with the given menu items.
func NewContextMenu(menuItems ...MenuItemNode) {
	Dispatch(func() {
		newContextMenu(menuItems...)
	})
}

// Dispatch executes the given function on the UI goroutine.
func Dispatch(f func()) {
	uiChan <- f
}
