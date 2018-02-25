package yaml

import (
	"io"

	"github.com/segmentio/objconv"
)

// Codec for the YAML format.
var Codec = objconv.Codec{
	NewEmitter: func(w io.Writer) objconv.Emitter { return NewEmitter(w) },
	NewParser:  func(r io.Reader) objconv.Parser { return NewParser(r) },
}

func init() {
	for _, name := range [...]string{
		"application/yaml",
		"text/yaml",
		"yaml",
	} {
		objconv.Register(name, Codec)
	}
}
