package btreeset

import (
	"github.com/qntx/gods/btree"
	"github.com/qntx/gods/container"
)

var _ container.EnumerableWithIndex[int] = (*Set[int])(nil)

// Each calls the given function once for each element, passing that element's index and value.
func (s *Set[T]) Each(f func(index int, value T)) {
	index := 0
	for value := range s.tree.Iter() {
		f(index, value)

		index++
	}
}

// Map invokes the given function once for each element and returns a
// container containing the values returned by the given function.
func (s *Set[T]) Map(f func(index int, value T) T) *Set[T] {
	newSet := &Set[T]{tree: btree.NewWith[T, struct{}](s.tree.MaxChildren(), s.tree.Comparator)}
	index := 0

	for value := range s.tree.Iter() {
		newSet.Add(f(index, value))

		index++
	}

	return newSet
}

// Select returns a new container containing all elements for which the given function returns a true value.
func (s *Set[T]) Select(f func(index int, value T) bool) *Set[T] {
	newSet := &Set[T]{tree: btree.NewWith[T, struct{}](s.tree.MaxChildren(), s.tree.Comparator)}
	index := 0

	for value := range s.tree.Iter() {
		if f(index, value) {
			newSet.Add(value)
		}

		index++
	}

	return newSet
}

// Any passes each element of the container to the given function and
// returns true if the function ever returns true for any element.
func (s *Set[T]) Any(f func(index int, value T) bool) bool {
	index := 0

	for value := range s.tree.Iter() {
		if f(index, value) {
			return true
		}

		index++
	}

	return false
}

// All passes each element of the container to the given function and
// returns true if the function returns true for all elements.
func (s *Set[T]) All(f func(index int, value T) bool) bool {
	index := 0

	for value := range s.tree.Iter() {
		if !f(index, value) {
			return false
		}

		index++
	}

	return true
}

// Find passes each element of the container to the given function and returns
// the first (index,value) for which the function is true or -1,nil otherwise
// if no element matches the criteria.
func (s *Set[T]) Find(f func(index int, value T) bool) (int, T) {
	index := 0

	for value := range s.tree.Iter() {
		if f(index, value) {
			return index, value
		}

		index++
	}

	var t T

	return -1, t
}
