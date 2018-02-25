package objconv

import (
	"encoding"
	"fmt"
	"io"
	"reflect"
	"time"
	"unsafe"
)

// An Encoder implements the high-level encoding algorithm that inspect encoded
// values and drive the use of an Emitter to create a serialized representation
// of the data.
//
// Instances of Encoder are not safe for use by multiple goroutines.
type Encoder struct {
	Emitter     Emitter // the emitter used by this encoder
	SortMapKeys bool    // whether map keys should be sorted
	key         bool
}

// NewEncoder returns a new encoder that outputs values to e.
//
// Encoders created by this function use the default encoder configuration,
// which is equivalent to using a zero-value EncoderConfig with only the Emitter
// field set.
//
// The function panics if e is nil.
func NewEncoder(e Emitter) *Encoder {
	if e == nil {
		panic("objconv: the emitter is nil")
	}
	return &Encoder{Emitter: e}
}

// Encode encodes the generic value v.
func (e Encoder) Encode(v interface{}) (err error) {
	if err = e.encodeMapValueMaybe(); err != nil {
		return
	}

	// This type switch optimizes encoding of common value types, it prevents
	// the use of reflection to identify the type of the value, which saves a
	// dynamic memory allocation.
	switch x := v.(type) {
	case nil:
		return e.Emitter.EmitNil()

	case bool:
		return e.Emitter.EmitBool(x)

	case int:
		return e.Emitter.EmitInt(int64(x), 0)

	case int8:
		return e.Emitter.EmitInt(int64(x), 8)

	case int16:
		return e.Emitter.EmitInt(int64(x), 16)

	case int32:
		return e.Emitter.EmitInt(int64(x), 32)

	case int64:
		return e.Emitter.EmitInt(x, 64)

	case uint8:
		return e.Emitter.EmitUint(uint64(x), 8)

	case uint16:
		return e.Emitter.EmitUint(uint64(x), 16)

	case uint32:
		return e.Emitter.EmitUint(uint64(x), 32)

	case uint64:
		return e.Emitter.EmitUint(x, 64)

	case string:
		return e.Emitter.EmitString(x)

	case []byte:
		return e.Emitter.EmitBytes(x)

	case time.Time:
		return e.Emitter.EmitTime(x)

	case time.Duration:
		return e.Emitter.EmitDuration(x)

	case []string:
		return e.encodeSliceOfString(x)

	case []interface{}:
		return e.encodeSliceOfInterface(x)

	case map[string]string:
		return e.encodeMapStringString(x)

	case map[string]interface{}:
		return e.encodeMapStringInterface(x)

	case map[interface{}]interface{}:
		return e.encodeMapInterfaceInterface(x)

		// Also checks for pointer types so the program can use this as a way
		// to avoid the dynamic memory allocation done by runtime.convT2E for
		// converting non-pointer types to empty interfaces.
	case *bool:
		if x == nil {
			return e.Emitter.EmitNil()
		}
		return e.Emitter.EmitBool(*x)

	case *int:
		if x == nil {
			return e.Emitter.EmitNil()
		}
		return e.Emitter.EmitInt(int64(*x), int(8*unsafe.Sizeof(0)))

	case *int8:
		if x == nil {
			return e.Emitter.EmitNil()
		}
		return e.Emitter.EmitInt(int64(*x), 8)

	case *int16:
		if x == nil {
			return e.Emitter.EmitNil()
		}
		return e.Emitter.EmitInt(int64(*x), 16)

	case *int32:
		if x == nil {
			return e.Emitter.EmitNil()
		}
		return e.Emitter.EmitInt(int64(*x), 32)

	case *int64:
		if x == nil {
			return e.Emitter.EmitNil()
		}
		return e.Emitter.EmitInt(*x, 64)

	case *uint8:
		if x == nil {
			return e.Emitter.EmitNil()
		}
		return e.Emitter.EmitUint(uint64(*x), 8)

	case *uint16:
		if x == nil {
			return e.Emitter.EmitNil()
		}
		return e.Emitter.EmitUint(uint64(*x), 16)

	case *uint32:
		if x == nil {
			return e.Emitter.EmitNil()
		}
		return e.Emitter.EmitUint(uint64(*x), 32)

	case *uint64:
		if x == nil {
			return e.Emitter.EmitNil()
		}
		return e.Emitter.EmitUint(*x, 64)

	case *string:
		if x == nil {
			return e.Emitter.EmitNil()
		}
		return e.Emitter.EmitString(*x)

	case *[]byte:
		if x == nil {
			return e.Emitter.EmitNil()
		}
		return e.Emitter.EmitBytes(*x)

	case *time.Time:
		if x == nil {
			return e.Emitter.EmitNil()
		}
		return e.Emitter.EmitTime(*x)

	case *time.Duration:
		if x == nil {
			return e.Emitter.EmitNil()
		}
		return e.Emitter.EmitDuration(*x)

	case *[]string:
		if x == nil {
			return e.Emitter.EmitNil()
		}
		return e.encodeSliceOfString(*x)

	case *[]interface{}:
		if x == nil {
			return e.Emitter.EmitNil()
		}
		return e.encodeSliceOfInterface(*x)

	case *map[string]string:
		if x == nil {
			return e.Emitter.EmitNil()
		}
		return e.encodeMapStringString(*x)

	case *map[string]interface{}:
		if x == nil {
			return e.Emitter.EmitNil()
		}
		return e.encodeMapStringInterface(*x)

	case *map[interface{}]interface{}:
		if x == nil {
			return e.Emitter.EmitNil()
		}
		return e.encodeMapInterfaceInterface(*x)

	case ValueEncoder:
		return x.EncodeValue(e)

	default:
		return e.encode(reflect.ValueOf(v))
	}
}

