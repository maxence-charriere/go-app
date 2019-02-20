package app

import (
	"fmt"
	"strings"
)

// attrTransform perform transformation for a given attribute.
type attrTransform func(name, value string) (string, string)

// jsToGoHandler convert a javascript handler to a go component handler.
func jsToGoHandler(name, value string) (string, string) {
	if !strings.HasPrefix(name, "on") {
		return name, value
	}

	if strings.HasPrefix(value, "js:") {
		return name, strings.TrimPrefix(value, "js:")
	}

	return name, fmt.Sprintf("callCompoHandler(this, event, '%s')", value)
}
