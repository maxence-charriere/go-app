package html

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var attrs = map[string]string{
	"hello": "world",
}

func testNode(t *testing.T, n node) {
	parent := newElemNode("div")

	n.SetParent(parent)
	assert.Equal(t, parent, n.Parent())
	assert.Equal(t, "compo", n.CompoID())
	assert.Equal(t, "control", n.ControlID())
}

func TestAttrsEqual(t *testing.T) {
	tests := []struct {
		scenario string
		a        map[string]string
		b        map[string]string
		equals   bool
	}{
		{
			scenario: "emptys",
			equals:   true,
		},
		{
			scenario: "equals",
			a: map[string]string{
				"a": "foo",
				"b": "bar",
				"c": "boo",
			},
			b: map[string]string{
				"b": "bar",
				"c": "boo",
				"a": "foo",
			},
			equals: true,
		},
		{
			scenario: "different lengths",
			a: map[string]string{
				"a": "foo",
				"b": "bar",
				"c": "boo",
			},
			b: map[string]string{
				"a": "foo",
				"b": "bar",
			},
			equals: false,
		},
		{
			scenario: "different values",
			a: map[string]string{
				"a": "foo",
			},
			b: map[string]string{
				"a": "bar",
			},
			equals: false,
		},
		{
			scenario: "different keys",
			a: map[string]string{
				"a": "foo",
			},
			b: map[string]string{
				"b": "foo",
			},
			equals: false,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			equals := attrsEqual(test.a, test.b)
			assert.Equal(t, test.equals, equals)
		})
	}
}
