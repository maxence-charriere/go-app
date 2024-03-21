package app

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestObserverObserving(t *testing.T) {
	t.Run("observing when source is mounted", func(t *testing.T) {
		var m nodeManager
		div, err := m.Mount(makeTestContext(), 1, Div())
		require.NoError(t, err)

		o := Observer{source: div}
		require.True(t, o.observing())
	})

	t.Run("not observing when source is nil", func(t *testing.T) {
		o := Observer{}
		require.False(t, o.observing())
	})

	t.Run("not observing when source is not mounted", func(t *testing.T) {
		o := Observer{source: Div()}
		require.False(t, o.observing())
	})

	t.Run("observing when condition is true", func(t *testing.T) {
		var m nodeManager
		div, err := m.Mount(makeTestContext(), 1, Div())
		require.NoError(t, err)

		o := Observer{
			source:    div,
			condition: func() bool { return true },
		}
		require.True(t, o.observing())
	})

	t.Run("not observing when condition is false", func(t *testing.T) {
		var m nodeManager
		div, err := m.Mount(makeTestContext(), 1, Div())
		require.NoError(t, err)

		o := Observer{
			source:    div,
			condition: func() bool { return false },
		}
		require.False(t, o.observing())
	})
}

func TestStateManagerObserve(t *testing.T) {
	t.Run("observer is set", func(t *testing.T) {
		stateName := uuid.NewString()

		var nm nodeManager
		compo, err := nm.Mount(makeTestContext(), 1, &hello{})
		require.NoError(t, err)
		ctx := nm.context(makeTestContext(), compo)

		var sm stateManager
		var receiver int
		observer := sm.Observe(ctx, stateName, &receiver)
		require.Equal(t, compo, observer.source)
		require.Equal(t, &receiver, observer.receiver)
		require.Nil(t, observer.condition)
		require.Nil(t, observer.changeHandler)
		require.Equal(t, stateName, observer.state)
		require.Equal(t, compo, observer.source)
		require.NotNil(t, observer.setObserver)

		require.Len(t, sm.observers, 1)
		require.Len(t, sm.observers[stateName], 1)
		observer = sm.observers[stateName][compo]
		require.NotZero(t, observer)
		require.Equal(t, compo, observer.source)
		require.Equal(t, &receiver, observer.receiver)
		require.Nil(t, observer.condition)
		require.Nil(t, observer.changeHandler)
		require.Empty(t, observer.state)
		require.Nil(t, observer.setObserver)
	})

	t.Run("observer is set with a while condition", func(t *testing.T) {
		stateName := uuid.NewString()

		var nm nodeManager
		compo, err := nm.Mount(makeTestContext(), 1, &hello{})
		require.NoError(t, err)
		ctx := nm.context(makeTestContext(), compo)

		var sm stateManager
		var receiver int
		observer := sm.Observe(ctx, stateName, &receiver).
			While(func() bool { return true })
		require.Equal(t, compo, observer.source)
		require.Equal(t, &receiver, observer.receiver)
		require.NotNil(t, observer.condition)
		require.Nil(t, observer.changeHandler)
		require.Equal(t, stateName, observer.state)
		require.Equal(t, compo, observer.source)
		require.NotNil(t, observer.setObserver)

		require.Len(t, sm.observers, 1)
		require.Len(t, sm.observers[stateName], 1)
		observer = sm.observers[stateName][compo]
		require.NotZero(t, observer)
		require.Equal(t, compo, observer.source)
		require.Equal(t, &receiver, observer.receiver)
		require.NotNil(t, observer.condition)
		require.Nil(t, observer.changeHandler)
		require.Empty(t, observer.state)
		require.Nil(t, observer.setObserver)
	})

	t.Run("observer is set with a change handler", func(t *testing.T) {
		stateName := uuid.NewString()

		var nm nodeManager
		compo, err := nm.Mount(makeTestContext(), 1, &hello{})
		require.NoError(t, err)
		ctx := nm.context(makeTestContext(), compo)

		var sm stateManager
		var receiver int
		observer := sm.Observe(ctx, stateName, &receiver).
			OnChange(func() {})
		require.Equal(t, compo, observer.source)
		require.Equal(t, &receiver, observer.receiver)
		require.Nil(t, observer.condition)
		require.NotNil(t, observer.changeHandler)
		require.Equal(t, stateName, observer.state)
		require.Equal(t, compo, observer.source)
		require.NotNil(t, observer.setObserver)

		require.Len(t, sm.observers, 1)
		require.Len(t, sm.observers[stateName], 1)
		observer = sm.observers[stateName][compo]
		require.NotZero(t, observer)
		require.Equal(t, compo, observer.source)
		require.Equal(t, &receiver, observer.receiver)
		require.Nil(t, observer.condition)
		require.NotNil(t, observer.changeHandler)
		require.Empty(t, observer.state)
		require.Nil(t, observer.setObserver)
	})

	t.Run("observer is set with a while condition and a change handler", func(t *testing.T) {
		stateName := uuid.NewString()

		var nm nodeManager
		compo, err := nm.Mount(makeTestContext(), 1, &hello{})
		require.NoError(t, err)
		ctx := nm.context(makeTestContext(), compo)

		var sm stateManager
		var receiver int
		observer := sm.Observe(ctx, stateName, &receiver).
			While(func() bool { return true }).
			OnChange(func() {})
		require.Equal(t, compo, observer.source)
		require.Equal(t, &receiver, observer.receiver)
		require.NotNil(t, observer.condition)
		require.NotNil(t, observer.changeHandler)
		require.Equal(t, stateName, observer.state)
		require.Equal(t, compo, observer.source)
		require.NotNil(t, observer.setObserver)

		require.Len(t, sm.observers, 1)
		require.Len(t, sm.observers[stateName], 1)
		observer = sm.observers[stateName][compo]
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
		stateName := uuid.NewString()

		var m stateManager
		ctx := makeTestContext()
		m.Set(ctx, stateName, 42)

		var number int
		m.Get(ctx, stateName, &number)
		require.Equal(t, 42, number)
	})

	t.Run("getting a state from local storage succeeds", func(t *testing.T) {
		stateName := uuid.NewString()

		var m stateManager
		ctx := makeTestContext()
		m.Set(ctx, stateName, 42).Persist()
		delete(m.states, stateName)
		require.Empty(t, m.states)

		var number int
		m.Get(ctx, stateName, &number)
		require.Equal(t, 42, number)
	})

	t.Run("getting an encrypted state from local storage succeeds", func(t *testing.T) {
		stateName := uuid.NewString()

		var m stateManager
		ctx := makeTestContext()
		m.Set(ctx, stateName, 42).PersistWithEncryption()
		delete(m.states, stateName)
		require.Empty(t, m.states)

		var number int
		m.Get(ctx, stateName, &number)
		require.Equal(t, 42, number)
	})

	t.Run("getting an expired state removes the state from state manager and local storage", func(t *testing.T) {
		stateName := uuid.NewString()

		var m stateManager
		ctx := makeTestContext()
		m.Set(ctx, stateName, 42).
			ExpiresIn(-time.Second).
			Persist()
		require.NotEmpty(t, m.states[stateName])

		var number int
		m.Get(ctx, stateName, &number)
		require.Zero(t, number)
		require.Empty(t, m.states[stateName])

		err := ctx.LocalStorage().Get(stateName, &number)
		require.NoError(t, err)
		require.Zero(t, number)
	})

	t.Run("getting a persisted expired state removes the state from local storage", func(t *testing.T) {
		stateName := uuid.NewString()

		var m stateManager
		ctx := makeTestContext()
		m.Set(ctx, stateName, 42).
			ExpiresIn(-time.Second).
			Persist()
		require.NotEmpty(t, m.states[stateName])
		delete(m.states, stateName)
		require.Empty(t, m.states)

		var number int
		m.Get(ctx, stateName, &number)
		require.Zero(t, number)
		require.Empty(t, m.states[stateName])

		err := ctx.LocalStorage().Get(stateName, &number)
		require.NoError(t, err)
		require.Zero(t, number)
	})

	t.Run("storing a state value into a wrong type logs an error", func(t *testing.T) {
		stateName := uuid.NewString()

		var m stateManager
		ctx := makeTestContext()
		m.Set(ctx, stateName, 42)

		var number string
		m.Get(ctx, stateName, &number)
		require.Zero(t, number)
	})

	t.Run("storing a state value from local storage into a wrong type logs an error", func(t *testing.T) {
		stateName := uuid.NewString()

		var m stateManager
		ctx := makeTestContext()
		m.Set(ctx, stateName, 42).Persist()
		delete(m.states, stateName)
		require.Empty(t, m.states)

		var number string
		m.Get(ctx, stateName, &number)
		require.Zero(t, number)
	})

	t.Run("getting a non existing state let receiver with current value", func(t *testing.T) {
		stateName := uuid.NewString()

		var m stateManager
		ctx := makeTestContext()

		number := 42
		m.Get(ctx, stateName, &number)
		require.Equal(t, 42, number)
	})
}

