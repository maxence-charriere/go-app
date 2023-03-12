package app

import (
	"sync"
)

type cacheItem struct {
	// The request path.
	Path string

	// The response content type.
	ContentType string

	// The response content encoding.
	ContentEncoding string

	// The response body.
	Body []byte
}

func (i cacheItem) Len() int {
	return len(i.Body)
}

type memoryCache struct {
	mu    sync.RWMutex
	items map[string]cacheItem
}

func newMemoryCache(size int) *memoryCache {
	return &memoryCache{
		items: make(map[string]cacheItem, size),
	}
}

func (c *memoryCache) Set(i cacheItem) {
	c.mu.Lock()
	c.items[i.Path] = i
	c.mu.Unlock()
}

func (c *memoryCache) Get(path string) (cacheItem, bool) {
	c.mu.Lock()
	i, ok := c.items[path]
	c.mu.Unlock()
	return i, ok
}
