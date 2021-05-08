package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLRU(t *testing.T) {
	testCache(t, &LRU{
		ItemTTL: time.Minute,
	})
}

func TestLRUEvict(t *testing.T) {
	ctx := context.TODO()
	isHelloEvicted := false

	c := LRU{
		MaxSize: 12,
		ItemTTL: time.Minute,
		OnEvict: func(key string, i Item) {
			isHelloEvicted = key == "/hello"
			require.Equal(t, String("hello"), i)
		},
	}

	c.Set(ctx, "/hello", String("hello"))
	require.Len(t, c.priority, 1)
	require.Equal(t, 1, c.Len())
	require.Equal(t, 5, c.Size())

	c.Set(ctx, "/world", String("world"))
	require.Len(t, c.priority, 2)
	require.Equal(t, 2, c.Len())
	require.Equal(t, 10, c.Size())

	c.Get(ctx, "/world")
	c.Set(ctx, "/goodbye", String("goodbye"))
	require.Len(t, c.priority, 2)
	require.Equal(t, 2, c.Len())
	require.Equal(t, 12, c.Size())
	require.True(t, isHelloEvicted)

	hello, isCached := c.Get(ctx, "/hello")
	require.False(t, isCached)
	require.Nil(t, hello)

	world, isCached := c.Get(ctx, "/world")
	require.True(t, isCached)
	require.Equal(t, String("world"), world)

	goodbye, isCached := c.Get(ctx, "/goodbye")
	require.True(t, isCached)
	require.Equal(t, String("goodbye"), goodbye)
}

func TestLRUSetSameKey(t *testing.T) {
	ctx := context.TODO()

	c := LRU{
		ItemTTL: time.Minute,
	}

	c.Set(ctx, "/test", String("test"))
	require.Len(t, c.priority, 1)
	require.Equal(t, 1, c.Len())
	require.Equal(t, 4, c.Size())

	c.Set(ctx, "/test", String("unit-test"))
	require.Len(t, c.priority, 2)
	require.Equal(t, 1, c.Len())
	require.Equal(t, 13, c.Size())
}

func TestSortLRUItems(t *testing.T) {
	now := time.Now()

	utests := []struct {
		scenario string
		now      time.Time
		in       []*lruItem
		out      []*lruItem
	}{
		{
			scenario: "nil",
			now:      now,
		},
		{
			scenario: "empty",
			now:      now,
			in:       []*lruItem{},
			out:      []*lruItem{},
		},
		{
			scenario: "1 item",
			now:      now,
			in: []*lruItem{
				{
					count:     1,
					expiresAt: now.Add(time.Second),
				},
			},
			out: []*lruItem{
				{
					count:     1,
					expiresAt: now.Add(time.Second),
				},
			},
		},
		{
			scenario: "2 items",
			now:      now,
			in: []*lruItem{
				{
					count:     1,
					expiresAt: now.Add(time.Second),
				},
				{
					count:     42,
					expiresAt: now.Add(time.Second),
				},
			},
			out: []*lruItem{
				{
					count:     42,
					expiresAt: now.Add(time.Second),
				},
				{
					count:     1,
					expiresAt: now.Add(time.Second),
				},
			},
		},
		{
			scenario: "multiple items",
			now:      now,
			in: []*lruItem{
				{
					count:     1,
					expiresAt: now.Add(time.Second),
				},
				{
					count:     42,
					expiresAt: now.Add(time.Second),
				},
				{
					count:     1,
					expiresAt: now.Add(time.Second),
				},
				{
					count:     7,
					expiresAt: now.Add(time.Second),
				},
				{
					count:     -21,
					expiresAt: now.Add(time.Second),
				},
			},
			out: []*lruItem{
				{
					count:     42,
					expiresAt: now.Add(time.Second),
				},
				{
					count:     7,
					expiresAt: now.Add(time.Second),
				},
				{
					count:     1,
					expiresAt: now.Add(time.Second),
				},
				{
					count:     1,
					expiresAt: now.Add(time.Second),
				},
				{
					count:     -21,
					expiresAt: now.Add(time.Second),
				},
			},
		},
		{
			scenario: "multiple items with expired ones",
			now:      now,
			in: []*lruItem{
				{
					count:     1,
					expiresAt: now.Add(time.Second),
				},
				{
					count:     42,
					expiresAt: now.Add(-time.Second),
				},
				{
					count:     1,
					expiresAt: now.Add(time.Second),
				},
				{
					count:     7,
					expiresAt: now.Add(time.Second),
				},
				{
					count:     -21,
					expiresAt: now.Add(time.Second),
				},
			},
			out: []*lruItem{
				{
					count:     7,
					expiresAt: now.Add(time.Second),
				},
				{
					count:     1,
					expiresAt: now.Add(time.Second),
				},
				{
					count:     1,
					expiresAt: now.Add(time.Second),
				},
				{
					count:     42,
					expiresAt: now.Add(-time.Second),
				},
				{
					count:     -21,
					expiresAt: now.Add(time.Second),
				},
			},
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			sortLRUItems(u.now, u.in)
			require.Equal(t, u.out, u.in)
		})
	}
}
