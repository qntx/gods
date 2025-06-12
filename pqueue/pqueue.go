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
	"container/heap"
	"errors"
	"fmt"

	"github.com/qntx/gods/cmp"
	"github.com/qntx/gods/container"
)

const defaultCapacity = 16

var (
	ErrInvalidItemType = errors.New("invalid item type")
)

type HeapKind int

const (
	// MinHeap yields items with the smallest priority first.
	MinHeap HeapKind = iota
	// MaxHeap yields items with the largest priority first.
	MaxHeap
)

// Item represents an element in the priority queue with a value and priority.
type Item[T comparable, V any] struct {
	index    int // index is used internally by the heap.Interface.
	Value    T   // Value identifies the item.
	Priority V   // Priority determines the item's order in the queue.
}

var _ container.PQueue[int, int] = (*PriorityQueue[int, int])(nil)

// var _ json.Marshaler = (*PriorityQueue[int, int])(nil)
// var _ json.Unmarshaler = (*PriorityQueue[int, int])(nil)

// PriorityQueue is a generic priority queue implementation using a heap.
// It maintains items with associated priorities, supporting both min-heap and max-heap
// configurations. The queue uses a map for O(1) value lookups and supports custom
// comparators for priority ordering.
type PriorityQueue[T comparable, V cmp.Ordered] struct {
	kind HeapKind
	heap []*Item[T, V]
	idx  map[T]*Item[T, V]
	cmp  cmp.Comparator[V]
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
func NewWith[T comparable, V cmp.Ordered](kind HeapKind, cmp cmp.Comparator[V]) *PriorityQueue[T, V] {
	pq := &PriorityQueue[T, V]{
		kind: kind,
		heap: make([]*Item[T, V], 0, defaultCapacity), // Pre-allocate for efficiency.
		idx:  make(map[T]*Item[T, V], defaultCapacity),
		cmp:  cmp,
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
	pq.idx[item.Value] = item
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
	delete(pq.idx, item.Value)

	return item
}

// Enqueue adds a value with the specified priority to the queue.
// If the value already exists, it updates the priority.
//
// Time complexity: O(log n).
func (pq *PriorityQueue[T, V]) Enqueue(value T, priority V) {
	if _, exists := pq.idx[value]; exists {
		pq.Set(value, priority)

		return
	}

	item := &Item[T, V]{
		Value:    value,
		Priority: priority,
	}
	heap.Push(pq, item)
}

// Dequeue removes and returns the item with the highest/lowest priority, based on the heap kind.
// Returns nil if the queue is empty.
// Time complexity: O(log n).
func (pq *PriorityQueue[T, V]) Dequeue() (value T, priority V, ok bool) {
	if pq.IsEmpty() {
		return
	}

	item := heap.Pop(pq)
	if item == nil {
		return
	}

	return item.(*Item[T, V]).Value, item.(*Item[T, V]).Priority, true
}

// Peek returns the item with the highest/lowest priority, based on the heap kind.
// Returns nil if the queue is empty.
// Time complexity: O(1).
func (pq *PriorityQueue[T, V]) Peek() (value T, priority V, ok bool) {
	if pq.IsEmpty() {
		return
	}

	return pq.heap[0].Value, pq.heap[0].Priority, true
}

// Set changes the priority of an existing value in the queue.
//
// Time complexity: O(log n).
func (pq *PriorityQueue[T, V]) Set(value T, priority V) bool {
	item, exists := pq.idx[value]
	if !exists {
		return false
	}

	item.Priority = priority
	heap.Fix(pq, item.index)

	return true
}

// Remove removes the item with the specified value from the queue.
// Returns true if the item was removed, false otherwise.
// Time complexity: O(log n).
func (pq *PriorityQueue[T, V]) Remove(value T) bool {
	item, exists := pq.idx[value]
	if !exists {
		return false
	}

	heap.Remove(pq, item.index)
	delete(pq.idx, value)

	return true
}

// Clear removes all items from the queue and resets its internal state.
// Time complexity: O(1).
func (pq *PriorityQueue[T, V]) Clear() {
	pq.heap = pq.heap[:0]
	pq.idx = make(map[T]*Item[T, V], defaultCapacity)
	heap.Init(pq)
}

// IsEmpty checks if the queue contains no items.
// Time complexity: O(1).
func (pq *PriorityQueue[T, V]) IsEmpty() bool {
	return len(pq.heap) == 0
}

// Values returns a copy of the values in the queue.
// This is a safe operation that doesn't expose the internal heap structure.
// Time complexity: O(n).
func (pq *PriorityQueue[T, V]) Values() []T {
	result := make([]T, len(pq.heap))
	for i, item := range pq.heap {
		result[i] = item.Value
	}

	return result
}

// ToSlice returns a copy of the heap slice containing all queue items.
// This is a safe operation that doesn't expose the internal heap structure.
// Time complexity: O(n).
func (pq *PriorityQueue[T, V]) ToSlice() []T {
	return pq.Values()
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

// String returns a string representation of the queue.
func (pq *PriorityQueue[T, V]) String() string {
	return fmt.Sprint(pq.heap)
}

// // MarshalJSON implements the json.Marshaler interface.
// func (pq *PriorityQueue[T, V]) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(pq.heap)
// }

// // UnmarshalJSON implements the json.Unmarshaler interface.
// func (pq *PriorityQueue[T, V]) UnmarshalJSON(data []byte) error {
// 	var items []*Item[T, V]
// 	if err := json.Unmarshal(data, &items); err != nil {
// 		return err
// 	}

// 	pq.Clear()
// 	pq.heap = items
// 	for _, item := range items {
// 		pq.idx[item.Value] = item
// 	}
// 	heap.Init(pq)

// 	return nil
// }
