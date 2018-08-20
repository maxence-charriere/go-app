package dom

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertChangeEqual(t *testing.T, expected, actual Change) {
	require.Equal(t, expected.Type, actual.Type)

	assertIDs := func(expected, actual string) {
		sep := strings.IndexByte(expected, ':')
		require.Equal(t, expected[:sep], actual[:sep])
	}

	switch expected.Type {
	case setText:
		exp := expected.Value.(textValue)
		act := actual.Value.(textValue)

		assertIDs(exp.ID, act.ID)
		assert.Equal(t, exp.Text, act.Text)

	case createElem:
		exp := expected.Value.(elemValue)
		act := actual.Value.(elemValue)

		assertIDs(exp.ID, act.ID)
		assert.Equal(t, exp.TagName, act.TagName)

	case setAttrs:
		exp := expected.Value.(elemValue)
		act := actual.Value.(elemValue)

		assertIDs(exp.ID, act.ID)
		assert.Equal(t, exp.Attrs, act.Attrs)

	case appendChild:
		exp := expected.Value.(childValue)
		act := actual.Value.(childValue)

		assertIDs(exp.ParentID, act.ParentID)
		assertIDs(exp.ChildID, act.ChildID)

	case removeChild:
		exp := expected.Value.(childValue)
		act := actual.Value.(childValue)

		assertIDs(exp.ParentID, act.ParentID)
		assertIDs(exp.ChildID, act.ChildID)

	case replaceChild:
		exp := expected.Value.(childValue)
		act := actual.Value.(childValue)

		assertIDs(exp.ParentID, act.ParentID)
		assertIDs(exp.ChildID, act.ChildID)
		assertIDs(exp.OldID, act.OldID)

	case mountElem:
		exp := expected.Value.(elemValue)
		act := actual.Value.(elemValue)

		assertIDs(exp.ID, act.ID)

	case createCompo:
		exp := expected.Value.(compoValue)
		act := actual.Value.(compoValue)

		assertIDs(exp.ID, act.ID)
		assert.Equal(t, exp.Name, act.Name)

	case setCompoRoot:
		exp := expected.Value.(compoValue)
		act := actual.Value.(compoValue)

		assertIDs(exp.ID, act.ID)
		assertIDs(exp.RootID, act.RootID)
		assert.Equal(t, exp.Name, act.Name)

	case deleteNode:
		exp := expected.Value.(deleteValue)
		act := actual.Value.(deleteValue)

		assertIDs(exp.ID, act.ID)
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
