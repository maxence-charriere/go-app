package objconv

import "time"

// The Parser interface must be implemented by types that provide decoding of a
// specific format (like json, resp, ...).
//
// Parsers are not expected to be safe for use by multiple goroutines.
type Parser interface {
	// ParseType is called by a decoder to ask the parser what is the type of
	// the next value that can be parsed.
	//
	// ParseType must be idempotent, it must be possible to call it multiple
	// without actually changing the state of the parser.
	ParseType() (Type, error)

	// ParseNil parses a nil value.
	ParseNil() error

	// ParseBool parses a boolean value.
	ParseBool() (bool, error)

	// ParseInt parses an integer value.
	ParseInt() (int64, error)

	// ParseUint parses an unsigned integer value.
	ParseUint() (uint64, error)

	// ParseFloat parses a floating point value.
	ParseFloat() (float64, error)

	// ParseString parses a string value.
	//
	// The string is returned as a byte slice because it is expected to be
	// pointing at an internal memory buffer, the decoder will make a copy of
	// the value. This design allows more memory allocation optimizations.
	ParseString() ([]byte, error)

	// ParseBytes parses a byte array value.
	//
	// The returned byte slice is expected to be pointing at an internal memory
	// buffer, the decoder will make a copy of the value. This design allows more
	// memory allocation optimizations.
	ParseBytes() ([]byte, error)

	// ParseTime parses a time value.
	ParseTime() (time.Time, error)

	// ParseDuration parses a duration value.
	ParseDuration() (time.Duration, error)

	// ParseError parses an error value.
	ParseError() (error, error)

	// ParseArrayBegin is called by the array-decoding algorithm when it starts.
	//
	// The method should return the length of the array being decoded, or a
	// negative value if it is unknown (some formats like json don't keep track
	// of the length of the array).
	ParseArrayBegin() (int, error)

	// ParseArrayEnd is called by the array-decoding algorithm when it
	// completes.
	//
	// The method receives the iteration counter as argument, which indicates
	// how many values were decoded from the array.
	ParseArrayEnd(int) error

	// ParseArrayNext is called by the array-decoding algorithm between each
	// value parsed in the array.
	//
	// The method receives the iteration counter as argument, which indicates
	// how many values were decoded from the array.
	//
	// If the ParseArrayBegin method returned a negative value this method
	// should return objconv.End to indicated that there is no more elements to
	// parse in the array. In this case the method is also called right before
	// decoding the first element ot handle the case where the array is empty
	// and the end-of-array marker can be read right away.
	ParseArrayNext(int) error

	// ParseMapBegin is called by the map-decoding algorithm when it starts.
	//
	// The method should return the length of the map being decoded, or a
	// negative value if it is unknown (some formats like json don't keep track
	// of the length of the map).
	ParseMapBegin() (int, error)

	// ParseMapEnd is called by the map-decoding algorithm when it completes.
	//
	// The method receives the iteration counter as argument, which indicates
	// how many values were decoded from the map.
	ParseMapEnd(int) error

	// ParseMapValue is called by the map-decoding algorithm after parsing a key
	// but before parsing the associated value.
	//
	// The method receives the iteration counter as argument, which indicates
	// how many values were decoded from the map.
	ParseMapValue(int) error

	// ParseMapNext is called by the map-decoding algorithm between each
	// value parsed in the map.
	//
	// The method receives the iteration counter as argument, which indicates
	// how many values were decoded from the map.
	//
	// If the ParseMapBegin method returned a negative value this method should
	// return objconv.End to indicated that there is no more elements to parse
	// in the map. In this case the method is also called right before decoding
	// the first element ot handle the case where the array is empty and the
	// end-of-map marker can be read right away.
	ParseMapNext(int) error
}

// The bytesDecoder interface may optionnaly be implemented by a Parser to
// provide an extra step in decoding a byte slice. This is sometimes necessary
// if the associated Emitter has transformed bytes slices because the format is
// not capable of representing binary data.
type bytesDecoder interface {
	// DecodeBytes is called when the destination variable for a string or a
	// byte slice is a byte slice, allowing the parser to apply a transformation
	// before the value is stored.
	DecodeBytes([]byte) ([]byte, error)
}

// The textParser interface may be implemented by parsers of human-readable
// formats. Such parsers instruct the encoder to prefer using
// encoding.TextUnmarshaler over encoding.BinaryUnmarshaler for example.
type textParser interface {
	// EmitsText returns true if the parser produces a human-readable format.
	TextParser() bool
}

func isTextParser(parser Parser) bool {
	p, _ := parser.(textParser)
	return p != nil && p.TextParser()
}
