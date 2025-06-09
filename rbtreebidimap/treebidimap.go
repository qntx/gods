// Package treebidimap implements a bidirectional map backed by two red-black tree.
//
// This structure guarantees that the map will be in both ascending key and value order.
//
// Other than key and value ordering, the goal with this structure is to avoid duplication of elements, which can be significant if contained elements are large.
//
// A bidirectional map, or hash bag, is an associative data structure in which the (key,value) pairs form a one-to-one correspondence.
// Thus the binary relation is functional in each direction: value can also act as a key to key.
// A pair (a,b) thus provides a unique coupling between 'a' and 'b' so that 'b' can be found when 'a' is used as a key and 'a' can be found when 'b' is used as a key.
//
// Structure is not thread safe.
//
// Reference: https://en.wikipedia.org/wiki/Bidirectional_map
package rbtreebidimap

import (
	"fmt"

	"strings"

	"github.com/qntx/gods/cmp"
	"github.com/qntx/gods/rbtree"
)

// Map holds the elements in two red-black trees.
type Map[K, V comparable] struct {
	forwardMap rbtree.Tree[K, V]
	inverseMap rbtree.Tree[V, K]
}

// New instantiates a bidirectional map.
func New[K, V cmp.Ordered]() *Map[K, V] {
	return &Map[K, V]{
		forwardMap: *rbtree.New[K, V](),
		inverseMap: *rbtree.New[V, K](),
	}
}

// NewWith instantiates a bidirectional map.
func NewWith[K, V comparable](keyComparator cmp.Comparator[K], valueComparator cmp.Comparator[V]) *Map[K, V] {
	return &Map[K, V]{
		forwardMap: *rbtree.NewWith[K, V](keyComparator),
		inverseMap: *rbtree.NewWith[V, K](valueComparator),
	}
}

// Put inserts element into the map.
func (m *Map[K, V]) Put(key K, value V) {
	if v, ok := m.forwardMap.Get(key); ok {
		m.inverseMap.Delete(v)
	}
	if k, ok := m.inverseMap.Get(value); ok {
		m.forwardMap.Delete(k)
	}
	m.forwardMap.Put(key, value)
	m.inverseMap.Put(value, key)
}

// Get searches the element in the map by key and returns its value or nil if key is not found in map.
// Second return parameter is true if key was found, otherwise false.
func (m *Map[K, V]) Get(key K) (value V, found bool) {
	return m.forwardMap.Get(key)
}

// GetKey searches the element in the map by value and returns its key or nil if value is not found in map.
// Second return parameter is true if value was found, otherwise false.
func (m *Map[K, V]) GetKey(value V) (key K, found bool) {
	return m.inverseMap.Get(value)
}

// Remove removes the element from the map by key.
func (m *Map[K, V]) Remove(key K) {
	if v, found := m.forwardMap.Get(key); found {
		m.forwardMap.Delete(key)
		m.inverseMap.Delete(v)
	}
}

// Empty returns true if map does not contain any elements
func (m *Map[K, V]) Empty() bool {
	return m.Len() == 0
}

// Len returns number of elements in the map.
func (m *Map[K, V]) Len() int {
	return m.forwardMap.Len()
}

// Keys returns all keys (ordered).
func (m *Map[K, V]) Keys() []K {
	return m.forwardMap.Keys()
}

// Values returns all values (ordered).
func (m *Map[K, V]) Values() []V {
	return m.inverseMap.Keys()
}

// Clear removes all elements from the map.
func (m *Map[K, V]) Clear() {
	m.forwardMap.Clear()
	m.inverseMap.Clear()
}

// String returns a string representation of container
func (m *Map[K, V]) String() string {
	str := "TreeBidiMap\nmap["
	it := m.forwardMap.Iterator()
	for it.Next() {
		str += fmt.Sprintf("%v:%v ", it.Key(), it.Value())
	}
	return strings.TrimRight(str, " ") + "]"
}
