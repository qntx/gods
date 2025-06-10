package linkedhashmap

import (
	"container/list" // Standard library list

	"github.com/qntx/gods/container"
)

// Assert Iterator implementation.
var _ container.ReverseIteratorWithKey[string, int] = (*Iterator[string, int])(nil)

// Iterator holding the iterator's state.
type Iterator[K comparable, V any] struct {
	m       *Map[K, V]    // Reference to the map instance
	current *list.Element // Current element in the ordering list (stores keys K)
}

// Iterator returns a stateful iterator whose elements are key/value pairs.
func (m *Map[K, V]) Iterator() *Iterator[K, V] {
	return &Iterator[K, V]{
		m:       m,
		current: nil, // Iterator starts before the first element or after the last.
	}
}

// Next moves the iterator to the next element and returns true if there was a next element in the container.
// If Next() returns true, then next element's key and value can be retrieved by Key() and Value().
// If Next() was called for the first time after Begin() or instantiation, it points to the first element.
// Modifies the state of the iterator.
func (it *Iterator[K, V]) Next() bool {
	if it.m.Empty() {
		return false
	}

	if it.current == nil { // If at beginning (or after End() and then Prev() went out of bounds)
		it.current = it.m.ordering.Front()
	} else {
		it.current = it.current.Next()
	}

	return it.current != nil
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's key and value can be retrieved by Key() and Value().
// If Prev() was called for the first time after End(), it points to the last element.
// Modifies the state of the iterator.
func (it *Iterator[K, V]) Prev() bool {
	if it.m.Empty() {
		return false
	}

	if it.current == nil { // If at end (or after Begin() and then Next() went out of bounds)
		it.current = it.m.ordering.Back()
	} else {
		it.current = it.current.Prev()
	}

	return it.current != nil
}

// Value returns the current element's value.
// Does not modify the state of the iterator.
// Should be called only after a successful Next(), Prev(), First(), or Last().
func (it *Iterator[K, V]) Value() V {
	if it.current == nil {
		var zeroV V
		// Depending on library's error handling, could panic here.
		// Returning zero value if current is nil.
		return zeroV
	}

	key := it.current.Value.(K) // Key is stored in the list element
	// The map's table stores element[V] which contains the actual value V
	if elemData, found := it.m.table[key]; found {
		return elemData.value
	}

	var zeroV V // Should not happen in a consistent map

	return zeroV
}

// Key returns the current element's key.
// Does not modify the state of the iterator.
// Should be called only after a successful Next(), Prev(), First(), or Last().
func (it *Iterator[K, V]) Key() K {
	if it.current == nil {
		var zeroK K
		// Depending on library's error handling, could panic here.
		return zeroK
	}

	return it.current.Value.(K) // Key is stored directly in the list.Element's Value
}

// Begin resets the iterator to its initial state (one-before-first).
// Call Next() to fetch the first element if any.
func (it *Iterator[K, V]) Begin() {
	it.current = nil
}

// End moves the iterator past the last element (one-past-the-end).
// Call Prev() to fetch the last element if any.
func (it *Iterator[K, V]) End() {
	it.current = nil
}

// First moves the iterator to the first element and returns true if there was a first element in the container.
// If First() returns true, then first element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator.
func (it *Iterator[K, V]) First() bool {
	if it.m.Empty() {
		it.current = nil

		return false
	}

	it.current = it.m.ordering.Front()

	return it.current != nil
}

// Last moves the iterator to the last element and returns true if there was a last element in the container.
// If Last() returns true, then last element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator.
func (it *Iterator[K, V]) Last() bool {
	if it.m.Empty() {
		it.current = nil

		return false
	}

	it.current = it.m.ordering.Back()

	return it.current != nil
}

// NextTo moves the iterator to the next element from current position that satisfies the condition given by the
// passed function, and returns true if there was a next element in the container.
// If NextTo() returns true, then next element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator.
func (it *Iterator[K, V]) NextTo(f func(key K, value V) bool) bool {
	for it.Next() {
		key, value := it.Key(), it.Value()
		if f(key, value) {
			return true
		}
	}

	return false
}

// PrevTo moves the iterator to the previous element from current position that satisfies the condition given by the
// passed function, and returns true if there was a next element in the container.
// If PrevTo() returns true, then next element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator.
func (it *Iterator[K, V]) PrevTo(f func(key K, value V) bool) bool {
	for it.Prev() {
		key, value := it.Key(), it.Value()
		if f(key, value) {
			return true
		}
	}

	return false
}
