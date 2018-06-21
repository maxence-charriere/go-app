package html

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestElemNode(t *testing.T) {
	e := newElemNode("node", "div")
	e.compoID = "compo"
	e.controlID = "control"
	e.attrs = attrs

	assert.Len(t, e.changes, 1)
	assert.Equal(t, createElem, e.changes[0].Type)
	testNode(t, e)

	c1 := newElemNode("c1", "h1")
	e.appendChild(c1)
	assert.Len(t, e.children, 1)
	assert.Equal(t, c1, e.children[0])
	assert.Equal(t, c1.Parent(), e)
	assert.Len(t, e.changes, 2)
	assert.Equal(t, appendChild, e.changes[1].Type)
	assert.Equal(t, childValue{ParentID: "node", ChildID: "c1"}, e.changes[1].Value)

	c2 := newElemNode("c2", "p")
	e.appendChild(c2)
	assert.Len(t, e.children, 2)
	assert.Equal(t, c2, e.children[1])
	assert.Len(t, e.changes, 3)
	assert.Equal(t, appendChild, e.changes[2].Type)
	assert.Equal(t, childValue{ParentID: "node", ChildID: "c2"}, e.changes[2].Value)

	e.removeChild(c1)
	assert.Len(t, e.children, 1)
	assert.Equal(t, c2, e.children[0])
	assert.Nil(t, c1.Parent())
	assert.Len(t, e.changes, 5)
	assert.Equal(t, removeChild, e.changes[3].Type)
	assert.Equal(t, childValue{ParentID: "node", ChildID: "c1"}, e.changes[3].Value)
	assert.Equal(t, deleteNode, e.changes[4].Type)
	assert.Equal(t, "c1", e.changes[4].Value)

	changes := e.ConsumeChanges()
	assert.Len(t, changes, 6)
	assert.Empty(t, e.changes)
	assert.Equal(t, createElem, changes[0].Type)

	e.Close()
	assert.Len(t, e.changes, 2)
	assert.Equal(t, deleteNode, e.changes[0].Type)
	assert.Equal(t, "c2", e.changes[0].Value)
	assert.Equal(t, deleteNode, e.changes[1].Type)
	assert.Equal(t, "node", e.changes[1].Value)
}
