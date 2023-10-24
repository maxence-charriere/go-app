package app

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestObserverObserving(t *testing.T) {
	t.Run("observing when source is mounted", func(t *testing.T) {
		var m nodeManager
		div, err := m.Mount(makeTestContext(), 1, Div())
		require.NoError(t, err)

		o := ObserverX{source: div}
		require.True(t, o.observing())
	})

	t.Run("not observing when source is nil", func(t *testing.T) {
		o := ObserverX{}
		require.False(t, o.observing())
	})

	t.Run("not observing when source is not mounted", func(t *testing.T) {
		o := ObserverX{source: Div()}
		require.False(t, o.observing())
	})

	t.Run("observing when condition is true", func(t *testing.T) {
		var m nodeManager
		div, err := m.Mount(makeTestContext(), 1, Div())
		require.NoError(t, err)

		o := ObserverX{
			source:    div,
			condition: func() bool { return true },
		}
		require.True(t, o.observing())
	})

	t.Run("not observing when condition is false", func(t *testing.T) {
		var m nodeManager
		div, err := m.Mount(makeTestContext(), 1, Div())
		require.NoError(t, err)

		o := ObserverX{
			source:    div,
			condition: func() bool { return false },
		}
		require.False(t, o.observing())
	})
}

func TestStateManagerObserve(t *testing.T) {
	t.Run("observer is set", func(t *testing.T) {
		var nm nodeManager
		compo, err := nm.Mount(makeTestContext(), 1, &hello{})
		require.NoError(t, err)
		ctx := nm.context(makeTestContext(), compo)

		var sm stateManager
		var receiver int
		observer := sm.Observe(ctx, "test", &receiver)
		require.Equal(t, compo, observer.source)
		require.Equal(t, &receiver, observer.receiver)
		require.Nil(t, observer.condition)
		require.Nil(t, observer.changeHandler)
		require.Equal(t, "test", observer.state)
		require.Equal(t, compo, observer.source)
		require.NotNil(t, observer.setObserver)

		require.Len(t, sm.observers, 1)
		require.Len(t, sm.observers["test"], 1)
		observer = sm.observers["test"][compo]
		require.NotZero(t, observer)
		require.Equal(t, compo, observer.source)
		require.Equal(t, &receiver, observer.receiver)
		require.Nil(t, observer.condition)
		require.Nil(t, observer.changeHandler)
		require.Empty(t, observer.state)
		require.Nil(t, observer.setObserver)
	})

	t.Run("observer is set with a while condition", func(t *testing.T) {
		var nm nodeManager
		compo, err := nm.Mount(makeTestContext(), 1, &hello{})
		require.NoError(t, err)
		ctx := nm.context(makeTestContext(), compo)

		var sm stateManager
		var receiver int
		observer := sm.Observe(ctx, "test", &receiver).
			While(func() bool { return true })
		require.Equal(t, compo, observer.source)
		require.Equal(t, &receiver, observer.receiver)
		require.NotNil(t, observer.condition)
		require.Nil(t, observer.changeHandler)
		require.Equal(t, "test", observer.state)
		require.Equal(t, compo, observer.source)
		require.NotNil(t, observer.setObserver)

		require.Len(t, sm.observers, 1)
		require.Len(t, sm.observers["test"], 1)
		observer = sm.observers["test"][compo]
		require.NotZero(t, observer)
		require.Equal(t, compo, observer.source)
		require.Equal(t, &receiver, observer.receiver)
		require.NotNil(t, observer.condition)
		require.Nil(t, observer.changeHandler)
		require.Empty(t, observer.state)
		require.Nil(t, observer.setObserver)
	})

	t.Run("observer is set with a change handler", func(t *testing.T) {
		var nm nodeManager
		compo, err := nm.Mount(makeTestContext(), 1, &hello{})
		require.NoError(t, err)
		ctx := nm.context(makeTestContext(), compo)

		var sm stateManager
		var receiver int
		observer := sm.Observe(ctx, "test", &receiver).
			OnChange(func() {})
		require.Equal(t, compo, observer.source)
		require.Equal(t, &receiver, observer.receiver)
		require.Nil(t, observer.condition)
		require.NotNil(t, observer.changeHandler)
		require.Equal(t, "test", observer.state)
		require.Equal(t, compo, observer.source)
		require.NotNil(t, observer.setObserver)

		require.Len(t, sm.observers, 1)
		require.Len(t, sm.observers["test"], 1)
		observer = sm.observers["test"][compo]
		require.NotZero(t, observer)
		require.Equal(t, compo, observer.source)
		require.Equal(t, &receiver, observer.receiver)
		require.Nil(t, observer.condition)
		require.NotNil(t, observer.changeHandler)
		require.Empty(t, observer.state)
		require.Nil(t, observer.setObserver)
	})

	t.Run("observer is set with a while condition and a change handler", func(t *testing.T) {
		var nm nodeManager
		compo, err := nm.Mount(makeTestContext(), 1, &hello{})
		require.NoError(t, err)
		ctx := nm.context(makeTestContext(), compo)

		var sm stateManager
		var receiver int
		observer := sm.Observe(ctx, "test", &receiver).
			While(func() bool { return true }).
			OnChange(func() {})
		require.Equal(t, compo, observer.source)
		require.Equal(t, &receiver, observer.receiver)
		require.NotNil(t, observer.condition)
		require.NotNil(t, observer.changeHandler)
		require.Equal(t, "test", observer.state)
		require.Equal(t, compo, observer.source)
		require.NotNil(t, observer.setObserver)

		require.Len(t, sm.observers, 1)
		require.Len(t, sm.observers["test"], 1)
		observer = sm.observers["test"][compo]
		require.NotZero(t, observer)
		require.Equal(t, compo, observer.source)
		require.Equal(t, &receiver, observer.receiver)
		require.NotNil(t, observer.condition)
		require.NotNil(t, observer.changeHandler)
		require.Empty(t, observer.state)
		require.Nil(t, observer.setObserver)
	})
}

