package app

import (
	"context"
	"net/url"
)

const (
	dispatcherSize = 4096
)

// Dispatcher is the interface that describes an environment that synchronizes
// UI instructions and UI elements lifecycle.
type Dispatcher interface {
	// Dispatch executes the given function on the UI goroutine and notifies the
	// source's nearest component to update its state.
	Dispatch(src UI, fn func(Context))

	// Defer executes the given function on the UI goroutine after notifying the
	// source's nearest component to update its state.
	Defer(src UI, fn func(Context))

	// 	Async launches the given function on a new goroutine.
	//
	// The difference versus just launching a goroutine is that it ensures that
	// the asynchronous instructions are called before the dispatcher is closed.
	//
	// This is important during component prerendering since asynchronous
	// operations need to complete before sending a pre-rendered page over HTTP.
	Async(fn func())

	// Wait waits for the asynchronous operations launched with Async() to
	// complete.
	Wait()

	start(context.Context)
	currentPage() Page
	localStorage() BrowserStorage
	sessionStorage() BrowserStorage
	runsInServer() bool
	resolveStaticResource(string) string
}

// ClientDispatcher is the interface that describes a dispatcher that emulates a
// client environment.
type ClientDispatcher interface {
	Dispatcher

	// Context returns the context associated with the root element.
	Context() Context

	// Consume executes all the remaining UI instructions.
	Consume()

	// Close consumes all the remaining UI instruction and releases allocated
	// resources.
	Close()

	// Mounts the given component as root element.
	Mount(UI)

	// Triggers OnNav from the root component.
	Nav(*url.URL)

	// Triggers OnAppUpdate from the root component.
	AppUpdate()

	// Triggers OnAppResize from the root component.
	AppResize()
}

// NewClientTester creates a testing dispatcher that simulates a
// client environment. The given UI element is mounted upon creation.
func NewClientTester(n UI) ClientDispatcher {
	e := &engine{}
	e.init()
	e.Mount(n)
	e.Consume()
	return e
}

// ServerDispatcher is the interface that describes a dispatcher that emulates a server environment.
type ServerDispatcher interface {
	Dispatcher

	// Context returns the context associated with the root element.
	Context() Context

	// Consume executes all the remaining UI instructions.
	Consume()

	// Close consumes all the remaining UI instruction and releases allocated
	// resources.
	Close()

	// Pre-renders the given component.
	PreRender()
}

// NewServerTester creates a testing dispatcher that simulates a
// client environment.
func NewServerTester(n UI) ServerDispatcher {
	e := &engine{RunsInServer: false}
	e.init()
	e.Mount(n)
	e.Consume()
	return e
}
