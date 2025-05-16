// Package deque implements a generic, non-thread-safe double-ended queue (deque)
// using a circular buffer.
//
// A circular buffer leverages an array with logically connected ends, enabling
// efficient O(1) insertion and removal from both front and back. It supports two
// modes: overwrite (fixed-size, overwrites oldest elements when full) and expansion
// (growable, doubles capacity when full), selected via a boolean flag.
// Ideal for scenarios requiring bounded or dynamic deques with fast access at both ends.
//
// Reference: https://en.wikipedia.org/wiki/Circular_buffer
package deque

import (
	"errors"
	"fmt"
	"strings"
)

// --------------------------------------------------------------------------------
// Constants and Errors

const (
	minCapacity = 1 // Minimum allowed capacity for the deque.
)

// Predefined errors for deque operations.
var (
	ErrInvalidCapacity = errors.New("capacity must be at least 1")
)

// --------------------------------------------------------------------------------
// Types and Interfaces

// Deque represents a double-ended queue implemented as a circular buffer.
//
// The type parameter T must be comparable to support equality checks if required.
// It supports two modes (overwrite or expansion, via growable) with efficient
// O(1) amortized operations for most methods.
type Deque[T comparable] struct {
	buf      []T  // Underlying buffer (fixed or growable).
	start    int  // Index of the front element.
	end      int  // Index of the next available slot at the back.
	capacity int  // Current capacity of the buffer.
	len      int  // Current number of elements.
	growable bool // True for expansion mode, false for overwrite mode.
}

// --------------------------------------------------------------------------------
// Constructors

// New initializes a new Deque with the given capacity in overwrite mode.
//
// Returns a pointer to a Deque instance.
//
// Example:
//
//	d := deque.New[int](5) // Creates a deque with capacity 5, overwrite mode.
func New[T comparable](capacity int) *Deque[T] {
	return NewWith[T](capacity, false)
}

// NewWith initializes a new Deque with the given capacity and mode.
//
// If growable is true, the deque grows when full; if false, it overwrites the
// oldest elements. Panics if capacity is less than 1.
//
// Example:
//
//	d := deque.NewWith[int](5, true) // Expansion mode.
func NewWith[T comparable](capacity int, growable bool) *Deque[T] {
	if capacity < minCapacity {
		panic(ErrInvalidCapacity)
	}

	return &Deque[T]{
		buf:      make([]T, capacity),
		capacity: capacity,
		growable: growable,
	}
}

// NewFrom creates a new Deque initialized with elements from the provided slice.
//
// The capacity parameter specifies the initial capacity. The actual capacity will be
// max(capacity, len(values), minCapacity) to ensure all elements fit.
// If growable is true, the deque will expand when full; otherwise, it will overwrite
// the oldest elements when full.
//
// Example:
//
//	d := deque.NewFrom([]int{1, 2, 3}, 10, true) // Creates a growable deque with capacity 10 and initial elements [1,2,3]
func NewFrom[T comparable](values []T, capacity int, growable bool) *Deque[T] {
	// Ensure capacity is at least the size of input slice and respects minimum capacity
	actualCapacity := max(max(capacity, len(values)), minCapacity)

	d := NewWith[T](actualCapacity, growable)

	// Add all elements from the slice
	for _, v := range values {
		d.PushBack(v)
	}

	return d
}

// --------------------------------------------------------------------------------
// Public Methods

// PushFront inserts an element at the front of the deque.
//
// In overwrite mode (growable=false), overwrites the oldest element (back) if full.
// In expansion mode (growable=true), doubles the capacity if full.
// Time complexity: O(1) amortized.
func (d *Deque[T]) PushFront(val T) {
	if d.Full() {
		if d.growable {
			d.grow()
		} else {
			// Overwrite mode: Move end backward to overwrite oldest
			d.end = d.prev(d.end)
		}
	}

	d.start = d.prev(d.start)
	d.buf[d.start] = val
	if !d.Full() || d.growable {
		d.len++
	}
}

// PushBack inserts an element at the back of the deque.
//
// In overwrite mode (growable=false), overwrites the oldest element (front) if full.
// In expansion mode (growable=true), doubles the capacity if full.
// Time complexity: O(1) amortized.
func (d *Deque[T]) PushBack(val T) {
	if d.Full() {
		if d.growable {
			d.grow()
		} else {
			// Overwrite mode: Move start forward to overwrite oldest
			d.start = d.next(d.start)
		}
	}

	d.buf[d.end] = val
	d.end = d.next(d.end)
	if !d.Full() || d.growable {
		d.len++
	}
}

