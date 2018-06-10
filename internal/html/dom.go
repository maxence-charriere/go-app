package html

import (
	"bytes"
	"encoding/json"
	"html/template"
	"time"

	"github.com/murlokswarm/app"
)

// NewDOM create a document object model store.
func NewDOM(f app.Factory, controlID string) app.DOM {
	return &dom{
		factory:   f,
		controlID: controlID,
	}
}

type dom struct {
	factory       app.Factory
	controlID     string
	rootComponent app.Component
	components    map[string]app.Component
	compoData     map[app.Component]compoData
}

// ComponentByID returns the component with the given identifier.
// It satisfies the app.DOM interface.
func (d *dom) ComponentByID(id string) (app.Component, error) {
	c, ok := d.components[id]
	if !ok {
		return nil, app.ErrNotFound
	}
	return c, nil
}

func (d *dom) insertComponent(id string, c app.Component) {
}

func (d *dom) deleteComponent(id string) {
	if c, err := d.ComponentByID(id); err == nil {
		delete(d.compoData, c)
	}
	delete(d.components, id)
}

// ContainsComponent reports whether the given component is in the dom.
func (d *dom) ContainsComponent(c app.Component) bool {
	_, ok := d.compoData[c]
	return ok
}

// Render create or update the given component.
// It satisfies the app.DOM interface.
func (d *dom) Render(c app.Component) ([]app.DOMChange, error) {
	panic("not implemented")
}

func (d *dom) mountComponent(c app.Component) ([]app.DOMChange, error) {
	panic("not implemented")
}

func (d *dom) mountNode(n app.DOMNode, compoID string) ([]app.DOMChange, error) {
	panic("not implemented")
}

func (d *dom) syncNodes(current, new node) ([]app.DOMChange, error) {
	panic("not implemented")
}

type compoData struct {
	root   node
	events app.EventSubscriber
}

func decodeComponent(c app.Component) (node, error) {
	var funcs template.FuncMap

	if compoExtRend, ok := c.(app.ComponentWithExtendedRender); ok {
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
	return decodeNodes(w.String())
}
