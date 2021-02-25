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

// Async launches the given function on a new goroutine.
//
// The difference versus just launching a goroutine is that it ensures that the
// asynchronous function is called before a page is fully pre-rendered and
// served over HTTP.
func (ctx Context) Async(fn func()) {
	ctx.dispatcher.Async(fn)
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
		navigateTo(ctx.dispatcher, u, true)
	})
}

// ResolveStaticResource resolves the given path to make it point to the right
// location whether static resources are located on a local directory or a
// remote bucket.
func (ctx Context) ResolveStaticResource(path string) string {
	return ctx.dispatcher.resolveStaticResource(path)
}

// LocalStorage returns a storage that uses the browser local storage associated
// to the document origin. Data stored has no expiration time.
func (ctx Context) LocalStorage() BrowserStorage {
	return ctx.dispatcher.localStorage()
}

// SessionStorage returns a storage that uses the browser session storage
// associated to the document origin. Data stored expire when the page
// session ends.
func (ctx Context) SessionStorage() BrowserStorage {
	return ctx.dispatcher.sessionStorage()
}

// ScrollTo scrolls to the HTML element with the given id.
func (ctx Context) ScrollTo(id string) {
	ctx.Dispatch(func() {
		Window().ScrollToID(id)
	})
}

func makeContext(src UI) Context {
	return Context{
		Context:            src.context(),
		Src:                src,
		JSSrc:              src.JSValue(),
		AppUpdateAvailable: appUpdateAvailable,
		Page:               src.dispatcher().currentPage(),
		dispatcher:         src.dispatcher(),
	}
}
