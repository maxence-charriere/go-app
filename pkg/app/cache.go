package app

import (
	"context"
	"sync"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/cache"
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
func (r PreRenderedItem) Size() int {
	return len(r.Body)
}

// NewPreRenderLRUCache creates an in memory LRU cache that stores items for the
// given duration. If provided, on eviction functions are called when item are
// evicted.
func NewPreRenderLRUCache(size int, itemTTL time.Duration, onEvict ...func(path string, i PreRenderedItem)) PreRenderCache {
	return &preRenderLRUCache{
		LRU: cache.LRU{
			MaxSize: size,
			ItemTTL: itemTTL,
			OnEvict: func(path string, i cache.Item) {
				item := i.(PreRenderedItem)
				for _, fn := range onEvict {
					fn(path, item)
				}
			},
		},
	}

}

type preRenderLRUCache struct {
	cache.LRU
}

func (c *preRenderLRUCache) Get(ctx context.Context, path string) (PreRenderedItem, bool) {
	i, ok := c.LRU.Get(ctx, path)
	if !ok {
		return PreRenderedItem{}, false
	}
	return i.(PreRenderedItem), true
}

func (c *preRenderLRUCache) Set(ctx context.Context, i PreRenderedItem) {
	c.LRU.Set(ctx, i.Path, i)
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
