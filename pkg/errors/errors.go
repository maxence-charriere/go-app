package errors

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
)

func init() {
	SetInlineEncoder()
}

// Unwrap returns the result of calling the Unwrap method on err, if err's type
// contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// Is reports whether any error in err's chain matches target.
//
// The chain consists of err itself followed by the sequence of errors obtained
// by repeatedly calling Unwrap.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
//
// An error type might provide an Is method so it can be treated as equivalent
// to an existing error. For example, if MyError defines
//
//	func (m MyError) Is(target error) bool { return target == fs.ErrExist }
//
// then Is(MyError{}, fs.ErrExist) returns true. See syscall.Errno.Is for an
// example in the standard library. An Is method should only shallowly compare
// err and the target and not call Unwrap on either.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target, and if one is
// found, sets target to that error value and returns true. Otherwise, it
// returns false.
//
// The chain consists of err itself followed by the sequence of errors obtained
// by repeatedly calling Unwrap.
//
// An error matches target if the error's concrete value is assignable to the
// value pointed to by target, or if the error has a method As(any) bool
// such that As(target) returns true. In the latter case, the As method is
// responsible for setting target.
//
// An error type might provide an As method so it can be treated as if it were a
// different error type.
//
// As panics if target is not a non-nil pointer to either a type that implements
// error, or to any interface type.
func As(err error, target any) bool {
	return errors.As(err, target)
}

// Type returns the type of the error.
//
// Uses reflect.TypeOf() when the given error does not implements the Error
// interface.
func Type(err error) string {
	if err == nil {
		return ""
	}

	if err, ok := err.(interface{ Type() string }); ok {
		return err.Type()
	}

	return reflect.TypeOf(err).String()
}

// HasType reports whether any error in err's chain matches the given type.
//
// The chain consists of err itself followed by the sequence of errors obtained
// by repeatedly calling Unwrap.
//
// An error matches the given type if the error has a method Type() string such
// that Type() returns a string equal to the given type.
func HasType(err error, v string) bool {
	for {
		if v == Type(err) {
			return true
		}

		if err = Unwrap(err); err == nil {
			return false
		}
	}
}

// UIError returns the UI message attached to err if present in the error chain.
// Otherwise it returns the default error string.
func UIError(err error) string {
	if err == nil {
		return ""
	}

	for current := err; current != nil; current = Unwrap(current) {
		if err, ok := current.(interface{ UIError() string }); ok {
			return err.UIError()
		}
	}
	return err.Error()
}

// Tag returns the first tag value in err's chain that matches the given key.
//
// The chain consists of err itself followed by the sequence of errors obtained
// by repeatedly calling Unwrap.
//
// An error has a tag when it has a method Tag(string) any such that Tag(k)
// returns a non-nil value.
func Tag(err error, k string) any {
	for {
		if err, ok := err.(interface{ Tag(string) any }); ok {
			if v := err.Tag(k); v != nil {
				return v
			}
		}

		if err = Unwrap(err); err == nil {
			return nil
		}
	}
}

// An enriched error.
type Error struct {
	Line        string
	Message     string
	DefinedType string
	UIMessage   string
	Tags        map[string]any
	WrappedErr  error
}

// New returns an error with the given message that can be enriched with a type
// and tags.
func New(msg string) Error {
	return makeError(msg)
}

// Newf returns an error with the given formatted message that can be enriched
// with a type and tags.
func Newf(msgFormat string, v ...any) Error {
	return makeError(fmt.Sprintf(msgFormat, v...))
}

func makeError(v string) Error {
	_, filename, line, _ := runtime.Caller(2)

	err := Error{
		Line:    filepath.Base(filename) + ":" + strconv.Itoa(line),
		Message: v,
	}
	return err
}

// WithType sets the application-defined type of the error.
func (e Error) WithType(v string) Error {
	e.DefinedType = v
	return e
}

