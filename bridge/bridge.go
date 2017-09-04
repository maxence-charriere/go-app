package bridge

import "encoding/json"

// Payload is the interface that describes a payload.
type Payload interface {
	// Bytes returns the payload bytes.
	Bytes() []byte

	// String returns the payload as a string.
	String() string

	// Unmarshal parses the encoded response and stores the result in the value
	// pointed to by v.
	// If v is nil or not a pointer, Unmarshal returns an error.
	Unmarshal(v interface{}) error
}

type jsonPayload []byte

func (p jsonPayload) Bytes() []byte {
	return p
}

func (p jsonPayload) String() string {
	return string(p)
}

func (p jsonPayload) Unmarshal(v interface{}) error {
	return json.Unmarshal(p, v)
}
