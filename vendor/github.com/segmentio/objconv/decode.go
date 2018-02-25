package objconv

import (
	"encoding"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/segmentio/objconv/objutil"
)

// A Decoder implements the algorithms for building data structures from their
// serialized forms.
//
// Decoders are not safe for use by multiple goroutines.
type Decoder struct {
	// Parser to use to load values.
	Parser Parser

	// MapType is used to override the type of maps produced by the decoder when
	// there is not destination type (when decoding to an empty interface).
	MapType reflect.Type

	off int // offset of the value when decoding a map
}

// NewDecoder returns a decoder object that uses p, will panic if p is nil.
func NewDecoder(p Parser) *Decoder {
	if p == nil {
		panic("objconv: the parser is nil")
	}
	return &Decoder{Parser: p}
}

// Decode expects v to be a pointer to a value in which the decoder will load
// the next parsed data.
//
// The method panics if v is neither a pointer type nor implements the
// ValueDecoder interface, or if v is a nil pointer.
func (d Decoder) Decode(v interface{}) error {
	to := reflect.ValueOf(v)

	if d.off != 0 {
		var err error
		if d.off, err = 0, d.Parser.ParseMapValue(d.off-1); err != nil {
			return err
		}
	}

	if !to.IsValid() {
		// This special case for a nil value is used to make it possible to
		// discard decoded values.
		_, err := d.decodeInterface(to)
		return err
	}

	// Optimization for ValueDecoder, in practice tho it's also handled in the
	// methods that are based on reflection.
	switch x := v.(type) {
	case ValueDecoder:
		return x.DecodeValue(d)
	}

	if to.Kind() == reflect.Ptr {
		// In most cases the method receives a pointer, but we may also have to
		// support types that aren't pointers but implement ValueDecoder, or
		// types that have got adapters set.
		// If we're not in either of those cases the code will likely panic when
		// the value is set because it won't be addressable.
		to = to.Elem()
	}

	_, err := d.decode(to)
	return err
}

func (d Decoder) decode(to reflect.Value) (Type, error) {
	return decodeFuncOf(to.Type())(d, to)
}

func (d Decoder) decodeBool(to reflect.Value) (t Type, err error) {
	if t, err = d.Parser.ParseType(); err == nil {
		err = d.decodeBoolFromType(t, to)
	}
	return
}

func (d Decoder) decodeBoolFromType(t Type, to reflect.Value) (err error) {
	var v bool

	switch t {
	case Nil:
		err = d.Parser.ParseNil()

	case Bool:
		v, err = d.Parser.ParseBool()

	default:
		err = typeConversionError(t, Bool)
	}

	if err != nil {
		return
	}

	if to.IsValid() {
		to.SetBool(v)
	}
	return
}

func (d Decoder) decodeInt(to reflect.Value) (t Type, err error) {
	if t, err = d.Parser.ParseType(); err == nil {
		err = d.decodeIntFromType(t, to)
	}
	return
}

func (d Decoder) decodeIntFromType(t Type, to reflect.Value) (err error) {
	var valid = to.IsValid()
	var i int64
	var u uint64

	switch t {
	case Nil:
		err = d.Parser.ParseNil()

	case Int:
		if i, err = d.Parser.ParseInt(); err != nil {
			return
		}

		if valid {
			switch t := to.Type(); t.Kind() {
			case reflect.Int:
				err = objutil.CheckInt64Bounds(i, int64(objutil.IntMin), uint64(objutil.IntMax), t)
			case reflect.Int8:
				err = objutil.CheckInt64Bounds(i, objutil.Int8Min, objutil.Int8Max, t)
			case reflect.Int16:
				err = objutil.CheckInt64Bounds(i, objutil.Int16Min, objutil.Int16Max, t)
			case reflect.Int32:
				err = objutil.CheckInt64Bounds(i, objutil.Int32Min, objutil.Int32Max, t)
			}
		}

	case Uint:
		if u, err = d.Parser.ParseUint(); err != nil {
			return
		}

		if valid {
			switch t := to.Type(); t.Kind() {
			case reflect.Int:
				err = objutil.CheckUint64Bounds(u, uint64(objutil.IntMax), t)
			case reflect.Int8:
				err = objutil.CheckUint64Bounds(u, objutil.Int8Max, t)
			case reflect.Int16:
				err = objutil.CheckUint64Bounds(u, objutil.Int16Max, t)
			case reflect.Int32:
				err = objutil.CheckUint64Bounds(u, objutil.Int32Max, t)
			case reflect.Int64:
				err = objutil.CheckUint64Bounds(u, objutil.Int64Max, t)
			}
		}

		i = int64(u)

	case String:
		var b []byte

		if b, err = d.Parser.ParseString(); err != nil {
			return
		}

		i, err = strconv.ParseInt(string(b), 10, 64)

	case Bytes:
		var b []byte

		if b, err = d.Parser.ParseBytes(); err != nil {
			return
		}

		i, err = strconv.ParseInt(string(b), 10, 64)

	default:
		err = typeConversionError(t, Int)
	}

	if err != nil {
		return
	}

	if valid {
		to.SetInt(i)
	}
	return
}

