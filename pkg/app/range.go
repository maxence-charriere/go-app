package app

import (
	"context"
	"io"
	"reflect"
	"sort"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
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
}

// Range returns a range loop that iterates within the given source. Source must
// be a slice, an array or a map with strings as keys.
func Range(src any) RangeLoop {
	return rangeLoop{source: src}
}

type rangeLoop struct {
	body   []UI
	source any
}

func (r rangeLoop) Slice(f func(int) UI) RangeLoop {
	src := reflect.ValueOf(r.source)
	if src.Kind() != reflect.Slice && src.Kind() != reflect.Array {
		panic(errors.New("range loop source is not a slice or array").
			Tag("src-type", src.Type),
		)
	}

	body := make([]UI, 0, src.Len())
	for i := 0; i < src.Len(); i++ {
		body = append(body, FilterUIElems(f(i))...)
	}

	r.body = body
	return r
}

func (r rangeLoop) Map(f func(string) UI) RangeLoop {
	src := reflect.ValueOf(r.source)
	if src.Kind() != reflect.Map {
		panic(errors.New("range loop source is not a map").
			Tag("src-type", src.Type),
		)
	}

	if keyType := src.Type().Key(); keyType.Kind() != reflect.String {
		panic(errors.New("range loop source keys are not strings").
			Tag("src-type", src.Type).
			Tag("key-type", keyType),
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

	r.body = body
	return r
}

func (r rangeLoop) Kind() Kind {
	return Selector
}

func (r rangeLoop) JSValue() Value {
	return nil
}

func (r rangeLoop) Mounted() bool {
	return false
}

func (r rangeLoop) name() string {
	return "range"
}

func (r rangeLoop) self() UI {
	return r
}

func (r rangeLoop) setSelf(UI) {
}

func (r rangeLoop) getContext() context.Context {
	return nil
}

func (r rangeLoop) getDispatcher() Dispatcher {
	return nil
}

func (r rangeLoop) getAttributes() attributes {
	return nil
}

func (r rangeLoop) getEventHandlers() eventHandlers {
	return nil
}

func (r rangeLoop) getParent() UI {
	return nil
}

func (r rangeLoop) setParent(UI) {
}

func (r rangeLoop) getChildren() []UI {
	return r.body
}

func (r rangeLoop) mount(Dispatcher) error {
	return errors.New("range loop is not mountable").
		Tag("name", r.name()).
		Tag("kind", r.Kind())
}

func (r rangeLoop) dismount() {
}

func (r rangeLoop) canUpdateWith(UI) bool {
	return false
}

func (r rangeLoop) updateWith(UI) error {
	return errors.New("range loop cannot be updated").
		Tag("name", r.name()).
		Tag("kind", r.Kind())
}

func (r rangeLoop) preRender(Page) {
}

func (r rangeLoop) onComponentEvent(any) {
}

func (r rangeLoop) html(w io.Writer) {
	panic("should not be called")
}

func (r rangeLoop) htmlWithIndent(w io.Writer, indent int) {
	panic("should not be called")
}
