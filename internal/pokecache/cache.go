package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	mu sync.Mutex
	entries map[string]CacheEntry
}

type CacheEntry struct {
	createdAt time.Time
	val []byte
}

func NewCache(interval time.Duration) *Cache {
	ent := make(map[string]CacheEntry) 
	cache := Cache{entries: ent}
  go cache.reapLoop(interval)
	return &cache
}

func (c *Cache) Add(key string, val []byte) {
  entry := CacheEntry{createdAt: time.Now(), val: val}
	c.mu.Lock()
	c.entries[key] = entry
	c.mu.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
  c.mu.Lock()
  defer c.mu.Unlock()
	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
  ticker := time.NewTicker(interval)
  defer ticker.Stop()

  for range ticker.C {
		c.mu.Lock()
		for key, entry := range c.entries {
			if time.Since(entry.createdAt) > interval {
				fmt.Println("REMOVED AN ENTRY")
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
  }
}


func InitializeCache(duration time.Duration) *Cache {
	return NewCache(duration)
}