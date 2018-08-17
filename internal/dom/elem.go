package dom

import (
	"github.com/google/uuid"
	"github.com/murlokswarm/app"
)

type elem struct {
	id       string
	compoID  string
	parent   node
	changes  changer
	tagName  string
	attrs    map[string]string
	children []node
}

func newElem(ch changer, tagName string) *elem {
	e := &elem{
		id:      tagName + "-" + uuid.New().String(),
		changes: ch,
		tagName: tagName,
	}

	e.changes.appendChanges(createElemChange(e.id, e.tagName))
	return e
}

func (e *elem) Close() {
	for _, c := range e.children {
		c.Close()
	}

	e.SetParent(nil)
	e.changes.appendChanges(deleteNodeChange(e.id))
}

func (e *elem) ID() string {
	return e.id
}

func (e *elem) CompoID() string {
	return e.compoID
}

func (e *elem) Parent() app.Node {
	return e.parent
}

func (e *elem) SetParent(p node) {
	e.parent = p
}

func (e *elem) TagName() string {
	return e.tagName
}

func (e *elem) SetAttrs(attrs map[string]string) {
	e.attrs = attrs
	e.changes.appendChanges(setAttrsChange(e.id, attrs))
}

func (e *elem) appendChild(child node) {
	e.children = append(e.children, child)
	child.SetParent(e)
	e.changes.appendChanges(appendChildChange(e.id, child.ID()))
}

func (e *elem) removeChild(child node) {
	for i, c := range e.children {
		if c == child {
			e.changes.appendChanges(removeChildChange(e.id, c.ID()))
			c.Close()

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
			e.changes.appendChanges(replaceChildChange(e.id, old.ID(), new.ID()))
			old.Close()
			return
		}
	}
}
