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
	changes   []Change
}

func newCompoNode(id string, name string, fields map[string]string) *compoNode {
	return &compoNode{
		id:     id,
		name:   name,
		fields: fields,
	}
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
	// HANDLE PARRENT
	r.SetParent(c)
	c.root = r
}

func (c *compoNode) UnsetRoot(r node) {
	// HANDLE PARRENT
	r.SetParent(nil)
	c.root = nil
}

func (c *compoNode) RemoveRoot() {
	c.root.SetParent(nil)
	c.root = nil
}

func (c *compoNode) Close() {}

func (c *compoNode) ConsumeChanges() []Change {
	return c.root.ConsumeChanges()
}
