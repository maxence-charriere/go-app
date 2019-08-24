package maestro

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"
)

// Compo is the interface that describes a component.
type Compo interface {
	// Render must return a HTML5 string. It supports standard Go html/template
	// API. The pipeline is based on the component struct.
	// See https://golang.org/pkg/text/template and
	// https://golang.org/pkg/html/template for template usage.
	Render() string
}

type mounter interface {
	OnMount()
}

type dismounter interface {
	OnDismount()
}

type compoWithExtendedRender interface {
	Funcs() map[string]interface{}
}

// CompoBuilder is a factory that can create components.
type CompoBuilder map[string]reflect.Type

// Import creates a builder that can build the given component.
func (b CompoBuilder) Import(c Compo) error {
	v := reflect.ValueOf(c)
	if v.Kind() != reflect.Ptr {
		return errors.New("component is not a pointer")
	}
	if v = v.Elem(); v.Kind() != reflect.Struct {
		return errors.New("component is not implemented on a struct")
	}
	if v.NumField() == 0 {
		return errors.New("component is based on a struct without field. use ZeroCompo instead of struct{}")
	}

	b[compoName(c)] = v.Type()
	return nil
}

// IsImported report wether a component has been imported.
func (b CompoBuilder) IsImported(name string) bool {
	_, ok := b[name]
	return ok
}

// New creates the named component.
func (b CompoBuilder) New(name string) (Compo, error) {
	t, ok := b[name]
	if !ok {
		return nil, errors.New("component " + name + " is not imported")
	}
	return reflect.New(t).Interface().(Compo), nil
}

func compoName(c Compo) string {
	v := reflect.ValueOf(c)
	v = reflect.Indirect(v)
	name := strings.ToLower(v.Type().String())
	return strings.TrimPrefix(name, "main.")
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
		addr := field.Addr()
		i := addr.Interface()
		if err := json.Unmarshal([]byte(value), i); err != nil {
			return err
		}
	}
	return nil
}
