package html

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextNode(t *testing.T) {
	text := newTextNode("node")
	text.compoID = "compo"
	text.controlID = "control"

	assert.Len(t, text.changes, 1)
	assert.Equal(t, createText, text.changes[0].Type)
	testNode(t, text)

	text.SetText("hello")
	assert.Equal(t, "hello", text.Text())
	assert.Len(t, text.changes, 2)
	assert.Equal(t, setText, text.changes[1].Type)
	assert.Equal(t, textValue{"node", "hello"}, text.changes[1].Value)

	changes := text.ConsumeChanges()
	assert.Len(t, changes, 2)
	assert.Empty(t, text.changes)

	text.Close()
	assert.Len(t, text.changes, 1)
	assert.Equal(t, deleteNode, text.changes[0].Type)
	assert.Equal(t, "node", text.changes[0].Value)
}
