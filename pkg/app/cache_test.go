package app

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSortPreRenderLRUCacheItem(t *testing.T) {
	utests := []struct {
		scenario string
		in       []*preRenderLRUCacheItem
		out      []*preRenderLRUCacheItem
	}{
		{
			scenario: "nil slice",
		},
		{
			scenario: "empty slice",
			in:       []*preRenderLRUCacheItem{},
			out:      []*preRenderLRUCacheItem{},
		},
		{
			scenario: "1 item",
			in: []*preRenderLRUCacheItem{
				{Count: 42},
			},
			out: []*preRenderLRUCacheItem{
				{Count: 42},
			},
		},
		{
			scenario: "2 items asc",
			in: []*preRenderLRUCacheItem{
				{Count: 21},
				{Count: 42},
			},
			out: []*preRenderLRUCacheItem{
				{Count: 42},
				{Count: 21},
			},
		},
		{
			scenario: "multiple items",
			in: []*preRenderLRUCacheItem{
				{Count: 42},
				{Count: 21},
				{Count: 2},
				{Count: 836},
				{Count: 14},
				{Count: 5},
				{Count: 0},
				{Count: 98},
				{Count: 1},
				{Count: 21},
			},
			out: []*preRenderLRUCacheItem{
				{Count: 836},
				{Count: 98},
				{Count: 42},
				{Count: 21},
				{Count: 21},
				{Count: 14},
				{Count: 5},
				{Count: 2},
				{Count: 1},
				{Count: 0},
			},
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			sortPreRenderLRUCacheItem(u.in)
			require.Equal(t, u.out, u.in)
		})
	}
}

func TestPreRenderLRUCache(t *testing.T) {
	testPreRenderCache(t, NewPreRenderLRUCache(100, time.Second))
}

func TestPreRenderCache(t *testing.T) {
	testPreRenderCache(t, newPreRenderCache(1))
}

func testPreRenderCache(t *testing.T, c PreRenderCache) {
	ctx := context.TODO()

	i := PreRenderedItem{
		Path:            "/test",
		ContentType:     "text/html",
		ContentEncoding: "gzip",
		Body:            []byte("test"),
	}

	ic, ok := c.Get(ctx, i.Path)
	require.Zero(t, ic)
	require.False(t, ok)

	c.Set(ctx, i)
	ic, ok = c.Get(ctx, i.Path)
	require.True(t, ok)
	require.Equal(t, i, ic)
}

func TestPreRenderLRUCacheExpire(t *testing.T) {
	ctx := context.TODO()
	evictCalled := false
	onEvict := func(int, int) { evictCalled = true }

	c := NewPreRenderLRUCache(16, -time.Second, onEvict).(*preRenderLRUCache)

	items := []PreRenderedItem{
		{
			Path: "/test1",
			Body: []byte("test"),
		},
		{
			Path: "/test2",
			Body: []byte("test"),
		},
		{
			Path: "/test3",
			Body: []byte("test"),
		},
		{
			Path: "/test4",
			Body: []byte("test"),
		},
	}

	for _, i := range items {
		c.Set(ctx, i)
	}
	require.Len(t, c.items, 4)
	require.Len(t, c.priorities, 4)
	require.Equal(t, 16, c.size)

	for _, i := range items {
		ic, ok := c.Get(ctx, i.Path)
		require.Zero(t, ic)
		require.False(t, ok)
	}
	require.Len(t, c.items, 4)
	require.Len(t, c.priorities, 4)
	require.Equal(t, 16, c.size)

	c.Set(ctx, PreRenderedItem{
		Path: "/test5",
		Body: []byte("test"),
	})
	require.Len(t, c.items, 1)
	require.Len(t, c.priorities, 1)
	require.Equal(t, 4, c.size)
	require.False(t, evictCalled)
}

func TestPreRenderLRUCacheEvict(t *testing.T) {
	ctx := context.TODO()

	evictCount := 0
	evictSize := 0
	onEvict := func(count int, size int) {
		evictCount += count
		evictSize += size
	}

	c := NewPreRenderLRUCache(8, time.Second, onEvict).(*preRenderLRUCache)

	items := []PreRenderedItem{
		{
			Path: "/test1",
			Body: []byte("test"),
		},
		{
			Path: "/test2",
			Body: []byte("test"),
		},
	}

	for _, i := range items {
		c.Set(ctx, i)
	}
	require.Len(t, c.items, 2)
	require.Len(t, c.priorities, 2)
	require.Equal(t, 8, c.size)
	require.Equal(t, 0, evictCount)
	require.Equal(t, 0, evictSize)

	c.Set(ctx, PreRenderedItem{
		Path: "/test3",
		Body: []byte("test"),
	})
	require.Len(t, c.items, 2)
	require.Len(t, c.priorities, 2)
	require.Equal(t, 8, c.size)
	require.Equal(t, 1, evictCount)
	require.Equal(t, 4, evictSize)

	c.Set(ctx, PreRenderedItem{
		Path: "/test4",
		Body: []byte("testbig"),
	})
	require.Len(t, c.items, 1)
	require.Len(t, c.priorities, 1)
	require.Equal(t, 7, c.size)
	require.Equal(t, 3, evictCount)
	require.Equal(t, 12, evictSize)
}
