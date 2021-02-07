package app

import (
	"context"
	"net/url"
	"sync"

	"github.com/maxence-charriere/go-app/v7/pkg/errors"
)

const (
	dispatcherSize = 4096
)

// Dispatcher represents a dispatcher that queues functions to be executed on a
// goroutine dedicated to performing UI instructions.
type Dispatcher interface {
	// Dispatch enqueues the given function to be executed on a goroutine
	// dedicated to managing UI modifications.
	Dispatch(func())

	start(context.Context)
}

// TestingDispatcher represents a dispatcher to use for testing purposes.
type TestingDispatcher interface {
	Dispatcher

	// Pre-renders the given component.
	PreRender(v UI)

	// Mounts the given component as root component.
	Mount(v UI)

	// Triggers OnNav from the root component.
	Nav(u *url.URL)

	// Triggers OnAppUpdate from the root component.
	AppUpdate()

	// Triggers OnAppResize from the root component.
	AppResize()

	// Close execute the remaining UI instructions and releases the allocated
	// resources.
	Consume()
}

// NewTestingDispatcher creates a testing dispatcher.
func NewTestingDispatcher(v UI) TestingDispatcher {
	disp := &uiDispatcher{
		ui: make(chan func(), dispatcherSize),
		body: Body().Body(
			Div(),
		).(*htmlBody),
	}

	if err := mount(disp, v); err != nil {
		panic(errors.New("mounting body failed").Wrap(err))
	}

	return disp
}

type uiDispatcher struct {
	startOnce   sync.Once
	concumeOnce sync.Once
	ui          chan func()
	body        *htmlBody
}

func newUIDispatcher(body *htmlBody) *uiDispatcher {
	return &uiDispatcher{
		ui:   make(chan func(), dispatcherSize),
		body: body,
	}
}

func (d *uiDispatcher) Dispatch(fn func()) {
	d.ui <- fn
}

func (d *uiDispatcher) PreRender(v UI) {
	d.Mount(v)
	d.Dispatch(func() {
		// preRender(, ??)
	})
}

func (d *uiDispatcher) Mount(v UI) {
	d.Dispatch(func() {
		if err := d.body.replaceChildAt(0, v); err != nil {
			panic(errors.New("mounting ui element failed").Wrap(err))
		}
	})
}

func (d *uiDispatcher) Nav(u *url.URL) {
	d.Dispatch(func() {
		d.body.onNav(u)
	})
}

func (d *uiDispatcher) AppUpdate() {
	d.Dispatch(func() {
		d.body.onAppUpdate()
	})
}

func (d *uiDispatcher) AppResize() {
	d.Dispatch(func() {
		d.body.onAppResize()
	})
}

func (d *uiDispatcher) Consume() {
	d.concumeOnce.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		d.Dispatch(cancel)
		d.start(ctx)
	})
}

func (d *uiDispatcher) start(ctx context.Context) {
	d.startOnce.Do(func() {
		for {
			select {
			case fn := <-d.ui:
				fn()

			case <-ctx.Done():
				break
			}
		}

		close(d.ui)
	})
}
