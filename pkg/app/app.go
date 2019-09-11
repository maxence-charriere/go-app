package app

import (
	"reflect"

	"github.com/maxence-charriere/app/pkg/log"
)

var (
	// DefaultPath is the path to the component to be  loaded when no path is
	// specified.
	DefaultPath string

	// NotFoundPath is the path to the component to be  loaded when an non
	// imported component is requested.
	NotFoundPath = "/app.notfound"

	ui         = make(chan func(), 256)
	components = make(compoBuilder)
	msgs       = &messenger{
		callExec: func(f func(...interface{}), args ...interface{}) {
			go f(args...)
		},
		callOnUI: UI,
	}
)

// Import imports the given components into the app.
// Components must be imported in order the be used by the app package.
// This allows components to be created dynamically when they are found into
// markup.
func Import(c ...Compo) {
	for _, compo := range c {
		if err := components.imports(compo); err != nil {
			log.Error("importing component failed").
				T("reason", err).
				T("html tag", compoName(compo)).
				T("component type", reflect.TypeOf(compo)).
				T("fix", "rename component").
				Panic()
		}
	}
}

// Run runs the app with the loaded URL.
func Run() {
	run()
}

// Render renders the given component. It should be called whenever a component
// is modified. Render is always executed on the UI goroutine.
//
// It panics if called before Run.
func Render(c Compo) {
	UI(func() { render(c) })
}

// UI calls a function on the UI goroutine.
func UI(f func()) {
	ui <- f
}

// Reload reloads the current page.
func Reload() {
	UI(func() { reload() })
}

// Bind creates a binding between a message and the given component.
func Bind(msg string, c Compo) *Binding {
	return bind(msg, c)
}

// Emit emits a message that triggers the associated bindings.
func Emit(msg string, args ...interface{}) {
	go msgs.emit(msg, args...)
}

// WindowSize returns the window width and height.
func WindowSize() (w, h int) {
	return windowSize()
}
