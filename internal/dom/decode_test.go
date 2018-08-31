package dom

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeNodes(t *testing.T) {
	root, err := decodeNodes(`
		<div>
			<!-- Comment -->	
			<h1>hello</h1>
			<br>
			<input type="text" required onchange="Test">
			<lib.Foo Bar="42">
			<lib.bar />
			<svg>
				<path d="M 42.42 Z "></path>
				<path d="M 21.21 Z " />
			</svg>
			<a href="html.Foo"></a>
		</div>
		`, JsToGoHandler, HrefCompoFmt)
	require.NoError(t, err)

	div := root.(*elem)
	require.Equal(t, "div", div.TagName())
	require.Empty(t, div.namespace)
	require.Len(t, div.children, 7)

	h1 := div.children[0].(*elem)
	require.Equal(t, "h1", h1.TagName())
	require.Empty(t, h1.namespace)
	require.Len(t, h1.children, 1)

	text := h1.children[0].(*text)
	require.Equal(t, "hello", text.text)

	br := div.children[1].(*elem)
	require.Equal(t, "br", br.TagName())
	require.Empty(t, br.namespace)
	require.Empty(t, br.children)

	input := div.children[2].(*elem)
	require.Equal(t, "input", input.TagName())
	require.Empty(t, input.namespace)
	require.Empty(t, input.children)

	foo := div.children[3].(*compo)
	require.Equal(t, "lib.foo", foo.name)
	require.Equal(t, "42", foo.fields["bar"])

	bar := div.children[4].(*compo)
	require.Equal(t, "lib.bar", bar.name)

	svg := div.children[5].(*elem)
	require.Equal(t, "svg", svg.TagName())
	require.Equal(t, svgNamespace, svg.namespace)
	require.Len(t, svg.children, 2)

	pathA := svg.children[0].(*elem)
	require.Equal(t, "path", pathA.TagName())
	require.Equal(t, svgNamespace, pathA.namespace)
	require.Equal(t, "M 42.42 Z ", pathA.attrs["d"])
	require.Empty(t, pathA.children)

	pathB := svg.children[1].(*elem)
	require.Equal(t, "path", pathB.TagName())
	require.Equal(t, svgNamespace, pathB.namespace)
	require.Equal(t, "M 21.21 Z ", pathB.attrs["d"])
	require.Empty(t, pathB.children)

	a := div.children[6].(*elem)
	require.Equal(t, "a", a.TagName())
	require.Empty(t, a.namespace)
	require.Equal(t, "compo:///html.Foo", a.attrs["href"])
}

func TestDecodeNodesError(t *testing.T) {
	_, err := decodeNodes(`
		<div>
			<div %error>
		</div>
		`, JsToGoHandler, HrefCompoFmt)
	assert.Error(t, err)
}

func TestIsHTMLTagName(t *testing.T) {
	assert.True(t, isHTMLTagName("div"))
	assert.False(t, isHTMLTagName("foo"))
}

func TestIsCompoTagName(t *testing.T) {
	assert.True(t, isCompoTagName("hello", false))
	assert.False(t, isCompoTagName("hello", true))
	assert.False(t, isCompoTagName("div", false))
	assert.False(t, isCompoTagName("div", true))
}

func TestTagName(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "div",
			expected: "div",
		},
		{
			name:     "viewbox",
			expected: "viewBox",
		},
	}

	for _, test := range tests {
		n := tagName(test.name)
		assert.Equal(t, test.expected, n)
	}
}
