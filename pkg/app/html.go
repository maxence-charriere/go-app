package app

import (
	"context"
	"io"
	"reflect"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

type htmlElement struct {
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

func (e *htmlElement) Kind() Kind {
	return HTML
}

func (e *htmlElement) JSValue() Value {
	return e.jsElement
}

func (e *htmlElement) Mounted() bool {
	return e.context != nil && e.context.Err() == nil
}

func (e *htmlElement) name() string {
	return e.tag
}

func (e *htmlElement) self() UI {
	return e.this
}

func (e *htmlElement) setSelf(v UI) {
	e.this = v
}

func (e *htmlElement) getContext() context.Context {
	return e.context
}

func (e *htmlElement) getDispatcher() Dispatcher {
	return e.dispatcher
}

func (e *htmlElement) getAttributes() attributes {
	return e.attributes
}

func (e *htmlElement) getEventHandlers() eventHandlers {
	return e.eventHandlers
}

func (e *htmlElement) getParent() UI {
	return e.parent
}

func (e *htmlElement) setParent(p UI) {
	e.parent = p
}

func (e *htmlElement) getChildren() []UI {
	return e.children
}

func (e *htmlElement) mount(d Dispatcher) error {
	if e.Mounted() {
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

func (e *htmlElement) dismount() {
	for _, c := range e.children {
		dismount(c)
	}

	for _, eh := range e.eventHandlers {
		eh.Dismount()
	}

	e.contextCancel()
}

func (e *htmlElement) canUpdateWith(v UI) bool {
	return e.Mounted() &&
		e.Kind() == v.Kind() &&
		e.name() == v.name()
}

func (e *htmlElement) updateWith(v UI) error {
	if !e.canUpdateWith(v) {
		return errors.New("cannot update html element with given element").
			Tag("current", reflect.TypeOf(e.self())).
			Tag("new", reflect.TypeOf(v))
	}

	if e.attributes == nil && v.getAttributes() != nil {
		e.attributes = v.getAttributes()
		e.attributes.Mount(e.jsElement, e.dispatcher.resolveStaticResource)
	} else if e.attributes != nil {
		e.attributes.Update(
			e.jsElement,
			v.getAttributes(),
			e.getDispatcher().resolveStaticResource,
		)
	}

	if e.eventHandlers == nil && v.getEventHandlers() != nil {
		e.eventHandlers = v.getEventHandlers()
		e.eventHandlers.Mount(e)
	} else if e.eventHandlers != nil {
		e.eventHandlers.Update(e, v.getEventHandlers())
	}

	childrenA := e.children
	childrenB := v.getChildren()
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

func (e *htmlElement) replaceChildAt(idx int, new UI) error {
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

func (e *htmlElement) removeChildAt(i int) error {
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

func (e *htmlElement) appendChild(v UI) error {
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

func (e *htmlElement) setAttr(name string, value any) {
	if e.attributes == nil {
		e.attributes = make(attributes)
	}
	e.attributes.Set(name, value)
}

func (e *htmlElement) setEventHandler(event string, h EventHandler, scope ...any) {
	if e.eventHandlers == nil {
		e.eventHandlers = make(eventHandlers)
	}
	e.eventHandlers.Set(event, h, scope...)
}

func (e *htmlElement) setChildren(v ...UI) {
	if e.isSelfClosing {
		panic(errors.New("cannot set children of a self closing element").
			Tag("element", e.tag),
		)
	}

	e.children = FilterUIElems(v...)
}

func (e *htmlElement) preRender(p Page) {
	for _, c := range e.getChildren() {
		c.preRender(p)
	}
}

func (e *htmlElement) onComponentEvent(le any) {
	for _, c := range e.getChildren() {
		c.onComponentEvent(le)
	}
}

func (e *htmlElement) html(w io.Writer) {
	io.WriteString(w, "<")
	io.WriteString(w, e.tag)

	for k, v := range e.attributes {
		io.WriteString(w, " ")
		io.WriteString(w, k)

		if v != "" {
			io.WriteString(w, `="`)
			io.WriteString(w, resolveAttributeURLValue(k, v, func(s string) string {
				if e.dispatcher != nil {
					return e.dispatcher.resolveStaticResource(v)
				}
				return v
			}))
			io.WriteString(w, `"`)
		}
	}

	io.WriteString(w, ">")

	if e.isSelfClosing {
		return
	}

	for _, c := range e.children {
		io.WriteString(w, "\n")
		if c.self() == nil {
			c.setSelf(c)
		}
		c.html(w)
	}

	if len(e.children) != 0 {
		io.WriteString(w, "\n")
	}

	io.WriteString(w, "</")
	io.WriteString(w, e.tag)
	io.WriteString(w, ">")
}

func (e *htmlElement) htmlWithIndent(w io.Writer, indent int) {
	writeIndent(w, indent)
	io.WriteString(w, "<")
	io.WriteString(w, e.tag)

	for k, v := range e.attributes {
		io.WriteString(w, " ")
		io.WriteString(w, k)

		if v != "" {
			io.WriteString(w, `="`)
			io.WriteString(w, resolveAttributeURLValue(k, v, func(s string) string {
				if e.dispatcher != nil {
					return e.dispatcher.resolveStaticResource(v)
				}
				return v
			}))
			io.WriteString(w, `"`)
		}
	}

	io.WriteString(w, ">")

	if e.isSelfClosing {
		return
	}

	for _, c := range e.children {
		io.WriteString(w, "\n")
		if c.self() == nil {
			c.setSelf(c)
		}
		c.htmlWithIndent(w, indent+1)
	}

	if len(e.children) != 0 {
		io.WriteString(w, "\n")
		writeIndent(w, indent)
	}

	io.WriteString(w, "</")
	io.WriteString(w, e.tag)
	io.WriteString(w, ">")
}
