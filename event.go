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

// Subscribe subscribes a function to the given event. Emit fails if the
// subscribbed func have more arguments than the emitted event.
//
// Panics if f is not a func.
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

// Emit emits the event with the given arguments.
func (r *EventRegistry) Emit(e Event, args ...interface{}) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, h := range r.handlers[e] {
		if err := r.callHandler(h.Handler, args...); err != nil {
			Logf("emitting %s failed: %s", e, err)
		}
	}
}

func (r *EventRegistry) callHandler(h interface{}, args ...interface{}) error {
	v := reflect.ValueOf(h)
	t := v.Type()

	argsv := make([]reflect.Value, t.NumIn())

	for i := 0; i < t.NumIn(); i++ {
		argt := t.In(i)

		if i >= len(args) {
			return errors.Errorf("missing %v at index %v", argt, i)
		}

		argv := reflect.ValueOf(args[i])
		if !argv.Type().ConvertibleTo(argt) {
			return errors.Errorf("arg at index %v is not a %v: %v", i, argt, argv.Type())
		}

		argsv[i] = argv.Convert(argt)
	}

	r.ui <- func() {
		v.Call(argsv)
	}

	return nil
}
