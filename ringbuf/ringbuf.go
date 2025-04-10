// Package ringbuf implements a generic, non-thread-safe double-ended queue (deque)
// using a circular buffer.
//
// A circular buffer leverages a fixed-size array with logically connected ends,
// enabling efficient O(1) insertion and removal from both front and back.
// Ideal for scenarios requiring bounded queues with fast access at both ends.
//
// Reference: https://en.wikipedia.org/wiki/Circular_buffer
package ringbuf

import (
	"errors"
	"fmt"
	"strings"
)

// --------------------------------------------------------------------------------
// Constants and Errors

const (
	minCapacity = 1 // Minimum allowed capacity for the queue.
)

// Predefined errors for queue operations.
var (
	ErrInvalidCapacity = errors.New("capacity must be at least 1")
)

// --------------------------------------------------------------------------------
// Types and Interfaces

// Queue represents a double-ended queue implemented as a circular buffer.
//
// The type parameter T must be comparable to support equality checks if required.
// It maintains a fixed capacity with efficient O(1) operations for most methods.
type Queue[T comparable] struct {
	buf      []T // Underlying fixed-size buffer.
	start    int // Index of the front element.
	end      int // Index of the next available slot at the back.
	capacity int // Maximum number of elements.
	len      int // Current number of elements.
}

// Deque defines the interface for a double-ended queue.
//
// This allows for dependency injection and testing with alternative implementations.
type Deque[T comparable] interface {
	PushFront(T)
	PushBack(T)
	PopFront() (T, bool)
	PopBack() (T, bool)
	Front() (T, bool)
	Back() (T, bool)
	Len() int
	Capacity() int
	Empty() bool
	Full() bool
}

// --------------------------------------------------------------------------------
// Constructors

// New initializes a new Queue with the given capacity.
//
// Returns a pointer to a Queue instance. Panics if capacity is less than 1.
//
// Example:
//
//	q := ringbuf.New[int](5) // Creates a queue with capacity 5.
func New[T comparable](capacity int) *Queue[T] {
	if capacity < minCapacity {
		panic(ErrInvalidCapacity)
	}

	return &Queue[T]{
		buf:      make([]T, capacity),
		capacity: capacity,
	}
}

// --------------------------------------------------------------------------------
// Public Methods

// PushFront inserts an element at the front of the queue.
//
// Overwrites the oldest element (back) if the queue is full.
// Time complexity: O(1).
func (q *Queue[T]) PushFront(val T) {
	q.start = q.prev(q.start)
	q.buf[q.start] = val

	if !q.Full() {
		q.len++
	} else {
		q.end = q.prev(q.end)
	}
}

// PushBack inserts an element at the back of the queue.
//
// Overwrites the oldest element (front) if the queue is full.
// Time complexity: O(1).
func (q *Queue[T]) PushBack(val T) {
	q.buf[q.end] = val
	q.end = q.next(q.end)

	if !q.Full() {
		q.len++
	} else {
		q.start = q.next(q.start)
	}
}

// PopFront removes and returns the front element.
//
// Returns the zero value of T and false if the queue is empty.
// Time complexity: O(1).
func (q *Queue[T]) PopFront() (val T, ok bool) {
	if q.Empty() {
		return val, false
	}

	val = q.buf[q.start]
	q.start = q.next(q.start)
	q.len--

	return val, true
}

// PopBack removes and returns the back element.
//
// Returns the zero value of T and false if the queue is empty.
// Time complexity: O(1).
func (q *Queue[T]) PopBack() (val T, ok bool) {
	if q.Empty() {
		return val, false
	}

	q.end = q.prev(q.end)
	val = q.buf[q.end]
	q.len--

	return val, true
}

// Front retrieves the front element without removing it.
//
// Returns the zero value of T and false if the queue is empty.
// Time complexity: O(1).
func (q *Queue[T]) Front() (val T, ok bool) {
	if q.Empty() {
		return val, false
	}

	return q.buf[q.start], true
}

// Back retrieves the back element without removing it.
//
// Returns the zero value of T and false if the queue is empty.
// Time complexity: O(1).
func (q *Queue[T]) Back() (val T, ok bool) {
	if q.Empty() {
		return val, false
	}

	return q.buf[q.prev(q.end)], true
}

// Peek retrieves the element at the specified index.
//
// Index 0 is the front, Len()-1 is the back. Returns the zero value of T and false
// if the index is invalid. Time complexity: O(1).
func (q *Queue[T]) Peek(idx int) (val T, ok bool) {
	if idx < 0 || idx >= q.len {
		return val, false
	}

	return q.buf[q.wrap(q.start+idx)], true
}

// Empty checks if the queue has no elements.
//
// Time complexity: O(1).
func (q *Queue[T]) Empty() bool {
	return q.len == 0
}

// Full checks if the queue is at maximum capacity.
//
// Time complexity: O(1).
func (q *Queue[T]) Full() bool {
	return q.len == q.capacity
}

// Len returns the current number of elements.
//
// Time complexity: O(1).
func (q *Queue[T]) Len() int {
	return q.len
}

// Capacity returns the maximum number of elements the queue can hold.
//
// Time complexity: O(1).
func (q *Queue[T]) Capacity() int {
	return q.capacity
}

// Clear resets the queue to an empty state.
//
// Preserves capacity but reinitializes the buffer. Time complexity: O(n).
func (q *Queue[T]) Clear() {
	*q = *New[T](q.capacity)
}

// Values returns a slice of all elements in FIFO order.
//
// Returns nil if the queue is empty. Time complexity: O(n).
func (q *Queue[T]) Values() []T {
	if q.Empty() {
		return nil
	}

	vals := make([]T, q.len)
	for i := range q.len {
		vals[i] = q.buf[q.wrap(q.start+i)]
	}

	return vals
}

// String returns a string representation of the queue in FIFO order.
//
// Time complexity: O(n).
func (q *Queue[T]) String() string {
	var sb strings.Builder

	sb.WriteString("Queue[")

	for i := range q.len {
		if i > 0 {
			sb.WriteString(", ")
		}

		fmt.Fprintf(&sb, "%v", q.buf[q.wrap(q.start+i)])
	}

	sb.WriteString("]")

	return sb.String()
}

// --------------------------------------------------------------------------------
// Private Helpers

// next calculates the next index in the circular buffer.
func (q *Queue[T]) next(idx int) int {
	return (idx + 1) % q.capacity
}

// prev calculates the previous index in the circular buffer.
func (q *Queue[T]) prev(idx int) int {
	return (idx - 1 + q.capacity) % q.capacity
}

// wrap ensures the index stays within buffer bounds.
func (q *Queue[T]) wrap(idx int) int {
	return idx % q.capacity
}
