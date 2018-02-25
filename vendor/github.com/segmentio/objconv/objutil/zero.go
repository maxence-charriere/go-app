package objutil

import (
	"reflect"
	"unsafe"
)

// IsZero returns true if the value given as argument is the zero-value of
// the type of v.
func IsZero(v interface{}) bool {
	return IsZeroValue(reflect.ValueOf(v))
}

func IsZeroValue(v reflect.Value) bool {
	if !v.IsValid() {
		return true // nil interface{}
	}
	switch v.Kind() {
	case reflect.Map, reflect.Slice, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.String:
		return v.Len() == 0
	case reflect.UnsafePointer:
		return unsafe.Pointer(v.Pointer()) == nil
	case reflect.Array:
		return isZeroArray(v)
	case reflect.Struct:
		return isZeroStruct(v)
	}
	return false
}

func isZeroArray(v reflect.Value) bool {
	for i, n := 0, v.Len(); i != n; i++ {
		if !IsZeroValue(v.Index(i)) {
			return false
		}
	}
	return true
}

func isZeroStruct(v reflect.Value) bool {
	for i, n := 0, v.NumField(); i != n; i++ {
		if !IsZeroValue(v.Field(i)) {
			return false
		}
	}
	return true
}
