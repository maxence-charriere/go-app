//go:generate go run gen/html.go
//go:generate go run gen/scripts.go
//go:generate go fmt

// Package app is a package to build progressive web apps (PWA) with Go
// programming language and WebAssembly.
// It uses a declarative syntax that allows creating and dealing with HTML
// elements only by using Go, and without writing any HTML markup.
// The package also provides an http.handler ready to serve all the required
// resources to run Go-based progressive web apps.
package app

import (
	"os"
	"runtime"
	"strings"
)

const (
	// IsAppWASM reports whether the code is running in the WebAssembly binary
	// (app.wasm).
	IsAppWASM = runtime.GOARCH == "wasm" && runtime.GOOS == "js"
)

var (
	staticResourcesURL string
	appUpdateAvailable bool
)

// Getenv retrieves the value of the environment variable named by the key. It
// returns the value, which will be empty if the variable is not present.
func Getenv(k string) string {
	if !IsAppWASM {
		os.Getenv(k)
	}

	env := Window().Call("goappGetenv", k)
	if !env.Truthy() {
		return ""
	}
	return env.String()
}

// KeepBodyClean prevents third-party Javascript libraries to add nodes to the
// body element.
func KeepBodyClean() (close func()) {
	if !IsAppWASM {
		return func() {}
	}

	release := Window().Call("goappKeepBodyClean")
	return func() {
		release.Invoke()
	}
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

// RunWhenOnBrowser starts the app, displaying the component associated with the
// current URL path.
//
// This call is skipped when the program is not run on a web browser. This
// allows writing client and server-side code without separation or
// pre-compilation flags.
//
// Eg:
//  func main() {
// 		// Define app routes.
// 		app.Route("/", myComponent{})
// 		app.Route("/other-page", myOtherComponent{})
//
// 		// Run the application when on a web browser (only executed on client side).
// 		app.RunWhenOnBrowser()
//
// 		// Launch the server that serves the app (only executed on server side):
// 		http.Handle("/", &app.Handler{Name: "My app"})
// 		http.ListenAndServe(":8080", nil)
//  }
func RunWhenOnBrowser() {
	if !IsAppWASM {
		return
	}

	panic("not implemented")
}
