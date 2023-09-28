package app

import (
	"context"
	"io"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

// Condition represents a control structure for conditionally displaying UI
// elements. It extends the UI interface to include methods for handling
// conditional logic.
type Condition interface {
	UI

	// ElseIf sets a UI element to be displayed when the given boolean
	// expression is true and all previous conditions have been false.
	//
	// expr: Boolean expression to evaluate.
	// elem: Function that returns the UI element to display.
	ElseIf(expr bool, elem func() UI) Condition

	// ElseIfSlice sets multiple UI elements to be displayed when the given
	// boolean expression is true and all previous conditions have been false.
	//
	// expr: Boolean expression to evaluate.
	// elems: Function that returns a slice of UI elements to display.
	ElseIfSlice(expr bool, elems func() []UI) Condition

	// Else sets a UI element to be displayed as a fallback when all previous
	// conditions have been false.
	//
	// elem: Function that returns the UI element to display.
	Else(elem func() UI) Condition

	// ElseSlice sets multiple UI elements to be displayed as a fallback when
	// all previous conditions have been false.
	//
	// expr: Boolean expression to evaluate.
	// elems: Function that returns a slice of UI elements to display.
	ElseSlice(elems func() []UI) Condition
}

// If returns a Condition that will display the given UI element based on the
// evaluation of the provided boolean expression.
func If(expr bool, elem func() UI) Condition {
	return IfSlice(expr, func() []UI {
		return []UI{elem()}
	})
}

// IfSlice returns a Condition that will display the given slice of UI elements
// based on the evaluation of the provided boolean expression.
func IfSlice(expr bool, elems func() []UI) Condition {
	if !expr {
		return condition{}
	}
	return condition{body: FilterUIElems(elems()...)}
}

type condition struct {
	body []UI
}

func (c condition) ElseIf(expr bool, elem func() UI) Condition {
	return c.ElseIfSlice(expr, func() []UI {
		return []UI{elem()}
	})
}

func (c condition) ElseIfSlice(expr bool, elems func() []UI) Condition {
	if len(c.body) != 0 || !expr {
		return c
	}

	c.body = FilterUIElems(elems()...)
	return c
}

func (c condition) Else(elem func() UI) Condition {
	return c.ElseSlice(func() []UI {
		return []UI{elem()}
	})
}

func (c condition) ElseSlice(elems func() []UI) Condition {
	return c.ElseIfSlice(true, elems)
}

func (c condition) JSValue() Value {
	return nil
}

func (c condition) Mounted() bool {
	return false
}

func (c condition) name() string {
	return "if.else"
}

func (c condition) self() UI {
	return c
}

func (c condition) setSelf(UI) {
}

func (c condition) getContext() context.Context {
	return nil
}

func (c condition) getDispatcher() Dispatcher {
	return nil
}

func (c condition) getAttributes() attributes {
	return nil
}

func (c condition) getEventHandlers() eventHandlers {
	return nil
}

func (c condition) getParent() UI {
	return nil
}

func (c condition) setParent(UI) {
}

func (c condition) getChildren() []UI {
	return c.body
}

func (c condition) mount(Dispatcher) error {
	return errors.New("condition is not mountable").WithTag("name", c.name())
}

func (c condition) dismount() {
}

func (c condition) canUpdateWith(UI) bool {
	return false
}

func (c condition) updateWith(UI) error {
	return errors.New("condition cannot be updated").WithTag("name", c.name())
}

func (c condition) onComponentEvent(any) {
}

func (c condition) html(w io.Writer) {
	panic("shoulnd not be called")
}

func (c condition) htmlWithIndent(w io.Writer, indent int) {
	panic("shoulnd not be called")
}
