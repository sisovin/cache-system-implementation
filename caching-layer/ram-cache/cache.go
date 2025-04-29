package ramcache

import (
	"sync"
	"time"
)

type Cache struct {
	data      sync.Map
	ttl       time.Duration
	stats     CacheStats
	evictChan chan struct{}
}

type CacheStats struct {
	Hits   int
	Misses int
}

func NewCache(ttl time.Duration) *Cache {
	cache := &Cache{
		ttl:       ttl,
		evictChan: make(chan struct{}),
	}
	go cache.evictExpiredItems()
	return cache
}

func (c *Cache) Set(key string, value interface{}) {
	c.data.Store(key, value)
}

func (c *Cache) Get(key string) (interface{}, bool) {
	value, ok := c.data.Load(key)
	if ok {
		c.stats.Hits++
	} else {
		c.stats.Misses++
	}
	return value, ok
}

func (c *Cache) evictExpiredItems() {
	ticker := time.NewTicker(c.ttl)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.data.Range(func(key, value interface{}) bool {
				c.data.Delete(key)
				return true
			})
		case <-c.evictChan:
			return
		}
	}
}

func (c *Cache) StopEviction() {
	close(c.evictChan)
}

func (c *Cache) Stats() CacheStats {
	return c.stats
}
