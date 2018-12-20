// Package app is a package to build GUI apps with Go, HTML and CSS.
package app

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

var (
	// ErrNotSupported describes an error that occurs when an unsupported
	// feature is used.
	ErrNotSupported = errors.New("not supported")

	// ErrElemNotSet describes an error that reports if an element is set.
	ErrElemNotSet = errors.New("element not set")

	// ErrCompoNotMounted describes an error that reports whether a component
	// is mounted.
	ErrCompoNotMounted = errors.New("component not mounted")

	// Logger is a function that formats using the default formats for its
	// operands and logs the resulting string.
	// It is used by Log, Logf, Panic and Panicf to generate logs.
	Logger func(format string, a ...interface{})

	// Kind describes the app kind (desktop|mobile|web).
	Kind string

	driver    Driver
	ui        = make(chan func(), 4096)
	factory   = NewFactory()
	events    = NewEventRegistry(ui)
	messages  = newMsgRegistry()
	whenDebug func(func())
)

const (
	// Running is the event emitted when the app starts to run.
	Running Event = "app.running"

	// Reopened is the event emitted when the app is reopened.
	Reopened Event = "app.reopened"

	// Focused is the event emitted when the app gets focus.
	Focused Event = "app.focused"

	// Blurred is the event emitted when the app loses focus.
	Blurred Event = "app.blurred"
)

func init() {
	EnableDebug(false)
}

// Import imports the given components into the app.
// Components must be imported in order the be used by the app package.
// This allows components to be created dynamically when they are found into
// markup.
func Import(c ...Compo) {
	for _, compo := range c {
		if _, err := factory.RegisterCompo(compo); err != nil {
			Panicf("import component failed: %s", err)
		}
	}
}

// Run runs the app with the given driver as backend.
func Run(d Driver, addons ...Addon) error {
	if len(addons) == 0 {
		addons = append(addons, Logs())
	}

	for _, addon := range addons {
		d = addon(d)
	}

	driver = d

	return driver.Run(DriverConfig{
		UI:      ui,
		Factory: factory,
		Events:  events,
	})
}

// CurrentDriver returns the current driver.
func CurrentDriver() Driver {
	return driver
}

// Name returns the application name.
//
// It panics if called before Run.
func Name() string {
	return driver.AppName()
}

// Resources returns the given path prefixed by the resources directory
// location.
// Resources should be used only for read only operations.
//
// It panics if called before Run.
func Resources(path ...string) string {
	return driver.Resources(path...)
}

// Storage returns the given path prefixed by the storage directory
// location.
//
// It panics if called before Run.
func Storage(path ...string) string {
	return driver.Storage(path...)
}

// Render renders the given component.
// It should be called when the display of component c have to be updated.
//
// It panics if called before Run.
func Render(c Compo) {
	driver.UI(func() {
		driver.Render(c)
	})
}

// ElemByCompo returns the element where the given component is mounted.
//
// It panics if called before Run.
func ElemByCompo(c Compo) Elem {
	return driver.ElemByCompo(c)
}

// NewWindow creates and displays the window described by the given
// configuration.
//
// It panics if called before Run.
func NewWindow(c WindowConfig) Window {
	return driver.NewWindow(c)
}

// NewPage creates the page described by the given configuration.
//
// It panics if called before Run.
func NewPage(c PageConfig) Page {
	return driver.NewPage(c)
}

// NewContextMenu creates and displays the context menu described by the
// given configuration.
//
// It panics if called before Run.
func NewContextMenu(c MenuConfig) Menu {
	return driver.NewContextMenu(c)
}

// NewController creates the controller described by the given configuration.
//
// It panics if called before Run.
func NewController(c ControllerConfig) Controller {
	return driver.NewController(c)
}

// NewFilePanel creates and displays the file panel described by the given
// configuration.
//
// It panics if called before Run.
func NewFilePanel(c FilePanelConfig) Elem {
	return driver.NewFilePanel(c)
}

// NewSaveFilePanel creates and displays the save file panel described by the
// given configuration.
//
// It panics if called before Run.
func NewSaveFilePanel(c SaveFilePanelConfig) Elem {
	return driver.NewSaveFilePanel(c)
}

// NewShare creates and display the share pannel to share the given value.
//
// It panics if called before Run.
func NewShare(v interface{}) Elem {
	return driver.NewShare(v)
}

// NewNotification creates and displays the notification described in the
// given configuration.
//
// It panics if called before Run.
func NewNotification(c NotificationConfig) Elem {
	return driver.NewNotification(c)
}

// MenuBar returns the menu bar.
//
// It panics if called before Run.
func MenuBar() Menu {
	return driver.MenuBar()
}

// NewStatusMenu creates and displays the status menu described in the given
// configuration.
//
// It panics if called before Run.
func NewStatusMenu(c StatusMenuConfig) StatusMenu {
	return driver.NewStatusMenu(c)
}

// Dock returns the dock tile.
//
// It panics if called before Run.
func Dock() DockTile {
	return driver.DockTile()
}

// Stop stops the app.
// Calling stop make Run return an error.
//
// It panics if called before Run.
func Stop() {
	driver.Stop()
}

// UI calls a function on the UI goroutine.
func UI(f func()) {
	driver.UI(f)
}

// Handle handles the message for the given key.
func Handle(key string, h Handler) {
	messages.handle(key, h)
}

// Post posts the given messages.
// Messages are handled in another goroutine.
func Post(msgs ...Msg) {
	messages.post(msgs...)
}

// NewMsg creates a message.
func NewMsg(key string) Msg {
	return &msg{key: key}
}

// Emit emits the event with the given value.
func Emit(e Event, value interface{}) {
	events.Emit(e, value)
}

// NewSubscriber creates an event subscriber to return when implementing the
// app.EventSubscriber interface.
func NewSubscriber() *Subscriber {
	return &Subscriber{
		Events: events,
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

// EnableDebug is a function that set whether debug mode is enabled.
func EnableDebug(v bool) {
	whenDebug = func(f func()) {}

	if v {
		whenDebug = func(f func()) {
			f()
		}
	}
}

// WhenDebug execute the given function when debug mode is enabled.
func WhenDebug(f func()) {
	whenDebug(f)
}

// CompoName returns the name of the given component.
// The returned name is the one to use in html tags.
func CompoName(c Compo) string {
	v := reflect.ValueOf(c)
	v = reflect.Indirect(v)

	name := strings.ToLower(v.Type().String())
	return strings.TrimPrefix(name, "main.")
}
