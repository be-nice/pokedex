package cache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries    map[string]cacheEntry
	reapTicker *time.Ticker
	done       chan struct{}
	interval   time.Duration
	mutex      sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		entries:    make(map[string]cacheEntry),
		interval:   interval,
		reapTicker: time.NewTicker(interval),
		done:       make(chan struct{}),
	}

	go cache.reapLoop()

	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.entries[key]
	return entry.val, exists
}

func (c *Cache) reapLoop() {
	for {
		select {
		case <-c.reapTicker.C:
			c.reapEntries()
		case <-c.done:
			return
		}
	}
}

func (c *Cache) reapEntries() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	expirationTime := time.Now().Add(-c.interval)
	for key, entry := range c.entries {
		if entry.createdAt.Before(expirationTime) {
			delete(c.entries, key)
		}
	}
}

func (c *Cache) Stop() {
	close(c.done)
	c.reapTicker.Stop()
}
