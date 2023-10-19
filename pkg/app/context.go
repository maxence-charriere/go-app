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

// Context represents a UI element-associated environment enabling interactions
// with the browser, page navigation, concurrency, and component communication.
type Context interface {
	context.Context

	// Src retrieves the linked UI element of the context.
	Src() UI

	// JSSrc fetches the JavaScript representation of the associated UI element.
	JSSrc() Value

	// AppUpdateAvailable checks if there's a pending app update.
	AppUpdateAvailable() bool

	// IsAppInstallable verifies if the app is eligible for installation.
	IsAppInstallable() bool

	// ShowAppInstallPrompt initiates the app installation process.
	ShowAppInstallPrompt()

	// DeviceID fetches a distinct identifier for the app on the present device.
	DeviceID() string

	// Page retrieves the current active page.
	Page() Page

	// Reload refreshes the present page.
	Reload()

	// Navigate transitions to the given URL string.
	Navigate(url string)

	// NavigateTo transitions to the provided URL.
	NavigateTo(u *url.URL)

	// ResolveStaticResource adjusts a given path to point to the correct static
	// resource location.
	ResolveStaticResource(string) string

	// ScrollTo adjusts the scrollbar to target an HTML element by its ID.
	ScrollTo(id string)

	// LocalStorage accesses the browser's local storage tied to the document
	// origin.
	LocalStorage() BrowserStorage

	// SessionStorage accesses the browser's session storage tied to the
	// document origin.
	SessionStorage() BrowserStorage

	// Encrypt enciphers a value using AES encryption.
	Encrypt(v any) ([]byte, error)

	// Decrypt deciphers encrypted data into a given reference value.
	Decrypt(crypted []byte, v any) error

	// Notifications accesses the notifications service.
	Notifications() NotificationService

	// Dispatch prompts the execution of a function on the UI goroutine,
	// flagging the enclosing component for an update, respecting any
	// implemented UpdateNotifier behavior.
	Dispatch(fn func(Context))

	// Defer postpones the function execution on the UI goroutine until the
	// current update cycle completes.
	Defer(fn func(Context))

	// Async initiates a function asynchronously. It enables go-app to monitor
	// goroutines, ensuring they conclude when rendering server-side.
	Async(fn func())

	// After pauses for a determined span, then triggers a specified function.
	After(d time.Duration, fn func(Context))

	// PreventUpdate halts updates for the enclosing component, respecting any
	// implemented UpdateNotifier behavior.
	PreventUpdate()

	// Handle designates a handler for a particular action, set to run on the UI
	// goroutine.
	Handle(action string, h ActionHandler)

	// NewAction generates a new action for handling.
	NewAction(name string, tags ...Tagger)

	// NewActionWithValue crafts an action with a given value for processing.
	NewActionWithValue(name string, v any, tags ...Tagger)

	// SetState modifies a state with the provided value.
	SetState(state string, v any, opts ...StateOption)

	// GetState fetches the value of a particular state.
	GetState(state string, recv any)

	// DelState erases a state, halting all associated observations.
	DelState(state string)

	// ObserveState establishes an observer for a state, tracking its changes.
	ObserveState(state string) Observer

	// Dispatcher fetches the app's dispatcher.
	Dispatcher() Dispatcher
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

func (ctx uiContext) IsAppInstallable() bool {
	if Window().Get("goappIsAppInstallable").Truthy() {
		return Window().Call("goappIsAppInstallable").Bool()
	}
	return false
}

func (ctx uiContext) IsAppInstalled() bool {
	if Window().Get("goappIsAppInstalled").Truthy() {
		return Window().Call("goappIsAppInstalled").Bool()
	}
	return false
}

func (ctx uiContext) ShowAppInstallPrompt() {
	if ctx.IsAppInstallable() {
		Window().Call("goappShowInstallPrompt")
	}
}

func (ctx uiContext) Page() Page {
	return ctx.page
}

func (ctx uiContext) Dispatch(fn func(Context)) {
	ctx.Dispatcher().Dispatch(Dispatch{
		Mode:     Update,
		Source:   ctx.Src(),
		Function: fn,
	})
}

func (ctx uiContext) Defer(fn func(Context)) {
	ctx.Dispatcher().Dispatch(Dispatch{
		Mode:     Defer,
		Source:   ctx.Src(),
		Function: fn,
	})
}

func (ctx uiContext) Handle(actionName string, h ActionHandler) {
	ctx.Dispatcher().Handle(actionName, ctx.Src(), h)
}

func (ctx uiContext) NewAction(name string, tags ...Tagger) {
	ctx.NewActionWithValue(name, nil, tags...)
}

func (ctx uiContext) NewActionWithValue(name string, v any, tags ...Tagger) {
	var tagMap Tags
	for _, t := range tags {
		if tagMap == nil {
			tagMap = t.Tags()
			continue
		}
		for k, v := range t.Tags() {
			tagMap[k] = v
		}
	}

	ctx.Dispatcher().Post(Action{
		Name:  name,
		Value: v,
		Tags:  tagMap,
	})
}

func (ctx uiContext) Async(fn func()) {
	ctx.Dispatcher().Async(fn)
}

func (ctx uiContext) After(d time.Duration, fn func(Context)) {
	ctx.Async(func() {
		time.Sleep(d)
		ctx.Dispatch(fn)
	})
}

func (ctx uiContext) Emit(fn func()) {
	ctx.Dispatcher().Emit(ctx.Src(), fn)
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
		navigate(ctx.Dispatcher(), rawURL)
	})
}

