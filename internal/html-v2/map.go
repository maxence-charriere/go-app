package html

import (
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

// Call represents a component method call.
type Call struct {
	// The component identifier.
	CompoID string

	// A dot separated string that points to a component field or method.
	FieldOrMethod string

	// The JSON value to map to a field or method's first argument.
	JSONValue string

	// A string that describes a field that may required override.
	Override string
}

// func (c *Call) Call() error {

// }

func pipeline(fieldOrMethod string) ([]string, error) {
	if len(fieldOrMethod) == 0 {
		return nil, errors.New("empty")
	}

	p := strings.Split(fieldOrMethod, ".")

	for _, e := range p {
		if len(e) == 0 {
			return nil, errors.Errorf("%s: contains an empty element", fieldOrMethod)
		}
	}

	return p, nil
}

func isExported(fieldOrMethod string) bool {
	return !unicode.IsLower(rune(fieldOrMethod[0]))
}
