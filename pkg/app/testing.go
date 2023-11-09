package app

import (
	"reflect"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

// Match compares the expected UI element with another UI element at a specified
// location in a UI tree. It is the preferred function for matching UI elements
// in tests due to its simplified usage.
//
// Example usage adapted for Match function:
//
//	tree := app.Div().Body(
//	    app.H2().Body(app.Text("foo")),
//	    app.P().Body(app.Text("bar")),
//	)
//
//	err := app.Match(app.Div(), tree)
//	// err == nil if the root matches a Div element
//
//	err := app.Match(app.H3(), tree, 0)
//	// err != nil because the first child is not an H3 element but a H2.
//
//	err = app.Match(app.Text("bar"), tree, 1, 0)
//	// err == nil if the text of the first child of the second element is "bar"
func Match(expected UI, root UI, path ...int) error {
	return TestMatch(root, TestUIDescriptor{
		Path:     TestPath(path...),
		Expected: expected,
	})
}

// TestUIDescriptor describes a UI element and its hierarchical location
// relative to parent elements for the purpose of testing.
type TestUIDescriptor struct {
	// Path represents the sequence of child indices to navigate through the UI
	// tree to reach the element to be tested. An empty path implies the root.
	Path []int

	// Expected is the UI element that is expected to be found at the location
	// specified by Path. The comparison behavior varies depending on the type
	// of element; simple text elements are compared by text value, HTML
	// elements by attributes and event handlers, and components by the values
	// of their exported fields.
	Expected UI
}

// TestPath is a utility function that constructs a path, represented as a slice
// of integers, for use in a TestUIDescriptor.
func TestPath(p ...int) []int {
	return p
}

// TestMatch searches for a UI element within a tree as described by a
// TestUIDescriptor and verifies if it matches the Expected element. It returns
// an error if the match is unsuccessful or if the path is invalid. Prefer using
// the Match function for a simpler API.
func TestMatch(root UI, d TestUIDescriptor) error {
	if len(d.Path) != 0 {
		index := d.Path[0]

		switch root := root.(type) {
		case HTML:
			children := root.body()
			if index < 0 || index >= len(children) {
				return errors.New("element to match is out of range").
					WithTag("type", reflect.TypeOf(d.Expected)).
					WithTag("parent-type", reflect.TypeOf(root)).
					WithTag("parent-children-count", len(children)).
					WithTag("index", index)
			}
			d.Path = d.Path[1:]
			return TestMatch(children[index], d)

		case Composer:
			if index != 0 {
				return errors.New("element to match is out of range").
					WithTag("type", reflect.TypeOf(d.Expected)).
					WithTag("parent-type", reflect.TypeOf(root)).
					WithTag("parent-children-count", 1).
					WithTag("index", index)
			}
			d.Path = d.Path[1:]
			return TestMatch(root.root(), d)
		}
	}

	return match(root, d)
}

func match(n UI, d TestUIDescriptor) error {
	if a, b := reflect.TypeOf(d.Expected), reflect.TypeOf(n); a != b {
		return errors.New("types are not matching").
			WithTag("type", a).
			WithTag("expected-type", b)
	}

	switch d.Expected.(type) {
	case *text:
		return matchText(n.(*text), d)

	case HTML:
		return matchHTML(n.(HTML), d)

	case Composer:
		return matchComponent(n.(Composer), d)

	case *raw:
		return matchRaw(n.(*raw), d)

	default:
		return errors.New("unsupported element").
			WithTag("type", reflect.TypeOf(n))
	}
}

func matchText(n *text, d TestUIDescriptor) error {
	a := n
	b := d.Expected.(*text)

	if a.value != b.value {
		return errors.New("text does not match").
			WithTag("type", reflect.TypeOf(a)).
			WithTag("expected-value", b.value).
			WithTag("current-value", a.value)
	}
	return nil
}

func matchHTML(n HTML, d TestUIDescriptor) error {
	a := n
	b := d.Expected.(HTML)

	if typeA, typeB := reflect.TypeOf(a), reflect.TypeOf(b); typeA != typeB || a.Tag() != b.Tag() {
		return errors.New("types are not matching").
			WithTag("type", typeA).
			WithTag("expected-type", typeB)
	}

	if err := matchHTMLAttributes(a.attrs(), b.attrs()); err != nil {
		return errors.New("attributes does not match").
			WithTag("type", reflect.TypeOf(a)).
			Wrap(err)
	}

	if err := matchHTMLEventHandlers(a.events(), b.events()); err != nil {
		return errors.New("event handlers does not match").
			WithTag("type", reflect.TypeOf(a)).
			Wrap(err)
	}

	return nil
}

func matchHTMLAttributes(a, b attributes) error {
	for key, expectedValue := range b {
		value, exists := a[key]
		if !exists {
			return errors.New("expected attribute not found").
				WithTag("name", key)
		}

		if value != expectedValue {
			return errors.New("value does not match").
				WithTag("name", key).
				WithTag("value", value).
				WithTag("expected-value", expectedValue)
		}
	}

	for key, value := range a {
		if _, exists := b[key]; !exists {
			return errors.New("attribute is not expected").
				WithTag("name", key).
				WithTag("value", value)
		}
	}

	return nil
}

func matchHTMLEventHandlers(a, b eventHandlers) error {
	for key := range b {
		if _, exists := a[key]; !exists {
			return errors.New("expected event handler not found").
				WithTag("event", key)
		}
	}

	for key := range a {
		if _, exists := b[key]; !exists {
			return errors.New("event handler is not expected").
				WithTag("event", key)
		}
	}

	return nil
}

func matchComponent(n Composer, d TestUIDescriptor) error {
	a := reflect.Indirect(reflect.ValueOf(n))
	b := reflect.Indirect(reflect.ValueOf(d.Expected))

	for i := 0; i < b.NumField(); i++ {
		fieldA := a.Field(i)
		fieldB := b.Field(i)

		if !fieldA.CanSet() {
			continue
		}

		if _, ok := fieldA.Interface().(Compo); ok {
			continue
		}

		if !reflect.DeepEqual(fieldA.Interface(), fieldB.Interface()) {
			return errors.New("field are not matching").
				WithTag("type", reflect.TypeOf(n)).
				WithTag("field", a.Type().Field(i).Name).
				WithTag("value", fieldA.Interface()).
				WithTag("expected-value", fieldB.Interface())
		}
	}

	return nil
}

func matchRaw(n *raw, d TestUIDescriptor) error {
	a := n
	b := d.Expected.(*raw)

	if a.value != b.value {
		return errors.New("the raw html element is not matching with the descriptor").
			WithTag("type", reflect.TypeOf(n)).
			WithTag("reason", "unexpected value").
			WithTag("expected-value", b.value).
			WithTag("current-value", a.value)
	}

	return nil
}
