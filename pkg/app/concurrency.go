package app

import (
	"context"
	"net/url"
	"reflect"

	"github.com/maxence-charriere/go-app/v7/pkg/errors"
)

const (
	dispatcherSize = 4096
)

// Dispatcher is the inerface that describes an environment that synchronizes UI
// instructions and components lifecycle.
type Dispatcher interface {
	// Dispatch enqueues the given function to be executed on a goroutine
	// dedicated to managing UI modifications.
	Dispatch(func())

	start(context.Context)
	isServerSideMode() bool
}

// TestingDispatcher represents a dispatcher to use for testing purposes.
type TestingDispatcher interface {
	Dispatcher

	// Pre-renders the given component.
	PreRender(Page)

	// Mounts the given component as root component.
	Mount(UI)

	// Triggers OnNav from the root component.
	Nav(*url.URL)

	// Triggers OnAppUpdate from the root component.
	AppUpdate()

	// Triggers OnAppResize from the root component.
	AppResize()

	// Consume executes all the remaining UI instructions.
	Consume()

	// Close consumes all the remaining UI instruction and releases allocated
	// resources.
	Close()
}

// NewClientTestingDispatcher creates a testing dispatcher that simulates a
// client environment. The given UI element is mounted upon creation.
func NewClientTestingDispatcher(v UI) TestingDispatcher {
	return newTestingDispatcher(v, false)
}

// NewServerTestingDispatcher creates a testing dispatcher that simulates a
// client environment. The given UI element is mounted upon creation.
func NewServerTestingDispatcher(v UI) TestingDispatcher {
	return newTestingDispatcher(v, false)
}

func newTestingDispatcher(v UI, serverSide bool) TestingDispatcher {
	disp := &uiDispatcher{
		ui:             make(chan func(), dispatcherSize),
		serverSideMode: serverSide,
		body: Body().Body(
			Div(),
		).(*htmlBody),
	}

	if err := mount(disp, disp.body); err != nil {
		panic(errors.New("mounting body failed").
			Tag("server-side-mode", disp.isServerSideMode()).
			Tag("body-type", reflect.TypeOf(disp.body)).
			Tag("ui-len", len(disp.ui)).
			Tag("ui-cap", cap(disp.ui)).
			Wrap(err))
	}

	disp.Mount(v)
	return disp
}

type uiDispatcher struct {
	ui             chan func()
	body           *htmlBody
	serverSideMode bool
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

func (d *uiDispatcher) PreRender(p Page) {
	d.Dispatch(func() {
		d.body.preRender(p)
	})
}

func (d *uiDispatcher) Mount(v UI) {
	d.Dispatch(func() {
		if err := d.body.replaceChildAt(0, v); err != nil {
			panic(errors.New("mounting ui element failed").
				Tag("server-side-mode", d.isServerSideMode()).
				Tag("body-type", reflect.TypeOf(d.body)).
				Tag("ui-len", len(d.ui)).
				Tag("ui-cap", cap(d.ui)).
				Wrap(err))
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
	for {
		select {
		case fn := <-d.ui:
			fn()

		default:
			return
		}
	}
}

func (d *uiDispatcher) Close() {
	if len(d.ui) != 0 {
		d.Consume()
	}

	dismount(d.body)
	close(d.ui)
}

func (d *uiDispatcher) start(ctx context.Context) {
	for {
		select {
		case fn := <-d.ui:
			fn()

		case <-ctx.Done():
			return
		}
	}
}

func (d *uiDispatcher) isServerSideMode() bool {
	return d.serverSideMode
}
