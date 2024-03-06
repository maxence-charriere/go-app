package app

import (
	"reflect"
	"sort"

	"github.com/maxence-charriere/go-app/v10/pkg/errors"
)

// RangeLoop represents a control structure that iterates within a slice, an
// array or a map.
type RangeLoop interface {
	UI

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

	body() []UI
}

// Range returns a range loop that iterates within the given source. Source must
// be a slice, an array or a map with strings as keys.
func Range(src any) RangeLoop {
	return rangeLoop{source: src}
}

type rangeLoop struct {
	children []UI
	source   any
}

func (r rangeLoop) Slice(f func(int) UI) RangeLoop {
	src := reflect.ValueOf(r.source)
	if src.Kind() != reflect.Slice && src.Kind() != reflect.Array {
		panic(errors.New("range loop source is not a slice or array").
			WithTag("src-type", src.Type),
		)
	}

	body := make([]UI, 0, src.Len())
	for i := 0; i < src.Len(); i++ {
		body = append(body, FilterUIElems(f(i))...)
	}

	r.children = body
	return r
}

func (r rangeLoop) Map(f func(string) UI) RangeLoop {
	src := reflect.ValueOf(r.source)
	if src.Kind() != reflect.Map {
		panic(errors.New("range loop source is not a map").
			WithTag("src-type", src.Type),
		)
	}

	if keyType := src.Type().Key(); keyType.Kind() != reflect.String {
		panic(errors.New("range loop source keys are not strings").
			WithTag("src-type", src.Type).
			WithTag("key-type", keyType),
		)
	}

	body := make([]UI, 0, src.Len())
	keys := make([]string, 0, src.Len())

	for _, k := range src.MapKeys() {
		keys = append(keys, k.String())
	}
	sort.Strings(keys)

	for _, k := range keys {
		body = append(body, FilterUIElems(f(k))...)
	}

	r.children = body
	return r
}

func (r rangeLoop) JSValue() Value {
	return nil
}

func (r rangeLoop) Mounted() bool {
	return false
}

func (r rangeLoop) setParent(UI) UI {
	return nil
}

func (r rangeLoop) parent() UI {
	return nil
}

func (r rangeLoop) body() []UI {
	return r.children
}
