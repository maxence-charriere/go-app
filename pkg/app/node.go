package app

import (
	"io"
	"reflect"
	"strings"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

// UI is the interface that describes a user interface element such as
// components and HTML elements.
type UI interface {
	// JSValue returns the javascript value linked to the element.
	JSValue() Value

	// Reports whether the element is mounted.
	Mounted() bool

	name() string
	self() UI
	setSelf(UI)
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

// nodeManager manages the lifecycle of UI elements. It handles the logic for
// mounting, dismounting, and updating nodes based on their type.
type nodeManager struct {
}

// Mount mounts a UI element based on its type and the specified depth.
// It returns the mounted UI element and any potential error during the process.
func (m nodeManager) Mount(depth uint, v UI) (UI, error) {
	switch v := v.(type) {
	case *text:
		return m.mountText(depth, v)

	case HTML:
		return m.mountHTMLElement(depth, v)

	case Composer:
		return m.mountComponent(depth, v)

	case *raw:
		return m.mountRawHTMLElement(depth, v)

	default:
		return nil, errors.New("unsupported element").
			WithTag("type", reflect.TypeOf(v)).
			WithTag("depth", depth)
	}
}

func (m nodeManager) mountText(depth uint, v *text) (UI, error) {
	panic("not implemented")
}

func (m nodeManager) mountHTMLElement(depth uint, v HTML) (UI, error) {
	panic("not implemented")
}

func (m nodeManager) mountComponent(depth uint, v Composer) (UI, error) {
	panic("not implemented")
}

func (m nodeManager) mountRawHTMLElement(depth uint, v *raw) (UI, error) {
	panic("not implemented")
}

// Dismount removes a UI element based on its type.
func (m nodeManager) Dismount(v UI) {
	switch v := v.(type) {
	case *text:
		m.dismountText(v)

	case HTML:
		m.dismountHTMLElement(v)

	case Composer:
		m.dismountComponent(v)

	case *raw:
		m.dismountRawHTMLElement(v)
	}
}

func (m nodeManager) dismountText(v *text) {
	panic("not implemented")
}

func (m nodeManager) dismountHTMLElement(v HTML) {
	panic("not implemented")
}

func (m nodeManager) dismountComponent(v Composer) error {
	panic("not implemented")
}

func (m nodeManager) dismountRawHTMLElement(v *raw) error {
	panic("not implemented")
}
