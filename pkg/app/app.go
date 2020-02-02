//go:generate go run gen/html.go
//go:generate go run gen/scripts.go
//go:generate go fmt

package app

import (
	"net/url"
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
	NotFound ValueNode = &notFound{}

	routes = make(map[string]ValueNode)
)

// EventHandler represents a function that can handle HTML events.
type EventHandler func(src Value, e Event)

// Implement Raw(string)
// Implement Svg(string)

// Route binds the requested path to the given UI node.
func Route(path string, n ValueNode) {
	routes[path] = n
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
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}

	if err = navigate(u, true); err != nil {
		panic(err)
	}
}

// Reload reloads the current page.
func Reload() {
	reload()
}

// Window returns the JavaScript "window" object.
func Window() BrowserWindow {
	return window
}

// NewContextMenu displays a context menu filled with the given menu items.
func NewContextMenu(menuItems ...MenuItemNode) {
	newContextMenu(menuItems...)
}
