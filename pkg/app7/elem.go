package app

import (
	"context"
	"fmt"

	"github.com/maxence-charriere/go-app/v6/pkg/errors"
)

type elem struct {
	attrs         map[string]string
	body          []UI
	ctx           context.Context
	ctxCancel     func()
	eventHandlers map[string]elemEventHandler
	jsvalue       Value
	parentElem    UI
	self          UI
	selfClosing   bool
	tag           string
}

func (e *elem) Kind() Kind {
	return HTML
}

func (e *elem) JSValue() Value {
	return e.jsvalue
}

func (e *elem) Mounted() bool {
	return e.self != nil && e.ctx != nil && e.jsvalue != nil
}

func (e *elem) setSelf(n UI) {
	e.self = n
}

func (e *elem) context() context.Context {
	return e.ctx
}

func (e *elem) parent() UI {
	return e.parentElem
}

func (e *elem) setParent(p UI) {
	e.parentElem = p
}

func (e *elem) children() []UI {
	return e.body
}

func (e *elem) appendChild(c UI) {
	e.body = append(e.body, c)
	e.JSValue().Call("appendChild", c)
}

func (e *elem) removeChild(c UI) {
	body := e.body
	for i, n := range body {
		if n == c {
			copy(body[i:], body[i+1:])
			body[len(body)-1] = nil
			body = body[:len(body)-1]
			e.body = body

			e.JSValue().Call("removeChild", c)
			return
		}
	}
}

func (e *elem) mount() error {
	if e.Mounted() {
		return errors.New("mounting ui element failed").
			Tag("reason", "already mounted").
			Tag("kind", e.Kind()).
			Tag("tag", e.tag)
	}

	e.ctx, e.ctxCancel = context.WithCancel(context.Background())

	v := Window().Get("document").Call("createElement", e.tag)
	if !v.Truthy() {
		return errors.New("mounting ui element failed").
			Tag("reason", "creating javascript node return nil").
			Tag("kind", e.Kind()).
			Tag("tag", e.tag)
	}
	e.jsvalue = v

	for k, v := range e.attrs {
		e.setJsAttr(k, v)
	}

	for k, v := range e.eventHandlers {
		e.setJsEventHandler(k, v)
	}

	for _, c := range e.children() {
		if err := c.mount(); err != nil {
			return err
		}

		e.appendChild(c)
	}

	return nil
}

func (e *elem) dismount() {
	for _, c := range e.children() {
		c.dismount()
	}

	for k, v := range e.eventHandlers {
		e.delJsEventHandler(k, v)
	}

	e.ctxCancel()
	e.jsvalue = nil
}

func (e *elem) update(n UI) error {
	panic("not implemented")
}

func (e *elem) setAttr(k string, v interface{}) {
	if e.attrs == nil {
		e.attrs = make(map[string]string)
	}

	switch k {
	case "style":
		s := e.attrs[k] + toString(v) + ";"
		e.attrs[k] = s
		return

	case "class":
		s := e.attrs[k]
		if s != "" {
			s += " "
		}
		s += toString(v)
		e.attrs[k] = s
		return
	}

	switch v := v.(type) {
	case bool:
		if !v {
			delete(e.attrs, k)
			return
		}
		e.attrs[k] = ""

	default:
		e.attrs[k] = toString(v)
	}
}

func (e *elem) setJsAttr(k, v string) {
	e.JSValue().Call("setAttribute", k, v)
}

func (e *elem) setEventHandler(k string, h EventHandler) {
	if e.eventHandlers == nil {
		e.eventHandlers = make(map[string]elemEventHandler)
	}

	e.eventHandlers[k] = elemEventHandler{
		event: k,
		value: h,
	}
}

func (e *elem) setJsEventHandler(k string, h elemEventHandler) {
	jshandler := makeJsEventHandler(e.self, h.value)
	h.jsvalue = jshandler
	e.eventHandlers[k] = h
	e.JSValue().Call("addEventListener", k, jshandler)
}

func (e *elem) delJsEventHandler(k string, h elemEventHandler) {
	e.JSValue().Call("addEventListener", k, h.jsvalue)
	h.jsvalue.Release()
	delete(e.eventHandlers, k)
}

func (e *elem) setBody(self UI, body ...UI) {
	if e.selfClosing {
		panic(errors.New("setting html element body failed").
			Tag("reason", "self closing element can't have children").
			Tag("tag", e.tag),
		)
	}

	body = FilterUIElems(body...)
	for _, n := range body {
		n.setParent(self)
	}
	e.body = body
}

type elemEventHandler struct {
	event   string
	jsvalue Func
	value   EventHandler
}

func (h elemEventHandler) equal(o elemEventHandler) bool {
	return h.event == o.event &&
		fmt.Sprintf("%p", h.value) == fmt.Sprintf("%p", o.value)
}

func makeJsEventHandler(src UI, h EventHandler) Func {
	return FuncOf(func(this Value, args []Value) interface{} {
		dispatch(func() {
			if !src.Mounted() {
				return
			}

			ctx := Context{
				Context: src.context(),
				Src:     src,
				JSSrc:   src.JSValue(),
			}

			event := Event{
				Value: args[0],
			}

			h(ctx, event)
		})

		return nil
	})
}
