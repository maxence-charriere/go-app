package errors

import (
	"encoding/json"
	"sync/atomic"
)

type encoderFunc func(any) ([]byte, error)

var encoder atomic.Value

func init() {
	SetInlineEncoder()
}

// SetEncoder sets the function used to encode errors and their tags.
//
// It is intended to be configured once during program startup, before the
// package is used concurrently. It should not be changed after concurrent use
// begins.
func SetEncoder(fn func(any) ([]byte, error)) {
	if fn == nil {
		panic("errors: nil encoder")
	}

	encoder.Store(encoderFunc(fn))
}

// SetInlineEncoder is a helper function that set the error encoder to
// json.Marshal.
func SetInlineEncoder() {
	SetEncoder(json.Marshal)
}

// SetIndentEncoder is a helper function that set the error encoder to a
// function that uses json.MarshalIndent.
func SetIndentEncoder() {
	SetEncoder(func(v any) ([]byte, error) {
		return json.MarshalIndent(v, "", "  ")
	})
}

func getEncoder() encoderFunc {
	return encoder.Load().(encoderFunc)
}
