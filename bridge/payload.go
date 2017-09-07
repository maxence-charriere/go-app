package bridge

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// Payload is the interface that describes a payload.
type Payload interface {
	// Len returns the payload size.
	Len() int

	// Bytes returns the payload bytes.
	Bytes() []byte

	// String returns the payload as a string.
	String() string

	// Unmarshal parses the encoded response and stores the result in the value
	// pointed to by v.
	// If v is nil or not a pointer, Unmarshal panics.
	Unmarshal(v interface{})
}

// NewPayload creates a payload from v.
func NewPayload(v interface{}) Payload {
	return makeJSONPayload(v)
}

// PayloadFromBytes creates a payload from an encoded byte slice.
func PayloadFromBytes(b []byte) Payload {
	return jsonPayload(b)
}

type jsonPayload []byte

func makeJSONPayload(v interface{}) (p jsonPayload) {
	var err error
	if p, err = json.Marshal(v); err != nil {
		panic(errors.Wrap(err, "making JSON payload failed"))
	}
	return
}

func (p jsonPayload) Len() int {
	return len(p)
}

func (p jsonPayload) Bytes() []byte {
	return p
}

func (p jsonPayload) String() string {
	return string(p)
}

func (p jsonPayload) Unmarshal(v interface{}) {
	if err := json.Unmarshal(p, v); err != nil {
		panic(errors.Wrap(err, "unmarshalling JSON payload failed"))
	}
}