func (e *Encoder) encodeMapValueMaybe() (err error) {
	if e.key {
		e.key, err = false, e.Emitter.EmitMapValue()
	}
	return
}

func (e Encoder) encode(v reflect.Value) error {
	return encodeFuncOf(v.Type())(e, v)
}

func (e Encoder) encodeBool(v reflect.Value) error {
	return e.Emitter.EmitBool(v.Bool())
}

func (e Encoder) encodeInt(v reflect.Value) error {
	return e.Emitter.EmitInt(v.Int(), 0)
}

func (e Encoder) encodeInt8(v reflect.Value) error {
	return e.Emitter.EmitInt(v.Int(), 8)
}

func (e Encoder) encodeInt16(v reflect.Value) error {
	return e.Emitter.EmitInt(v.Int(), 16)
}

func (e Encoder) encodeInt32(v reflect.Value) error {
	return e.Emitter.EmitInt(v.Int(), 32)
}

func (e Encoder) encodeInt64(v reflect.Value) error {
	return e.Emitter.EmitInt(v.Int(), 64)
}

func (e Encoder) encodeUint(v reflect.Value) error {
	return e.Emitter.EmitUint(v.Uint(), 0)
}

func (e Encoder) encodeUint8(v reflect.Value) error {
	return e.Emitter.EmitUint(v.Uint(), 8)
}

func (e Encoder) encodeUint16(v reflect.Value) error {
	return e.Emitter.EmitUint(v.Uint(), 16)
}

func (e Encoder) encodeUint32(v reflect.Value) error {
	return e.Emitter.EmitUint(v.Uint(), 32)
}

func (e Encoder) encodeUint64(v reflect.Value) error {
	return e.Emitter.EmitUint(v.Uint(), 64)
}

func (e Encoder) encodeUintptr(v reflect.Value) error {
	return e.Emitter.EmitUint(v.Uint(), 0)
}

func (e Encoder) encodeFloat32(v reflect.Value) error {
	return e.Emitter.EmitFloat(v.Float(), 32)
}