func (d Decoder) decodeUint(to reflect.Value) (t Type, err error) {
	if t, err = d.Parser.ParseType(); err == nil {
		err = d.decodeUintFromType(t, to)
	}
	return
}

func (d Decoder) decodeUintFromType(t Type, to reflect.Value) (err error) {
	var valid = to.IsValid()
	var i int64
	var u uint64

	switch t {
	case Nil:
		err = d.Parser.ParseNil()

	case Int:
		if i, err = d.Parser.ParseInt(); err != nil {
			return
		}

		if valid {
			switch t := to.Type(); t.Kind() {
			case reflect.Uint:
				err = objutil.CheckInt64Bounds(i, 0, uint64(objutil.UintMax), t)
			case reflect.Uint8:
				err = objutil.CheckInt64Bounds(i, 0, objutil.Uint8Max, t)
			case reflect.Uint16:
				err = objutil.CheckInt64Bounds(i, 0, objutil.Uint16Max, t)
			case reflect.Uint32:
				err = objutil.CheckInt64Bounds(i, 0, objutil.Uint32Max, t)
			case reflect.Uint64:
				err = objutil.CheckInt64Bounds(i, 0, objutil.Uint64Max, t)
			}
		}

		u = uint64(i)

	case Uint:
		if u, err = d.Parser.ParseUint(); err != nil {
			return
		}

		if valid {
			switch t := to.Type(); t.Kind() {
			case reflect.Uint:
				err = objutil.CheckUint64Bounds(u, uint64(objutil.UintMax), t)
			case reflect.Uint8:
				err = objutil.CheckUint64Bounds(u, objutil.Uint8Max, t)
			case reflect.Uint16:
				err = objutil.CheckUint64Bounds(u, objutil.Uint16Max, t)
			case reflect.Uint32:
				err = objutil.CheckUint64Bounds(u, objutil.Uint32Max, t)
			}
		}

	case String:
		var b []byte

		if b, err = d.Parser.ParseString(); err != nil {
			return
		}

		u, err = strconv.ParseUint(string(b), 10, 64)

	case Bytes:
		var b []byte

		if b, err = d.Parser.ParseBytes(); err != nil {
			return
		}

		u, err = strconv.ParseUint(string(b), 10, 64)

	default:
		err = typeConversionError(t, Uint)
	}

	if err != nil {
		return
	}

	if valid {
		to.SetUint(u)
	}
	return
}

func (d Decoder) decodeFloat(to reflect.Value) (t Type, err error) {
	if t, err = d.Parser.ParseType(); err == nil {
		err = d.decodeFloatFromType(t, to)
	}
	return
}

func (d Decoder) decodeFloatFromType(t Type, to reflect.Value) (err error) {
	var i int64
	var u uint64
	var f float64

	switch t {
	case Nil:
		err = d.Parser.ParseNil()

	case Int:
		if i, err = d.Parser.ParseInt(); err == nil {
			if err = objutil.CheckInt64Bounds(i, objutil.Float64IntMin, objutil.Float64IntMax, int64Type); err == nil {
				f = float64(i)
			}
		}

	case Uint:
		if u, err = d.Parser.ParseUint(); err == nil {
			if err = objutil.CheckUint64Bounds(u, objutil.Float64IntMax, uint64Type); err == nil {
				f = float64(u)
			}
		}

	case Float:
		f, err = d.Parser.ParseFloat()

	case String:
		var b []byte

		if b, err = d.Parser.ParseString(); err != nil {
			return
		}

		f, err = strconv.ParseFloat(string(b), 64)

	case Bytes:
		var b []byte

		if b, err = d.Parser.ParseBytes(); err != nil {
			return
		}

		f, err = strconv.ParseFloat(string(b), 64)

	default:
		err = typeConversionError(t, Float)
	}

	if err != nil {
		return
	}

	if to.IsValid() {
		to.SetFloat(f)
	}
	return
}

func (d Decoder) decodeString(to reflect.Value) (t Type, err error) {
	if t, err = d.Parser.ParseType(); err == nil {
		err = d.decodeStringFromType(t, to)
	}
	return
}

func (d Decoder) decodeStringFromType(t Type, to reflect.Value) (err error) {
	var a [64]byte
	var b []byte

	switch t {
	case Nil:
		err = d.Parser.ParseNil()

	case String:
		b, err = d.Parser.ParseString()

	case Bytes:
		b, err = d.Parser.ParseBytes()

	case Bool:
		var v bool
		if v, err = d.Parser.ParseBool(); err == nil {
			if v {
				b = append(a[:0], "true"...)
			} else {
				b = append(a[:0], "false"...)
			}
		}

	case Int:
		var v int64
		if v, err = d.Parser.ParseInt(); err == nil {
			b = strconv.AppendInt(a[:0], v, 10)
		}

	case Uint:
		var v uint64
		if v, err = d.Parser.ParseUint(); err == nil {
			b = strconv.AppendUint(a[:0], v, 10)
		}

	case Float:
		var v float64
		if v, err = d.Parser.ParseFloat(); err == nil {
			b = strconv.AppendFloat(a[:0], v, 'g', -1, 64)
		}

	case Time:
		var v time.Time
		if v, err = d.Parser.ParseTime(); err == nil {
			b = v.AppendFormat(a[:0], time.RFC3339Nano)
		}

	case Duration:
		var v time.Duration
		if v, err = d.Parser.ParseDuration(); err == nil {
			b = objutil.AppendDuration(a[:0], v)
		}

	case Error:
		var v error
		if v, err = d.Parser.ParseError(); err == nil {
			b = append(a[:0], v.Error()...)
		}

	default:
		err = typeConversionError(t, String)
	}

	if err != nil {
		return
	}

	if to.IsValid() {
		to.SetString(string(b))
	}
	return
}

