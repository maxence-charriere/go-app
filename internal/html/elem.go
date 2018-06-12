package html

import (
	"github.com/google/uuid"
	"github.com/murlokswarm/app"
)

type elemNode struct {
	id        string
	compoID   string
	controlID string
	tagName   string
	attrs     map[string]string
	parent    node
	children  []node
	changes   []Change
}

func newElemNode() *elemNode {
	n := &elemNode{
		id: uuid.New().String(),
	}

	n.changes = []Change{
		createElemChange(n),
	}
	return n
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

func (e *elemNode) SetAttrs(a map[string]string) {
	e.attrs = a
	e.changes = append(e.changes, setAttrsChange(a))
}

func (e *elemNode) Parent() app.DOMNode {
	return e.parent
}

func (e *elemNode) SetParent(p node) {
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
	c.SetParent(e)
	e.changes = append(e.changes, appendChildChange(e.ID(), c.ID()))
}

func (e *elemNode) removeChild(c node) {
	for i, child := range e.children {
		if c == child {
			children := e.children
			copy(children[i:], children[i+1:])
			children[len(children)-1] = nil
			e.children = children[:len(children)-1]

			c.SetParent(nil)
			e.changes = append(e.changes, removeChildChange(e.ID(), c.ID()))
			return
		}
	}
}

func (e *elemNode) ConsumeChanges() []Change {
	changes := e.changes
	e.changes = nil
	return changes
}
