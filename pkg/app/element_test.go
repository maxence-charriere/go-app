package app

import (
	"testing"
)

func TestElemMountDismount(t *testing.T) {
	testMountDismount(t, []mountTest{
		{
			scenario: "html element",
			node: Div().
				Class("hello").
				OnClick(func(Context, Event) {}),
		},
	})
}

func TestElemUpdate(t *testing.T) {
	testUpdate(t, []updateTest{
		{
			scenario: "html element attributes are updated",
			a: Div().
				ID("max").
				Class("foo").
				AccessKey("test"),
			b: Div().
				ID("max").
				Class("bar").
				Lang("fr"),
			matches: []TestUIDescriptor{
				{
					Expected: Div().
						ID("max").
						Class("bar").
						Lang("fr"),
				},
			},
		},
		{
			scenario: "html element event handlers are added",
			a:        Div(),
			b: Div().
				OnClick(func(Context, Event) {}).
				OnChange(func(Context, Event) {}),
			matches: []TestUIDescriptor{
				{
					Expected: Div().
						OnClick(func(Context, Event) {}).
						OnChange(func(Context, Event) {}),
				},
			},
		},
		{
			scenario: "html element event handlers are updated",
			a: Div().
				OnClick(func(Context, Event) {}).
				OnBlur(func(Context, Event) {}),
			b: Div().
				OnClick(func(Context, Event) {}).
				OnChange(func(Context, Event) {}),
			matches: []TestUIDescriptor{
				{
					Expected: Div().
						OnClick(func(Context, Event) {}).
						OnChange(func(Context, Event) {}),
				},
			},
		},
		{
			scenario: "html element event handlers are removed",
			a: Div().
				OnClick(func(Context, Event) {}).
				OnBlur(func(Context, Event) {}),
			b: Div(),
			matches: []TestUIDescriptor{
				{
					Expected: Div(),
				},
			},
		},
		{
			scenario: "html element is replaced by a text",
			a: Div().Body(
				H2().Text("hello"),
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
			scenario: "html element is replaced by a component",
			a: Div().Body(
				H2().Text("hello"),
			),
			b: Div().Body(
				&hello{},
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},
				{
					Path:     TestPath(0),
					Expected: &hello{},
				},
				{
					Path:     TestPath(0, 0, 0),
					Expected: H1(),
				},
				{
					Path:     TestPath(0, 0, 0, 0),
					Expected: Text("hello, "),
				},
			},
		},
		{
			scenario: "html element is replaced by another html element",
			a: Div().Body(
				H2(),
			),
			b: Div().Body(
				H1(),
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},
				{
					Path:     TestPath(0),
					Expected: H1(),
				},
			},
		},
		{
			scenario: "html element is replaced by raw html element",
			a: Div().Body(
				H2().Text("hello"),
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
	})
}
