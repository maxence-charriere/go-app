package dom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsToGoHandler(t *testing.T) {
	n, v := JsToGoHandler("change", "OnChange")
	assert.Equal(t, "change", n)
	assert.Equal(t, "OnChange", v)

	n, v = JsToGoHandler("onchange", "js:alert('change')")
	assert.Equal(t, "onchange", n)
	assert.Equal(t, "alert('change')", v)

	n, v = JsToGoHandler("onchange", "OnChange")
	assert.Equal(t, "onchange", n)
	assert.Equal(t, "callCompoHandler(this, event, 'OnChange')", v)
}

func TestHrefCompoFmt(t *testing.T) {
	n, v := HrefCompoFmt("link", "hello")
	assert.Equal(t, "link", n)
	assert.Equal(t, "hello", v)

	n, v = HrefCompoFmt("href", "http://localhost:\\hello/world")
	assert.Equal(t, "href", n)
	assert.Equal(t, "http://localhost:\\hello/world", v)

	n, v = HrefCompoFmt("href", "http://hello")
	assert.Equal(t, "href", n)
	assert.Equal(t, "http://hello", v)

	n, v = HrefCompoFmt("href", "hello")
	assert.Equal(t, "href", n)
	assert.Equal(t, "compo:///hello", v)

	n, v = HrefCompoFmt("href", "/hello")
	assert.Equal(t, "href", n)
	assert.Equal(t, "compo:///hello", v)
}
