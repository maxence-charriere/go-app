package app

import (
	"bytes"
	"context"
	"encoding/json"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/errors"
)

type engineX struct {
	ctx context.Context

	localStorage   BrowserStorage
	sessionStorage BrowserStorage
	browser        browser

	routes         *router
	internalURLs   []string
	resolveURL     func(string) string
	originPage     *requestPage
	lastVisitedURL *url.URL

	nodes   nodeManager
	updates updateManager
	body    HTMLBody

	dispatches chan func()
	defers     chan func()
	goroutines sync.WaitGroup

	asynchronousActionHandlers map[string]ActionHandler
	actions                    actionManager
	states                     stateManager
}

func newEngine(ctx context.Context, routes *router, resolveURL func(string) string, originPage *requestPage, actionHandlers map[string]ActionHandler) *engineX {
	var localStorage BrowserStorage
	var sessionStorage BrowserStorage
	if IsServer {
		localStorage = newMemoryStorage()
		sessionStorage = newMemoryStorage()
	} else {
		localStorage = newJSStorage("localStorage")
		sessionStorage = newJSStorage("sessionStorage")
	}

	if resolveURL == nil {
		resolveURL = func(v string) string { return v }
	}
	originPage.resolveURL = resolveURL

	engine := &engineX{
		ctx:                        ctx,
		routes:                     routes,
		resolveURL:                 resolveURL,
		originPage:                 originPage,
		localStorage:               localStorage,
		lastVisitedURL:             &url.URL{},
		sessionStorage:             sessionStorage,
		nodes:                      nodeManager{},
		dispatches:                 make(chan func(), 4096),
		defers:                     make(chan func(), 4096),
		asynchronousActionHandlers: actionHandlers,
	}

	engine.initBrowser()
	return engine
}

func (e *engineX) baseContext() Context {
	return Context{
		Context:               e.ctx,
		resolveURL:            e.resolveURL,
		appUpdatable:          e.browser.AppUpdatable,
		page:                  e.page,
		navigate:              e.Navigate,
		localStorage:          e.localStorage,
		sessionStorage:        e.sessionStorage,
		dispatch:              e.dispatch,
		defere:                e.defere,
		async:                 e.async,
		addComponentUpdate:    e.updates.Add,
		removeComponentUpdate: e.updates.Done,
		handleAction:          e.actions.Handle,
		postAction:            e.actions.Post,
		observeState:          e.states.Observe,
		unObserveState:        e.states.UnObserve,
		getState:              e.states.Get,
		setState:              e.states.Set,
		delState:              e.states.Delete,

		notifyComponentEvent: e.nodes.NotifyComponentEvent,
	}
}

// Navigate directs the engine to the specified URL destination, which might be
// an internal page within the app, an external link outside the app, or a
// mailto link. If the 'updateHistory' flag is true, the destination is added to
// the browser's history.
func (e *engineX) Navigate(destination *url.URL, updateHistory bool) {
	if destination.Host == "" {
		destination.Host = e.originPage.URL().Host
	}

	switch {
	case e.internalURL(destination),
		e.mailTo(destination):
		Window().Get("location").Set("href", destination.String())
		return

	case e.externalNavigation(destination):
		Window().Call("open", destination.String())
		return

	case destination.String() == e.lastVisitedURL.String():
		return
	}

	defer func() {
		if updateHistory {
			Window().addHistory(destination)
		}
		e.lastVisitedURL = destination

		e.nodes.NotifyComponentEvent(e.baseContext(), e.body, nav{})

		if destination.Fragment != "" {
			e.defere(func() {
				Window().ScrollToID(destination.Fragment)
			})
		}
	}()

	if destination.Path == e.lastVisitedURL.Path &&
		destination.Fragment != e.lastVisitedURL.Fragment {
		return
	}

	path := strings.TrimPrefix(destination.Path, Getenv("GOAPP_ROOT_PREFIX"))
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	root, ok := e.routes.createComponent(path)
	if !ok {
		root = &notFound{}
	}

	if err := e.Load(root); err != nil {
		panic(errors.New("loading component failed").
			WithTag("component-type", reflect.TypeOf(root)).
			Wrap(err))
	}
}

func (e *engineX) initBrowser() {
	if IsServer {
		return
	}
	e.browser.HandleEvents(e.baseContext(), e.notifyComponentEvent)
}

func (e *engineX) notifyComponentEvent(event any) {
	e.nodes.NotifyComponentEvent(e.baseContext(), e.body, event)
}

func (e *engineX) externalNavigation(v *url.URL) bool {
	return v.Host != e.originPage.URL().Host
}

func (e *engineX) mailTo(v *url.URL) bool {
	return v.Scheme == "mailto"
}

