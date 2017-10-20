package html

import (
	"bytes"
	"encoding/json"
	"html/template"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/pkg/errors"
)

// Markup implements the app.Markup interface.
type Markup struct {
	components map[uuid.UUID]app.Component
	roots      map[app.Component]app.Tag
	factory    app.Factory
}

// NewMarkup creates a markup with the given factory.
func NewMarkup(factory app.Factory) *Markup {
	return &Markup{
		components: make(map[uuid.UUID]app.Component),
		roots:      make(map[app.Component]app.Tag),
		factory:    factory,
	}
}

// Component satisfies the app.Markup interface.
func (m *Markup) Component(id uuid.UUID) (compo app.Component, err error) {
	var ok bool
	if compo, ok = m.components[id]; !ok {
		err = errors.New("component not mounted")
	}
	return
}

// Contains satisfies the app.Markup interface.
func (m *Markup) Contains(compo app.Component) bool {
	_, ok := m.roots[compo]
	return ok
}

// Root satisfies the app.Markup interface.
func (m *Markup) Root(compo app.Component) (root app.Tag, err error) {
	var ok bool
	if root, ok = m.roots[compo]; !ok {
		err = errors.New("component not mounted")
	}
	return
}

func (m *Markup) Mount(compo app.Component) (root app.Tag, err error) {
	return m.mount(compo, uuid.New())
}

func (m *Markup) mount(compo app.Component, compoID uuid.UUID) (tag app.Tag, err error) {
	if m.Contains(compo) {
		err = errors.New("component is already mounted")
		return
	}

	if err = decodeComponent(compo, &tag); err != nil {
		err = errors.Wrap(err, "decoding component failed")
		return
	}

	if err = m.mountTag(&tag, compoID); err != nil {
		return
	}

	m.components[compoID] = compo
	m.roots[compo] = tag

	if mounter, ok := compo.(app.Mounter); ok {
		mounter.OnMount()
	}
	return
}

func decodeComponent(compo app.Component, tag *app.Tag) error {
	var funcs template.FuncMap
	if compoExtRend, ok := compo.(app.ComponentWithExtendedRender); ok {
		funcs = compoExtRend.Funcs()
	} else {
		funcs = make(template.FuncMap, 3)
	}

	funcs["raw"] = func(s string) template.HTML {
		return template.HTML(s)
	}

	funcs["time"] = func(t time.Time, layout string) string {
		return t.Format(layout)
	}

	funcs["json"] = func(v interface{}) string {
		b, _ := json.Marshal(v)
		return string(b)
	}

	rendering := compo.Render()
	tmpl := template.Must(template.New("").Funcs(funcs).Parse(rendering))

	var buff bytes.Buffer
	if err := tmpl.Execute(&buff, compo); err != nil {
		return err
	}

	dec := NewDecoder(&buff)
	return dec.Decode(tag)
}

func (m *Markup) mountTag(tag *app.Tag, compoID uuid.UUID) error {
	tag.ID = uuid.New()
	tag.CompoID = compoID

	if tag.Is(app.TextTag) {
		return nil
	}

	if tag.Is(app.CompoTag) {
		compo, err := m.factory.NewComponent(tag.Name)
		if err != nil {
			return err
		}

		if err = mapComponentFields(compo, tag.Attributes); err != nil {
			return err
		}

		_, err = m.mount(compo, tag.ID)
		return err
	}

	for i := range tag.Children {
		if err := m.mountTag(&tag.Children[i], compoID); err != nil {
			return errors.Wrap(err, "mounting children failed")
		}
	}
	return nil
}

func mapComponentFields(compo app.Component, attrs app.AttributeMap) error {
	if len(attrs) == 0 {
		return nil
	}

	val := reflect.ValueOf(compo).Elem()
	typ := val.Type()

	for i, numfields := 0, typ.NumField(); i < numfields; i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		if field.Anonymous {
			continue
		}

		if len(field.PkgPath) != 0 {
			continue
		}

		attrName := strings.ToLower(field.Name)
		attrVal, ok := attrs[attrName]

		if !ok {
			if fieldVal.Kind() == reflect.Bool {
				fieldVal.SetBool(false)
			}
			continue
		}

		if err := mapComponentField(fieldVal, attrVal); err != nil {
			return errors.Wrapf(err, "mapping attribute %s to field %s failed", attrName, field.Name)
		}
	}
	return nil
}

func mapComponentField(field reflect.Value, attr string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(attr)

	case reflect.Bool:
		if len(attr) == 0 {
			attr = "true"
		}
		b, err := strconv.ParseBool(attr)
		if err != nil {
			return err
		}
		field.SetBool(b)

	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		n, err := strconv.ParseInt(attr, 0, 64)
		if err != nil {
			return err
		}
		field.SetInt(n)

	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uintptr:
		n, err := strconv.ParseUint(attr, 0, 64)
		if err != nil {
			return err
		}
		field.SetUint(n)

	case reflect.Float64, reflect.Float32:
		n, err := strconv.ParseFloat(attr, 64)
		if err != nil {
			return err
		}
		field.SetFloat(n)

	default:
		addr := field.Addr()
		i := addr.Interface()
		if err := json.Unmarshal([]byte(attr), i); err != nil {
			return err
		}
	}
	return nil
}

func (m *Markup) Dismount(compo Component) {

}

func (m *Markup) Update(compo app.Component) (syncs []app.TagSync, err error) {
	panic("not implemented")
}
