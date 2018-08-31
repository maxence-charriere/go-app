package dom

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
)

type compo struct {
	id      string
	compoID string
	name    string
	fields  map[string]string
	parent  node
	root    node
	changes []Change
}

func newCompo(name string, fields map[string]string) *compo {
	c := &compo{
		id:     name + ":" + uuid.New().String(),
		name:   name,
		fields: fields,
	}

	c.changes = append(c.changes, createCompoChange(c.id, name))
	return c
}

func (c *compo) ID() string {
	return c.id
}

func (c *compo) CompoID() string {
	return c.compoID
}

func (c *compo) Parent() app.Node {
	return c.parent
}

func (c *compo) SetParent(p node) {
	c.parent = p
}

func (c *compo) SetRoot(root node) {
	root.SetParent(c)
	c.root = root
	c.changes = append(c.changes, setCompoRootChange(c.id, root.ID()))
}

func (c *compo) RemoveRoot() {
	c.root.Close()
	c.changes = append(c.changes, c.root.Flush()...)
	c.root = nil
}

func (c *compo) Flush() []Change {
	changes := make([]Change, 0, len(c.changes))

	if c.root != nil {
		changes = append(changes, c.root.Flush()...)
	}

	changes = append(changes, c.changes...)
	c.changes = c.changes[:0]
	return changes
}

func (c *compo) Close() {
	if c.root != nil {
		c.root.Close()
		c.changes = append(c.changes, c.root.Flush()...)
	}

	c.SetParent(nil)
	c.changes = append(c.changes, deleteNodeChange(c.id))
}

func validateCompo(c app.Compo) error {
	v := reflect.ValueOf(c)
	if v.Kind() != reflect.Ptr {
		return errors.New("compo is not a pointer")
	}

	v = v.Elem()
	if v.NumField() == 0 {
		return errors.New("compo is based on a struct without field. use app.ZeroCompo instead of struct{}")
	}
	return nil
}

func decodeCompo(c app.Compo, t ...Transform) (node, error) {
	var funcs template.FuncMap

	if compoExtRend, ok := c.(app.CompoWithExtendedRender); ok {
		funcs = compoExtRend.Funcs()
	}

	if len(funcs) == 0 {
		funcs = make(template.FuncMap, 4)
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

	tmpl, err := template.
		New("").
		Funcs(funcs).
		Parse(c.Render())
	if err != nil {
		return nil, err
	}

	var w bytes.Buffer
	if err = tmpl.Execute(&w, c); err != nil {
		return nil, err
	}

	return decodeNodes(w.String(), t...)
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
