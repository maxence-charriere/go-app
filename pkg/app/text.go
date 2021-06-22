package app

import (
	"context"
	"html"
	"io"
	"net/url"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

// Text creates a simple text element.
func Text(v interface{}) UI {
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
	return t.jsvalue != nil && t.dispatcher() != nil
}

func (t *text) name() string {
	return "text"
}

func (t *text) self() UI {
	return t
}

func (t *text) setSelf(n UI) {
}

func (t *text) context() context.Context {
	return context.TODO()
}

func (t *text) dispatcher() Dispatcher {
	return t.disp
}

func (t *text) attributes() map[string]string {
	return nil
}

func (t *text) eventHandlers() map[string]eventHandler {
	return nil
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

func (t *text) update(n UI) error {
	if !t.Mounted() {
		return nil
	}

	o, isText := n.(*text)
	if !isText {
		return errors.New("updating ui element failed").
			Tag("replace", true).
			Tag("reason", "different element types").
			Tag("current-kind", t.Kind()).
			Tag("current-name", t.name()).
			Tag("updated-kind", n.Kind()).
			Tag("updated-name", n.name())
	}

	if t.value != o.value {
		t.value = o.value
		t.JSValue().setNodeValue(o.value)
	}

	return nil
}

func (t *text) onNav(*url.URL) {
}

func (t *text) onAppUpdate() {
}

func (t *text) onAppInstallChange() {
}

func (t *text) onResize() {
}

func (t *text) preRender(Page) {
}

func (t *text) html(w io.Writer) {
	w.Write(stob(html.EscapeString(t.value)))
}

func (t *text) htmlWithIndent(w io.Writer, indent int) {
	writeIndent(w, indent)
	w.Write(stob(html.EscapeString(t.value)))
}
