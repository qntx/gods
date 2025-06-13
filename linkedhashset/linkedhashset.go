// Package linkedhashset is a set that preserves insertion-order.
//
// It is backed by a hash table to store values and doubly-linked list to store ordering.
//
// Note that insertion-order is not affected if an element is re-inserted into the set.
//
// Structure is not thread safe.
//
// References: http://en.wikipedia.org/wiki/Set_%28abstract_data_type%29
package linkedhashset

import (
	"container/list"
	"encoding/json"
	"fmt"
	"iter"
	"strings"

	"github.com/qntx/gods/container"
)

const defaultSize = 16

var _ container.Set[int] = (*Set[int])(nil)
var _ json.Marshaler = (*Set[int])(nil)
var _ json.Unmarshaler = (*Set[int])(nil)

// Set holds elements in go's native map.
type Set[T comparable] struct {
	table    map[T]*list.Element
	ordering *list.List
}

// New instantiates a new empty set and adds the passed values, if any, to the set.
func New[T comparable]() *Set[T] {
	return NewWith[T](defaultSize)
}

// NewFrom instantiates a new empty set and adds the passed values, if any, to the set.
func NewFrom[T comparable](values ...T) *Set[T] {
	return NewWith[T](defaultSize, values...)
}

// NewWith instantiates a new empty set and adds the passed values, if any, to the set.
func NewWith[T comparable](size int, values ...T) *Set[T] {
	set := &Set[T]{
		table:    make(map[T]*list.Element, size),
		ordering: list.New(),
	}

	set.Append(values...)

	return set
}

// Add adds the items (one or more) to the set.
// Note that insertion-order is not affected if an element is re-inserted into the set.
func (set *Set[T]) Add(item T) bool {
	if _, contains := set.table[item]; !contains {
		element := set.ordering.PushBack(item)
		set.table[item] = element

		return true
	}

	return false
}

// Append adds the items (one or more) to the set.
// Note that insertion-order is not affected if an element is re-inserted into the set.
func (set *Set[T]) Append(items ...T) int {
	prevlen := set.Len()

	for _, item := range items {
		if _, contains := set.table[item]; !contains {
			element := set.ordering.PushBack(item)
			set.table[item] = element
		}
	}

	return set.Len() - prevlen
}

// Remove removes the item from the set.
// This operation is now O(1) on average due to direct element access.
func (set *Set[T]) Remove(item T) {
	if element, contains := set.table[item]; contains {
		set.ordering.Remove(element)
		delete(set.table, item)
	}
}

// RemoveAll removes the items (one or more) from the set.
// This operation is now O(1) on average due to direct element access.
func (set *Set[T]) RemoveAll(items ...T) {
	for _, item := range items {
		if element, contains := set.table[item]; contains {
			set.ordering.Remove(element)
			delete(set.table, item)
		}
	}
}

// Pop removes and returns an arbitrary item from the set.
func (set *Set[T]) Pop() (v T, ok bool) {
	if element := set.ordering.Front(); element != nil {
		set.Remove(element.Value.(T))

		return element.Value.(T), true
	}

	return
}

// ContainsOne check if item is present in the set.
func (set *Set[T]) ContainsOne(item T) bool {
	_, contains := set.table[item]

	return contains
}

// Contains check if items (one or more) are present in the set.
// All items have to be present in the set for the method to return true.
// Returns true if no arguments are passed at all, i.e. set is always superset of empty set.
func (set *Set[T]) Contains(items ...T) bool {
	for _, item := range items {
		if _, contains := set.table[item]; !contains {
			return false
		}
	}

	return true
}

// ContainsAny returns whether at least one of the
// given items are in the set.
func (set *Set[T]) ContainsAny(val ...T) bool {
	for _, item := range val {
		if _, contains := set.table[item]; contains {
			return true
		}
	}

	return false
}

// ContainsAnyElement returns whether at least one of the
// given element are in the set.
func (set *Set[T]) ContainsAnyElement(other container.Set[T]) bool {
	for item := range other.Iter() {
		if _, contains := set.table[item]; contains {
			return true
		}
	}

	return false
}

// Equal returns whether the given set is equal to this set.
func (set *Set[T]) Equal(other container.Set[T]) bool {
	if set.Len() != other.Len() {
		return false
	}

	for item := range set.Iter() {
		if contains := other.ContainsOne(item); !contains {
			return false
		}
	}

	return true
}

// IsSubset returns whether this set is a subset of the other set.
func (set *Set[T]) IsSubset(other container.Set[T]) bool {
	for item := range set.Iter() {
		if contains := other.ContainsOne(item); !contains {
			return false
		}
	}

	return true
}

