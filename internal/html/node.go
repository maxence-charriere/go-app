package html

import (
	"github.com/murlokswarm/app"
)

type node interface {
	app.DOMNode
	SetParent(node)
	ConsumeChanges() []Change
}

func attrsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}

	for k, va := range a {
		if vb, ok := b[k]; !ok || va != vb {
			return false
		}
	}
	return true
}
