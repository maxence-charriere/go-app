package dom

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func requireChangeMatch(t *testing.T, expected, actual change) {
	requireIDsMatch := func(expected, actual string) {
		expDelimiter := strings.IndexByte(expected, ':')
		actDelimiter := strings.IndexByte(actual, ':')

		if expDelimiter != actDelimiter {
			t.Fatal("bad id set")
		}

		if expDelimiter < 0 {
			return
		}

		require.Equal(t, expected[:expDelimiter], actual[:actDelimiter])
	}

	require.Equal(t, expected.Action, actual.Action)
	requireIDsMatch(expected.NodeID, actual.NodeID)
	requireIDsMatch(expected.CompoID, actual.CompoID)
	require.Equal(t, expected.Type, actual.Type)
	require.Equal(t, expected.Namespace, actual.Namespace)
	require.Equal(t, expected.Key, actual.Key)
	require.Equal(t, expected.Value, actual.Value)
	requireIDsMatch(expected.ChildID, actual.ChildID)
	requireIDsMatch(expected.NewChildID, actual.NewChildID)
	require.Equal(t, expected.IsCompo, actual.IsCompo)
}

func requireChangesMatches(t *testing.T, expected, actual []change) {
	require.Len(t, actual, len(expected))

	for i := range expected {
		requireChangeMatch(t, expected[i], actual[i])
	}
}

func TestIsHTMLNode(t *testing.T) {
	assert.True(t, isHTMLNode("div"))
	assert.False(t, isHTMLNode("foo"))
}

func TestIsCompoNode(t *testing.T) {
	assert.True(t, isCompoNode("hello", ""))
	assert.False(t, isCompoNode("hello", svg))
	assert.False(t, isCompoNode("div", ""))
	assert.False(t, isCompoNode("div", svg))
}

func TestNodeType(t *testing.T) {
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
		n := nodeType(test.name)
		assert.Equal(t, test.expected, n)
	}
}
