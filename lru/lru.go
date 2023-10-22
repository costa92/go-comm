package lru

import "errors"

// OnEvicted  is used to get a callback when a cache entry is evicted
type OnEvicted[K comparable, V any] func(key K, value V)

// LRU is an LRU cache. It is not safe for concurrent access.
type LRU[K comparable, V any] struct {
	size  int // 当前多大
	ll    *List[K, V]
	cache map[K]*Entry[K, V]
	// optional and executed when an entry is purged.
	OnEvicted OnEvicted[K, V]
}

type Value interface {
	Len() int
}

// NewLRU New Len returns the number of cache entries.
func NewLRU[K comparable, V any](size int, onEvicted OnEvicted[K, V]) (*LRU[K, V], error) {
	if size <= 0 {
		return nil, errors.New("must provide a positive size")
	}
	return &LRU[K, V]{
		size:      size,                     // 允许使用的最大内存
		ll:        NewList[K, V](),          // 双向链表
		cache:     make(map[K]*Entry[K, V]), // 字典
		OnEvicted: onEvicted,                // 某条记录被移除时的回调函数，可以为 nil
	}, nil
}

// Purge is used to completely clear the cache.
func (c *LRU[K, V]) Purge() {
	for k, v := range c.cache {
		if c.OnEvicted != nil {
			c.OnEvicted(k, v.Value)
		}
		delete(c.cache, k)
	}
	c.ll.Init()
}

func (c *LRU[K, V]) Add(key K, value V) bool {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		ele.Value = value
		return false
	}
	// add new entry
	ent := c.ll.PushFront(key, value)
	c.cache[key] = ent
	ok := c.ll.Len() > c.size
	if ok {
		c.removeOldest()
	}
	return ok
}

// Get looks up a key's value from the cache.
func (c *LRU[K, V]) Get(key K) (any, bool) {
	// 从字典中找到对应的双向链表的节点
	if ele, ok := c.cache[key]; ok {
		// 将该节点移动到队头
		c.ll.MoveToFront(ele)
		return ele.Value, true
	}
	return nil, false
}

// RemoveOldest removes the oldest item from the cache.
func (c *LRU[K, V]) RemoveOldest() (key K, value V, ok bool) {
	// 从缓存队列的尾部获取缓存元素
	if ele := c.ll.Back(); ele != nil {
		c.removeElement(ele)
		return ele.Key, ele.Value, true
	}
	return
}

// GetOldest returns the oldest entry
func (c *LRU[K, V]) GetOldest() (key K, value V, ok bool) {
	if ent := c.ll.Back(); ent != nil {
		return ent.Key, ent.Value, true
	}
	return
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
func (c *LRU[K, V]) Keys() []K {
	keys := make([]K, c.ll.Len())
	i := 0
	for ent := c.ll.Back(); ent != nil; ent = ent.PrevEntry() {
		keys[i] = ent.Key
		i++
	}
	return keys
}

// Values returns a slice of the values in the cache, from oldest to newest.
func (c *LRU[K, V]) Values() []V {
	values := make([]V, len(c.cache))
	i := 0
	for ent := c.ll.Back(); ent != nil; ent = ent.PrevEntry() {
		values[i] = ent.Value
		i++
	}
	return values
}

// removeElement is used to remove a given list element from the cache
func (c *LRU[K, V]) removeElement(e *Entry[K, V]) {
	c.ll.Remove(e)
	delete(c.cache, e.Key)
	if c.OnEvicted != nil {
		c.OnEvicted(e.Key, e.Value)
	}
}

// Len returns the number of cache entries.
func (c *LRU[K, V]) Len() int {
	return c.ll.Len()
}

// Resize changes the cache size.
func (c *LRU[K, V]) Resize(size int) (evicted int) {
	diff := c.Len() - size
	if diff < 0 {
		diff = 0
	}
	for i := 0; i < diff; i++ {
		c.removeOldest()
	}
	c.size = size
	return diff
}

// removeOldest removes the oldest item from the cache.
func (c *LRU[K, V]) removeOldest() {
	if ent := c.ll.Back(); ent != nil {
		c.removeElement(ent)
	}
}

func (c *LRU[K, V]) Remove(key K) bool {
	if ent, ok := c.cache[key]; ok {
		c.removeElement(ent)
		return true
	}
	return false
}
