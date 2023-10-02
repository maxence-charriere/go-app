package app

import (
	"io"
	"reflect"
	"strconv"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

// HTML represents an interface for HTML elements.
type HTML interface {
	UI

	// Returns the name of the HTML tag represented by the element.
	Tag() string

	// Returns the XML namespace of the HTML element.
	XMLNamespace() string

	// Indicates whether the HTML element is self-closing (like <img> or <br>).
	// Returns true for self-closing elements, otherwise false.
	SelfClosing() bool

	// Returns the nesting level of the HTML element. A higher value
	// indicates deeper nesting within the document.
	Depth() uint

	attrs() attributes
	events() eventHandlers
	setDepth(uint)
	setJSElement(Value)
	parent() UI
	body() []UI
}

type htmlElement struct {
	attributes    attributes
	eventHandlers eventHandlers
	children      []UI

	dispatcher Dispatcher
	this       UI

	tag           string
	xmlns         string
	depth         uint
	isSelfClosing bool
	jsElement     Value
	parentElement UI
}

func (e *htmlElement) JSValue() Value {
	return e.jsElement
}

func (e *htmlElement) Mounted() bool {
	return e.jsElement != nil
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

func (e *htmlElement) getDispatcher() Dispatcher {
	return e.dispatcher
}

func (e *htmlElement) getAttributes() attributes {
	return e.attrs()
}

func (e *htmlElement) getEventHandlers() eventHandlers {
	return e.events()
}

func (e *htmlElement) getParent() UI {
	return e.parentElement
}

func (e *htmlElement) getChildren() []UI {
	return e.body()
}

func (e *htmlElement) mount(d Dispatcher) error {
	if e.Mounted() {
		return errors.New("html element is already mounted").WithTag("tag", e.tag)
	}

	e.dispatcher = d

	jsElement, err := Window().createElement(e.tag, e.xmlns)
	if err != nil {
		return errors.New("mounting js element failed").
			WithTag("tag", e.tag).
			WithTag("xmlns", e.xmlns).
			Wrap(err)
	}
	e.jsElement = jsElement

	e.attributes.Mount(jsElement, d.resolveStaticResource)
	e.eventHandlers.Mount(e)

	for i, c := range e.children {
		if err := mount(d, c); err != nil {
			return errors.New("mounting child failed").
				WithTag("index", i).
				WithTag("child", c.name()).
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

	e.jsElement = nil
}

func (e *htmlElement) canUpdateWith(v UI) bool {
	return e.Mounted() && e.name() == v.name()
}

func (e *htmlElement) updateWith(v UI) error {
	if !e.canUpdateWith(v) {
		return errors.New("cannot update html element with given element").
			WithTag("current", reflect.TypeOf(e.self())).
			WithTag("new", reflect.TypeOf(v))
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
					WithTag("child", reflect.TypeOf(a)).
					WithTag("new-child", reflect.TypeOf(b)).
					WithTag("index", i).
					Wrap(err)
			}
		} else {
			if err := e.replaceChildAt(i, b); err != nil {
				return errors.New("replacing child failed").
					WithTag("child", reflect.TypeOf(a)).
					WithTag("new-child", reflect.TypeOf(b)).
					WithTag("index", i).
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
				WithTag("child", reflect.TypeOf(childrenA[0])).
				WithTag("index", i).
				Wrap(err)
		}

		childrenA = childrenA[1:]
	}

	for len(childrenB) != 0 {
		b := childrenB[0]

		if err := e.appendChild(b); err != nil {
			return errors.New("appending child failed").
				WithTag("child", reflect.TypeOf(b)).
				WithTag("index", i).
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
			WithTag("name", e.name()).
			WithTag("index", idx).
			WithTag("old-name", old.name()).
			WithTag("new-name", new.name()).
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
			WithTag("index", i).
			WithTag("children-count", len(e.children))
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
			WithTag("element", reflect.TypeOf(v)).
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

// TODO: Remove
func (e *htmlElement) setParent(v UI) UI {
	e.parentElement = v
	return nil
}

func (e *htmlElement) setChildren(v ...UI) {
	if e.isSelfClosing {
		panic(errors.New("cannot set children of a self closing element").
			WithTag("element", e.tag),
		)
	}

	e.children = FilterUIElems(v...)
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
		e.writeHTMLAttribute(w, k, v)
	}

	io.WriteString(w, ">")

	if e.isSelfClosing {
		return
	}

	hasNewLineChildren := len(e.children) > 1

	for _, c := range e.children {
		if hasNewLineChildren {
			io.WriteString(w, "\n")
		}

		if c.self() == nil {
			c.setSelf(c)
		}
		c.html(w)
	}

	if hasNewLineChildren {
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
		e.writeHTMLAttribute(w, k, v)
	}

	io.WriteString(w, ">")

	if e.isSelfClosing {
		return
	}

	var hasNewLineChildren bool
	if len(e.children) > 0 {
		_, isText := e.children[0].(*text)
		hasNewLineChildren = len(e.children) > 1 || !isText
	}

	for _, c := range e.children {
		if hasNewLineChildren {
			io.WriteString(w, "\n")
		}

		if c.self() == nil {
			c.setSelf(c)
		}
		c.htmlWithIndent(w, indent+1)
	}

	if hasNewLineChildren {
		io.WriteString(w, "\n")
		writeIndent(w, indent)
	}

	io.WriteString(w, "</")
	io.WriteString(w, e.tag)
	io.WriteString(w, ">")
}

func (e *htmlElement) writeHTMLAttribute(w io.Writer, k, v string) {
	if (k == "id" || k == "class") && v == "" {
		return
	}

	io.WriteString(w, " ")
	io.WriteString(w, k)

	if v != "" && v != "true" {
		io.WriteString(w, `=`)
		io.WriteString(w, strconv.Quote(resolveAttributeURLValue(k, v, func(s string) string {
			if e.dispatcher != nil {
				return e.dispatcher.resolveStaticResource(s)
			}
			return s
		})))
	}
}

// -----------------------------------------------------------------------------

func (e *htmlElement) Tag() string {
	return e.tag
}

func (e *htmlElement) XMLNamespace() string {
	return e.xmlns
}

func (e *htmlElement) SelfClosing() bool {
	return e.isSelfClosing
}

func (e *htmlElement) Depth() uint {
	return e.depth
}

func (e *htmlElement) attrs() attributes {
	return e.attributes
}

func (e *htmlElement) events() eventHandlers {
	return e.eventHandlers
}

func (e *htmlElement) setDepth(v uint) {
	e.depth = v
}

func (e *htmlElement) setJSElement(v Value) {
	e.jsElement = v
}

func (e *htmlElement) parent() UI {
	return e.parentElement
}

func (e *htmlElement) body() []UI {
	return e.children
}
