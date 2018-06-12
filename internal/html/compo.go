package html

import "github.com/murlokswarm/app"

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

func (c *compoNode) SetParent(p node) {
	c.parent = p
}

func (c *compoNode) Root() node {
	return c.root
}

func (c *compoNode) SetRoot(r node) {
	r.SetParent(c)
	c.root = r
}

func (c *compoNode) RemoveRoot() {
	c.root.SetParent(nil)
	c.root = nil
}
