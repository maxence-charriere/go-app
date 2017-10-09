package html

import (
	"bytes"
	"encoding/json"
	"html/template"
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
		components: make(map[uuid.Component]app.Component),
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
		return errors.New("component not mounted")
	}
}

// Mount satisfies the app.Markup interface.
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

	if err = m.mountTag(tag, compoID); err != nil {
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

	switch tag.Type {
	case app.TextTag:
		return nil
	case app.CompoTag:

	case app.SimpleTag:
	default:
		return errors.Errorf("tag named %s: type %v is not supported", tag.Name, tag.Type)
	}
}

// Dismount satisfies the app.Markup interface.
func (m *Markup) Dismount(compo app.Component) {
	panic("not implemented")
}

// Update satisfies the app.Markup interface.
func (m *Markup) Update(compo app.Component) (syncs []app.TagSync, err error) {
	panic("not implemented")
}
