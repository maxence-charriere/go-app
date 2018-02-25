package objconv

import (
	"encoding"
	"errors"
	"reflect"
	"sync"
	"time"
	"unsafe"
)

// Type is an enumeration that represent all the base types supported by the
// emitters and parsers.
type Type int

const (
	Unknown Type = iota
	Nil
	Bool
	Int
	Uint
	Float
	String
	Bytes
	Time
	Duration
	Error
	Array
	Map
)

// String returns a human readable representation of the type.
func (t Type) String() string {
	switch t {
	case Nil:
		return "nil"
	case Bool:
		return "bool"
	case Int:
		return "int"
	case Uint:
		return "uint"
	case Float:
		return "float"
	case String:
		return "string"
	case Bytes:
		return "bytes"
	case Time:
		return "time"
	case Duration:
		return "duration"
	case Error:
		return "error"
	case Array:
		return "array"
	case Map:
		return "map"
	default:
		return "<type>"
	}
}

var (
	zeroCache = make(map[reflect.Type]reflect.Value)
	zeroMutex sync.RWMutex
)

// zeroValueOf and the related cache is used to keep the zero values so they
// don't need to be reallocated every time they're used.
func zeroValueOf(t reflect.Type) reflect.Value {
	zeroMutex.RLock()
	v, ok := zeroCache[t]
	zeroMutex.RUnlock()

	if !ok {
		v = reflect.Zero(t)
		zeroMutex.Lock()
		zeroCache[t] = v
		zeroMutex.Unlock()
	}

	return v
}

var (
	// basic types
	boolType           = reflect.TypeOf(false)
	intType            = reflect.TypeOf(int(0))
	int8Type           = reflect.TypeOf(int8(0))
	int16Type          = reflect.TypeOf(int16(0))
	int32Type          = reflect.TypeOf(int32(0))
	int64Type          = reflect.TypeOf(int64(0))
	uintType           = reflect.TypeOf(uint(0))
	uint8Type          = reflect.TypeOf(uint8(0))
	uint16Type         = reflect.TypeOf(uint16(0))
	uint32Type         = reflect.TypeOf(uint32(0))
	uint64Type         = reflect.TypeOf(uint64(0))
	uintptrType        = reflect.TypeOf(uintptr(0))
	float32Type        = reflect.TypeOf(float32(0))
	float64Type        = reflect.TypeOf(float64(0))
	stringType         = reflect.TypeOf("")
	bytesType          = reflect.TypeOf([]byte(nil))
	timeType           = reflect.TypeOf(time.Time{})
	durationType       = reflect.TypeOf(time.Duration(0))
	sliceInterfaceType = reflect.TypeOf(([]interface{})(nil))
	timePtrType        = reflect.PtrTo(timeType)

	// interfaces
	errorInterface             = elemTypeOf((*error)(nil))
	valueEncoderInterface      = elemTypeOf((*ValueEncoder)(nil))
	valueDecoderInterface      = elemTypeOf((*ValueDecoder)(nil))
	binaryMarshalerInterface   = elemTypeOf((*encoding.BinaryMarshaler)(nil))
	binaryUnmarshalerInterface = elemTypeOf((*encoding.BinaryUnmarshaler)(nil))
	textMarshalerInterface     = elemTypeOf((*encoding.TextMarshaler)(nil))
	textUnmarshalerInterface   = elemTypeOf((*encoding.TextUnmarshaler)(nil))
	emptyInterface             = elemTypeOf((*interface{})(nil))

	// common map types, used for optimization for map encoding algorithms
	mapStringStringType       = reflect.TypeOf((map[string]string)(nil))
	mapStringInterfaceType    = reflect.TypeOf((map[string]interface{})(nil))
	mapInterfaceInterfaceType = reflect.TypeOf((map[interface{}]interface{})(nil))
)

func elemTypeOf(v interface{}) reflect.Type {
	return reflect.TypeOf(v).Elem()
}

func stringNoCopy(b []byte) string {
	n := len(b)
	if n == 0 {
		return ""
	}
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(&b[0])),
		Len:  n,
	}))
}

