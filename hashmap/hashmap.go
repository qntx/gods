// Package hashmap implements a map backed by a hash table.
//
// Elements are unordered in the map.
//
// Structure is not thread safe.
//
// Reference: http://en.wikipedia.org/wiki/Associative_array
package hashmap

import (
	"encoding/json"
	"fmt"
	"iter"
	"maps"
	"slices"

	"github.com/qntx/gods/cmp"
	"github.com/qntx/gods/container"
)

var _ container.Map[int, int] = (*Map[int, int])(nil)
var _ json.Marshaler = (*Map[int, int])(nil)
var _ json.Unmarshaler = (*Map[int, int])(nil)

type Map[K comparable, V any] struct {
	m   map[K]V
	cmp cmp.Comparator[K] // Comparator for key sorting
}

// New instantiates a hash map with no key comparator.
func New[K cmp.Ordered, V any]() *Map[K, V] {
	return &Map[K, V]{
		m:   make(map[K]V),
		cmp: cmp.Compare[K],
	}
}

// NewWith instantiates a hash map with the specified size and key comparator.
// If cmp is nil, SortedKeys will panic or produce undefined behavior.
func NewWith[K comparable, V any](cmp cmp.Comparator[K], size int) *Map[K, V] {
	return &Map[K, V]{
		m:   make(map[K]V, size),
		cmp: cmp,
	}
}

// Put associates the specified value with the given key in the map.
// If the key already exists, its value is updated with the new value.
func (m *Map[K, V]) Put(key K, value V) {
	m.m[key] = value
}

// Get retrieves the value associated with the specified key.
// Returns the value and true if the key was found, or the zero value of V and false if not.
func (m *Map[K, V]) Get(key K) (V, bool) {
	value, found := m.m[key]

	return value, found
}

// Has returns true if the specified key is present in the map, false otherwise.
func (m *Map[K, V]) Has(key K) bool {
	_, found := m.m[key]

	return found
}

// Delete removes the key-value pair associated with the specified key.
// Returns the value and true if the key was found and removed, false if the key was not present.
func (m *Map[K, V]) Delete(key K) (value V, found bool) {
	if _, found := m.m[key]; found {
		delete(m.m, key)

		return value, true
	}

	return value, false
}

// Len returns the number of key-value pairs in the map.
func (m *Map[K, V]) Len() int {
	return len(m.m)
}

// IsEmpty returns true if the map contains no key-value pairs.
func (m *Map[K, V]) IsEmpty() bool {
	return len(m.m) == 0
}

// Clear removes all key-value pairs from the map.
func (m *Map[K, V]) Clear() {
	clear(m.m)
}

// UnsortedKeys returns a slice containing all keys in the map (in random order).
func (m *Map[K, V]) UnsortedKeys() []K {
	k := make([]K, 0, len(m.m))
	for key := range m.m {
		k = append(k, key)
	}

	return k
}

// Keys returns a slice containing all keys in the map, sorted using the comparator.
func (m *Map[K, V]) Keys() []K {
	return slices.SortedFunc(maps.Keys(m.m), m.cmp)
}

// UnsortedValues returns a slice containing all values in the map (in random order).
func (m *Map[K, V]) UnsortedValues() []V {
	v := make([]V, 0, len(m.m))
	for _, value := range m.m {
		v = append(v, value)
	}

	return v
}

// Values returns a slice containing all values in the map, sorted by their corresponding keys using m.cmp.
func (m *Map[K, V]) Values() []V {
	k := m.Keys()
	v := make([]V, len(k))

	for i, key := range k {
		v[i] = m.m[key]
	}

	return v
}

// ToSlice returns a slice containing all values in the map (in random order).
func (m *Map[K, V]) ToSlice() []V {
	return m.Values()
}

// UnsortedEntries returns two slices containing all keys and values in the map (in random order).
func (m *Map[K, V]) UnsortedEntries() (keys []K, values []V) {
	k := make([]K, 0, len(m.m))
	v := make([]V, 0, len(m.m))

	for key, value := range m.m {
		k = append(k, key)
		v = append(v, value)
	}

	return k, v
}

// SortedEntries returns two slices containing all keys and values in the map,
// sorted by keys using m.cmp.
func (m *Map[K, V]) Entries() (keys []K, values []V) {
	return m.Keys(), m.Values()
}

// Iter returns an iterator over key-value pairs from m.
// The iteration order is not specified and is not guaranteed to be the same from one call to the next.
func (m *Map[K, V]) Iter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range m.m {
			if !yield(k, v) {
				return
			}
		}
	}
}

// Clone returns a new Map containing all key-value pairs from the current map,
// with the same comparator.
func (m *Map[K, V]) Clone() container.Map[K, V] {
	clone := NewWith[K, V](m.cmp, len(m.m))
	maps.Copy(clone.m, m.m)

	return clone
}

// MarshalJSON outputs the JSON representation of the map.
func (m *Map[K, V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.m)
}

// UnmarshalJSON populates the map from the input JSON representation.
// The comparator is preserved.
func (m *Map[K, V]) UnmarshalJSON(data []byte) error {
	temp := make(map[K]V)
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	clear(m.m)
	maps.Copy(m.m, temp)

	return nil
}

// String returns a string representation of the map.
func (m *Map[K, V]) String() string {
	return fmt.Sprintf("HashMap\n%v", m.m)
}
