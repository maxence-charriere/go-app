package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsToGoHandler(t *testing.T) {
	n, v := jsToGoHandler("change", "OnChange")
	assert.Equal(t, "change", n)
	assert.Equal(t, "OnChange", v)

	n, v = jsToGoHandler("onchange", "js:alert('change')")
	assert.Equal(t, "onchange", n)
	assert.Equal(t, "alert('change')", v)

	n, v = jsToGoHandler("onchange", "OnChange")
	assert.Equal(t, "onchange", n)
	assert.Equal(t, "callCompoHandler(this, event, 'OnChange')", v)
}