// ValueParser is parser that uses "natural" in-memory representation of data
// structures.
//
// This is mainly useful for testing the decoder algorithms.
type ValueParser struct {
	stack []reflect.Value
	ctx   []valueParserContext
}

type valueParserContext struct {
	value  reflect.Value
	keys   []reflect.Value
	fields []structField
}

// NewValueParser creates a new parser that exposes the value v.
func NewValueParser(v interface{}) *ValueParser {
	return &ValueParser{
		stack: []reflect.Value{reflect.ValueOf(v)},
	}
}

func (p *ValueParser) ParseType() (Type, error) {
	v := p.value()

	if !v.IsValid() {
		return Nil, nil
	}

	switch v.Interface().(type) {
	case time.Time:
		return Time, nil

	case time.Duration:
		return Duration, nil

	case error:
		return Error, nil
	}

	switch v.Kind() {
	case reflect.Bool:
		return Bool, nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Int, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return Uint, nil

	case reflect.Float32, reflect.Float64:
		return Float, nil

	case reflect.String:
		return String, nil

	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return Bytes, nil
		}
		return Array, nil

	case reflect.Array:
		return Array, nil

	case reflect.Map:
		return Map, nil

	case reflect.Struct:
		return Map, nil

	case reflect.Interface:
		if v.IsNil() {
			return Nil, nil
		}
	}

	return Nil, errors.New("objconv: unsupported type found in value parser: " + v.Type().String())
}

func (p *ValueParser) ParseNil() (err error) {
	return
}

func (p *ValueParser) ParseBool() (v bool, err error) {
	v = p.value().Bool()
	return
}

func (p *ValueParser) ParseInt() (v int64, err error) {
	v = p.value().Int()
	return
}

func (p *ValueParser) ParseUint() (v uint64, err error) {
	v = p.value().Uint()
	return
}

func (p *ValueParser) ParseFloat() (v float64, err error) {
	v = p.value().Float()
	return
}

func (p *ValueParser) ParseString() (v []byte, err error) {
	v = []byte(p.value().String())
	return
}

func (p *ValueParser) ParseBytes() (v []byte, err error) {
	v = p.value().Bytes()
	return
}

func (p *ValueParser) ParseTime() (v time.Time, err error) {
	v = p.value().Interface().(time.Time)
	return
}

func (p *ValueParser) ParseDuration() (v time.Duration, err error) {
	v = p.value().Interface().(time.Duration)
	return
}

func (p *ValueParser) ParseError() (v error, err error) {
	v = p.value().Interface().(error)
	return
}

func (p *ValueParser) ParseArrayBegin() (n int, err error) {
	v := p.value()
	n = v.Len()
	p.pushContext(valueParserContext{value: v})

	if n != 0 {
		p.push(v.Index(0))
	}

	return
}

func (p *ValueParser) ParseArrayEnd(n int) (err error) {
	if n != 0 {
		p.pop()
	}
	p.popContext()
	return
}

func (p *ValueParser) ParseArrayNext(n int) (err error) {
	ctx := p.context()
	p.pop()
	p.push(ctx.value.Index(n))
	return
}

func (p *ValueParser) ParseMapBegin() (n int, err error) {
	v := p.value()

	if v.Kind() == reflect.Map {
		n = v.Len()
		k := v.MapKeys()
		p.pushContext(valueParserContext{value: v, keys: k})
		if n != 0 {
			p.push(k[0])
		}
	} else {
		c := valueParserContext{value: v}
		s := structCache.lookup(v.Type())

		for _, f := range s.fields {
			if !f.omit(v.FieldByIndex(f.index)) {
				c.fields = append(c.fields, f)
				n++
			}
		}

		p.pushContext(c)
		if n != 0 {
			p.push(reflect.ValueOf(c.fields[0].name))
		}
	}

	return
}

func (p *ValueParser) ParseMapEnd(n int) (err error) {
	if n != 0 {
		p.pop()
	}
	p.popContext()
	return
}

func (p *ValueParser) ParseMapValue(n int) (err error) {
	ctx := p.context()
	p.pop()

	if ctx.keys != nil {
		p.push(ctx.value.MapIndex(ctx.keys[n]))
	} else {
		p.push(ctx.value.FieldByIndex(ctx.fields[n].index))
	}

	return
}

