package app

import (
	"reflect"
	"sync"

	"github.com/google/uuid"
	"github.com/murlokswarm/app/key"
	"github.com/pkg/errors"
)

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

func newEventRegistry(dispatcher func(func())) *eventRegistry {
	return &eventRegistry{
		handlers:   make(map[string][]eventHandler),
		dispatcher: dispatcher,
	}
}

type eventRegistry struct {
	mutex      sync.RWMutex
	handlers   map[string][]eventHandler
	dispatcher func(f func())
}

func (m *eventRegistry) Subscribe(name string, handler interface{}) (unsuscribe func()) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if reflect.ValueOf(handler).Kind() != reflect.Func {
		Panic(errors.Errorf("can't subscribe event %s: handler is not a func: %T",
			name,
			handler,
		))
	}

	id := uuid.New().String()

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
	m.mutex.Lock()
	defer m.mutex.Unlock()

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
	m.mutex.RLock()
	defer m.mutex.RUnlock()

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

// EventSubscriber represents an event subscriber.
type EventSubscriber struct {
	registry    *eventRegistry
	unsuscribes []func()
}

// Subscribe subscribes a function to the named event.
// It panics if f is not a func.
func (s *EventSubscriber) Subscribe(name string, f interface{}) *EventSubscriber {
	unsubscribe := s.registry.Subscribe(name, f)
	s.unsuscribes = append(s.unsuscribes, unsubscribe)
	return s
}

// Close closes the event handler and unsubscribe all its events.
func (s *EventSubscriber) Close() {
	for _, unsuscribe := range s.unsuscribes {
		unsuscribe()
	}
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
	Source    EventSource
}

// WheelEvent represents an onwheel event arg.
type WheelEvent struct {
	DeltaX    float64
	DeltaY    float64
	DeltaZ    float64
	DeltaMode DeltaMode
	Source    EventSource
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
	Source    EventSource
}

// DragAndDropEvent represents an ondrop event arg.
type DragAndDropEvent struct {
	Files         []string
	Data          string
	DropEffect    string
	EffectAllowed string
	Source        EventSource
}

// EventSource represents a descriptor to an event source.
type EventSource struct {
	GoappID string
	CompoID string
	ID      string
	Class   string
	Data    map[string]string
	Value   string
}
