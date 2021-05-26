package app

import (
	"reflect"
	"sync"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

// Observer is an observer that observes changes for a given state.
type Observer interface {
	// Defines a condition that reports whether the observer keeps observing the
	// associated state. Multiple conditions can be defined by successively
	// calling While().
	While(condition func() bool) Observer

	// Stores the value associated with the observed state into the given
	// receiver. Panics when the receiver is not a pointer or nil.
	//
	// The receiver is updated each time the associated state changes. It is
	// unchanged when its pointed value has a different type than the associated
	// state value.
	Value(recv interface{})
}

type observer struct {
	element    UI
	subscribe  func(*observer)
	conditions []func() bool
	receiver   interface{}
}

func newObserver(elem UI, subscribe func(*observer)) *observer {
	return &observer{
		element:   elem,
		subscribe: subscribe,
	}
}

func (o *observer) While(fn func() bool) Observer {
	o.conditions = append(o.conditions, fn)
	return o
}

func (o *observer) Value(recv interface{}) {
	if reflect.ValueOf(recv).Kind() != reflect.Ptr {
		panic(errors.New("observer value receiver is not a pointer"))
	}

	o.receiver = recv
	o.subscribe(o)
}

func (o *observer) isObserving() bool {
	if !o.element.Mounted() {
		return false
	}

	for _, c := range o.conditions {
		if !c() {
			return false
		}
	}

	return true
}

type state struct {
	value     interface{}
	observers map[*observer]struct{}
}

type store struct {
	mutex  sync.RWMutex
	states map[string]state
	disp   Dispatcher
}

func makeStore(d Dispatcher) store {
	return store{
		states: make(map[string]state),
		disp:   d,
	}
}

func (s *store) Set(key string, v interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	state := s.get(key)
	state.value = v
	s.states[key] = state

	for o := range state.observers {
		if !o.isObserving() {
			delete(state.observers, o)
			continue
		}

		elem := o.element
		recv := o.receiver
		s.disp.Dispatch(elem, func(ctx Context) {
			if err := storeValue(recv, v); err != nil {
				Log(errors.New("notifying observer failed").
					Tag("state", key).
					Tag("element", reflect.TypeOf(elem)).
					Wrap(err))
			}
		})
	}
}

func (s *store) Get(key string, recv interface{}) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	state, exists := s.states[key]
	if !exists {
		return
	}

	if err := storeValue(recv, state.value); err != nil {
		Log(errors.New("storing state value failed").
			Tag("state", key).
			Wrap(err))
	}
}

func (s *store) Observe(key string, elem UI) Observer {
	return newObserver(elem, func(o *observer) {
		if err := s.subscribe(key, o); err != nil {
			Log(errors.New("notifying observer failed").
				Tag("state", key).
				Tag("element", reflect.TypeOf(elem)).
				Wrap(err))
		}
	})
}

func (s *store) subscribe(key string, o *observer) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	state := s.get(key)
	state.observers[o] = struct{}{}
	s.states[key] = state
	return storeValue(o.receiver, state.value)
}

func (s *store) get(key string) state {
	st, exists := s.states[key]
	if !exists {
		st = state{observers: make(map[*observer]struct{})}
	}
	return st
}

func storeValue(recv, v interface{}) error {
	dst := reflect.ValueOf(recv)
	if dst.Kind() != reflect.Ptr {
		return errors.New("receiver is not a pointer")
	}
	dst = dst.Elem()

	src := reflect.ValueOf(v)
	switch {
	case src == reflect.Value{}:
		dst.Set(reflect.Zero(dst.Type()))
		return nil

	case src.Kind() == reflect.Ptr:
		src = src.Elem()
	}

	if src.Type() != dst.Type() {
		return errors.New("value and receiver are not of the same type").
			Tag("value", src.Type()).
			Tag("receiver", dst.Type())
	}

	dst.Set(src)
	return nil
}
