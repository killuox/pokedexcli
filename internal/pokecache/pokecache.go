package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	mu   sync.Mutex
	data map[string]cacheEntry
}

func NewCache(cleanupInterval time.Duration) *Cache {
	cache := &Cache{
		data: make(map[string]cacheEntry),
	}
	go cache.readLoop(cleanupInterval)

	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = cacheEntry{
		val:       val,
		createdAt: time.Now(),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.data[key]

	if ok {
		return v.val, true
	} else {
		return nil, false
	}
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

func (c *Cache) readLoop(cleanupInterval time.Duration) {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()
	for range ticker.C {
		for key, d := range c.data {

			if d.createdAt.Before(time.Now()) {
				c.Delete(key)
			}
		}
	}
}
