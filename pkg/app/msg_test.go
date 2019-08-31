package app

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMessenger(t *testing.T) {
	foo := &Foo{}
	m := messenger{}
	bindingCalled := false

	bind, close := m.bind("test", foo)
	require.Len(t, m.bindings, 1)
	require.Len(t, m.bindings["test"], 1)

	bind.Do(func(ctx context.Context) {
		bindingCalled = true
	})

	m.emit(context.TODO(), "test")
	require.True(t, bindingCalled)

	close()
	require.Len(t, m.bindings, 1)
	require.Empty(t, m.bindings["test"])
}

func TestBindingDo(t *testing.T) {
	tests := []struct {
		scenario string
		function interface{}
		panic    bool
	}{
		{
			scenario: "function is added to binding",
			function: func(context.Context) {},
		},
		{
			scenario: "function with args is added to binding",
			function: func(context.Context, int, bool) {},
		},
		{
			scenario: "non function added to binding panics",
			function: 42,
			panic:    true,
		},
		{
			scenario: "function without context added to binding panics",
			function: func() {},
			panic:    true,
		},
		{
			scenario: "function without context as 1st arg added to binding panics",
			function: func(int) {},
			panic:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			b := Binding{msg: "test"}

			if test.panic {
				require.Panics(t, func() { b.Do(test.function) })
				return
			}

			b.Do(test.function)
			require.Len(t, b.funcs, 1)
		})
	}
}
