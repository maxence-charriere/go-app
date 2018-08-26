package dom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestElem(t *testing.T) {
	p := newElem("p", "")

	e := newElem("div", "")
	e.SetAttrs(map[string]string{"foo": "bar"})
	assert.NotEmpty(t, e.ID())
	assert.Empty(t, e.CompoID())
	assert.Nil(t, e.Parent())

	e.SetParent(p)
	assert.Equal(t, p, e.Parent())

	c1 := newElem("h1", "")
	e.appendChild(c1)
	assert.Len(t, e.children, 1)
	assert.Equal(t, c1, e.children[0])
	assert.Equal(t, e, c1.Parent())

	c2 := newElem("p", "")
	e.appendChild(c2)
	assert.Len(t, e.children, 2)
	assert.Equal(t, c2, e.children[1])
	assert.Equal(t, e, c2.Parent())

	c3 := newElem("span", "")
	e.replaceChild(c2, c3)
	assert.Len(t, e.children, 2)
	assert.Equal(t, c3, e.children[1])
	assert.Equal(t, e, c3.Parent())

	e.removeChild(c1)
	assert.Len(t, e.children, 1)

	e.Close()
}
