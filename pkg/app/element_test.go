package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestElemSetEventHandler(t *testing.T) {
	e := &elem{}
	h := func(Context, Event) {}
	e.setEventHandler("click", h)

	expectedHandler := eventHandler{
		event:     "click",
		goHandler: h,
	}

	registeredHandler := e.events["click"]
	require.True(t, expectedHandler.Equal(registeredHandler))
}

func TestElemUpdateEventHandlers(t *testing.T) {
	utests := []struct {
		scenario string
		current  EventHandler
		incoming EventHandler
	}{
		{
			scenario: "handler is removed",
			current:  func(Context, Event) {},
			incoming: nil,
		},
		{
			scenario: "handler is added",
			current:  nil,
			incoming: func(Context, Event) {},
		},
		{
			scenario: "handler is updated",
			current:  func(Context, Event) {},
			incoming: func(Context, Event) {},
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			var current map[string]eventHandler
			var incoming map[string]eventHandler

			if u.current != nil {
				current = map[string]eventHandler{
					"click": {
						event:     "click",
						goHandler: u.current,
					},
				}
			}

			if u.incoming != nil {
				incoming = map[string]eventHandler{
					"click": {
						event:     "click",
						goHandler: u.incoming,
					},
				}
			}

			n := Div().(*htmlDiv)
			n.events = current

			d := NewClientTester(n)
			defer d.Close()

			d.Consume()

			n.updateEventHandler(incoming)

			if len(incoming) == 0 {
				require.Empty(t, n.getAttributes())
				return
			}

			h := n.getEventHandlers()["click"]
			require.True(t, h.Equal(incoming["click"]))
		})
	}
}

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
						OnClick(nil).
						OnChange(nil),
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

func TestIsURLAttrValue(t *testing.T) {
	utests := []struct {
		name     string
		expected bool
	}{
		{
			name:     "cite",
			expected: true,
		},
		{
			name:     "data",
			expected: true,
		},
		{
			name:     "href",
			expected: true,
		},
		{
			name:     "src",
			expected: true,
		},
		{
			name:     "srcset",
			expected: true,
		},
		{
			name:     "data-test",
			expected: false,
		},
	}

	for _, u := range utests {
		t.Run(u.name, func(t *testing.T) {
			require.Equal(t, u.expected, isURLAttrValue(u.name))
		})
	}
}
