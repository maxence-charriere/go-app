package app

import "reflect"

// EventHandler represents a function that can handle HTML events. They are
// always called on the UI goroutine.
type EventHandler func(ctx Context, e Event)

// Event is the interface that describes a javascript event.
type Event struct {
	Value
}

// PreventDefault cancels the event if it is cancelable. The default action that
// belongs to the event will not occur.
func (e Event) PreventDefault() {
	e.Call("preventDefault")
}

// If several listeners are attached to the same element for the same event type,
// they are called in the order in which they were added. If
// StopImmediatePropagation() is invoked during one such call, no remaining
// listeners will be called, either on that element or any other element.
func (e Event) StopImmediatePropagation() {
	e.Call("stopImmediatePropagation")
}

type eventHandlers map[string]eventHandler

func (h eventHandlers) Set(event string, eh EventHandler, scope ...any) {
	if eh != nil {
		h[event] = makeEventHandler(event, eh, scope...)
	}
}

type eventHandler struct {
	event     string
	scope     string
	goHandler EventHandler
	jsHandler Func
	close     func()
}

func makeEventHandler(event string, h EventHandler, scope ...any) eventHandler {
	return eventHandler{
		event:     event,
		scope:     toPath(scope...),
		goHandler: h,
	}
}

func (h eventHandler) Equal(v eventHandler) bool {
	return h.event == v.event &&
		h.scope == v.scope &&
		reflect.ValueOf(h.goHandler).Pointer() == reflect.ValueOf(v.goHandler).Pointer()
}

func trackMousePosition(e Event) {
	x := e.Get("clientX")
	if !x.Truthy() {
		return
	}

	y := e.Get("clientY")
	if !y.Truthy() {
		return
	}

	Window().setCursorPosition(x.Int(), y.Int())
}