func (d Decoder) decodeBytes(to reflect.Value) (t Type, err error) {
	if t, err = d.Parser.ParseType(); err == nil {
		err = d.decodeBytesFromType(t, to)
	}
	return
}

func (d Decoder) decodeBytesFromType(t Type, to reflect.Value) (err error) {
	var b []byte

	switch t {
	case Nil:
		err = d.Parser.ParseNil()

	case String:
		b, err = d.Parser.ParseString()

	case Bytes:
		b, err = d.Parser.ParseBytes()

	default:
		err = typeConversionError(t, String)
	}

	if err != nil {
		return
	}

	if bd, ok := d.Parser.(bytesDecoder); ok {
		if b, err = bd.DecodeBytes(b); err != nil {
			return
		}
	}

	if to.IsValid() {
		if t == Nil {
			to.SetBytes(nil)
		} else {
			v := make([]byte, len(b))
			copy(v, b)
			to.SetBytes(v)
		}
	}
	return
}

func (d Decoder) decodeTime(to reflect.Value) (t Type, err error) {
	if t, err = d.Parser.ParseType(); err == nil {
		err = d.decodeTimeFromType(t, to)
	}
	return
}

func (d Decoder) decodeTimeFromType(t Type, to reflect.Value) (err error) {
	var s []byte
	var v time.Time

	switch t {
	case Nil:
		err = d.Parser.ParseNil()

	case String:
		s, err = d.Parser.ParseString()

	case Bytes:
		s, err = d.Parser.ParseBytes()

	case Time:
		v, err = d.Parser.ParseTime()
	}

	if err != nil {
		return
	}

	if to.IsValid() {
		if t == String || t == Bytes {
			v, err = time.Parse(time.RFC3339Nano, string(s))
		}
		*(to.Addr().Interface().(*time.Time)) = v
	}
	return
}

func (d Decoder) decodeDuration(to reflect.Value) (t Type, err error) {
	if t, err = d.Parser.ParseType(); err == nil {
		err = d.decodeDurationFromType(t, to)
	}
	return
}

func (d Decoder) decodeDurationFromType(t Type, to reflect.Value) (err error) {
	var s []byte
	var v time.Duration

	switch t {
	case Nil:
		err = d.Parser.ParseNil()

	case String:
		s, err = d.Parser.ParseString()

	case Bytes:
		s, err = d.Parser.ParseBytes()

	case Duration:
		v, err = d.Parser.ParseDuration()
	}

	if err != nil {
		return
	}

	if t == String || t == Bytes {
		v, err = time.ParseDuration(string(s))
	}

	if to.IsValid() {
		to.SetInt(int64(v))
	}
	return
}

func (d Decoder) decodeError(to reflect.Value) (t Type, err error) {
	if t, err = d.Parser.ParseType(); err == nil {
		err = d.decodeErrorFromType(t, to)
	}
	return
}

func (d Decoder) decodeErrorFromType(t Type, to reflect.Value) (err error) {
	var s []byte
	var v error

	switch t {
	case Nil:
		err = d.Parser.ParseNil()

	case String:
		s, err = d.Parser.ParseString()

	case Bytes:
		s, err = d.Parser.ParseBytes()

	case Error:
		v, err = d.Parser.ParseError()
	}

	if err != nil {
		return
	}

	if to.IsValid() {
		if t == String || t == Bytes {
			v = errors.New(string(s))
		}
		to.Set(reflect.ValueOf(v))
	}
	return
}

func (d Decoder) decodeSlice(to reflect.Value) (t Type, err error) {
	return d.decodeSliceWith(to, decodeFuncOf(to.Type().Elem()))
}

func (d Decoder) decodeSliceWith(to reflect.Value, f decodeFunc) (t Type, err error) {
	if t, err = d.Parser.ParseType(); err == nil {
		err = d.decodeSliceFromTypeWith(t, to, f)
	}
	return
}

func (d Decoder) decodeSliceFromType(typ Type, to reflect.Value) (err error) {
	f := Decoder.decodeInterface
	if to.IsValid() {
		f = decodeFuncOf(to.Type().Elem())
	}
	return d.decodeSliceFromTypeWith(typ, to, f)
}

