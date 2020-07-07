package app

import (
	"context"
	"io"
	"net/url"

	"github.com/maxence-charriere/go-app/v7/pkg/errors"
)

type elem struct {
	attrs       map[string]string
	body        []UI
	ctx         context.Context
	ctxCancel   func()
	events      map[string]eventHandler
	jsvalue     Value
	parentElem  UI
	selfClosing bool
	tag         string
	this        UI
}

func (e *elem) Kind() Kind {
	return HTML
}

func (e *elem) JSValue() Value {
	return e.jsvalue
}

func (e *elem) Mounted() bool {
	return e.ctx != nil && e.ctx.Err() == nil &&
		e.self() != nil &&
		e.jsvalue != nil
}

func (e *elem) name() string {
	return e.tag
}

func (e *elem) self() UI {
	return e.this
}

func (e *elem) setSelf(n UI) {
	e.this = n
}

func (e *elem) context() context.Context {
	return e.ctx
}

func (e *elem) attributes() map[string]string {
	return e.attrs
}

func (e *elem) eventHandlers() map[string]eventHandler {
	return e.events
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

func (e *elem) mount() error {
	if e.Mounted() {
		return errors.New("mounting ui element failed").
			Tag("reason", "already mounted").
			Tag("name", e.name()).
			Tag("kind", e.Kind())
	}

	e.ctx, e.ctxCancel = context.WithCancel(context.Background())

	v := Window().Get("document").Call("createElement", e.tag)
	if !v.Truthy() {
		return errors.New("mounting ui element failed").
			Tag("reason", "create javascript node returned nil").
			Tag("name", e.name()).
			Tag("kind", e.Kind())
	}
	e.jsvalue = v

	for k, v := range e.attrs {
		e.setJsAttr(k, v)
	}

	for k, v := range e.events {
		e.setJsEventHandler(k, v)
	}

	for _, c := range e.children() {
		if err := e.appendChild(c, true); err != nil {
			return errors.New("mounting ui element failed").
				Tag("name", e.name()).
				Tag("kind", e.Kind()).
				Wrap(err)
		}
	}

	return nil
}

func (e *elem) dismount() {
	for _, c := range e.children() {
		dismount(c)
	}

	for k, v := range e.events {
		e.delJsEventHandler(k, v)
	}

	e.ctxCancel()
	e.jsvalue = nil
}

func (e *elem) update(n UI) error {
	if !e.Mounted() {
		return nil
	}

	if n.Kind() != e.Kind() || n.name() != e.name() {
		return errors.New("updating ui element failed").
			Tag("replace", true).
			Tag("reason", "different element types").
			Tag("current-kind", e.Kind()).
			Tag("current-name", e.name()).
			Tag("updated-kind", n.Kind()).
			Tag("updated-name", n.name())
	}

	e.updateAttrs(n.attributes())
	e.updateEventHandler(n.eventHandlers())

	achildren := e.children()
	bchildren := n.children()
	i := 0

	// Update children:
	for len(achildren) != 0 && len(bchildren) != 0 {
		a := achildren[0]
		b := bchildren[0]

		err := update(a, b)
		if isErrReplace(err) {
			err = e.replaceChildAt(i, b)
		}

		if err != nil {
			return errors.New("updating ui element failed").
				Tag("kind", e.Kind()).
				Tag("name", e.name()).
				Wrap(err)
		}

		achildren = achildren[1:]
		bchildren = bchildren[1:]
		i++
	}

	// Remove children:
	for len(achildren) != 0 {
		if err := e.removeChildAt(i); err != nil {
			return errors.New("updating ui element failed").
				Tag("kind", e.Kind()).
				Tag("name", e.name()).
				Wrap(err)
		}

		achildren = achildren[1:]
	}

	// Add children:
	for len(bchildren) != 0 {
		c := bchildren[0]

		if err := e.appendChild(c, false); err != nil {
			return errors.New("updating ui element failed").
				Tag("kind", e.Kind()).
				Tag("name", e.name()).
				Wrap(err)
		}

		bchildren = bchildren[1:]
	}

	return nil
}

func (e *elem) appendChild(c UI, onlyJsValue bool) error {
	if err := mount(c); err != nil {
		return errors.New("appending child failed").
			Tag("name", e.name()).
			Tag("kind", e.Kind()).
			Tag("child-name", c.name()).
			Tag("child-kind", c.Kind()).
			Wrap(err)
	}

	if !onlyJsValue {
		e.body = append(e.body, c)
	}

	c.setParent(e.self())
	e.JSValue().Call("appendChild", c)
	return nil
}

func (e *elem) replaceChildAt(idx int, new UI) error {
	old := e.body[idx]

	if err := mount(new); err != nil {
		return errors.New("replacing child failed").
			Tag("name", e.name()).
			Tag("kind", e.Kind()).
			Tag("index", idx).
			Tag("old-name", old.name()).
			Tag("old-kind", old.Kind()).
			Tag("new-name", new.name()).
			Tag("new-kind", new.Kind()).
			Wrap(err)
	}

	e.body[idx] = new
	new.setParent(e.self())
	e.JSValue().Call("replaceChild", new, old)

	dismount(old)
	return nil
}

func (e *elem) removeChildAt(idx int) error {
	body := e.body
	if idx < 0 || idx >= len(body) {
		return errors.New("removing child failed").
			Tag("reason", "index out of range").
			Tag("index", idx).
			Tag("name", e.name()).
			Tag("kind", e.Kind())
	}

	c := body[idx]

	copy(body[idx:], body[idx+1:])
	body[len(body)-1] = nil
	body = body[:len(body)-1]
	e.body = body

	e.JSValue().Call("removeChild", c)
	dismount(c)
	return nil
}

func (e *elem) updateAttrs(attrs map[string]string) {
	for k := range e.attrs {
		if _, exist := attrs[k]; !exist {
			e.delAttr(k)
		}
	}

	if e.attrs == nil && len(attrs) != 0 {
		e.attrs = make(map[string]string, len(attrs))
	}

	for k, v := range attrs {
		if curval := e.attrs[k]; curval != v {
			e.attrs[k] = v
			e.setJsAttr(k, v)
		}
	}
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

func (e *elem) delAttr(k string) {
	e.JSValue().Call("removeAttribute", k)
	delete(e.attrs, k)
}

func (e *elem) updateEventHandler(handlers map[string]eventHandler) {
	for k, current := range e.events {
		if _, exists := handlers[k]; !exists {
			e.delJsEventHandler(k, current)
		}
	}

	if e.events == nil && len(handlers) != 0 {
		e.events = make(map[string]eventHandler, len(handlers))
	}

	for k, new := range handlers {
		if current, exists := e.events[k]; !current.equal(new) {
			if exists {
				e.delJsEventHandler(k, current)
			}

			e.events[k] = new
			e.setJsEventHandler(k, new)
		}
	}
}

func (e *elem) setEventHandler(k string, h EventHandler) {
	if e.events == nil {
		e.events = make(map[string]eventHandler)
	}

	e.events[k] = eventHandler{
		event: k,
		value: h,
	}
}

func (e *elem) setJsEventHandler(k string, h eventHandler) {
	jshandler := makeJsEventHandler(e.self(), h.value)
	h.jsvalue = jshandler
	e.events[k] = h
	e.JSValue().Call("addEventListener", k, jshandler)
}

func (e *elem) delJsEventHandler(k string, h eventHandler) {
	e.JSValue().Call("removeEventListener", k, h.jsvalue)
	h.jsvalue.Release()
	delete(e.events, k)
}

func (e *elem) setBody(body ...UI) {
	if e.selfClosing {
		panic(errors.New("setting html element body failed").
			Tag("reason", "self closing element can't have children").
			Tag("name", e.name()),
		)
	}

	e.body = FilterUIElems(body...)
}

func (e *elem) onNav(u *url.URL) {
	for _, c := range e.children() {
		c.onNav(u)
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
