package cache

import (
	"context"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

func testCache(t *testing.T, c Cache) {
	ctx := context.TODO()

	_, isCached := c.Get(ctx, "/foo")
	require.False(t, isCached)
	require.Zero(t, c.Len())
	require.Zero(t, c.Size())

	c.Set(ctx, "/foo", Bytes("foo"))
	require.Equal(t, 1, c.Len())
	require.Equal(t, 3, c.Size())

	c.Set(ctx, "/bar", String("bar"))
	require.Equal(t, 2, c.Len())
	require.Equal(t, 6, c.Size())

	foo, isCached := c.Get(ctx, "/foo")
	require.True(t, isCached)
	require.Equal(t, Bytes("foo"), foo.(Bytes))

	bar, isCached := c.Get(ctx, "/bar")
	require.True(t, isCached)
	require.Equal(t, String("bar"), bar.(String))

	c.Del(ctx, "/foo")
	require.Equal(t, 1, c.Len())
	require.Equal(t, 3, c.Size())

	foo, isCached = c.Get(ctx, "/foo")
	require.False(t, isCached)
	require.Nil(t, foo)
}

func TestItemSize(t *testing.T) {
	utests := []struct {
		scenario     string
		item         Item
		expectedSize int
	}{
		{
			scenario:     "bytes",
			item:         Bytes("boo"),
			expectedSize: 3,
		},
		{
			scenario:     "string",
			item:         String("hello"),
			expectedSize: 5,
		},
		{
			scenario:     "int",
			item:         Int(42),
			expectedSize: int(unsafe.Sizeof(42)),
		},
		{
			scenario:     "float",
			item:         Float(42.0),
			expectedSize: int(unsafe.Sizeof(42.1)),
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			require.Equal(t, u.expectedSize, u.item.Size())
		})
	}
}
