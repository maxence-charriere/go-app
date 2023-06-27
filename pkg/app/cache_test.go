package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemoryCache(t *testing.T) {
	c := newMemoryCache(5)

	i := cacheItem{
		Path:            "/test",
		ContentType:     "text/html",
		ContentEncoding: "gzip",
		Body:            []byte("test"),
	}

	ic, ok := c.Get(i.Path)
	require.Zero(t, ic)
	require.False(t, ok)

	c.Set(i)
	ic, ok = c.Get(i.Path)
	require.True(t, ok)
	require.Equal(t, i, ic)
}
