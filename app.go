package app

import "github.com/pkg/errors"

var (
	// Loggers contains the loggers used by the app.
	Loggers []Logger

	// ErrNotFound is the error returned when fetching nonexistent elements or
	// components.
	ErrNotFound = errors.New("not found")

	components = NewFactory()
	backends   = make([]Backend, 0, 5)
	uiChan     = make(chan func(), 1024)
	events     = newEventRegistry(CallOnUIGoroutine)
	actions    = newActionRegistry(events)
)

// Import imports the given component type into the app.
// Components must be imported in order the be used by the app package.
// This allows components to be created dynamically when they are found into
// markup.
func Import(c Component) {
	if _, err := components.Register(c); err != nil {
		Log("import %T failed: %s", c, err)
		panic(err)
	}
}

// Run runs the app.
// It blocks until the app is closed.
func Run() error {
	for _, b := range backends {
		if err := b.Run(components, uiChan); err != nil {
			return err
		}
	}
	return nil
}

// Render the given value.
// Valid values are:
//   - app.Element
//   - app.Component
// Backends decide whether they render a value.
// Other types of value are ignored.
func Render(v interface{}) {
	for _, b := range backends {
		if err := b.Render(v); err != nil {
			Log("rendering %T failed: %s", v)
		}
	}
}

// CallOnUIGoroutine calls the given function on the UI dedicated goroutine.
func CallOnUIGoroutine(f func()) {
	uiChan <- f
}

// PostAction creates and posts the named action with the given arg.
// The action is then handled in a separate goroutine.
func PostAction(name string, arg interface{}) {
	actions.Post(name, arg)
}

// PostActions posts a batch of actions.
// All the actions are handled sequentially in a separate goroutine.
func PostActions(a ...Action) {
	actions.PostBatch(a...)
}

// HandleAction handles the named action with the given action handler.
func HandleAction(name string, h ActionHandler) {
	actions.Handle(name, h)
}

// NewEventSubscriber creates an event subscriber.
func NewEventSubscriber() EventSubscriber {
	return &eventSubscriber{
		registry: events,
	}
}

// Log logs a message according to a format specifier.
// It is a helper function that calls Log() for all the loggers set in
// app.Loggers.
func Log(format string, v ...interface{}) {
	for _, l := range Loggers {
		l.Log(format, v...)
	}
}

// Debug logs a debug message according to a format specifier.
// It is a helper function that calls Debug() for all the loggers set in
// app.Loggers.
func Debug(format string, v ...interface{}) {
	for _, l := range Loggers {
		l.Debug(format, v...)
	}
}

// WhenDebug execute the given function when debug mode is enabled.
// It is a helper function that calls WhenDebug() for all the loggers set in
// app.Loggers.
func WhenDebug(f func()) {
	for _, l := range Loggers {
		l.WhenDebug(f)
	}
}
