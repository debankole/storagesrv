package collections

import "sync"

type Node[K comparable, V any] struct {
	key   K
	value V
	next  *Node[K, V]
	prev  *Node[K, V]
}

// OrderedMap is a map that maintains the order of items and concurrency safe
type OrderedMap[K comparable, V any] struct {
	m       map[K]*Node[K, V]
	head    *Node[K, V]
	tail    *Node[K, V]
	rwMutex sync.RWMutex
}

type kv[K comparable, V any] struct {
	Key   K
	Value V
}

func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		m: make(map[K]*Node[K, V]),
	}
}

func (o *OrderedMap[K, V]) AddItem(key K, value V) {
	o.rwMutex.Lock()
	defer o.rwMutex.Unlock()

	if node, exists := o.m[key]; exists {
		node.value = value
		return
	}

	node := &Node[K, V]{key: key, value: value}
	o.m[key] = node

	if o.tail != nil {
		o.tail.next = node
		node.prev = o.tail
	} else {
		o.head = node
	}
	o.tail = node
}

func (o *OrderedMap[K, V]) RemoveItem(key K) {
	o.rwMutex.Lock()
	defer o.rwMutex.Unlock()

	if node, exists := o.m[key]; exists {
		delete(o.m, key)
		if node.prev != nil {
			node.prev.next = node.next
		} else {
			o.head = node.next
		}
		if node.next != nil {
			node.next.prev = node.prev
		} else {
			o.tail = node.prev
		}
	}
}

func (o *OrderedMap[K, V]) GetItem(key K) (V, bool) {
	o.rwMutex.RLock()
	defer o.rwMutex.RUnlock()

	var val V
	var ok bool
	var node *Node[K, V]
	if node, ok = o.m[key]; ok {
		val = node.value
	}

	return val, ok
}

func (o *OrderedMap[K, V]) GetAllItems() []kv[K, V] {
	o.rwMutex.RLock()
	defer o.rwMutex.RUnlock()

	var items []kv[K, V]
	for node := o.head; node != nil; node = node.next {
		items = append(items, kv[K, V]{Key: node.key, Value: node.value})
	}
	return items
}
