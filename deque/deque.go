// Package deque implements a generic, non-thread-safe double-ended queue (deque)
// using a circular buffer.
//
// A circular buffer leverages an array with logically connected ends, enabling
// efficient O(1) insertion and removal from both front and back. It supports two
// modes: overwrite (fixed-size, overwrites oldest elements when full) and expansion
// (growable, doubles capacity when full), selected via a boolean flag.
// Ideal for scenarios requiring bounded or dynamic deques with fast access at both ends.
//
// Reference:
// - https://en.wikipedia.org/wiki/Circular_buffer
// - https://en.wikipedia.org/wiki/Double-ended_queue
package deque

import (
	"errors"
	"fmt"
	"strings"
)

// --------------------------------------------------------------------------------
// Constants and Errors

const (
	minCapacity  = 1 // Minimum allowed capacity for the deque.
	growthFactor = 2 // Factor by which capacity grows when deque is full in expansion mode.
)

// Predefined errors for deque operations.
var (
	ErrInvalidCapacity = errors.New("capacity must be at least 1")
	ErrIndexOutOfRange = errors.New("index out of range")
	ErrEmptyDeque      = errors.New("deque is empty")
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
// The capacity parameter specifies the initial capacity.
// If growable is true, the deque will expand when full; otherwise, it will overwrite
// the oldest elements when full.
//
// Example:
//
//	d := deque.NewFrom([]int{1, 2, 3}, 10, true) // Creates a growable deque with capacity 10 and initial elements [1,2,3]
func NewFrom[T comparable](values []T, capacity int, growable bool) *Deque[T] {
	d := NewWith[T](capacity, growable)

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
//
// Time complexity: O(1) amortized.
func (d *Deque[T]) PushFront(val T) {
	if d.Full() {
		if d.growable {
			d.Grow(d.Cap() * growthFactor)
		} else {
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
//
// Time complexity: O(1) amortized.
func (d *Deque[T]) PushBack(val T) {
	if d.Full() {
		if d.growable {
			d.Grow(d.Cap() * growthFactor)
		} else {
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
//
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
//
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

// Insert adds an element at the specified index, shifting subsequent elements toward the back.
// Index 0 inserts at the front, Len() inserts at the back. If the deque is full and growable,
// the capacity is doubled. Panics if the index is invalid (out of range [0, Len()]).
//
// Time complexity: O(n) where n is the number of elements after the insertion point.
func (d *Deque[T]) Insert(idx int, val T) {
	if idx < 0 || idx > d.len {
		panic(fmt.Errorf("%w [0,%d]: %d", ErrIndexOutOfRange, d.len, idx))
	}

	if idx == 0 {
		d.PushFront(val)
		return
	}
	if idx == d.len {
		d.PushBack(val)
		return
	}

	if d.Full() {
		if d.growable {
			d.Grow(d.Cap() * growthFactor)
		} else {
			d.start = d.next(d.start)
			d.len--
		}
	}

	if idx <= d.len/2 {
		// Shift [0..idx-1] left
		newStart := d.prev(d.start)
		d.buf[newStart] = d.buf[d.start]
		for i := range idx - 1 {
			d.buf[d.wrap(d.start+i)] = d.buf[d.wrap(d.start+i+1)]
		}
		d.buf[d.wrap(d.start+idx-1)] = val
		d.start = newStart
	} else {
		// Shift [idx..len-1] right
		newEnd := d.next(d.end)
		d.buf[d.end] = d.buf[d.prev(d.end)]
		for i := d.len - 1; i > idx; i-- {
			d.buf[d.wrap(d.start+i)] = d.buf[d.wrap(d.start+i-1)]
		}
		d.buf[d.wrap(d.start+idx)] = val
		d.end = newEnd
	}
	d.len++
}

// Remove removes and returns the element at the specified index.
//
// Shifts subsequent elements toward the front to fill the gap. Panics if the index
// is invalid (out of range [0, Len()-1]). Time complexity: O(n) where n is the
// number of elements after the removal point.
func (d *Deque[T]) Remove(idx int) (val T, ok bool) {
	if idx < 0 || idx >= d.len {
		panic(fmt.Errorf("%w [0,%d): %d", ErrIndexOutOfRange, d.len, idx))
	}

	if idx == 0 {
		return d.PopFront()
	}

	if idx == d.len-1 {
		return d.PopBack()
	}

	val = d.buf[d.wrap(d.start+idx)]

	if idx < d.len/2 {
		// Shift [idx+1..len-1] left
		for i := idx; i < d.len-1; i++ {
			d.buf[d.wrap(d.start+i)] = d.buf[d.wrap(d.start+i+1)]
		}
		d.end = d.prev(d.end)
	} else {
		// Shift [0..idx-1] right
		for i := idx; i > 0; i-- {
			d.buf[d.wrap(d.start+i)] = d.buf[d.wrap(d.start+i-1)]
		}
		d.start = d.next(d.start)
	}
	d.len--

	return val, true
}

// Swap exchanges the elements at indices i and j.
//
// Panics if either index is invalid (out of range [0, Len()-1]).
// Time complexity: O(1).
func (d *Deque[T]) Swap(i, j int) {
	if i < 0 || i >= d.len {
		panic(fmt.Errorf("%w: i [0,%d): %d", ErrIndexOutOfRange, d.len, i))
	}
	if j < 0 || j >= d.len {
		panic(fmt.Errorf("%w: j [0,%d): %d", ErrIndexOutOfRange, d.len, j))
	}

	if i == j {
		return
	}
	iPos := d.wrap(d.start + i)
	jPos := d.wrap(d.start + j)

	d.buf[iPos], d.buf[jPos] = d.buf[jPos], d.buf[iPos]
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

// Get retrieves the element at the specified index.
//
// Index 0 is the front, Len()-1 is the back. Panics if the index is invalid.
// Time complexity: O(1).
func (d *Deque[T]) Get(idx int) (val T) {
	if idx < 0 || idx >= d.len {
		panic(fmt.Errorf("%w: idx [0,%d): %d", ErrIndexOutOfRange, d.len, idx))
	}

	return d.buf[d.wrap(d.start+idx)]
}

// Set sets the element at the specified index.
//
// Index 0 is the front, Len()-1 is the back. Panics if the index is invalid.
// Time complexity: O(1).
func (d *Deque[T]) Set(idx int, val T) {
	if idx < 0 || idx >= d.len {
		panic(fmt.Errorf("%w: idx [0,%d): %d", ErrIndexOutOfRange, d.len, idx))
	}

	d.buf[d.wrap(d.start+idx)] = val
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

// Cap returns the current capacity of the deque.
//
// Time complexity: O(1).
func (d *Deque[T]) Cap() int {
	return d.capacity
}

// Clear resets the deque to an empty state.
//
// Preserves capacity and mode but reinitializes the buffer. Time complexity: O(n).
func (d *Deque[T]) Clear() {
	*d = *NewWith[T](d.capacity, d.growable)
}

// Grow doubles the capacity of the deque when full (expansion mode only).
//
// Copies existing elements to a new buffer in FIFO order. Time complexity: O(n).
func (d *Deque[T]) Grow(n int) {
	if n < minCapacity {
		panic(fmt.Errorf("%w: n >= %d: %d", ErrInvalidCapacity, minCapacity, n))
	}

	c := d.Cap()
	l := d.Len()
	// If already big enough.
	if n <= c {
		return
	}

	buf := make([]T, n)
	for i := range l {
		buf[i] = d.buf[d.wrap(d.start+i)]
	}

	d.buf = buf
	d.start = 0
	d.end = l
	d.capacity = n
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
