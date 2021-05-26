package app

import (
	"reflect"
	"testing"

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
		Value(&v)

	require.True(t, isSubscribeCalled)
	require.Equal(t, elem, o.element)
	require.Len(t, o.conditions, 1)
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
}

func TestStoreGetSet(t *testing.T) {
	d := NewClientTester(Div())
	defer d.Close()

	s := newStore(d)
	key := "/test/getSet"

	var v int
	err := s.Get(key, &v)
	require.NoError(t, err)
	require.Zero(t, v)

	s.Set(key, 42)
	err = s.Get(key, &v)
	require.NoError(t, err)
	require.Equal(t, 42, v)
}

func TestStoreObserve(t *testing.T) {
	source := &foo{}
	d := NewClientTester(source)
	defer d.Close()

	s := newStore(d)
	key := "/test/observe"

	s.Observe(key, source).Value(&source.Bar)
	require.Equal(t, "", source.Bar)
	require.Len(t, s.states, 1)
	require.Len(t, s.states[key].observers, 1)

	s.Set(key, "hello")
	d.Consume()
	require.Equal(t, "hello", source.Bar)

	s.Set(key, nil)
	d.Consume()
	require.Equal(t, "", source.Bar)

	s.Set(key, 42)
	d.Consume()
	require.Equal(t, "", source.Bar)

	d.Mount(Div())
	d.Consume()
	s.Set(key, "hi")
	require.Empty(t, s.states[key].observers)

	s.Set(key, 42)
	s.Observe(key, source).Value(&source.Bar)
	require.Equal(t, "", source.Bar)
}

func TestStoreValue(t *testing.T) {
	nb := 42
	c := copyTester{pointer: &nb}

	utests := []struct {
		scenario string
		src      interface{}
		recv     interface{}
		expected interface{}
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
