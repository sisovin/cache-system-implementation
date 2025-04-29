package cache

import (
	"container/list"
	"sync"
	"time"
)

type CacheItem struct {
	key        string
	value      interface{}
	expiration int64
}

type LRUCache struct {
	capacity int
	items    map[string]*list.Element
	evictList *list.List
	mutex    sync.Mutex
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		evictList: list.New(),
	}
}

func (c *LRUCache) Set(key string, value interface{}, duration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if item, ok := c.items[key]; ok {
		c.evictList.MoveToFront(item)
		item.Value.(*CacheItem).value = value
		item.Value.(*CacheItem).expiration = time.Now().Add(duration).UnixNano()
		return
	}

	if c.evictList.Len() >= c.capacity {
		c.evict()
	}

	item := &CacheItem{
		key:        key,
		value:      value,
		expiration: time.Now().Add(duration).UnixNano(),
	}
	entry := c.evictList.PushFront(item)
	c.items[key] = entry
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if item, ok := c.items[key]; ok {
		if time.Now().UnixNano() > item.Value.(*CacheItem).expiration {
			c.removeElement(item)
			return nil, false
		}
		c.evictList.MoveToFront(item)
		return item.Value.(*CacheItem).value, true
	}
	return nil, false
}

func (c *LRUCache) evict() {
	item := c.evictList.Back()
	if item != nil {
		c.removeElement(item)
	}
}

func (c *LRUCache) removeElement(e *list.Element) {
	c.evictList.Remove(e)
	delete(c.items, e.Value.(*CacheItem).key)
}
