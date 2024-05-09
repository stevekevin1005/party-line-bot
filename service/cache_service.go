package service

import (
	"sync"
	"time"
)

type Cache struct {
	data *sync.Map
}

type cacheItem struct {
	value  interface{}
	expiry time.Time
}

func NewCache() *Cache {
	return &Cache{
		data: &sync.Map{},
	}
}

func (c *Cache) Set(key string, value interface{}, expiration time.Duration) {
	expiry := time.Now().Add(expiration)
	c.data.Store(key, cacheItem{
		value:  value,
		expiry: expiry,
	})
	// 启动定时器，在过期时间后自动删除缓存项
	time.AfterFunc(expiration, func() {
		c.Delete(key)
	})
}

func (c *Cache) Get(key string) (interface{}, bool) {
	item, ok := c.data.Load(key)
	if !ok {
		return nil, false
	}

	cachedItem := item.(cacheItem)
	if time.Now().After(cachedItem.expiry) {
		// 如果缓存项已过期，则删除并返回false
		c.Delete(key)
		return nil, false
	}

	return cachedItem.value, true
}

func (c *Cache) Delete(key string) {
	c.data.Delete(key)
}
