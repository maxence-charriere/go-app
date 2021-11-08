// Package errors implements functions to manipulate errors.
//
// Errors created are taggable and wrappable.
//
//   errWithTags := errors.New("an error with tags").
//       Tag("a", 42).
// 	     Tag("b", 21)
//
//   errWithWrap := errors.New("error").
//       Tag("a", 42).
// 	     Wrap(errors.New("wrapped error"))
//
// The package mirrors https://golang.org/pkg/errors package.
package errors

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"unsafe"
)

// As is documented at https://golang.org/pkg/errors/#As.
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Is is documented at https://golang.org/pkg/errors/#Is.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// Unwrap is documented at https://golang.org/pkg/errors/#Unwrap.
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// Tag retrieves the value of the tag named by the key. If the tag exists,
// its value (which may be empty) is returned and the boolean is true. Otherwise
// the returned value will be empty and the boolean will be false.
func Tag(err error, k string) (string, bool) {
	ierr, ok := err.(Error)
	if !ok {
		return "", false
	}
	return ierr.Lookup(k)
}

// New returns an error with the given description that can be tagged.
func New(v string) Error {
	return Error{
		description: v,
	}
}

// Newf returns an error with the given formatted description that can be
// tagged.
func Newf(format string, v ...interface{}) Error {
	return New(fmt.Sprintf(format, v...))
}

// Error is an error implementation that supports tagging and wrapping.
type Error struct {
	description string
	tags        []tag
	maxKeyLen   int
	wrap        error
}

// Tag sets the named tag with the given value.
func (e Error) Tag(k string, v interface{}) Error {
	if e.tags == nil {
		e.tags = make([]tag, 0, 8)
	}

	if l := len(k); l > e.maxKeyLen {
		e.maxKeyLen = l
	}

	switch v := v.(type) {
	case string:
		e.tags = append(e.tags, tag{key: k, value: v})

	default:
		e.tags = append(e.tags, tag{key: k, value: fmt.Sprintf("%+v", v)})
	}

	return e
}

// Lookup retrieves the value of the tag named by the key. If the tag exists,
// its value (which may be empty) is returned and the boolean is true. Otherwise
// the returned value will be empty and the boolean will be false.
func (e Error) Lookup(tag string) (string, bool) {
	for _, t := range e.tags {
		if t.key == tag {
			return t.value, true
		}
	}

	if w, ok := e.wrap.(Error); ok {
		return w.Lookup(tag)
	}
	return "", false
}

// Wrap wraps the given error. Nil errors are ingnored.
func (e Error) Wrap(err error) Error {
	if err == nil {
		return e
	}

	if e.maxKeyLen < 5 {
		e.maxKeyLen = 5
	}

	if e.wrap == nil {
		e.wrap = err
		return e
	}

	description := ""
	if perr, ok := err.(Error); ok {
		description = perr.description
	} else {
		description = err.Error()
	}

	e.wrap = New(description).Wrap(e.wrap)
	return e
}

// Unwrap unwraps the given error.
func (e Error) Unwrap() error {
	return e.wrap
}

// Is reports if the target matches the error or its wrapped values.
func (e Error) Is(target error) bool {
	o, ok := target.(Error)
	if !ok {
		return false
	}

	return e.description == o.description &&
		reflect.DeepEqual(e.tags, o.tags)
}

func (e Error) Error() string {
	w := bytes.NewBuffer(make([]byte, 0, len(e.description)+len(e.tags)*(e.maxKeyLen+11)))
	e.format(w, 0)
	return bytesToString(w.Bytes())
}

func (e Error) format(w *bytes.Buffer, indent int) {
	w.WriteString(e.description)
	if e.wrap != nil || len(e.tags) != 0 {
		w.WriteByte(':')
	}

	tags := e.tags
	sort.Slice(tags, func(a, b int) bool {
		return strings.Compare(tags[a].key, tags[b].key) < 0
	})

	for _, t := range e.tags {
		k := t.key
		v := t.value

		w.WriteByte('\n')
		e.indent(w, indent+4)
		w.WriteString(k)
		w.WriteByte(':')
		e.indent(w, e.maxKeyLen-len(k)+1)
		w.WriteString(v)
	}

	if e.wrap == nil {
		return
	}

	w.WriteByte('\n')
	e.indent(w, indent+4)
	w.WriteString("error")
	w.WriteByte(':')
	e.indent(w, e.maxKeyLen-5+1)

	if err, ok := e.wrap.(Error); ok {
		err.format(w, indent+4)
		return
	}

	w.WriteString(e.wrap.Error())
}

func (e Error) indent(w *bytes.Buffer, n int) {
	for i := 0; i < n; i++ {
		w.WriteByte(' ')
	}
}

type tag struct {
	key   string
	value string
}

func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
