package html

import (
	"github.com/google/uuid"
	"github.com/murlokswarm/app"
)

type textNode struct {
	id        string
	compoID   string
	controlID string
	text      string
	parent    node
	changes   []Change
}

func newTextNode() *textNode {
	n := &textNode{
		id: uuid.New().String(),
	}

	n.changes = []Change{
		createTextChange(n),
	}
	return n
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

func (t *textNode) SetText(text string) {
	t.text = text
	t.changes = append(t.changes, setTextChange(text))
}

func (t *textNode) Parent() app.DOMNode {
	return t.parent
}

func (t *textNode) SetParent(p node) {
	t.parent = p
}

func (t *textNode) ConsumeChanges() []Change {
	changes := t.changes
	t.changes = nil
	return changes
}