func TestStateManagerSet(t *testing.T) {
	t.Run("state is set", func(t *testing.T) {
		stateName := uuid.NewString()

		var m stateManager
		ctx := makeTestContext()

		state := m.Set(ctx, stateName, 42)
		require.Equal(t, 42, state.value)
		require.NotNil(t, state.ctx)
		require.Equal(t, stateName, state.name)
		require.NotNil(t, state.expire)
		require.NotNil(t, state.persist)
		require.NotNil(t, state.broadcast)

		require.NotEmpty(t, m.states)
		state = m.states[stateName]
		require.Equal(t, 42, state.value)
		require.Zero(t, state.ctx)
		require.Empty(t, state.name)
		require.Nil(t, state.expire)
		require.Nil(t, state.persist)
		require.Nil(t, state.broadcast)
	})

	t.Run("state is set and notified to observers", func(t *testing.T) {
		stateName := uuid.NewString()

		e := newTestEngine()
		ctx := e.baseContext()

		var nm nodeManager
		compo, err := nm.Mount(ctx, 1, &hello{})
		require.NoError(t, err)
		ctx = nm.context(ctx, compo)

		var sm stateManager
		var number int
		sm.Observe(ctx, stateName, &number)

		sm.Set(ctx, stateName, 42)
		e.ConsumeAll()
		require.Equal(t, 42, number)
	})

	t.Run("set state removes a non observing observer", func(t *testing.T) {
		stateName := uuid.NewString()

		e := newTestEngine()
		ctx := e.baseContext()

		var nm nodeManager
		compo, err := nm.Mount(ctx, 1, &hello{})
		require.NoError(t, err)
		ctx = nm.context(ctx, compo)

		var sm stateManager
		var number int
		sm.Observe(ctx, stateName, &number).While(func() bool {
			return false
		})

		sm.Set(ctx, stateName, 42)
		e.ConsumeAll()
		require.Zero(t, number)
		require.Empty(t, sm.observers[stateName])
	})

	t.Run("set state log an error when the value cannot be stored in observer receiver", func(t *testing.T) {
		stateName := uuid.NewString()

		e := newTestEngine()
		ctx := e.baseContext()

		var nm nodeManager
		compo, err := nm.Mount(ctx, 1, &hello{})
		require.NoError(t, err)
		ctx = nm.context(ctx, compo)

		var sm stateManager
		var number string
		sm.Observe(ctx, stateName, &number)

		sm.Set(ctx, stateName, 42)
		e.ConsumeAll()
		require.Zero(t, number)
	})

	t.Run("set state trigger observer change handler", func(t *testing.T) {
		stateName := uuid.NewString()

		e := newTestEngine()
		ctx := e.baseContext()

		var nm nodeManager
		compo, err := nm.Mount(ctx, 1, &hello{})
		require.NoError(t, err)
		ctx = nm.context(ctx, compo)

		var sm stateManager
		var number int
		sm.Observe(ctx, stateName, &number).OnChange(func() {
			number = 21
		})

		sm.Set(ctx, stateName, 42)
		e.ConsumeAll()
		require.Equal(t, 21, number)
	})

	t.Run("set state persists a state in local storage", func(t *testing.T) {
		stateName := uuid.NewString()

		var m stateManager
		ctx := makeTestContext()

		state := m.Set(ctx, stateName, 42).Persist()
		require.Equal(t, 42, state.value)
		require.NotNil(t, state.ctx)
		require.Equal(t, stateName, state.name)
		require.NotNil(t, state.expire)
		require.NotNil(t, state.persist)
		require.NotNil(t, state.broadcast)

		require.NotEmpty(t, m.states)
		state = m.states[stateName]
		require.Equal(t, 42, state.value)
		require.Zero(t, state.ctx)
		require.Empty(t, state.name)
		require.Nil(t, state.expire)
		require.Nil(t, state.persist)
		require.Nil(t, state.broadcast)

		delete(m.states, stateName)
		require.Empty(t, m.states)

		var number int
		m.Get(ctx, stateName, &number)
		require.Equal(t, 42, number)
	})

	t.Run("set non encodable state in local storage logs an error", func(t *testing.T) {
		stateName := uuid.NewString()

		var m stateManager
		ctx := makeTestContext()

		m.Set(ctx, stateName, func() {}).Persist()
		delete(m.states, stateName)
		require.Empty(t, m.states)

		var number int
		m.Get(ctx, stateName, &number)
		require.Zero(t, number)
	})

	t.Run("set state persists an encrypted state in local storage", func(t *testing.T) {
		stateName := uuid.NewString()

		var m stateManager
		ctx := makeTestContext()

		state := m.Set(ctx, stateName, 42).PersistWithEncryption()
		require.Equal(t, 42, state.value)
		require.NotNil(t, state.ctx)
		require.Equal(t, stateName, state.name)
		require.NotNil(t, state.expire)
		require.NotNil(t, state.persist)
		require.NotNil(t, state.broadcast)

		require.NotEmpty(t, m.states)
		state = m.states[stateName]
		require.Equal(t, 42, state.value)
		require.Zero(t, state.ctx)
		require.Empty(t, state.name)
		require.Nil(t, state.expire)
		require.Nil(t, state.persist)
		require.Nil(t, state.broadcast)

		delete(m.states, stateName)
		require.Empty(t, m.states)

		var number int
		m.Get(ctx, stateName, &number)
		require.Equal(t, 42, number)
	})

	t.Run("set non encodable encrypted state in local storage logs an error", func(t *testing.T) {
		stateName := uuid.NewString()

		var m stateManager
		ctx := makeTestContext()

		m.Set(ctx, stateName, func() {}).PersistWithEncryption()
		delete(m.states, stateName)
		require.Empty(t, m.states)

		var number int
		m.Get(ctx, stateName, &number)
		require.Zero(t, number)
	})

	t.Run("set state set an expiration duration", func(t *testing.T) {
		stateName := uuid.NewString()

		var m stateManager
		ctx := makeTestContext()

		state := m.Set(ctx, stateName, 42).ExpiresIn(time.Minute)
		require.Equal(t, 42, state.value)
		require.NotZero(t, state.expiresAt)
		require.NotNil(t, state.ctx)
		require.Equal(t, stateName, state.name)
		require.NotNil(t, state.expire)
		require.NotNil(t, state.persist)
		require.NotNil(t, state.broadcast)

		require.NotEmpty(t, m.states)
		state = m.states[stateName]
		require.Equal(t, 42, state.value)
		require.NotZero(t, state.expiresAt)
		require.Zero(t, state.ctx)
		require.Empty(t, state.name)
		require.Nil(t, state.expire)
		require.Nil(t, state.persist)
		require.Nil(t, state.broadcast)

		var number int
		m.Get(ctx, stateName, &number)
		require.Equal(t, 42, number)
	})

	t.Run("set state set an expiration time", func(t *testing.T) {
		stateName := uuid.NewString()

		var m stateManager
		ctx := makeTestContext()

		state := m.Set(ctx, stateName, 42).ExpiresAt(time.Now().Add(time.Minute))
		require.Equal(t, 42, state.value)
		require.NotZero(t, state.expiresAt)
		require.NotNil(t, state.ctx)
		require.Equal(t, stateName, state.name)
		require.NotNil(t, state.expire)
		require.NotNil(t, state.persist)
		require.NotNil(t, state.broadcast)

		require.NotEmpty(t, m.states)
		state = m.states[stateName]
		require.Equal(t, 42, state.value)
		require.NotZero(t, state.expiresAt)
		require.Zero(t, state.ctx)
		require.Empty(t, state.name)
		require.Nil(t, state.expire)
		require.Nil(t, state.persist)
		require.Nil(t, state.broadcast)

		var number int
		m.Get(ctx, stateName, &number)
		require.Equal(t, 42, number)
	})

	t.Run("expired state is not notified to observers", func(t *testing.T) {
		stateName := uuid.NewString()

		e := newTestEngine()
		ctx := e.baseContext()

		var nm nodeManager
		compo, err := nm.Mount(ctx, 1, &hello{})
		require.NoError(t, err)
		ctx = nm.context(ctx, compo)

		var sm stateManager
		var number int
		sm.Observe(ctx, stateName, &number)

		sm.Set(ctx, stateName, 42).ExpiresIn(-time.Hour)
		e.ConsumeAll()
		require.Zero(t, number)
	})

	t.Run("set state broadcasts not supported", func(t *testing.T) {
		if IsClient {
			t.Skip()
		}

		stateName := uuid.NewString()

		var m stateManager
		ctx := makeTestContext()

		state := m.Set(ctx, stateName, 42).Broadcast()
		require.Equal(t, 42, state.value)
		require.NotNil(t, state.ctx)
		require.Equal(t, stateName, state.name)
		require.NotNil(t, state.expire)
		require.NotNil(t, state.persist)
		require.NotNil(t, state.broadcast)

		require.NotEmpty(t, m.states)
		state = m.states[stateName]
		require.Equal(t, 42, state.value)
		require.Zero(t, state.ctx)
		require.Empty(t, state.name)
		require.Nil(t, state.expire)
		require.Nil(t, state.persist)
		require.Nil(t, state.broadcast)
	})

	t.Run("set state broadcast a value", func(t *testing.T) {
		stateName := uuid.NewString()

		if IsServer {
			t.Skip()
		}

		m1 := newTestEngine()
		m2 := newTestEngine()
		m3 := newTestEngine()

		compo3 := &hello{}
		err := m3.Load(compo3)
		require.NoError(t, err)
		ctx3 := m2.nodes.context(m3.baseContext(), compo3)
		var value3 int
		broadcasted3 := false
		ctx3.ObserveState(stateName, &value3).
			OnChange(func() {
				broadcasted3 = true
			})

		compo2 := &hello{}
		err = m2.Load(compo2)
		require.NoError(t, err)
		ctx2 := m2.nodes.context(m2.baseContext(), compo2)
		var value2 int
		broadcasted2 := false
		ctx2.ObserveState(stateName, &value2).
			WithBroadcast().
			OnChange(func() {
				broadcasted2 = true
			})

		compo1 := &hello{}
		err = m1.Load(compo1)
		require.NoError(t, err)
		ctx1 := m1.nodes.context(m1.baseContext(), compo1)
		ctx1.SetState(stateName, 42).Broadcast()
		m1.ConsumeAll()

		for !broadcasted2 {
			m2.ConsumeAll()
			time.Sleep(time.Millisecond * 5)
		}

		require.Equal(t, 42, value2)
		require.Zero(t, value3)
		require.False(t, broadcasted3)
	})
}

