package template

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

type Encoder func(interface{}) ([]byte, error)

func Gen(encode Encoder, tmpl string, data interface{}) ([]byte, error) {
	md, err := mapData(encode, data)
	if err != nil {
		return nil, errors.New("mapping data failed").Wrap(err)
	}

	for k, v := range md {
		tmpl = strings.ReplaceAll(tmpl, k, v)
	}

	return []byte(tmpl), nil
}

func mapData(encode Encoder, data interface{}) (map[string]string, error) {
	v := reflect.Indirect(reflect.ValueOf(data))
	t := v.Type()
	switch v.Kind() {
	case reflect.Struct, reflect.Map:

	default:
		return nil, errors.New("data is not a struct or map")
	}

	m := make(map[string]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}

		val, err := encode(v.Field(i).Interface())
		if err != nil {
			return nil, errors.New("encoding value failed").
				Tag("key", f.Name).
				Wrap(err)
		}

		key := fmt.Sprintf(`"T(%s)"`, f.Name)
		m[key] = string(val)
	}
	return m, nil
}
