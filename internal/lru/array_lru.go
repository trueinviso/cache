package lru

import (
	"cache/internal/array_linked_list"
	"sync"
)

type ArrayLRU struct {
	ll    *array_linked_list.LinkedList
	store map[int64]int64
	mu    sync.Mutex
}

type ShardedArrayLRU struct {
	lruList []*ArrayLRU
}

func NewArrayLRU(maxSize int64) *ArrayLRU {
	ll := array_linked_list.NewLinkedList(maxSize)
	store := make(map[int64]int64)

	return &ArrayLRU{
		ll:    ll,
		store: store,
	}
}

func NewShardedArrayLRU(maxSize int64) *ShardedArrayLRU {
	ll := make([]*ArrayLRU, 0, 16)
	var lru *ArrayLRU
	for i := 0; i < 16; i++ {
		lru = NewArrayLRU(int64(maxSize / 16))
		ll = append(ll, lru)
	}
	sl := ShardedArrayLRU{
		lruList: ll,
	}
	return &sl
}

func (sl *ShardedArrayLRU) getShard(key int64) *ArrayLRU {
	if key < 0 {
		key = -key
	}
	index := key % 16
	return sl.lruList[index]
}

func (sl *ShardedArrayLRU) Add(key int64, data any) {
	sl.getShard(key).Add(key, data)
}

func (sl *ShardedArrayLRU) Read(key int64) *array_linked_list.Node {
	return sl.getShard(key).Read(key)
}

func (l *ArrayLRU) Add(key int64, data any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	idx, exists := l.store[key]
	if exists {
		l.ll.NodeList[idx].Data = data
		l.ll.MoveToFront(idx)
	} else {
		node := &array_linked_list.Node{
			Id:   key,
			Data: data,
		}

		var idx int64
		if l.ll.Size >= l.ll.MaxSize-2 {
			idx, oldKey := l.ll.AddToFullCache(node)
			delete(l.store, oldKey)
			l.store[key] = idx
		} else {
			idx = l.ll.Insert(node)
			l.store[key] = idx
		}
	}
}

func (l *ArrayLRU) Read(key int64) *array_linked_list.Node {
	l.mu.Lock()
	defer l.mu.Unlock()

	idx, exists := l.store[key]
	if exists {
		node := &l.ll.NodeList[idx]
		l.ll.MoveToFront(idx)

		return node
	}
	return nil
}
