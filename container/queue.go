package container

import "github.com/qntx/gods/cmp"

// Queue is a generic interface for a first-in, first-out (FIFO) data structure.
// It supports adding elements to the back and removing them from the front.
// Implementations (e.g., array-based or linked-list queues) must provide all
// operations defined here, including those inherited from Container[T] (e.g., Len, IsEmpty, Clear).
// Type parameter T must be comparable to enable equality checks for elements.
type Queue[T comparable] interface {
	Container[T]

	// Enqueue adds an element to the back of the queue.
	Enqueue(val T)

	// Dequeue removes and returns the front element of the queue.
	// Returns the element and true if the queue is non-empty,
	// or the zero value of T and false if the queue is empty.
	Dequeue() (val T, ok bool)

	// Peek returns the front element of the queue without removing it.
	// Returns the element and true if the queue is non-empty,
	// or the zero value of T and false if the queue is empty.
	Peek() (val T, ok bool)
}

// PQueue is a generic interface for a priority queue.
// It supports adding elements to the back and removing them from the front.
// Implementations (e.g., array-based or linked-list priority queues) must provide all
// operations defined here, including those inherited from Container[T] (e.g., Len, IsEmpty, Clear).
// Type parameter T must be comparable to enable equality checks for elements.
type PQueue[T comparable, V cmp.Ordered] interface {
	Container[T]

	// Enqueue adds an element to the queue.
	Enqueue(val T, priority V)

	// Dequeue removes and returns the front element of the queue.
	// Returns the element and true if the queue is non-empty,
	// or the zero value of T and false if the queue is empty.
	Dequeue() (val T, priority V, ok bool)

	// Peek returns the front element of the queue without removing it.
	// Returns the element and true if the queue is non-empty,
	// or the zero value of T and false if the queue is empty.
	Peek() (val T, priority V, ok bool)
}

// Deque is a generic interface for a double-ended queue, allowing
// addition and removal of elements at both the front and back.
// Implementations (e.g., circular buffers, doubly-linked lists) must
// support all operations defined here, including those inherited from
// Container[T] (e.g., Len, IsEmpty, Clear).
// Type parameter T must be comparable to enable equality checks for elements.
type Deque[T comparable] interface {
	Container[T]

	// PushFront adds an element to the front of the deque.
	PushFront(val T)

	// PushBack adds an element to the back of the deque.
	PushBack(val T)

	// PopFront removes and returns the front element of the deque.
	// Returns the element and true if the deque is non-empty,
	// or the zero value of T and false if the deque is empty.
	PopFront() (val T, ok bool)

	// PopBack removes and returns the back element of the deque.
	// Returns the element and true if the deque is non-empty,
	// or the zero value of T and false if the deque is empty.
	PopBack() (val T, ok bool)

	// Front returns the front element of the deque without removing it.
	// Returns the element and true if the deque is non-empty,
	// or the zero value of T and false if the deque is empty.
	Front() (val T, ok bool)

	// Back returns the back element of the deque without removing it.
	// Returns the element and true if the deque is non-empty,
	// or the zero value of T and false if the deque is empty.
	Back() (val T, ok bool)

	// Capacity returns the maximum number of elements the deque can hold,
	// or -1 if the capacity is unbounded (e.g., for linked-list implementations).
	Capacity() int

	// Get returns the element at the specified index, where 0 is the front.
	// Returns the element and true if the index is valid,
	// or the zero value of T and false if the index is out of bounds.
	Get(idx int) (val T, ok bool)

	// Set updates the element at the specified index, where 0 is the front.
	// Panics if the index is out of bounds.
	Set(idx int, val T)

	// Insert adds an element at the specified index, shifting subsequent elements.
	// Index 0 inserts at the front, and Len() inserts at the back.
	// Panics if the index is out of bounds (i.e., idx < 0 or idx > Len()).
	Insert(idx int, val T)

	// Remove removes the element at the specified index, shifting subsequent elements.
	// Panics if the index is out of bounds.
	Remove(idx int) (val T, ok bool)

	// Swap exchanges the elements at the specified indices.
	// Panics if either index is out of bounds.
	Swap(idx1, idx2 int)
}
