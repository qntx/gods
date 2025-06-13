// Package linkedhashmap is a map that preserves insertion-order.
//
// It is backed by a hash table to store values and doubly-linked list to store ordering.
//
// Structure is not thread safe.
//
// Reference: http://en.wikipedia.org/wiki/Associative_array
package linkedhashmap

import (
	"bytes"
	"cmp"
	"container/list"
	"encoding/json"
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/qntx/gods/container"
)

var _ container.Map[int, int] = (*Map[int, int])(nil)
var _ json.Marshaler = (*Map[int, int])(nil)
var _ json.Unmarshaler = (*Map[int, int])(nil)

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

	return
}

// Has returns true if the specified key is present in the map, false otherwise.
func (m *Map[K, V]) Has(key K) bool {
	_, ok := m.table[key]

	return ok
}

// Delete removes the element from the map by key.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (m *Map[K, V]) Delete(key K) (value V, found bool) {
	if elem, contains := m.table[key]; contains {
		delete(m.table, key)
		m.ordering.Remove(elem.listElem) // O(1) removal from list

		return elem.value, true
	}

	return
}

func (m *Map[K, V]) IsEmpty() bool {
	return m.Len() == 0
}

func (m *Map[K, V]) Len() int {
	return m.ordering.Len() // Use Len() for container/list
}

func (m *Map[K, V]) Keys() []K {
	keys := make([]K, 0, m.ordering.Len())
	for e := m.ordering.Front(); e != nil; e = e.Next() {
		keys = append(keys, e.Value.(K)) // Type assertion needed
	}

	return keys
}

func (m *Map[K, V]) Values() []V {
	values := make([]V, m.Len())
	count := 0

	for _, v := range m.Iter() {
		values[count] = v
		count++
	}

	return values
}

func (m *Map[K, V]) Entries() (keys []K, values []V) {
	return m.Keys(), m.Values()
}

func (m *Map[K, V]) ToSlice() []V {
	return m.Values()
}

func (m *Map[K, V]) Iter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for e := m.ordering.Front(); e != nil; e = e.Next() {
			key := e.Value.(K)       // Key from the ordered list
			elemData := m.table[key] // Get the element data (value + listElem pointer)

			value := elemData.value // Extract the actual value
			if !yield(key, value) {
				return
			}
		}
	}
}

func (m *Map[K, V]) Clear() {
	// For Go versions < 1.21, or if specific deallocation logic per element was needed:
	// for k := range m.table {
	// 	 delete(m.table, k)
	// }
	clear(m.table)    // Efficiently clears the map (Go 1.21+)
	m.ordering.Init() // Reinitializes the list, making it empty
}

func (m *Map[K, V]) Clone() container.Map[K, V] {
	clone := New[K, V]()
	for k, v := range m.Iter() {
		clone.Put(k, v)
	}

	return clone
}

// MarshalJSON @implements json.Marshaler.
func (m *Map[K, V]) MarshalJSON() ([]byte, error) {
	var b []byte
	buf := bytes.NewBuffer(b)

	buf.WriteRune('{')

	count := 0
	lastIndex := m.Len() - 1

	for k, v := range m.Iter() {
		km, err := json.Marshal(k)
		if err != nil {
			return nil, err
		}

		buf.Write(km)

		buf.WriteRune(':')

		vm, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		buf.Write(vm)

		if count != lastIndex {
			buf.WriteRune(',')
		}

		count++
	}

	buf.WriteRune('}')

	return buf.Bytes(), nil
}

// UnmarshalJSON @implements json.Unmarshaler.
func (m *Map[K, V]) UnmarshalJSON(data []byte) error {
	elements := make(map[K]V)

	err := json.Unmarshal(data, &elements)
	if err != nil {
		return err
	}

	index := make(map[K]int)

	var keys []K
	for key := range elements {
		keys = append(keys, key)
		esc, _ := json.Marshal(key)
		index[key] = bytes.Index(data, esc)
	}

	byIndex := func(key1, key2 K) int {
		return cmp.Compare(index[key1], index[key2])
	}

	slices.SortFunc(keys, byIndex)

	m.Clear()

	for _, key := range keys {
		m.Put(key, elements[key])
	}

	return nil
}

// String returns a string representation of container.
func (m *Map[K, V]) String() string {
	str := "LinkedHashMap\nmap["

	for k, v := range m.Iter() {
		str += fmt.Sprintf("%v:%v ", k, v)
	}

	return strings.TrimRight(str, " ") + "]"
}
