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
	mutex    sync.Mutex
	interval time.Duration
	entries  map[string]cacheEntry
}

func NewCache(interval time.Duration) *Cache {
	var c Cache
	c.interval = interval
	c.entries = make(map[string]cacheEntry) // maps in go are like dicts in python
	go c.reapLoop()
	return &c
}

func (c *Cache) Add(key string, val []byte) {
	var entry cacheEntry
	entry.createdAt = time.Now()
	entry.val = val
	c.mutex.Lock()
	// also, `defer c.mutex.Unlock()` here defers the unlock operation until end of function
	c.entries[key] = entry
	c.mutex.Unlock() // remove if using `defer`
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop() {

	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		for mapKey, mapVal := range c.entries {
			if time.Since(mapVal.createdAt) > c.interval {
				delete(c.entries, mapKey)
			}
		}
		c.mutex.Unlock()
	}
}
