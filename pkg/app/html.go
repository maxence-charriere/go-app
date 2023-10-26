package app

import (
	"io"
	"strconv"
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

	depth() uint
	attrs() attributes
	setAttrs(attributes) HTML
	events() eventHandlers
	setEvents(eventHandlers) HTML
	setDepth(uint) UI
	setJSElement(Value) HTML
	parent() UI
	body() []UI
	setBody([]UI) HTML
}

type htmlElement struct {
	dispatcher Dispatcher
	this       UI

	tag           string
	xmlns         string
	treeDepth     uint
	isSelfClosing bool
	jsElement     Value
	attributes    attributes
	eventHandlers eventHandlers
	parentElement UI
	children      []UI
}

func (e *htmlElement) JSValue() Value {
	return e.jsElement
}

func (e *htmlElement) Mounted() bool {
	return e.jsElement != nil
}

func (e *htmlElement) Tag() string {
	return e.tag
}

func (e *htmlElement) XMLNamespace() string {
	return e.xmlns
}

func (e *htmlElement) SelfClosing() bool {
	return e.isSelfClosing
}

func (e *htmlElement) depth() uint {
	return e.treeDepth
}

func (e *htmlElement) attrs() attributes {
	return e.attributes
}

func (e *htmlElement) setAttr(name string, value any) {
	if e.attributes == nil {
		e.attributes = make(attributes)
	}
	e.attributes.Set(name, value)
}

func (e *htmlElement) events() eventHandlers {
	return e.eventHandlers
}

func (e *htmlElement) setEventHandler(event string, h EventHandler, scope ...any) {
	if e.eventHandlers == nil {
		e.eventHandlers = make(eventHandlers)
	}
	e.eventHandlers.Set(event, h, scope...)
}

func (e *htmlElement) parent() UI {
	return e.parentElement
}

func (e *htmlElement) body() []UI {
	return e.children
}

func (e *htmlElement) html(w io.Writer) {
	panic("not implemented")

	// io.WriteString(w, "<")
	// io.WriteString(w, e.tag)

	// for k, v := range e.attributes {
	// 	e.writeHTMLAttribute(w, k, v)
	// }

	// io.WriteString(w, ">")

	// if e.isSelfClosing {
	// 	return
	// }

	// hasNewLineChildren := len(e.children) > 1

	// for _, c := range e.children {
	// 	if hasNewLineChildren {
	// 		io.WriteString(w, "\n")
	// 	}

	// 	if c.self() == nil {
	// 		c.setSelf(c)
	// 	}
	// 	c.html(w)
	// }

	// if hasNewLineChildren {
	// 	io.WriteString(w, "\n")
	// }

	// io.WriteString(w, "</")
	// io.WriteString(w, e.tag)
	// io.WriteString(w, ">")
}

func (e *htmlElement) htmlWithIndent(w io.Writer, indent int) {
	panic("not implemented")

	// writeIndent(w, indent)
	// io.WriteString(w, "<")
	// io.WriteString(w, e.tag)

	// for k, v := range e.attributes {
	// 	e.writeHTMLAttribute(w, k, v)
	// }

	// io.WriteString(w, ">")

	// if e.isSelfClosing {
	// 	return
	// }

	// var hasNewLineChildren bool
	// if len(e.children) > 0 {
	// 	_, isText := e.children[0].(*text)
	// 	hasNewLineChildren = len(e.children) > 1 || !isText
	// }

	// for _, c := range e.children {
	// 	if hasNewLineChildren {
	// 		io.WriteString(w, "\n")
	// 	}

	// 	if c.self() == nil {
	// 		c.setSelf(c)
	// 	}
	// 	c.htmlWithIndent(w, indent+1)
	// }

	// if hasNewLineChildren {
	// 	io.WriteString(w, "\n")
	// 	writeIndent(w, indent)
	// }

	// io.WriteString(w, "</")
	// io.WriteString(w, e.tag)
	// io.WriteString(w, ">")
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
