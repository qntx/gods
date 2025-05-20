// Package pqueue provides a generic priority queue implementation based on container/heap.
// It supports both min-heap and max-heap configurations, with efficient operations for
// adding, updating, and removing items based on their priorities. The queue uses a map for
// O(1) lookups and supports custom comparators for flexible priority ordering.
//
// Example usage:
//
//	pq := pqueue.New[int, int](pqueue.MinHeap)
//	pq.Put(1, 10) // Add item with value 1 and priority 10
//	item := pq.Get() // Retrieve item with minimum priority
//	fmt.Println(item.Value, item.Priority) // Output: 1 10
package pqueue

import (
	"cmp"
	"container/heap"
	"errors"

	godscmp "github.com/qntx/gods/cmp"
)

// Error messages defined as constants.
var (
	// ErrValueNotFound indicates that a value was not found in the queue.
	ErrValueNotFound = errors.New("value not found in queue")
	// ErrNilComparator indicates the comparator function is nil.
	ErrNilComparator = errors.New("comparator cannot be nil")
	// ErrInvalidItemType indicates an incorrect type was provided to Push.
	ErrInvalidItemType = errors.New("invalid item type")
)

// HeapKind specifies the type of heap: min-heap or max-heap.
type HeapKind int

const (
	// MinHeap yields items with the smallest priority first.
	MinHeap HeapKind = iota
	// MaxHeap yields items with the largest priority first.
	MaxHeap
)

// Item represents an element in the priority queue with a value and priority.
type Item[T comparable, V any] struct {
	Value    T   // Value identifies the item.
	Priority V   // Priority determines the item's order in the queue.
	index    int // index is used internally by the heap.Interface.
}

// PriorityQueue is a generic priority queue implementation using a heap.
// It maintains items with associated priorities, supporting both min-heap and max-heap
// configurations. The queue uses a map for O(1) value lookups and supports custom
// comparators for priority ordering.
type PriorityQueue[T comparable, V cmp.Ordered] struct {
	kind   HeapKind
	heap   []*Item[T, V]
	idxMap map[T]*Item[T, V]
	cmp    godscmp.Comparator[V]
}

// New creates a new priority queue with the default comparator for ordered types.
// It initializes an empty queue with the specified heap kind (MinHeap or MaxHeap).
// The type V must implement cmp.Ordered (e.g., int, float64, string).
//
// Args:
//
//	kind: The heap type (MinHeap or MaxHeap).
//
// Returns:
//
//	A pointer to an initialized PriorityQueue.
//
// Example:
//
//	pq := New[int, int](MinHeap)
//	pq.Put(1, 10)
func New[T comparable, V cmp.Ordered](kind HeapKind) *PriorityQueue[T, V] {
	return NewWith[T](kind, cmp.Compare[V])
}

// NewWith creates a new priority queue with a custom comparator for priorities.
// It initializes an empty queue with the specified heap kind and comparator.
//
// Args:
//
//	kind: The heap type (MinHeap or MaxHeap).
//	cmp: A comparator function for priorities.
//
// Returns:
//
//	A pointer to an initialized PriorityQueue.
//
// Example:
//
//	pq := NewWith[string, int](MaxHeap, cmp.Compare[int])
//	pq.Put("task1", 5)
func NewWith[T comparable, V cmp.Ordered](kind HeapKind, cmp godscmp.Comparator[V]) *PriorityQueue[T, V] {
	if cmp == nil {
		panic(ErrNilComparator)
	}
	pq := &PriorityQueue[T, V]{
		kind:   kind,
		heap:   make([]*Item[T, V], 0, 16), // Pre-allocate for efficiency.
		idxMap: make(map[T]*Item[T, V], 16),
		cmp:    cmp,
	}
	heap.Init(pq)
	return pq
}

// Len returns the number of items in the queue.
// Time complexity: O(1).
func (pq *PriorityQueue[T, V]) Len() int {
	return len(pq.heap)
}

// Less determines the ordering of items based on their priorities and heap kind.
// Time complexity: O(1).
func (pq *PriorityQueue[T, V]) Less(i, j int) bool {
	c := pq.cmp(pq.heap[i].Priority, pq.heap[j].Priority)
	return (pq.kind == MinHeap && c < 0) || (pq.kind == MaxHeap && c > 0)
}

