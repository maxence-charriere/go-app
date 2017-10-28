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

// Mount satisfies the app.Markup interface.
func (m *Markup) Mount(compo app.Component) (root app.Tag, err error) {
	return m.mount(compo, uuid.New())
}

func (m *Markup) mount(compo app.Component, compoID uuid.UUID) (root app.Tag, err error) {
	if m.Contains(compo) {
		err = errors.New("component is already mounted")
		return
	}

	if err = decodeComponent(compo, &root); err != nil {
		return
	}

	if err = m.mountTag(&root, uuid.New(), compoID); err != nil {
		return
	}

	m.components[compoID] = compo
	m.roots[compo] = root

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

func (m *Markup) mountTag(tag *app.Tag, id uuid.UUID, compoID uuid.UUID) error {
	tag.ID = id
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
		if err := m.mountTag(&tag.Children[i], uuid.New(), compoID); err != nil {
			return err
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
			return err
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

// Dismount satisfies the app.Markup interface.
func (m *Markup) Dismount(compo app.Component) {
	root, ok := m.roots[compo]
	if !ok {
		return
	}

	m.dismountTag(root)
	delete(m.components, root.CompoID)
	delete(m.roots, compo)

	if dismounter, ok := compo.(app.Dismounter); ok {
		dismounter.OnDismount()
	}
}

func (m *Markup) dismountTag(tag app.Tag) {
	if tag.Is(app.CompoTag) {
		// Sub component are registered under the id of the tag that targets
		// them.
		compo, err := m.Component(tag.ID)
		if err != nil {
			return
		}

		m.Dismount(compo)
		return
	}

	for _, child := range tag.Children {
		m.dismountTag(child)
	}
}

// Update satisfies the app.Markup interface.
func (m *Markup) Update(compo app.Component) (syncs []app.TagSync, err error) {
	syncs, _, err = m.update(compo)
	return
}

func (m *Markup) update(compo app.Component) (syncs []app.TagSync, replaceParent bool, err error) {
	var root app.Tag
	var newRoot app.Tag

	if root, err = m.Root(compo); err != nil {
		return
	}

	if err = decodeComponent(compo, &newRoot); err != nil {
		return
	}

	syncs, replaceParent, err = m.syncTags(&root, &newRoot)
	return
}

func (m *Markup) syncTags(current, new *app.Tag) (syncs []app.TagSync, replaceParent bool, err error) {
	if current.Name != new.Name {
		return m.mergeTags(current, new)
	}

	if current.Is(app.TextTag) {
		replaceParent = m.syncTextTags(current, new)
		return
	}

	if current.Is(app.CompoTag) {
		return m.syncComponentTags(current, new)
	}

	attrEquals := attributesEquals(current.Name, current.Attributes, new.Attributes)
	if !attrEquals {
		current.Attributes = new.Attributes
	}

	var replace bool
	var childSyncs []app.TagSync

	if childSyncs, replace, err = m.syncChildTags(current, new); err != nil {
		return
	}

	if replace {
		syncs = append(syncs, app.TagSync{
			Tag:     *current,
			Replace: true,
		})
		return
	}

	syncs = append(syncs, childSyncs...)

	if !attrEquals {
		syncs = append(syncs, app.TagSync{
			Tag: *current,
		})
	}
	return
}

func (m *Markup) mergeTags(current, new *app.Tag) (syncs []app.TagSync, replaceParent bool, err error) {
	m.dismountTag(*current)

	if err = m.mountTag(new, current.ID, current.CompoID); err != nil {
		return
	}

	*current = *new

	if current.Is(app.TextTag) {
		replaceParent = true
		return
	}

	syncs = append(syncs, app.TagSync{
		Tag:     *current,
		Replace: true,
	})
	return
}

func (m *Markup) syncTextTags(current, new *app.Tag) (replaceParent bool) {
	if current.Text != new.Text {
		current.Text = new.Text
		replaceParent = true
	}
	return
}

func (m *Markup) syncComponentTags(current, new *app.Tag) (syncs []app.TagSync, replaceParent bool, err error) {
	if attributesEquals(current.Name, current.Attributes, new.Attributes) {
		return
	}

	current.Attributes = new.Attributes

	var compo app.Component
	if compo, err = m.Component(current.ID); err != nil {
		return
	}

	if err = mapComponentFields(compo, current.Attributes); err != nil {
		return
	}

	syncs, replaceParent, err = m.update(compo)
	return
}

func attributesEquals(tagname string, current, new app.AttributeMap) bool {
	if len(current) != len(new) {
		return false
	}

	for name, val := range current {
		newVal, ok := new[name]
		if !ok {
			return false
		}
		if val != newVal {
			return false
		}
		if tagname == "input" && name == "value" && len(val) == 0 {
			return false
		}
	}
	return true
}

func (m *Markup) syncChildTags(current, new *app.Tag) (syncs []app.TagSync, replaceParent bool, err error) {
	curChildren := current.Children
	newChildren := new.Children
	i := 0

	if len(curChildren) != len(newChildren) {
		replaceParent = true
	}

	for len(curChildren) != 0 && len(newChildren) != 0 {
		var childSyncs []app.TagSync
		var replace bool

		if childSyncs, replace, err = m.syncTags(&curChildren[0], &newChildren[0]); err != nil {
			return
		}

		if replace {
			replaceParent = true
			syncs = nil
		}

		if !replaceParent {
			syncs = append(syncs, childSyncs...)
		}

		curChildren = curChildren[1:]
		newChildren = newChildren[1:]
		i++
	}

	current.Children = current.Children[:i]

	for len(curChildren) != 0 {
		m.dismountTag(curChildren[0])
		curChildren = curChildren[1:]
	}

	for len(newChildren) != 0 {
		child := &newChildren[0]
		if err = m.mountTag(child, uuid.New(), current.CompoID); err != nil {
			return
		}
		current.Children = append(current.Children, *child)
		newChildren = newChildren[1:]
	}
	return
}
