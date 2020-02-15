package app

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/maxence-charriere/app/pkg/log"
)

// Raw returns a node from the given raw value.
//
// Note that it is not recommended to use this kind of node since there is no
// check on the raw string content.
func Raw(v string) UI {
	v = strings.TrimSpace(v)

	tag := rawOpenTag(v)
	if tag == "" {
		log.Error("creating raw node failed").
			T("error", "no opening tag").
			Panic()
		return nil
	}

	return &raw{
		tagName:   tag,
		outerHTML: v,
	}
}

type raw struct {
	parentNode nodeWithChildren
	jsValue    Value
	tagName    string
	outerHTML  string
}

func (r *raw) nodeType() reflect.Type {
	return reflect.TypeOf(r)
}

func (r *raw) JSValue() Value {
	return r.jsValue
}

func (r *raw) parent() nodeWithChildren {
	return r.parentNode
}

func (r *raw) setParent(p nodeWithChildren) {
	r.parentNode = p
}

func (r *raw) dismount() {
	r.jsValue = nil
}

func (r *raw) raw() string {
	return r.outerHTML
}

func (r *raw) mount() error {
	if r.jsValue != nil {
		return fmt.Errorf("node already mounted: %+v", r)
	}

	var v Value

	switch r.tagName {
	case "svg":
		v = Window().
			Get("document").
			Call("createElementNS", "http://www.w3.org/2000/svg", r.tagName)

	default:
		v = Window().Get("document").Call("createElement", r.tagName)
	}

	tmpParent := Window().Get("document").Call("createElement", "div")
	tmpParent.Call("appendChild", v)
	v.Set("outerHTML", r.outerHTML)

	r.jsValue = tmpParent.Get("firstChild")
	return nil
}

func (r *raw) html(w io.Writer) {
	r.htmlWithIndent(w, 0)
}

func (r *raw) htmlWithIndent(w io.Writer, indent int) {
	writeIndent(w, indent)
	w.Write(stob(r.outerHTML))
	w.Write(ln())
}

func rawOpenTag(raw string) string {
	raw = strings.TrimSpace(raw)

	if strings.HasPrefix(raw, "</") || !strings.HasPrefix(raw, "<") {
		return ""
	}

	end := -1
	for i := 1; i < len(raw); i++ {
		if raw[i] == ' ' ||
			raw[i] == '\t' ||
			raw[i] == '\n' ||
			raw[i] == '>' {
			end = i
			break
		}
	}
	if end <= 0 {
		return ""
	}

	return raw[1:end]
}
