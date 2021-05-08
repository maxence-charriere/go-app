package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExpire(t *testing.T) {
	testCache(t, &Expire{
		ItemTTL: time.Minute,
	})
}

func TestExpireExpire(t *testing.T) {
	ctx := context.TODO()

	c := Expire{
		ItemTTL: time.Minute,
	}

	c.Set(ctx, "/hello", String("hello"))
	require.Equal(t, 1, c.Len())
	require.Equal(t, 5, c.Size())

	c.Set(ctx, "/world", String("world"))
	require.Equal(t, 2, c.Len())
	require.Equal(t, 10, c.Size())

	c.Set(ctx, "/goodbye", String("goodbye"))
	require.Equal(t, 3, c.Len())
	require.Equal(t, 17, c.Size())

	c.items["/hello"].expiresAt = time.Now().Add(-time.Second)
	c.items["/world"].expiresAt = time.Now().Add(-time.Second)

	c.Set(ctx, "/goodmorning", String("goodmorning"))
	require.Equal(t, 2, c.Len())
	require.Equal(t, 18, c.Size())
	require.Len(t, c.queue, 2)
}

func TestExpireSetSameKey(t *testing.T) {
	ctx := context.TODO()

	c := Expire{
		ItemTTL: time.Minute,
	}

	c.Set(ctx, "/foo", String("foo"))
	require.Len(t, c.queue, 1)
	require.Equal(t, 1, c.Len())
	require.Equal(t, 3, c.Size())

	c.Set(ctx, "/bar", String("bar"))
	require.Len(t, c.queue, 2)
	require.Equal(t, 2, c.Len())
	require.Equal(t, 6, c.Size())

	c.Set(ctx, "/bar", String("barre"))
	require.Len(t, c.queue, 3)
	require.Equal(t, 2, c.Len())
	require.Equal(t, 8, c.Size())

	c.Set(ctx, "/foo", String("fooo"))
	require.Len(t, c.queue, 2)
	require.Equal(t, 2, c.Len())
	require.Equal(t, 9, c.Size())
}
