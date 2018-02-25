package objutil

import "strings"

// Tag represents the result of parsing the tag of a struct field.
type Tag struct {
	// Name is the field name that should be used when serializing.
	Name string

	// Omitempty is true if the tag had `omitempty` set.
	Omitempty bool

	// Omitzero is true if the tag had `omitzero` set.
	Omitzero bool
}

// ParseTag parses a raw tag obtained from a struct field, returning the results
// as a tag value.
func ParseTag(s string) Tag {
	var name string
	var omitzero bool
	var omitempty bool

	name, s = parseNextTagToken(s)

	for len(s) != 0 {
		var token string
		switch token, s = parseNextTagToken(s); token {
		case "omitempty":
			omitempty = true
		case "omitzero":
			omitzero = true
		}
	}

	return Tag{
		Name:      name,
		Omitempty: omitempty,
		Omitzero:  omitzero,
	}
}

// ParseTagJSON is similar to ParseTag but only supports features supported by
// the standard encoding/json package.
func ParseTagJSON(s string) Tag {
	var name string
	var omitempty bool

	name, s = parseNextTagToken(s)

	for len(s) != 0 {
		var token string
		switch token, s = parseNextTagToken(s); token {
		case "omitempty":
			omitempty = true
		}
	}

	return Tag{
		Name:      name,
		Omitempty: omitempty,
	}
}

func parseNextTagToken(s string) (token string, next string) {
	if split := strings.IndexByte(s, ','); split < 0 {
		token = s
	} else {
		token, next = s[:split], s[split+1:]
	}
	return
}
