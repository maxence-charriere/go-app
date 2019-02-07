package app

import (
	"encoding/json"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// Compo is the interface that describes a component.
// Must be implemented on a non empty struct pointer.
type Compo interface {
	// Render must return HTML 5.
	// It supports standard Go html/template API.
	// The pipeline is based on the component struct.
	// See https://golang.org/pkg/text/template and
	// https://golang.org/pkg/html/template for template usage.
	Render() string
}

// Mounter is the interface that wraps OnMount method.
type Mounter interface {
	Compo

	// OnMount is called when a component is mounted.
	// App.Render should not be called inside.
	OnMount()
}

// Dismounter is the interface that wraps OnDismount method.
type Dismounter interface {
	Compo

	// OnDismount is called when a component is dismounted.
	// App.Render should not be called inside.
	OnDismount()
}

// Navigable is the interface that wraps OnNavigate method.
type Navigable interface {
	Compo

	// OnNavigate is called when a component is loaded or navigated to.
	// It is called just after the component is mounted.
	OnNavigate(u *url.URL)
}

// EventSubscriber is the interface that describes a component that subscribes
// to events emitted from messages.
type EventSubscriber interface {
	// Subscribe is called when a component is mounted.
	// The returned subscriber is used to subscribe to events emitted from
	// messages.
	// All the event subscribed are automatically unsuscribed when the component
	// is dismounted.
	Subscribe() *Subscriber
}

// CompoWithExtendedRender is the interface that wraps Funcs method.
type CompoWithExtendedRender interface {
	Compo

	// Funcs returns a map of funcs to use when rendering a component.
	// Funcs named raw, json and time are reserved.
	// They handle raw html code, json conversions and time format.
	// They can't be overloaded.
	// See https://golang.org/pkg/text/template/#Template.Funcs for more details.
	Funcs() map[string]interface{}
}

// ZeroCompo is the type to use as base for empty components.
// Every instances of an empty struct is given the same memory address, which
// causes problem for indexing components.
// ZeroCompo have a placeholder field to avoid that.
type ZeroCompo struct {
	placeholder byte
}

type compo struct {
	ID       string
	ParentID string
	Compo    Compo
	Events   *Subscriber
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
