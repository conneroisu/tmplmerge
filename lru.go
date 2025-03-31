package twerge

import (
	"sync"
)

// Make creates a new LRU cache
func Make(maxCapacity int) ICache {
	head := &node{}
	tail := &node{}
	tail.next = head
	head.prev = tail
	return &lru{
		maxCapacity: maxCapacity,
		capacity:    0,
		cache:       make(map[string]*node),
		head:        head,
		tail:        tail,
	}
}

// ICache is the interface for a LRU cache
type ICache interface {
	Get(string) string
	Set(string, string)
}

type node struct {
	key  string
	val  string
	prev *node
	next *node
}

type lru struct {
	maxCapacity int
	capacity    int
	cache       map[string]*node
	head        *node
	tail        *node
	cacheMutex  sync.RWMutex
	listMutex   sync.Mutex
}

func (l *lru) Get(key string) string {
	l.cacheMutex.RLock()
	n := l.cache[key]
	if n == nil {
		l.cacheMutex.RUnlock()
		return ""
	}
	l.cacheMutex.RUnlock()

	l.listMutex.Lock()
	l.remove(n)
	l.insertRight(n)
	l.listMutex.Unlock()

	return n.val
}

func (l *lru) Set(key, value string) {
	l.cacheMutex.Lock()
	if n := l.cache[key]; n != nil {
		l.remove(n)
	}
	n := &node{key: key, val: value}
	l.cache[key] = n
	l.cacheMutex.Unlock()
	l.listMutex.Lock()
	l.insertRight(n)
	l.listMutex.Unlock()
	// evict

	l.listMutex.Lock()
	if l.capacity > l.maxCapacity {
		delete(l.cache, l.tail.next.key)
		l.remove(l.tail.next)
	}
	l.listMutex.Unlock()

}

func (l *lru) insertRight(n *node) {
	prev := l.head.prev
	prev.next = n
	n.prev = prev
	n.next = l.head
	l.head.prev = n
}

func (l *lru) remove(n *node) {
	prev := n.prev
	nxt := n.next
	prev.next = nxt
	nxt.prev = prev
	n.prev = nil
	n.next = nil
	l.capacity--
}
