package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestElemSetAttr(t *testing.T) {
	utests := []struct {
		scenario       string
		key            string
		value          interface{}
		explectedValue string
		valueNotSet    bool
	}{
		{
			scenario:       "string",
			key:            "title",
			value:          "test",
			explectedValue: "test",
		},
		{
			scenario:       "int",
			key:            "max",
			value:          42,
			explectedValue: "42",
		},
		{
			scenario:       "bool true",
			key:            "hidden",
			value:          true,
			explectedValue: "",
		},
		{
			scenario:    "bool false",
			key:         "hidden",
			value:       false,
			valueNotSet: true,
		},
		{
			scenario:       "style",
			key:            "style",
			value:          "margin:42",
			explectedValue: "margin:42;",
		},
		{
			scenario:       "set successive styles",
			key:            "style",
			value:          "padding:42",
			explectedValue: "margin:42;padding:42;",
		},
		{
			scenario:       "class",
			key:            "class",
			value:          "hello",
			explectedValue: "hello",
		},
		{
			scenario:       "set successive classes",
			key:            "class",
			value:          "world",
			explectedValue: "hello world",
		},
	}

	e := &elem{}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			e.setAttr(u.key, u.value)
			v, exists := e.attrs[u.key]
			require.Equal(t, u.explectedValue, v)
			require.Equal(t, u.valueNotSet, !exists)
		})
	}
}

func TestElemUpdateAttrs(t *testing.T) {
	utests := []struct {
		scenario string
		current  map[string]string
		incoming map[string]string
	}{
		{
			scenario: "attributes are removed",
			current: map[string]string{
				"foo":   "bar",
				"hello": "world",
			},
			incoming: nil,
		},
		{
			scenario: "attributes are added",
			current:  nil,
			incoming: map[string]string{
				"foo":   "bar",
				"hello": "world",
			},
		},
		{
			scenario: "attributes are updated",
			current: map[string]string{
				"foo":   "bar",
				"hello": "world",
			},
			incoming: map[string]string{
				"foo":   "boo",
				"hello": "there",
			},
		},
		{
			scenario: "attributes are synced",
			current: map[string]string{
				"foo":   "bar",
				"hello": "world",
			},
			incoming: map[string]string{
				"foo":     "boo",
				"goodbye": "world",
			},
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			testSkipNoWasm(t)

			n := Div().(*htmlDiv)
			err := n.mount()
			require.NoError(t, err)
			defer n.dismount()

			n.attrs = u.current
			n.updateAttrs(u.incoming)

			if len(u.incoming) == 0 {
				require.Empty(t, n.attributes())
				return
			}

			require.Equal(t, u.incoming, n.attributes())
		})
	}
}

func TestElemSetEventHandler(t *testing.T) {
	e := &elem{}
	h := func(Context, Event) {}
	e.setEventHandler("click", h)

	expectedHandler := eventHandler{
		event: "click",
		value: h,
	}

	registeredHandler := e.events["click"]
	require.True(t, expectedHandler.equal(registeredHandler))
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
			testSkipNoWasm(t)

			var current map[string]eventHandler
			var incoming map[string]eventHandler

			if u.current != nil {
				current = map[string]eventHandler{
					"click": {
						event: "click",
						value: u.current,
					},
				}
			}

			if u.incoming != nil {
				incoming = map[string]eventHandler{
					"click": {
						event: "click",
						value: u.incoming,
					},
				}
			}

			n := Div().(*htmlDiv)
			n.events = current
			err := n.mount()
			require.NoError(t, err)
			defer n.dismount()

			n.updateEventHandler(incoming)

			if len(incoming) == 0 {
				require.Empty(t, n.attributes())
				return
			}

			h := n.eventHandlers()["click"]
			require.True(t, h.equal(incoming["click"]))
		})
	}
}
