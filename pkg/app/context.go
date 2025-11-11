package app

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/maxence-charriere/go-app/v10/pkg/errors"
)

// Context represents a UI element-associated environment enabling interactions
// with the browser, page navigation, concurrency, and component communication.
type Context struct {
	context.Context

	page                  func() Page
	appUpdatable          bool
	resolveURL            func(string) string
	navigate              func(*url.URL, bool)
	localStorage          BrowserStorage
	sessionStorage        BrowserStorage
	dispatch              func(func())
	defere                func(func())
	async                 func(func())
	addComponentUpdate    func(Composer, int)
	removeComponentUpdate func(Composer)
	handleAction          func(string, UI, bool, ActionHandler)
	postAction            func(Context, Action)
	observeState          func(Context, string, any) Observer
	unObserveState        func(Context, string)
	getState              func(Context, string, any)
	setState              func(Context, string, any) State
	delState              func(Context, string)

	sourceElement        UI
	notifyComponentEvent func(Context, UI, any)
}

// Src retrieves the linked UI element of the context.
func (ctx Context) Src() UI {
	return ctx.sourceElement
}

// JSSrc fetches the JavaScript representation of the associated UI element.
func (ctx Context) JSSrc() Value {
	return ctx.sourceElement.JSValue()
}

// AppUpdateAvailable checks if there's a pending app update.
func (ctx Context) AppUpdateAvailable() bool {
	return ctx.appUpdatable
}

// IsAppInstallable verifies if the app is eligible for installation.
func (ctx Context) IsAppInstallable() bool {
	if Window().Get("goappIsAppInstallable").Truthy() {
		return Window().Call("goappIsAppInstallable").Bool()
	}
	return false
}

// IsAppleBrowser reports whether the app is running on an Apple browser.
func (ctx Context) IsAppleBrowser() bool {
	if Window().Get("goappIsAppleBrowser").Truthy() {
		return Window().Call("goappIsAppleBrowser").Bool()
	}
	return false
}

// ShowAppInstallPrompt initiates the app installation process.
func (ctx Context) ShowAppInstallPrompt() {
	if ctx.IsAppInstallable() {
		Window().Call("goappShowInstallPrompt")
	}
}

// DeviceID fetches a distinct identifier for the app on the present device.
func (ctx Context) DeviceID() string {
	var id string
	if err := ctx.localStorage.Get("/go-app/deviceID", &id); err != nil {
		panic(errors.New("retrieving device id failed").Wrap(err))
	}
	if id != "" {
		return id
	}

	id = uuid.NewString()
	if err := ctx.localStorage.Set("/go-app/deviceID", id); err != nil {
		panic(errors.New("creating device id failed").Wrap(err))
	}
	return id
}

// Page retrieves the current active page.
func (ctx Context) Page() Page {
	return ctx.page()
}

// Reload refreshes the present page.
func (ctx Context) Reload() {
	if IsServer {
		return
	}
	Window().Get("location").Call("reload")
}

// Navigate transitions to the given URL string.
func (ctx Context) Navigate(rawURL string) {
	u, err := url.Parse(rawURL)
	if err != nil {
		Log(errors.New("navigating to URL failed").
			WithTag("url", rawURL).
			Wrap(err))
		return
	}
	ctx.NavigateTo(u)
}

// NavigateTo transitions to the provided URL.
func (ctx Context) NavigateTo(u *url.URL) {
	ctx.navigate(u, true)
}

// ResolveStaticResource adjusts a given path to point to the correct static
// resource location.
func (ctx Context) ResolveStaticResource(v string) string {
	return ctx.resolveURL(v)
}

// ScrollTo adjusts the scrollbar to target an HTML element by its ID.
func (ctx Context) ScrollTo(id string) {
	ctx.Defer(func(ctx Context) {
		Window().ScrollToID(id)
	})
}

// LocalStorage accesses the browser's local storage tied to the document
// origin.
func (ctx Context) LocalStorage() BrowserStorage {
	return ctx.localStorage
}

// SessionStorage accesses the browser's session storage tied to the
// document origin.
func (ctx Context) SessionStorage() BrowserStorage {
	return ctx.sessionStorage
}

