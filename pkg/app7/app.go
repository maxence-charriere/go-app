package app

import "log"

var (
	// Logger is the logger used to log info and errors.
	Logger func(format string, v ...interface{}) = log.Printf
)

// EventHandler represents a function that can handle HTML events.
type EventHandler func(ctx Context, e Event)

// Window returns the JavaScript "window" object.
func Window() BrowserWindow {
	return window
}
