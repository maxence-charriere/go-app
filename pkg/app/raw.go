package app

import (
	"context"
	"io"
	"net/url"
	"strings"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

// Raw returns a ui element from the given raw value. HTML raw value must have a
// single root.
//
// It is not recommended to use this kind of node since there is no check on the
// raw string content.
func Raw(v string) UI {
	v = strings.TrimSpace(v)

	tag := rawRootTagName(v)
	if tag == "" {
		v = "<div></div>"
	}

	return &raw{
		value: v,
		tag:   tag,
	}
}

type raw struct {
	disp       Dispatcher
	jsvalue    Value
	parentElem UI
	tag        string
	value      string
}

func (r *raw) kind() Kind {
	return RawHTML
}

func (r *raw) JSValue() Value {
	return r.jsvalue
}

func (r *raw) IsMounted() bool {
	return r.jsvalue != nil && r.getDispatcher() != nil
}

func (r *raw) name() string {
	return "raw." + r.tag
}

func (r *raw) self() UI {
	return r
}

func (r *raw) setSelf(UI) {
}

func (r *raw) getContext() context.Context {
	return nil
}

func (r *raw) getDispatcher() Dispatcher {
	return r.disp
}

func (r *raw) getAttributes() map[string]string {
	return nil
}

func (r *raw) getEventHandlers() map[string]eventHandler {
	return nil
}

func (r *raw) getParent() UI {
	return r.parentElem
}

func (r *raw) setParent(p UI) {
	r.parentElem = p
}

func (r *raw) getChildren() []UI {
	return nil
}

func (r *raw) mount(d Dispatcher) error {
	if r.IsMounted() {
		return errors.New("mounting raw html element failed").
			Tag("reason", "already mounted").
			Tag("name", r.name()).
			Tag("kind", r.kind())
	}

	r.disp = d

	wrapper, err := Window().createElement("div", "")
	if err != nil {
		return errors.New("creating raw node wrapper failed").Wrap(err)
	}

	if IsServer {
		r.jsvalue = wrapper
		return nil
	}

	wrapper.setInnerHTML(r.value)
	value := wrapper.firstChild()
	if !value.Truthy() {
		return errors.New("mounting raw html element failed").
			Tag("reason", "converting raw html to html elements returned nil").
			Tag("name", r.name()).
			Tag("kind", r.kind()).
			Tag("raw-html", r.value)
	}
	wrapper.removeChild(value)
	r.jsvalue = value
	return nil
}

func (r *raw) dismount() {
	r.jsvalue = nil
}

func (r *raw) canUpdateWith(n UI) bool {
	if n, ok := n.(*raw); ok {
		return r.value == n.value
	}
	return false
}

func (r *raw) updateWith(n UI) error {
	return nil
}

func (r *raw) onNav(*url.URL) {
}

func (r *raw) onAppUpdate() {
}

func (r *raw) onAppInstallChange() {
}

func (r *raw) onResize() {
}

func (r *raw) preRender(Page) {
}

func (r *raw) html(w io.Writer) {
	w.Write([]byte(r.value))
}

func (r *raw) htmlWithIndent(w io.Writer, indent int) {
	writeIndent(w, indent)
	w.Write([]byte(r.value))
}

func rawRootTagName(raw string) string {
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
