// Package linkedhashmap is a map that preserves insertion-order.
//
// It is backed by a hash table to store values and doubly-linked list to store ordering.
//
// Structure is not thread safe.
//
// Reference: http://en.wikipedia.org/wiki/Associative_array
package linkedhashmap

import (
	"container/list"
	"fmt"
	"strings"
)

// element is a helper struct to store the value and a pointer to the list element.
// This allows for O(1) removal from the list.
type element[V any] struct {
	value    V
	listElem *list.Element
}

// Map holds the elements in a regular hash table, and uses a doubly-linked list from container/list to store key ordering.
type Map[K comparable, V any] struct {
	table    map[K]element[V] // Stores key -> {value, *list.Element}
	ordering *list.List       // Stores keys (K) in insertion order
}

// New instantiates a linked-hash-map.
func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		table:    make(map[K]element[V]),
		ordering: list.New(),
	}
}

// Put inserts key-value pair into the map.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (m *Map[K, V]) Put(key K, value V) {
	if elem, contains := m.table[key]; contains {
		elem.value = value
		m.table[key] = elem
		// Note: If LRU behavior (move to back on access/update) is desired, add:
		// m.ordering.MoveToBack(elem.listElem)
		// For strict insertion-order-only, existing elements' order doesn't change on update.
	} else {
		listElement := m.ordering.PushBack(key) // Store the key in the list
		m.table[key] = element[V]{value: value, listElem: listElement}
	}
}

// Get searches the element in the map by key and returns its value or nil if key is not found in tree.
// Second return parameter is true if key was found, otherwise false.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (m *Map[K, V]) Get(key K) (value V, found bool) {
	if elem, ok := m.table[key]; ok {
		return elem.value, true
	}

	var zeroV V // Required to return the zero value for V if not found

	return zeroV, false
}

// Remove removes the element from the map by key.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (m *Map[K, V]) Remove(key K) {
	if elem, contains := m.table[key]; contains {
		delete(m.table, key)
		m.ordering.Remove(elem.listElem) // O(1) removal from list
	}
}

// Empty returns true if map does not contain any elements.
func (m *Map[K, V]) Empty() bool {
	return m.Size() == 0
}

// Size returns number of elements in the map.
func (m *Map[K, V]) Size() int {
	return m.ordering.Len() // Use Len() for container/list
}

// Keys returns all keys in-order.
func (m *Map[K, V]) Keys() []K {
	keys := make([]K, 0, m.ordering.Len())
	for e := m.ordering.Front(); e != nil; e = e.Next() {
		keys = append(keys, e.Value.(K)) // Type assertion needed
	}

	return keys
}

// Values returns all values in-order based on the key.
func (m *Map[K, V]) Values() []V {
	values := make([]V, m.Size())
	count := 0

	it := m.Iterator() // Iterator will use the updated Keys() and Get()
	for it.Next() {
		values[count] = it.Value()
		count++
	}

	return values
}

// Clear removes all elements from the map.
func (m *Map[K, V]) Clear() {
	// For Go versions < 1.21, or if specific deallocation logic per element was needed:
	// for k := range m.table {
	// 	 delete(m.table, k)
	// }
	clear(m.table)    // Efficiently clears the map (Go 1.21+)
	m.ordering.Init() // Reinitializes the list, making it empty
}

// String returns a string representation of container.
func (m *Map[K, V]) String() string {
	str := "LinkedHashMap\nmap["
	it := m.Iterator() // Iterator will use the updated Keys() and Get()

	for it.Next() {
		str += fmt.Sprintf("%v:%v ", it.Key(), it.Value())
	}

	return strings.TrimRight(str, " ") + "]"
}
