package app

import (
	"context"
	"net/url"
	"reflect"

	"github.com/maxence-charriere/go-app/v8/pkg/errors"
)

const (
	dispatcherSize = 4096
)

// Dispatcher is the interface that describes an environment that synchronizes
// UI instructions and UI elements lifecycle.
type Dispatcher interface {
	// Dispatch enqueues the given function to be executed on a goroutine
	// dedicated to managing UI modifications.
	Dispatch(func())

	start(context.Context)
	currentPage() Page
	isServerSideMode() bool
	resolveStaticResource(string) string
}

// TestingDispatcher represents a dispatcher to use for testing purposes.
type TestingDispatcher interface {
	Dispatcher

	// Pre-renders the given component.
	PreRender()

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
	u, _ := url.Parse("https://localhost")

	disp := newUIDispatcher(serverSide, &requestPage{url: u}, func(url string) string {
		return url
	})
	disp.body = Body().Body(
		Div(),
	).(elemWithChildren)

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
	ui                        chan func()
	body                      elemWithChildren
	page                      Page
	mountedOnce               bool
	serverSideMode            bool
	resolveStaticResourceFunc func(string) string
}

func newUIDispatcher(serverSide bool, p Page, resolveStaticResource func(string) string) *uiDispatcher {
	return &uiDispatcher{
		ui:                        make(chan func(), dispatcherSize),
		page:                      p,
		serverSideMode:            serverSide,
		resolveStaticResourceFunc: resolveStaticResource,
	}
}

func (d *uiDispatcher) Dispatch(fn func()) {
	d.ui <- fn
}

func (d *uiDispatcher) PreRender() {
	d.Dispatch(func() {
		d.body.preRender(d.currentPage())
	})
}

func (d *uiDispatcher) Mount(v UI) {
	d.Dispatch(func() {
		if !d.mountedOnce {
			if err := d.body.(elemWithChildren).replaceChildAt(0, v); err != nil {
				panic(errors.New("mounting ui element failed").
					Tag("server-side-mode", d.isServerSideMode()).
					Tag("body-type", reflect.TypeOf(d.body)).
					Tag("ui-len", len(d.ui)).
					Tag("ui-cap", cap(d.ui)).
					Wrap(err))
			}
			d.mountedOnce = true
			return
		}

		err := update(d.body.children()[0], v)
		if err == nil {
			return
		}
		if !isErrReplace(err) {
			panic(errors.New("mounting ui element failed").
				Tag("server-side-mode", d.isServerSideMode()).
				Tag("body-type", reflect.TypeOf(d.body)).
				Tag("ui-len", len(d.ui)).
				Tag("ui-cap", cap(d.ui)).
				Wrap(err))
		}

		if err := d.body.(elemWithChildren).replaceChildAt(0, v); err != nil {
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
	if p, ok := d.currentPage().(*requestPage); ok {
		p.url = u
	}

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
	d.Consume()
	dismount(d.body)
	d.body = nil
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

func (d *uiDispatcher) currentPage() Page {
	return d.page
}

func (d *uiDispatcher) isServerSideMode() bool {
	return d.serverSideMode
}

func (d *uiDispatcher) resolveStaticResource(url string) string {
	return d.resolveStaticResourceFunc(url)
}
