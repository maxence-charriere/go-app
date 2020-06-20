package app

import (
	"context"
	"strings"

	"github.com/maxence-charriere/go-app/v6/pkg/errors"
)

// Raw returns a ui element from the given raw value.
//
// It is not recommended to use this kind of node since there is no check on the
// raw string content.
func Raw(v string) UI {
	v = strings.TrimSpace(v)

	tag := rawRootTagName(v)
	if tag == "" {
		panic(errors.New("creating raw element failed").
			Tag("reason", "opening tag not found"))
	}

	return &raw{
		outerHTML: v,
		tag:       tag,
	}
}

type raw struct {
	jsvalue    Value
	outerHTML  string
	parentElem UI
	tag        string
}

func (r *raw) Kind() Kind {
	return RawHTML
}

func (r *raw) JSValue() Value {
	return r.jsvalue
}

func (r *raw) Mounted() bool {
	return r.jsvalue != nil
}

func (r *raw) name() string {
	return "raw." + r.tag
}

func (r *raw) self() UI {
	return r
}

func (r *raw) setSelf(UI) {
}

func (r *raw) context() context.Context {
	return nil
}

func (r *raw) attributes() map[string]string {
	return nil
}

func (r *raw) eventHandlers() map[string]eventHandler {
	return nil
}

func (r *raw) parent() UI {
	return r.parentElem
}

func (r *raw) setParent(p UI) {
	r.parentElem = p
}

func (r *raw) children() []UI {
	return nil
}

func (r *raw) mount() error {
	panic("not implemented")
}

func (r *raw) dismount() {
}

func (r *raw) update(n UI) error {
	panic("not implemented")
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
