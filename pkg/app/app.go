//go:generate go run gen/html.go
//go:generate go run gen/scripts.go
//go:generate go fmt

package app

import (
	"net/url"
	"strings"

	"github.com/maxence-charriere/go-app/v7/pkg/errors"
)

var (
	staticResourcesURL string
)

// Getenv retrieves the value of the environment variable named by the key. It
// returns the value, which will be empty if the variable is not present.
func Getenv(k string) string {
	return getenv(k)
}

// KeepBodyClean prevents third-party Javascript libraries to add nodes to the
// body element.
func KeepBodyClean() (close func()) {
	return keepBodyClean()
}

// Navigate navigates to the given URL.
func Navigate(rawurl string) {
	dispatch(func() {
		u, err := url.Parse(rawurl)
		if err != nil {
			panic(errors.New("navigating to page failed").
				Tag("url", rawurl).
				Wrap(err),
			)
		}

		if u.String() == Window().URL().String() {
			return
		}

		if err = navigate(u, true); err != nil {
			panic(errors.New("navigating to page failed").
				Tag("url", rawurl).
				Wrap(err),
			)
		}
	})
}

// NewContextMenu displays a context menu filled with the given menu items.
func NewContextMenu(menuItems ...MenuItemNode) {
	dispatch(func() {
		newContextMenu(menuItems...)
	})
}

// Reload reloads the current page.
func Reload() {
	dispatch(func() {
		reload()
	})
}

// Run starts the wasm app and displays the UI node associated with the
// requested URL path.
//
// It panics if Go architecture is not wasm.
func Run() {
	run()
}

// StaticResource makes a static resource path point to the right
// location whether the root directory is remote or not.
//
// Static resources are resources located in the web directory.
//
// This call is used internally to resolve paths within Cite, Data, Href, Src,
// and SrcSet. Paths already resolved are skipped.
func StaticResource(path string) string {
	if !strings.HasPrefix(path, "/web/") &&
		!strings.HasPrefix(path, "web/") {
		return path
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return staticResourcesURL + path
}

// Window returns the JavaScript "window" object.
func Window() BrowserWindow {
	return window
}
