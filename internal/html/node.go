package html

import "github.com/murlokswarm/app"

type node interface {
	app.DOMNode
	SetParent(app.DOMNode)
}

type textNode struct {
	id        string
	compoID   string
	controlID string
	text      string
	parent    app.DOMNode
}

func (t *textNode) ID() string {
	return t.id
}

func (t *textNode) CompoID() string {
	return t.compoID
}

func (t *textNode) ControlID() string {
	return t.controlID
}

func (t *textNode) Text() string {
	return t.text
}

func (t *textNode) Parent() app.DOMNode {
	return t.parent
}

func (t *textNode) SetParent(p app.DOMNode) {
	t.parent = p
}

type elemNode struct {
	id        string
	compoID   string
	controlID string
	tagName   string
	attrs     map[string]string
	parent    app.DOMNode
	children  []app.DOMNode
}

func (e *elemNode) ID() string {
	return e.id
}

func (e *elemNode) CompoID() string {
	return e.compoID
}

func (e *elemNode) ControlID() string {
	return e.controlID
}

func (e *elemNode) TagName() string {
	return e.tagName
}

func (e *elemNode) Attrs() map[string]string {
	return e.attrs
}

func (e *elemNode) Parent() app.DOMNode {
	return e.parent
}

func (e *elemNode) SetParent(p app.DOMNode) {
	e.parent = p
}

func (e *elemNode) Children() []app.DOMNode {
	return e.children
}

func (e *elemNode) appendChild(c node) {
	e.children = append(e.children, c)
	c.SetParent(e)
}

type compoNode struct {
	id        string
	compoID   string
	controlID string
	name      string
	fields    map[string]string
	parent    app.DOMNode
}

func (c *compoNode) ID() string {
	return c.id
}

func (c *compoNode) CompoID() string {
	return c.compoID
}

func (c *compoNode) ControlID() string {
	return c.controlID
}

func (c *compoNode) Name() string {
	return c.name
}

func (c *compoNode) Fields() map[string]string {
	return c.fields
}

func (c *compoNode) Parent() app.DOMNode {
	return c.parent
}

func (c *compoNode) SetParent(p app.DOMNode) {
	c.parent = p
}
