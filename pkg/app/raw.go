package app

import (
	"strings"
)

// Raw creates a UI element from the given raw HTML string. The raw HTML must
// have a single root element. If the root tag cannot be determined, it defaults
// to a div element.
//
// Using Raw can be risky since there's no validation of the provided
// string content. Ensure that the content is safe and sanitized before use.
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
	jsElement     Value
	parentElement UI
	treeDepth     uint
	tag           string
	value         string
}

func (r *raw) JSValue() Value {
	return r.jsElement
}

func (r *raw) Mounted() bool {
	return r.jsElement != nil
}

func (r *raw) depth() uint {
	return r.treeDepth
}

func (r *raw) parent() UI {
	return r.parentElement
}

func (r *raw) setParent(p UI) UI {
	r.parentElement = p
	return r
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
