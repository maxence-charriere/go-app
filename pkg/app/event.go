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

// EventOption represents an option for configuring event handlers, such as
// setting scopes or marking them as passive.
type EventOption struct {
	name  string
	value string
}

// EventScope returns an EventOption that adds a scope to an event handler.
// This is useful in dynamic UI contexts to ensure correct mounting and
// dismounting of handlers as UI elements are reordered. The scope is defined
// by concatenating the provided path arguments into a unique identifier.
func EventScope(v ...any) EventOption {
	return EventOption{
		name:  "scope",
		value: toPath(v...),
	}
}

// PassiveEvent returns an EventOption that marks an event handler as passive.
// Passive handlers improve performance for high-frequency events by signaling
// that they will not call Event.PreventDefault(), allowing the browser to
// optimize event processing.
// More on passive listeners: https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/addEventListener#using_passive_listeners
func PassiveEvent() EventOption {
	return EventOption{
		name: "passive",
	}
}

type eventHandlers map[string]eventHandler

func (h eventHandlers) Set(event string, eh EventHandler, options ...EventOption) {
	if eh != nil {
		h[event] = makeEventHandler(event, eh, options...)
	}
}

type eventHandler struct {
	event     string
	scope     string
	passive   bool
	goHandler EventHandler
	jsHandler Func
	close     func()
}

func makeEventHandler(event string, h EventHandler, options ...EventOption) eventHandler {
	handler := eventHandler{
		event:     event,
		goHandler: h,
	}

	for _, option := range options {
		switch option.name {
		case "scope":
			handler.scope = option.value

		case "passive":
			handler.passive = true
		}
	}
	return handler
}

func (h eventHandler) Equal(v eventHandler) bool {
	return h.event == v.event &&
		h.scope == v.scope &&
		h.passive == v.passive &&
		reflect.ValueOf(h.goHandler).Pointer() == reflect.ValueOf(v.goHandler).Pointer()
}

func (h eventHandler) options() map[string]any {
	if h.passive {
		return map[string]any{
			"passive": true,
		}
	}
	return nil
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