// Swap exchanges two items in the heap and updates their indices.
// Time complexity: O(1).
func (pq *PriorityQueue[T, V]) Swap(i, j int) {
	pq.heap[i], pq.heap[j] = pq.heap[j], pq.heap[i]
	pq.heap[i].index = i
	pq.heap[j].index = j
}

// Push adds an item to the heap.
// Time complexity: O(log n).
func (pq *PriorityQueue[T, V]) Push(x any) {
	item, ok := x.(*Item[T, V])
	if !ok {
		panic(ErrInvalidItemType)
	}
	item.index = len(pq.heap)
	pq.heap = append(pq.heap, item)
	pq.idxMap[item.Value] = item
}

// Pop removes and returns the top item from the heap.
// Time complexity: O(log n).
func (pq *PriorityQueue[T, V]) Pop() any {
	n := len(pq.heap)
	if n == 0 {
		return nil
	}
	item := pq.heap[n-1]
	pq.heap = pq.heap[:n-1]
	delete(pq.idxMap, item.Value)
	return item
}

// Put adds a value with the specified priority to the queue.
// If the value already exists, it updates the priority.
//
// Args:
//
//	value: The item value.
//	priority: The priority associated with the value.
//
// Returns:
//
//	true if the operation was successful, false otherwise.
//
// Time complexity: O(log n).
func (pq *PriorityQueue[T, V]) Put(value T, priority V) {
	if _, exists := pq.idxMap[value]; exists {
		pq.Set(value, priority)
		return
	}
	item := &Item[T, V]{
		Value:    value,
		Priority: priority,
	}
	heap.Push(pq, item)
}

// Set changes the priority of an existing value in the queue.
//
// Args:
//
//	value: The item value to update.
//	priority: The new priority.
//
// Returns:
//
//	true if the operation was successful, false otherwise.
//
// Time complexity: O(log n).
func (pq *PriorityQueue[T, V]) Set(value T, priority V) bool {
	item, exists := pq.idxMap[value]
	if !exists {
		return false
	}
	item.Priority = priority
	heap.Fix(pq, item.index)
	return true
}

// Get removes and returns the item with the highest/lowest priority, based on the heap kind.
// Returns nil if the queue is empty.
// Time complexity: O(log n).
func (pq *PriorityQueue[T, V]) Get() (*Item[T, V], bool) {
	if pq.Empty() {
		return nil, false
	}
	item := heap.Pop(pq)
	if item == nil {
		return nil, false
	}
	return item.(*Item[T, V]), true
}

// Peek returns the item with the highest/lowest priority, based on the heap kind.
// Returns nil if the queue is empty.
// Time complexity: O(1).
func (pq *PriorityQueue[T, V]) Peek() (*Item[T, V], bool) {
	if pq.Empty() {
		return nil, false
	}
	return pq.heap[0], true
}

// Remove removes the item with the specified value from the queue.
// Returns true if the item was removed, false otherwise.
// Time complexity: O(log n).
func (pq *PriorityQueue[T, V]) Remove(value T) bool {
	item, exists := pq.idxMap[value]
	if !exists {
		return false
	}
	heap.Remove(pq, item.index)
	return true
}

// Clear removes all items from the queue and resets its internal state.
// Time complexity: O(1).
func (pq *PriorityQueue[T, V]) Clear() {
	pq.heap = pq.heap[:0]
	pq.idxMap = make(map[T]*Item[T, V], 16)
	heap.Init(pq)
}

// Empty checks if the queue contains no items.
// Time complexity: O(1).
func (pq *PriorityQueue[T, V]) Empty() bool {
	return len(pq.heap) == 0
}

// Items returns a copy of the heap slice containing all queue items.
// This is a safe operation that doesn't expose the internal heap structure.
// Time complexity: O(n).
func (pq *PriorityQueue[T, V]) Items() []*Item[T, V] {
	result := make([]*Item[T, V], len(pq.heap))
	copy(result, pq.heap)
	return result
}

// UnsafeItems returns a direct reference to the internal heap slice.
// WARNING: This is unsafe and should only be used for read-only operations.
// Modifying the returned slice directly may corrupt the heap structure.
// Time complexity: O(1).
func (pq *PriorityQueue[T, V]) UnsafeItems() []*Item[T, V] {
	return pq.heap
}