func (d Decoder) decodeSliceFromTypeWith(typ Type, to reflect.Value, f decodeFunc) (err error) {
	if !to.IsValid() {
		return d.decodeArrayImpl(typ, func(d Decoder) (err error) {
			_, err = f(d, reflect.Value{})
			return
		})
	}

	t := to.Type()
	s := reflect.MakeSlice(t, 0, 0)
	i := 0
	n := 0

	if err = d.decodeArrayImpl(typ, func(d Decoder) (err error) {
		if i == n {
			if n *= 5; n == 0 {
				n = 10
			}
			sc := reflect.MakeSlice(t, n, n)
			reflect.Copy(sc, s)
			s = sc
		}
		if _, err = f(d, s.Index(i)); err != nil {
			return
		}
		i++
		return
	}); err != nil {
		return
	}

	if typ == Nil {
		to.Set(zeroValueOf(t))
	} else {
		if i != n {
			s = s.Slice(0, i)
		}
		to.Set(s)
	}
	return
}

func (d Decoder) decodeArray(to reflect.Value) (t Type, err error) {
	return d.decodeArrayWith(to, decodeFuncOf(to.Type().Elem()))
}

func (d Decoder) decodeArrayWith(to reflect.Value, f decodeFunc) (t Type, err error) {
	if t, err = d.Parser.ParseType(); err == nil {
		err = d.decodeArrayFromTypeWith(t, to, f)
	}
	return
}

func (d Decoder) decodeArrayFromTypeWith(typ Type, to reflect.Value, f decodeFunc) (err error) {
	n := to.Len()       // len(to)
	t := to.Type()      // [...]T
	e := t.Elem()       // T
	z := zeroValueOf(e) // T{}

	for i := 0; i != n; i++ {
		to.Index(i).Set(z) // reset to the zero-value
	}

	i := 0

	if err = d.decodeArrayImpl(typ, func(d Decoder) (err error) {
		if i < n {
			if _, err = f(d, to.Index(i)); err != nil {
				return
			}
		}
		i++
		return
	}); err != nil {
		return
	}

	if typ == Nil {
		to.Set(zeroValueOf(t))
	} else if i != n {
		err = fmt.Errorf("objconv: array length mismatch, expected %d but only %d elements were decoded", n, i)
	}

	return
}

func (d Decoder) decodeMap(to reflect.Value) (Type, error) {
	t := to.Type()
	return d.decodeMapWith(to, decodeFuncOf(t.Key()), decodeFuncOf(t.Elem()))
}

func (d Decoder) decodeMapWith(to reflect.Value, kf decodeFunc, vf decodeFunc) (t Type, err error) {
	if t, err = d.Parser.ParseType(); err == nil {
		err = d.decodeMapFromTypeWith(t, to, kf, vf)
	}
	return
}

func (d Decoder) decodeMapFromType(typ Type, to reflect.Value) (err error) {
	kf := Decoder.decodeInterface
	vf := Decoder.decodeInterface
	if to.IsValid() {
		t := to.Type()
		kf = decodeFuncOf(t.Key())
		vf = decodeFuncOf(t.Elem())
	}
	return d.decodeMapFromTypeWith(typ, to, kf, vf)
}

func (d Decoder) decodeMapFromTypeWith(typ Type, to reflect.Value, kf decodeFunc, vf decodeFunc) (err error) {
	if !to.IsValid() {
		return d.decodeMapImpl(typ, func(kd Decoder, vd Decoder) (err error) {
			if _, err = d.decodeInterface(reflect.Value{}); err != nil {
				return
			}
			if err = d.Parser.ParseMapValue(vd.off - 1); err != nil {
				return
			}
			_, err = d.decodeInterface(reflect.Value{})
			return
		})
	}

	t := to.Type() // map[K]V

	switch t {
	case mapInterfaceInterfaceType:
		return d.decodeMapInterfaceInterface(typ, to)

	case mapStringInterfaceType:
		return d.decodeMapStringInterface(typ, to)

	case mapStringStringType:
		return d.decodeMapStringString(typ, to)
	}

	m := reflect.MakeMap(t) // make(map[K]V)

	kt := t.Key()                // K
	kz := zeroValueOf(kt)        // K{}
	kv := reflect.New(kt).Elem() // &K{}

	vt := t.Elem()               // V
	vz := zeroValueOf(vt)        // V{}
	vv := reflect.New(vt).Elem() // &V{}

	if err = d.decodeMapImpl(typ, func(kd Decoder, vd Decoder) (err error) {
		kv.Set(kz) // reset the key to its zero-value
		vv.Set(vz) // reset the value to its zero-value
		if _, err = kf(d, kv); err != nil {
			return
		}
		if err = d.Parser.ParseMapValue(vd.off - 1); err != nil {
			return
		}
		if _, err = vf(d, vv); err != nil {
			return
		}
		m.SetMapIndex(kv, vv)
		return
	}); err != nil {
		return
	}

	if typ == Nil {
		to.Set(zeroValueOf(t))
	} else {
		to.Set(m)
	}
	return
}