func (e *engineX) internalURL(v *url.URL) bool {
	if e.internalURLs == nil {
		json.Unmarshal([]byte(Getenv("GOAPP_INTERNAL_URLS")), &e.internalURLs)
	}

	url := v.String()
	for _, u := range e.internalURLs {
		if strings.HasPrefix(url, u) {
			return true
		}
	}
	return false
}

func (e *engineX) page() Page {
	if IsClient {
		return makeBrowserPage(e.resolveURL)
	}
	return e.originPage
}

func (e *engineX) Load(v Composer) error {
	if e.body == nil {
		body := Body()
		body = body.setJSElement(Window().Get("document").Get("body")).(HTMLBody)

		firstChild := Div()
		firstChild = firstChild.setJSElement(body.JSValue().firstElementChild()).(HTMLDiv)
		firstChild = firstChild.setParent(body).(HTMLDiv)

		body = body.setBody([]UI{firstChild}).(HTMLBody)
		e.body = body

		for action, handler := range e.asynchronousActionHandlers {
			e.actions.Handle(action, body, true, handler)
		}
	}

	body, err := e.nodes.Update(e.baseContext(), e.body, Body().privateBody(v))
	if err != nil {
		return errors.New("updating root failed").Wrap(err)
	}
	e.body = body.(HTMLBody)
	return nil
}

// Start initiates the main event loop of the engine at the specified framerate.
// The loop efficiently manages dispatches, component updates, and deferred
// actions.
func (e *engineX) Start(framerate int) {
	if framerate <= 0 {
		framerate = 30
	}

	iddleFrameDuration := time.Hour
	activeFrameDuration := time.Second / time.Duration(framerate)
	currentFrameDuration := time.Nanosecond
	frames := time.NewTicker(currentFrameDuration)
	defer frames.Stop()

	e.states.CleanupExpiredPersistedStates(e.baseContext())

	for {
		select {
		case dispatch := <-e.dispatches:
			if currentFrameDuration != activeFrameDuration {
				frames.Reset(activeFrameDuration)
				currentFrameDuration = activeFrameDuration
			}
			dispatch()

		case <-frames.C:
			e.processFrame()
			frames.Reset(iddleFrameDuration)
			currentFrameDuration = iddleFrameDuration

		case <-e.ctx.Done():
			return
		}
	}
}

func (e *engineX) processFrame() {
	e.updates.UpdateForEach(func(c Composer) {
		if !c.Mounted() {
			return
		}

		if _, err := e.nodes.UpdateComponentRoot(e.baseContext(), c); err != nil {
			panic(errors.New("updating component failed").Wrap(err))
		}
	})
	e.executeDefers()
	e.actions.Cleanup()
	e.states.Cleanup()
}

func (e *engineX) executeDefers() {
	for {
		select {
		case defere := <-e.defers:
			defere()

		default:
			return
		}
	}
}

// ConsumeNext waits for any ongoing goroutines to finish, then executes the
// next dispatch in the queue. After executing the dispatch, it processes a
// frame.
func (e *engineX) ConsumeNext() {
	e.goroutines.Wait()
	dispatch := <-e.dispatches
	dispatch()
	e.processFrame()
}

// ConsumeAll continuously waits for ongoing goroutines to finish, executes all
// available dispatches in the queue until none are left, and then processes a
// frame.
func (e *engineX) ConsumeAll() {
	for {
		select {
		case dispatch := <-e.dispatches:
			dispatch()

		default:
			e.processFrame()
			e.goroutines.Wait()
			if len(e.dispatches) == 0 {
				return
			}
		}
	}
}

// Encode serializes the given HTML element, integrating the engine's root
// component as the initial child within the document's body. The final HTML
// content, including the standard DOCTYPE declaration, is written  to the
// provided buffer.
func (e *engineX) Encode(w *bytes.Buffer, document HTMLHtml) error {
	if e.body == nil {
		return errors.New("no component loaded")
	}
	root := e.body.body()[0]

	var body HTML
	for _, child := range document.(HTML).body() {
		if child, isBody := child.(HTMLBody); isBody {
			body = child
			break
		}
	}
	if body == nil {
		return errors.New("document does not have a body")
	}

	children := make([]UI, 0, len(body.body())+1)
	children = append(children, root)
	children = append(children, body.body()...)
	body.setBody(children)

	w.WriteString("<!doctype html>\n")
	e.nodes.Encode(e.baseContext(), w, document)
	return nil
}

func (e *engineX) dispatch(v func()) {
	e.dispatches <- v
}

func (e *engineX) defere(v func()) {
	e.defers <- v
}

func (e *engineX) async(v func()) {
	e.goroutines.Add(1)
	go func() {
		v()
		e.goroutines.Done()
	}()
}
