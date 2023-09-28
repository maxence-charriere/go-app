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

// FilterUIElems processes and returns a filtered list of the provided UI
// elements.
//
// Specifically, it:
// - Interprets and removes selector elements such as Condition and RangeLoop.
// - Eliminates nil elements and nil pointers.
// - Flattens and includes the children of recognized selector elements.
//
// This function is primarily intended for components that accept ui elements as
// variadic arguments or slice, such as the Body method of HTML elements.
func FilterUIElems(v ...UI) []UI {
	if len(v) == 0 {
		return nil
	}

	removeELemAt := func(i int) {
		copy(v[i:], v[i+1:])
		v[len(v)-1] = nil
		v = v[:len(v)-1]
	}

	var trailing []UI
	replaceElemAt := func(i int, elems ...UI) {
		trailing = append(trailing, v[i+1:]...)
		v = append(v[:i], elems...)
		v = append(v, trailing...)
		trailing = trailing[:0]
	}

	for i := len(v) - 1; i >= 0; i-- {
		elem := v[i]
		if elem == nil {
			removeELemAt(i)
		}
		if elemValue := reflect.ValueOf(elem); elemValue.Kind() == reflect.Pointer && elemValue.IsNil() {
			removeELemAt(i)
		}

		switch elem.(type) {
		case Condition, RangeLoop:
			replaceElemAt(i, elem.getChildren()...)
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
