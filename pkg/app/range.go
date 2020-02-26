package app

import (
	"reflect"

	"github.com/maxence-charriere/app/pkg/log"
)

// RangeCondition represents a control structure that iterates within a slice,
// an array or a map.
type RangeCondition struct {
	source interface{}
	body   []UI
}

func (c RangeCondition) nodeType() reflect.Type {
	return reflect.TypeOf(c)
}

func (c RangeCondition) nodes() []UI {
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
