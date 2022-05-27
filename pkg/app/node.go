package app

import (
	"context"
	"io"
	"net/url"
	"reflect"
	"strings"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

// UI is the interface that describes a user interface element such as
// components and HTML elements.
type UI interface {
	// Kind represents the specific kind of a UI element.
	Kind() Kind

	// JSValue returns the javascript value linked to the element.
	JSValue() Value

	// Reports whether the element is mounted.
	IsMounted() bool

	name() string
	self() UI
	setSelf(UI)
	getContext() context.Context
	getDispatcher() Dispatcher
	getAttributes() map[string]string
	eventHandlers() map[string]eventHandler
	parent() UI
	setParent(UI)
	getChildren() []UI
	mount(Dispatcher) error
	dismount()
	update(UI) error
	onNav(*url.URL)
	onAppUpdate()
	onAppInstallChange()
	onResize()
	preRender(Page)
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
func FilterUIElems(uis ...UI) []UI {
	if len(uis) == 0 {
		return nil
	}

	elems := make([]UI, 0, len(uis))

	for _, n := range uis {
		// Ignore nil elements:
		if v := reflect.ValueOf(n); n == nil ||
			v.Kind() == reflect.Ptr && v.IsNil() {
			continue
		}

		switch n.Kind() {
		case SimpleText, HTML, Component, RawHTML:
			elems = append(elems, n)

		case Selector:
			elems = append(elems, n.getChildren()...)

		default:
			panic(errors.New("filtering ui elements failed").
				Tag("reason", "unexpected element type found").
				Tag("kind", n.Kind()).
				Tag("name", n.name()),
			)
		}
	}

	return elems
}

func isErrReplace(err error) bool {
	_, replace := errors.Tag(err, "replace")
	return replace
}

func mount(d Dispatcher, n UI) error {
	n.setSelf(n)
	return n.mount(d)
}

func dismount(n UI) {
	n.dismount()
	n.setSelf(nil)
}

func update(a, b UI) error {
	a.setSelf(a)
	b.setSelf(b)
	return a.update(b)
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
	if !ui.IsMounted() {
		ui.setSelf(ui)
	}
	ui.html(w)
}

// PrintHTMLWithIndent writes an idented HTML representation of the UI element
// into the given writer.
func PrintHTMLWithIndent(w io.Writer, ui UI) {
	if !ui.IsMounted() {
		ui.setSelf(ui)
	}
	ui.htmlWithIndent(w, 0)
}
