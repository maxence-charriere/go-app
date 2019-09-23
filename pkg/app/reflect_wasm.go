package app

import (
	"encoding/json"
	"errors"
	"net/url"
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

func mapCompoFields(c Compo, fields map[string]string) error {
	v := reflect.ValueOf(c).Elem()
	t := v.Type()

	for i, numfields := 0, t.NumField(); i < numfields; i++ {
		fv := v.Field(i)
		ft := t.Field(i)

		if ft.Anonymous {
			continue
		}

		// Ignore non exported field.
		if len(ft.PkgPath) != 0 {
			continue
		}

		name := strings.ToLower(ft.Name)
		value, ok := fields[name]

		// Remove not set boolean.
		if !ok && fv.Kind() == reflect.Bool {
			fv.SetBool(false)
			continue
		} else if !ok {
			continue
		}

		if err := mapCompoField(fv, value); err != nil {
			return err
		}
	}
	return nil
}

func mapCompoField(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)

	case reflect.Bool:
		if len(value) == 0 {
			value = "true"
		}
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(b)

	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		n, err := strconv.ParseInt(value, 0, 64)
		if err != nil {
			return err
		}
		field.SetInt(n)

	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uintptr:
		n, err := strconv.ParseUint(value, 0, 64)
		if err != nil {
			return err
		}
		field.SetUint(n)

	case reflect.Float64, reflect.Float32:
		n, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(n)

	default:
		var err error
		if value, err = url.QueryUnescape(value); err != nil {
			return err
		}
		addr := field.Addr()
		i := addr.Interface()
		if err = json.Unmarshal([]byte(value), i); err != nil {
			return err
		}
	}
	return nil
}

func mapCompoFieldFromURLQuery(c Compo, query url.Values) error {
	attrs := make(map[string]string, len(query))
	for k := range query {
		v := query.Get(k)
		k = strings.ToLower(k)
		attrs[k] = v
	}
	return mapCompoFields(c, attrs)
}
