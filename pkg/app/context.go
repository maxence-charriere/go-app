package app

import (
	"context"
)

// Context represents a context that is tied to a UI element. It is canceled
// when the element is dismounted.
//
// It implements the context.Context interface.
//  https://golang.org/pkg/context/#Context
type Context struct {
	context.Context

	// The UI element tied to the context.
	Src UI

	// The JavaScript value of the element tied to the context. This is a
	// shorthand for:
	//  ctx.Src.JSValue()
	JSSrc Value

	// Reports whether the app has been updated in background. Use app.Reload()
	// to load the updated version.
	AppUpdateAvailable bool

	// The info about the current page.
	Page Page

	dispatcher Dispatcher
}

// Dispatch executes the given function on the goroutine dedicated to updating
// the UI.
func (ctx Context) Dispatch(fn func()) {
	ctx.dispatcher.Dispatch(fn)
}

func makeContext(src UI, p Page) Context {
	return Context{
		Context:            src.context(),
		Src:                src,
		JSSrc:              src.JSValue(),
		AppUpdateAvailable: appUpdateAvailable,
		Page:               p,
		dispatcher:         src.dispatcher(),
	}
}
