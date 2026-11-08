package lru

import (
	"lru/internal/linked_list"
	"sync"
)

type LRU struct {
	ll    *linked_list.LinkedList
	store map[int64]*linked_list.Node
	mu    sync.Mutex
}

type ShardedLRU struct {
	lruList []*LRU
}

func NewLRU(maxSize int64) *LRU {
	ll := linked_list.NewLinkedList(maxSize)
	store := make(map[int64]*linked_list.Node)

	return &LRU{
		ll:    ll,
		store: store,
	}
}

func NewShardedLRU(maxSize int64) *ShardedLRU {
	ll := make([]*LRU, 0, 16)
	var lru *LRU
	for i := 0; i < 16; i++ {
		lru = NewLRU(int64(maxSize / 16))
		ll = append(ll, lru)
	}
	sl := ShardedLRU{
		lruList: ll,
	}
	return &sl
}

func (sl *ShardedLRU) getShard(key int64) *LRU {
	if key < 0 {
		key = -key
	}
	index := key % 16
	return sl.lruList[index]
}

func (sl *ShardedLRU) Add(key int64, data any) {
	sl.getShard(key).Add(key, data)
}

func (sl *ShardedLRU) Read(key int64) *linked_list.Node {
	return sl.getShard(key).Read(key)
}

func (l *LRU) Add(key int64, data any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	node, exists := l.store[key]
	if exists {
		node.Data = data
		l.ll.MoveToFront(node)
	} else {
		if l.ll.Size >= l.ll.MaxSize-2 {
			dk := l.ll.Pop()
			delete(l.store, dk)
		}
		node = &linked_list.Node{
			Id:   key,
			Data: data,
		}
		l.ll.Insert(node)
		l.store[key] = node
	}
}

func (l *LRU) Read(key int64) *linked_list.Node {
	l.mu.Lock()
	defer l.mu.Unlock()

	node, exists := l.store[key]
	if exists {
		l.ll.MoveToFront(node)

		return node
	}
	return nil
}
