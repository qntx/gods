package linkedhashset

import "github.com/qntx/gods/container"

// Assert Enumerable implementation
// Note: container.Enumerable[T] expects Map and Select to return container.Container[T].
// Since *Set[T] implements container.Container[T], returning *Set[T] is valid.
var _ container.EnumerableWithIndex[int] = (*Set[int])(nil)

// Each calls the given function once for each element, passing that element's index and value.
func (set *Set[T]) Each(f func(index int, value T)) {
	i := 0
	for v := range set.Iter() {
		f(i, v)

		i++
	}
}

// Map invokes the given function once for each element and returns a new Set
// containing the values returned by the given function.
// The new Set will contain the mapped values, preserving the original order if mapped values are unique.
// If mapped values are not unique, only the first occurrence (based on iteration order) will be kept.
func (set *Set[T]) Map(f func(index int, value T) T) *Set[T] {
	newSet := New[T]()
	i := 0

	for v := range set.Iter() {
		newSet.Add(f(i, v))

		i++
	}

	return newSet
}

// Select returns a new Set containing all elements for which the given function returns a true value.
// The new Set will preserve the original insertion order of the selected elements.
func (set *Set[T]) Select(f func(index int, value T) bool) *Set[T] {
	newSet := New[T]()

	i := 0
	for v := range set.Iter() {
		if f(i, v) {
			newSet.Add(v)
		}

		i++
	}

	return newSet
}

// Any returns true if the given function returns true for one or more elements.
// This function is short-circuiting, i.e., it returns as soon as the first true is found.
func (set *Set[T]) Any(f func(index int, value T) bool) bool {
	i := 0
	for v := range set.Iter() {
		if f(i, v) {
			return true
		}

		i++
	}

	return false
}

// All returns true if the given function returns true for all elements.
// This function is short-circuiting, i.e., it returns false as soon as the first false is found.
func (set *Set[T]) All(f func(index int, value T) bool) bool {
	i := 0
	for v := range set.Iter() {
		if !f(i, v) {
			return false
		}

		i++
	}

	return true
}

// Find returns the first index and value for which the provided function returns true.
// If no element satisfies the condition, it returns -1 and the zero value of T.
func (set *Set[T]) Find(f func(index int, value T) bool) (int, T) {
	i := 0
	for v := range set.Iter() {
		if f(i, v) {
			return i, v
		}

		i++
	}

	var zero T // zero value for type T

	return -1, zero
}
