package app

import "testing"

func TestTextMountDismout(t *testing.T) {
	testMountDismount(t, []mountTest{
		{
			scenario: "text",
			node:     Text("hello"),
		},
		{
			scenario: "textf",
			node:     Textf("hello %s", "Maxence"),
		},
	})
}

func TestTextUpdate(t *testing.T) {
	testUpdate(t, []updateTest{
		{
			scenario: "text element is updated",
			a:        Text("hello"),
			b:        Text("world"),
			matches: []TestUIDescriptor{
				{
					Expected: Text("world"),
				},
			},
		},

		{
			scenario: "text is replaced by a html elem",
			a: Div().Body(
				Text("hello"),
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
			scenario: "text is replaced by a component",
			a: Div().Body(
				Text("hello"),
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
			scenario: "text is replaced by a raw html element",
			a: Div().Body(
				Text("hello"),
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
