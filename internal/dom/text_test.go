package dom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestText(t *testing.T) {
	dom := NewDOM()
	p := newElem(dom, "p")
	dom.ReadChanges()

	text := newText(dom)

	assert.NotEmpty(t, text.ID())
	assert.Empty(t, text.CompoID())
	assert.Nil(t, text.Parent())

	text.SetParent(p)
	assert.Equal(t, p, text.Parent())

	text.SetText("hello")
	text.Close()

	changes := dom.ReadChanges()
	t.Log(prettyChanges(changes))

	assertChangesEqual(t, []Change{
		createTextChange(""),
		setTextChange("", "hello"),
		deleteNodeChange(""),
	}, changes)
}
