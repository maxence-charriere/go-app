package dom

import (
	"github.com/google/uuid"
	"github.com/murlokswarm/app"
)

type text struct {
	id      string
	compoID string
	parent  node
	changes changer
	text    string
}

func newText(ch changer) *text {
	t := &text{
		id:      "text-" + uuid.New().String(),
		changes: ch,
	}

	t.changes.appendChanges(createTextChange(t.id))
	return t
}

func (t *text) Close() {
	t.SetParent(nil)
	t.changes.appendChanges(deleteNodeChange(t.id))
}

func (t *text) ID() string {
	return t.id
}

func (t *text) CompoID() string {
	return t.compoID
}

func (t *text) Parent() app.Node {
	return t.parent
}

func (t *text) SetParent(p node) {
	t.parent = p
}

func (t *text) SetText(text string) {
	t.text = text
	t.SetParent(nil)
	t.changes.appendChanges(setTextChange(t.id, t.text))
}
