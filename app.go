// Package app is a package to build GUI apps with Go, HTML and CSS.
package app

import (
	"reflect"

	"github.com/maxence-charriere/app/internal/maestro"
	"github.com/maxence-charriere/app/pkg/log"
	"github.com/pkg/errors"
)

var (
	// ErrCompoNotMounted describes an error that reports whether a component
	// is mounted.
	ErrCompoNotMounted = errors.New("component not mounted")

	// ErrElemNotSet describes an error that reports if an element is set.
	ErrElemNotSet = errors.New("element not set")

	// ErrNotSupported describes an error that occurs when an unsupported
	// feature is used.
	ErrNotSupported = errors.New("not supported")

	// ErrNoWasm describes an error that occurs when Run or Render are called
	// in a non wasm environment.
	ErrNoWasm = errors.New("go architecture is not wasm")

	// DefaultPath is the path to the component to be  loaded when no path is
	// specified.
	DefaultPath string

	// NotFoundPath is the path to the component to be  loaded when an non
	// imported component is requested.
	NotFoundPath = "/app.notfound"

	components = make(maestro.CompoBuilder)
	ui         = make(chan func(), 4096)
	whenDebug  func(func())
)

func init() {
	EnableDebug(false)
}

// EnableDebug is a function that set whether debug mode is enabled.
func EnableDebug(v bool) {
	whenDebug = func(f func()) {}

	if v {
		whenDebug = func(f func()) {
			f()
		}
	}
}

// Import imports the given components into the app.
// Components must be imported in order the be used by the app package.
// This allows components to be created dynamically when they are found into
// markup.
func Import(c ...Compo) {
	for _, compo := range c {
		// if _, err := components.register(compo); err != nil {
		// 	Panicf("import component failed: %s", err)
		// }

		if err := components.Import(compo); err != nil {
			panic(err)
		}
	}
}

// Navigate navigates to the given URL.
func Navigate(url string) {
	navigate(url)
}

// Path returns the path to the given component.
func Path(c Compo) string {
	return "/" + compoName(c)
}

// Render renders the given component.
// It should be called whenever a component is modified.
//
// It panics if called before Run.
func Render(c Compo) {
	if err := render(c); err != nil {
		log.Error("rendering component failed").
			T("reason", err).
			T("component", reflect.TypeOf(c))
	}
}

// Run runs the app with the loaded URL.
func Run() error {
	return run()
}

// UI calls a function on the UI goroutine.
func UI(f func()) {
	ui <- f
}

// WhenDebug execute the given function when debug mode is enabled.
func WhenDebug(f func()) {
	whenDebug(f)
}
