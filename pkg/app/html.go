package app

import (
	"context"
	"io"
	"net/url"
	"reflect"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

type htmlElement[T any] struct {
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
}

func (e *htmlElement[T]) Attr(name string, value any) T {
	if e.attributes == nil {
		e.attributes = make(attributes)
	}
	e.attributes.Set(name, value)

	return e.toHTMLInterface()
}

func (e *htmlElement[T]) On(event string, h EventHandler, scope ...any) T {
	if e.eventHandlers == nil {
		e.eventHandlers = make(eventHandlers)
	}
	e.eventHandlers.Set(event, h, scope...)

	return e.toHTMLInterface()
}

func (e *htmlElement[T]) Text(v any) T {
	switch e.tag {
	case "textarea":
		return e.Attr("value", v)

	default:
		return e.Body(Text(v))
	}
}

func (e *htmlElement[T]) Body(v ...UI) T {
	return e.setChildren(v...)
}

func (e *htmlElement[T]) JSValue() Value {
	return e.jsElement
}

func (e *htmlElement[T]) IsMounted() bool {
	return e.context != nil && e.context.Err() == nil
}

func (e *htmlElement[T]) toHTMLInterface() T {
	var i any = e
	return i.(T)
}

func (e *htmlElement[T]) getContext() context.Context {
	return e.context
}

func (e *htmlElement[T]) getDispatcher() Dispatcher {
	return e.dispatcher
}

func (e *htmlElement[T]) getAttributes() map[string]string {
	return e.attributes
}

func (e *htmlElement[T]) getEventHandlers() map[string]eventHandler {
	return e.eventHandlers
}

func (e *htmlElement[T]) getParent() UI {
	return e.parent
}

func (e *htmlElement[T]) setParent(v UI) {
	e.parent = v
}

func (e *htmlElement[T]) getChildren() []UI {
	return e.children
}

func (e *htmlElement[T]) setChildren(v ...UI) T {
	if e.isSelfClosing {
		panic(errors.New("cannot set children of a self closing element").
			Tag("tag", e.tag),
		)
	}

	e.children = FilterUIElems(v...)
	return e.toHTMLInterface()
}

func (e *htmlElement[T]) mount(d Dispatcher) error {
	if e.IsMounted() {
		return errors.New("html element is already mounted").Tag("tag", e.tag)
	}

	e.context, e.contextCancel = context.WithCancel(context.Background())
	e.dispatcher = d

	jsElement, err := Window().createElement(e.tag, e.xmlns)
	if err != nil {
		return errors.New("creating javascript element failed").
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
				Tag("child-type", reflect.TypeOf(c)).
				Wrap(err)
		}

		c.setParent(e)
		e.JSValue().appendChild(c)
	}

	return nil
}

func (e *htmlElement[T]) dismount() {
	for _, c := range e.children {
		dismount(c)
	}

	for _, eh := range e.eventHandlers {
		eh.Dismount()
	}

	e.contextCancel()
	e.jsElement = nil
}

func (e *htmlElement[T]) canUpdateWith(v UI) bool {
	if v, ok := v.(*htmlElement[T]); ok {
		return ok && e.tag == v.tag
	}
	return false
}

func (e *htmlElement[T]) updateWith(v UI) error {
	if e.IsMounted() {
		return nil
	}

	newElement, ok := v.(*htmlElement[T])
	if !ok {
		return errors.New("new element is not an html element").
			Tag("new-element-type", reflect.TypeOf(v))
	}

	if e.attributes == nil && newElement.attributes != nil {
		e.attributes = newElement.attributes
		e.attributes.Mount(e.jsElement, e.dispatcher.resolveStaticResource)
	} else if e.attributes != nil && newElement.attributes != nil {
		e.attributes.Update(
			e.jsElement,
			newElement.attributes,
			e.getDispatcher().resolveStaticResource,
		)
	}

	if e.eventHandlers == nil && newElement.eventHandlers != nil {
		e.eventHandlers = newElement.eventHandlers
		e.eventHandlers.Mount(e)
	} else if e.eventHandlers != nil && newElement.eventHandlers != nil {
		e.eventHandlers.Update(e, newElement.eventHandlers)
	}

	childrenA := e.children
	childrenB := newElement.children
	i := 0

	for len(childrenA) != 0 && len(childrenB) != 0 {
		a := childrenA[0]
		b := childrenB[0]

		if canUpdate(a, b) {
			if err := update(a, b); err != nil {
				return errors.New("updating child failed").
					Tag("child-type", reflect.TypeOf(a)).
					Tag("new-child-type", reflect.TypeOf(b)).
					Tag("index", i).
					Wrap(err)
			}
		} else {
			if err := e.replaceChildAt(i, b); err != nil {
				return errors.New("replacing child failed").
					Tag("child-type", reflect.TypeOf(a)).
					Tag("new-child-type", reflect.TypeOf(b)).
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
				Tag("child-type", reflect.TypeOf(childrenA[0])).
				Wrap(err)
		}

		childrenA = childrenA[1:]
	}

	for len(childrenB) != 0 {
		b := childrenB[0]

		if err := e.appendChild(b); err != nil {
			return errors.New("appending child failed").
				Tag("child-type", reflect.TypeOf(b)).
				Wrap(err)
		}

		childrenB = childrenB[1:]
	}

	return nil
}

func (e *htmlElement[T]) replaceChildAt(i int, new UI) error {
	if i < 0 || i >= len(e.children) {
		return errors.New("index out of range").
			Tag("index", i).
			Tag("children-count", len(e.children))
	}

	if err := mount(e.dispatcher, new); err != nil {
		return errors.New("mounting new element failed").
			Tag("element-type", reflect.TypeOf(new)).
			Wrap(err)
	}

	old := e.children[i]
	defer dismount(old)

	new.setParent(e)
	e.children[i] = new
	e.jsElement.replaceChild(new, old)
	return nil
}

func (e *htmlElement[T]) removeChildAt(i int) error {
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

func (e *htmlElement[T]) appendChild(v UI) error {
	if err := mount(e.dispatcher, v); err != nil {
		return errors.New("mounting element failed").
			Tag("element-type", reflect.TypeOf(v)).
			Wrap(err)
	}

	v.setParent(e)
	e.jsElement.appendChild(v)
	e.children = append(e.children, v)
	return nil
}

func (e *htmlElement[T]) onNav(u *url.URL) {
	for _, c := range e.children {
		c.onNav(u)
	}
}

func (e *htmlElement[T]) onAppUpdate() {
	for _, c := range e.children {
		c.onAppUpdate()
	}
}

func (e *htmlElement[T]) onAppInstallChange() {
	for _, c := range e.children {
		c.onAppInstallChange()
	}
}

func (e *htmlElement[T]) onResize() {
	for _, c := range e.children {
		c.onResize()
	}
}

func (e *htmlElement[T]) preRender(p Page) {
	for _, c := range e.children {
		c.preRender(p)
	}
}

func (e *htmlElement[T]) html(w io.Writer) {
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

func (e *htmlElement[T]) htmlWithIndent(w io.Writer, indent int) {
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

// -----------------------------------------------------------------------------
// The method below might be removed in later versions.
// -----------------------------------------------------------------------------
func (e *htmlElement[T]) kind() Kind {
	return HTML
}

func (e *htmlElement[T]) name() string {
	return e.tag
}

func (e *htmlElement[T]) self() UI {
	if e.IsMounted() {
		return e
	}
	return nil
}

func (e *htmlElement[T]) setSelf(UI) {
}
