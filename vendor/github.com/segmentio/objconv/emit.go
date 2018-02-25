package objconv

import "time"

// The Emitter interface must be implemented by types that provide encoding
// of a specific format (like json, resp, ...).
//
// Emitters are not expected to be safe for use by multiple goroutines.
type Emitter interface {
	// EmitNil writes a nil value to the writer.
	EmitNil() error

	// EmitBool writes a boolean value to the writer.
	EmitBool(bool) error

	// EmitInt writes an integer value to the writer.
	EmitInt(v int64, bitSize int) error

	// EmitUint writes an unsigned integer value to the writer.
	EmitUint(v uint64, bitSize int) error

	// EmitFloat writes a floating point value to the writer.
	EmitFloat(v float64, bitSize int) error

	// EmitString writes a string value to the writer.
	EmitString(string) error

	// EmitBytes writes a []byte value to the writer.
	EmitBytes([]byte) error

	// EmitTime writes a time.Time value to the writer.
	EmitTime(time.Time) error

	// EmitDuration writes a time.Duration value to the writer.
	EmitDuration(time.Duration) error

	// EmitError writes an error value to the writer.
	EmitError(error) error

	// EmitArrayBegin writes the beginning of an array value to the writer.
	// The method receives the length of the array.
	EmitArrayBegin(int) error

	// EmitArrayEnd writes the end of an array value to the writer.
	EmitArrayEnd() error

	// EmitArrayNext is called after each array value except to the last one.
	EmitArrayNext() error

	// EmitMapBegin writes the beginning of a map value to the writer.
	// The method receives the length of the map.
	EmitMapBegin(int) error

	// EmitMapEnd writes the end of a map value to the writer.
	EmitMapEnd() error

	// EmitMapValue is called after each map key was written.
	EmitMapValue() error

	// EmitMapNext is called after each map value was written except the last one.
	EmitMapNext() error
}

// The PrettyEmitter interface may be implemented by emitters supporting a more
// human-friendlly format.
type PrettyEmitter interface {
	// PrettyEmitter returns a new emitter that outputs to the same writer in a
	// pretty format.
	PrettyEmitter() Emitter
}

// The textEmitter interface may be implemented by emitters of human-readable
// formats. Such emitters instruct the encoder to prefer using
// encoding.TextMarshaler over encoding.BinaryMarshaler for example.
type textEmitter interface {
	// EmitsText returns true if the emitter produces a human-readable format.
	TextEmitter() bool
}

func isTextEmitter(emitter Emitter) bool {
	e, _ := emitter.(textEmitter)
	return e != nil && e.TextEmitter()
}

type discardEmitter struct{}

func (e discardEmitter) EmitNil() error                     { return nil }
func (e discardEmitter) EmitBool(v bool) error              { return nil }
func (e discardEmitter) EmitInt(v int64, _ int) error       { return nil }
func (e discardEmitter) EmitUint(v uint64, _ int) error     { return nil }
func (e discardEmitter) EmitFloat(v float64, _ int) error   { return nil }
func (e discardEmitter) EmitString(v string) error          { return nil }
func (e discardEmitter) EmitBytes(v []byte) error           { return nil }
func (e discardEmitter) EmitTime(v time.Time) error         { return nil }
func (e discardEmitter) EmitDuration(v time.Duration) error { return nil }
func (e discardEmitter) EmitError(v error) error            { return nil }
func (e discardEmitter) EmitArrayBegin(v int) error         { return nil }
func (e discardEmitter) EmitArrayEnd() error                { return nil }
func (e discardEmitter) EmitArrayNext() error               { return nil }
func (e discardEmitter) EmitMapBegin(v int) error           { return nil }
func (e discardEmitter) EmitMapEnd() error                  { return nil }
func (e discardEmitter) EmitMapNext() error                 { return nil }
func (e discardEmitter) EmitMapValue() error                { return nil }

var (
	// Discard is a special emitter that outputs nothing and simply discards
	// the values.
	//
	// This emitter is mostly useful to benchmark the encoder, but it can also be
	// used to disable an encoder output if necessary.
	Discard Emitter = discardEmitter{}
)
