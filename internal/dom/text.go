package dom

import (
	"github.com/google/uuid"
	"github.com/murlokswarm/app"
)

type text struct {
	id      string
	compoID string
	text    string
	parent  node
	changes []Change
}

func newText() *text {
	t := &text{
		id: "text:" + uuid.New().String(),
	}

	t.changes = append(t.changes, createTextChange(t.id))
	return t
}

func (t *text) ID() string {
	return t.id
}

func (t *text) CompoID() string {
	return t.compoID
}

func (t *text) SetText(text string) {
	t.text = text
	t.SetParent(nil)
	t.changes = append(t.changes, setTextChange(t.id, t.text))
}

func (t *text) Parent() app.Node {
	return t.parent
}

func (t *text) SetParent(p node) {
	t.parent = p
}

func (t *text) Flush() []Change {
	c := t.changes
	t.changes = nil
	return c
}

func (t *text) Close() {
	t.SetParent(nil)
	t.changes = append(t.changes, deleteNodeChange(t.id))
}