func (d Decoder) decodeMapInterfaceInterface(typ Type, to reflect.Value) error {
	m := to.Interface().(map[interface{}]interface{})

	if m == nil {
		m = make(map[interface{}]interface{})
		to.Set(reflect.ValueOf(m))
	}

	for k := range m {
		delete(m, k)
	}

	return d.decodeMapImpl(typ, func(kd Decoder, vd Decoder) (err error) {
		var k interface{}
		var v interface{}

		if err = kd.Decode(&k); err != nil {
			return
		}
		if err = vd.Decode(&v); err != nil {
			return
		}

		m[k] = v
		return
	})
}

func (d Decoder) decodeMapStringInterface(typ Type, to reflect.Value) (err error) {
	m := to.Interface().(map[string]interface{})

	if m == nil {
		m = make(map[string]interface{})
		to.Set(reflect.ValueOf(m))
	}

	for k := range m {
		delete(m, k)
	}

	return d.decodeMapImpl(typ, func(kd Decoder, vd Decoder) (err error) {
		var b []byte
		var k string
		var v interface{}

		if _, b, err = d.decodeTypeAndString(); err != nil {
			return
		}
		k = string(b)

		if err = vd.Decode(&v); err != nil {
			return
		}

		m[k] = v
		return
	})
}

func (d Decoder) decodeMapStringString(typ Type, to reflect.Value) (err error) {
	m := to.Interface().(map[string]string)

	if m == nil {
		m = make(map[string]string)
		to.Set(reflect.ValueOf(m))
	}

	for k := range m {
		delete(m, k)
	}

	return d.decodeMapImpl(typ, func(kd Decoder, vd Decoder) (err error) {
		var b []byte
		var k string
		var v string

		if _, b, err = d.decodeTypeAndString(); err != nil {
			return
		}
		k = string(b)

		if err = d.Parser.ParseMapValue(vd.off - 1); err != nil {
			return
		}

		if _, b, err = d.decodeTypeAndString(); err != nil {
			return
		}
		v = string(b)

		m[k] = v
		return
	})
}

func (d Decoder) decodeStruct(to reflect.Value) (Type, error) {
	return d.decodeStructWith(to, structCache.lookup(to.Type()))
}

func (d Decoder) decodeStructWith(to reflect.Value, s *structType) (t Type, err error) {
	if t, err = d.Parser.ParseType(); err == nil {
		err = d.decodeStructFromTypeWith(t, to, s)
	}
	return
}

func (d Decoder) decodeStructFromTypeWith(typ Type, to reflect.Value, s *structType) (err error) {
	if err = d.decodeMapImpl(typ, func(kd Decoder, vd Decoder) (err error) {
		var b []byte

		if _, b, err = d.decodeTypeAndString(); err != nil {
			return
		}
		f := s.fieldsByName[string(b)]

		if err = d.Parser.ParseMapValue(vd.off - 1); err != nil {
			return
		}

		if f == nil {
			_, err = d.decodeInterface(reflect.Value{}) // discard
			return
		}

		_, err = f.decode(d, to.FieldByIndex(f.index))
		return
	}); err != nil {
		to.Set(zeroValueOf(to.Type()))
	}
	return
}

func (d Decoder) decodePointer(to reflect.Value) (Type, error) {
	return d.decodePointerWith(to, decodeFuncOf(to.Type().Elem()))
}

func (d Decoder) decodePointerWith(to reflect.Value, f decodeFunc) (typ Type, err error) {
	var t = to.Type()
	var v reflect.Value

	if to.IsNil() {
		v = reflect.New(t.Elem())
	} else {
		v = to
	}

	if typ, err = f(d, v.Elem()); err != nil {
		return
	}

	switch {
	case typ == Nil:
		to.Set(zeroValueOf(t))
	case to.IsNil():
		to.Set(v)
	}

	return
}

func (d Decoder) decodeDecoderPointer(to reflect.Value) (Type, error) {
	return Unknown /* just needs to not be Nil */, to.Addr().Interface().(ValueDecoder).DecodeValue(d)
}

func (d Decoder) decodeDecoder(to reflect.Value) (Type, error) {
	return Unknown /* just needs to not be Nil */, to.Interface().(ValueDecoder).DecodeValue(d)
}

func (d Decoder) decodeUnmarshalerPointer(to reflect.Value) (Type, error) {
	return d.decodeUnmarshaler(to.Addr())
}

func (d Decoder) decodeUnmarshaler(to reflect.Value) (Type, error) {
	if isTextParser(d.Parser) {
		return d.decodeTextUnmarshaler(to)
	}
	return d.decodeBinaryUnmarshaler(to)
}

func (d Decoder) decodeBinaryUnmarshalerPointer(to reflect.Value) (Type, error) {
	return d.decodeBinaryUnmarshaler(to.Addr())
}

func (d Decoder) decodeBinaryUnmarshaler(to reflect.Value) (t Type, err error) {
	var b []byte
	var v = reflect.ValueOf(&b).Elem()

	if t, err = d.decodeBytes(v); err != nil {
		return
	}

	if to.Kind() == reflect.Ptr && to.IsNil() {
		to.Set(reflect.New(to.Type().Elem()))
	}

	err = to.Interface().(encoding.BinaryUnmarshaler).UnmarshalBinary(b)
	return
}

