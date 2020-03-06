package app

import (
	"reflect"
	"sort"

	"github.com/maxence-charriere/go-app/pkg/log"
)

// RangeLoop represents a control structure that iterates within a slice, an
// array or a map.
type RangeLoop interface {
	Node

	nodes() []UI

	// Slice sets the loop content by repeating the given function for the
	// number of elements in the source.
	//
	// It panics if the range source is not a slice or an array.
	Slice(f func(int) UI) RangeLoop

	// Map sets the loop content by repeating the given function for the number
	// of elements in the source. Elements are ordered by keys.
	//
	// It panics if the range source is not a map or if map keys are not strings.
	Map(f func(string) UI) RangeLoop
}

// Range returns a range loop that iterates within the given source. Source must
// be a slice, an array or a map with strings as keys.
func Range(src interface{}) RangeLoop {
	return rangeCondition{source: src}
}

type rangeCondition struct {
	source interface{}
	body   []UI
}

func (c rangeCondition) nodeType() reflect.Type {
	return reflect.TypeOf(c)
}

func (c rangeCondition) nodes() []UI {
	return c.body
}

func (c rangeCondition) Slice(f func(int) UI) RangeLoop {
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

func (c rangeCondition) Map(f func(string) UI) RangeLoop {
	v := reflect.ValueOf(c.source)
	if v.Kind() != reflect.Map {
		log.Error("range source is not a map").
			T("source-type", v.Type()).
			Panic()
	}
	if keyType := v.Type().Key(); keyType.Kind() != reflect.String {
		log.Error("range source keys are not strings").
			T("key-type", keyType).
			Panic()
	}

	c.body = nil

	keys := make([]string, 0, v.Len())
	for _, k := range v.MapKeys() {
		keys = append(keys, k.String())
	}
	sort.Strings(keys)

	for _, k := range keys {
		c.body = append(c.body, indirect(f(k))...)
	}
	return c
}
