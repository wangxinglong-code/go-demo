package cache

//参照https://github.com/patrickmn/go-cache

import (
	"runtime"
	"sync"
	"time"
)

const (
	DefaultExpiration      time.Duration = 5 * time.Minute
	DefaultCleanupInterval time.Duration = 10 * time.Minute
	NoExpiration           time.Duration = -1
)

var cacheClient *cache

type Cache struct {
	*cache
}

type cache struct {
	defaultExpiration time.Duration
	items             map[string]Item
	mu                sync.RWMutex
	onEvicted         func(string, interface{})
	watchmen          *watchmen
}

type Item struct {
	Value      interface{}
	Expiration int64
}

type watchmen struct {
	Interval time.Duration
	stop     chan bool
}

type keyAndValue struct {
	key   string
	value interface{}
}

//初始化缓存
func MemCache() *Cache {
	items := make(map[string]Item)
	return newCacheWithWatchmen(DefaultExpiration, DefaultCleanupInterval, items)
}

//创建守护者
func newCacheWithWatchmen(de time.Duration, ci time.Duration, m map[string]Item) *Cache {
	c := newCache(de, m)
	C := &Cache{c}
	if ci > 0 {
		runWatchmen(c, ci)
		runtime.SetFinalizer(C, stopWatchmen)
	}
	return C
}

//创建缓存
func newCache(de time.Duration, m map[string]Item) *cache {
	if de == 0 {
		de = -1
	}
	c := &cache{
		defaultExpiration: de,
		items:             m,
	}
	return c
}

//守护者启动
func runWatchmen(c *cache, ci time.Duration) {
	j := &watchmen{
		Interval: ci,
		stop:     make(chan bool),
	}
	c.watchmen = j
	go j.Run(c)
}

//守护者停止
func stopWatchmen(c *Cache) {
	c.watchmen.stop <- true
}

//守护者启动定时器
func (j *watchmen) Run(c *cache) {
	ticker := time.NewTicker(j.Interval)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-j.stop:
			ticker.Stop()
			return
		}
	}
}

// 删除所有过期缓存
func (c *cache) DeleteExpired() {
	var evictedItems []keyAndValue
	now := time.Now().UnixNano()
	c.mu.Lock()
	for k, v := range c.items {
		// "Inlining" of expired
		if v.Expiration > 0 && now > v.Expiration {
			ov, evicted := c.delete(k)
			if evicted {
				evictedItems = append(evictedItems, keyAndValue{k, ov})
			}
		}
	}
	c.mu.Unlock()
	for _, v := range evictedItems {
		c.onEvicted(v.key, v.value)
	}
}

//删除key
func (c *cache) delete(k string) (interface{}, bool) {
	if c.onEvicted != nil {
		if v, found := c.items[k]; found {
			delete(c.items, k)
			return v.Value, true
		}
	}
	delete(c.items, k)
	return nil, false
}

// 强制删除key
func (c *cache) Delete(k string) {
	c.mu.Lock()
	v, evicted := c.delete(k)
	c.mu.Unlock()
	if evicted {
		c.onEvicted(k, v)
	}
}

//设置缓存
func (c *cache) Set(k string, x interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	c.mu.Lock()
	c.items[k] = Item{
		Value:      x,
		Expiration: e,
	}
	c.mu.Unlock()
}

//获取缓存
func (c *cache) Get(k string) (interface{}, bool) {
	c.mu.RLock()
	item, found := c.items[k]
	if !found {
		c.mu.RUnlock()
		return nil, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			c.mu.RUnlock()
			return nil, false
		}
	}
	c.mu.RUnlock()
	return item.Value, true
}
