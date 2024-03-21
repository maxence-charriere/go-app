//go:generate go run gen/html.go
//go:generate go run gen/scripts.go
//go:generate go fmt

// Package app is a package to build progressive web apps (PWA) with Go
// programming language and WebAssembly.
// It uses a declarative syntax that allows creating and dealing with HTML
// elements only by using Go, and without writing any HTML markup.
// The package also provides an http.Handler ready to serve all the required
// resources to run Go-based progressive web apps.
package app

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"runtime"
)

const (
	// IsClient reports whether the code is running as a client in the
	// WebAssembly binary (app.wasm).
	IsClient = runtime.GOARCH == "wasm" && runtime.GOOS == "js"

	// IsServer reports whether the code is running on a server for
	// pre-rendering purposes.
	IsServer = runtime.GOARCH != "wasm" || runtime.GOOS != "js"
)

var (
	routes = makeRouter()
	window = newBrowserWindow()
)

// Getenv retrieves the value of the environment variable named by the key. It
// returns the value, which will be empty if the variable is not present.
func Getenv(k string) string {
	if IsServer || !Window().Get("goappGetenv").Truthy() {
		return os.Getenv(k)
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
	if IsServer {
		return func() {}
	}

	release := Window().Call("goappKeepBodyClean")
	return func() {
		release.Invoke()
	}
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
//
//	 func main() {
//			// Define app routes.
//			app.Route("/", myComponent{})
//			app.Route("/other-page", myOtherComponent{})
//
//			// Run the application when on a web browser (only executed on client side).
//			app.RunWhenOnBrowser()
//
//			// Launch the server that serves the app (only executed on server side):
//			http.Handle("/", &app.Handler{Name: "My app"})
//			http.ListenAndServe(":8080", nil)
//	 }
func RunWhenOnBrowser() {
	if IsServer {
		return
	}

	defer func() {
		err := recover()
		displayLoadError(err)
		panic(err)
	}()

	resolveURL := clientResourceResolver(Getenv("GOAPP_STATIC_RESOURCES_URL"))
	originPage := makeRequestPage(Window().URL(), resolveURL)

	engine := newEngine(context.Background(),
		&routes,
		resolveURL,
		&originPage,
		actionHandlers,
	)

	engine.Navigate(window.URL(), false)
	engine.Start(120)
}

func displayLoadError(err any) {
	loadingLabel := Window().
		Get("document").
		Call("getElementById", "app-wasm-loader-label")
	if !loadingLabel.Truthy() {
		return
	}
	loadingLabel.setInnerText(fmt.Sprint(err))
}

// Route associates a given path with a function that generates a new Composer
// component. When a user navigates to the specified path, the function
// newComponent is invoked to create and mount the associated component.
//
// Example:
//
//	Route("/home", func() Composer {
//	    return NewHomeComponent()
//	})
func Route(path string, newComponent func() Composer) {
	routes.route(path, newComponent)
}

// RouteWithRegexp associates a URL path pattern with a function that generates
// a new Composer component. When a user navigates to a URL path that matches
// the given regular expression pattern, the function newComponent is invoked to
// create and mount the associated component.
//
// Example:
//
//	RouteWithRegexp("^/users/[0-9]+$", func() Composer {
//	    return NewUserComponent()
//	})
func RouteWithRegexp(pattern string, newComponent func() Composer) {
	routes.routeWithRegexp(pattern, newComponent)
}

// NewZeroComponentFactory returns a function that, when invoked, creates and
// returns a new instance of the same type as the provided component. The new
// instance is initialized with zero values for all its fields.
//
// The function uses reflection to determine the type of the provided Composer
// and to create new instances of that type.
//
// Example:
//
//	componentFunc := NewZeroComponentFactory(MyComponent{})
//	newComponent := componentFunc()
func NewZeroComponentFactory(c Composer) func() Composer {
	componentType := reflect.TypeOf(c)

	return func() Composer {
		return reflect.New(componentType.Elem()).Interface().(Composer)
	}
}

// TryUpdate attempts to update the application in the browser. On success, it
// notifies components implementing the AppUpdater interface that an update is
// ready.
func TryUpdate() {
	if tryUpdate := Window().Get("goappTryUpdate"); IsClient && tryUpdate.Truthy() {
		tryUpdate.Invoke()
	}
}
