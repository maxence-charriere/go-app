package dom

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertChangeEqual(t *testing.T, expected, actual Change) {
	require.Equal(t, expected.Type, actual.Type)

	switch expected.Type {
	case setText:
		assert.Equal(t, expected.Value.(textValue).Text, actual.Value.(textValue).Text)

	case createElem:
		assert.Equal(t, expected.Value.(elemValue).TagName, actual.Value.(elemValue).TagName)

	case setAttrs:
		assert.Equal(t, expected.Value.(elemValue).Attrs, actual.Value.(elemValue).Attrs)

	case createCompo:
		assert.Equal(t, expected.Value.(compoValue).Name, actual.Value.(compoValue).Name)

	default:
	}
}

func assertChangesEqual(t *testing.T, expected, actual []Change) {
	assert.Len(t, actual, len(expected))

	for i := range actual {
		assertChangeEqual(t, expected[i], actual[i])
	}
}

func prettyChanges(c []Change) string {
	b, _ := json.MarshalIndent(c, "", "    ")
	return string(b)
}
