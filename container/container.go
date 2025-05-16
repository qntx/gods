// Package container provides core interfaces and utilities for working with data structures.
//
// It defines the base Container interface for all container types, along with utility functions
// for sorting and manipulating container elements. The package integrates with standard Go
// libraries and supports generic types for flexibility and type safety.
//
// Key features:
//   - Container: Base interface for all data structures.
//   - Iterators: Stateful iteration (defined separately).
//   - Enumerable: Ruby-inspired container functions (defined separately).
//   - Serialization: JSON marshalers and unmarshalers (defined separately).
package container

import (
	"cmp"
	"slices"

	godscmp "github.com/qntx/gods/cmp"
)

// --------------------------------------------------------------------------------
// Base Container Interface

// Container defines the fundamental interface for all container data structures.
//
// This interface provides basic operations for querying and manipulating a container's
// elements, using a generic type T to support any data type.
//
// Example usage:
//
//	type IntList []int
//	func (l IntList) Empty() bool { return len(l) == 0 }
//	func (l IntList) Size() int { return len(l) }
//	func (l IntList) Clear() { l = nil }
//	func (l IntList) Values() []int { return l }
//	func (l IntList) String() string { return fmt.Sprint(l) }
type Container[T any] interface {
	// Empty returns true if the container has no elements.
	Empty() bool

	// Size returns the number of elements in the container.
	Size() int

	// Clear removes all elements from the container, resetting it to an empty state.
	Clear()

	// Values returns a slice containing all elements in the container.
	// The order of elements is implementation-dependent.
	Values() []T

	// String returns a string representation of the container's elements,
	// suitable for logging or debugging.
	String() string
}

// --------------------------------------------------------------------------------
// Sorting Utilities

// GetSortedValues returns a sorted slice of the container's elements for ordered types.
//
// It uses the natural ordering of type T, as defined by the cmp.Ordered constraint.
// The original container remains unchanged, and the returned slice is a new copy.
//
// Returns the original values slice if it has fewer than 2 elements, as sorting is unnecessary.
func GetSortedValues[T cmp.Ordered](c Container[T]) []T {
	values := c.Values()
	if len(values) < 2 {
		return values
	}
	// Create a copy to avoid modifying the original slice
	sorted := make([]T, len(values))
	copy(sorted, values)
	slices.Sort(sorted)

	return sorted
}

// GetSortedValuesFunc returns a sorted slice of the container's elements using a custom comparator.
//
// It is designed for types that do not implement cmp.Ordered, allowing flexible sorting logic
// via the provided comparator function. The original container remains unchanged.
//
// Returns the original values slice if it has fewer than 2 elements, as sorting is unnecessary.
func GetSortedValuesFunc[T any](c Container[T], cmp godscmp.Comparator[T]) []T {
	values := c.Values()
	if len(values) < 2 {
		return values
	}
	// Create a copy to avoid modifying the original slice
	sorted := make([]T, len(values))
	copy(sorted, values)
	slices.SortFunc(sorted, cmp)

	return sorted
}
