package app

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"reflect"

	"github.com/maxence-charriere/go-app/v7/pkg/errors"
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
	context() context.Context
	attributes() map[string]string
	eventHandlers() map[string]eventHandler
	parent() UI
	setParent(UI)
	children() []UI
	mount() error
	dismount()
	update(UI) error
	onNav(*url.URL)
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
			elems = append(elems, n.children()...)

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

// EventHandler represents a function that can handle HTML events. They are
// always called on the UI goroutine.
type EventHandler func(ctx Context, e Event)

type eventHandler struct {
	event   string
	jsvalue Func
	value   EventHandler
}

func (h eventHandler) equal(o eventHandler) bool {
	return h.event == o.event &&
		fmt.Sprintf("%p", h.value) == fmt.Sprintf("%p", o.value)
}

func makeJsEventHandler(src UI, h EventHandler) Func {
	return FuncOf(func(this Value, args []Value) interface{} {
		dispatch(func() {
			if !src.Mounted() {
				return
			}

			ctx := Context{
				Context: src.context(),
				Src:     src,
				JSSrc:   src.JSValue(),
			}

			event := Event{
				Value: args[0],
			}

			trackMousePosition(event)
			h(ctx, event)
		})

		return nil
	})
}

func trackMousePosition(e Event) {
	x := e.Get("clientX")
	if !x.Truthy() {
		return
	}

	y := e.Get("clientY")
	if !y.Truthy() {
		return
	}

	Window().setCursorPosition(x.Int(), y.Int())
}

func isErrReplace(err error) bool {
	_, replace := errors.Tag(err, "replace")
	return replace
}

func mount(n UI) error {
	n.setSelf(n)
	return n.mount()
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

type writableNode interface {
	html(w io.Writer)
	htmlWithIndent(w io.Writer, indent int)
}
