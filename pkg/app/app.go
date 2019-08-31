// Package app is a package to build GUI apps with Go, HTML and CSS.
package app

import (
	"reflect"

	"github.com/maxence-charriere/app/internal/maestro"
	"github.com/maxence-charriere/app/pkg/log"
	"github.com/pkg/errors"
)

var (
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
	ui         = make(chan func(), 256)
)

// Import imports the given components into the app.
// Components must be imported in order the be used by the app package.
// This allows components to be created dynamically when they are found into
// markup.
func Import(c ...Compo) {
	for _, compo := range c {
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
	return "/" + maestro.CompoName(c)
}

// Render renders the given component.
// It should be called whenever a component is modified.
// Render is always excecuted on the UI goroutine.
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
