// Package container provides a set of generic interfaces for working with container data structures.
// It defines common operations for indexed and key-value based collections, enabling consistent
// iteration, filtering, and querying capabilities across different container implementations.
package container

// --------------------------------------------------------------------------------
// Indexed Container Interface

// EnumerableWithIndex defines a generic interface for ordered containers whose elements
// can be accessed by an integer index.
//
// This interface provides methods for iterating, transforming, and querying elements in a
// container. It is designed to work with any type T, making it reusable across various data
// structures like slices, arrays, or custom indexed collections.
//
// Example usage:
//
//	type IntSlice []int
//	func (s IntSlice) Each(f func(int, int)) {
//	    for i, v := range s {
//	        f(i, v)
//	    }
//	}
//	// Implement other methods similarly...
type EnumerableWithIndex[T any] interface {
	// Each invokes the provided function once for each element, passing the element's
	// index and value. The iteration order is implementation-dependent but typically
	// follows the natural order of the container.
	Each(fn func(index int, value T))

	// Any returns true if the provided function returns true for at least one element.
	// It stops iteration as soon as a match is found, optimizing for early exits.
	Any(fn func(index int, value T) bool) bool

	// All returns true if the provided function returns true for every element in the
	// container. It stops and returns false on the first failure.
	All(fn func(index int, value T) bool) bool

	// Find returns the first index and value for which the provided function returns true.
	// If no element satisfies the condition, it returns -1 and the zero value of T.
	Find(fn func(index int, value T) bool) (int, T)
}

// --------------------------------------------------------------------------------
// Key-Value Container Interface

// EnumerableWithKey defines a generic interface for containers whose elements are key-value pairs.
//
// This interface supports iteration and querying over key-value collections, such as maps or
// custom associative data structures. It uses type parameters K and V for keys and values,
// providing type safety and flexibility.
//
// Example usage:
//
//	type StringMap map[string]int
//	func (m StringMap) Each(f func(string, int)) {
//	    for k, v := range m {
//	        f(k, v)
//	    }
//	}
//	// Implement other methods similarly...
type EnumerableWithKey[K, V any] interface {
	// Each invokes the provided function once for each element, passing the element's
	// key and value. The iteration order is implementation-dependent (e.g., maps may be unordered).
	Each(fn func(key K, value V))

	// Any returns true if the provided function returns true for at least one key-value pair.
	// It stops iteration as soon as a match is found, optimizing for early exits.
	Any(fn func(key K, value V) bool) bool

	// All returns true if the provided function returns true for every key-value pair in the
	// container. It stops and returns false on the first failure.
	All(fn func(key K, value V) bool) bool

	// Find returns the first key and value for which the provided function returns true.
	// If no element satisfies the condition, it returns the zero values of K and V.
	Find(fn func(key K, value V) bool) (K, V)
}
