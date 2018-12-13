package dom

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	"github.com/murlokswarm/app"
)

type compo struct {
	ID       string
	ParentID string
	Compo    app.Compo
	Events   app.Subscriber
}

func mapCompoFields(c app.Compo, fields map[string]string) error {
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
		addr := field.Addr()
		i := addr.Interface()
		if err := json.Unmarshal([]byte(value), i); err != nil {
			return err
		}
	}
	return nil
}
