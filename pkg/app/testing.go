package app

import (
	"fmt"
	"reflect"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

// TestUIDescriptor represents a descriptor that describes a UI element and its
// location from its parents.
type TestUIDescriptor struct {
	// The location of the node. It is used by the TestMatch to find the
	// element to test.
	//
	// If empty, the expected UI element is compared with the root of the tree.
	//
	// Otherwise, each integer represents the index of the element to traverse,
	// from the root's children to the element to compare
	Path []int

	// The element to compare with the element targeted by Path. Compare
	// behavior varies depending on the element kind.
	//
	// Simple text elements only have their text value compared.
	//
	// HTML elements have their attribute compared and check if their event
	// handlers are set.
	//
	// Components have their exported field values compared.
	Expected UI
}

// TestPath is a helper function that returns a path to use in a
// TestUIDescriptor.
func TestPath(p ...int) []int {
	return p
}

// TestMatch looks for the element targeted by the descriptor in the given tree
// and reports whether it matches with the expected element.
//
// Eg:
//
//	tree := app.Div().Body(
//	    app.H2().Body(
//	        app.Text("foo"),
//	    ),
//	    app.P().Body(
//	        app.Text("bar"),
//	    ),
//	)
//
//	// Testing root:
//	err := app.TestMatch(tree, app.TestUIDescriptor{
//	    Path:     TestPath(),
//	    Expected: app.Div(),
//	})
//	// OK => err == nil
//
//	// Testing h2:
//	err := app.TestMatch(tree, app.TestUIDescriptor{
//	    Path:     TestPath(0),
//	    Expected: app.H3(),
//	})
//	// KO => err != nil because we ask h2 to match with h3
//
//	// Testing text from p:
//	err = app.TestMatch(tree, app.TestUIDescriptor{
//	    Path:     TestPath(1, 0),
//	    Expected: app.Text("bar"),
//	})
//	// OK => err == nil
func TestMatch(tree UI, d TestUIDescriptor) error {
	if d.Expected != nil {
		d.Expected.setSelf(d.Expected)
	}

	if len(d.Path) != 0 {
		idx := d.Path[0]

		if idx < 0 || idx >= len(tree.getChildren()) {
			// Check that the element does not exists.
			if d.Expected == nil {
				return nil
			}

			return errors.New("ui element to match is out of range").
				WithTag("name", d.Expected.name()).
				WithTag("kind", d.Expected.Kind()).
				WithTag("parent-name", tree.name()).
				WithTag("parent-kind", tree.Kind()).
				WithTag("parent-children-count", len(tree.getChildren())).
				WithTag("index", idx)
		}

		c := tree.getChildren()[idx]
		p := c.getParent()

		if p != tree {
			return errors.New("unexpected ui element parent").
				WithTag("name", d.Expected.name()).
				WithTag("kind", d.Expected.Kind()).
				WithTag("parent-name", p.name()).
				WithTag("parent-kind", p.Kind()).
				WithTag("parent-addr", fmt.Sprintf("%p", p)).
				WithTag("expected-parent-name", tree.name()).
				WithTag("expected-parent-kind", tree.Kind()).
				WithTag("expected-parent-addr", fmt.Sprintf("%p", tree))
		}

		d.Path = d.Path[1:]
		return TestMatch(c, d)
	}

	if d.Expected.name() != tree.name() || d.Expected.Kind() != tree.Kind() {
		return errors.New("the UI element is not matching the descriptor").
			WithTag("expected-name", d.Expected.name()).
			WithTag("expected-kind", d.Expected.Kind()).
			WithTag("current-name", tree.name()).
			WithTag("current-kind", tree.Kind())
	}

	switch d.Expected.Kind() {
	case SimpleText:
		return matchText(tree, d)

	case HTML:
		if err := matchHTMLElemAttrs(tree, d); err != nil {
			return err
		}
		return matchHTMLElemEventHandlers(tree, d)

	case Component:
		return matchComponent(tree, d)

	case RawHTML:
		return matchRaw(tree, d)

	default:
		return errors.New("the UI element is not matching the descriptor").
			WithTag("reason", "unavailable matching for the kind").
			WithTag("kind", d.Expected.Kind())
	}
}

func matchText(n UI, d TestUIDescriptor) error {
	a := n.(*text)
	b := d.Expected.(*text)

	if a.value != b.value {
		return errors.New("the text element is not matching the descriptor").
			WithTag("name", a.name()).
			WithTag("reason", "unexpected text value").
			WithTag("expected-value", b.value).
			WithTag("current-value", a.value)
	}
	return nil
}

func matchHTMLElemAttrs(n UI, d TestUIDescriptor) error {
	aAttrs := n.getAttributes()
	bAttrs := d.Expected.getAttributes()

	if len(aAttrs) != len(bAttrs) {
		return errors.New("the html element is not matching the descriptor").
			WithTag("name", n.name()).
			WithTag("reason", "unexpected attributes length").
			WithTag("expected-attributes-length", len(bAttrs)).
			WithTag("current-attributes-length", len(aAttrs))
	}

	for k, b := range bAttrs {
		a, exists := aAttrs[k]
		if !exists {
			return errors.New("the html element is not matching the descriptor").
				WithTag("name", n.name()).
				WithTag("reason", "an attribute is missing").
				WithTag("attribute", k)
		}

		if a != b {
			return errors.New("the html element is not matching the descriptor").
				WithTag("name", n.name()).
				WithTag("reason", "unexpected attribute value").
				WithTag("attribute", k).
				WithTag("expected-value", b).
				WithTag("current-value", a)
		}
	}

	for k := range bAttrs {
		_, exists := bAttrs[k]
		if !exists {
			return errors.New("the html element is not matching the descriptor").
				WithTag("name", n.name()).
				WithTag("reason", "an unexpected attribute is present").
				WithTag("attribute", k)
		}
	}

	return nil
}

func matchHTMLElemEventHandlers(n UI, d TestUIDescriptor) error {
	aevents := n.getEventHandlers()
	bevents := d.Expected.getEventHandlers()

	if len(aevents) != len(bevents) {
		return errors.New("the html element is not matching the descriptor").
			WithTag("name", n.name()).
			WithTag("reason", "unexpected event handlers length").
			WithTag("expected-event-handlers-length", len(bevents)).
			WithTag("current-event-handlers-length", len(aevents))
	}

	for k := range bevents {
		_, exists := aevents[k]
		if !exists {
			return errors.New("the html element is not matching the descriptor").
				WithTag("name", n.name()).
				WithTag("reason", "an event handler is missing").
				WithTag("event-handler", k)
		}
	}

	for k := range bevents {
		_, exists := aevents[k]
		if !exists {
			return errors.New("the html element is not matching the descriptor").
				WithTag("name", n.name()).
				WithTag("reason", "an unexpected event handler is present").
				WithTag("event-handler", k)
		}
	}

	return nil

}

func matchComponent(n UI, d TestUIDescriptor) error {
	aval := reflect.ValueOf(n).Elem()
	bval := reflect.ValueOf(d.Expected).Elem()

	compotype := reflect.TypeOf(Compo{})

	for i := 0; i < bval.NumField(); i++ {
		a := aval.Field(i)
		b := bval.Field(i)

		if a.Type() == compotype {
			continue
		}

		if !a.CanSet() {
			continue
		}

		if !reflect.DeepEqual(a.Interface(), b.Interface()) {
			return errors.New("the component is not matching with the descriptor").
				WithTag("name", n.name()).
				WithTag("reason", "unexpected field value").
				WithTag("field", bval.Type().Field(i).Name).
				WithTag("expected-value", b.Interface()).
				WithTag("current-value", a.Interface())
		}
	}

	return nil
}

func matchRaw(n UI, d TestUIDescriptor) error {
	a := n.(*raw)
	b := d.Expected.(*raw)

	if a.value != b.value {
		return errors.New("the raw html element is not matching with the descriptor").
			WithTag("name", n.name()).
			WithTag("reason", "unexpected value").
			WithTag("expected-value", b.value).
			WithTag("current-value", a.value)
	}

	return nil
}
