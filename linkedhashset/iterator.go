package linkedhashset

import (
	"container/list"

	"github.com/qntx/gods/container"
)

// Assert Iterator implementation
var _ container.ReverseIteratorWithIndex[int] = (*Iterator[int])(nil)

// Iterator holds the iterator's state
type Iterator[T comparable] struct {
	set     *Set[T]
	element *list.Element // current element
	index   int           // current index, -1 before first, set.Size() after last
}

// Iterator returns a new iterator for the set.
func (set *Set[T]) Iterator() *Iterator[T] {
	return &Iterator[T]{set: set, element: nil, index: -1} // Initial state: before first
}

// Next moves the iterator to the next element and returns true if there was a next element in the container.
// If Next() returns true, then next element's index and value can be retrieved by Index() and Value().
// If Next() was called for the first time, then it will point the iterator to the first element if it exists.
// Modifies the state of the iterator.
func (iterator *Iterator[T]) Next() bool {
	if iterator.index >= iterator.set.Size()-1 { // Already at or past the last element
		iterator.element = nil
		iterator.index = iterator.set.Size() // Ensure index is set.Size() if past the end
		return false
	}

	if iterator.index == -1 { // First call to Next() or after Begin()
		iterator.element = iterator.set.ordering.Front()
	} else {
		iterator.element = iterator.element.Next()
	}

	if iterator.element != nil {
		iterator.index++
		return true
	}

	// Should not happen if logic above is correct and index < set.Size()-1
	// but as a safeguard, if element becomes nil unexpectedly:
	iterator.index = iterator.set.Size()
	return false
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (iterator *Iterator[T]) Prev() bool {
	if iterator.index <= 0 { // Already at or before the first element
		iterator.element = nil
		iterator.index = -1 // Ensure index is -1 if before the start
		return false
	}

	if iterator.index == iterator.set.Size() { // First call to Prev() after End()
		iterator.element = iterator.set.ordering.Back()
	} else {
		iterator.element = iterator.element.Prev()
	}

	if iterator.element != nil {
		iterator.index--
		return true
	}

	// Should not happen if logic above is correct and index > 0
	// but as a safeguard, if element becomes nil unexpectedly:
	iterator.index = -1
	return false
}

// Value returns the current element's value.
// Does not modify the state of the iterator.
// Panics if the iterator is not pointing to a valid element.
func (iterator *Iterator[T]) Value() T {
	if iterator.element == nil {
		panic("Iterator: Value() called on an invalid element")
	}
	return iterator.element.Value.(T)
}

// Index returns the current element's index.
// Does not modify the state of the iterator.
func (iterator *Iterator[T]) Index() int {
	return iterator.index
}

// Begin resets the iterator to its initial state (one-before-first)
// Call Next() to advance to the first element if any.
func (iterator *Iterator[T]) Begin() {
	iterator.element = nil
	iterator.index = -1
}

// End moves the iterator past the last element (one-past-last).
// Call Prev() to advance to the last element if any.
func (iterator *Iterator[T]) End() {
	iterator.element = nil
	iterator.index = iterator.set.Size()
}

// First moves the iterator to the first element and returns true if there was a first element in the container.
// If First() returns true, then first element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (iterator *Iterator[T]) First() bool {
	first := iterator.set.ordering.Front()
	if first != nil {
		iterator.element = first
		iterator.index = 0
		return true
	}
	iterator.element = nil
	iterator.index = -1 // No first element, so index is -1 (like Begin)
	return false
}

// Last moves the iterator to the last element and returns true if there was a last element in the container.
// If Last() returns true, then last element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (iterator *Iterator[T]) Last() bool {
	last := iterator.set.ordering.Back()
	if last != nil {
		iterator.element = last
		iterator.index = iterator.set.Size() - 1
		return true
	}
	iterator.element = nil
	iterator.index = -1 // No last element, so index is -1 (like Begin, or Size if End for empty)
	// To be consistent with First on empty, -1 is better.
	// However, if set is empty, Size() is 0. End() sets index to 0.
	// If set is empty, Size()-1 is -1. This seems consistent.
	return false
}

// NextTo moves the iterator to the next element from current position that satisfies the condition given by the
// passed function, and returns true if there was a next element in the container.
// If NextTo() returns true, then next element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (iterator *Iterator[T]) NextTo(f func(index int, value T) bool) bool {
	for iterator.Next() {
		index, value := iterator.Index(), iterator.Value()
		if f(index, value) {
			return true
		}
	}
	return false
}

// PrevTo moves the iterator to the previous element from current position that satisfies the condition given by the
// passed function, and returns true if there was a next element in the container.
// If PrevTo() returns true, then next element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (iterator *Iterator[T]) PrevTo(f func(index int, value T) bool) bool {
	for iterator.Prev() {
		index, value := iterator.Index(), iterator.Value()
		if f(index, value) {
			return true
		}
	}
	return false
}
