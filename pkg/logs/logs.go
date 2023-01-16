// Package logs implements functions to manipulate logs.
//
// Logs created are taggable.
//
//	logs.WithTags := logs.New("a log with tags").
//	    WithTag("a", 42).
//	    WithTag("b", 21)
package logs

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

var (
	// The function used to encode log entries and their tags.
	Encoder func(any) ([]byte, error)
)

func init() {
	SetInlineEncoder()
}

// SetInlineEncoder is a helper function that set the logs encoder to
// json.Marshal.
func SetInlineEncoder() {
	Encoder = json.Marshal
}

// SetIndentEncoder is a helper function that set the logs encoder to a
// function that uses json.MarshalIndent.
func SetIndentEncoder() {
	Encoder = func(v any) ([]byte, error) {
		return json.MarshalIndent(v, "", "  ")
	}
}

// A log entry.
type Entry struct {
	Line    string         `json:"line,omitempty"`
	Message string         `json:"message"`
	Tags    map[string]any `json:"tags,omitempty"`
}

// New returns a log with the given description that can be tagged.
func New(v string) Entry {
	return makeEntry(v)
}

// Newf returns a log with the given formatted description that can be tagged.
func Newf(msgFormat string, v ...any) Entry {
	return makeEntry(msgFormat, v...)
}

func makeEntry(msgFormat string, v ...any) Entry {
	_, filename, line, _ := runtime.Caller(2)

	return Entry{
		Line:    fmt.Sprintf("%s:%v", filepath.Base(filename), line),
		Message: fmt.Sprintf(msgFormat, v...),
	}
}

// WithTag sets the named tag with the given value.
func (e Entry) WithTag(k string, v any) Entry {
	if e.Tags == nil {
		e.Tags = make(map[string]any)
	}

	e.Tags[k] = v
	return e
}

func (e Entry) String() string {
	s, err := Encoder(e)
	if err != nil {
		return errors.Newf("encoding log entry failed").Wrap(err).Error()
	}
	return string(s)
}