func (d Decoder) decodeTextUnmarshalerPointer(to reflect.Value) (Type, error) {
	return d.decodeTextUnmarshaler(to.Addr())
}

func (d Decoder) decodeTextUnmarshaler(to reflect.Value) (t Type, err error) {
	var s string
	var v = reflect.ValueOf(&s).Elem()

	if t, err = d.decodeString(v); err != nil {
		return
	}

	if to.Kind() == reflect.Ptr && to.IsNil() {
		to.Set(reflect.New(to.Type().Elem()))
	}

	err = to.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(s))
	return
}

func (d Decoder) decodeInterface(to reflect.Value) (t Type, err error) {
	if t, err = d.Parser.ParseType(); err == nil {
		err = d.decodeInterfaceFromType(t, to)
	}
	return
}

func (d Decoder) decodeInterfaceFromType(t Type, to reflect.Value) (err error) {
	switch t {
	case Nil:
		err = d.decodeInterfaceFromNil(to)
	case Bool:
		err = d.decodeInterfaceFrom(boolType, t, to, Decoder.decodeBoolFromType)
	case Int:
		err = d.decodeInterfaceFrom(int64Type, t, to, Decoder.decodeIntFromType)
	case Uint:
		err = d.decodeInterfaceFrom(uint64Type, t, to, Decoder.decodeUintFromType)
	case Float:
		err = d.decodeInterfaceFrom(float64Type, t, to, Decoder.decodeFloatFromType)
	case String:
		err = d.decodeInterfaceFrom(stringType, t, to, Decoder.decodeStringFromType)
	case Bytes:
		err = d.decodeInterfaceFrom(bytesType, t, to, Decoder.decodeBytesFromType)
	case Time:
		err = d.decodeInterfaceFrom(timeType, t, to, Decoder.decodeTimeFromType)
	case Duration:
		err = d.decodeInterfaceFrom(durationType, t, to, Decoder.decodeDurationFromType)
	case Error:
		err = d.decodeInterfaceFrom(errorInterface, t, to, Decoder.decodeErrorFromType)
	case Array:
		err = d.decodeInterfaceFrom(sliceInterfaceType, t, to, Decoder.decodeSliceFromType)
	case Map:
		if to.IsValid() && d.MapType != nil {
			v := reflect.New(d.MapType).Elem()
			_, err = d.decode(v)
			to.Set(v)
		} else {
			err = d.decodeInterfaceFrom(mapInterfaceInterfaceType, t, to, Decoder.decodeMapFromType)
		}
	default:
		panic("objconv: parser returned an unsupported value type: " + t.String())
	}
	return
}

func (d Decoder) decodeInterfaceFromNil(to reflect.Value) (err error) {
	if err = d.Parser.ParseNil(); err == nil {
		if to.IsValid() {
			to.Set(zeroValueOf(to.Type()))
		}
	}
	return
}

func (d Decoder) decodeInterfaceFrom(from reflect.Type, t Type, to reflect.Value, decode func(Decoder, Type, reflect.Value) error) (err error) {
	if !to.IsValid() {
		return decode(d, t, reflect.Value{})
	}

	v := reflect.New(from).Elem()

	if err = decode(d, t, v); err != nil {
		return
	}

	to.Set(v)
	return
}

func (d Decoder) decodeUnsupported(to reflect.Value) (Type, error) {
	return Nil, fmt.Errorf("objconv: the decoder doesn't support values of type %s", to.Type())
}

func (d Decoder) decodeTypeAndString() (t Type, b []byte, err error) {
	if t, err = d.Parser.ParseType(); err == nil {
		// This algorithm is the same than the one used in
		// decodeStringWithType, and should be kept in sync.
		switch t {
		case Nil:
			err = d.Parser.ParseNil()
		case String:
			b, err = d.Parser.ParseString()
		case Bytes:
			b, err = d.Parser.ParseBytes()
		default:
			err = typeConversionError(t, String)
		}
	}
	return
}

// DecodeArray provides the implementation of the algorithm for decoding arrays,
// where f is called to decode each element of the array.
func (d Decoder) DecodeArray(f func(Decoder) error) (err error) {
	var typ Type

	if d.off != 0 {
		if d.off, err = 0, d.Parser.ParseMapValue(d.off-1); err != nil {
			return
		}
	}

	if typ, err = d.Parser.ParseType(); err != nil {
		return
	}

	err = d.decodeArrayImpl(typ, f)
	return
}

func (d Decoder) decodeArrayImpl(t Type, f func(Decoder) error) (err error) {
	var n int

	switch t {
	case Nil:
		err = d.Parser.ParseNil()
		return

	case Array:
		n, err = d.Parser.ParseArrayBegin()

	default:
		err = typeConversionError(t, Array)
	}

	if err != nil {
		return
	}

	i := 0

	for n < 0 || i < n {
		if n < 0 || i != 0 {
			if err = d.Parser.ParseArrayNext(i); err != nil {
				if err == End {
					err = nil
					break
				}
				return
			}
		}
		if err = f(d); err != nil {
			return
		}
		i++
	}

	err = d.Parser.ParseArrayEnd(i)
	return
}

