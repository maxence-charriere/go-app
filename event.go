package app

import (
	"reflect"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Event is a string that identifies an app event.
type Event string

// Subscriber is a struct to subscribe to events emitted by a event registry.
type Subscriber struct {
	// The event restry that emits events. By default tt uses the app default
	// event registry.
	// This should be set only for testing purpose.
	Events *EventRegistry

	unsuscribes []func()
}

// Subscribe subscribes a function to the given key.
// It panics if f is not a func.
func (s *Subscriber) Subscribe(e Event, f interface{}) *Subscriber {
	unsubscribe := s.Events.subscribe(e, f)
	s.unsuscribes = append(s.unsuscribes, unsubscribe)
	return s
}

// Close unsubscribes all the subscriptions.
func (s *Subscriber) Close() {
	for _, unsuscribe := range s.unsuscribes {
		unsuscribe()
	}
}

type eventHandler struct {
	ID      string
	Handler interface{}
}

// EventRegistry is a struct that manages event flow.
type EventRegistry struct {
	mutex    sync.RWMutex
	handlers map[Event][]eventHandler
	ui       chan func()
}

// NewEventRegistry creates a event registry.
func NewEventRegistry(ui chan func()) *EventRegistry {
	return &EventRegistry{
		handlers: make(map[Event][]eventHandler),
		ui:       ui,
	}
}

func (r *EventRegistry) subscribe(e Event, handler interface{}) func() {
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

func (r *EventRegistry) unsubscribe(e Event, id string) {
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

// Emit emits the event with the given value.
func (r *EventRegistry) Emit(e Event, v interface{}) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, h := range r.handlers[e] {
		val := reflect.ValueOf(h.Handler)
		typ := val.Type()

		if typ.NumIn() == 0 {
			r.ui <- func() { val.Call(nil) }
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

		r.ui <- func() {
			val.Call([]reflect.Value{
				argVal.Convert(argTyp),
			})
		}
	}
}
