package app

import (
	"reflect"
)

// Condition represents a control structure that displays nodes depending on a
// given expression.
type Condition interface {
	Node

	// ElseIf sets the condition with the given nodes if previous expressions
	// were not met and given expression is true.
	ElseIf(expr bool, nodes ...Node) Condition

	// Else sets the condition with the given UI elements if previous
	// expressions were not met.
	Else(nodes ...Node) Condition

	isSatisfied() bool
	nodes() []UI
}

// If returns a condition that contains the given nodes depending on the given
// expression.
func If(expr bool, nodes ...Node) Condition {
	if !expr {
		nodes = nil
	}

	return condition{
		body:      Indirect(nodes...),
		satisfied: !expr,
	}
}

type condition struct {
	body      []UI
	satisfied bool
}

func (c condition) nodeType() reflect.Type {
	return reflect.TypeOf(c)
}

func (c condition) isSatisfied() bool {
	return c.satisfied
}

func (c condition) nodes() []UI {
	return c.body
}

func (c condition) ElseIf(expr bool, nodes ...Node) Condition {
	if !c.satisfied {
		return c
	}

	if expr {
		c.body = Indirect(nodes...)
		c.satisfied = false
	}

	return c
}

func (c condition) Else(nodes ...Node) Condition {
	return c.ElseIf(true, nodes...)
}