func (e Encoder) encodeFloat64(v reflect.Value) error {
	return e.Emitter.EmitFloat(v.Float(), 64)
}

func (e Encoder) encodeString(v reflect.Value) error {
	return e.Emitter.EmitString(v.String())
}

func (e Encoder) encodeBytes(v reflect.Value) error {
	return e.Emitter.EmitBytes(v.Bytes())
}

func (e Encoder) encodeTime(v reflect.Value) error {
	var t time.Time

	// Here we may receive either a pointer or a plain value because there is a
	// special case for *time.Time in the encoder to avoid having it match the
	// encoding.TextMarshaler interface and instead treat it the same way than
	// if we had gotten the plain value right away.
	//
	// As a side effect, this also sometimes permit more optimizations because
	// having a pointer will likely avoid a memory allocation when calling
	// Interface on the value.
	if v.Kind() != reflect.Ptr {
		t = v.Interface().(time.Time)
	} else {
		if ptr := v.Interface().(*time.Time); ptr == nil {
			return e.Emitter.EmitNil()
		} else {
			t = *ptr
		}
	}

	return e.Emitter.EmitTime(t)
}

func (e Encoder) encodeDuration(v reflect.Value) error {
	return e.Emitter.EmitDuration(time.Duration(v.Int()))
}

func (e Encoder) encodeError(v reflect.Value) error {
	return e.Emitter.EmitError(v.Interface().(error))
}

func (e Encoder) encodeArray(v reflect.Value) error {
	return e.encodeArrayWith(v, encodeFuncOf(v.Type().Elem()))
}

func (e Encoder) encodeArrayWith(v reflect.Value, f encodeFunc) error {
	i := 0
	return e.EncodeArray(v.Len(), func(e Encoder) (err error) {
		err = f(e, v.Index(i))
		i++
		return
	})
}

func (e Encoder) encodeSliceOfString(a []string) error {
	i := 0
	return e.EncodeArray(len(a), func(e Encoder) (err error) {
		err = e.Emitter.EmitString(a[i])
		i++
		return
	})
}

func (e Encoder) encodeSliceOfInterface(a []interface{}) error {
	i := 0
	return e.EncodeArray(len(a), func(e Encoder) (err error) {
		err = e.Encode(a[i])
		i++
		return
	})
}

func (e Encoder) encodeMap(v reflect.Value) error {
	t := v.Type()
	kf := encodeFuncOf(t.Key())
	vf := encodeFuncOf(t.Elem())
	return e.encodeMapWith(v, kf, vf)
}

func (e Encoder) encodeMapWith(v reflect.Value, kf encodeFunc, vf encodeFunc) error {
	t := v.Type()

	if !e.SortMapKeys {
		switch {
		case t.ConvertibleTo(mapInterfaceInterfaceType):
			return e.encodeMapInterfaceInterfaceValue(v.Convert(mapInterfaceInterfaceType))

		case t.ConvertibleTo(mapStringInterfaceType):
			return e.encodeMapStringInterfaceValue(v.Convert(mapStringInterfaceType))

		case t.ConvertibleTo(mapStringStringType):
			return e.encodeMapStringStringValue(v.Convert(mapStringStringType))
		}
	}

	var k []reflect.Value
	var n = v.Len()
	var i = 0

	if n != 0 {
		k = v.MapKeys()

		if e.SortMapKeys {
			sortValues(t.Key(), k)
		}
	}

	return e.EncodeMap(n, func(ke Encoder, ve Encoder) (err error) {
		if err = kf(e, k[i]); err != nil {
			return
		}
		if err = e.Emitter.EmitMapValue(); err != nil {
			return
		}
		if err = vf(e, v.MapIndex(k[i])); err != nil {
			return
		}
		i++
		return
	})
}

