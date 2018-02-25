package json

import (
	"io"

	"github.com/segmentio/objconv"
)

// Codec for the JSON format.
var Codec = objconv.Codec{
	NewEmitter: func(w io.Writer) objconv.Emitter { return NewEmitter(w) },
	NewParser:  func(r io.Reader) objconv.Parser { return NewParser(r) },
}

// PrettyCodec for the JSON format.
var PrettyCodec = objconv.Codec{
	NewEmitter: func(w io.Writer) objconv.Emitter { return NewPrettyEmitter(w) },
	NewParser:  func(r io.Reader) objconv.Parser { return NewParser(r) },
}

func init() {
	for _, name := range [...]string{
		"application/json",
		"text/json",
		"json",
	} {
		objconv.Register(name, Codec)
	}
}