// PopFront removes and returns the front element.
//
// Returns the zero value of T and false if the deque is empty.
// Time complexity: O(1).
func (d *Deque[T]) PopFront() (val T, ok bool) {
	if d.Empty() {
		return val, false
	}

	val = d.buf[d.start]
	d.start = d.next(d.start)
	d.len--

	return val, true
}

// PopBack removes and returns the back element.
//
// Returns the zero value of T and false if the deque is empty.
// Time complexity: O(1).
func (d *Deque[T]) PopBack() (val T, ok bool) {
	if d.Empty() {
		return val, false
	}

	d.end = d.prev(d.end)
	val = d.buf[d.end]
	d.len--

	return val, true
}

// Front retrieves the front element without removing it.
//
// Returns the zero value of T and false if the deque is empty.
// Time complexity: O(1).
func (d *Deque[T]) Front() (val T, ok bool) {
	if d.Empty() {
		return val, false
	}

	return d.buf[d.start], true
}

// Back retrieves the back element without removing it.
//
// Returns the zero value of T and false if the deque is empty.
// Time complexity: O(1).
func (d *Deque[T]) Back() (val T, ok bool) {
	if d.Empty() {
		return val, false
	}

	return d.buf[d.prev(d.end)], true
}

// Peek retrieves the element at the specified index.
//
// Index 0 is the front, Len()-1 is the back. Returns the zero value of T and false
// if the index is invalid. Time complexity: O(1).
func (d *Deque[T]) Peek(idx int) (val T, ok bool) {
	if idx < 0 || idx >= d.len {
		return val, false
	}

	return d.buf[d.wrap(d.start+idx)], true
}

// Empty checks if the deque has no elements.
//
// Time complexity: O(1).
func (d *Deque[T]) Empty() bool {
	return d.len == 0
}

// Full checks if the deque is at current capacity.
//
// Time complexity: O(1).
func (d *Deque[T]) Full() bool {
	return d.len == d.capacity
}

// Growable returns true if the deque is in expansion mode, false if in overwrite mode.
//
// Time complexity: O(1).
func (d *Deque[T]) Growable() bool {
	return d.growable
}

// Len returns the current number of elements.
//
// Time complexity: O(1).
func (d *Deque[T]) Len() int {
	return d.len
}

// Capacity returns the current capacity of the deque.
//
// Time complexity: O(1).
func (d *Deque[T]) Capacity() int {
	return d.capacity
}

// Clear resets the deque to an empty state.
//
// Preserves capacity and mode but reinitializes the buffer. Time complexity: O(n).
func (d *Deque[T]) Clear() {
	*d = *NewWith[T](d.capacity, d.growable)
}

// Values returns a slice of all elements in FIFO order.
//
// Returns nil if the deque is empty. Time complexity: O(n).
func (d *Deque[T]) Values() []T {
	if d.Empty() {
		return nil
	}

	vals := make([]T, d.len)
	for i := range d.len {
		vals[i] = d.buf[d.wrap(d.start+i)]
	}

	return vals
}

// String returns a string representation of the deque in FIFO order.
//
// Time complexity: O(n).
func (d *Deque[T]) String() string {
	var sb strings.Builder

	sb.WriteString("Deque[")

	for i := range d.len {
		if i > 0 {
			sb.WriteString(", ")
		}

		fmt.Fprintf(&sb, "%v", d.buf[d.wrap(d.start+i)])
	}

	sb.WriteString("]")

	return sb.String()
}

// --------------------------------------------------------------------------------
// Private Helpers

// next calculates the next index in the circular buffer.
func (d *Deque[T]) next(idx int) int {
	return (idx + 1) % d.capacity
}

// prev calculates the previous index in the circular buffer.
func (d *Deque[T]) prev(idx int) int {
	return (idx - 1 + d.capacity) % d.capacity
}

// wrap ensures the index stays within buffer bounds.
func (d *Deque[T]) wrap(idx int) int {
	return idx % d.capacity
}

// grow doubles the capacity of the deque when full (expansion mode only).
//
// Copies existing elements to a new buffer in FIFO order. Time complexity: O(n).
func (d *Deque[T]) grow() {
	newCapacity := max(d.capacity*2, minCapacity)

	newBuf := make([]T, newCapacity)
	// Copy elements in FIFO order
	for i := range d.len {
		newBuf[i] = d.buf[d.wrap(d.start+i)]
	}

	// Update deque state
	d.buf = newBuf
	d.start = 0
	d.end = d.len
	d.capacity = newCapacity
}
