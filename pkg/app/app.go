// +build !wasm

package app

// Import imports the given components into the app.
// Components must be imported in order the be used by the app package.
// This allows components to be created dynamically when they are found into
// markup.
func Import(c ...Compo) {
}

// Run runs the app with the loaded URL.
func Run() {}

// Render renders the given component. It should be called whenever a component
// is modified. Render is always excecuted on the UI goroutine.
//
// It panics if called before Run.
func Render(c Compo) {}

// UI calls a function on the UI goroutine.
func UI(f func()) {}