// Type returns the application-defined type of the error.
//
// When no explicit type is set, the wrapped error type is returned if present.
// Otherwise, it returns the Go type of Error.
func (e Error) Type() string {
	if e.DefinedType != "" {
		return e.DefinedType
	}

	if e.WrappedErr != nil {
		return Type(e.WrappedErr)
	}

	return reflect.TypeOf(e).String()
}

// WithTag sets the named tag with the given value.
func (e Error) WithTag(k string, v any) Error {
	if e.Tags == nil {
		e.Tags = make(map[string]any)
	}

	e.Tags[k] = v
	return e
}

// Tag returns the value associated with the given tag key.
func (e Error) Tag(k string) any {
	return e.Tags[k]
}

// WithUIError sets the message intended to be displayed in the UI.
func (e Error) WithUIError(msg string) Error {
	e.UIMessage = msg
	return e
}

// UIError returns the message intended to be displayed in the UI.
func (e Error) UIError() string {
	if e.UIMessage == "" {
		return e.Error()
	}
	return e.UIMessage
}

// Wrap sets the wrapped error.
func (e Error) Wrap(err error) Error {
	e.WrappedErr = err
	return e
}

// Unwrap returns the wrapped error.
func (e Error) Unwrap() error {
	return e.WrappedErr
}

// Error returns the string representation of the error.
func (e Error) Error() string {
	s, err := getEncoder()(makeJSONError(e))
	if err != nil {
		return fmt.Sprintf(`{"message": "encoding error failed: %s"}`, err)
	}
	return string(s)
}

// MarshalJSON returns the JSON representation of the error.
func (e Error) MarshalJSON() ([]byte, error) {
	return getEncoder()(makeJSONError(e))
}

// Is reports whether err matches the receiver.
func (e Error) Is(err error) bool {
	rerr, ok := err.(Error)
	if !ok {
		return false
	}

	return rerr.Line == e.Line &&
		rerr.Message == e.Message &&
		rerr.DefinedType == e.DefinedType &&
		reflect.DeepEqual(rerr.Tags, e.Tags) &&
		isSameErr(rerr.WrappedErr, e.WrappedErr)
}

func isSameErr(a, b error) bool {
	if a == nil || b == nil {
		return a == b
	}

	ta := reflect.TypeOf(a)
	tb := reflect.TypeOf(b)
	if ta != tb || !ta.Comparable() {
		return false
	}

	return a == b
}

type jsonError struct {
	Line        string         `json:"line,omitempty"`
	Message     string         `json:"message"`
	UIMessage   string         `json:"ui,omitempty"`
	DefinedType string         `json:"type,omitempty"`
	Tags        map[string]any `json:"tags,omitempty"`
	WrappedErr  any            `json:"wrap,omitempty"`
}

func makeJSONError(err Error) jsonError {
	tags := makeJSONTags(err.Tags)

	return jsonError{
		Line:        err.Line,
		Message:     err.Message,
		UIMessage:   err.UIMessage,
		DefinedType: err.DefinedType,
		Tags:        tags,
		WrappedErr:  makeJSONWrappedErr(err.WrappedErr),
	}
}

func makeJSONWrappedErr(err error) any {
	switch err := err.(type) {
	case nil:
		return nil

	case Error:
		return makeJSONError(err)

	case *Error:
		if err == nil {
			return nil
		}
		return makeJSONError(*err)

	default:
		return err.Error()
	}
}

func makeJSONTags(tags map[string]any) map[string]any {
	if len(tags) == 0 {
		return nil
	}

	requiresNormalization := false
normalizationLoop:
	for _, v := range tags {
		switch v.(type) {
		case reflect.Type:
			requiresNormalization = true
			break normalizationLoop
		}
	}
	if !requiresNormalization {
		return tags
	}

	jsonTags := make(map[string]any, len(tags))
	for k, v := range tags {
		switch v := v.(type) {
		case reflect.Type:
			jsonTags[k] = v.String()

		default:
			jsonTags[k] = v
		}
	}
	return jsonTags
}
