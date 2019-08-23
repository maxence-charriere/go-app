// Package app is a package to build GUI apps with Go, HTML and CSS.
package app

import (
	"fmt"
	"log"
	"strings"

	"github.com/maxence-charriere/app/internal/maestro"
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

	// Logger is a function that formats using the default formats for its
	// operands and logs the resulting string.
	// It is used by Log, Logf, Panic and Panicf to generate logs.
	Logger = log.Printf

	maestre    = maestro.NewMaestro()
	components = newCompoBuilder()
	messages   = newMsgRegistry()
	ui         = make(chan func(), 4096)
	events     = newEventRegistry(ui)
	whenDebug  func(func())
)

func init() {
	EnableDebug(false)
}

// Emit emits the event with the given arguments.
func Emit(e Event, args ...interface{}) {
	events.Emit(e, args...)
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

// Handle handles the message for the given key.
func Handle(key string, h MsgHandler) {
	messages.handle(key, h)
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

		if err := maestre.Import(compo); err != nil {
			panic(err)
		}
	}
}

// Log formats using the default formats for its operands and logs the resulting
// string.
// Spaces are always added between operands and a newline is appended.
func Log(a ...interface{}) {
	format := ""

	for range a {
		format += "%v "
	}

	format = format[:len(format)-1]
	Logger(format, a...)
}

// Logf formats according to a format specifier and logs the resulting string.
func Logf(format string, a ...interface{}) {
	Logger(format, a...)
}

// Navigate navigates to the given URL.
func Navigate(url string) {
	navigate(url)
}

// NewContextMenu displays a context menu filled with the given menu items.
//
// Context menu requires an app.contextmenu component in the loaded page.
// 	func (c *Compo) Render() string {
// 		return `
// 	<div>
// 		<!-- ... -->
// 		<app.contextmenu>
// 	</div>
// 		`
// 	}
func NewContextMenu(items ...MenuItem) {
	Emit("app.NewContextMenu", items)
}

// NewMsg creates a message.
func NewMsg(key string) Msg {
	return &msg{key: key}
}

// NewSubscriber creates an event subscriber to return when implementing the
// app.EventSubscriber interface.
func NewSubscriber() *Subscriber {
	return &Subscriber{events: events}
}

// Panic is equivalent to Log() followed by a call to panic().
func Panic(a ...interface{}) {
	Log(a...)
	panic(strings.TrimSpace(fmt.Sprintln(a...)))
}

// Panicf is equivalent to Logf() followed by a call to panic().
func Panicf(format string, a ...interface{}) {
	Logf(format, a...)
	panic(fmt.Sprintf(format, a...))
}

// Path returns the path to the given component.
func Path(c Compo) string {
	return "/" + compoName(c)
}

// Post posts the given messages.
// Messages are handled in another goroutine.
func Post(msgs ...Msg) {
	messages.post(msgs...)
}

// Reload reloads the current page.
func Reload() {
	reload()
}

// Render renders the given component.
// It should be called whenever a component is modified.
//
// It panics if called before Run.
func Render(c Compo) {
	WhenDebug(func() {
		Logf("rendering %T")
	})

	if err := render(c); err != nil {
		Log(err)
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