func TestStateManagerDelete(t *testing.T) {
	t.Run("state is deleted from memory", func(t *testing.T) {
		stateName := uuid.NewString()

		var m stateManager
		ctx := makeTestContext()

		m.Set(ctx, stateName, 42)
		require.NotEmpty(t, m.states)
		m.Delete(ctx, stateName)
		require.Empty(t, m.states)
	})

	t.Run("state is deleted from local storage", func(t *testing.T) {
		stateName := uuid.NewString()

		var m stateManager
		ctx := makeTestContext()

		m.Set(ctx, stateName, 42).Persist()
		delete(m.states, stateName)
		require.Empty(t, m.states)

		var number int
		m.Get(ctx, stateName, &number)
		require.Equal(t, 42, number)

		var number2 int
		m.Delete(ctx, stateName)
		m.Get(ctx, stateName, &number2)
		require.Zero(t, number2)
	})
}

func TestStateManagerCleanup(t *testing.T) {
	t.Run("non observing observers are removed", func(t *testing.T) {
		stateName := uuid.NewString()

		ctx := makeTestContext()
		var nm nodeManager
		compo, err := nm.Mount(ctx, 1, &hello{})
		require.NoError(t, err)
		ctx = nm.context(ctx, compo)

		var sm stateManager
		var number int
		sm.Observe(ctx, stateName, &number)
		require.NotEmpty(t, sm.observers)
		nm.Dismount(compo)

		sm.Cleanup()
		require.Empty(t, sm.observers)
	})
}

