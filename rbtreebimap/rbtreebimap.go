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
package rbtreebimap

import (
	"encoding/json"
	"fmt"
	"iter"

	"strings"

	"github.com/qntx/gods/cmp"
	"github.com/qntx/gods/container"
	"github.com/qntx/gods/rbtree"
)

var _ container.OrderedBiMap[string, int] = (*Map[string, int])(nil)
var _ json.Marshaler = (*Map[string, int])(nil)
var _ json.Unmarshaler = (*Map[string, int])(nil)

type Map[K, V comparable] struct {
	fwd rbtree.Tree[K, V]
	inv rbtree.Tree[V, K]
}

func New[K, V cmp.Ordered]() *Map[K, V] {
	return &Map[K, V]{
		fwd: *rbtree.New[K, V](),
		inv: *rbtree.New[V, K](),
	}
}

func NewWith[K, V comparable](kcmp cmp.Comparator[K], vcmp cmp.Comparator[V]) *Map[K, V] {
	return &Map[K, V]{
		fwd: *rbtree.NewWith[K, V](kcmp),
		inv: *rbtree.NewWith[V, K](vcmp),
	}
}

func (m *Map[K, V]) Put(k K, v V) {
	if val, ok := m.fwd.Get(k); ok {
		m.inv.Delete(val)
	}

	if key, ok := m.inv.Get(v); ok {
		m.fwd.Delete(key)
	}

	m.fwd.Put(k, v)
	m.inv.Put(v, k)
}

func (m *Map[K, V]) Has(k K) bool {
	return m.fwd.Has(k)
}

func (m *Map[K, V]) HasValue(v V) bool {
	return m.inv.Has(v)
}

func (m *Map[K, V]) Get(k K) (v V, ok bool) {
	return m.fwd.Get(k)
}

func (m *Map[K, V]) GetKey(v V) (k K, ok bool) {
	return m.inv.Get(v)
}

func (m *Map[K, V]) Delete(k K) (v V, ok bool) {
	if v, ok := m.fwd.Get(k); ok {
		m.fwd.Delete(k)
		m.inv.Delete(v)

		return v, true
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

func (m *Map[K, V]) Begin() (k K, v V, ok bool) {
	return m.fwd.Begin()
}

func (m *Map[K, V]) End() (k K, v V, ok bool) {
	return m.fwd.End()
}

func (m *Map[K, V]) DeleteBegin() (k K, v V, ok bool) {
	k, v, ok = m.fwd.DeleteBegin()
	if ok {
		m.inv.Delete(v)
	}

	return
}

func (m *Map[K, V]) DeleteEnd() (k K, v V, ok bool) {
	k, v, ok = m.fwd.DeleteEnd()
	if ok {
		m.inv.Delete(v)
	}

	return
}

func (m *Map[K, V]) Iter() iter.Seq2[K, V] {
	return m.fwd.Iter()
}

func (m *Map[K, V]) RIter() iter.Seq2[K, V] {
	return m.fwd.RIter()
}

func (m *Map[K, V]) IsEmpty() bool {
	return m.Len() == 0
}

func (m *Map[K, V]) Len() int {
	return m.fwd.Len()
}

func (m *Map[K, V]) Keys() []K {
	return m.fwd.Keys()
}

func (m *Map[K, V]) Values() []V {
	return m.inv.Keys()
}

func (m *Map[K, V]) ToSlice() []V {
	return m.fwd.ToSlice()
}

func (m *Map[K, V]) Entries() ([]K, []V) {
	return m.fwd.Entries()
}

func (m *Map[K, V]) Clear() {
	m.fwd.Clear()
	m.inv.Clear()
}

func (m *Map[K, V]) Clone() container.Map[K, V] {
	return &Map[K, V]{
		fwd: *(m.fwd.Clone().(*rbtree.Tree[K, V])),
		inv: *(m.inv.Clone().(*rbtree.Tree[V, K])),
	}
}

func (m *Map[K, V]) MarshalJSON() ([]byte, error) {
	return m.fwd.MarshalJSON()
}

func (m *Map[K, V]) UnmarshalJSON(data []byte) error {
	var elems map[K]V

	if err := json.Unmarshal(data, &elems); err != nil {
		return err
	}

	m.Clear()

	for k, v := range elems {
		m.Put(k, v)
	}

	return nil
}

func (m *Map[K, V]) String() string {
	s := "TreeBidiMap\nmap["

	for k, v := range m.fwd.Iter() {
		s += fmt.Sprintf("%v:%v ", k, v)
	}

	return strings.TrimRight(s, " ") + "]"
}