func (p *ValueParser) ParseMapNext(n int) (err error) {
	ctx := p.context()
	p.pop()

	if ctx.keys != nil {
		p.push(ctx.keys[n])
	} else {
		p.push(reflect.ValueOf(ctx.fields[n].name))
	}

	return
}

func (p *ValueParser) value() reflect.Value {
	v := p.stack[len(p.stack)-1]

	if !v.IsValid() {
		return v
	}

	switch v.Interface().(type) {
	case error:
		return v
	}

dereference:
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr:
		if !v.IsNil() {
			v = v.Elem()
			goto dereference
		}
	}

	return v
}

func (p *ValueParser) push(v reflect.Value) {
	p.stack = append(p.stack, v)
}

func (p *ValueParser) pop() {
	p.stack = p.stack[:len(p.stack)-1]
}

func (p *ValueParser) pushContext(ctx valueParserContext) {
	p.ctx = append(p.ctx, ctx)
}

func (p *ValueParser) popContext() {
	p.ctx = p.ctx[:len(p.ctx)-1]
}

func (p *ValueParser) context() *valueParserContext {
	return &p.ctx[len(p.ctx)-1]
}

// ValueEmitter is a special kind of emitter, instead of serializing the values
// it receives it builds an in-memory representation of the data.
//
// This is useful for testing the high-level API of the package without actually
// having to generate a serialized representation.
type ValueEmitter struct {
	stack []interface{}
	marks []int
}

// NewValueEmitter returns a pointer to a new ValueEmitter object.
func NewValueEmitter() *ValueEmitter {
	return &ValueEmitter{}
}

// Value returns the value built in the emitter.
func (e *ValueEmitter) Value() interface{} { return e.stack[0] }

func (e *ValueEmitter) EmitNil() error { return e.push(nil) }

func (e *ValueEmitter) EmitBool(v bool) error { return e.push(v) }

func (e *ValueEmitter) EmitInt(v int64, _ int) error { return e.push(v) }

func (e *ValueEmitter) EmitUint(v uint64, _ int) error { return e.push(v) }

func (e *ValueEmitter) EmitFloat(v float64, _ int) error { return e.push(v) }

func (e *ValueEmitter) EmitString(v string) error { return e.push(v) }

func (e *ValueEmitter) EmitBytes(v []byte) error { return e.push(v) }

func (e *ValueEmitter) EmitTime(v time.Time) error { return e.push(v) }

func (e *ValueEmitter) EmitDuration(v time.Duration) error { return e.push(v) }

func (e *ValueEmitter) EmitError(v error) error { return e.push(v) }

func (e *ValueEmitter) EmitArrayBegin(v int) error { return e.pushMark() }

func (e *ValueEmitter) EmitArrayEnd() error {
	v := e.pop(e.popMark())
	a := make([]interface{}, len(v))
	copy(a, v)
	return e.push(a)
}

func (e *ValueEmitter) EmitArrayNext() error { return nil }

func (e *ValueEmitter) EmitMapBegin(v int) error { return e.pushMark() }

func (e *ValueEmitter) EmitMapEnd() error {
	v := e.pop(e.popMark())
	n := len(v)
	m := make(map[interface{}]interface{}, n/2)

	for i := 0; i != n; i += 2 {
		m[v[i]] = v[i+1]
	}

	return e.push(m)
}

func (e *ValueEmitter) EmitMapValue() error { return nil }

func (e *ValueEmitter) EmitMapNext() error { return nil }

func (e *ValueEmitter) push(v interface{}) error {
	e.stack = append(e.stack, v)
	return nil
}

func (e *ValueEmitter) pop(n int) []interface{} {
	v := e.stack[n:]
	e.stack = e.stack[:n]
	return v
}

func (e *ValueEmitter) pushMark() error {
	e.marks = append(e.marks, len(e.stack))
	return nil
}

func (e *ValueEmitter) popMark() int {
	n := len(e.marks) - 1
	m := e.marks[n]
	e.marks = e.marks[:n]
	return m
}
