package app

import (
	"context"
	"net/url"
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

// Reload reloads the WebAssembly app at the current page.
func (ctx Context) Reload() {
	if IsServer {
		return
	}

	ctx.Dispatch(func() {
		Window().Get("location").Call("reload")
	})
}

// Navigate navigates to the given URL. This is a helper method that converts
// rawURL to an *url.URL and then calls ctx.NavigateTo under the hood.
func (ctx Context) Navigate(rawURL string) {
	ctx.Dispatch(func() {
		navigate(ctx.dispatcher, rawURL)
	})
}

// NavigateTo navigates to the given URL.
func (ctx Context) NavigateTo(u *url.URL) {
	ctx.Dispatch(func() {
		navigateTo(ctx.dispatcher, u)
	})
}

func makeContext(src UI) Context {
	return Context{
		Context:            src.context(),
		Src:                src,
		JSSrc:              src.JSValue(),
		AppUpdateAvailable: appUpdateAvailable,
		Page:               src.Dispatcher().currentPage(),
		dispatcher:         src.Dispatcher(),
	}
}
