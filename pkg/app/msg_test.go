package app

import (
	"context"
	"testing"

	"github.com/pkg/errors"
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

func TestBindingExec(t *testing.T) {
	tests := []struct {
		scenario  string
		args      []interface{}
		functions []interface{}
	}{
		{
			scenario: "execute function with matching args",
			args:     []interface{}{"hello", 42},
			functions: []interface{}{
				func(ctx context.Context, s string, i int) {
					require.Equal(t, "hello", s)
					require.Equal(t, 42, i)
				},
			},
		},
		{
			scenario: "execute function with less matching args",
			args:     []interface{}{"hello", 42},
			functions: []interface{}{
				func(ctx context.Context, s string) {
					require.Equal(t, "hello", s)
				},
			},
		},
		{
			scenario: "execute function with more matching args",
			args:     []interface{}{"hello"},
			functions: []interface{}{
				func(ctx context.Context, s string, i int) {
					require.Equal(t, "hello", s)
					require.Equal(t, 0, i)
				},
			},
		},
		{
			scenario: "execute function with non matching args",
			args:     []interface{}{"hello", 42},
			functions: []interface{}{
				func(ctx context.Context, s string, i int32) {
					require.Equal(t, "hello", s)
					require.Equal(t, 42, i)
				},
			},
		},
		{
			scenario: "execute multiple functions with matching args",
			args:     []interface{}{"hello", 42},
			functions: []interface{}{
				func(ctx context.Context, s string, i int) (bool, error) {
					return true, nil
				},
				func(ctx context.Context, b bool, err error) error {
					require.True(t, b)
					require.NoError(t, err)
					return errors.New("test")
				},
				func(ctx context.Context, err error) {
					require.Error(t, err)
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			b := Binding{msg: "test"}

			for _, f := range test.functions {
				b.Do(f)
			}

			b.exec(context.TODO(), test.args...)
		})
	}
}
