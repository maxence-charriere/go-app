package app

import (
	"context"
	"net/url"
	"sync"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

type engineX struct {
	// The number of frame per seconds.
	FrameRate int

	// The page.
	Page Page

	// Reports whether the engine runs on server-side.
	IsServerSide bool

	// The storage use as local storage.
	LocalStorage BrowserStorage

	// The storage used as session storage.
	SessionStorage BrowserStorage

	// The function used to resolve static resource paths.
	StaticResourceResolver func(string) string

	// The body of the page.
	Body HTMLBody

	// The action handlers that are not associated with a component and are
	// executed asynchronously.
	ActionHandlers map[string]ActionHandler

	initOnce  sync.Once
	closeOnce sync.Once
	wait      sync.WaitGroup

	dispatches       chan Dispatch
	componentUpdates map[Composer]bool
	deferables       []Dispatch
	actions          actionManager
	states           *store
	isFirstMount     bool
}

func (e *engineX) Context() Context {
	return makeContext(e.Body)
}

func (e *engineX) Dispatch(d Dispatch) {
	if d.Source == nil {
		d.Source = e.Body
	}
	if d.Function == nil {
		d.Function = func(Context) {}
	}
	e.dispatches <- d
}

func (e *engineX) Emit(src UI, fn func()) {
	e.Dispatch(Dispatch{
		Mode:   Next,
		Source: src,
		Function: func(ctx Context) {
			if fn != nil {
				fn()
			}

			for c := getParentComponent(src); c != nil; c = getParentComponent(c) {
				e.addComponentUpdate(c)
			}
		},
	})
}

func (e *engineX) Handle(actionName string, src UI, h ActionHandler) {
	e.actions.handle(actionName, false, src, h)
}

func (e *engineX) Post(a Action) {
	e.Async(func() {
		e.actions.post(a)
	})
}

func (e *engineX) SetState(state string, v any, opts ...StateOption) {
	e.states.Set(state, v, opts...)
}

func (e *engineX) GetState(state string, recv any) {
	e.states.Get(state, recv)
}

func (e *engineX) DelState(state string) {
	e.states.Del(state)
}

func (e *engineX) ObserveState(state string, elem UI) Observer {
	return e.states.Observe(state, elem)
}

func (e *engineX) Async(fn func()) {
	e.wait.Add(1)
	go func() {
		fn()
		e.wait.Done()
	}()
}

func (e *engineX) Wait() {
	e.wait.Wait()
}

func (e *engineX) Consume() {
	for {
		e.Wait()

		select {
		case d := <-e.dispatches:
			e.handleDispatch(d)

		default:
			e.handleFrame()
		}
	}
}

func (e *engineX) ConsumeNext() {
	e.Wait()
	e.handleDispatch(<-e.dispatches)
	e.handleFrame()
}

func (e *engineX) Close() {
	e.closeOnce.Do(func() {
		e.Consume()
		e.Wait()

		dismount(e.Body)
		e.Body = nil
		e.states.Close()
	})
}

func (e *engineX) PreRender() {
	e.Dispatch(Dispatch{
		Mode:   Update,
		Source: e.Body,
		Function: func(ctx Context) {
			ctx.Src().preRender(e.Page)
		},
	})
}

func (e *engineX) Mount(v UI) {
	e.Dispatch(Dispatch{
		Mode:   Update,
		Source: e.Body,
		Function: func(ctx Context) {
			if !e.isFirstMount {
				if err := e.Body.(*htmlBody).replaceChildAt(0, v); err != nil {
					panic(errors.New("mounting first ui element failed").Wrap(err))
				}

				e.isFirstMount = false
				return
			}

			if firstChild := e.Body.getChildren()[0]; canUpdate(firstChild, v) {
				if err := update(firstChild, v); err != nil {
					panic(errors.New("mounting ui element failed").Wrap(err))
				}
				return
			}

			if err := e.Body.(*htmlBody).replaceChildAt(0, v); err != nil {
				panic(errors.New("mounting ui element failed").Wrap(err))
			}
		},
	})
}

func (e *engineX) Nav(u *url.URL) {
	if p, ok := e.Page.(*requestPage); ok {
		p.ReplaceURL(u)
	}

	e.Dispatch(Dispatch{
		Mode:   Update,
		Source: e.Body,
		Function: func(ctx Context) {
			ctx.Src().onComponentEvent(nav{})
		},
	})
}

func (e *engineX) AppUpdate() {
	e.Dispatch(Dispatch{
		Mode:   Update,
		Source: e.Body,
		Function: func(ctx Context) {
			ctx.Src().onComponentEvent(appUpdate{})
		},
	})
}

func (e *engineX) AppInstallChange() {
	e.Dispatch(Dispatch{
		Mode:   Update,
		Source: e.Body,
		Function: func(ctx Context) {
			ctx.Src().onComponentEvent(appInstallChange{})
		},
	})
}

func (e *engineX) AppResize() {
	e.Dispatch(Dispatch{
		Mode:   Update,
		Source: e.Body,
		Function: func(ctx Context) {
			ctx.Src().onComponentEvent(resize{})
		},
	})
}

func (e *engineX) init() {
	if e.FrameRate <= 0 {
		e.FrameRate = 60
	}

	if e.Page == nil {
		u, _ := url.Parse("https://test.go-app.dev")
		e.Page = &requestPage{url: u}
	}

	if e.LocalStorage == nil {
		e.LocalStorage = newMemoryStorage()
	}

	if e.SessionStorage == nil {
		e.SessionStorage = newMemoryStorage()
	}

	if e.StaticResourceResolver == nil {
		e.StaticResourceResolver = func(path string) string {
			return path
		}
	}

	if e.Body == nil {
		body := Body().privateBody(Div())
		if err := mount(e, body); err != nil {
			panic(errors.New("mounting engine default body failed").Wrap(err))
		}
		e.Body = body
	}

	e.dispatches = make(chan Dispatch, 4096)
	e.componentUpdates = make(map[Composer]bool)
	e.deferables = make([]Dispatch, 32)
	e.states = newStore(e)
	e.isFirstMount = true

	for actionName, handler := range e.ActionHandlers {
		e.actions.handle(actionName, true, e.Body, handler)
	}
}

func (e *engineX) getCurrentPage() Page {
	return e.Page
}

func (e *engineX) getLocalStorage() BrowserStorage {
	return e.LocalStorage
}

func (e *engineX) getSessionStorage() BrowserStorage {
	return e.SessionStorage
}

func (e *engineX) isServerSide() bool {
	return e.IsServerSide
}

func (e *engineX) resolveStaticResource(path string) string {
	return e.StaticResourceResolver(path)
}

func (e *engineX) addComponentUpdate(c Composer) {
	if c == nil {
		return
	}
	if _, isAdded := e.componentUpdates[c]; isAdded {
		return
	}
	e.componentUpdates[c] = true
}

func (e *engineX) preventComponentUpdate(c Composer) {
	e.componentUpdates[c] = false
}

func (e *engineX) addDeferable(d Dispatch) {
	e.deferables = append(e.deferables, d)
}

func (e *engineX) start(ctx context.Context) {
	frameDuration := time.Second / time.Duration(e.FrameRate)
	frames := time.NewTicker(frameDuration)

	cleanups := time.NewTicker(time.Minute)
	defer cleanups.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case d := <-e.dispatches:
			e.handleDispatch(d)

		case <-frames.C:
			e.handleFrame()

		case <-cleanups.C:
			e.actions.closeUnusedHandlers()
			e.states.Cleanup()
		}
	}
}

func (e *engineX) handleDispatch(d Dispatch) {
	switch d.Mode {
	case Update:
		d.do()
		e.addComponentUpdate(getParentComponent(d.Source))

	case Defer:
		e.deferables = append(e.deferables, d)

	case Next:
		d.do()
	}
}

func (e *engineX) handleFrame() {
	e.handleComponentUpdates()
	e.handleDeferables()
}

func (e *engineX) handleComponentUpdates() {
	for component, canUppdate := range e.componentUpdates {
		if !component.Mounted() || !canUppdate {
			delete(e.componentUpdates, component)
			continue
		}

		if err := component.updateRoot(); err != nil {
			panic(err)
		}
		delete(e.componentUpdates, component)
	}
}

func (e *engineX) handleDeferables() {
	for i := range e.deferables {
		e.deferables[i].do()
		e.deferables[i] = Dispatch{}
	}

	e.deferables = e.deferables[:0]
}
