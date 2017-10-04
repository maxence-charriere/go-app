package markup

import (
	"net/url"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// CompoBuilder is the interface that describes a component factory.
type CompoBuilder interface {
	// Register registers component of type c into the builder.
	// Components must be registered to be used.
	// During a rendering, it allows to create components of same type as c when
	// a tag named like c is found.
	Register(c Component) error

	// New creates a component named n.
	New(n string) (Component, error)
}

// NewCompoBuilder creates a compo builder.
func NewCompoBuilder() CompoBuilder {
	return make(compoBuilder)
}

type compoBuilder map[string]reflect.Type

func (b compoBuilder) Register(c Component) error {
	if err := ensureValidComponent(c); err != nil {
		return err
	}

	v := reflect.ValueOf(c).Elem()
	t := v.Type()

	name := normalizeCompoName(t.String())
	b[name] = t
	return nil
}

func (b compoBuilder) New(name string) (c Component, err error) {
	t, ok := b[name]
	if !ok {
		err = errors.Errorf("component %s is not registered", name)
		return
	}
	v := reflect.New(t)
	c = v.Interface().(Component)
	return
}

func ensureValidComponent(c Component) error {
	v := reflect.ValueOf(c)
	if v.Kind() != reflect.Ptr {
		return errors.Errorf("%T must be implemented on a struct pointer", c)
	}

	if v = v.Elem(); v.Kind() != reflect.Struct {
		return errors.Errorf("%T must be implemented on a struct pointer", c)
	}

	if v.NumField() == 0 {
		return errors.Errorf("%T can't be implemented on an empty struct pointer", c)
	}
	return nil
}

func normalizeCompoName(name string) string {
	name = strings.ToLower(name)
	if pkgsep := strings.IndexByte(name, '.'); pkgsep != -1 {
		pkgname := name[:pkgsep]
		if pkgname == "main" {
			name = name[pkgsep+1:]
		}
	}
	return name
}

// ComponentNameFromURL returns the component name from URL.
// ok reports whether URL points a component.
func ComponentNameFromURL(u *url.URL) (name string, ok bool) {
	if len(u.Scheme) != 0 && u.Scheme != "component" {
		return
	}

	name = strings.TrimLeft(u.Path, "/")
	ok = len(name) != 0
	return
}
