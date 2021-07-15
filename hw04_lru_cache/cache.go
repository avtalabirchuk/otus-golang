package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш по ключу.
	Get(key Key) (interface{}, bool)     // Получить значение из кэша по ключу.
	Clear()                              // Очистить кэш.
}

type lruCache struct {
	capacity int               // емкость
	queue    List              // очередь последних элементов
	items    map[Key]*ListItem // словарь
	mutx     sync.Mutex        // Блокировка доступа на запись
}

// Задать значение ключа  и положить его в кеш.
func (lru *lruCache) Set(key Key, value interface{}) bool {
	lru.mutx.Lock()
	defer lru.mutx.Unlock()
	item, keyExists := lru.items[key]
	// Проверка на существование ключа в кеше, если есть то переместить на верх списка
	if keyExists {
		item.Value = cacheItem{Key: key, Value: value}
		lru.queue.MoveToFront(item)
	} else {
		// Если нет, содается новый элемент очереди и помещается в начало списка
		item = lru.queue.PushFront(cacheItem{Key: key, Value: value})
		lru.items[key] = item
		// если длинна очереди становится больше емкости, последний элемент очереди удаляется
		if lru.queue.Len() > lru.capacity {
			lastItem := lru.queue.Back()
			lru.queue.Remove(lastItem)
			delete(lru.items, Key(lastItem.Value.(cacheItem).Key))
		}
	}
	return keyExists
}

// Получаем значение из кэша.
func (lru *lruCache) Get(key Key) (interface{}, bool) {
	lru.mutx.Lock()
	defer lru.mutx.Unlock()
	item, keyExist := lru.items[key]
	if keyExist {
		lru.queue.MoveToFront(item)
		return item.Value.(cacheItem).Value, true
	}
	// ничего не нашли
	return nil, false
}

func (lru *lruCache) Clear() {
	lru.items = make(map[Key]*ListItem)
	lru.queue = &list{}
}

type cacheItem struct {
	Key   Key
	Value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		items:    make(map[Key]*ListItem),
		queue:    &list{},
	}
}
