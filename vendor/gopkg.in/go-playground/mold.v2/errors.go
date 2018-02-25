package mold

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	// ErrInvalidDive describes an invalid dive tag configuration
	ErrInvalidDive = errors.New("Invalid dive tag configuration")

	// ErrUndefinedKeysTag describes an undefined keys tag when and endkeys tag defined
	ErrUndefinedKeysTag = errors.New("'" + endKeysTag + "' tag encountered without a corresponding '" + keysTag + "' tag")

	// ErrInvalidKeysTag describes a misuse of the keys tag
	ErrInvalidKeysTag = errors.New("'" + keysTag + "' tag must be immediately preceeded by the '" + diveTag + "' tag")
)

// ErrUndefinedTag defines a tag that does not exist
type ErrUndefinedTag struct {
	tag   string
	field string
}

// Error returns the UndefinedTag error text
func (e *ErrUndefinedTag) Error() string {
	return strings.TrimSpace(fmt.Sprintf("unregistered/undefined transformation '%s' found on field %s", e.tag, e.field))
}

// ErrInvalidTag defines a bad value for a tag being used
type ErrInvalidTag struct {
	tag   string
	field string
}

// Error returns the InvalidTag error text
func (e *ErrInvalidTag) Error() string {
	return fmt.Sprintf("invalid tag '%s' found on field %s", e.tag, e.field)
}

// An ErrInvalidTransformValue describes an invalid argument passed to Struct or Var.
// (The argument passed must be a non-nil pointer.)
type ErrInvalidTransformValue struct {
	typ reflect.Type
	fn  string
}

func (e *ErrInvalidTransformValue) Error() string {
	if e.typ == nil {
		return fmt.Sprintf("mold: %s(nil)", e.fn)
	}

	if e.typ.Kind() != reflect.Ptr {
		return fmt.Sprintf("mold: %s(non-pointer %s)", e.fn, e.typ.String())
	}

	return fmt.Sprintf("mold: %s(nil %s)", e.fn, e.typ.String())
}

// ErrInvalidTransformation describes an invalid argument passed to
// `Struct` or `Field`
type ErrInvalidTransformation struct {
	typ reflect.Type
}

// Error returns ErrInvalidTransformation message
func (e *ErrInvalidTransformation) Error() string {
	return "mold: (nil " + e.typ.String() + ")"
}
