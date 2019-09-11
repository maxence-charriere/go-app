package app

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestMessenger(t *testing.T) {
	m := messenger{
		callExec: func(f func(a ...interface{}), a ...interface{}) {
			f(a...)
		},
	}

	foo := &Foo{}
	bindingCalled := false

	bind, close := m.bind("test", foo)
	require.Len(t, m.bindings, 1)
	require.Len(t, m.bindings["test"], 1)

	bind.Do(func(ctx context.Context) {
		bindingCalled = true
	})

	m.emit("test")
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
			function: func() {},
		},
		{
			scenario: "function with args is added to binding",
			function: func(int, bool) {},
		},
		{
			scenario: "non function added to binding panics",
			function: 42,
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
			require.Len(t, b.actions, 1)
		})
	}
}

func TestBindingExec(t *testing.T) {
	tests := []struct {
		scenario string
		args     []interface{}
		actions  []interface{}
	}{
		{
			scenario: "execute function with matching args",
			args:     []interface{}{"hello", 42},
			actions: []interface{}{
				makeDoTest(false, func(s string, i int) {
					require.Equal(t, "hello", s)
					require.Equal(t, 42, i)
				}),
			},
		},
		{
			scenario: "execute function with matching args on ui",
			args:     []interface{}{"hello", 42},
			actions: []interface{}{
				makeDoTest(true, func(s string, i int) {
					require.Equal(t, "hello", s)
					require.Equal(t, 42, i)
				}),
			},
		},
		{
			scenario: "execute function with wait and matching args",
			args:     []interface{}{"hello", 42},
			actions: []interface{}{
				time.Millisecond,
				makeDoTest(false, func(s string, i int) {
					require.Equal(t, "hello", s)
					require.Equal(t, 42, i)
				}),
			},
		},
		{
			scenario: "execute function with negative wait and matching args",
			args:     []interface{}{"hello", 42},
			actions: []interface{}{
				-time.Millisecond,
				makeDoTest(false, func(s string, i int) {
					require.Equal(t, "hello", s)
					require.Equal(t, 42, i)
				}),
			},
		},
		{
			scenario: "execute function with less matching args",
			args:     []interface{}{"hello", 42},
			actions: []interface{}{
				makeDoTest(false, func(s string) {
					require.Equal(t, "hello", s)
				}),
			},
		},
		{
			scenario: "execute function with more matching args",
			args:     []interface{}{"hello"},
			actions: []interface{}{
				makeDoTest(false, func(s string, i int) {
					require.Equal(t, "hello", s)
					require.Equal(t, 0, i)
				}),
			},
		},
		{
			scenario: "execute function with non matching args",
			args:     []interface{}{"hello", 42},
			actions: []interface{}{
				makeDoTest(false, func(s string, i int32) {
					require.Equal(t, "hello", s)
					require.Equal(t, 42, i)
				}),
			},
		},
		{
			scenario: "execute multiple functions with matching args",
			args:     []interface{}{"hello", 42},
			actions: []interface{}{
				makeDoTest(false, func(s string, i int) (bool, error) {
					return true, nil
				}),
				makeDoTest(false, func(b bool, err error) error {
					require.True(t, b)
					require.NoError(t, err)
					return errors.New("test")
				}),
				makeDoTest(false, func(err error) {
					require.Error(t, err)
				}),
			},
		},
		{
			scenario: "execute multiple functions with matching args and wait",
			args:     []interface{}{"hello", 42},
			actions: []interface{}{
				makeDoTest(false, func(s string, i int) (bool, error) {
					return true, nil
				}),
				time.Millisecond,
				makeDoTest(true, func(b bool, err error) error {
					require.True(t, b)
					require.NoError(t, err)
					return errors.New("test")
				}),
				time.Millisecond,
				makeDoTest(false, func(err error) {
					require.Error(t, err)
				}),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			b := Binding{
				msg:      "test",
				callOnUI: func(f func()) { f() },
			}

			for _, a := range test.actions {
				switch i := a.(type) {
				case time.Duration:
					b.Wait(i)

				case doTest:
					if i.callOnUI {
						b.DoOnUI(i.function)
					} else {
						b.Do(i.function)
					}

				}
			}

			b.exec(test.args...)
		})
	}
}

type doTest struct {
	callOnUI bool
	function interface{}
}

func makeDoTest(callOnUI bool, f interface{}) doTest {
	return doTest{
		callOnUI: callOnUI,
		function: f,
	}
}

func TestBindingDefer(t *testing.T) {
	b := Binding{
		callOnUI: func(f func()) { f() },
	}

	b.Defer(time.Millisecond * 200).
		Do(func(i int) {
			require.Equal(t, 99, i)
		})

	for i := 0; i < 100; i++ {
		go b.exec(i)
		time.Sleep(time.Millisecond)
	}
}
