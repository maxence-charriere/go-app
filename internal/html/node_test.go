package html

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/stretchr/testify/assert"
)

var changes = []app.DOMChange{app.DOMChange{
	Type:  app.DOMNoChanges,
	Value: 42,
}}

var attrs = map[string]string{
	"hello": "world",
}

func TestTextNode(t *testing.T) {
	text := &textNode{
		id:        "node",
		compoID:   "compo",
		controlID: "control",
		text:      "hello",
		changes:   changes,
	}

	assert.Equal(t, "hello", text.Text())
	testNode(t, text)
}

func TestElemNode(t *testing.T) {
	elem := &elemNode{
		id:        "node",
		compoID:   "compo",
		controlID: "control",
		tagName:   "img",
		attrs:     attrs,
		changes:   changes,
	}

	assert.Equal(t, "img", elem.TagName())
	assert.Equal(t, attrs, elem.Attrs())

	childA := &textNode{
		id:   "text",
		text: "hello",
	}

	childB := &textNode{
		id:   "text",
		text: "world",
	}

	elem.appendChild(childA)
	elem.appendChild(childB)
	assert.Equal(t, []app.DOMNode{childA, childB}, elem.Children())
	testNode(t, elem)
}

func TestCompoNode(t *testing.T) {
	compo := &compoNode{
		id:        "node",
		compoID:   "compo",
		controlID: "control",
		name:      "foo",
		fields:    attrs,
		changes:   changes,
	}

	assert.Equal(t, "foo", compo.Name())
	assert.Equal(t, attrs, compo.Fields())
	testNode(t, compo)
}

func testNode(t *testing.T, n node) {
	parent := &elemNode{
		id:      "parent",
		tagName: "div",
	}
	n.SetParent(parent)

	assert.Equal(t, parent, n.Parent())
	assert.Equal(t, "node", n.ID())
	assert.Equal(t, "compo", n.CompoID())
	assert.Equal(t, "control", n.ControlID())
	assert.Equal(t, changes, n.Changes())
}