func (e Encoder) encodeMapInterfaceInterfaceValue(v reflect.Value) error {
	return e.encodeMapInterfaceInterface(v.Interface().(map[interface{}]interface{}))
}

func (e Encoder) encodeMapInterfaceInterface(m map[interface{}]interface{}) (err error) {
	n := len(m)
	i := 0

	if err = e.Emitter.EmitMapBegin(n); err != nil {
		return
	}

	for k, v := range m {
		if i != 0 {
			if err = e.Emitter.EmitMapNext(); err != nil {
				return
			}
		}
		if err = e.Encode(k); err != nil {
			return
		}
		if err = e.Emitter.EmitMapValue(); err != nil {
			return
		}
		if err = e.Encode(v); err != nil {
			return
		}
		i++
	}

	return e.Emitter.EmitMapEnd()
}

func (e Encoder) encodeMapStringInterfaceValue(v reflect.Value) error {
	return e.encodeMapStringInterface(v.Interface().(map[string]interface{}))
}

func (e Encoder) encodeMapStringInterface(m map[string]interface{}) (err error) {
	n := len(m)
	i := 0

	if err = e.Emitter.EmitMapBegin(n); err != nil {
		return
	}

	for k, v := range m {
		if i != 0 {
			if err = e.Emitter.EmitMapNext(); err != nil {
				return
			}
		}
		if err = e.Emitter.EmitString(k); err != nil {
			return
		}
		if err = e.Emitter.EmitMapValue(); err != nil {
			return
		}
		if err = e.Encode(v); err != nil {
			return
		}
		i++
	}

	return e.Emitter.EmitMapEnd()
}

func (e Encoder) encodeMapStringStringValue(v reflect.Value) error {
	return e.encodeMapStringString(v.Interface().(map[string]string))
}

func (e Encoder) encodeMapStringString(m map[string]string) (err error) {
	n := len(m)
	i := 0

	if err = e.Emitter.EmitMapBegin(n); err != nil {
		return
	}

	for k, v := range m {
		if i != 0 {
			if err = e.Emitter.EmitMapNext(); err != nil {
				return
			}
		}
		if err = e.Emitter.EmitString(k); err != nil {
			return
		}
		if err = e.Emitter.EmitMapValue(); err != nil {
			return
		}
		if err = e.Emitter.EmitString(v); err != nil {
			return
		}
		i++
	}

	return e.Emitter.EmitMapEnd()
}

func (e Encoder) encodeStruct(v reflect.Value) error {
	return e.encodeStructWith(v, structCache.lookup(v.Type()))
}

func (e Encoder) encodeStructWith(v reflect.Value, s *structType) (err error) {
	n := 0

	for i := range s.fields {
		f := &s.fields[i]
		if !f.omit(v.FieldByIndex(f.index)) {
			n++
		}
	}

	if err = e.Emitter.EmitMapBegin(n); err != nil {
		return
	}
	n = 0

	for i := range s.fields {
		f := &s.fields[i]
		if fv := v.FieldByIndex(f.index); !f.omit(fv) {
			if n != 0 {
				if err = e.Emitter.EmitMapNext(); err != nil {
					return
				}
			}
			if err = e.Emitter.EmitString(f.name); err != nil {
				return
			}
			if err = e.Emitter.EmitMapValue(); err != nil {
				return
			}
			if err = f.encode(e, fv); err != nil {
				return
			}
			n++
		}
	}

	return e.Emitter.EmitMapEnd()
}

func (e Encoder) encodePointer(v reflect.Value) error {
	return e.encodePointerWith(v, encodeFuncOf(v.Type().Elem()))
}

func (e Encoder) encodePointerWith(v reflect.Value, f encodeFunc) error {
	if v.IsNil() {
		return e.Emitter.EmitNil()
	}
	return f(e, v.Elem())
}

func (e Encoder) encodeInterface(v reflect.Value) error {
	if v.IsNil() {
		return e.Emitter.EmitNil()
	}
	return e.encode(v.Elem())
}

