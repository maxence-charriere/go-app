package app

import (
	"fmt"
	"io"
)

type elem struct {
	parentNode  nodeWithChildren
	value       Value
	body        []ValueNode
	tag         string
	attrs       map[string]string
	events      map[string]eventHandler
	selfClosing bool
}

func (e *elem) JSValue() Value {
	return e.value
}

func (e *elem) parent() nodeWithChildren {
	return e.parentNode
}

func (e *elem) setParent(p nodeWithChildren) {
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

func (e *elem) replaceChild(old, new ValueNode) {
	for i, c := range e.body {
		if c == old {
			e.body[i] = new
			return
		}
	}
}

func (e *elem) setBody(parent nodeWithChildren, body []Node) {
	ibody := indirect(body...)

	for _, n := range ibody {
		n.setParent(parent.(nodeWithChildren))
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
		e.attrs[k] = t

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
		event := Event{Value: args[0]}
		trackMousePosition(event)
		h.function(this, event)
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

func (e *elem) children() []ValueNode {
	return e.body
}

func (e *elem) appendChild(child ValueNode) {
	e.body = append(e.body, child)
}

func (e *elem) appendChildValue(child ValueNode) {
	e.value.Call("appendChild", child)
}

func (e *elem) removeChild(child ValueNode) {
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

func (e *elem) removeChildValue(child ValueNode) {
	e.value.Call("removeChild", child)
}

func (e *elem) replaceChildValue(old, new ValueNode) {
	e.value.Call("replaceChild", new, old)
}

func (e *elem) update(n standardNode) {
	for k := range e.attrs {
		if _, ok := n.attributes()[k]; !ok {
			e.removeAttributeValue(k)
			delete(e.attrs, k)
		}
	}

	for k, v := range n.attributes() {
		if v != e.attrs[k] {
			e.setAttribute(k, v)
			e.setAttributeValue(k, v)
		}
	}

	for k, v := range e.events {
		if _, ok := n.eventHandlers()[k]; !ok {
			e.removeEventHandlerValue(k, v)
			delete(e.events, k)
		}
	}

	for k, v := range n.eventHandlers() {
		if current := e.events[k]; !v.equals(current) {
			e.removeEventHandlerValue(k, current)
			e.setEventHandler(k, v.function)
			e.setEventHandlerValue(k, v)
		}
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
