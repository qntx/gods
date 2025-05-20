// Package container provides generic iterator interfaces for traversing container data structures.
// It includes forward and reverse iterators for both indexed and key-value based collections,
// enabling flexible and type-safe iteration over various container implementations.
package container

// IteratorWithIndex defines a generic, stateful iterator for ordered containers with indexed elements.
//
// This interface allows forward traversal of a container using integer indices. It maintains an
// internal cursor that can be moved to specific positions or advanced incrementally.
//
// Example usage:
//
//	type IntSlice []int
//	func (s IntSlice) Next() bool { ... }
//	func (s IntSlice) Value() int { ... }
//	// Implement other methods similarly...
type IteratorWithIndex[T any] interface {
	// Next advances the iterator to the next element and returns true if a next element exists.
	// On the first call, it positions the iterator at the first element if the container is non-empty.
	// The current index and value can then be retrieved with Index() and Value().
	Next() bool

	// Value returns the current element's value without modifying the iterator's state.
	Value() T

	// Index returns the current element's index without modifying the iterator's state.
	Index() int

	// Begin resets the iterator to its initial state, positioning it before the first element.
	// Call Next() to move to the first element if it exists.
	Begin()

	// First moves the iterator directly to the first element and returns true if one exists.
	// The first element's index and value can then be retrieved with Index() and Value().
	First() bool

	// NextTo advances the iterator to the next element that satisfies the given condition,
	// returning true if such an element is found. The matching element's index and value
	// can then be retrieved with Index() and Value().
	NextTo(fn func(index int, value T) bool) bool
}

// IteratorWithKey defines a generic, stateful iterator for containers with key-value pairs.
//
// This interface enables forward traversal of key-value collections, such as maps or custom
// associative structures, using type parameters K and V for type safety.
//
// Example usage:
//
//	type StringMap map[string]int
//	func (m StringMap) Next() bool { ... }
//	func (m StringMap) Key() string { ... }
//	// Implement other methods similarly...
type IteratorWithKey[K, V any] interface {
	// Next advances the iterator to the next element and returns true if a next element exists.
	// On the first call, it positions the iterator at the first element if the container is non-empty.
	// The current key and value can then be retrieved with Key() and Value().
	Next() bool

	// Value returns the current element's value without modifying the iterator's state.
	Value() V

	// Key returns the current element's key without modifying the iterator's state.
	Key() K

	// Begin resets the iterator to its initial state, positioning it before the first element.
	// Call Next() to move to the first element if it exists.
	Begin()

	// First moves the iterator directly to the first element and returns true if one exists.
	// The first element's key and value can then be retrieved with Key() and Value().
	First() bool

	// NextTo advances the iterator to the next element that satisfies the given condition,
	// returning true if such an element is found. The matching element's key and value
	// can then be retrieved with Key() and Value().
	NextTo(fn func(key K, value V) bool) bool
}

// ReverseIteratorWithIndex extends IteratorWithIndex with reverse traversal capabilities.
//
// This interface adds methods for backward iteration, including moving to the last element
// and traversing to previous elements that satisfy specific conditions.
//
// It embeds IteratorWithIndex[T] to inherit its forward traversal methods.
type ReverseIteratorWithIndex[T any] interface {
	// Prev moves the iterator to the previous element and returns true if a previous element exists.
	// The previous element's index and value can then be retrieved with Index() and Value().
	Prev() bool

	// End positions the iterator past the last element (one-past-the-end).
	// Call Prev() to move to the last element if it exists.
	End()

	// Last moves the iterator directly to the last element and returns true if one exists.
	// The last element's index and value can then be retrieved with Index() and Value().
	Last() bool

	// PrevTo moves the iterator to the previous element that satisfies the given condition,
	// returning true if such an element is found. The matching element's index and value
	// can then be retrieved with Index() and Value().
	PrevTo(fn func(index int, value T) bool) bool

	IteratorWithIndex[T]
}

// ReverseIteratorWithKey extends IteratorWithKey with reverse traversal capabilities.
//
// This interface adds methods for backward iteration over key-value pairs, including moving
// to the last element and traversing to previous elements that satisfy specific conditions.
//
// It embeds IteratorWithKey[K, V] to inherit its forward traversal methods.
type ReverseIteratorWithKey[K, V any] interface {
	// Prev moves the iterator to the previous element and returns true if a previous element exists.
	// The previous element's key and value can then be retrieved with Key() and Value().
	Prev() bool

	// End positions the iterator past the last element (one-past-the-end).
	// Call Prev() to move to the last element if it exists.
	End()

	// Last moves the iterator directly to the last element and returns true if one exists.
	// The last element's key and value can then be retrieved with Key() and Value().
	Last() bool

	// PrevTo moves the iterator to the previous element that satisfies the given condition,
	// returning true if such an element is found. The matching element's key and value
	// can then be retrieved with Key() and Value().
	PrevTo(fn func(key K, value V) bool) bool

	IteratorWithKey[K, V]
}
