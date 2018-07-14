package app

import (
	"reflect"
	"sync"

	"github.com/google/uuid"
	"github.com/murlokswarm/app/key"
	"github.com/pkg/errors"
)

var (
	// DefaultEventRegistry is the default event registry.
	DefaultEventRegistry EventRegistry
)

// EventRegistry is the interface that describes an event registry.
type EventRegistry interface {
	// Subscribe subscribes the given handler to the named event.
	// It panics if handler is not a func.
	Subscribe(name string, handler interface{}) (unsuscribe func())

	EventDispatcher
}

// EventDispatcher is the interface that describes an event dispatcher.
type EventDispatcher interface {
	// Dispatch dispatches the named event with the given argument.
	// It is done on the UI goroutine.
	Dispatch(name string, arg interface{})
}

type eventHandler struct {
	ID      string
	Handler interface{}
}

// NewEventRegistry creates an event registry.
func NewEventRegistry(dispatcher func(func())) EventRegistry {
	return &eventRegistry{
		handlers:   make(map[string][]eventHandler),
		dispatcher: dispatcher,
	}
}

type eventRegistry struct {
	handlers   map[string][]eventHandler
	dispatcher func(f func())
}

func (m *eventRegistry) Subscribe(name string, handler interface{}) (unsuscribe func()) {
	if reflect.ValueOf(handler).Kind() != reflect.Func {
		panic(errors.Errorf("can't subscribe event %s: handler is not a func: %T",
			name,
			handler,
		))
	}

	id := uuid.New()

	handlers := m.handlers[name]
	handlers = append(handlers, eventHandler{
		ID:      id,
		Handler: handler,
	})
	m.handlers[name] = handlers

	return func() {
		m.Unsubscribe(name, id)
	}
}

func (m *eventRegistry) Unsubscribe(name string, id string) {
	handlers := m.handlers[name]

	for i, h := range handlers {
		if h.ID == id {
			end := len(handlers) - 1
			handlers[i] = handlers[end]
			handlers[end] = eventHandler{}
			handlers = handlers[:end]

			m.handlers[name] = handlers
			return
		}
	}
}

func (m *eventRegistry) Dispatch(name string, arg interface{}) {
	for _, h := range m.handlers[name] {
		val := reflect.ValueOf(h.Handler)
		typ := val.Type()

		if typ.NumIn() == 0 {
			m.dispatcher(func() {
				val.Call(nil)
			})
			return
		}

		argVal := reflect.ValueOf(arg)
		argTyp := typ.In(0)

		if !argVal.Type().ConvertibleTo(argTyp) {
			Log("dispatching event %s failed: %s",
				name,
				errors.Errorf("can't convert %s to %s", argVal.Type(), argTyp),
			)
			return
		}

		m.dispatcher(func() {
			val.Call([]reflect.Value{
				argVal.Convert(argTyp),
			})
		})
	}
}

// ConcurrentEventRegistry returns a decorated version of the given event
// registry that ensure concurrency safety.
func ConcurrentEventRegistry(r EventRegistry) EventRegistry {
	return &concurrentEventRegistry{
		base: r,
	}
}

type concurrentEventRegistry struct {
	mutex sync.RWMutex
	base  EventRegistry
}

func (r *concurrentEventRegistry) Subscribe(name string, handler interface{}) func() {
	r.mutex.Lock()
	unsubscribe := r.base.Subscribe(name, handler)
	r.mutex.Unlock()

	return func() {
		r.mutex.Lock()
		unsubscribe()
		r.mutex.Unlock()
	}
}

func (r *concurrentEventRegistry) Dispatch(name string, arg interface{}) {
	r.mutex.RLock()
	r.base.Dispatch(name, arg)
	r.mutex.RUnlock()
}

// EventSubscriber is the interface that describes an event subscriber.
type EventSubscriber interface {
	// Subscribe subscribes the given handler to the named event.
	// It panics if handler is not a func.
	Subscribe(name string, handler interface{})

	// Close closes the event handler and unsubscribe all its events.
	Close() error
}

// NewEventSubscriber creates an event subscriber.
func NewEventSubscriber() EventSubscriber {
	return &eventSubscriber{
		registry: DefaultEventRegistry,
	}
}

type eventSubscriber struct {
	registry    EventRegistry
	unsuscribes []func()
}

func (s *eventSubscriber) Subscribe(name string, handler interface{}) {
	unsubscribe := s.registry.Subscribe(name, handler)
	s.unsuscribes = append(s.unsuscribes, unsubscribe)
}

func (s *eventSubscriber) Close() error {
	for _, unsuscribe := range s.unsuscribes {
		unsuscribe()
	}
	return nil
}

// MouseEvent represents an onmouse event arg.
type MouseEvent struct {
	ClientX   float64
	ClientY   float64
	PageX     float64
	PageY     float64
	ScreenX   float64
	ScreenY   float64
	Button    int
	Detail    int
	AltKey    bool
	CtrlKey   bool
	MetaKey   bool
	ShiftKey  bool
	InnerText string
}

// WheelEvent represents an onwheel event arg.
type WheelEvent struct {
	DeltaX    float64
	DeltaY    float64
	DeltaZ    float64
	DeltaMode DeltaMode
}

// DeltaMode is an indication of the units of measurement for a delta value.
type DeltaMode uint64

// KeyboardEvent represents an onkey event arg.
type KeyboardEvent struct {
	CharCode  rune
	KeyCode   key.Code
	Location  key.Location
	AltKey    bool
	CtrlKey   bool
	MetaKey   bool
	ShiftKey  bool
	InnerText string
}

// DragAndDropEvent represents an ondrop event arg.
type DragAndDropEvent struct {
	Files         []string
	Data          string
	DropEffect    string
	EffectAllowed string
}