// IsProperSubset returns whether this set is a proper subset of the other set.
func (set *Set[T]) IsProperSubset(other container.Set[T]) bool {
	return set.IsSubset(other) && set.Len() < other.Len()
}

// IsSuperset returns whether this set is a superset of the other set.
func (set *Set[T]) IsSuperset(other container.Set[T]) bool {
	for item := range other.Iter() {
		if _, contains := set.table[item]; !contains {
			return false
		}
	}

	return true
}

// IsProperSuperset returns whether this set is a proper superset of the other set.
func (set *Set[T]) IsProperSuperset(other container.Set[T]) bool {
	return set.IsSuperset(other) && set.Len() > other.Len()
}

// Union returns the union of two sets.
// The new set consists of all elements that are in "set" or "another" (possibly both).
// Ref: https://en.wikipedia.org/wiki/Union_(set_theory)
func (set *Set[T]) Union(another container.Set[T]) container.Set[T] {
	result := New[T]()

	for item := range set.Iter() {
		result.Add(item)
	}

	for item := range another.Iter() {
		result.Add(item)
	}

	return result
}

// Intersect returns the intersection between two sets.
// The new set consists of all elements that are both in "set" and "another".
// Ref: https://en.wikipedia.org/wiki/Intersection_(set_theory)
func (set *Set[T]) Intersect(another container.Set[T]) container.Set[T] {
	result := New[T]()

	if set.Len() <= another.Len() {
		for item := range set.Iter() {
			if contains := another.ContainsOne(item); contains {
				result.Add(item)
			}
		}
	} else {
		for item := range another.Iter() {
			if contains := set.ContainsOne(item); contains {
				result.Add(item)
			}
		}
	}

	return result
}

// Difference returns the difference between two sets.
// The new set consists of all elements that are in "set" but not in "another".
// Ref: https://proofwiki.org/wiki/Definition:Set_Difference
func (set *Set[T]) Difference(another container.Set[T]) container.Set[T] {
	result := New[T]()

	for item := range set.Iter() {
		if contains := another.ContainsOne(item); !contains {
			result.Add(item)
		}
	}

	return result
}

// SymmetricDifference returns the symmetric difference between two sets.
// The new set consists of all elements that are in "set" or "another" but not in both.
// Ref: https://proofwiki.org/wiki/Definition:Set_Difference
func (set *Set[T]) SymmetricDifference(another container.Set[T]) container.Set[T] {
	result := New[T]()

	for item := range set.Iter() {
		if contains := another.ContainsOne(item); !contains {
			result.Add(item)
		}
	}

	for item := range another.Iter() {
		if contains := set.ContainsOne(item); !contains {
			result.Add(item)
		}
	}

	return result
}

// IsEmpty returns true if set does not contain any elements.
func (set *Set[T]) IsEmpty() bool {
	return set.Len() == 0
}

// Len returns number of elements within the set.
func (set *Set[T]) Len() int {
	return set.ordering.Len()
}

// Clear clears all values in the set.
func (set *Set[T]) Clear() {
	set.table = make(map[T]*list.Element)
	set.ordering.Init()
}

// Values returns all items in the set.
func (set *Set[T]) Values() []T {
	values := make([]T, 0, set.Len())
	for element := set.ordering.Front(); element != nil; element = element.Next() {
		values = append(values, element.Value.(T))
	}

	return values
}

// ToSlice returns all items in the set as a slice.
func (set *Set[T]) ToSlice() []T {
	return set.Values()
}

// Iter returns an iterator over the values of the set, in insertion order.
func (set *Set[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for e := set.ordering.Front(); e != nil; e = e.Next() {
			if !yield(e.Value.(T)) {
				return
			}
		}
	}
}

// Clone returns a clone of the set using the same
// implementation, duplicating all keys.
func (set *Set[T]) Clone() container.Set[T] {
	return NewFrom(set.Values()...)
}

// MarshalJSON outputs the JSON representation of the set.
func (set *Set[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(set.Values())
}

// UnmarshalJSON populates the set from the input JSON representation.
func (set *Set[T]) UnmarshalJSON(data []byte) error {
	var elements []T

	err := json.Unmarshal(data, &elements)
	if err == nil {
		set.Clear()
		set.Append(elements...)
	}

	return err
}

// String returns a string representation of container.
func (set *Set[T]) String() string {
	str := "LinkedHashSet\n"

	items := []string{}
	for element := set.ordering.Front(); element != nil; element = element.Next() {
		items = append(items, fmt.Sprintf("%v", element.Value.(T)))
	}

	str += strings.Join(items, ", ")

	return str
}
