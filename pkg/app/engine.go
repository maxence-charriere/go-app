package app

import (
	"context"
	"net/url"
	"sort"
	"sync"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

const (
	eventBufferSize  = 4096
	updateBufferSize = 64
	deferBufferSize  = 64
)

type engine struct {
	// The rate where component updates are performed (per seconds).
	UpdateRate int

	// The page.
	Page Page

	// Reports whether the engine runs in a server.
	RunsInServer bool

	// The storage use as local storage.
	LocalStorage BrowserStorage

	// The storage used as session storage.
	SessionStorage BrowserStorage

	// The function used to resolve static resource paths.
	ResolveStaticResources func(string) string

	// The body of the page.
	Body HTMLBody

	// The action handlers that are not associated with a component and are
	// executed asynchronously.
	ActionHandlers map[string]ActionHandler

	initOnce  sync.Once
	startOnce sync.Once
	closeOnce sync.Once
	wait      sync.WaitGroup

	isMountedOnce bool
	dispatches    chan Dispatch
	updates       map[Composer]struct{}
	updateQueue   []updateDescriptor
	defers        []Dispatch
	actions       actionManager
	states        *store
}

func (e *engine) Dispatch(d Dispatch) {
	if d.Source == nil {
		d.Source = e.Body
	}
	if d.Function == nil {
		d.Function = func(Context) {}
	}
	e.dispatches <- d
}

func (e *engine) Emit(src UI, fn func() bool) {
	if !src.Mounted() {
		return
	}

	if fn != nil {
		if fn() {
			// is true if updates should be skipped
			return
		}
	}

	compoCount := 0
	for n := src; n != nil; n = n.parent() {
		compo, ok := n.(Composer)
		if !ok {
			continue
		}

		compoCount++
		if compoCount > 1 {
			e.Dispatch(Dispatch{
				Source: compo,
				Mode:   Update,
			})
		}
	}
}

func (e *engine) Handle(actionName string, src UI, h ActionHandler) {
	e.actions.handle(actionName, false, src, h)
}

func (e *engine) SetState(state string, v interface{}, opts ...StateOption) {
	e.states.Set(state, v, opts...)
}

func (e *engine) GetState(state string, recv interface{}) {
	e.states.Get(state, recv)
}

func (e *engine) DelState(state string) {
	e.states.Del(state)
}

func (e *engine) ObserveState(state string, elem UI) Observer {
	return e.states.Observe(state, elem)
}

func (e *engine) Post(a Action) {
	e.Async(func() {
		e.actions.post(a)
	})
}

func (e *engine) Async(fn func()) {
	e.wait.Add(1)
	go func() {
		fn()
		e.wait.Done()
	}()
}

func (e *engine) Wait() {
	e.wait.Wait()
}

func (e *engine) Context() Context {
	return makeContext(e.Body)
}

func (e *engine) Consume() {
	for {
		e.Wait()

		select {
		case d := <-e.dispatches:
			e.handleDispatch(d)

		default:
			e.updateComponents()
			e.execDeferableEvents()
			return
		}
	}
}

func (e *engine) ConsumeNext() {
	e.Wait()

	select {
	case d := <-e.dispatches:
		e.handleDispatch(d)
		e.updateComponents()
		e.execDeferableEvents()

	default:
	}
}

func (e *engine) Close() {
	e.closeOnce.Do(func() {
		e.Consume()
		e.Wait()

		dismount(e.Body)
		e.Body = nil
		close(e.dispatches)

		e.states.Close()
	})
}

func (e *engine) PreRender() {
	e.Dispatch(Dispatch{
		Mode:   Update,
		Source: e.Body,
		Function: func(ctx Context) {
			ctx.Src().preRender(e.Page)
		},
	})
}

func (e *engine) Mount(n UI) {
	e.Dispatch(Dispatch{
		Mode:   Update,
		Source: e.Body,
		Function: func(ctx Context) {
			if !e.isMountedOnce {
				if err := e.Body.(elemWithChildren).replaceChildAt(0, n); err != nil {
					panic(errors.New("mounting ui element failed").
						Tag("dispatches-count", len(e.dispatches)).
						Tag("dispatches-capacity", cap(e.dispatches)).
						Tag("updates-count", len(e.updates)).
						Tag("updates-queue-len", len(e.updateQueue)).
						Wrap(err))
				}

				e.isMountedOnce = true
				return
			}

			err := update(e.Body.children()[0], n)
			if err == nil {
				return
			}
			if !isErrReplace(err) {
				panic(errors.New("mounting ui element failed").
					Tag("dispatches-count", len(e.dispatches)).
					Tag("dispatches-capacity", cap(e.dispatches)).
					Tag("updates-count", len(e.updates)).
					Tag("updates-queue-len", len(e.updateQueue)).
					Wrap(err))
			}

			if err := e.Body.(elemWithChildren).replaceChildAt(0, n); err != nil {
				panic(errors.New("mounting ui element failed").
					Tag("dispatches-count", len(e.dispatches)).
					Tag("dispatches-capacity", cap(e.dispatches)).
					Tag("updates-count", len(e.updates)).
					Tag("updates-queue-len", len(e.updateQueue)).
					Wrap(err))
			}
		},
	})
}

func (e *engine) Nav(u *url.URL) {
	if p, ok := e.Page.(*requestPage); ok {
		p.ReplaceURL(u)
	}

	e.Dispatch(Dispatch{
		Mode:   Update,
		Source: e.Body,
		Function: func(ctx Context) {
			ctx.Src().onNav(u)
		},
	})
}

func (e *engine) AppUpdate() {
	e.Dispatch(Dispatch{
		Mode:   Update,
		Source: e.Body,
		Function: func(ctx Context) {
			ctx.Src().onAppUpdate()
		},
	})
}

func (e *engine) AppInstallChange() {
	e.Dispatch(Dispatch{
		Mode:   Update,
		Source: e.Body,
		Function: func(ctx Context) {
			ctx.Src().onAppInstallChange()
		},
	})
}

func (e *engine) AppResize() {
	e.Dispatch(Dispatch{
		Mode:   Update,
		Source: e.Body,
		Function: func(ctx Context) {
			ctx.Src().onResize()
		},
	})
}

func (e *engine) init() {
	e.initOnce.Do(func() {
		e.dispatches = make(chan Dispatch, eventBufferSize)
		e.updates = make(map[Composer]struct{})
		e.updateQueue = make([]updateDescriptor, 0, updateBufferSize)
		e.defers = make([]Dispatch, 0, deferBufferSize)
		e.states = newStore(e)

		if e.UpdateRate <= 0 {
			e.UpdateRate = 60
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

		if e.ResolveStaticResources == nil {
			e.ResolveStaticResources = func(path string) string {
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

		for actionName, handler := range e.ActionHandlers {
			e.actions.handle(actionName, true, e.Body, handler)
		}
	})
}

func (e *engine) start(ctx context.Context) {
	e.startOnce.Do(func() {
		updateInterval := time.Second / time.Duration(e.UpdateRate)
		currentInterval := time.Duration(updateInterval)

		updates := time.NewTicker(currentInterval)
		defer updates.Stop()

		cleanup := time.NewTicker(time.Minute)
		defer cleanup.Stop()

		for {
			select {
			case <-ctx.Done():
				return

			case d := <-e.dispatches:
				if currentInterval != updateInterval {
					currentInterval = updateInterval
					updates.Reset(currentInterval)
				}

				e.handleDispatch(d)

			case <-updates.C:
				e.updateComponents()
				e.execDeferableEvents()

				if len(e.dispatches) == 0 {
					currentInterval = time.Hour
					updates.Reset(currentInterval)
				}

			case <-cleanup.C:
				e.actions.closeUnusedHandlers()
				e.states.Cleanup()
			}
		}
	})
}

func (e *engine) handleDispatch(d Dispatch) {
	switch d.Mode {
	case Next:
		d.Function(makeContext(d.Source))

	case Update:
		if d.Source.Mounted() {
			ctx := makeContext(d.Source).(uiContext)
			d.Function(ctx)
			if *ctx.skipUpdates < 2 {
				e.scheduleComponentUpdate(d.Source)
			}
		}

	case Defer:
		if d.Source.Mounted() {
			e.defers = append(e.defers, d)
		}
	}
}

func (e *engine) scheduleComponentUpdate(n UI) {
	if !n.Mounted() {
		return
	}

	c := nearestCompo(n)
	if c == nil {
		return
	}

	if _, isScheduled := e.updates[c]; isScheduled {
		return
	}

	e.updates[c] = struct{}{}
	e.updateQueue = append(e.updateQueue, updateDescriptor{
		compo:    c,
		priority: compoPriority(c),
	})
}

func (e *engine) updateComponents() {
	if len(e.updates) == 0 {
		return
	}

	sortUpdateDescriptors(e.updateQueue)
	for _, ud := range e.updateQueue {
		compo := ud.compo
		if !compo.Mounted() {
			e.removeFromUpdates(compo)
			continue
		}

		if _, requiresUpdate := e.updates[compo]; !requiresUpdate {
			continue
		}

		if err := compo.updateRoot(); err != nil {
			panic(err)
		}
		e.removeFromUpdates(compo)
	}

	e.updateQueue = e.updateQueue[:0]
}

func (e *engine) removeFromUpdates(c Composer) {
	delete(e.updates, c)
}

func (e *engine) execDeferableEvents() {
	for _, d := range e.defers {
		if d.Source.Mounted() {
			d.Function(makeContext(d.Source))
		}
	}
	e.defers = e.defers[:0]
}

func (e *engine) currentPage() Page {
	return e.Page
}

func (e *engine) localStorage() BrowserStorage {
	return e.LocalStorage
}

func (e *engine) sessionStorage() BrowserStorage {
	return e.SessionStorage
}

func (e *engine) runsInServer() bool {
	return e.RunsInServer
}

func (e *engine) resolveStaticResource(path string) string {
	return e.ResolveStaticResources(path)
}

type updateDescriptor struct {
	compo    Composer
	priority int
}

func sortUpdateDescriptors(d []updateDescriptor) {
	sort.Slice(d, func(a, b int) bool {
		return d[a].priority < d[b].priority
	})
}

func nearestCompo(n UI) Composer {
	for node := n; node != nil; node = node.parent() {
		if c, isCompo := node.(Composer); isCompo {
			return c
		}
	}
	return nil
}

func compoPriority(c Composer) int {
	depth := 1
	for parent := c.parent(); parent != nil; parent = parent.parent() {
		depth++
	}
	return depth
}

type msgHandler struct {
	src      UI
	function MsgHandler
}
