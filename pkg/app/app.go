//go:generate go run gen/html.go
//go:generate go run gen/scripts.go
//go:generate go fmt

package app

import (
	"net/url"

	"github.com/maxence-charriere/go-app/v7/pkg/errors"
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
