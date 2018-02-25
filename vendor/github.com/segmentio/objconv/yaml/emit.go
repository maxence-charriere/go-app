package yaml

import (
	"encoding/base64"
	"io"
	"time"

	yaml "gopkg.in/yaml.v2"
)

// Emitter implements a YAML emitter that satisfies the objconv.Emitter
// interface.
type Emitter struct {
	w io.Writer
	// The stack is used to keep track of the container being built by the
	// emitter, which may be an arrayEmitter or mapEmitter.
	stack []emitter
}

func NewEmitter(w io.Writer) *Emitter {
	return &Emitter{w: w}
}

func (e *Emitter) Reset(w io.Writer) {
	e.w = w
	e.stack = e.stack[:0]
}

func (e *Emitter) EmitNil() error {
	return e.emit(nil)
}

func (e *Emitter) EmitBool(v bool) error {
	return e.emit(v)
}

func (e *Emitter) EmitInt(v int64, _ int) error {
	return e.emit(v)
}

func (e *Emitter) EmitUint(v uint64, _ int) error {
	return e.emit(v)
}

func (e *Emitter) EmitFloat(v float64, _ int) error {
	return e.emit(v)
}

func (e *Emitter) EmitString(v string) error {
	return e.emit(v)
}

func (e *Emitter) EmitBytes(v []byte) error {
	return e.emit(base64.StdEncoding.EncodeToString(v))
}

func (e *Emitter) EmitTime(v time.Time) error {
	return e.emit(v.Format(time.RFC3339Nano))
}

func (e *Emitter) EmitDuration(v time.Duration) error {
	return e.emit(v.String())
}

func (e *Emitter) EmitError(v error) error {
	return e.emit(v.Error())
}

func (e *Emitter) EmitArrayBegin(_ int) (err error) {
	e.push(&arrayEmitter{})
	return
}

func (e *Emitter) EmitArrayEnd() (err error) {
	e.emit(e.pop().value())
	return
}

func (e *Emitter) EmitArrayNext() (err error) {
	return
}

func (e *Emitter) EmitMapBegin(_ int) (err error) {
	e.push(&mapEmitter{})
	return
}

func (e *Emitter) EmitMapEnd() (err error) {
	e.emit(e.pop().value())
	return
}

func (e *Emitter) EmitMapValue() (err error) {
	return
}

func (e *Emitter) EmitMapNext() (err error) {
	return
}

func (e *Emitter) TextEmitter() bool {
	return true
}

func (e *Emitter) emit(v interface{}) (err error) {
	var b []byte

	if n := len(e.stack); n != 0 {
		e.stack[n-1].emit(v)
		return
	}

	if b, err = yaml.Marshal(v); err != nil {
		return
	}

	_, err = e.w.Write(b)
	return
}

func (e *Emitter) push(v emitter) {
	e.stack = append(e.stack, v)
}

func (e *Emitter) pop() emitter {
	i := len(e.stack) - 1
	v := e.stack[i]
	e.stack = e.stack[:i]
	return v
}

type emitter interface {
	emit(interface{})
	value() interface{}
}

type arrayEmitter struct {
	self []interface{}
}

func (e *arrayEmitter) emit(v interface{}) {
	e.self = append(e.self, v)
}

func (e *arrayEmitter) value() interface{} {
	return e.self
}

type mapEmitter struct {
	self yaml.MapSlice
	val  bool
}

func (e *mapEmitter) emit(v interface{}) {
	if e.val {
		e.val = false
		e.self[len(e.self)-1].Value = v
	} else {
		e.val = true
		e.self = append(e.self, yaml.MapItem{Key: v})
	}
}

func (e *mapEmitter) value() interface{} {
	return e.self
}
