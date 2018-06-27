package html

import (
	"github.com/google/uuid"
	"github.com/murlokswarm/app"
)

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

func newCompoNode(name string, fields map[string]string) *compoNode {
	return &compoNode{
		id:     "compo-" + uuid.New().String(),
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
	r.SetParent(c)
	c.root = r
}

func (c *compoNode) RemoveRoot() {
	root := c.root
	root.Close()
	c.changes = append(c.changes, root.ConsumeChanges()...)
	c.root = nil
}

func (c *compoNode) Close() {
	c.changes = c.changes[:0]
	if c.root != nil {
		c.root.Close()
		c.changes = append(c.changes, c.root.ConsumeChanges()...)
	}
	c.SetParent(nil)
}

func (c *compoNode) ConsumeChanges() []Change {
	var changes []Change
	if c.root != nil {
		changes = append(changes, c.root.ConsumeChanges()...)
	}
	changes = append(changes, c.changes...)
	c.changes = c.changes[:0]
	return changes
}
