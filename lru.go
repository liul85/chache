package chache

import "container/list"

// Cache is the LRU cache
type Cache struct {
	capacity  int64
	used      int64
	cache     map[string]*list.Element
	l         *list.List
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

// Value is any type
type Value interface {
	Len() int64
}

//New can be used to init a cache
func New(capacity int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		capacity:  capacity,
		l:         list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

//Get return kv if a key exist else nil
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.l.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return nil, false
}

//RemoveOldest removes oldest kv from cache
func (c *Cache) RemoveOldest() {
	ele := c.l.Back()

	if ele != nil {
		c.l.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.used -= int64(len(kv.key)) + kv.value.Len()
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

//Add method adds a key value pair to the cache
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.l.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.used += value.Len() - kv.value.Len()
		kv.value = value
	} else {
		ele := c.l.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.used += int64(len(key)) + value.Len()
	}

	for c.capacity != 0 && c.used > c.capacity {
		c.RemoveOldest()
	}
}

//Len returns the number of cache entries
func (c *Cache) Len() int {
	return c.l.Len()
}
