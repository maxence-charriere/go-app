package mold

import (
	"reflect"
)

// extractType gets the actual underlying type of field value.
func extractType(current reflect.Value) (reflect.Value, reflect.Kind) {
	switch current.Kind() {
	case reflect.Ptr:
		if current.IsNil() {
			return current, reflect.Ptr
		}
		return extractType(current.Elem())

	case reflect.Interface:
		if current.IsNil() {
			return current, reflect.Interface
		}
		return extractType(current.Elem())

	default:
		return current, current.Kind()
	}
}

// HasValue determines if a reflect.Value is it's default value
func HasValue(field reflect.Value) bool {
	switch field.Kind() {
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return !field.IsNil()
	default:
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}