func (ctx uiContext) NavigateTo(u *url.URL) {
	ctx.Defer(func(ctx Context) {
		navigateTo(ctx.Dispatcher(), u, true)
	})
}

func (ctx uiContext) ResolveStaticResource(path string) string {
	return ctx.Dispatcher().resolveStaticResource(path)
}

func (ctx uiContext) LocalStorage() BrowserStorage {
	return ctx.Dispatcher().getLocalStorage()
}

func (ctx uiContext) SessionStorage() BrowserStorage {
	return ctx.Dispatcher().getSessionStorage()
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

	id = uuid.NewString()
	if err := ctx.LocalStorage().Set("/go-app/deviceID", id); err != nil {
		panic(errors.New("creating device id failed").Wrap(err))
	}
	return id
}

func (ctx uiContext) Encrypt(v any) ([]byte, error) {
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

func (ctx uiContext) Decrypt(crypted []byte, v any) error {
	b, err := decrypt(ctx.cryptoKey(), crypted)
	if err != nil {
		return errors.New("decrypting value failed").Wrap(err)
	}

	if err := json.Unmarshal(b, v); err != nil {
		return errors.New("decoding value failed").Wrap(err)
	}
	return nil
}

func (ctx uiContext) SetState(state string, v any, opts ...StateOption) {
	ctx.Dispatcher().SetState(state, v, opts...)
}

func (ctx uiContext) GetState(state string, recv any) {
	ctx.Dispatcher().GetState(state, recv)
}

func (ctx uiContext) DelState(state string) {
	ctx.Dispatcher().DelState(state)
}

func (ctx uiContext) ObserveState(state string) Observer {
	return ctx.Dispatcher().ObserveState(state, ctx.src)
}

func (ctx uiContext) Dispatcher() Dispatcher {
	return ctx.disp
}

func (ctx uiContext) Notifications() NotificationService {
	return NotificationService{dispatcher: ctx.Dispatcher()}
}

func (ctx uiContext) PreventUpdate() {
	ctx.Dispatcher().preventComponentUpdate(getComponent(ctx.src))
}

func (ctx uiContext) cryptoKey() string {
	return strings.ReplaceAll(ctx.DeviceID(), "-", "")
}

func makeContext(src UI) Context {
	return uiContext{
		Context:            context.Background(),
		src:                src,
		jsSrc:              src.JSValue(),
		appUpdateAvailable: appUpdateAvailable,
		page:               src.getDispatcher().getCurrentPage(),
		disp:               src.getDispatcher(),
	}
}

type nodeContext struct {
	context.Context

	sourceElement             UI
	page                      func() Page
	appUpdatable              bool
	resolveURL                func(string) string
	navigate                  func(*url.URL, bool)
	localStorage              BrowserStorage
	sessionStorage            BrowserStorage
	dispatch                  func(func())
	defere                    func(func())
	async                     func(func())
	addComponentUpdate        func(Composer)
	removeComponentUpdate     func(Composer)
	foreachUpdatableComponent func(UI, func(Composer))
	handleAction              func(string, UI, bool, ActionHandler)
	postAction                func(Context, Action)
}

func (ctx nodeContext) Src() UI {
	return ctx.sourceElement
}

func (ctx nodeContext) JSSrc() Value {
	return ctx.sourceElement.JSValue()
}

func (ctx nodeContext) AppUpdateAvailable() bool {
	return ctx.appUpdatable
}

func (ctx nodeContext) IsAppInstallable() bool {
	if Window().Get("goappIsAppInstallable").Truthy() {
		return Window().Call("goappIsAppInstallable").Bool()
	}
	return false
}

func (ctx nodeContext) ShowAppInstallPrompt() {
	if ctx.IsAppInstallable() {
		Window().Call("goappShowInstallPrompt")
	}
}

func (ctx nodeContext) DeviceID() string {
	var id string
	if err := ctx.localStorage.Get("/go-app/deviceID", &id); err != nil {
		panic(errors.New("retrieving device id failed").Wrap(err))
	}
	if id != "" {
		return id
	}

	id = uuid.NewString()
	if err := ctx.sessionStorage.Set("/go-app/deviceID", id); err != nil {
		panic(errors.New("creating device id failed").Wrap(err))
	}
	return id
}

func (ctx nodeContext) Page() Page {
	return ctx.page()
}

func (ctx nodeContext) Reload() {
	if IsServer {
		return
	}
	Window().Get("location").Call("reload")
}

func (ctx nodeContext) Navigate(rawURL string) {
	u, err := url.Parse(rawURL)
	if err != nil {
		Log(errors.New("navigating to URL failed").
			WithTag("url", rawURL).
			Wrap(err))
		return
	}
	ctx.NavigateTo(u)
}

func (ctx nodeContext) NavigateTo(u *url.URL) {
	ctx.navigate(u, true)
}

func (ctx nodeContext) ResolveStaticResource(v string) string {
	return ctx.resolveURL(v)
}

func (ctx nodeContext) ScrollTo(id string) {
	ctx.Defer(func(ctx Context) {
		Window().ScrollToID(id)
	})
}

func (ctx nodeContext) LocalStorage() BrowserStorage {
	return ctx.localStorage
}

func (ctx nodeContext) SessionStorage() BrowserStorage {
	return ctx.sessionStorage
}

func (ctx nodeContext) Encrypt(v any) ([]byte, error) {
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

func (ctx nodeContext) Decrypt(crypted []byte, v any) error {
	b, err := decrypt(ctx.cryptoKey(), crypted)
	if err != nil {
		return errors.New("decrypting value failed").Wrap(err)
	}

	if err := json.Unmarshal(b, v); err != nil {
		return errors.New("decoding value failed").Wrap(err)
	}
	return nil
}

func (ctx nodeContext) cryptoKey() string {
	return strings.ReplaceAll(ctx.DeviceID(), "-", "")
}

func (ctx nodeContext) Notifications() NotificationService {
	panic("not implemented")
}

func (ctx nodeContext) Dispatch(v func(Context)) {
	ctx.dispatch(func() {
		if !ctx.sourceElement.Mounted() {
			return
		}
		ctx.foreachUpdatableComponent(ctx.sourceElement, ctx.addComponentUpdate)
		v(ctx)
	})
}

func (ctx nodeContext) Defer(v func(Context)) {
	ctx.defere(func() {
		if !ctx.sourceElement.Mounted() {
			return
		}
		v(ctx)
	})
}

func (ctx nodeContext) Async(v func()) {
	ctx.async(v)
}

func (ctx nodeContext) After(d time.Duration, f func(Context)) {
	ctx.async(func() {
		time.Sleep(d)
		ctx.Dispatch(f)
	})
}

func (ctx nodeContext) PreventUpdate() {
	ctx.foreachUpdatableComponent(ctx.sourceElement, ctx.removeComponentUpdate)
}

func (ctx nodeContext) Handle(action string, h ActionHandler) {
	ctx.handleAction(action, ctx.sourceElement, false, h)
}

func (ctx nodeContext) NewAction(action string, tags ...Tagger) {
	ctx.NewActionWithValue(action, nil, tags...)

}

func (ctx nodeContext) NewActionWithValue(action string, v any, tags ...Tagger) {
	var tagMap Tags
	for _, tag := range tags {
		if tagMap == nil {
			tagMap = make(Tags)
		}

		for k, v := range tag.Tags() {
			tagMap[k] = v
		}
	}

	ctx.postAction(ctx, Action{
		Name:  action,
		Value: v,
		Tags:  tagMap,
	})
}

func (ctx nodeContext) SetState(state string, v any, opts ...StateOption) {
	panic("not implemented")
}

func (ctx nodeContext) GetState(state string, recv any) {
	panic("not implemented")
}

func (ctx nodeContext) DelState(state string) {
	panic("not implemented")
}

func (ctx nodeContext) ObserveState(state string) Observer {
	panic("not implemented")
}

func (ctx nodeContext) Dispatcher() Dispatcher {
	panic("to deprecate")
}
