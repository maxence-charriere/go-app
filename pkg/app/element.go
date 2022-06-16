package app

import (
	"context"
	"io"
	"net/url"
	"reflect"
	"strconv"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

type elemWithChildren interface {
	UI

	replaceChildAt(idx int, new UI) error
}

type elem struct {
	tag           string
	xmlns         string
	isSelfClosing bool
	attributes    attributes
	eventHandlers eventHandlers
	parent        UI
	children      []UI

	context       context.Context
	contextCancel func()
	dispatcher    Dispatcher
	jsElement     Value
	this          UI
}

func (e *elem) Kind() Kind {
	return HTML
}

func (e *elem) JSValue() Value {
	return e.jsElement
}

func (e *elem) IsMounted() bool {
	return e.context != nil && e.context.Err() == nil
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

func (e *elem) getContext() context.Context {
	return e.context
}

func (e *elem) getDispatcher() Dispatcher {
	return e.dispatcher
}

func (e *elem) getAttributes() attributes {
	return e.attributes
}

func (e *elem) getEventHandlers() eventHandlers {
	return e.eventHandlers
}

func (e *elem) getParent() UI {
	return e.parent
}

func (e *elem) setParent(p UI) {
	e.parent = p
}

func (e *elem) getChildren() []UI {
	return e.children
}

func (e *elem) mount(d Dispatcher) error {
	if e.IsMounted() {
		return errors.New("html element is already mounted").
			Tag("tag", e.tag).
			Tag("kind", e.Kind())
	}

	e.context, e.contextCancel = context.WithCancel(context.Background())
	e.dispatcher = d

	jsElement, err := Window().createElement(e.tag, e.xmlns)
	if err != nil {
		return errors.New("mounting js element failed").
			Tag("kind", e.Kind()).
			Tag("tag", e.tag).
			Tag("xmlns", e.xmlns).
			Wrap(err)
	}
	e.jsElement = jsElement

	e.attributes.Mount(jsElement, d.resolveStaticResource)
	e.eventHandlers.Mount(e)

	for i, c := range e.children {
		if err := mount(d, c); err != nil {
			return errors.New("mounting child failed").
				Tag("index", i).
				Tag("child", c.name()).
				Tag("child-kind", c.Kind()).
				Wrap(err)
		}

		c.setParent(e.self())
		e.JSValue().appendChild(c)
	}

	return nil
}

func (e *elem) dismount() {
	for _, c := range e.children {
		dismount(c)
	}

	for _, eh := range e.eventHandlers {
		eh.Dismount()
	}

	e.contextCancel()
}

func (e *elem) canUpdateWith(n UI) bool {
	return n.Kind() == e.Kind() && n.name() == e.name()
}

func (e *elem) updateWith(n UI) error {
	if !e.IsMounted() {
		return errors.New("cannot update a non mounted html element").
			Tag("element", reflect.TypeOf(e.self()))
	}

	if !e.canUpdateWith(n) {
		return errors.New("cannot update html element").
			Tag("current-element", reflect.TypeOf(e.self())).
			Tag("new-element", reflect.TypeOf(n))
	}

	if e.attributes == nil && n.getAttributes() != nil {
		e.attributes = n.getAttributes()
		e.attributes.Mount(e.jsElement, e.dispatcher.resolveStaticResource)
	} else if e.attributes != nil {
		e.attributes.Update(
			e.jsElement,
			n.getAttributes(),
			e.getDispatcher().resolveStaticResource,
		)
	}

	if e.eventHandlers == nil && n.getEventHandlers() != nil {
		e.eventHandlers = n.getEventHandlers()
		e.eventHandlers.Mount(e)
	} else if e.eventHandlers != nil {
		e.eventHandlers.Update(e, n.getEventHandlers())
	}

	childrenA := e.children
	childrenB := n.getChildren()
	i := 0

	for len(childrenA) != 0 && len(childrenB) != 0 {
		a := childrenA[0]
		b := childrenB[0]

		if canUpdate(a, b) {
			if err := update(a, b); err != nil {
				return errors.New("updating child failed").
					Tag("child", reflect.TypeOf(a)).
					Tag("new-child", reflect.TypeOf(b)).
					Tag("index", i).
					Wrap(err)
			}
		} else {
			if err := e.replaceChildAt(i, b); err != nil {
				return errors.New("replacing child failed").
					Tag("child", reflect.TypeOf(a)).
					Tag("new-child", reflect.TypeOf(b)).
					Tag("index", i).
					Wrap(err)
			}
		}

		childrenA = childrenA[1:]
		childrenB = childrenB[1:]
		i++
	}

	for len(childrenA) != 0 {
		if err := e.removeChildAt(i); err != nil {
			return errors.New("removing child failed").
				Tag("child", reflect.TypeOf(childrenA[0])).
				Tag("index", i).
				Wrap(err)
		}

		childrenA = childrenA[1:]
	}

	for len(childrenB) != 0 {
		b := childrenB[0]

		if err := e.appendChild(b); err != nil {
			return errors.New("appending child failed").
				Tag("child", reflect.TypeOf(b)).
				Tag("index", i).
				Wrap(err)
		}

		childrenB = childrenB[1:]
	}

	return nil
}

func (e *elem) replaceChildAt(idx int, new UI) error {
	old := e.children[idx]

	if err := mount(e.getDispatcher(), new); err != nil {
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

	e.children[idx] = new
	new.setParent(e.self())
	e.JSValue().replaceChild(new, old)

	dismount(old)
	return nil
}

func (e *elem) removeChildAt(i int) error {
	if i < 0 || i >= len(e.children) {
		return errors.New("index out of range").
			Tag("index", i).
			Tag("children-count", len(e.children))
	}

	child := e.children[i]
	e.jsElement.removeChild(child)
	dismount(child)

	children := e.children
	copy(children[i:], children[i+1:])
	children[len(children)-1] = nil
	e.children = children[:len(children)-1]
	return nil
}

func (e *elem) appendChild(v UI) error {
	if err := mount(e.getDispatcher(), v); err != nil {
		return errors.New("mounting element failed").
			Tag("element", reflect.TypeOf(v)).
			Wrap(err)
	}

	v.setParent(e.self())
	e.JSValue().appendChild(v)
	e.children = append(e.children, v)
	return nil
}

func (e *elem) setAttr(name string, value interface{}) {
	if e.attributes == nil {
		e.attributes = make(attributes)
	}
	e.attributes.Set(name, value)
}

func (e *elem) resolveURLAttr(k, v string) string {
	if !isURLAttrValue(k) {
		return v
	}
	return e.dispatcher.resolveStaticResource(v)
}

func (e *elem) setJsAttr(k, v string) {
	switch k {
	case "value":
		e.JSValue().Set("value", v)

	case "class":
		e.JSValue().Set("className", v)

	case "contenteditable":
		e.JSValue().Set("contentEditable", v)

	case "async",
		"autofocus",
		"autoplay",
		"checked",
		"default",
		"defer",
		"disabled",
		"hidden",
		"ismap",
		"loop",
		"multiple",
		"muted",
		"open",
		"readonly",
		"required",
		"reversed",
		"selected":
		switch k {
		case "ismap":
			k = "isMap"
		case "readonly":
			k = "readOnly"
		}
		v, _ := strconv.ParseBool(v)
		e.JSValue().Set(k, v)

	default:
		if isURLAttrValue(k) {
			v = e.getDispatcher().resolveStaticResource(v)
		}
		e.JSValue().setAttr(k, v)
	}
}

func (e *elem) setEventHandler(event string, h EventHandler, scope ...interface{}) {
	if e.eventHandlers == nil {
		e.eventHandlers = make(eventHandlers)
	}
	e.eventHandlers.Set(event, h, scope...)
}

func (e *elem) setBody(body ...UI) {
	if e.isSelfClosing {
		panic(errors.New("setting html element body failed").
			Tag("reason", "self closing element can't have children").
			Tag("name", e.name()),
		)
	}

	e.children = FilterUIElems(body...)
}

func (e *elem) onNav(u *url.URL) {
	for _, c := range e.getChildren() {
		c.onNav(u)
	}
}

func (e *elem) onAppUpdate() {
	for _, c := range e.getChildren() {
		c.onAppUpdate()
	}
}

func (e *elem) onAppInstallChange() {
	for _, c := range e.getChildren() {
		c.onAppInstallChange()
	}
}

func (e *elem) onResize() {
	for _, c := range e.getChildren() {
		c.onResize()
	}
}

func (e *elem) preRender(p Page) {
	for _, c := range e.getChildren() {
		c.preRender(p)
	}
}

func (e *elem) html(w io.Writer) {
	w.Write([]byte("<"))
	w.Write([]byte(e.tag))

	for k, v := range e.attributes {
		w.Write([]byte(" "))
		w.Write([]byte(k))

		if v != "" {
			w.Write([]byte(`="`))
			w.Write([]byte(resolveAttributeURLValue(k, v, func(s string) string {
				if e.dispatcher != nil {
					return e.dispatcher.resolveStaticResource(v)
				}
				return v
			})))
			w.Write([]byte(`"`))
		}
	}

	w.Write([]byte(">"))

	if e.isSelfClosing {
		return
	}

	for _, c := range e.children {
		w.Write(ln())
		if c.self() == nil {
			c.setSelf(c)
		}
		c.html(w)
	}

	if len(e.children) != 0 {
		w.Write(ln())
	}

	w.Write([]byte("</"))
	w.Write([]byte(e.tag))
	w.Write([]byte(">"))
}

func (e *elem) htmlWithIndent(w io.Writer, indent int) {
	writeIndent(w, indent)
	w.Write([]byte("<"))
	w.Write([]byte(e.tag))

	for k, v := range e.attributes {
		w.Write([]byte(" "))
		w.Write([]byte(k))

		if v != "" {
			w.Write([]byte(`="`))
			w.Write([]byte(resolveAttributeURLValue(k, v, func(s string) string {
				if e.dispatcher != nil {
					return e.dispatcher.resolveStaticResource(v)
				}
				return v
			})))
			w.Write([]byte(`"`))
		}
	}

	w.Write([]byte(">"))

	if e.isSelfClosing {
		return
	}

	for _, c := range e.children {
		w.Write(ln())
		if c.self() == nil {
			c.setSelf(c)
		}
		c.htmlWithIndent(w, indent+1)
	}

	if len(e.children) != 0 {
		w.Write(ln())
		writeIndent(w, indent)
	}

	w.Write([]byte("</"))
	w.Write([]byte(e.tag))
	w.Write([]byte(">"))
}

func isURLAttrValue(k string) bool {
	switch k {
	case "cite",
		"data",
		"href",
		"src",
		"srcset":
		return true

	default:
		return false
	}
}
