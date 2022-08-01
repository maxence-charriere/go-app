package app

import (
	"context"
	"html"
	"io"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

// Text creates a simple text element.
func Text(v any) UI {
	return &text{value: toString(v)}
}

type text struct {
	disp       Dispatcher
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
	return t.jsvalue != nil && t.getDispatcher() != nil
}

func (t *text) name() string {
	return "text"
}

func (t *text) self() UI {
	return t
}

func (t *text) setSelf(n UI) {
}

func (t *text) getContext() context.Context {
	return context.TODO()
}

func (t *text) getDispatcher() Dispatcher {
	return t.disp
}

func (t *text) getAttributes() attributes {
	return nil
}

func (t *text) getEventHandlers() eventHandlers {
	return nil
}

func (t *text) getParent() UI {
	return t.parentElem
}

func (t *text) setParent(p UI) {
	t.parentElem = p
}

func (t *text) getChildren() []UI {
	return nil
}

func (t *text) mount(d Dispatcher) error {
	if t.Mounted() {
		return errors.New("mounting ui element failed").
			Tag("reason", "already mounted").
			Tag("kind", t.Kind()).
			Tag("name", t.name()).
			Tag("value", t.value)
	}

	t.disp = d
	t.jsvalue = Window().createTextNode(t.value)
	return nil
}

func (t *text) dismount() {
	t.jsvalue = nil
}

func (t *text) canUpdateWith(n UI) bool {
	_, ok := n.(*text)
	return ok
}

func (t *text) updateWith(n UI) error {
	if !t.Mounted() {
		return nil
	}

	o, _ := n.(*text)
	if t.value != o.value {
		t.value = o.value
		t.JSValue().setNodeValue(o.value)
	}

	return nil
}

func (t *text) preRender(Page) {
}

func (t *text) onComponentEvent(any) {
}

func (t *text) html(w io.Writer) {
	w.Write([]byte(html.EscapeString(t.value)))
}

func (t *text) htmlWithIndent(w io.Writer, indent int) {
	writeIndent(w, indent)
	w.Write([]byte(html.EscapeString(t.value)))
}
