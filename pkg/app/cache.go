package app

import (
	"context"
	"sync"
	"time"
)

// PreRenderCache is the interface that describes a cache that stores
// pre-rendered resources.
type PreRenderCache interface {
	// Get returns the item at the given path.
	Get(ctx context.Context, path string) (PreRenderedItem, bool)

	// Set stored the item at the given path.
	Set(ctx context.Context, i PreRenderedItem)
}

// PreRenderedItem represent an item that is stored in a PreRenderCache.
type PreRenderedItem struct {
	// The request path.
	Path string

	// The response content type.
	ContentType string

	// The response content encoding.
	ContentEncoding string

	// The response body.
	Body []byte
}

// Len return the body length.
func (r PreRenderedItem) Len() int {
	return len(r.Body)
}

// NewPreRenderLRUCache creates an in memory LRU cache that stores items for the
// given duration. If provided, on eviction functions are called when item are
// evicted.
func NewPreRenderLRUCache(size int, itemTTL time.Duration, onEvict ...func(count, size int)) PreRenderCache {
	return &preRenderLRUCache{
		itemTTL:    itemTTL,
		onEvict:    onEvict,
		priorities: make([]*preRenderLRUCacheItem, 0, 64),
		items:      make(map[string]*preRenderLRUCacheItem, 64),
		maxSize:    size,
	}
}

type preRenderLRUCache struct {
	itemTTL time.Duration
	onEvict []func(int, int)

	mu         sync.Mutex
	priorities []*preRenderLRUCacheItem
	items      map[string]*preRenderLRUCacheItem
	size       int
	maxSize    int
}

func (c *preRenderLRUCache) Get(ctx context.Context, path string) (PreRenderedItem, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	i, ok := c.items[path]
	if !ok {
		return PreRenderedItem{}, false
	}

	if time.Now().After(i.ExpiresAt) {
		i.Count = 0
		return PreRenderedItem{}, false
	}

	i.Count++
	return i.Item, true
}

func (c *preRenderLRUCache) Set(ctx context.Context, i PreRenderedItem) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.evict(i.Len())

	if c.size+i.Len() <= c.maxSize {
		item := &preRenderLRUCacheItem{
			ExpiresAt: time.Now().Add(c.itemTTL),
			Item:      i,
			Count:     1,
		}

		c.size += i.Len()
		c.items[i.Path] = item
		c.priorities = append(c.priorities, item)
	}
}

func (c *preRenderLRUCache) evict(size int) {
	if len(c.priorities) == 0 || c.size+size <= c.maxSize {
		return
	}

	sortPreRenderLRUCacheItem(c.priorities)

	for len(c.priorities) != 0 && c.last().Count == 0 {
		c.removeLast()
	}

	evictedCount := 0
	evictedSize := 0

	for len(c.priorities) != 0 && c.size+size > c.maxSize {
		last := c.last()
		c.removeLast()
		evictedCount++
		evictedSize += last.Item.Len()
	}

	if evictedCount == 0 {
		return
	}

	for _, fn := range c.onEvict {
		fn(evictedCount, evictedSize)
	}
}

func (c *preRenderLRUCache) removeLast() {
	item := c.last()
	c.size -= item.Item.Len()
	delete(c.items, item.Item.Path)

	end := c.lastIdx()
	c.priorities[end] = nil
	c.priorities = c.priorities[:end]
}

func (c *preRenderLRUCache) lastIdx() int {
	return len(c.priorities) - 1
}

func (c *preRenderLRUCache) last() *preRenderLRUCacheItem {
	return c.priorities[c.lastIdx()]
}

type preRenderLRUCacheItem struct {
	Count     uint
	ExpiresAt time.Time
	Item      PreRenderedItem
}

func sortPreRenderLRUCacheItem(s []*preRenderLRUCacheItem) {
	if len(s) < 2 {
		return
	}

	pi := sortPreRenderLRUCacheItemPartition(s)
	sortPreRenderLRUCacheItem(s[:pi])
	sortPreRenderLRUCacheItem(s[pi+1:])
}

func sortPreRenderLRUCacheItemPartition(s []*preRenderLRUCacheItem) int {
	piIdx := len(s) - 1
	pi := s[piIdx]

	i := -1
	for j, item := range s {
		if item.Count > pi.Count {
			i++
			s[i], s[j] = s[j], s[i]
		}
	}

	i++
	s[i], s[piIdx] = s[piIdx], s[i]
	return i
}

type preRenderCache struct {
	mu    sync.RWMutex
	items map[string]PreRenderedItem
}

func newPreRenderCache(size int) *preRenderCache {
	return &preRenderCache{
		items: make(map[string]PreRenderedItem, size),
	}
}

func (c *preRenderCache) Set(ctx context.Context, i PreRenderedItem) {
	c.mu.Lock()
	c.items[i.Path] = i
	c.mu.Unlock()
}

func (c *preRenderCache) Get(ctx context.Context, path string) (PreRenderedItem, bool) {
	c.mu.Lock()
	i, ok := c.items[path]
	c.mu.Unlock()
	return i, ok
}
