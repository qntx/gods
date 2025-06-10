// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pqueue

import (
	"cmp"

	"github.com/qntx/gods/container"
)

// --------------------------------------------------------------------------------
// Iterator Interface Assertions
// --------------------------------------------------------------------------------

// Verify Iterator implements required interfaces at compile time.
var _ container.ReverseIteratorWithIndex[string] = (*Iterator[string, int])(nil)

// --------------------------------------------------------------------------------
// Iterator Type
// --------------------------------------------------------------------------------

// Iterator provides ordered access to the elements in a priority queue.
// Note: Iterator does not guarantee the heap property during traversal.
// Elements are accessed in the underlying slice order, not by priority.
type Iterator[T comparable, V cmp.Ordered] struct {
	queue  *PriorityQueue[T, V] // Reference to the original queue.
	index  int                  // Current position.
	values []*Item[T, V]        // Copy of queue items to allow safe iteration.
}

// --------------------------------------------------------------------------------
// Constructor
// --------------------------------------------------------------------------------

// Iterator returns a stateful iterator for the priority queue.
// Complexity: O(n) as it makes a copy of the queue's items for safe iteration.
func (pq *PriorityQueue[T, V]) Iterator() *Iterator[T, V] {
	// Make a copy of the queue items to allow safe iteration
	values := make([]*Item[T, V], len(pq.heap))
	copy(values, pq.heap)

	return &Iterator[T, V]{
		queue:  pq,
		index:  -1, // Start before first element
		values: values,
	}
}

// --------------------------------------------------------------------------------
// Forward Iteration Methods
// --------------------------------------------------------------------------------

// Complexity: O(1).
func (iterator *Iterator[T, V]) Next() bool {
	if iterator.index < len(iterator.values)-1 {
		iterator.index++

		return true
	}

	return false
}

// Complexity: O(1).
func (iterator *Iterator[T, V]) Value() T {
	if iterator.index < 0 || iterator.index >= len(iterator.values) {
		panic("Invalid iterator state: call Next() or Prev() before Value()")
	}

	return iterator.values[iterator.index].Value
}

// Complexity: O(1).
func (iterator *Iterator[T, V]) Index() int {
	return iterator.index
}

// Complexity: O(1).
func (iterator *Iterator[T, V]) Begin() {
	iterator.index = -1
}

// Complexity: O(1).
func (iterator *Iterator[T, V]) First() bool {
	if len(iterator.values) > 0 {
		iterator.index = 0

		return true
	}

	return false
}

// Complexity: O(n).
func (iterator *Iterator[T, V]) NextTo(f func(index int, value T) bool) bool {
	for iterator.index < len(iterator.values)-1 {
		iterator.index++
		if f(iterator.index, iterator.values[iterator.index].Value) {
			return true
		}
	}

	return false
}

// --------------------------------------------------------------------------------
// Reverse Iteration Methods
// --------------------------------------------------------------------------------

// Complexity: O(1).
func (iterator *Iterator[T, V]) Prev() bool {
	if iterator.index > 0 {
		iterator.index--

		return true
	}

	return false
}

// Complexity: O(1).
func (iterator *Iterator[T, V]) End() {
	iterator.index = len(iterator.values)
}

// Complexity: O(1).
func (iterator *Iterator[T, V]) Last() bool {
	if len(iterator.values) > 0 {
		iterator.index = len(iterator.values) - 1

		return true
	}

	return false
}

// Complexity: O(n).
func (iterator *Iterator[T, V]) PrevTo(f func(index int, value T) bool) bool {
	for iterator.index > 0 {
		iterator.index--
		if f(iterator.index, iterator.values[iterator.index].Value) {
			return true
		}
	}

	return false
}

// --------------------------------------------------------------------------------
// Additional Methods (Optional)
// --------------------------------------------------------------------------------

// Complexity: O(1).
func (iterator *Iterator[T, V]) Priority() V {
	if iterator.index < 0 || iterator.index >= len(iterator.values) {
		panic("Invalid iterator state: call Next() or Prev() before Priority()")
	}

	return iterator.values[iterator.index].Priority
}
