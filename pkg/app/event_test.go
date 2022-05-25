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
