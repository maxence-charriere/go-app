package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"time"
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
func Tag(err error, k string) string {
	for {
		if err, ok := err.(interface{ Tag(string) string }); ok {
			if v := err.Tag(k); v != "" {
				return v
			}
		}

		if err = Unwrap(err); err == nil {
			return ""
		}
	}
}

// Message returns the error message.
func Message(err error) string {
	if err, ok := err.(Error); ok {
		return err.Message()
	}

	return err.Error()
}

// Error is the interface that describes an enriched error.
type Error interface {
	error

	// Returns the file line where the error was created.
	Line() string

	// Returns the error message.
	Message() string

	// Sets the given type to the error.
	WithType(v string) Error

	// Returns the type of the error.
	Type() string

	// Sets the tag key with the given value. The value is converted to a
	// string.
	WithTag(k string, v any) Error

	// Return the tag value associated with the given key.
	Tag(k string) string

	// Returns the tags as a list of key-value pairs.
	Tags() map[string]string

	// Wraps the given error.
	Wrap(err error) Error

	// Returns the wrapped error. Returns nil when there is no wrapped error.
	Unwrap() error
}

// New returns an error with the given message that can be enriched with a type
// and tags.
func New(msg string) Error {
	return makeRichError(msg)
}

// Newf returns an error with the given formatted message that can be enriched
// with a type and tags.
func Newf(msgFormat string, v ...any) Error {
	return makeRichError(fmt.Sprintf(msgFormat, v...))
}

type richError struct {
	line        string
	message     string
	definedType string
	tags        map[string]string
	wrappedErr  error
}

func makeRichError(msg string) richError {
	_, filename, line, _ := runtime.Caller(2)

	return richError{
		message: msg,
		line:    fmt.Sprintf("%s:%v", filepath.Base(filename), line),
	}
}

func (e richError) Line() string {
	return e.line
}

func (e richError) Message() string {
	return e.message
}

func (e richError) WithType(v string) Error {
	e.definedType = v
	return e
}

func (e richError) Type() string {
	if e.definedType != "" {
		return e.definedType
	}

	if e.wrappedErr != nil {
		return Type(e.wrappedErr)
	}

	return reflect.TypeOf(richError{}).String()
}

func (e richError) WithTag(k string, v any) Error {
	if e.tags == nil {
		e.tags = make(map[string]string)
	}

	e.tags[k] = toString(v)
	return e
}

func (e richError) Tag(k string) string {
	return e.tags[k]
}

func (e richError) Tags() map[string]string {
	return e.tags
}

func (e richError) Wrap(err error) Error {
	e.wrappedErr = err
	return e
}

func (e richError) Unwrap() error {
	return e.wrappedErr
}

func (e richError) Error() string {
	b, _ := e.MarshalJSON()
	return string(b)
}

func (e richError) MarshalJSON() ([]byte, error) {
	werr := e.wrappedErr
	if _, ok := werr.(Error); !ok && werr != nil {
		werr = richError{
			message:     werr.Error(),
			definedType: Type(werr),
		}
	}

	return Encoder(struct {
		Line    string            `json:"line,omitempty"`
		Message string            `json:"message"`
		Type    string            `json:"type"`
		Tags    map[string]string `json:"tags,omitempty"`
		Wrap    error             `json:"wrap,omitempty"`
	}{
		Line:    e.line,
		Message: e.message,
		Type:    e.Type(),
		Tags:    e.tags,
		Wrap:    werr,
	})
}

func (e richError) Is(err error) bool {
	rerr, ok := err.(richError)
	if !ok {
		return false
	}

	return rerr.line == e.line &&
		rerr.message == e.message &&
		rerr.definedType == e.definedType &&
		reflect.DeepEqual(rerr.tags, e.tags) &&
		rerr.wrappedErr == e.wrappedErr
}

func toString(v any) string {
	switch v := v.(type) {
	case string:
		return v

	case int:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)

	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)

	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)

	case bool:
		return strconv.FormatBool(v)

	case time.Duration:
		return v.String()

	case []byte:
		return string(v)

	default:
		b, _ := Encoder(v)
		return string(b)
	}
}
