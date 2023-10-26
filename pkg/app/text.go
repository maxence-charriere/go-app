package app

import (
	"fmt"
	"html"
	"io"
)

// Text creates a simple text element.
func Text(v any) UI {
	return &text{value: toString(v)}
}

// Text creates a simple text element with the given format and values.
func Textf(format string, v ...any) UI {
	return &text{value: fmt.Sprintf(format, v...)}
}

type text struct {
	disp          Dispatcher
	jsvalue       Value
	parentElement UI
	value         string
}

func (t *text) JSValue() Value {
	return t.jsvalue
}

func (t *text) Mounted() bool {
	return t.jsvalue != nil
}

func (t *text) getParent() UI {
	return t.parentElement
}

func (t *text) setParent(p UI) UI {
	t.parentElement = p
	return t
}

func (t *text) html(w io.Writer) {
	w.Write([]byte(html.EscapeString(t.value)))
}

func (t *text) htmlWithIndent(w io.Writer, indent int) {
	writeIndent(w, indent)
	w.Write([]byte(html.EscapeString(t.value)))
}

func (t *text) parent() UI {
	return t.parentElement
}
