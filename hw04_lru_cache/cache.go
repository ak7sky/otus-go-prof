package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mtx      *sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	cache.mtx.Lock()
	defer cache.mtx.Unlock()

	if li, exist := cache.items[key]; exist {
		li.Value = cacheItem{key: key, value: value}
		cache.queue.MoveToFront(li)
		return exist
	}

	li := cache.queue.PushFront(cacheItem{key: key, value: value})
	cache.items[key] = li
	if cache.queue.Len() > cache.capacity {
		back := cache.queue.Back()
		backKey := back.Value.(cacheItem).key
		cache.queue.Remove(back)
		delete(cache.items, backKey)
	}

	return false
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cache.mtx.Lock()
	defer cache.mtx.Unlock()

	if li, exist := cache.items[key]; exist {
		cache.queue.MoveToFront(li)
		return li.Value.(cacheItem).value, true
	}
	return nil, false
}

func (cache *lruCache) Clear() {
	cache.queue = NewList()
	cache.items = make(map[Key]*ListItem, cache.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		mtx:      new(sync.Mutex),
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
