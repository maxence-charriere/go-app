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

func BenchmarkMemoryCacheRead(b *testing.B) {
	c := newMemoryCache(1)
	item := cacheItem{Path: "/benchmark", Body: []byte("cached response")}
	c.Set(item)
	got, ok := c.Get(item.Path)
	require.True(b, ok)
	require.Equal(b, item, got)
	for _, parallel := range []bool{false, true} {
		name := "serial"
		if parallel {
			name = "parallel"
		}
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			if parallel {
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						c.Get(item.Path)
					}
				})
				return
			}
			for b.Loop() {
				c.Get(item.Path)
			}
		})
	}
}
