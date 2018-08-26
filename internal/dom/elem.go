package dom

import (
	"github.com/google/uuid"
	"github.com/murlokswarm/app"
)

type elem struct {
	id        string
	compoID   string
	tagName   string
	namespace string
	attrs     map[string]string
	parent    node
	children  []node
	changes   []Change
}

func newElem(tagName string, namespace string) *elem {
	e := &elem{
		id:        tagName + ":" + uuid.New().String(),
		tagName:   tagName,
		namespace: namespace,
	}

	e.changes = append(e.changes, createElemChange(e.id, e.tagName, e.namespace))
	return e
}

func (e *elem) ID() string {
	return e.id
}

func (e *elem) CompoID() string {
	return e.compoID
}

func (e *elem) TagName() string {
	return e.tagName
}

func (e *elem) SetAttrs(attrs map[string]string) {
	e.attrs = attrs
	e.changes = append(e.changes, setAttrsChange(e.id, attrs))
}

func (e *elem) Parent() app.Node {
	return e.parent
}

func (e *elem) SetParent(p node) {
	e.parent = p
}

func (e *elem) appendChild(child node) {
	e.children = append(e.children, child)
	child.SetParent(e)

	e.changes = append(e.changes, appendChildChange(e.id, child.ID()))
}

func (e *elem) removeChild(child node) {
	for i, c := range e.children {
		if c == child {
			e.changes = append(e.changes, removeChildChange(e.id, c.ID()))

			c.Close()
			e.changes = append(e.changes, c.Flush()...)

			children := e.children
			copy(children[i:], children[i+1:])
			children[len(children)-1] = nil
			e.children = children[:len(children)-1]
			return
		}
	}
}

func (e *elem) replaceChild(old, new node) {
	for i, c := range e.children {
		if c == old {
			e.children[i] = new
			new.SetParent(e)
			e.changes = append(e.changes, new.Flush()...)
			e.changes = append(e.changes, replaceChildChange(e.id, old.ID(), new.ID()))

			old.Close()
			e.changes = append(e.changes, old.Flush()...)
			return
		}
	}
}

func (e *elem) Flush() []Change {
	changes := make([]Change, 0, len(e.changes))

	for _, c := range e.children {
		changes = append(changes, c.Flush()...)
	}

	changes = append(changes, e.changes...)
	e.changes = e.changes[:0]
	return changes
}

func (e *elem) Close() {
	for _, c := range e.children {
		c.Close()
		e.changes = append(e.changes, c.Flush()...)
	}

	e.SetParent(nil)
	e.changes = append(e.changes, deleteNodeChange(e.id))
}
