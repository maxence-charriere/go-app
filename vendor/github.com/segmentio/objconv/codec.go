package objconv

import (
	"io"
	"sync"
)

// A Codec is a factory for encoder and decoders that work on byte streams.
type Codec struct {
	NewEmitter func(io.Writer) Emitter
	NewParser  func(io.Reader) Parser
}

// NewEncoder returns a new encoder that outputs to w.
func (c Codec) NewEncoder(w io.Writer) *Encoder {
	return NewEncoder(c.NewEmitter(w))
}

// NewDecoder returns a new decoder that takes input from r.
func (c Codec) NewDecoder(r io.Reader) *Decoder {
	return NewDecoder(c.NewParser(r))
}

// NewStreamEncoder returns a new stream encoder that outputs to w.
func (c Codec) NewStreamEncoder(w io.Writer) *StreamEncoder {
	return NewStreamEncoder(c.NewEmitter(w))
}

// NewStreamDecoder returns a new stream decoder that takes input from r.
func (c Codec) NewStreamDecoder(r io.Reader) *StreamDecoder {
	return NewStreamDecoder(c.NewParser(r))
}

// A Registry associates mime types to codecs.
//
// It is safe to use a registry concurrently from multiple goroutines.
type Registry struct {
	mutex  sync.RWMutex
	codecs map[string]Codec
}

// Register adds a codec for a mimetype to r.
func (reg *Registry) Register(mimetype string, codec Codec) {
	defer reg.mutex.Unlock()
	reg.mutex.Lock()

	if reg.codecs == nil {
		reg.codecs = make(map[string]Codec)
	}

	reg.codecs[mimetype] = codec
}

// Unregister removes the codec for a mimetype from r.
func (reg *Registry) Unregister(mimetype string) {
	defer reg.mutex.Unlock()
	reg.mutex.Lock()

	delete(reg.codecs, mimetype)
}

// Lookup returns the codec associated with mimetype, ok is set to true or false
// based on whether a codec was found.
func (reg *Registry) Lookup(mimetype string) (codec Codec, ok bool) {
	reg.mutex.RLock()
	codec, ok = reg.codecs[mimetype]
	reg.mutex.RUnlock()
	return
}

// Codecs returns a map of all codecs registered in reg.
func (reg *Registry) Codecs() (codecs map[string]Codec) {
	codecs = make(map[string]Codec)
	reg.mutex.RLock()
	for mimetype, codec := range reg.codecs {
		codecs[mimetype] = codec
	}
	reg.mutex.RUnlock()
	return
}

// The global registry to which packages add their codecs.
var registry Registry

// Register adds a codec for a mimetype to the global registry.
func Register(mimetype string, codec Codec) {
	registry.Register(mimetype, codec)
}

// Unregister removes the codec for a mimetype from the global registry.
func Unregister(mimetype string) {
	registry.Unregister(mimetype)
}

// Lookup returns the codec associated with mimetype, ok is set to true or false
// based on whether a codec was found.
func Lookup(mimetype string) (Codec, bool) {
	return registry.Lookup(mimetype)
}

// Codecs returns a map of all codecs registered in the global registry.
func Codecs() map[string]Codec {
	return registry.Codecs()
}
