package app

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

// Context is the interface that describes a context tied to a UI element.
//
// A context provides mechanisms to deal with the browser, the current page,
// navigation, concurrency, and component communication.
//
// It is canceled when its associated UI element is dismounted.
type Context interface {
	context.Context

	// Returns the UI element tied to the context.
	Src() UI

	// Returns the associated JavaScript value. The is an helper method for:
	//  ctx.Src.JSValue()
	JSSrc() Value

	// Reports whether the app has been updated in background. Use app.Reload()
	// to load the updated version.
	AppUpdateAvailable() bool

	// Returns the current page.
	Page() Page

	// Executes the given function on the UI goroutine and notifies the
	// context's nearest component to update its state.
	Dispatch(fn func(Context))

	// Executes the given function on the UI goroutine after notifying the
	// context's nearest component to update its state.
	Defer(fn func(Context))

	// Registers the handler for the given action name. When an action occurs,
	// the handler is executed on the UI goroutine.
	Handle(actionName string, h ActionHandler)

	// Returns a builder that creates an action with the given name. The action
	// can be then posted with the Post() method:
	//  ctx.NewAction(""myAction").Post()
	NewAction(name string) ActionBuilder

	// Executes the given function on a new goroutine.
	//
	// The difference versus just launching a goroutine is that it ensures that
	// the asynchronous function is called before a page is fully pre-rendered
	// and served over HTTP.
	Async(fn func())

	// Asynchronously waits for the given duration and dispatches the given
	// function.
	After(d time.Duration, fn func(Context))

	// Executes the given function and notifies the parent components to update
	// their state. It should be used to launch component custom event handlers.
	Emit(fn func())

	// Reloads the WebAssembly app to the current page. It is like refreshing
	// the browser page.
	Reload()

	// Navigates to the given URL. This is a helper method that converts url to
	// an *url.URL and then calls ctx.NavigateTo under the hood.
	Navigate(url string)

	// Navigates to the given URL.
	NavigateTo(u *url.URL)

	// Resolves the given path to make it point to the right location whether
	// static resources are located on a local directory or a remote bucket.
	ResolveStaticResource(string) string

	// Returns a storage that uses the browser local storage associated to the
	// document origin. Data stored has no expiration time.
	LocalStorage() BrowserStorage

	// Returns a storage that uses the browser session storage associated to the
	// document origin. Data stored expire when the page session ends.
	SessionStorage() BrowserStorage

	// Scrolls to the HTML element with the given id.
	ScrollTo(id string)

	// Returns a UUID that identifies the app on the current device.
	DeviceID() string

	// Encrypts the given value using AES encryption.
	Encrypt(v interface{}) ([]byte, error)

	// Decrypts the given encrypted bytes and stores them in the given value.
	Decrypt(crypted []byte, v interface{}) error

	dispatcher() Dispatcher
}

type uiContext struct {
	context.Context

	src                UI
	jsSrc              Value
	appUpdateAvailable bool
	page               Page
	disp               Dispatcher
}

func (ctx uiContext) Src() UI {
	return ctx.src
}

func (ctx uiContext) JSSrc() Value {
	return ctx.jsSrc
}

func (ctx uiContext) AppUpdateAvailable() bool {
	return ctx.appUpdateAvailable
}

func (ctx uiContext) Page() Page {
	return ctx.page
}

func (ctx uiContext) Dispatch(fn func(Context)) {
	ctx.dispatcher().Dispatch(ctx.Src(), fn)
}

func (ctx uiContext) Defer(fn func(Context)) {
	ctx.dispatcher().Defer(ctx.Src(), fn)
}

func (ctx uiContext) Handle(actionName string, h ActionHandler) {
	ctx.dispatcher().Handle(actionName, ctx.Src(), h)
}

func (ctx uiContext) NewAction(name string) ActionBuilder {
	return newActionBuilder(ctx.dispatcher(), name)
}

func (ctx uiContext) Async(fn func()) {
	ctx.dispatcher().Async(fn)
}

func (ctx uiContext) After(d time.Duration, fn func(Context)) {
	ctx.Async(func() {
		time.Sleep(d)
		ctx.Dispatch(fn)
	})
}

func (ctx uiContext) Emit(fn func()) {
	ctx.dispatcher().Emit(ctx.Src(), fn)
}

func (ctx uiContext) Reload() {
	if IsServer {
		return
	}
	ctx.Defer(func(ctx Context) {
		Window().Get("location").Call("reload")
	})
}

func (ctx uiContext) Navigate(rawURL string) {
	ctx.Defer(func(ctx Context) {
		navigate(ctx.dispatcher(), rawURL)
	})
}

func (ctx uiContext) NavigateTo(u *url.URL) {
	ctx.Defer(func(ctx Context) {
		navigateTo(ctx.dispatcher(), u, true)
	})
}

func (ctx uiContext) ResolveStaticResource(path string) string {
	return ctx.dispatcher().resolveStaticResource(path)
}

func (ctx uiContext) LocalStorage() BrowserStorage {
	return ctx.dispatcher().localStorage()
}

func (ctx uiContext) SessionStorage() BrowserStorage {
	return ctx.dispatcher().sessionStorage()
}

func (ctx uiContext) ScrollTo(id string) {
	ctx.Defer(func(ctx Context) {
		Window().ScrollToID(id)
	})
}

func (ctx uiContext) DeviceID() string {
	var id string
	if err := ctx.LocalStorage().Get("/go-app/deviceID", &id); err != nil {
		panic(errors.New("retrieving device id failed").Wrap(err))
	}
	if id != "" {
		return id
	}

	id = uuid.New().String()
	if err := ctx.LocalStorage().Set("/go-app/deviceID", id); err != nil {
		panic(errors.New("creating device id failed").Wrap(err))
	}
	return id
}

func (ctx uiContext) Encrypt(v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, errors.New("encoding value failed").Wrap(err)
	}

	b, err = encrypt(ctx.cryptoKey(), b)
	if err != nil {
		return nil, errors.New("encrypting value failed").Wrap(err)
	}
	return b, nil
}

func (ctx uiContext) Decrypt(crypted []byte, v interface{}) error {
	b, err := decrypt(ctx.cryptoKey(), crypted)
	if err != nil {
		return errors.New("decrypting value failed").Wrap(err)
	}

	if err := json.Unmarshal(b, v); err != nil {
		return errors.New("decoding value failed").Wrap(err)
	}
	return nil
}

func (ctx uiContext) cryptoKey() string {
	return strings.ReplaceAll(ctx.DeviceID(), "-", "")
}

func (ctx uiContext) dispatcher() Dispatcher {
	return ctx.disp
}

func makeContext(src UI) Context {
	return uiContext{
		Context:            src.context(),
		src:                src,
		jsSrc:              src.JSValue(),
		appUpdateAvailable: appUpdateAvailable,
		page:               src.dispatcher().currentPage(),
		disp:               src.dispatcher(),
	}
}
