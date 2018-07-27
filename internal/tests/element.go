package tests

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
)

type element struct {
	id string
}

func newElement() *element {
	elem := &element{
		id: uuid.New().String(),
	}
	return elem
}

func (e *element) ID() string {
	return e.id
}

type elemWithCompo struct {
	id        string
	factory   app.Factory
	lastFocus time.Time
	component app.Compo
}

func newElemWithCompo() *elemWithCompo {
	factory := app.NewFactory()
	factory.Register(&Foo{})

	return &elemWithCompo{
		id:        uuid.New().String(),
		factory:   factory,
		lastFocus: time.Now(),
	}
}

func (e *elemWithCompo) ID() string {
	return e.id
}

func (e *elemWithCompo) Load(rawurl string, v ...interface{}) error {
	rawurl = fmt.Sprintf(rawurl, v...)

	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}

	compo, err := e.factory.New(app.CompoNameFromURL(u))
	if err != nil {
		return err
	}

	e.component = compo
	return nil
}

func (e *elemWithCompo) Compo() app.Compo {
	return e.component
}

func (e *elemWithCompo) Contains(c app.Compo) bool {
	return e.component == c
}

func (e *elemWithCompo) Render(c app.Compo) error {
	e.component = c
	return nil
}

func (e *elemWithCompo) LastFocus() time.Time {
	return e.lastFocus
}

func testElemWithCompo(t *testing.T, newElem func() (app.ElemWithCompo, error)) {
	tests := []struct {
		scenario string
		function func(t *testing.T, elem app.ElemWithCompo)
	}{
		{
			scenario: "load a component",
			function: testElemWithCompoLoadSuccess,
		},
		{
			scenario: "load a component fails",
			function: testElemWithCompoLoadFail,
		},
		{
			scenario: "render a component",
			function: testElemWithCompoRenderSuccess,
		},
		{
			scenario: "render a component fails",
			function: testElemWithCompoRenderFail,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			elem, err := newElem()
			if err == app.ErrNotSupported {
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if closer, ok := elem.(app.Closer); ok {
				defer closer.Close()
			}

			test.function(t, elem)
		})
	}
}

func testElemWithCompoLoadSuccess(t *testing.T, e app.ElemWithCompo) {
	if err := e.Load("tests.hello"); err != nil {
		t.Fatal(err)
	}
}

func testElemWithCompoLoadFail(t *testing.T, e app.ElemWithCompo) {
	err := e.Load("tests.abracadabra")
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testElemWithCompoRenderSuccess(t *testing.T, e app.ElemWithCompo) {
	if err := e.Load("tests.hello"); err != nil {
		t.Fatal(err)
	}

	compo := e.Compo()
	if compo == nil {
		t.Fatal("component is nil")
	}

	hello := compo.(*Hello)
	hello.Name = "Maxence"

	if err := e.Render(hello); err != nil {
		t.Fatal(err)
	}
}

func testElemWithCompoRenderFail(t *testing.T, e app.ElemWithCompo) {
	if err := e.Load("tests.hello"); err != nil {
		t.Fatal(err)
	}

	compo := e.Compo()
	if compo == nil {
		t.Fatal("component is nil")
	}

	hello := compo.(*Hello)
	hello.TmplErr = true

	err := e.Render(hello)
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testElementWithNavigation(t *testing.T, newElem func() (app.Navigator, error)) {
	tests := []struct {
		scenario string
		function func(t *testing.T, elem app.Navigator)
	}{
		{
			scenario: "reload a component",
			function: testElementWithNavigationReloadSuccess,
		},
		{
			scenario: "reload a component fails",
			function: testElementWithNavigationReloadFail,
		},
		{
			scenario: "load previous component",
			function: testElementWithNavigationPreviousSuccess,
		},
		{
			scenario: "load previous component fails",
			function: testElementWithNavigationPreviousFail,
		},
		{
			scenario: "load next component",
			function: testElementWithNavigationNextSuccess,
		},
		{
			scenario: "load next component fails",
			function: testElementWithNavigationNextFail,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			elem, err := newElem()
			if err == app.ErrNotSupported {
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if closer, ok := elem.(app.Closer); ok {
				defer closer.Close()
			}

			test.function(t, elem)
		})
	}
}

func testElementWithNavigationReloadSuccess(t *testing.T, e app.Navigator) {
	if err := e.Load("tests.hello"); err != nil {
		t.Fatal(err)
	}

	if err := e.Reload(); err != nil {
		t.Fatal(err)
	}
}

func testElementWithNavigationReloadFail(t *testing.T, e app.Navigator) {
	err := e.Reload()
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testElementWithNavigationPreviousSuccess(t *testing.T, e app.Navigator) {
	if err := e.Load("tests.hello"); err != nil {
		t.Fatal(err)
	}

	if e.CanPrevious() {
		t.Fatal("can previous is true")
	}

	if err := e.Load("tests.world"); err != nil {
		t.Fatal(err)
	}

	if !e.CanPrevious() {
		t.Fatal("can previous is false")
	}

	if err := e.Previous(); err != nil {
		t.Fatal(err)
	}
}

func testElementWithNavigationPreviousFail(t *testing.T, e app.Navigator) {
	if err := e.Load("tests.hello"); err != nil {
		t.Fatal(err)
	}

	if e.CanPrevious() {
		t.Fatal("can previous is true")
	}

	err := e.Previous()
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testElementWithNavigationNextSuccess(t *testing.T, e app.Navigator) {
	if err := e.Load("tests.hello"); err != nil {
		t.Fatal(err)
	}

	if e.CanNext() {
		t.Fatal("can next is true")
	}

	if err := e.Load("tests.world"); err != nil {
		t.Fatal(err)
	}

	if e.CanNext() {
		t.Fatal("can next is true")
	}

	if err := e.Previous(); err != nil {
		t.Fatal(err)
	}

	if !e.CanNext() {
		t.Fatal("can next is false")
	}

	if err := e.Next(); err != nil {
		t.Fatal(err)
	}
}

func testElementWithNavigationNextFail(t *testing.T, e app.Navigator) {
	if err := e.Load("tests.hello"); err != nil {
		t.Fatal(err)
	}

	if e.CanNext() {
		t.Fatal("can next is true")
	}

	err := e.Next()
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}
