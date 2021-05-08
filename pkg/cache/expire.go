package cache

import (
	"context"
	"sync"
	"time"
)

type Expire struct {
	// The duration while an item is cached.
	ItemTTL time.Duration

	once  sync.Once
	mutex sync.RWMutex
	size  int
	items map[string]*memItem
	queue []*memItem
}

func (c *Expire) Get(ctx context.Context, key string) (Item, bool) {
	c.once.Do(c.init)
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	i, isCached := c.items[key]
	if !isCached || i.expiresAt.Before(time.Now()) {
		return nil, false
	}
	return i.value, true
}

func (c *Expire) Set(ctx context.Context, key string, i Item) {
	c.once.Do(c.init)
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if i, isCached := c.items[key]; isCached {
		c.del(i)
	}

	c.expire()
	c.add(&memItem{
		key:       key,
		expiresAt: time.Now().Add(c.ItemTTL),
		value:     i,
	})
}

func (c *Expire) Del(ctx context.Context, key string) {
	c.once.Do(c.init)
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if i, isCached := c.items[key]; isCached {
		c.del(i)
		c.expire()
	}
}

func (c *Expire) Len() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return len(c.items)
}

func (c *Expire) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.size
}

func (c *Expire) init() {
	c.items = make(map[string]*memItem)
}

func (c *Expire) add(i *memItem) {
	c.items[i.key] = i
	c.queue = append(c.queue, i)
	c.size += i.value.Size()
}

func (c *Expire) del(i *memItem) {
	if i.value != nil {
		delete(c.items, i.key)
		c.size -= i.value.Size()
		i.value = nil
	}
}

func (c *Expire) expire() {
	now := time.Now()

	i := 0
	for i < len(c.queue) && c.queue[i].isExpired(now) {
		c.del(c.queue[i])
		i++
	}

	copy(c.queue, c.queue[i:])
	c.queue = c.queue[:len(c.queue)-i]
}

type memItem struct {
	key       string
	expiresAt time.Time
	value     Item
}

func (i *memItem) isExpired(now time.Time) bool {
	return i.value == nil || i.expiresAt.Before(now)
}
