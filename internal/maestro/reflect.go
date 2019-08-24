package maestro

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

func getReceiver(v interface{}, target string) (reflect.Value, error) {
	return getReceiverFromValue(reflect.ValueOf(v), target)
}

func getReceiverFromValue(v reflect.Value, target string) (reflect.Value, error) {
	switch v.Kind() {
	case reflect.Ptr:
		return getReceiverFromPtr(v, target)

	case reflect.Struct:
		return getReceiverFromStruct(v, target)

	case reflect.Map:
		return getReceiverFromMap(v, target)

	case reflect.Slice, reflect.Array:
		return getReceiverFromSlice(v, target)

	case reflect.Interface:
		return getReceiverFromValue(v.Elem(), target)

	default:
		return getReveiverFromDefault(v, target)
	}
}

func getReceiverFromPtr(v reflect.Value, target string) (reflect.Value, error) {
	method, nextTarget := parseTarget(target)

	if method != "" && nextTarget == "" && isExported(method) {
		res := v.MethodByName(method)
		if res.IsValid() {
			return res, nil
		}
	}

	return getReceiverFromValue(v.Elem(), target)
}

func getReceiverFromStruct(v reflect.Value, target string) (reflect.Value, error) {
	target, nextTarget := parseTarget(target)
	if !isExported(target) {
		return reflect.Value{}, errors.New("non exported field or method named " + target)
	}

	res := v.MethodByName(target)
	if !res.IsValid() {
		res = v.FieldByName(target)
	}
	if !res.IsValid() {
		return reflect.Value{}, errors.New("no field or method named " + target)
	}

	if nextTarget == "" {
		return res, nil
	}
	return getReceiverFromValue(res, nextTarget)
}

func getReceiverFromMap(v reflect.Value, target string) (reflect.Value, error) {
	target, nextTarget := parseTarget(target)

	res := v.MethodByName(target)
	if !res.IsValid() || !isExported(target) {
		res = v.MapIndex(reflect.ValueOf(target))
	}
	if !res.IsValid() {
		return reflect.Value{}, errors.New("no key or non exported method named " + target)
	}

	if nextTarget == "" {
		return res, nil
	}
	return getReceiverFromValue(res, nextTarget)
}

func getReceiverFromSlice(v reflect.Value, target string) (reflect.Value, error) {
	target, nextTarget := parseTarget(target)

	res := v.MethodByName(target)
	if !res.IsValid() || !isExported(target) {
		idx, err := strconv.Atoi(target)
		if err != nil {
			return reflect.Value{}, err
		}
		if idx < v.Len() {
			res = v.Index(idx)
		}
	}
	if !res.IsValid() {
		return reflect.Value{}, errors.New("out of range index or non exported method named " + target)
	}

	if nextTarget == "" {
		return res, nil
	}
	return getReceiverFromValue(res, nextTarget)
}

func getReveiverFromDefault(v reflect.Value, target string) (reflect.Value, error) {
	target, nextTarget := parseTarget(target)
	if nextTarget != "" {
		return reflect.Value{}, errors.New(v.Type().String() + "can't contain other value")
	}
	if target == "" {
		return v, nil
	}

	res := v.MethodByName(target)
	if !res.IsValid() || !isExported(target) {
		return reflect.Value{}, errors.New("invalid value kind or non exported method name " + target)
	}
	return res, nil
}

func parseTarget(target string) (current, next string) {
	end := strings.IndexRune(target, '.')
	if end != -1 {
		return target[:end], target[end+1:]
	}
	return target, ""
}

func isExported(fieldOrMethod string) bool {
	return !unicode.IsLower(rune(fieldOrMethod[0]))
}
