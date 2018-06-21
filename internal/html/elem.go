package html

import (
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

func newElemNode(id string, tagName string) *elemNode {
	n := &elemNode{
		id:      id,
		tagName: tagName,
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
			c.Close()
			e.changes = append(e.changes, c.ConsumeChanges()...)
			return
		}
	}
}

func (e *elemNode) Close() {
	e.changes = e.changes[:0]
	for _, c := range e.children {
		c.Close()
		e.changes = append(e.changes, c.ConsumeChanges()...)
	}
	e.changes = append(e.changes, deleteNodeChange(e.ID()))
}

func (e *elemNode) ConsumeChanges() []Change {
	changes := make([]Change, 0, len(e.changes))
	for _, c := range e.children {
		changes = append(changes, c.ConsumeChanges()...)
	}
	changes = append(changes, e.changes...)
	e.changes = e.changes[:0]
	return changes
}