// DecodeMap provides the implementation of the algorithm for decoding maps,
// where f is called to decode each pair of key and value.
//
// The function f is expected to decode two values from the map, the first one
// being the key and the second the associated value. The first decoder must be
// used to decode the key, the second one for the value.
func (d Decoder) DecodeMap(f func(Decoder, Decoder) error) (err error) {
	var typ Type

	if d.off != 0 {
		if d.off, err = 0, d.Parser.ParseMapValue(d.off-1); err != nil {
			return
		}
	}

	if typ, err = d.Parser.ParseType(); err != nil {
		return
	}

	err = d.decodeMapImpl(typ, f)
	return
}

func (d Decoder) decodeMapImpl(t Type, f func(Decoder, Decoder) error) (err error) {
	var n int

	switch t {
	case Nil:
		err = d.Parser.ParseNil()
		return

	case Map:
		n, err = d.Parser.ParseMapBegin()

	default:
		err = typeConversionError(t, Map)
	}

	if err != nil {
		return
	}

	i := 0

	for n < 0 || i < n {
		if n < 0 || i != 0 {
			if err = d.Parser.ParseMapNext(i); err != nil {
				if err == End {
					err = nil
					break
				}
				return
			}
		}

		d1 := d
		d2 := d
		d2.off = i + 1

		if err = f(d1, d2); err != nil {
			return
		}

		i++
	}

	err = d.Parser.ParseMapEnd(i)
	return
}

// StreamDecoder decodes values in a streaming fashion, allowing an array to be
// consumed without loading it fully in memory.
//
// Instances of StreamDecoder are not safe for use by multiple goroutines.
type StreamDecoder struct {
	// Parser to use to load values.
	Parser Parser

	// MapType is used to override the type of maps produced by the decoder when
	// there is not destination type (when decoding to an empty interface).
	MapType reflect.Type

	err error
	typ Type
	cnt int
	max int
}

// NewStreamDecoder returns a new stream decoder that takes input from p.
//
// The function panics if p is nil.
func NewStreamDecoder(p Parser) *StreamDecoder {
	if p == nil {
		panic("objconv: the parser is nil")
	}
	return &StreamDecoder{Parser: p}
}

// Len returns the number of values remaining to be read from the stream, which
// may be -1 if the underlying format doesn't provide this information. If an
// error occurred while decoding the stream the method returns zero because no
// more values can be read.
func (d *StreamDecoder) Len() int {
	if d.err != nil {
		return 0
	}

	if d.typ == Unknown {
		if d.init() != nil {
			return 0
		}
	}

	return d.max - d.cnt
}

// Err returns the last error returned by the Decode method.
//
// The method returns nil if the stream reached its natural end.
func (d *StreamDecoder) Err() error {
	if d.err == End {
		return nil
	}
	return d.err
}

// Decodes the next value from the stream into v.
func (d *StreamDecoder) Decode(v interface{}) error {
	if d.err != nil {
		return d.err
	}

	err := error(nil)
	cnt := d.cnt
	max := d.max
	dec := Decoder{
		Parser:  d.Parser,
		MapType: d.MapType,
	}

	switch d.typ {
	case Unknown:
		err = d.init()
		max = d.max
	case Array:
		if cnt == max {
			err = dec.Parser.ParseArrayEnd(cnt)
		} else if cnt != 0 {
			err = dec.Parser.ParseArrayNext(cnt)
		}
	}

	if err == nil {
		if cnt == max {
			err = End
		} else {
			switch err = dec.Decode(v); err {
			case nil:
				cnt++
			case End:
				cnt++
				max = cnt
			default:
				if max < 0 && dec.Parser.ParseArrayEnd(cnt) == nil {
					err = End
				}
			}
		}
	}

	d.err = err
	d.cnt = cnt
	d.max = max
	return err
}

// Encoder returns a new StreamEncoder which can be used to re-encode the stream
// decoded by d into e.
//
// The method panics if e is nil.
func (d *StreamDecoder) Encoder(e Emitter) (enc *StreamEncoder, err error) {
	var typ Type

	if typ, err = d.Parser.ParseType(); err == nil {
		enc = NewStreamEncoder(e)
		enc.oneshot = typ != Array
	}

	return
}

func (d *StreamDecoder) init() error {
	err := error(nil)
	typ := Unknown
	max := 0

	if typ, err = d.Parser.ParseType(); err == nil {
		switch typ {
		default:
			max = 1
		case Array:
			max, err = d.Parser.ParseArrayBegin()
		}
	}

	d.err = err
	d.typ = typ
	d.max = max
	return err
}

// ValueDecoder is the interface that can be implemented by types that wish to
// provide their own decoding algorithms.
//
// The DecodeValue method is called when the value is found by a decoding
// algorithm.
type ValueDecoder interface {
	DecodeValue(Decoder) error
}

// ValueDecoderFunc allows the use of regular functions or methods as value
// decoders.
type ValueDecoderFunc func(Decoder) error

// DecodeValue calls f(d).
func (f ValueDecoderFunc) DecodeValue(d Decoder) error { return f(d) }

