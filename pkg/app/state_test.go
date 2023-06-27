package app

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestObserver(t *testing.T) {
	elem := &hello{}
	c := NewClientTester(elem)
	defer c.Close()

	isSubscribeCalled := false
	isObserving := true
	v := 42

	o := newObserver(elem, func(o *observer) {
		isSubscribeCalled = true
	})
	o.While(func() bool { return isObserving }).
		OnChange(func() {}).
		Value(&v)

	require.True(t, isSubscribeCalled)
	require.Equal(t, elem, o.element)
	require.Len(t, o.conditions, 1)
	require.Len(t, o.onChanges, 1)
	require.NotNil(t, o.receiver)
	require.True(t, o.isObserving())

	c.Mount(Div())
	c.Consume()
	require.False(t, o.isObserving())

	c.Mount(elem)
	c.Consume()
	require.True(t, o.isObserving())

	isObserving = false
	require.False(t, o.isObserving())

	require.Panics(t, func() {
		var s string
		newObserver(elem, func(*observer) {}).Value(s)
	})
}

func TestStateIsExpired(t *testing.T) {
	utests := []struct {
		scenario  string
		state     State
		isExpired bool
	}{
		{
			scenario:  "state without expiration",
			state:     State{},
			isExpired: false,
		},
		{
			scenario:  "state is not expired",
			state:     State{ExpiresAt: time.Now().Add(time.Minute)},
			isExpired: false,
		},
		{
			scenario:  "state is expired",
			state:     State{ExpiresAt: time.Now().Add(-time.Minute)},
			isExpired: true,
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			require.Equal(t, u.isExpired, u.state.isExpired(time.Now()))
		})
	}
}

func TestStore(t *testing.T) {
	d := NewClientTester(Div())
	defer d.Close()

	s := newStore(d)
	defer s.Close()
	defer s.Cleanup()
	key := "/test/store"

	var v int
	s.Get(key, &v)
	require.Zero(t, v)

	s.Set(key, 42)
	s.Get(key, &v)
	require.Equal(t, 42, v)

	s.Set(key, "21")
	s.Get(key, &v)
	require.Equal(t, 42, v)

	s.Del(key)
	require.Empty(t, s.states)
}

func TestStorePersist(t *testing.T) {
	d := NewClientTester(Div())
	defer d.Close()

	s := newStore(d)
	defer s.Close()
	key := "/test/store/persist"

	t.Run("value is pesisted", func(t *testing.T) {
		var v int

		s.Set(key, 42, Persist)
		s.Get(key, &v)
		require.Equal(t, 42, v)
		require.Equal(t, 1, d.getLocalStorage().Len())
	})

	t.Run("value is not pesisted", func(t *testing.T) {
		var v int

		s.Set(key, struct {
			Func func()
		}{}, Persist)
		s.Get(key, &v)
		require.Equal(t, 0, v)
	})

	t.Run("value is obtained from local storage", func(t *testing.T) {
		var v int

		s.Set(key, 21, Persist)
		delete(s.states, key)
		require.Empty(t, s.states)

		s.Get(key, &v)
		require.Equal(t, 21, v)
		require.Equal(t, 1, d.getLocalStorage().Len())
	})

	t.Run("value is observed from local storage", func(t *testing.T) {
		var v int

		s.Set(key, 84, Persist)
		delete(s.states, key)
		require.Empty(t, s.states)

		s.Observe(key, Div()).Value(&v)
		require.Equal(t, 84, v)
		require.Equal(t, 1, d.getLocalStorage().Len())
	})

	t.Run("value is deleted", func(t *testing.T) {
		var v int

		s.Set(key, 1977, Persist)
		s.Del(key)

		require.Empty(t, s.states)
		s.Get(key, &v)
		require.Equal(t, 0, v)
		require.Equal(t, 0, d.getLocalStorage().Len())
	})
}

func TestStoreEncrypt(t *testing.T) {
	d := NewClientTester(Div())
	defer d.Close()

	s := newStore(d)
	defer s.Close()
	key := "/test/store/crypt"

	t.Run("value is encrypted and decrypted", func(t *testing.T) {
		var v int

		s.Set(key, 42, Persist, Encrypt)
		s.Get(key, &v)
		require.Equal(t, 42, v)
		require.Equal(t, 2, d.getLocalStorage().Len(), d.getLocalStorage()) // Contain app ID.
	})

	t.Run("value is decrypted from local storage", func(t *testing.T) {
		var v int

		s.Set(key, 43, Persist)
		delete(s.states, key)

		s.Get(key, &v)
		require.Equal(t, 43, v)
		require.Equal(t, 2, d.getLocalStorage().Len()) // Contain app ID.
	})
}

