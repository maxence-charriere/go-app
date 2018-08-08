package html

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
			<input type="text" required>
			<lib.Foo Bar="42">
			<lib.bar />
			<svg>
				<path d="M 42.42 Z "></path>
				<path d="M 21.21 Z " />
			</svg>
		</div>
		`)
	require.NoError(t, err)

	div := root.(*elemNode)
	require.Equal(t, "div", div.TagName())
	require.Len(t, div.Children(), 6)

	h1 := div.Children()[0].(*elemNode)
	require.Equal(t, "h1", h1.TagName())
	require.Len(t, h1.Children(), 1)

	text := h1.Children()[0].(*textNode)
	require.Equal(t, "hello", text.Text())

	br := div.Children()[1].(*elemNode)
	require.Equal(t, "br", br.TagName())
	require.Empty(t, br.Children())

	input := div.Children()[2].(*elemNode)
	require.Equal(t, "input", input.TagName())
	require.Empty(t, input.Children())

	foo := div.Children()[3].(*compoNode)
	require.Equal(t, "lib.foo", foo.Name())
	require.Equal(t, "42", foo.Fields()["bar"])

	bar := div.Children()[4].(*compoNode)
	require.Equal(t, "lib.bar", bar.Name())

	svg := div.Children()[5].(*elemNode)
	require.Equal(t, "svg", svg.TagName())
	require.Len(t, svg.Children(), 2)

	pathA := svg.Children()[0].(*elemNode)
	require.Equal(t, "path", pathA.TagName())
	require.Equal(t, "M 42.42 Z ", pathA.Attrs()["d"])
	require.Empty(t, pathA.Children())

	pathB := svg.Children()[1].(*elemNode)
	require.Equal(t, "path", pathB.TagName())
	require.Equal(t, "M 21.21 Z ", pathB.Attrs()["d"])
	require.Empty(t, pathB.Children())
}

func TestDecodeNodesError(t *testing.T) {
	_, err := decodeNodes(`
		<div>
			<div %error>
		</div>
		`)
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