func TestStateManagerCleanupExpiredPersistedStates(t *testing.T) {
	if IsServer {
		t.Skip()
	}

	t.Run("expired states are removed", func(t *testing.T) {
		stateName := uuid.NewString()
		ctx := makeTestContext()

		var sm stateManager
		sm.Set(ctx, stateName, 42).
			ExpiresIn(-time.Second).
			Persist()

		var state storableState
		ctx.LocalStorage().Get(stateName, &state)
		require.NotZero(t, state)

		sm.CleanupExpiredPersistedStates(ctx)

		var expiredState storableState
		ctx.LocalStorage().Get(stateName, &expiredState)
		require.Zero(t, expiredState)
	})

	t.Run("non expired states are not removed", func(t *testing.T) {
		stateName := uuid.NewString()
		ctx := makeTestContext()

		var sm stateManager
		sm.Set(ctx, stateName, 42).
			ExpiresIn(time.Minute).
			Persist()

		var state storableState
		ctx.LocalStorage().Get(stateName, &state)
		require.NotZero(t, state)

		sm.CleanupExpiredPersistedStates(ctx)

		var nonExpiredState storableState
		ctx.LocalStorage().Get(stateName, &nonExpiredState)
		require.NotZero(t, nonExpiredState)
	})

	t.Run("non state are not removed", func(t *testing.T) {
		stateName := uuid.NewString()
		ctx := makeTestContext()

		ctx.LocalStorage().Set(stateName, 42)

		var sm stateManager
		sm.CleanupExpiredPersistedStates(ctx)

		var nonStateValue int
		ctx.LocalStorage().Get(stateName, &nonStateValue)
		require.Equal(t, 42, nonStateValue)
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
