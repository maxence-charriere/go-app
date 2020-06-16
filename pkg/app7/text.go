package app

import (
	"reflect"

	"github.com/maxence-charriere/go-app/v6/pkg/errors"
)

// Text creates a simple text element.
func Text(v interface{}) UI {
	return &text{value: toString(v)}
}

type text struct {
	jsvalue    Value
	parentElem UI
	value      string
}

func (t *text) Kind() Kind {
	return SimpleText
}

func (t *text) JSValue() Value {
	return t.jsvalue
}

func (t *text) Mounted() bool {
	return t.jsvalue != nil
}

func (t *text) parent() UI {
	return t.parentElem
}

func (t *text) setParent(p UI) {
	t.parentElem = p
}

func (t *text) children() []UI {
	return nil
}

func (t *text) appendChild(UI) {
	panic("text can't have children")
}

func (t *text) removeChild(UI) {
}

func (t *text) mount() error {
	if t.Mounted() {
		return errors.New("mounting text failed").
			Tag("reason", "text is already mounted").
			Tag("value", t.value)
	}

	t.jsvalue = Window().
		Get("document").
		Call("createTextNode", t.value)

	return nil
}

func (t *text) update(n UI) error {
	if !t.Mounted() {
		return nil
	}

	o, isText := n.(*text)
	if !isText {
		return errors.New("updating text failed").
			Tag("reason", "updated node is not a text").
			Tag("current-value", t.value).
			Tag("updated-kind", n.Kind()).
			Tag("updated-type", reflect.TypeOf(n))
	}

	if t.value != o.value {
		t.value = o.value
		t.jsvalue.Set("nodeValue", o.value)
	}

	return nil
}

func (t *text) dismount() {
	t.jsvalue = nil
}
