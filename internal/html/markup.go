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
	components       map[string]app.Compo
	roots            map[app.Compo]*app.Tag
	eventSubscribers map[app.Compo]*app.EventSubscriber
	factory          *app.Factory
}

// NewMarkup creates a markup with the given factory.
func NewMarkup(factory *app.Factory) *Markup {
	return &Markup{
		components:       make(map[string]app.Compo),
		roots:            make(map[app.Compo]*app.Tag),
		eventSubscribers: make(map[app.Compo]*app.EventSubscriber),
		factory:          factory,
	}
}

// Len satisfies the app.Markup interface.
func (m *Markup) Len() int {
	return len(m.components)
}

// Factory satisfies the app.Markup interface.
func (m *Markup) Factory() *app.Factory {
	return m.factory
}

// Compo satisfies the app.Markup interface.
func (m *Markup) Compo(id string) (compo app.Compo, err error) {
	var ok bool
	if compo, ok = m.components[id]; !ok {
		err = errors.New("component not mounted")
	}
	return
}

// Contains satisfies the app.Markup interface.
func (m *Markup) Contains(compo app.Compo) bool {
	_, ok := m.roots[compo]
	return ok
}

// Root satisfies the app.Markup interface.
func (m *Markup) Root(compo app.Compo) (root app.Tag, err error) {
	rootPtr, ok := m.roots[compo]
	if !ok {
		err = errors.New("component not mounted")
		return
	}

	root = *rootPtr
	return
}

// FullRoot satisfies the app.Markup interface.
func (m *Markup) FullRoot(tag app.Tag) (root app.Tag, err error) {
	root = tag
	root.Children = make([]app.Tag, len(tag.Children))
	copy(root.Children, tag.Children)

	for i, child := range tag.Children {
		if !child.Is(app.CompoTag) {
			continue
		}

		var compo app.Compo
		if compo, err = m.Compo(child.ID); err != nil {
			return root, err
		}

		// The err checking is ignored here because the err would be the same as
		// m.Compo call.
		child, _ = m.Root(compo)

		if child, err = m.FullRoot(child); err != nil {
			return root, err
		}

		root.Children[i] = child
	}

	return root, nil
}

// Mount satisfies the app.Markup interface.
func (m *Markup) Mount(compo app.Compo) (root app.Tag, err error) {
	return m.mount(compo, uuid.New().String())
}

func (m *Markup) mount(compo app.Compo, compoID string) (root app.Tag, err error) {
	if m.Contains(compo) {
		err = errors.New("component is already mounted")
		return
	}

	if err = decodeCompo(compo, &root); err != nil {
		return
	}

	if err = m.mountTag(&root, uuid.New().String(), compoID); err != nil {
		return
	}

	m.components[compoID] = compo
	m.roots[compo] = &root

	if mounter, ok := compo.(app.Mounter); ok {
		mounter.OnMount()
	}

	if subscriber, ok := compo.(app.Subscriber); ok {
		m.eventSubscribers[compo] = subscriber.Subscribe()
	}
	return
}

func decodeCompo(compo app.Compo, tag *app.Tag) error {
	var funcs template.FuncMap
	if compoExtRend, ok := compo.(app.CompoWithExtendedRender); ok {
		funcs = compoExtRend.Funcs()
	} else {
		funcs = make(template.FuncMap, 3)
	}

	funcs["raw"] = func(s string) template.HTML {
		return template.HTML(s)
	}

	funcs["compo"] = func(s string) template.HTML {
		return template.HTML("<" + s + ">")
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

func (m *Markup) mountTag(tag *app.Tag, id string, compoID string) error {
	tag.ID = id
	tag.CompoID = compoID

	if tag.Is(app.TextTag) {
		return nil
	}

	if tag.Is(app.CompoTag) {
		compo, err := m.factory.NewCompo(tag.Name)
		if err != nil {
			return err
		}

		if err = mapCompoFields(compo, tag.Attributes); err != nil {
			return err
		}

		_, err = m.mount(compo, tag.ID)
		return err
	}

	for i := range tag.Children {
		if err := m.mountTag(&tag.Children[i], uuid.New().String(), compoID); err != nil {
			return err
		}
	}
	return nil
}

func mapCompoFields(compo app.Compo, attrs app.AttributeMap) error {
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

		if err := mapCompoField(fieldVal, attrVal); err != nil {
			return err
		}
	}
	return nil
}

func mapCompoField(field reflect.Value, attr string) error {
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
func (m *Markup) Dismount(compo app.Compo) {
	root, ok := m.roots[compo]
	if !ok {
		return
	}

	m.dismountTag(*root)
	delete(m.components, root.CompoID)
	delete(m.roots, compo)

	if dismounter, ok := compo.(app.Dismounter); ok {
		dismounter.OnDismount()
	}

	if _, ok := compo.(app.Subscriber); ok {
		subscriber := m.eventSubscribers[compo]
		subscriber.Close()
		delete(m.eventSubscribers, compo)
	}
}

func (m *Markup) dismountTag(tag app.Tag) {
	if tag.Is(app.CompoTag) {
		// Sub component are registered under the id of the tag that targets
		// them.
		compo, err := m.Compo(tag.ID)
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
func (m *Markup) Update(compo app.Compo) (syncs []app.TagSync, err error) {
	syncs, _, err = m.update(compo)
	return
}

func (m *Markup) update(compo app.Compo) (syncs []app.TagSync, replaceParent bool, err error) {
	root, ok := m.roots[compo]
	if !ok {
		err = errors.New("component not mounted")
		return
	}

	var newRoot app.Tag
	if err = decodeCompo(compo, &newRoot); err != nil {
		return
	}

	syncs, replaceParent, err = m.syncTags(root, &newRoot)
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
		return m.syncCompoTags(current, new)
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

func (m *Markup) syncCompoTags(current, new *app.Tag) (syncs []app.TagSync, replaceParent bool, err error) {
	if attributesEquals(current.Name, current.Attributes, new.Attributes) {
		return
	}

	current.Attributes = new.Attributes

	var compo app.Compo
	if compo, err = m.Compo(current.ID); err != nil {
		return
	}

	if err = mapCompoFields(compo, current.Attributes); err != nil {
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
		if err = m.mountTag(child, uuid.New().String(), current.CompoID); err != nil {
			return
		}
		current.Children = append(current.Children, *child)
		newChildren = newChildren[1:]
	}
	return
}

// Map satisfies the app.Markup interface.
func (m *Markup) Map(mapping app.Mapping) (function func(), err error) {
	var pipeline []string
	if pipeline, err = app.ParseMappingTarget(mapping.Target); err != nil {
		return
	}

	var compo app.Compo
	if compo, err = m.Compo(mapping.CompoID); err != nil {
		return
	}

	mapper := newMapper(pipeline, mapping.JSONValue)
	function, err = mapper.MapTo(compo)
	return
}
