package dom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestText(t *testing.T) {
	p := newElem("p", "")

	text := newText()
	assert.NotEmpty(t, text.ID())
	assert.Empty(t, text.CompoID())
	assert.Nil(t, text.Parent())

	text.SetParent(p)
	assert.Equal(t, p, text.Parent())

	text.SetText("hello")
	text.Close()
}
