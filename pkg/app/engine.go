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
	Body *HTMLBody

	init sync.Once
	wait sync.WaitGroup

	events      chan event
	updates     map[Composer]struct{}
	updateQueue []updateDescriptor
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

func (e *engine) start(ctx context.Context) {
	e.init.Do(func() {
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
		}

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
				if len(e.updates) > 0 {
					e.updateComponents()
				}
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
