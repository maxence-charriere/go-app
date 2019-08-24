package maestro

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventTransform(t *testing.T) {
	k, v := eventTransform("change", "OnChange")
	assert.Equal(t, "change", k)
	assert.Equal(t, "OnChange", v)

	k, v = eventTransform("onchange", "js:alert('change')")
	assert.Equal(t, "onchange", k)
	assert.Equal(t, "alert('change')", v)

	k, v = eventTransform("onchange", "OnChange")
	assert.Equal(t, "onchange", k)
	assert.Equal(t, "//go: OnChange", v)
}
