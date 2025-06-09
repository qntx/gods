package linkedhashset

import (
	"container/list"

	"github.com/qntx/gods/container"
)

// Assert Iterator implementation
var _ container.ReverseIteratorWithIndex[int] = (*Iterator[int])(nil)

// Iterator holding the iterator's state
type Iterator[T comparable] struct {
	set     *Set[T]       // Reference to the set
	element *list.Element // Current element
}

// Iterator returns a stateful iterator whose values can be fetched by an index.
func (set *Set[T]) Iterator() Iterator[T] {
	return Iterator[T]{set: set, element: nil}
}

// Next moves the iterator to the next element and returns true if there was a next element in the container.
// If Next() returns true, then next element's index and value can be retrieved by Index() and Value().
// If Next() was called for the first time, then it will point the iterator to the first element if it exists.
// Modifies the state of the iterator.
func (iterator *Iterator[T]) Next() bool {
	if iterator.element == nil { // Initial call or after Begin()
		iterator.element = iterator.set.ordering.Front()
	} else {
		iterator.element = iterator.element.Next()
	}
	return iterator.element != nil
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (iterator *Iterator[T]) Prev() bool {
	if iterator.element == nil { // Initial call or after End()
		iterator.element = iterator.set.ordering.Back()
	} else {
		iterator.element = iterator.element.Prev()
	}
	return iterator.element != nil
}

// Value returns the current element's value.
// Does not modify the state of the iterator.
func (iterator *Iterator[T]) Value() T {
	if iterator.element == nil {
		panic("Iterator: Value() called on an invalid element")
	}
	return iterator.element.Value.(T)
}

// Index returns the current element's index.
// Returns -1 if the iterator is not pointing to a valid element.
// This is an O(N) operation for linked list.
func (iterator *Iterator[T]) Index() int {
	if iterator.element == nil {
		return -1
	}

	current := iterator.set.ordering.Front()
	index := 0
	for current != nil {
		if current == iterator.element {
			return index
		}
		index++
		current = current.Next()
	}

	// Should not happen if currentElement is part of the list and not nil
	// but as a fallback, or if currentElement was somehow invalidated externally.
	// However, the initial nil check should cover invalid states for Index().
	return -1 // Or panic, indicating an inconsistent iterator state if currentElement is non-nil but not found.
}

// Begin resets the iterator to its initial state (one-before-first)
// Call Next() to fetch the first element if any.
func (iterator *Iterator[T]) Begin() {
	iterator.element = nil
}

// End moves the iterator past the last element (one-past-the-end).
// Call Prev() to fetch the last element if any.
func (iterator *Iterator[T]) End() {
	iterator.element = nil
}

// First moves the iterator to the first element and returns true if there was a first element in the container.
// If First() returns true, then first element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (iterator *Iterator[T]) First() bool {
	iterator.element = iterator.set.ordering.Front()
	return iterator.element != nil
}

// Last moves the iterator to the last element and returns true if there was a last element in the container.
// If Last() returns true, then last element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (iterator *Iterator[T]) Last() bool {
	iterator.element = iterator.set.ordering.Back()
	return iterator.element != nil
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
