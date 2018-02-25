package yaml

import (
	"bytes"
	"io"
	"sync"

	"github.com/segmentio/objconv"
)

// NewEncoder returns a new YAML encoder that writes to w.
func NewEncoder(w io.Writer) *objconv.Encoder {
	return objconv.NewEncoder(NewEmitter(w))
}

// NewStreamEncoder returns a new YAML stream encoder that writes to w.
func NewStreamEncoder(w io.Writer) *objconv.StreamEncoder {
	return objconv.NewStreamEncoder(NewEmitter(w))
}

// Marshal writes the YAML representation of v to a byte slice returned in b.
func Marshal(v interface{}) (b []byte, err error) {
	m := marshalerPool.Get().(*marshaler)
	m.b.Truncate(0)

	if err = (objconv.Encoder{Emitter: m}).Encode(v); err == nil {
		b = make([]byte, m.b.Len())
		copy(b, m.b.Bytes())
	}

	marshalerPool.Put(m)
	return
}

var marshalerPool = sync.Pool{
	New: func() interface{} { return newMarshaler() },
}

type marshaler struct {
	Emitter
	b bytes.Buffer
}

func newMarshaler() *marshaler {
	m := &marshaler{}
	m.w = &m.b
	return m
}
