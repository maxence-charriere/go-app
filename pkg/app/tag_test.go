package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTags(t *testing.T) {
	tags := make(Tags)
	tags.Set("number", 42)
	require.Equal(t, "42", tags.Get("number"))
	require.Empty(t, tags["nothing"])
}
