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
		callOnUI: func(f func()) { f() },
		callExec: func(f func(a ...interface{}), a ...interface{}) {
			f(a...)
		},
	}

	foo := &Foo{}
	bindingCalled := false

	bind, close := m.bind("test", foo)
	require.Len(t, m.bindings, 1)
	require.Len(t, m.bindings["test"], 1)

	bind.Do(func() {
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
		actions  []actionTest
	}{
		{
			scenario: "execute function with matching args",
			args:     []interface{}{"hello", 42},
			actions: []actionTest{
				doTest(func(s string, i int) {
					require.Equal(t, "hello", s)
					require.Equal(t, 42, i)
				}),
			},
		},
		{
			scenario: "execute function with matching bind context and args",
			args:     []interface{}{"hello", 42},
			actions: []actionTest{
				doTest(func(ctx BindContext, s string, i int) {
					require.NotNil(t, ctx)
					require.Equal(t, "hello", s)
					require.Equal(t, 42, i)
				}),
			},
		},
		{
			scenario: "execute function with matching context and args",
			args:     []interface{}{"hello", 42},
			actions: []actionTest{
				doTest(func(ctx context.Context, s string, i int) {
					require.NotNil(t, ctx)
					require.Equal(t, "hello", s)
					require.Equal(t, 42, i)
				}),
			},
		},
		{
			scenario: "execute async function with matching args on ui",
			args:     []interface{}{"hello", 42},
			actions: []actionTest{
				doAsyncTest(func(s string, i int) {
					require.Equal(t, "hello", s)
					require.Equal(t, 42, i)
				}),
			},
		},
		{
			scenario: "execute function with wait and matching args",
			args:     []interface{}{"hello", 42},
			actions: []actionTest{
				waitTest(time.Millisecond),
				doTest(func(s string, i int) {
					require.Equal(t, "hello", s)
					require.Equal(t, 42, i)
				}),
			},
		},
		{
			scenario: "execute function with negative wait and matching args",
			args:     []interface{}{"hello", 42},
			actions: []actionTest{
				waitTest(-time.Millisecond),
				doTest(func(s string, i int) {
					require.Equal(t, "hello", s)
					require.Equal(t, 42, i)
				}),
			},
		},
		{
			scenario: "execute function with less matching args",
			args:     []interface{}{"hello", 42},
			actions: []actionTest{
				doTest(func(s string) {
					require.Equal(t, "hello", s)
				}),
			},
		},
		{
			scenario: "execute function with more matching args",
			args:     []interface{}{"hello"},
			actions: []actionTest{
				doTest(func(s string, i int) {
					require.Equal(t, "hello", s)
					require.Equal(t, 0, i)
				}),
			},
		},
		{
			scenario: "execute function with non matching args",
			args:     []interface{}{"hello", 42},
			actions: []actionTest{
				doTest(func(s string, i int32) {
					require.Equal(t, "hello", s)
					require.Equal(t, 42, i)
				}),
			},
		},
		{
			scenario: "execute multiple functions",
			args:     []interface{}{"hello", 42},
			actions: []actionTest{
				doTest(func(s string, i int) (bool, error) {
					return true, nil
				}),
				doTest(func(ctx BindContext) {
					require.NotNil(t, ctx)
				}),
				doTest(func() {
					require.True(t, true)
				}),
			},
		},
		{
			scenario: "execute multiple functions and wait",
			args:     []interface{}{"hello", 42},
			actions: []actionTest{
				doAsyncTest(func(s string, i int) (bool, error) {
					return true, nil
				}),
				waitTest(time.Millisecond),
				doAsyncTest(func(ctx BindContext) {
					require.NotNil(t, ctx)
				}),
				waitTest(time.Millisecond),
				doTest(func() {
					require.True(t, true)
				}),
			},
		},
		{
			scenario: "execute and cancel",
			args:     []interface{}{"hello", 42},
			actions: []actionTest{
				doTest(func(ctx BindContext) {
					ctx.Cancel(nil)
				}),
				doTest(func() {
					t.Fatal("action called in cancelled binding")
				}),
				whenCancelTest(func(ctx BindContext) {
					require.Equal(t, context.Canceled, ctx.Err())
				}),
			},
		},
		{
			scenario: "execute and cancel with err",
			args:     []interface{}{"hello", 42},
			actions: []actionTest{
				doTest(func(ctx BindContext) {
					ctx.Cancel(errors.New("test cancel"))
				}),
				doTest(func() {
					t.Fatal("action called in cancelled binding")
				}),
				whenCancelTest(func(ctx BindContext) {
					require.Error(t, ctx.Err())
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
				switch a.name {
				case "Wait":
					b.Wait(a.duration)

				case "Do":
					b.Do(a.function)

				case "DoAsync":
					b.DoAsync(a.function)

				case "WhenCancel":
					b.WhenCancel(a.whenCancel)
				}
			}

			b.exec(test.args...)
		})
	}
}

type actionTest struct {
	name       string
	function   interface{}
	whenCancel func(BindContext)
	duration   time.Duration
}

func doTest(f interface{}) actionTest {
	return actionTest{
		name:     "Do",
		function: f,
	}
}

func doAsyncTest(f interface{}) actionTest {
	return actionTest{
		name:     "DoAsync",
		function: f,
	}
}

func waitTest(d time.Duration) actionTest {
	return actionTest{
		name:     "Wait",
		duration: d,
	}
}

func whenCancelTest(f func(BindContext)) actionTest {
	return actionTest{
		name:       "WhenCancel",
		whenCancel: f,
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

func TestBindingContext(t *testing.T) {
	b := newBindContext()
	require.Nil(t, b.Get("test"))

	_, exists := b.Lookup("test")
	require.False(t, exists)

	b.Set("test", 42)
	require.Equal(t, 42, b.Get("test"))
	v, exists := b.Lookup("test")
	require.Equal(t, 42, v)
	require.True(t, exists)

	require.NoError(t, b.Err())
	b.Cancel(nil)
	require.Equal(t, context.Canceled, b.Err())
	b.Cancel(errors.New("test"))
	require.Error(t, b.Err())
}