func TestStateManagerGet(t *testing.T) {
	t.Run("getting a state from memory succeeds", func(t *testing.T) {
		var m stateManager
		ctx := makeTestContext()
		m.Set(ctx, "test", 42)

		var number int
		m.Get(ctx, "test", &number)
		require.Equal(t, 42, number)
	})

	t.Run("getting a state from local storage succeeds", func(t *testing.T) {
		var m stateManager
		ctx := makeTestContext()
		m.Set(ctx, "test", 42).Persist()
		delete(m.states, "test")
		require.Empty(t, m.states)

		var number int
		m.Get(ctx, "test", &number)
		require.Equal(t, 42, number)
	})

	t.Run("getting an encrypted state from local storage succeeds", func(t *testing.T) {
		var m stateManager
		ctx := makeTestContext()
		m.Set(ctx, "test", 42).PersistWithEncryption()
		delete(m.states, "test")
		require.Empty(t, m.states)

		var number int
		m.Get(ctx, "test", &number)
		require.Equal(t, 42, number)
	})

	t.Run("getting an expired state removes the state from state manager and local storage", func(t *testing.T) {
		var m stateManager
		ctx := makeTestContext()
		m.Set(ctx, "test", 42).
			ExpiresIn(-time.Second).
			Persist()
		require.NotEmpty(t, m.states["test"])

		var number int
		m.Get(ctx, "test", &number)
		require.Zero(t, number)
		require.Empty(t, m.states["test"])

		err := ctx.LocalStorage().Get("test", &number)
		require.NoError(t, err)
		require.Zero(t, number)
	})

	t.Run("getting a persisted expired state removes the state from local storage", func(t *testing.T) {
		var m stateManager
		ctx := makeTestContext()
		m.Set(ctx, "test", 42).
			ExpiresIn(-time.Second).
			Persist()
		require.NotEmpty(t, m.states["test"])
		delete(m.states, "test")
		require.Empty(t, m.states)

		var number int
		m.Get(ctx, "test", &number)
		require.Zero(t, number)
		require.Empty(t, m.states["test"])

		err := ctx.LocalStorage().Get("test", &number)
		require.NoError(t, err)
		require.Zero(t, number)
	})

	t.Run("storing a state value into a wrong type logs an error", func(t *testing.T) {
		var m stateManager
		ctx := makeTestContext()
		m.Set(ctx, "test", 42)

		var number string
		m.Get(ctx, "test", &number)
		require.Empty(t, number)
	})

	t.Run("getting a non existing state let receiver with current value", func(t *testing.T) {
		var m stateManager
		ctx := makeTestContext()

		number := 42
		m.Get(ctx, "test", &number)
		require.Equal(t, 42, number)
	})
}

func TestStoreValue(t *testing.T) {
	nb := 42
	c := copyTester{pointer: &nb}

	utests := []struct {
		scenario string
		src      any
		recv     any
		expected any
		err      bool
	}{
		{
			scenario: "value to exported field receiver",
			src:      42,
			recv:     &c.Exported,
			expected: 42,
		},
		{
			scenario: "value unexported field receiver",
			src:      21,
			recv:     &c.unexported,
			expected: 21,
		},
		{
			scenario: "nil to receiver",
			src:      nil,
			recv:     &c.unexported,
			expected: 0,
		},
		{
			scenario: "pointer to receiver",
			src:      new(int),
			recv:     &c.unexported,
			expected: 0,
		},
		{
			scenario: "nil to pointer receiver",
			src:      nil,
			recv:     &c.pointer,
			expected: (*int)(nil),
		},
		{
			scenario: "slice to receiver",
			src:      []int{14, 2, 86},
			recv:     &c.slice,
			expected: []int{14, 2, 86},
		},
		{
			scenario: "map to receiver",
			src:      map[string]int{"foo": 42},
			recv:     &c.mapp,
			expected: map[string]int{"foo": 42},
		},
		{
			scenario: "receiver have a different type",
			src:      "hello",
			recv:     &c.Exported,
			err:      true,
		},
		{
			scenario: "receiver is not a pointer",
			src:      51,
			recv:     c.Exported,
			err:      true,
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			err := storeValue(u.recv, u.src)
			if u.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			recv := reflect.ValueOf(u.recv).Elem().Interface()
			require.Equal(t, u.expected, recv)
		})
	}
}

type copyTester struct {
	Exported   int
	unexported int
	pointer    *int
	slice      []int
	mapp       map[string]int
}
