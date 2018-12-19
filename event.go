package app

import (
	"reflect"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Event is a string that identifies an app event.
type Event string

// Subscriber is the interface that describes an event subscriber.
type Subscriber interface {
	// Subscribe subscribes a function to the given key.
	// It panics if f is not a func.
	Subscribe(e Event, f interface{}) Subscriber

	// Close unsubscribes all the subscriptions.
	Close()
}

type subscriber struct {
	registry    *eventRegistry
	unsuscribes []func()
}

func (s *subscriber) Subscribe(e Event, f interface{}) Subscriber {
	unsubscribe := s.registry.subscribe(e, f)
	s.unsuscribes = append(s.unsuscribes, unsubscribe)
	return s
}

func (s *subscriber) Close() {
	for _, unsuscribe := range s.unsuscribes {
		unsuscribe()
	}
}

type eventHandler struct {
	ID      string
	Handler interface{}
}

type eventRegistry struct {
	mutex    sync.RWMutex
	handlers map[Event][]eventHandler
	callOnUI func(f func())
}

func newEventRegistry(callOnUI func(func())) *eventRegistry {
	return &eventRegistry{
		handlers: make(map[Event][]eventHandler),
		callOnUI: callOnUI,
	}
}

func (r *eventRegistry) subscribe(e Event, handler interface{}) func() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if reflect.ValueOf(handler).Kind() != reflect.Func {
		Panic(errors.Errorf("can't subscribe to event %s: handler is not a func: %T",
			e,
			handler,
		))
	}

	id := uuid.New().String()
	handlers := r.handlers[e]

	handlers = append(handlers, eventHandler{
		ID:      id,
		Handler: handler,
	})

	r.handlers[e] = handlers

	return func() {
		r.unsubscribe(e, id)
	}
}

func (r *eventRegistry) unsubscribe(e Event, id string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	handlers := r.handlers[e]

	for i, h := range handlers {
		if h.ID == id {
			end := len(handlers) - 1
			handlers[i] = handlers[end]
			handlers[end] = eventHandler{}
			handlers = handlers[:end]

			r.handlers[e] = handlers
			return
		}
	}
}

func (r *eventRegistry) Emit(e Event, v interface{}) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, h := range r.handlers[e] {
		val := reflect.ValueOf(h.Handler)
		typ := val.Type()

		if typ.NumIn() == 0 {
			r.callOnUI(func() {
				val.Call(nil)
			})

			return
		}

		argVal := reflect.ValueOf(v)
		argTyp := typ.In(0)

		if !argVal.Type().ConvertibleTo(argTyp) {
			Log("dispatching event %s failed: %s",
				e,
				errors.Errorf("can't convert %s to %s", argVal.Type(), argTyp),
			)
			return
		}

		r.callOnUI(func() {
			val.Call([]reflect.Value{
				argVal.Convert(argTyp),
			})
		})
	}
}