func (e Encoder) encodeEncoder(v reflect.Value) error {
	return v.Interface().(ValueEncoder).EncodeValue(e)
}

func (e Encoder) encodeMarshaler(v reflect.Value) error {
	if isTextEmitter(e.Emitter) {
		return e.encodeTextMarshaler(v)
	}
	return e.encodeBinaryMarshaler(v)
}

func (e Encoder) encodeBinaryMarshaler(v reflect.Value) error {
	b, err := v.Interface().(encoding.BinaryMarshaler).MarshalBinary()
	if err == nil {
		err = e.Emitter.EmitBytes(b)
	}
	return err
}

func (e Encoder) encodeTextMarshaler(v reflect.Value) error {
	b, err := v.Interface().(encoding.TextMarshaler).MarshalText()
	if err == nil {
		err = e.Emitter.EmitString(stringNoCopy(b))
	}
	return err
}

func (e Encoder) encodeUnsupported(v reflect.Value) error {
	return fmt.Errorf("objconv: the encoder doesn't support values of type %s", v.Type())
}

// EncodeArray provides the implementation of the array encoding algorithm,
// where n is the number of elements in the array, and f a function called to
// encode each element.
//
// The n argument can be set to a negative value to indicate that the program
// doesn't know how many elements it will output to the array. Be mindful that
// not all emitters support encoding arrays of unknown lengths.
//
// The f function is called to encode each element of the array.
func (e Encoder) EncodeArray(n int, f func(Encoder) error) (err error) {
	if e.key {
		if e.key, err = false, e.Emitter.EmitMapValue(); err != nil {
			return
		}
	}

	if err = e.Emitter.EmitArrayBegin(n); err != nil {
		return
	}

encodeArray:
	for i := 0; n < 0 || i < n; i++ {
		if i != 0 {
			if e.Emitter.EmitArrayNext(); err != nil {
				return
			}
		}
		switch err = f(e); err {
		case nil:
		case End:
			break encodeArray
		default:
			return
		}
	}

	return e.Emitter.EmitArrayEnd()
}

// EncodeMap provides the implementation of the map encoding algorithm, where n
// is the number of elements in the map, and f a function called to encode each
// element.
//
// The n argument can be set to a negative value to indicate that the program
// doesn't know how many elements it will output to the map. Be mindful that not
// all emitters support encoding maps of unknown length.
//
// The f function is called to encode each element of the map, it is expected to
// encode two values, the first one being the key, follow by the associated value.
// The first encoder must be used to encode the key, the second for the value.
func (e Encoder) EncodeMap(n int, f func(Encoder, Encoder) error) (err error) {
	if e.key {
		if e.key, err = false, e.Emitter.EmitMapValue(); err != nil {
			return
		}
	}

	if err = e.Emitter.EmitMapBegin(n); err != nil {
		return
	}

encodeMap:
	for i := 0; n < 0 || i < n; i++ {
		if i != 0 {
			if err = e.Emitter.EmitMapNext(); err != nil {
				return
			}
		}
		e.key = true
		err = f(
			Encoder{Emitter: e.Emitter, SortMapKeys: e.SortMapKeys},
			Encoder{Emitter: e.Emitter, SortMapKeys: e.SortMapKeys, key: true},
		)
		// Because internal calls don't use the exported methods they may not
		// reset this flag to false when expected, forcing the value here.
		e.key = false

		switch err {
		case nil:
		case End:
			break encodeMap
		default:
			return
		}
	}

	return e.Emitter.EmitMapEnd()
}

// A StreamEncoder encodes and writes a stream of values to an output stream.
//
// Instances of StreamEncoder are not safe for use by multiple goroutines.
type StreamEncoder struct {
	Emitter     Emitter // the emitter used by this encoder
	SortMapKeys bool    // whether map keys should be sorted

	err     error
	max     int
	cnt     int
	opened  bool
	closed  bool
	oneshot bool
}

