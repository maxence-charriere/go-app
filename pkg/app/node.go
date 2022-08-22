package app

import (
	"context"
	"io"
	"reflect"
	"strings"
)

// UI is the interface that describes a user interface element such as
// components and HTML elements.
type UI interface {
	// Kind represents the specific kind of a UI element.
	Kind() Kind

	// JSValue returns the javascript value linked to the element.
	JSValue() Value

	// Reports whether the element is mounted.
	Mounted() bool

	name() string
	self() UI
	setSelf(UI)
	getContext() context.Context
	getDispatcher() Dispatcher
	getAttributes() attributes
	getEventHandlers() eventHandlers
	getParent() UI
	setParent(UI)
	getChildren() []UI
	mount(Dispatcher) error
	dismount()
	canUpdateWith(UI) bool
	updateWith(UI) error
	preRender(Page)
	onComponentEvent(any)
	html(w io.Writer)
	htmlWithIndent(w io.Writer, indent int)
}

// Kind represents the specific kind of a user interface element.
type Kind uint

func (k Kind) String() string {
	switch k {
	case SimpleText:
		return "text"

	case HTML:
		return "html"

	case Component:
		return "component"

	case Selector:
		return "selector"

	case RawHTML:
		return "raw"

	default:
		return "undefined"
	}
}

const (
	// UndefinedElem represents an undefined UI element.
	UndefinedElem Kind = iota

	// SimpleText represents a simple text element.
	SimpleText

	// HTML represents an HTML element.
	HTML

	// Component represents a customized, independent and reusable UI element.
	Component

	// Selector represents an element that is used to select a subset of
	// elements within a given list.
	Selector

	// RawHTML represents an HTML element obtained from a raw HTML code snippet.
	RawHTML
)

// FilterUIElems returns a filtered version of the given UI elements where
// selector elements such as If and Range are interpreted and removed. It also
// remove nil elements.
//
// It should be used only when implementing components that can accept content
// with variadic arguments like HTML elements Body method.
func FilterUIElems(v ...UI) []UI {
	if len(v) == 0 {
		return nil
	}

	remove := func(i int) {
		copy(v[i:], v[i+1:])
		v[len(v)-1] = nil
		v = v[:len(v)-1]
	}

	var b []UI
	replaceAt := func(i int, s ...UI) {
		b = append(b, v[i+1:]...)
		v = append(v[:i], s...)
		v = append(v, b...)
		b = b[:0]
	}

	for i := len(v) - 1; i >= 0; i-- {
		e := v[i]
		if ev := reflect.ValueOf(e); e == nil || ev.Kind() == reflect.Pointer && ev.IsNil() {
			remove(i)
			continue
		}

		switch e.Kind() {
		case SimpleText, HTML, Component, RawHTML:

		case Selector:
			replaceAt(i, e.getChildren()...)

		default:
			remove(i)
		}
	}

	return v
}

func mount(d Dispatcher, n UI) error {
	n.setSelf(n)
	return n.mount(d)
}

func dismount(n UI) {
	n.dismount()
	n.setSelf(nil)
}

func canUpdate(a, b UI) bool {
	a.setSelf(a)
	b.setSelf(b)
	return a.canUpdateWith(b)
}

func update(a, b UI) error {
	a.setSelf(a)
	b.setSelf(b)
	return a.updateWith(b)
}

// HTMLString return an HTML string representation of the given UI element.
func HTMLString(ui UI) string {
	var w strings.Builder
	PrintHTML(&w, ui)
	return w.String()
}

// HTMLStringWithIndent return an indented HTML string representation of the
// given UI element.
func HTMLStringWithIndent(ui UI) string {
	var w strings.Builder
	PrintHTMLWithIndent(&w, ui)
	return w.String()
}

// PrintHTML writes an HTML representation of the UI element into the given
// writer.
func PrintHTML(w io.Writer, ui UI) {
	if !ui.Mounted() {
		ui.setSelf(ui)
	}
	ui.html(w)
}

// PrintHTMLWithIndent writes an idented HTML representation of the UI element
// into the given writer.
func PrintHTMLWithIndent(w io.Writer, ui UI) {
	if !ui.Mounted() {
		ui.setSelf(ui)
	}
	ui.htmlWithIndent(w, 0)
}
