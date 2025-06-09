package container

import "iter"

// Set is the primary interface provided by the mapset package. It
// represents an unordered set of data and a large number of
// operations that can be applied to that set.
type Set[T comparable] interface {
	Container[T]

	// Clone returns a clone of the set using the same
	// implementation, duplicating all keys.
	Clone() Set[T]

	// Add adds an element to the set. Returns whether
	// the item was added.
	Add(val T) bool

	// Append multiple elements to the set. Returns
	// the number of elements added.
	Append(val ...T) int

	// Remove removes a single element from the set.
	Remove(i T)

	// RemoveAll removes multiple elements from the set.
	RemoveAll(i ...T)

	// Pop removes and returns an arbitrary item from the set.
	Pop() (T, bool)

	// ContainsOne returns whether the given item
	// is in the set.
	//
	// Contains may cause the argument to escape to the heap.
	// See: https://github.com/deckarep/golang-set/issues/118
	ContainsOne(val T) bool

	// Contains returns whether the given items
	// are all in the set.
	Contains(val ...T) bool

	// ContainsAny returns whether at least one of the
	// given items are in the set.
	ContainsAny(val ...T) bool

	// ContainsAnyElement returns whether at least one of the
	// given element are in the set.
	ContainsAnyElement(other Set[T]) bool

	// Equal determines if two sets are equal to each
	// other. If they have the same cardinality
	// and contain the same elements, they are
	// considered equal. The order in which
	// the elements were added is irrelevant.
	//
	// Note that the argument to Equal must be
	// of the same type as the receiver of the
	// method. Otherwise, Equal will panic.
	Equal(other Set[T]) bool

	// IsSubset determines if every element in this set is in
	// the other set.
	//
	// Note that the argument to IsSubset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsSubset will panic.
	IsSubset(other Set[T]) bool

	// IsProperSubset determines if every element in this set is in
	// the other set but the two sets are not equal.
	//
	// Note that the argument to IsProperSubset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsProperSubset will panic.
	IsProperSubset(other Set[T]) bool

	// IsSuperset determines if every element in the other set
	// is in this set.
	//
	// Note that the argument to IsSuperset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsSuperset will panic.
	IsSuperset(other Set[T]) bool

	// IsProperSuperset determines if every element in the other set
	// is in this set but the two sets are not
	// equal.
	//
	// Note that the argument to IsSuperset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsSuperset will
	// panic.
	IsProperSuperset(other Set[T]) bool

	// Union returns a new set with all elements in both sets.
	//
	// Note that the argument to Union must be of the
	// same type as the receiver of the method.
	// Otherwise, Union will panic.
	Union(other Set[T]) Set[T]

	// Intersect returns a new set containing only the elements
	// that exist only in both sets.
	//
	// Note that the argument to Intersect
	// must be of the same type as the receiver
	// of the method. Otherwise, Intersect will
	// panic.
	Intersect(other Set[T]) Set[T]

	// Difference returns the difference between this set
	// and other. The returned set will contain
	// all elements of this set that are not also
	// elements of other.
	//
	// Note that the argument to Difference
	// must be of the same type as the receiver
	// of the method. Otherwise, Difference will
	// panic.
	Difference(other Set[T]) Set[T]

	// SymmetricDifference returns a new set with all elements which are
	// in either this set or the other set but not in both.
	//
	// Note that the argument to SymmetricDifference
	// must be of the same type as the receiver
	// of the method. Otherwise, SymmetricDifference will
	// panic.
	SymmetricDifference(other Set[T]) Set[T]

	// Each iterates over elements and executes the passed func against each element.
	// If passed func returns true, stop iteration at the time.
	Each(func(T) bool)

	// Iter returns a channel of elements that you can
	// range over.
	Iter() iter.Seq[T]
}
