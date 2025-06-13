// Package rbtreeset provides a set implementation using a red-black tree.
package rbtreeset

import (
	"encoding/json"
	"fmt"
	"iter"
	"strings"

	"slices"

	"github.com/qntx/gods/cmp"
	"github.com/qntx/gods/container"
	"github.com/qntx/gods/rbtree"
)

var _ container.Set[string] = (*Set[string])(nil)
var _ json.Marshaler = (*Set[string])(nil)
var _ json.Unmarshaler = (*Set[string])(nil)

type Set[T comparable] struct {
	tree *rbtree.Tree[T, struct{}]
}

func New[T cmp.Ordered](values ...T) *Set[T] {
	return NewWith(cmp.Compare[T], values...)
}

func NewWith[T comparable](cmp cmp.Comparator[T], values ...T) *Set[T] {
	s := &Set[T]{tree: rbtree.NewWith[T, struct{}](cmp)}
	for _, v := range values {
		s.tree.Put(v, struct{}{})
	}

	return s
}

func (s *Set[T]) Add(values T) bool {
	prevLen := s.tree.Len()
	s.tree.Put(values, struct{}{})

	return prevLen != s.tree.Len()
}

func (s *Set[T]) Append(values ...T) int {
	prevLen := s.tree.Len()

	for _, v := range values {
		s.tree.Put(v, struct{}{})
	}

	return s.tree.Len() - prevLen
}

func (s *Set[T]) Remove(values T) {
	s.tree.Delete(values)
}

func (s *Set[T]) RemoveAll(values ...T) {
	for _, v := range values {
		s.tree.Delete(v)
	}
}

func (s *Set[T]) Pop() (k T, ok bool) {
	k, _, ok = s.tree.DeleteBegin()

	return k, ok
}

func (s *Set[T]) PopEnd() (k T, ok bool) {
	k, _, ok = s.tree.DeleteEnd()

	return k, ok
}

func (s *Set[T]) ContainsOne(v T) bool {
	return s.tree.Has(v)
}

func (s *Set[T]) Contains(values ...T) bool {
	for _, v := range values {
		if !s.tree.Has(v) {
			return false
		}
	}

	return true
}

func (s *Set[T]) ContainsAny(val ...T) bool {
	return slices.ContainsFunc(val, s.tree.Has)
}

func (s *Set[T]) ContainsAnyElement(other container.Set[T]) bool {
	o := other.(*Set[T])

	return slices.ContainsFunc(o.tree.Keys(), s.tree.Has)
}

func (s *Set[T]) Equal(other container.Set[T]) bool {
	o := other.(*Set[T])

	if s.tree.Len() != o.tree.Len() {
		return false
	}

	for _, v := range s.tree.Keys() {
		if !o.tree.Has(v) {
			return false
		}
	}

	return true
}

func (s *Set[T]) IsSubset(other container.Set[T]) bool {
	o := other.(*Set[T])

	for _, v := range s.tree.Keys() {
		if !o.tree.Has(v) {
			return false
		}
	}

	return true
}

func (s *Set[T]) IsProperSubset(other container.Set[T]) bool {
	return s.IsSubset(other) && s.tree.Len() < other.(*Set[T]).tree.Len()
}

func (s *Set[T]) IsSuperset(other container.Set[T]) bool {
	o := other.(*Set[T])

	for _, v := range o.tree.Keys() {
		if !s.tree.Has(v) {
			return false
		}
	}

	return true
}

func (s *Set[T]) IsProperSuperset(other container.Set[T]) bool {
	return s.IsSuperset(other) && s.tree.Len() > other.(*Set[T]).tree.Len()
}

// Union returns a new set containing all elements from s or other.
// Returns an empty set if comparators differ.
// Ref: https://en.wikipedia.org/wiki/Union_(set_theory)
func (s *Set[T]) Union(other container.Set[T]) container.Set[T] {
	res := NewWith(s.tree.Comparator())

	for v := range s.Iter() {
		res.Add(v)
	}

	for v := range other.Iter() {
		res.Add(v)
	}

	return res
}

// Intersect returns a new set containing elements present in both s and other.
// Returns an empty set if comparators differ.
// Ref: https://en.wikipedia.org/wiki/Intersection_(set_theory)
func (s *Set[T]) Intersect(other container.Set[T]) container.Set[T] {
	res := NewWith(s.tree.Comparator())

	// Iterate over smaller set for efficiency.
	src, dst := s, other.(*Set[T])
	if s.Len() > other.Len() {
		src, dst = dst, src
	}

	for v := range src.Iter() {
		if dst.Contains(v) {
			res.Add(v)
		}
	}

	return res
}

// Difference returns a new set containing elements in s but not in other.
// Returns an empty set if comparators differ.
// Ref: https://proofwiki.org/wiki/Definition:Set_Difference
func (s *Set[T]) Difference(other container.Set[T]) container.Set[T] {
	res := NewWith(s.tree.Comparator())

	for v := range s.Iter() {
		if !other.Contains(v) {
			res.Add(v)
		}
	}

	return res
}

// SymmetricDifference returns a new set with all elements which are
// in either this set or the other set but not in both.
// Returns an empty set if comparators differ.
// Ref: https://en.wikipedia.org/wiki/Symmetric_difference
func (s *Set[T]) SymmetricDifference(other container.Set[T]) container.Set[T] {
	res := NewWith(s.tree.Comparator())

	for v := range s.Iter() {
		if !other.Contains(v) {
			res.Add(v)
		}
	}

	for v := range other.Iter() {
		if !s.Contains(v) {
			res.Add(v)
		}
	}

	return res
}

func (s *Set[T]) IsEmpty() bool {
	return s.tree.Len() == 0
}

func (s *Set[T]) Len() int {
	return s.tree.Len()
}

func (s *Set[T]) Clear() {
	s.tree.Clear()
}

func (s *Set[T]) Values() []T {
	return s.tree.Keys()
}

func (s *Set[T]) ToSlice() []T {
	return s.tree.Keys()
}

func (s *Set[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for k := range s.tree.Iter() {
			if !yield(k) {
				return
			}
		}
	}
}

func (s *Set[T]) Clone() container.Set[T] {
	return &Set[T]{tree: s.tree.Clone().(*rbtree.Tree[T, struct{}])}
}

func (s *Set[T]) Each(f func(T) bool) {
	for v := range s.Iter() {
		if !f(v) {
			return
		}
	}
}

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
		set.Append(elements...)
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
