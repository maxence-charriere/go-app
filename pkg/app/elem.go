package app

import (
	"fmt"
	"io"
	"net/url"

	"github.com/maxence-charriere/go-app/v6/pkg/log"
)

type standardNode interface {
	UI
	nodeWithChildren

	attributes() map[string]string
	setAttribute(k string, v interface{})
	setAttributeValue(k, v string)
	removeAttributeValue(k string)
	eventHandlers() map[string]eventHandler
	setEventHandler(k string, h EventHandler)
	setEventHandlerValue(k string, h eventHandler)
	removeEventHandlerValue(k string, h eventHandler)
	mount() error
	children() []UI
	appendChild(child UI)
	appendChildValue(child UI)
	removeChild(child UI)
	removeChildValue(child UI)
	replaceChildValue(old, new UI)
	update(n standardNode) (updated bool)
	triggerOnNav(u *url.URL)
}

type elem struct {
	parentNode  UI
	value       Value
	body        []UI
	tag         string
	attrs       map[string]string
	events      map[string]eventHandler
	selfClosing bool
}

func (e *elem) JSValue() Value {
	return e.value
}

func (e *elem) parent() UI {
	return e.parentNode
}

func (e *elem) setParent(p UI) {
	e.parentNode = p
}

func (e *elem) dismount() {
	for _, c := range e.body {
		c.dismount()
	}

	for k, h := range e.events {
		e.removeEventHandlerValue(k, h)
	}

	e.value = nil
}

func (e *elem) replaceChild(old, new UI) {
	if e.selfClosing {
		log.Error("replacing child failed").
			T("error", "self closing tag can't have children").
			T("tag", e.tag).
			Panic()
	}

	for i, c := range e.body {
		if c == old {
			e.body[i] = new
			return
		}
	}
}

func (e *elem) setBody(parent UI, body []Node) {
	if e.selfClosing {
		log.Error("set body failed").
			T("error", "self closing tag can't have children").
			T("tag", e.tag).
			Panic()
	}

	ibody := indirect(body...)

	for _, n := range ibody {
		n.setParent(parent)
	}

	e.body = ibody
}

func (e *elem) attributes() map[string]string {
	return e.attrs
}

func (e *elem) setAttribute(k string, v interface{}) {
	if e.attrs == nil {
		e.attrs = make(map[string]string)
	}

	switch t := v.(type) {
	case string:
		switch k {
		case "style":
			style := e.attrs[k]
			if style != "" {
				style += ";"
			}
			style += t
			e.attrs[k] = style

		default:
			e.attrs[k] = t
		}

	case bool:
		if t {
			e.attrs[k] = ""
		} else {
			delete(e.attrs, k)
		}

	default:
		e.attrs[k] = fmt.Sprintf("%v", t)
	}
}

func (e *elem) setAttributeValue(k, v string) {
	e.value.Call("setAttribute", k, v)
}

func (e *elem) removeAttributeValue(k string) {
	e.value.Call("removeAttribute", k)
}

func (e *elem) eventHandlers() map[string]eventHandler {
	return e.events
}

func (e *elem) setEventHandler(k string, h EventHandler) {
	if e.events == nil {
		e.events = make(map[string]eventHandler)
	}

	e.events[k] = eventHandler{function: h}
}

func (e *elem) setEventHandlerValue(k string, h eventHandler) {
	callback := FuncOf(func(this Value, args []Value) interface{} {
		dispatcher(func() {
			event := Event{Value: args[0]}
			trackMousePosition(event)
			h.function(this, event)
		})
		return nil
	})

	h.jsFunction = callback
	e.events[k] = h
	e.value.Call("addEventListener", k, callback)
}

func (e *elem) removeEventHandlerValue(k string, h eventHandler) {
	if h.jsFunction != nil {
		e.value.Call("removeEventListener", k, h.jsFunction)
		h.jsFunction.Release()
		h.jsFunction = nil
		e.events[k] = h
	}
}

func (e *elem) mount() error {
	if e.value != nil {
		return fmt.Errorf("node already mounted: %+v", e)
	}

	value := Window().
		Get("document").
		Call("createElement", e.tag)
	if !value.Truthy() {
		return fmt.Errorf("creating node failed: %+v", e)
	}
	e.value = value

	for k, v := range e.attrs {
		e.setAttributeValue(k, v)
	}

	for k, v := range e.events {
		e.setEventHandlerValue(k, v)
	}

	for _, c := range e.body {
		if err := mount(c); err != nil {
			return err
		}
		e.appendChildValue(c)
	}

	return nil
}

func (e *elem) children() []UI {
	return e.body
}

func (e *elem) appendChild(child UI) {
	e.body = append(e.body, child)
}

func (e *elem) appendChildValue(child UI) {
	e.value.Call("appendChild", child)
}

func (e *elem) removeChild(child UI) {
	children := e.body
	for i, c := range children {
		if c == child {
			copy(children[i:], children[i+1:])
			children[len(children)-1] = nil
			children = children[:len(children)-1]
			e.body = children
			return
		}
	}
}

func (e *elem) removeChildValue(child UI) {
	e.value.Call("removeChild", child)
}

func (e *elem) replaceChildValue(old, new UI) {
	e.value.Call("replaceChild", new, old)
}

func (e *elem) update(n standardNode) (updated bool) {
	for k := range e.attrs {
		if _, ok := n.attributes()[k]; !ok {
			e.removeAttributeValue(k)
			delete(e.attrs, k)
			updated = true
		}
	}

	for k, v := range n.attributes() {
		if v != e.attrs[k] {
			e.setAttribute(k, v)
			e.setAttributeValue(k, v)
			updated = true
		}
	}

	for k, v := range e.events {
		if _, ok := n.eventHandlers()[k]; !ok {
			e.removeEventHandlerValue(k, v)
			delete(e.events, k)
			updated = true
		}
	}

	for k, v := range n.eventHandlers() {
		if current := e.events[k]; !v.equals(current) {
			e.removeEventHandlerValue(k, current)
			e.setEventHandler(k, v.function)
			e.setEventHandlerValue(k, v)
			updated = true
		}
	}

	return updated
}

func (e *elem) triggerOnNav(u *url.URL) {
	for _, c := range e.body {
		triggerOnNav(c, u)
	}
}

func (e *elem) html(w io.Writer) {
	e.htmlWithIndent(w, 0)
}

func (e *elem) htmlWithIndent(w io.Writer, indent int) {
	writeIndent(w, indent)
	w.Write(stob("<"))
	w.Write(stob(e.tag))

	for k, v := range e.attrs {
		w.Write(stob(" "))
		w.Write(stob(k))

		if v != "" {
			w.Write(stob(`="`))
			w.Write(stob(v))
			w.Write(stob(`"`))
		}
	}

	w.Write(stob(">"))

	if e.selfClosing {
		return
	}

	for _, c := range e.body {
		w.Write(ln())
		c.(writableNode).htmlWithIndent(w, indent+1)
	}

	if len(e.body) != 0 {
		w.Write(ln())
		writeIndent(w, indent)
	}

	w.Write(stob("</"))
	w.Write(stob(e.tag))
	w.Write(stob(">"))
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

	window.setCursorPosition(x.Int(), y.Int())
}

type eventHandler struct {
	function   EventHandler
	jsFunction Func
}

func (h eventHandler) equals(o eventHandler) bool {
	return fmt.Sprintf("%p", h.function) == fmt.Sprintf("%p", o.function)
}