type decodeFuncOpts struct {
	recurse bool
	structs map[reflect.Type]*structType
}

type decodeFunc func(Decoder, reflect.Value) (Type, error)

func decodeFuncOf(t reflect.Type) decodeFunc {
	return makeDecodeFunc(t, decodeFuncOpts{})
}

func makeDecodeFunc(t reflect.Type, opts decodeFuncOpts) decodeFunc {
	if a, ok := AdapterOf(t); ok {
		decode := a.Decode
		return func(d Decoder, v reflect.Value) (Type, error) {
			err := decode(d, v)
			return Unknown /* just needs to not be Nil */, err
		}
	}

	// fast path: check if it's a basic go type
	switch t {
	case boolType:
		return Decoder.decodeBool

	case stringType:
		return Decoder.decodeString

	case bytesType:
		return Decoder.decodeBytes

	case timeType:
		return Decoder.decodeTime

	case durationType:
		return Decoder.decodeDuration

	case emptyInterface:
		return Decoder.decodeInterface

	case intType, int8Type, int16Type, int32Type, int64Type:
		return Decoder.decodeInt

	case uintType, uint8Type, uint16Type, uint32Type, uint64Type, uintptrType:
		return Decoder.decodeUint

	case float32Type, float64Type:
		return Decoder.decodeFloat
	}

	// check if it implements one of the special case interfaces, first on the
	// plain type, then on the pointer type
	switch {
	case t.Implements(valueDecoderInterface):
		return Decoder.decodeDecoder

	case t.Implements(errorInterface):
		return Decoder.decodeError

	case t.Implements(binaryUnmarshalerInterface) && t.Implements(textUnmarshalerInterface):
		return Decoder.decodeUnmarshaler

	case t.Implements(binaryUnmarshalerInterface):
		return Decoder.decodeBinaryUnmarshaler

	case t.Implements(textUnmarshalerInterface):
		return Decoder.decodeTextUnmarshaler
	}

	switch p := reflect.PtrTo(t); {
	case p.Implements(valueDecoderInterface):
		return Decoder.decodeDecoderPointer

	case p.Implements(binaryUnmarshalerInterface) && p.Implements(textUnmarshalerInterface):
		return Decoder.decodeUnmarshalerPointer

	case p.Implements(binaryUnmarshalerInterface):
		return Decoder.decodeBinaryUnmarshalerPointer

	case p.Implements(textUnmarshalerInterface):
		return Decoder.decodeTextUnmarshalerPointer
	}

	// check what kind is the type, potentially generate a decoder
	switch t.Kind() {
	case reflect.Struct:
		return makeDecodeStructFunc(t, opts)

	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			return Decoder.decodeBytes
		}
		return makeDecodeSliceFunc(t, opts)

	case reflect.Map:
		return makeDecodeMapFunc(t, opts)

	case reflect.Ptr:
		return makeDecodePtrFunc(t, opts)

	case reflect.Array:
		return makeDecodeArrayFunc(t, opts)

	case reflect.Bool:
		return Decoder.decodeBool

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Decoder.decodeInt

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return Decoder.decodeUint

	case reflect.Float32, reflect.Float64:
		return Decoder.decodeFloat

	case reflect.String:
		return Decoder.decodeString

	default:
		return Decoder.decodeUnsupported
	}
}

func makeDecodeSliceFunc(t reflect.Type, opts decodeFuncOpts) decodeFunc {
	if !opts.recurse {
		return Decoder.decodeSlice
	}
	f := makeDecodeFunc(t.Elem(), opts)
	return func(d Decoder, v reflect.Value) (Type, error) {
		return d.decodeSliceWith(v, f)
	}
}

func makeDecodeArrayFunc(t reflect.Type, opts decodeFuncOpts) decodeFunc {
	if !opts.recurse {
		return Decoder.decodeArray
	}
	f := makeDecodeFunc(t.Elem(), opts)
	return func(d Decoder, v reflect.Value) (Type, error) {
		return d.decodeArrayWith(v, f)
	}
}

func makeDecodeMapFunc(t reflect.Type, opts decodeFuncOpts) decodeFunc {
	if !opts.recurse {
		return Decoder.decodeMap
	}
	kf := makeDecodeFunc(t.Key(), opts)
	vf := makeDecodeFunc(t.Elem(), opts)
	return func(d Decoder, v reflect.Value) (Type, error) {
		return d.decodeMapWith(v, kf, vf)
	}
}

func makeDecodeStructFunc(t reflect.Type, opts decodeFuncOpts) decodeFunc {
	if !opts.recurse {
		return Decoder.decodeStruct
	}
	s := newStructType(t, opts.structs)
	return func(d Decoder, v reflect.Value) (Type, error) {
		return d.decodeStructWith(v, s)
	}
}

func makeDecodePtrFunc(t reflect.Type, opts decodeFuncOpts) decodeFunc {
	if !opts.recurse {
		return Decoder.decodePointer
	}
	f := makeDecodeFunc(t.Elem(), opts)
	return func(d Decoder, v reflect.Value) (Type, error) {
		return d.decodePointerWith(v, f)
	}
}
