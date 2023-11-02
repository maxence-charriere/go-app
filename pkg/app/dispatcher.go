package app

import (
	"context"
	"net/url"
)

// Dispatcher is the interface that describes an environment that synchronizes
// UI instructions and UI elements lifecycle.
type Dispatcher interface {
	// Context returns the context associated with the root element.
	Context() Context

	// Handle registers the handler for the given action name. When an action
	// occurs, the handler is executed on the UI goroutine.
	Handle(actionName string, src UI, h ActionHandler)

	// Wait waits for the asynchronous operations launched with Async() to
	// complete.
	Wait()

	start(context.Context)
	getCurrentPage() Page
	getLocalStorage() BrowserStorage
	getSessionStorage() BrowserStorage
	resolveStaticResource(string) string
	removeComponentUpdate(Composer)
	preventComponentUpdate(Composer)
}

// ClientDispatcher is the interface that describes a dispatcher that emulates a
// client environment.
type ClientDispatcher interface {
	Dispatcher

	// Consume executes all the remaining UI instructions.
	Consume()

	// ConsumeNext executes the next UI instructions.
	ConsumeNext()

	// Close consumes all the remaining UI instruction and releases allocated
	// resources.
	Close()

	// Mounts the given component as root element.
	Mount(UI)

	// Triggers OnNav from the root component.
	Nav(*url.URL)

	// Triggers OnAppUpdate from the root component.
	AppUpdate()

	// Triggers OnAppInstallChange from the root component.
	AppInstallChange()

	// Triggers OnAppResize from the root component.
	AppResize()
}

// NewClientTester creates a testing dispatcher that simulates a
// client environment. The given UI element is mounted upon creation.
// func NewClientTester(n UI) ClientDispatcher {
// panic("not implemented")
// e := &engine{
// 	ActionHandlers: actionHandlers,
// }

// if IsClient {
// 	e.LocalStorage = newJSStorage("localStorage")
// 	e.LocalStorage.Clear()

// 	e.SessionStorage = newJSStorage("sessionStorage")
// 	e.SessionStorage.Clear()
// }

// e.init()
// e.Mount(n)
// e.Consume()
// return e
// }

// ServerDispatcher is the interface that describes a dispatcher that emulates a server environment.
type ServerDispatcher interface {
	Dispatcher

	// Consume executes all the remaining UI instructions.
	Consume()

	// ConsumeNext executes the next UI instructions.
	ConsumeNext()

	// Close consumes all the remaining UI instruction and releases allocated
	// resources.
	Close()
}

// NewServerTester creates a testing dispatcher that simulates a
// client environment.
func NewServerTester(n UI) ServerDispatcher {
	panic("not implemented")

	// e := &engine{
	// 	ActionHandlers: actionHandlers,
	// }
	// e.init()
	// e.Mount(n)
	// e.Consume()
	// return e
}
