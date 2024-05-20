package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
)

var (
	// The function used to encode errors and their tags.
	Encoder func(any) ([]byte, error)
)

func init() {
	SetInlineEncoder()
}

// SetInlineEncoder is a helper function that set the error encoder to
// json.Marshal.
func SetInlineEncoder() {
	Encoder = json.Marshal
}

// SetIndentEncoder is a helper function that set the error encoder to a
// function that uses json.MarshalIndent.
func SetIndentEncoder() {
	Encoder = func(v any) ([]byte, error) {
		return json.MarshalIndent(v, "", "  ")
	}
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

// Tag returns the first tag value in err's chain that matches the given key.
//
// The chain consists of err itself followed by the sequence of errors obtained
// by repeatedly calling Unwrap.
//
// An error has a tag when it has a method Tag(string) string such that Tag(k)
// returns a non-empty string value.
func Tag(err error, k string) any {
	for {
		if err, ok := err.(Error); ok {
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
	return makeError(msgFormat, v...)
}

func makeError(msgFormat string, v ...any) Error {
	_, filename, line, _ := runtime.Caller(2)

	err := Error{
		Line:    fmt.Sprintf("%s:%v", filepath.Base(filename), line),
		Message: fmt.Sprintf(msgFormat, v...),
	}
	return err
}

func (e Error) WithType(v string) Error {
	e.DefinedType = v
	return e
}

func (e Error) Type() string {
	if e.DefinedType != "" {
		return e.DefinedType
	}

	if e.WrappedErr != nil {
		return Type(e.WrappedErr)
	}

	return reflect.TypeOf(e).String()
}

func (e Error) WithTag(k string, v any) Error {
	if e.Tags == nil {
		e.Tags = make(map[string]any)
	}

	e.Tags[k] = v
	return e
}

func (e Error) Tag(k string) any {
	return e.Tags[k]
}

func (e Error) Wrap(err error) Error {
	e.WrappedErr = err
	return e
}

func (e Error) Unwrap() error {
	return e.WrappedErr
}

func (e Error) Error() string {
	s, err := Encoder(e)
	if err != nil {
		return fmt.Sprintf(`{"message": "encoding error failed: %s"}`, err)
	}
	return string(s)
}

func (e Error) MarshalJSON() ([]byte, error) {
	var wrappedErr any = e.WrappedErr
	if _, ok := e.WrappedErr.(Error); !ok && e.WrappedErr != nil {
		wrappedErr = e.WrappedErr.Error()
	}

	var tags map[string]any
	if l := len(e.Tags); l != 0 {
		tags = make(map[string]any, l)
		for k, v := range e.Tags {
			switch v := v.(type) {
			case reflect.Type:
				tags[k] = v.String()

			default:
				tags[k] = v
			}
		}
	}

	return Encoder(struct {
		Line        string         `json:"line,omitempty"`
		Message     string         `json:"message"`
		DefinedType string         `json:"type,omitempty"`
		Tags        map[string]any `json:"tags,omitempty"`
		WrappedErr  any            `json:"wrap,omitempty"`
	}{
		Line:        e.Line,
		Message:     e.Message,
		DefinedType: e.DefinedType,
		Tags:        tags,
		WrappedErr:  wrappedErr,
	})
}

func (e Error) Is(err error) bool {
	rerr, ok := err.(Error)
	if !ok {
		return false
	}

	return rerr.Line == e.Line &&
		rerr.Message == e.Message &&
		rerr.DefinedType == e.DefinedType &&
		reflect.DeepEqual(rerr.Tags, e.Tags) &&
		rerr.WrappedErr == e.WrappedErr
}