// Encrypt enciphers a value using AES encryption.
func (ctx Context) Encrypt(v any) ([]byte, error) {
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

// Decrypt deciphers encrypted data into a given reference value.
func (ctx Context) Decrypt(crypted []byte, v any) error {
	b, err := decrypt(ctx.cryptoKey(), crypted)
	if err != nil {
		return errors.New("decrypting value failed").Wrap(err)
	}

	if err := json.Unmarshal(b, v); err != nil {
		return errors.New("decoding value failed").Wrap(err)
	}
	return nil
}

func (ctx Context) cryptoKey() string {
	return strings.ReplaceAll(ctx.DeviceID(), "-", "")
}

// Notifications accesses the notifications service.
func (ctx Context) Notifications() NotificationService {
	return NotificationService{}
}

// Dispatch prompts the execution of a function on the UI goroutine,
// flagging the enclosing component for an update.
func (ctx Context) Dispatch(v func(Context)) {
	ctx.dispatch(func() {
		if !ctx.sourceElement.Mounted() {
			return
		}

		for c, ok := component(ctx.sourceElement); ok; c, ok = component(c.parent()) {
			ctx.addComponentUpdate(c, 1)
		}

		if v != nil {
			v(ctx)
		}
	})
}

// Defer postpones the function execution on the UI goroutine until the
// current update cycle completes.
func (ctx Context) Defer(v func(Context)) {
	ctx.defere(func() {
		if !ctx.sourceElement.Mounted() {
			return
		}

		if v != nil {
			v(ctx)
		}
	})
}

// Async initiates a function asynchronously. It enables go-app to monitor
// goroutines, ensuring they conclude when rendering server-side.
func (ctx Context) Async(v func()) {
	ctx.async(v)
}

// After pauses for a determined span, then triggers a specified function.
func (ctx Context) After(d time.Duration, f func(Context)) {
	ctx.async(func() {
		time.Sleep(d)
		ctx.Dispatch(f)
	})
}

// PreventUpdate halts updates for the enclosing component.
func (ctx Context) PreventUpdate() {
	for c, ok := component(ctx.sourceElement); ok; c, ok = component(c.parent()) {
		ctx.addComponentUpdate(c, -1)
	}
}

// Update flags the enclosing component for an update.
func (ctx Context) Update() {
	ctx.Dispatch(nil)
}

// Handle designates a handler for a particular action, set to run on the UI
// goroutine.
func (ctx Context) Handle(action string, h ActionHandler) {
	ctx.handleAction(action, ctx.sourceElement, false, h)
}

// NewAction generates a new action for handling.
func (ctx Context) NewAction(action string, tags ...Tagger) {
	ctx.NewActionWithValue(action, nil, tags...)

}

// NewActionWithValue crafts an action with a given value for processing.
func (ctx Context) NewActionWithValue(action string, v any, tags ...Tagger) {
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

// ObserveState establishes an observer for a state, tracking its changes.
func (ctx Context) ObserveState(state string, recv any) Observer {
	return ctx.observeState(ctx, state, recv)
}

func (ctx Context) UnObserveState(state string) {
	ctx.unObserveState(ctx, state)
}

// GetState fetches the value of a particular state.
func (ctx Context) GetState(state string, recv any) {
	ctx.getState(ctx, state, recv)
}

// SetState modifies a state with the provided value.
func (ctx Context) SetState(state string, v any) State {
	return ctx.setState(ctx, state, v)
}

// DelState erases a state, halting all associated observations.
func (ctx Context) DelState(state string) {
	ctx.delState(ctx, state)
}

// ResizeContent notifies the children of the associated element that implement
// the Resizer interface about a resize event. It ensures that components can
// adjust their size and layout in response to changes. This method is typically
// used when the size of the container changes, requiring child components to
// update their dimensions accordingly.
func (ctx Context) ResizeContent() {
	ctx.Defer(func(ctx Context) {
		ctx.Dispatch(func(ctx Context) {
			ctx.notifyComponentEvent(ctx, ctx.Src(), resize{})
		})
	})
}
