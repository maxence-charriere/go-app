package app

import (
	"reflect"

	"github.com/maxence-charriere/app/pkg/log"
)

// IfCondition represents a control structure that displays nodes depending on a
// given expression.
type IfCondition struct {
	body []ValueNode
	eval bool
}

func (c IfCondition) nodeType() reflect.Type {
	return reflect.TypeOf(c)
}

func (c IfCondition) nodes() []ValueNode {
	return c.body
}

// If returns a condition that whether contains the given nodes depending on the
// given expression.
func If(expr bool, nodes ...Node) IfCondition {
	if !expr {
		nodes = nil
	}

	return IfCondition{
		body: indirect(nodes...),
		eval: !expr,
	}
}

// ElseIf sets the condition with the given nodes if previous expressions were
// not met and given expression is true.
func (c IfCondition) ElseIf(expr bool, nodes ...Node) IfCondition {
	if !c.eval {
		return c
	}

	if expr {
		c.body = indirect(nodes...)
		c.eval = false
	}

	return c
}

// Else sets the condition with the given UI elements if previous expressions
// were not met.
func (c IfCondition) Else(nodes ...Node) IfCondition {
	return c.ElseIf(true, nodes...)
}

// RangeCondition represents a control structure that iterates within a slice,
// an array or a map.
type RangeCondition struct {
	source interface{}
	body   []ValueNode
}

func (c RangeCondition) nodeType() reflect.Type {
	return reflect.TypeOf(c)
}

func (c RangeCondition) nodes() []ValueNode {
	return c.body
}

// Range returns a condition that iterates within the given source. Source must
// be a slice, an array or a map with strings as keys.
func Range(src interface{}) RangeCondition {
	return RangeCondition{source: src}
}

// Slice set the nodes of the condition by repeating the given function for
// the amount of elements in the source.
//
// It panics if the range source is not a map or an array.
func (c RangeCondition) Slice(f func(int) Node) RangeCondition {
	v := reflect.ValueOf(c.source)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		log.Error("range source is not a slice or array").
			T("source-type", v.Type()).
			Panic()
	}

	c.body = nil
	for i := 0; i < v.Len(); i++ {
		c.body = append(c.body, indirect(f(i))...)
	}
	return c
}

// Map set the nodes of the condition by repeating the given function for the
// amount of elements in the source.
//
// It panics if the range source is not a map or if map keys are not strings.
func (c RangeCondition) Map(f func(string) Node) RangeCondition {
	v := reflect.ValueOf(c.source)
	if v.Kind() != reflect.Map {
		log.Error("range source is not a map").
			T("source-type", v.Type()).
			Panic()
	}

	c.body = nil
	for _, key := range v.MapKeys() {
		if key.Kind() != reflect.String {
			log.Error("range source keys is not a string").
				T("key-type", key.Type()).
				Panic()
		}

		c.body = append(c.body, indirect(f(key.String()))...)
	}
	return c
}
