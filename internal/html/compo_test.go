package html

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompoNode(t *testing.T) {
	c := newCompoNode("foo", attrs)
	c.compoID = "compo"
	c.controlID = "control"

	assert.Equal(t, "foo", c.Name())
	assert.Equal(t, attrs, c.Fields())
	testNode(t, c)

	root := newTextNode()
	root.SetText("hello")

	c.SetRoot(root)
	assert.Equal(t, root, c.Root())
	assert.Equal(t, c, root.Parent())

	changes := c.ConsumeChanges()
	assert.Len(t, changes, 2)

	c.RemoveRoot()
	assert.Nil(t, root.Parent())
	assert.Nil(t, c.root)
	assert.Len(t, c.changes, 1)

	root = newTextNode()
	root.SetText("world")
	c.SetRoot(root)

	c.Close()
	assert.Len(t, c.changes, 1)
}