// NewStreamEncoder returns a new stream encoder that outputs to e.
//
// The function panics if e is nil.
func NewStreamEncoder(e Emitter) *StreamEncoder {
	if e == nil {
		panic("objconv.NewStreamEncoder: the emitter is nil")
	}
	return &StreamEncoder{Emitter: e}
}

// Open explicitly tells the encoder to start the stream, setting the number
// of values to n.
//
// Depending on the actual format that the stream is encoding to, n may or
// may not have to be accurate, some formats also support passing a negative
// value to indicate that the number of elements is unknown.
func (e *StreamEncoder) Open(n int) error {
	if err := e.err; err != nil {
		return err
	}

	if e.closed {
		return io.ErrClosedPipe
	}

	if !e.opened {
		e.max = n
		e.opened = true

		if !e.oneshot {
			e.err = e.Emitter.EmitArrayBegin(n)
		}
	}

	return e.err
}

// Close terminates the stream encoder.
func (e *StreamEncoder) Close() error {
	if !e.closed {
		if err := e.Open(-1); err != nil {
			return err
		}

		e.closed = true

		if !e.oneshot {
			e.err = e.Emitter.EmitArrayEnd()
		}
	}

	return e.err
}

// Encode writes v to the stream, encoding it based on the emitter configured
// on e.
func (e *StreamEncoder) Encode(v interface{}) error {
	if err := e.Open(-1); err != nil {
		return err
	}

	if e.max >= 0 && e.cnt >= e.max {
		return fmt.Errorf("objconv: too many values sent to a stream encoder exceed the configured limit of %d", e.max)
	}

	if !e.oneshot && e.cnt != 0 {
		e.err = e.Emitter.EmitArrayNext()
	}

	if e.err == nil {
		e.err = (Encoder{
			Emitter:     e.Emitter,
			SortMapKeys: e.SortMapKeys,
		}).Encode(v)

		if e.cnt++; e.max >= 0 && e.cnt >= e.max {
			e.Close()
		}
	}

	return e.err
}

// ValueEncoder is the interface that can be implemented by types that wish to
// provide their own encoding algorithms.
//
// The EncodeValue method is called when the value is found by an encoding
// algorithm.
type ValueEncoder interface {
	EncodeValue(Encoder) error
}

// ValueEncoderFunc allows the use of regular functions or methods as value
// encoders.
type ValueEncoderFunc func(Encoder) error

// EncodeValue calls f(e).
func (f ValueEncoderFunc) EncodeValue(e Encoder) error { return f(e) }

// encodeFuncOpts is used to configure how the encodeFuncOf behaves.
type encodeFuncOpts struct {
	recurse bool
	structs map[reflect.Type]*structType
}

// encodeFunc is the prototype of functions that encode values.
type encodeFunc func(Encoder, reflect.Value) error

// encodeFuncOf returns an encoder function for t.
func encodeFuncOf(t reflect.Type) encodeFunc {
	return makeEncodeFunc(t, encodeFuncOpts{})
}

