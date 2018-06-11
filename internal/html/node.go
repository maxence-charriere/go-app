package html

import "github.com/murlokswarm/app"

type node interface {
	app.DOMNode
	setParent(node)
}

type textNode struct {
	id        string
	compoID   string
	controlID string
	text      string
	parent    node
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

func (t *textNode) setParent(p node) {
	t.parent = p
}

type elemNode struct {
	id        string
	compoID   string
	controlID string
	tagName   string
	attrs     map[string]string
	parent    node
	children  []node
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

func (e *elemNode) setParent(p node) {
	e.parent = p
}

func (e *elemNode) Children() []app.DOMNode {
	children := make([]app.DOMNode, len(e.children))
	for i, c := range e.children {
		children[i] = c
	}
	return children
}

func (e *elemNode) appendChild(c node) {
	e.children = append(e.children, c)
	c.setParent(e)
}

func (e *elemNode) removeChild(c node) {
	for i, child := range e.children {
		if c == child {
			children := e.children
			copy(children[i:], children[i+1:])
			children[len(children)-1] = nil
			e.children = children[:len(children)-1]

			c.setParent(nil)
			return
		}
	}
}

type compoNode struct {
	id        string
	compoID   string
	controlID string
	name      string
	component app.Component
	fields    map[string]string
	parent    node
	root      node
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

func (c *compoNode) setParent(p node) {
	c.parent = p
}

func (c *compoNode) setRoot(r node) {
	r.setParent(c)
	c.root = r
}

func (c *compoNode) removeRoot() {
	c.root.setParent(nil)
	c.root = nil
}

func attrsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}

	for k, va := range a {
		if vb, ok := b[k]; !ok || va != vb {
			return false
		}
	}
	return true
}