func TestStoreExpiresIn(t *testing.T) {
	d := NewClientTester(Div())
	defer d.Close()

	s := newStore(d)
	defer s.Close()
	key := "/test/store/expiresIn"

	t.Run("value is not expired", func(t *testing.T) {
		var v int

		s.Set(key, 42, Persist, ExpiresIn(time.Minute))
		s.Get(key, &v)
		require.Equal(t, 42, v)
		require.Equal(t, 1, d.getLocalStorage().Len())
	})

	t.Run("get expired value", func(t *testing.T) {
		var v int

		s.Set(key, 21, Persist, ExpiresIn(-time.Minute))
		s.Get(key, &v)
		require.Equal(t, 0, v)
		require.Equal(t, 0, d.getLocalStorage().Len())
	})

	t.Run("get persisted expired value", func(t *testing.T) {
		var v int

		s.Del(key)
		delete(s.states, key)

		s.disp.getLocalStorage().Set(key, persistentState{
			ExpiresAt: time.Now().Add(-time.Minute),
		})

		s.Get(key, &v)
		require.Equal(t, 0, v)
		require.Equal(t, 0, d.getLocalStorage().Len())
	})

	t.Run("observe expired value", func(t *testing.T) {
		var v int

		s.Set(key, 21, Persist, ExpiresIn(-time.Minute))
		s.Observe(key, Div()).Value(&v)
		require.Equal(t, 0, v)
		require.Equal(t, 0, d.getLocalStorage().Len())
	})

	t.Run("expire expired values", func(t *testing.T) {
		s.Set(key, 99, Persist, ExpiresIn(time.Minute))
		require.Len(t, s.states, 1)
		require.Equal(t, 1, d.getLocalStorage().Len())

		state := s.states[key]
		state.ExpiresAt = time.Now().Add(-time.Minute)
		s.states[key] = state
		require.True(t, state.isExpired(time.Now()))
		require.Equal(t, 1, d.getLocalStorage().Len())

		s.expireExpiredValues()
		require.True(t, state.isExpired(time.Now()))
		require.Equal(t, 0, d.getLocalStorage().Len())
	})
}

func TestStoreBroadcast(t *testing.T) {
	d1 := NewClientTester(&foo{})
	s1 := newStore(d1)
	defer d1.Close()
	defer s1.Close()

	bar := &bar{}
	d2 := NewClientTester(bar)
	s2 := newStore(d2)
	defer d2.Close()
	defer s2.Close()

	require.NotEqual(t, s1.id, s2.id)

	key := "/test/store/broadcast"
	t.Run("state is not broadcasted", func(t *testing.T) {
		var v int
		s2.Observe(key, bar).Value(&v)

		s1.Set(key, func() {}, Broadcast)
		d2.Consume()
		require.Zero(t, v)
	})

	t.Run("state is broadcasted", func(t *testing.T) {
		if IsServer {
			t.Skip()
		}

		var v int

		s2.Observe(key, bar).Value(&v)
		s1.Set(key, 42, Broadcast)

		time.Sleep(time.Millisecond * 100)
		d2.Consume()
		require.Equal(t, 42, v)
	})

	t.Run("broadcasted state is not observed", func(t *testing.T) {
		if IsServer {
			t.Skip()
		}

		var v func()

		s2.Observe(key, bar).Value(&v)
		s1.Set(key, 42, Broadcast)

		time.Sleep(time.Millisecond * 50)
		d2.Consume()
		require.Zero(t, v)
	})
}

