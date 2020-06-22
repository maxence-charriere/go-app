package app

import (
	"net/url"
	"testing"
)

func TestCompoMountDismount(t *testing.T) {
	testMountDismount(t, []mountTest{
		{
			scenario: "component",
			node:     &hello{},
		},
	})
}

func TestCompoUpdate(t *testing.T) {
	testUpdate(t, []updateTest{
		{
			scenario: "component is updated",
			a:        &bar{Value: "rab"},
			b:        &bar{Value: "bar"},
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: &bar{Value: "bar"},
				},
				{
					Path:     TestPath(0),
					Expected: Text("bar"),
				},
			},
		},
		{
			scenario:   "component returns replace error when updated with a non component-element",
			a:          &hello{},
			b:          Text("hello"),
			replaceErr: true,
		},
		{
			scenario: "component is updated",
			a:        &hello{},
			b:        &hello{Greeting: "world"},
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: &hello{Greeting: "world"},
				},
				{
					Path:     TestPath(0),
					Expected: Div(),
				},
				{
					Path:     TestPath(0, 0),
					Expected: H1(),
				},
				{
					Path:     TestPath(0, 0, 0),
					Expected: Text("hello, "),
				},
				{
					Path:     TestPath(0, 0, 1),
					Expected: Text("world"),
				},
			},
		},
		{
			scenario: "component is replaced by a text",
			a: Div().Body(
				&hello{},
			),
			b: Div().Body(
				Text("hello"),
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},
				{
					Path:     TestPath(0),
					Expected: Text("hello"),
				},
			},
		},
		{
			scenario: "component is replaced by an html element",
			a: Div().Body(
				&hello{},
			),
			b: Div().Body(
				H2().Text("hello"),
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},
				{
					Path:     TestPath(0),
					Expected: H2(),
				},
				{
					Path:     TestPath(0, 0),
					Expected: Text("hello"),
				},
			},
		},
		{
			scenario: "component is replaced by a raw html element",
			a: Div().Body(
				&hello{},
			),
			b: Div().Body(
				Raw("<svg></svg>"),
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},
				{
					Path:     TestPath(0),
					Expected: Raw("<svg></svg>"),
				},
			},
		},
		{
			scenario: "component is replaced by another component",
			a: Div().Body(
				&hello{},
			),
			b: Div().Body(
				&bar{},
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},
				{
					Path:     TestPath(0),
					Expected: &bar{},
				},
				{
					Path:     TestPath(0, 0),
					Expected: Text(""),
				},
			},
		},
		{
			scenario: "component root is updated",
			a: Div().Body(
				&foo{Bar: "hello"},
			),
			b: Div().Body(
				&foo{Bar: "goodbye"},
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},
				{
					Path:     TestPath(0),
					Expected: &foo{Bar: "goodbye"},
				},
				{
					Path:     TestPath(0, 0),
					Expected: &bar{Value: "goodbye"},
				},
				{
					Path:     TestPath(0, 0, 0),
					Expected: Text("goodbye"),
				},
			},
		},
		{
			scenario: "component root is replaced by a component",
			a: Div().Body(
				&foo{},
			),
			b: Div().Body(
				&foo{Bar: "test"},
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},
				{
					Path:     TestPath(0),
					Expected: &foo{Bar: "test"},
				},
				{
					Path:     TestPath(0, 0),
					Expected: &bar{Value: "test"},
				},
				{
					Path:     TestPath(0, 0, 0),
					Expected: Text("test"),
				},
			},
		},
		{
			scenario: "component root is replaced by a non-component",
			a: Div().Body(
				&foo{Bar: "test"},
			),
			b: Div().Body(
				&foo{},
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},
				{
					Path:     TestPath(0),
					Expected: &foo{},
				},
				{
					Path:     TestPath(0, 0),
					Expected: Text("bar"),
				},
			},
		},
	})
}

type hello struct {
	Compo

	Greeting string
}

func (h *hello) OnMount(Context) {
}

func (h *hello) OnNav(Context, *url.URL) {
}

func (h *hello) OnDismount(Context) {
}

func (h *hello) Render() UI {
	return Div().Body(
		H1().Body(
			Text("hello, "),
			Text(h.Greeting),
		),
	)
}

type foo struct {
	Compo
	Bar string
}

func (f *foo) Render() UI {
	return If(f.Bar != "",
		&bar{Value: f.Bar},
	).Else(
		Text("bar"),
	)
}

type bar struct {
	Compo
	Value string
}

func (b *bar) Render() UI {
	return Text(b.Value)
}
