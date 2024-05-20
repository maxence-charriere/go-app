package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEventHandlersSet(t *testing.T) {
	t.Run("nil event handler is not set", func(t *testing.T) {
		eventHandlers := make(eventHandlers)
		eventHandlers.Set("click", nil)
		require.Empty(t, eventHandlers)
	})

	t.Run("event handler is set", func(t *testing.T) {
		eventHandlers := make(eventHandlers)
		eventHandlers.Set("click", func(ctx Context, e Event) {})
		require.NotZero(t, eventHandlers["click"])
	})
}

func TestMakeEventHandler(t *testing.T) {
	t.Run("without options", func(t *testing.T) {
		eh := makeEventHandler("click", func(ctx Context, e Event) {})
		require.Equal(t, "click", eh.event)
		require.Empty(t, eh.scope)
		require.False(t, eh.passive)
		require.NotNil(t, eh.goHandler)
		require.Nil(t, eh.jsHandler)
	})

	t.Run("with scope option", func(t *testing.T) {
		eh := makeEventHandler("click", func(ctx Context, e Event) {}, EventScope(1))
		require.Equal(t, "click", eh.event)
		require.Equal(t, "/1", eh.scope)
		require.False(t, eh.passive)
		require.NotNil(t, eh.goHandler)
		require.Nil(t, eh.jsHandler)
	})

	t.Run("with passive option", func(t *testing.T) {
		eh := makeEventHandler("click", func(ctx Context, e Event) {}, PassiveEvent())
		require.Equal(t, "click", eh.event)
		require.Empty(t, eh.scope)
		require.True(t, eh.passive)
		require.NotNil(t, eh.goHandler)
		require.Nil(t, eh.jsHandler)
	})

}

func TestEventHandlerEqual(t *testing.T) {
	funcA := func(Context, Event) {}
	funcB := func(Context, Event) {}

	utests := []struct {
		scenario string
		a        eventHandler
		b        eventHandler
		equals   bool
	}{
		{
			scenario: "same event with same func are equal",
			a: eventHandler{
				event:     "test",
				goHandler: funcA,
			},
			b: eventHandler{
				event:     "test",
				goHandler: funcA,
			},
			equals: true,
		},
		{
			scenario: "same event with different func are not equal",
			a: eventHandler{
				event:     "test",
				goHandler: funcA,
			},
			b: eventHandler{
				event:     "test",
				goHandler: funcB,
			},
			equals: false,
		},
		{
			scenario: "same event with a nil func are not equal",
			a: eventHandler{
				event:     "test",
				goHandler: funcA,
			},
			b: eventHandler{
				event:     "test",
				goHandler: nil,
			},
			equals: false,
		},
		{
			scenario: "same event with same func and same scope are equal",
			a: eventHandler{
				event:     "test",
				scope:     "/hello",
				goHandler: funcA,
			},
			b: eventHandler{
				event:     "test",
				scope:     "/hello",
				goHandler: funcA,
			},
			equals: true,
		},
		{
			scenario: "same event with same func and different scope are not equal",
			a: eventHandler{
				event:     "test",
				scope:     "/hello",
				goHandler: funcA,
			},
			b: eventHandler{
				event:     "test",
				scope:     "/bye",
				goHandler: funcA,
			},
			equals: false,
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			require.Equal(t, u.equals, u.a.Equal(u.b))
		})
	}
}

func BenchmarkEventHandlerEquality(b *testing.B) {
	funcA := func(Context, Event) {}
	funcB := func(Context, Event) {}

	for n := 0; n < b.N; n++ {
		a := eventHandler{
			event:     "test",
			goHandler: funcA,
		}

		b := eventHandler{
			event:     "test",
			goHandler: funcB,
		}

		a.Equal(b)
	}
}
