package aggregation

import (
	"sync"
	"time"
)

type ttlCacheItem struct {
	value     any
	expiresAt time.Time
}

type ttlCache struct {
	mu    sync.Mutex
	items map[string]ttlCacheItem
	ttl   time.Duration
}

func newTTLCache(ttl time.Duration) *ttlCache {
	return &ttlCache{
		items: make(map[string]ttlCacheItem),
		ttl:   ttl,
	}
}

func (c *ttlCache) Get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if !ok {
		return nil, false
	}
	if time.Now().After(item.expiresAt) {
		delete(c.items, key)
		return nil, false
	}
	return item.value, true
}

func (c *ttlCache) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = ttlCacheItem{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
}
