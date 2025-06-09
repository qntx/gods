package rbtreeset

import (
	"github.com/qntx/gods/container"
	"github.com/qntx/gods/rbtree"
)

var _ container.EnumerableWithIndex[int] = (*Set[int])(nil)

// Each calls the given function once for each element, passing that element's index and value.
func (s *Set[T]) Each(f func(index int, value T)) {
	iterator := s.Iterator()
	for iterator.Next() {
		f(iterator.Index(), iterator.Value())
	}
}

// Map invokes the given function once for each element and returns a
// container containing the values returned by the given function.
func (s *Set[T]) Map(f func(index int, value T) T) *Set[T] {
	newSet := &Set[T]{tree: rbtree.NewWith[T, struct{}](s.tree.Comparator())}
	iterator := s.Iterator()
	for iterator.Next() {
		newSet.Add(f(iterator.Index(), iterator.Value()))
	}
	return newSet
}

// Select returns a new container containing all elements for which the given function returns a true value.
func (s *Set[T]) Select(f func(index int, value T) bool) *Set[T] {
	newSet := &Set[T]{tree: rbtree.NewWith[T, struct{}](s.tree.Comparator())}
	iterator := s.Iterator()
	for iterator.Next() {
		if f(iterator.Index(), iterator.Value()) {
			newSet.Add(iterator.Value())
		}
	}
	return newSet
}

// Any passes each element of the container to the given function and
// returns true if the function ever returns true for any element.
func (s *Set[T]) Any(f func(index int, value T) bool) bool {
	iterator := s.Iterator()
	for iterator.Next() {
		if f(iterator.Index(), iterator.Value()) {
			return true
		}
	}
	return false
}

// All passes each element of the container to the given function and
// returns true if the function returns true for all elements.
func (s *Set[T]) All(f func(index int, value T) bool) bool {
	iterator := s.Iterator()
	for iterator.Next() {
		if !f(iterator.Index(), iterator.Value()) {
			return false
		}
	}
	return true
}

// Find passes each element of the container to the given function and returns
// the first (index,value) for which the function is true or -1,nil otherwise
// if no element matches the criteria.
func (s *Set[T]) Find(f func(index int, value T) bool) (int, T) {
	iterator := s.Iterator()
	for iterator.Next() {
		if f(iterator.Index(), iterator.Value()) {
			return iterator.Index(), iterator.Value()
		}
	}
	var t T
	return -1, t
}
