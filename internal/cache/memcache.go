package cache

import (
	"context"
	"sync"
	"time"
)

type InMemoryCache struct {
	mu    sync.RWMutex
	items map[string]item
}

type item struct {
	val    []byte
	expiry time.Time
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		items: make(map[string]item),
	}
}

// Get retrieves the value associated with the given key, returning a copy of the value if it exists and is not expired.
func (c *InMemoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	c.mu.RLock()
	it, ok := c.items[key]
	c.mu.RUnlock()

	if !ok || time.Now().After(it.expiry) {
		return nil, nil
	}

	// Return a copy to avoid mutation
	cp := make([]byte, len(it.val))
	copy(cp, it.val)
	return cp, nil
}

// Set stores a key-value pair in the cache with a specified time-to-live (ttl) duration.
func (c *InMemoryCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	c.mu.Lock()
	c.items[key] = item{
		val:    append([]byte{}, value...),
		expiry: time.Now().Add(ttl),
	}
	c.mu.Unlock()
	return nil
}
