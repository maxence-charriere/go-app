package app

import (
	"context"
	"net/url"
	"sort"
	"sync"
	"time"

	"github.com/maxence-charriere/go-app/v8/pkg/errors"
)

const (
	eventBufferSize  = 4096
	updateBufferSize = 64
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

	initOnce  sync.Once
	startOnce sync.Once
	closeOnce sync.Once
	wait      sync.WaitGroup

	isMountedOnce bool
	events        chan event
	updates       map[Composer]struct{}
	updateQueue   []updateDescriptor
}

func (e *engine) Dispatch(src UI, fn func()) {
	e.events <- event{
		source:   src,
		function: fn,
	}
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
		select {
		case ev := <-e.events:
			ev.function()
			e.scheduleComponentUpdate(ev.source)

		default:
			if len(e.updates) == 0 {
				return
			}
			e.updateComponents()
		}
	}
}

func (e *engine) Close() {
	e.closeOnce.Do(func() {
		e.Consume()
		e.Wait()

		dismount(e.Body)
		e.Body = nil
		close(e.events)
	})
}

func (e *engine) PreRender() {
	e.Dispatch(e.Body, func() {
		e.Body.preRender(e.Page)
	})
}

func (e *engine) Mount(n UI) {
	e.Dispatch(e.Body, func() {
		if !e.isMountedOnce {
			if err := e.Body.(elemWithChildren).replaceChildAt(0, n); err != nil {
				panic(errors.New("mounting ui element failed").
					Tag("events-count", len(e.events)).
					Tag("events-capacity", cap(e.events)).
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
				Tag("events-count", len(e.events)).
				Tag("events-capacity", cap(e.events)).
				Tag("updates-count", len(e.updates)).
				Tag("updates-queue-len", len(e.updateQueue)).
				Wrap(err))
		}

		if err := e.Body.(elemWithChildren).replaceChildAt(0, n); err != nil {
			panic(errors.New("mounting ui element failed").
				Tag("events-count", len(e.events)).
				Tag("events-capacity", cap(e.events)).
				Tag("updates-count", len(e.updates)).
				Tag("updates-queue-len", len(e.updateQueue)).
				Wrap(err))
		}
	})
}

func (e *engine) Nav(u *url.URL) {
	if p, ok := e.Page.(*requestPage); ok {
		p.ReplaceURL(u)
	}

	e.Dispatch(e.Body, func() {
		e.Body.onNav(u)
	})
}

func (e *engine) AppUpdate() {
	e.Dispatch(e.Body, func() {
		e.Body.onAppUpdate()
	})
}

func (e *engine) AppResize() {
	e.Dispatch(e.Body, func() {
		e.Body.onResize()
	})
}

func (e *engine) init() {
	e.initOnce.Do(func() {
		e.events = make(chan event, eventBufferSize)
		e.updates = make(map[Composer]struct{})
		e.updateQueue = make([]updateDescriptor, 0, updateBufferSize)

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
			body := Body().Body(Div())
			if err := mount(e, body); err != nil {
				panic(errors.New("mounting engine default body failed").Wrap(err))
			}
			e.Body = body
		}
	})
}

func (e *engine) start(ctx context.Context) {
	e.startOnce.Do(func() {
		updates := time.NewTicker(time.Second / time.Duration(e.UpdateRate))
		defer updates.Stop()

		for {
			select {
			case <-ctx.Done():
				return

			case ev := <-e.events:
				ev.function()
				e.scheduleComponentUpdate(ev.source)

			case <-updates.C:
				e.updateComponents()
			}
		}
	})
}

func (e *engine) scheduleComponentUpdate(n UI) {
	if !n.Mounted() {
		return
	}

	var compo Composer
	var depth = 0

	for {
		if c, isCompo := n.(Composer); compo == nil && isCompo {
			if _, isScheduled := e.updates[c]; isScheduled {
				return
			}
			compo = c
		}

		parent := n.parent()
		if parent == nil {
			break
		}

		depth++
		n = parent
	}

	if compo == nil {
		return
	}

	e.updates[compo] = struct{}{}
	e.updateQueue = append(e.updateQueue, updateDescriptor{
		compo:    compo,
		priority: depth + 1,
	})
}

func (e *engine) updateComponents() {
	if len(e.updates) == 0 {
		return
	}

	sort.Slice(e.updateQueue, func(a, b int) bool {
		return e.updateQueue[a].priority < e.updateQueue[b].priority
	})

	for _, ud := range e.updateQueue {
		compo := ud.compo
		if !compo.Mounted() {
			continue
		}

		if err := compo.updateRoot(); err != nil {
			panic(err)
		}
		delete(e.updates, compo)
	}

	e.updateQueue = e.updateQueue[:0]
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

type event struct {
	source   UI
	function func()
}

type updateDescriptor struct {
	compo    Composer
	priority int
}
