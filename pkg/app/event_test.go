package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEventHandlerEquality(t *testing.T) {
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
				event: "test",
				value: funcA,
			},
			b: eventHandler{
				event: "test",
				value: funcA,
			},
			equals: true,
		},
		{
			scenario: "same event with different func are not equal",
			a: eventHandler{
				event: "test",
				value: funcA,
			},
			b: eventHandler{
				event: "test",
				value: funcB,
			},
			equals: false,
		},
		{
			scenario: "same event with a nil func are not equal",
			a: eventHandler{
				event: "test",
				value: funcA,
			},
			b: eventHandler{
				event: "test",
				value: nil,
			},
			equals: false,
		},
		{
			scenario: "same event with same func and same scope are equal",
			a: eventHandler{
				event: "test",
				scope: "/hello",
				value: funcA,
			},
			b: eventHandler{
				event: "test",
				scope: "/hello",
				value: funcA,
			},
			equals: true,
		},
		{
			scenario: "same event with same func and different scope are not equal",
			a: eventHandler{
				event: "test",
				scope: "/hello",
				value: funcA,
			},
			b: eventHandler{
				event: "test",
				scope: "/bye",
				value: funcA,
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
			event: "test",
			value: funcA,
		}

		b := eventHandler{
			event: "test",
			value: funcB,
		}

		a.Equal(b)
	}
}
