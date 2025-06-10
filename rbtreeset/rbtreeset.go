// Package rbtreeset provides a set implementation using a red-black tree.
package rbtreeset

import (
	"encoding/json"
	"fmt"
	"iter"
	"reflect"
	"strings"

	"github.com/qntx/gods/cmp"
	"github.com/qntx/gods/rbtree"
)

// present is a marker for set membership.
var present = struct{}{}

// Set is a red-black tree-based set of comparable elements.
type Set[T comparable] struct {
	tree *rbtree.Tree[T, struct{}]
}

// New creates a new set for ordered types with optional initial values.
func New[T cmp.Ordered](values ...T) *Set[T] {
	return NewWith(cmp.GenericComparator[T], values...)
}

// NewWith creates a new set with a custom comparator and optional initial values.
func NewWith[T comparable](cmp cmp.Comparator[T], values ...T) *Set[T] {
	s := &Set[T]{tree: rbtree.NewWith[T, struct{}](cmp)}
	for _, v := range values {
		s.tree.Put(v, present)
	}

	return s
}

// Add inserts one or more elements into the set.
func (s *Set[T]) Add(values ...T) {
	for _, v := range values {
		s.tree.Put(v, present)
	}
}

// Remove deletes one or more elements from the set.
func (s *Set[T]) Remove(values ...T) {
	for _, v := range values {
		s.tree.Delete(v)
	}
}

// Contains checks if all specified elements are present in the set.
// Returns true if no elements are provided, as a set is a superset of an empty set.
func (s *Set[T]) Contains(values ...T) bool {
	for _, v := range values {
		if _, ok := s.tree.Get(v); !ok {
			return false
		}
	}

	return true
}

// Empty reports whether the set contains no elements.
func (s *Set[T]) Empty() bool {
	return s.tree.Len() == 0
}

// Len returns the number of elements in the set.
func (s *Set[T]) Len() int {
	return s.tree.Len()
}

// Clear removes all elements from the set.
func (s *Set[T]) Clear() {
	s.tree.Clear()
}

// Values returns a slice of all elements in the set.
func (s *Set[T]) Values() []T {
	return s.tree.Keys()
}

// All returns an iterator over all elements in the set in sorted order.
func (s *Set[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for k := range s.tree.Iter() {
			if !yield(k) {
				return
			}
		}
	}
}

var _ json.Marshaler = (*Set[string])(nil)
var _ json.Unmarshaler = (*Set[string])(nil)

// MarshalJSON outputs the JSON representation of the set.
func (set *Set[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(set.Values())
}

// UnmarshalJSON populates the set from the input JSON representation.
func (set *Set[T]) UnmarshalJSON(data []byte) error {
	var elements []T

	err := json.Unmarshal(data, &elements)
	if err == nil {
		set.Clear()
		set.Add(elements...)
	}

	return err
}

// String returns a string representation of the set.
func (s *Set[T]) String() string {
	var b strings.Builder

	b.WriteString("TreeSet\n")

	for v := range s.Iter() {
		fmt.Fprintf(&b, "%v", v)
	}

	return b.String()
}

// Intersection returns a new set containing elements present in both s and other.
// Returns an empty set if comparators differ.
// Ref: https://en.wikipedia.org/wiki/Intersection_(set_theory)
func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	res := NewWith(s.tree.Comparator())

	sCmp := reflect.ValueOf(s.tree.Comparator())
	oCmp := reflect.ValueOf(other.tree.Comparator())

	if sCmp.Pointer() != oCmp.Pointer() {
		return res
	}

	// Iterate over smaller set for efficiency.
	src, dst := s, other
	if s.Len() > other.Len() {
		src, dst = other, s
	}

	for v := range src.Iter() {
		if dst.Contains(v) {
			res.Add(v)
		}
	}

	return res
}

// Union returns a new set containing all elements from s or other.
// Returns an empty set if comparators differ.
// Ref: https://en.wikipedia.org/wiki/Union_(set_theory)
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	res := NewWith(s.tree.Comparator())

	sCmp := reflect.ValueOf(s.tree.Comparator())
	oCmp := reflect.ValueOf(other.tree.Comparator())

	if sCmp.Pointer() != oCmp.Pointer() {
		return res
	}

	for v := range s.Iter() {
		res.Add(v)
	}

	for v := range other.Iter() {
		res.Add(v)
	}

	return res
}

// Difference returns a new set containing elements in s but not in other.
// Returns an empty set if comparators differ.
// Ref: https://proofwiki.org/wiki/Definition:Set_Difference
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	res := NewWith(s.tree.Comparator())

	sCmp := reflect.ValueOf(s.tree.Comparator())
	oCmp := reflect.ValueOf(other.tree.Comparator())

	if sCmp.Pointer() != oCmp.Pointer() {
		return res
	}

	for v := range s.Iter() {
		if !other.Contains(v) {
			res.Add(v)
		}
	}

	return res
}