func TestStoreObserve(t *testing.T) {
	key := "/test/observe"

	t.Run("state is observed and stored in value", func(t *testing.T) {
		foo := &foo{}
		d := NewClientTester(foo)
		defer d.Close()

		s := newStore(d)
		defer s.Close()

		s.Observe(key, foo).Value(&foo.Bar)
		require.Equal(t, "", foo.Bar)

		s.Set(key, "hello")
		d.Consume()
		require.Equal(t, "hello", foo.Bar)
	})

	t.Run("function is called when observed value changes", func(t *testing.T) {
		foo := &foo{}
		d := NewClientTester(foo)
		defer d.Close()

		s := newStore(d)
		defer s.Close()

		isOnChangeCalled := false
		s.Observe(key, foo).
			OnChange(func() {
				isOnChangeCalled = true
			}).
			Value(&foo.Bar)

		s.Set(key, "hello")
		d.Consume()
		require.Equal(t, "hello", foo.Bar)
		require.True(t, isOnChangeCalled)
	})

	t.Run("zero value is set and stored in observed value", func(t *testing.T) {
		foo := &foo{}
		d := NewClientTester(foo)
		defer d.Close()

		s := newStore(d)
		defer s.Close()

		s.Observe(key, foo).Value(&foo.Bar)
		s.Set(key, "hi")
		d.Consume()
		require.Equal(t, "hi", foo.Bar)

		s.Set(key, nil)
		d.Consume()
		require.Equal(t, "", foo.Bar)
	})

	t.Run("observed value with different type is not stored", func(t *testing.T) {
		foo := &foo{}
		d := NewClientTester(foo)
		defer d.Close()

		s := newStore(d)
		defer s.Close()

		s.Observe(key, foo).Value(&foo.Bar)

		s.Set(key, 42)
		d.Consume()
		require.Equal(t, "", foo.Bar)
	})

	t.Run("observer that stop observing is removed from store", func(t *testing.T) {
		foo := &foo{}
		d := NewClientTester(foo)
		defer d.Close()

		s := newStore(d)
		defer s.Close()

		isObserving := true
		s.Observe(key, foo).
			While(func() bool {
				return isObserving
			}).
			Value(&foo.Bar)

		s.Set(key, "hi")
		d.Consume()
		require.Equal(t, "hi", foo.Bar)

		isObserving = false
		s.Set(key, "hey")
		d.Consume()
		require.Equal(t, "hi", foo.Bar)
		require.Empty(t, s.states[key].observers)
	})

	t.Run("observer created from unmounted component is removed", func(t *testing.T) {
		foo := &foo{}
		d := NewClientTester(foo)
		defer d.Close()

		s := newStore(d)
		defer s.Close()

		s.Observe(key, foo).Value(&foo.Bar)

		d.Mount(Div())
		d.Consume()
		require.False(t, foo.Mounted())

		s.Set(key, "hi")
		d.Consume()
		require.Empty(t, s.states[key].observers)
		require.Equal(t, "", foo.Bar)
	})

	t.Run("current value is stored in observer value", func(t *testing.T) {
		foo := &foo{}
		d := NewClientTester(foo)
		defer d.Close()

		s := newStore(d)
		defer s.Close()

		s.Set(key, "bye")
		s.Observe(key, foo).Value(&foo.Bar)
		require.Equal(t, "bye", foo.Bar)
	})

	t.Run("current value fails to be stored in observer value", func(t *testing.T) {
		foo := &foo{}
		d := NewClientTester(foo)
		defer d.Close()

		s := newStore(d)
		defer s.Close()

		s.Set(key, 42)
		s.Observe(key, foo).Value(&foo.Bar)
		require.Equal(t, "", foo.Bar)
	})
}

func TestRemoveUnusedObservers(t *testing.T) {
	source := &foo{}
	d := NewClientTester(source)
	defer d.Close()

	s := newStore(d)
	defer s.Close()
	key := "/test/observe/remove"

	var v int
	n := 5
	for i := 0; i < 5; i++ {
		s.Observe(key, source).
			While(func() bool { return false }).
			Value(&v)
	}
	state := s.states[key]
	require.Len(t, state.observers, n)

	s.removeUnusedObservers()
	require.Empty(t, state.observers)
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

func TestExpireExpriredPersistentValues(t *testing.T) {
	if IsServer {
		t.Skip()
	}

	d := NewClientTester(&foo{})
	defer d.Close()
	localStorage := d.getLocalStorage()

	s := newStore(d)
	defer s.Close()

	t.Run("non expired state is not removed", func(t *testing.T) {
		localStorage.Clear()
		s.setPersistent("/hello", false, time.Now().Add(time.Minute), "hello")
		require.Equal(t, 1, localStorage.Len())

		s.expireExpriredPersistentValues()
		require.Equal(t, 1, localStorage.Len())
	})

	t.Run("expired state is removed", func(t *testing.T) {
		localStorage.Clear()
		s.setPersistent("/bye", false, time.Now().Add(-time.Minute), "bye")
		require.Equal(t, 1, localStorage.Len())

		s.expireExpriredPersistentValues()
		require.Equal(t, 0, localStorage.Len())
	})

	t.Run("non state value is not removed", func(t *testing.T) {
		localStorage.Clear()
		localStorage.Set("/hi", "hi")
		require.Equal(t, 1, localStorage.Len())

		s.expireExpriredPersistentValues()
		require.Equal(t, 1, localStorage.Len())
	})
}
