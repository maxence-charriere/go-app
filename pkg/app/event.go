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

func (h eventHandler) Mount(src UI) eventHandler {
	jsHandler := makeJSEventHandler(src, h.goHandler)
	src.JSValue().addEventListener(h.event, jsHandler)

	close := func() {
		src.JSValue().removeEventListener(h.event, jsHandler)
		jsHandler.Release()
	}

	h.jsHandler = jsHandler
	h.close = close
	return h
}

func (h eventHandler) Dismount() {
	if h.close != nil {
		h.close()
	}
}

func makeJSEventHandler(src UI, h EventHandler) Func {
	return FuncOf(func(this Value, args []Value) interface{} {
		src.getDispatcher().Dispatch(Dispatch{
			Mode:   Update,
			Source: src,
			Function: func(ctx Context) {
				ctx.Emit(func() {
					event := Event{
						Value: args[0],
					}
					trackMousePosition(event)
					h(ctx, event)
				})
			},
		})
		return nil
	})
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