func makeEncodeFunc(t reflect.Type, opts encodeFuncOpts) encodeFunc {
	if adapter, ok := AdapterOf(t); ok {
		return adapter.Encode
	}

	switch t {
	case boolType:
		return Encoder.encodeBool

	case stringType:
		return Encoder.encodeString

	case bytesType:
		return Encoder.encodeBytes

	case timeType, timePtrType:
		return Encoder.encodeTime

	case durationType:
		return Encoder.encodeDuration

	case emptyInterface:
		return Encoder.encodeInterface

	case intType:
		return Encoder.encodeInt

	case int8Type:
		return Encoder.encodeInt8

	case int16Type:
		return Encoder.encodeInt16

	case int32Type:
		return Encoder.encodeInt32

	case int64Type:
		return Encoder.encodeInt64

	case uintType:
		return Encoder.encodeUint

	case uint8Type:
		return Encoder.encodeUint8

	case uint16Type:
		return Encoder.encodeUint16

	case uint32Type:
		return Encoder.encodeUint32

	case uint64Type:
		return Encoder.encodeUint64

	case uintptrType:
		return Encoder.encodeUintptr

	case float32Type:
		return Encoder.encodeFloat32

	case float64Type:
		return Encoder.encodeFloat64
	}

	switch {
	case t.Implements(valueEncoderInterface):
		return Encoder.encodeEncoder

	case t.Implements(binaryMarshalerInterface) && t.Implements(textMarshalerInterface):
		return Encoder.encodeMarshaler

	case t.Implements(binaryMarshalerInterface):
		return Encoder.encodeBinaryMarshaler

	case t.Implements(textMarshalerInterface):
		return Encoder.encodeTextMarshaler

	case t.Implements(errorInterface):
		return Encoder.encodeError
	}

	switch t.Kind() {
	case reflect.Struct:
		return makeEncodeStructFunc(t, opts)

	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			return Encoder.encodeBytes
		}
		return makeEncodeArrayFunc(t, opts)

	case reflect.Map:
		return makeEncodeMapFunc(t, opts)

	case reflect.Ptr:
		return makeEncodePtrFunc(t, opts)

	case reflect.Array:
		return makeEncodeArrayFunc(t, opts)

	case reflect.String:
		return Encoder.encodeString

	case reflect.Bool:
		return Encoder.encodeBool

	case reflect.Int:
		return Encoder.encodeInt

	case reflect.Int8:
		return Encoder.encodeInt8

	case reflect.Int16:
		return Encoder.encodeInt16

	case reflect.Int32:
		return Encoder.encodeInt32

	case reflect.Int64:
		return Encoder.encodeInt64

	case reflect.Uint:
		return Encoder.encodeUint

	case reflect.Uint8:
		return Encoder.encodeUint8

	case reflect.Uint16:
		return Encoder.encodeUint16

	case reflect.Uint32:
		return Encoder.encodeUint32

	case reflect.Uint64:
		return Encoder.encodeUint64

	case reflect.Uintptr:
		return Encoder.encodeUintptr

	case reflect.Float32:
		return Encoder.encodeFloat32

	case reflect.Float64:
		return Encoder.encodeFloat64

	default:
		return Encoder.encodeUnsupported
	}
}

func makeEncodeArrayFunc(t reflect.Type, opts encodeFuncOpts) encodeFunc {
	if !opts.recurse {
		return Encoder.encodeArray
	}
	f := makeEncodeFunc(t.Elem(), opts)
	return func(e Encoder, v reflect.Value) error {
		return e.encodeArrayWith(v, f)
	}
}

func makeEncodeMapFunc(t reflect.Type, opts encodeFuncOpts) encodeFunc {
	if !opts.recurse {
		return Encoder.encodeMap
	}
	kf := makeEncodeFunc(t.Key(), opts)
	vf := makeEncodeFunc(t.Elem(), opts)
	return func(e Encoder, v reflect.Value) error {
		return e.encodeMapWith(v, kf, vf)
	}
}

func makeEncodeStructFunc(t reflect.Type, opts encodeFuncOpts) encodeFunc {
	if !opts.recurse {
		return Encoder.encodeStruct
	}
	s := newStructType(t, opts.structs)
	return func(e Encoder, v reflect.Value) error {
		return e.encodeStructWith(v, s)
	}
}

func makeEncodePtrFunc(t reflect.Type, opts encodeFuncOpts) encodeFunc {
	if !opts.recurse {
		return Encoder.encodePointer
	}
	f := makeEncodeFunc(t.Elem(), opts)
	return func(e Encoder, v reflect.Value) error {
		return e.encodePointerWith(v, f)
	}
}
