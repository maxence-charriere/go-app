package app

import (
	"context"
	"io"
	"net/url"
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
	return e.getDispatcher() != nil &&
		e.context != nil &&
		e.context.Err() == nil &&
		e.self() != nil &&
		e.jsElement != nil
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

	e.dispatcher = d
	e.context, e.contextCancel = context.WithCancel(context.Background())

	jsElement, err := Window().createElement(e.tag, "")
	if err != nil {
		return errors.New("mounting ui element failed").
			Tag("name", e.name()).
			Tag("kind", e.Kind()).
			Wrap(err)
	}
	e.jsElement = jsElement

	e.attributes.Mount(jsElement, d.resolveStaticResource)
	e.eventHandlers.Mount(e)

	for _, c := range e.getChildren() {
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
	for _, c := range e.getChildren() {
		dismount(c)
	}

	for _, eh := range e.eventHandlers {
		eh.Dismount()
	}

	e.contextCancel()
	e.jsElement = nil
}

func (e *elem) canUpdateWith(n UI) bool {
	return n.Kind() == e.Kind() && n.name() == e.name()
}

func (e *elem) updateWith(n UI) error {
	if !e.IsMounted() {
		return nil
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

	achildren := e.getChildren()
	bchildren := n.getChildren()
	i := 0

	// Update children:
	for len(achildren) != 0 && len(bchildren) != 0 {
		a := achildren[0]
		b := bchildren[0]

		var err error
		if canUpdate(a, b) {
			err = update(a, b)
		} else {
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
	if err := mount(e.getDispatcher(), c); err != nil {
		return errors.New("appending child failed").
			Tag("name", e.name()).
			Tag("kind", e.Kind()).
			Tag("child-name", c.name()).
			Tag("child-kind", c.Kind()).
			Wrap(err)
	}

	if !onlyJsValue {
		e.children = append(e.children, c)
	}

	c.setParent(e.self())
	e.JSValue().appendChild(c)
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

func (e *elem) removeChildAt(idx int) error {
	body := e.children
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
	e.children = body

	e.JSValue().removeChild(c)
	dismount(c)
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
