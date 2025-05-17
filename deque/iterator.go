// Package deque provides an iterator for traversing the circular buffer queue.
//
// This file implements a stateful iterator for the Queue type, supporting both
// forward and reverse iteration with O(1) access to elements and indices.
package deque

import (
	"errors"

	"github.com/qntx/gods/container"
)

// --------------------------------------------------------------------------------
// Constants and Errors
// --------------------------------------------------------------------------------

// Predefined errors for iterator operations.
var (
	ErrInvalidIteratorPosition = errors.New("iterator value accessed at invalid position")
)

// --------------------------------------------------------------------------------
// Interface Assertions
// --------------------------------------------------------------------------------

// Ensure Iterator implements container.ReverseIteratorWithIndex at compile time.
var _ container.ReverseIteratorWithIndex[int] = (*Iterator[int])(nil)

// --------------------------------------------------------------------------------
// Types
// --------------------------------------------------------------------------------

// Iterator provides forward and reverse traversal over a Queue's elements.
//
// It maintains a position relative to the queue's front (index 0) to back (Len()-1).
// The iterator is read-only and does not modify the underlying queue. Most operations
// are O(1), except for NextTo and PrevTo which may be O(n).
type Iterator[T comparable] struct {
	queue *Deque[T] // Reference to the queue being iterated.
	index int       // Current position: -1 (before start), Len() (past end), or valid index.
}

// --------------------------------------------------------------------------------
// Constructor
// --------------------------------------------------------------------------------

// Iterator creates a new iterator for the queue.
//
// Starts in an invalid state (before the first element). Use Next() to reach the
// first element or End() followed by Prev() for the last. Time complexity: O(1).
func (q *Deque[T]) Iterator() *Iterator[T] {
	return &Iterator[T]{
		queue: q,
		index: -1, // Before-first state.
	}
}

// --------------------------------------------------------------------------------
// Public Methods
// --------------------------------------------------------------------------------

// Next advances the iterator to the next element.
//
// Returns true if the iterator is at a valid element after moving, false if it
// reaches or is already past the end. Time complexity: O(1).
func (it *Iterator[T]) Next() bool {
	if it.queue.Empty() || it.index >= it.queue.Len()-1 {
		it.index = it.queue.Len() // Move to past-end state.

		return false
	}

	it.index++

	return true
}

// Prev moves the iterator to the previous element.
//
// Returns true if the iterator is at a valid element after moving, false if it
// reaches or is already before the start. Time complexity: O(1).
func (it *Iterator[T]) Prev() bool {
	if it.queue.Empty() || it.index <= 0 {
		it.index = -1 // Move to before-start state.

		return false
	}

	it.index--

	return true
}

// Value returns the current element's value.
//
// Panics if the iterator is not at a valid position (before start or past end).
// Time complexity: O(1).
func (it *Iterator[T]) Value() T {
	if !it.valid() {
		panic("ringbuf: " + ErrInvalidIteratorPosition.Error())
	}

	return it.queue.buf[it.wrap(it.queue.start+it.index)]
}

// Index returns the current iterator position.
//
// Returns -1 if before the start, Len() if past the end, or a valid index (0 to Len()-1).
// Time complexity: O(1).
func (it *Iterator[T]) Index() int {
	return it.index
}

// Begin resets the iterator to before the first element.
//
// Use Next() to move to the first element. Time complexity: O(1).
func (it *Iterator[T]) Begin() {
	it.index = -1
}

// End positions the iterator past the last element.
//
// Use Prev() to move to the last element. Time complexity: O(1).
func (it *Iterator[T]) End() {
	it.index = it.queue.Len()
}

// First moves the iterator to the first element.
//
// Returns true if the queue is non-empty, false otherwise. Time complexity: O(1).
func (it *Iterator[T]) First() bool {
	if it.queue.Empty() {
		it.index = -1

		return false
	}

	it.index = 0

	return true
}

// Last moves the iterator to the last element.
//
// Returns true if the queue is non-empty, false otherwise. Time complexity: O(1).
func (it *Iterator[T]) Last() bool {
	if it.queue.Empty() {
		it.index = -1

		return false
	}

	it.index = it.queue.Len() - 1

	return true
}

// NextTo advances to the next element satisfying the given condition.
//
// Moves forward from the current position until an element matches the predicate
// or the end is reached. Returns true if a match is found. Time complexity: O(n)
// in the worst case.
func (it *Iterator[T]) NextTo(f func(index int, value T) bool) bool {
	for it.Next() {
		if f(it.index, it.Value()) {
			return true
		}
	}

	return false
}

// PrevTo moves to the previous element satisfying the given condition.
//
// Moves backward from the current position until an element matches the predicate
// or the start is reached. Returns true if a match is found. Time complexity: O(n)
// in the worst case.
func (it *Iterator[T]) PrevTo(f func(index int, value T) bool) bool {
	for it.Prev() {
		if f(it.index, it.Value()) {
			return true
		}
	}

	return false
}

// --------------------------------------------------------------------------------
// Private Helpers
// --------------------------------------------------------------------------------

// valid checks if the iterator is at a valid element position.
func (it *Iterator[T]) valid() bool {
	return it.index >= 0 && it.index < it.queue.Len()
}

// wrap ensures the index stays within the queue's buffer bounds.
func (it *Iterator[T]) wrap(idx int) int {
	return idx % it.queue.capacity
}
