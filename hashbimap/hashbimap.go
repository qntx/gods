// Package hashbidimap implements a bidirectional map backed by two hashmaps.
//
// A bidirectional map, or hash bag, is an associative data structure in which the (key,value) pairs form a one-to-one correspondence.
// Thus the binary relation is functional in each direction: value can also act as a key to key.
// A pair (a,b) thus provides a unique coupling between 'a' and 'b' so that 'b' can be found when 'a' is used as a key and 'a' can be found when 'b' is used as a key.
//
// Elements are unordered in the map.
//
// Structure is not thread safe.
//
// Reference: https://en.wikipedia.org/wiki/Bidirectional_map
package hashbimap

import (
	"encoding/json"
	"fmt"
	"iter"

	"github.com/qntx/gods/cmp"
	"github.com/qntx/gods/container"
	"github.com/qntx/gods/hashmap"
)

var _ container.BiMap[string, int] = (*Map[string, int])(nil)
var _ json.Marshaler = (*Map[string, int])(nil)
var _ json.Unmarshaler = (*Map[string, int])(nil)

type Map[K cmp.Ordered, V cmp.Ordered] struct {
	fwd hashmap.Map[K, V]
	inv hashmap.Map[V, K]
}

func New[K, V cmp.Ordered]() *Map[K, V] {
	return &Map[K, V]{
		fwd: *hashmap.New[K, V](),
		inv: *hashmap.New[V, K](),
	}
}

func NewWith[K, V cmp.Ordered](kcmp cmp.Comparator[K], vcmp cmp.Comparator[V], size int) *Map[K, V] {
	return &Map[K, V]{
		fwd: *hashmap.NewWith[K, V](kcmp, size),
		inv: *hashmap.NewWith[V, K](vcmp, size),
	}
}

// Put inserts element into the map.
func (m *Map[K, V]) Put(key K, value V) {
	if v, ok := m.fwd.Get(key); ok {
		m.inv.Delete(v)
	}

	if k, ok := m.inv.Get(value); ok {
		m.fwd.Delete(k)
	}

	m.fwd.Put(key, value)
	m.inv.Put(value, key)
}

// Get searches the element in the map by key and returns its value or nil if key is not found in map.
// Second return parameter is true if key was found, otherwise false.
func (m *Map[K, V]) Get(key K) (value V, found bool) {
	return m.fwd.Get(key)
}

// GetKey searches the element in the map by value and returns its key or nil if value is not found in map.
// Second return parameter is true if value was found, otherwise false.
func (m *Map[K, V]) GetKey(value V) (key K, found bool) {
	return m.inv.Get(value)
}

func (m *Map[K, V]) Has(k K) bool {
	return m.fwd.Has(k)
}

func (m *Map[K, V]) HasValue(v V) bool {
	return m.inv.Has(v)
}

// Delete removes the element from the map by key.
func (m *Map[K, V]) Delete(key K) (v V, ok bool) {
	if value, found := m.fwd.Get(key); found {
		m.fwd.Delete(key)
		m.inv.Delete(value)

		return value, true
	}

	return v, false
}

func (m *Map[K, V]) DeleteValue(v V) (k K, ok bool) {
	if k, ok := m.inv.Get(v); ok {
		m.fwd.Delete(k)
		m.inv.Delete(v)

		return k, true
	}

	return k, false
}

func (m *Map[K, V]) Iter() iter.Seq2[K, V] {
	return m.fwd.Iter()
}

// IsEmpty returns true if map does not contain any elements.
func (m *Map[K, V]) IsEmpty() bool {
	return m.Len() == 0
}

// Len returns number of elements in the map.
func (m *Map[K, V]) Len() int {
	return m.fwd.Len()
}

// Keys returns all keys (random order).
func (m *Map[K, V]) Keys() []K {
	return m.fwd.Keys()
}

// Values returns all values (random order).
func (m *Map[K, V]) Values() []V {
	return m.inv.Keys()
}

func (m *Map[K, V]) ToSlice() []V {
	return m.fwd.ToSlice()
}

func (m *Map[K, V]) Entries() ([]K, []V) {
	return m.fwd.Entries()
}

// Clear removes all elements from the map.
func (m *Map[K, V]) Clear() {
	m.fwd.Clear()
	m.inv.Clear()
}

func (m *Map[K, V]) Clone() container.Map[K, V] {
	return &Map[K, V]{
		fwd: *(m.fwd.Clone().(*hashmap.Map[K, V])),
		inv: *(m.inv.Clone().(*hashmap.Map[V, K])),
	}
}

// MarshalJSON outputs the JSON representation of the map.
func (m *Map[K, V]) MarshalJSON() ([]byte, error) {
	return m.fwd.MarshalJSON()
}

// UnmarshalJSON populates the map from the input JSON representation.
func (m *Map[K, V]) UnmarshalJSON(data []byte) error {
	var elems map[K]V

	err := json.Unmarshal(data, &elems)
	if err != nil {
		return err
	}

	m.Clear()

	for k, v := range elems {
		m.Put(k, v)
	}

	return nil
}

// String returns a string representation of container.
func (m *Map[K, V]) String() string {
	str := "HashBidiMap\n"
	str += fmt.Sprintf("%v", m.fwd)

	return str
}
