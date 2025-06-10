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
)

// Set holds elements in go's native map.
type Set[T comparable] struct {
	table    map[T]*list.Element
	ordering *list.List
}

// var _ container.Set = (*Set[string])(nil).
var _ json.Marshaler = (*Set[string])(nil)
var _ json.Unmarshaler = (*Set[string])(nil)

// New instantiates a new empty set and adds the passed values, if any, to the set.
func New[T comparable](values ...T) *Set[T] {
	set := &Set[T]{
		table:    make(map[T]*list.Element),
		ordering: list.New(),
	}
	if len(values) > 0 {
		set.Add(values...)
	}

	return set
}

// Add adds the items (one or more) to the set.
// Note that insertion-order is not affected if an element is re-inserted into the set.
func (set *Set[T]) Add(items ...T) {
	for _, item := range items {
		if _, contains := set.table[item]; !contains {
			element := set.ordering.PushBack(item)
			set.table[item] = element
		}
	}
}

// Remove removes the items (one or more) from the set.
// This operation is now O(1) on average due to direct element access.
func (set *Set[T]) Remove(items ...T) {
	for _, item := range items {
		if element, contains := set.table[item]; contains {
			set.ordering.Remove(element)
			delete(set.table, item)
		}
	}
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

// Empty returns true if set does not contain any elements.
func (set *Set[T]) Empty() bool {
	return set.Size() == 0
}

// Size returns number of elements within the set.
func (set *Set[T]) Size() int {
	return set.ordering.Len()
}

// Clear clears all values in the set.
func (set *Set[T]) Clear() {
	set.table = make(map[T]*list.Element)
	set.ordering.Init()
}

// Values returns all items in the set.
func (set *Set[T]) Values() []T {
	values := make([]T, 0, set.Size())
	for element := set.ordering.Front(); element != nil; element = element.Next() {
		values = append(values, element.Value.(T))
	}

	return values
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
		set.Add(elements...)
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

// Intersection returns the intersection between two sets.
// The new set consists of all elements that are both in "set" and "another".
// Ref: https://en.wikipedia.org/wiki/Intersection_(set_theory)
func (set *Set[T]) Intersection(another *Set[T]) *Set[T] {
	result := New[T]()

	if set.Size() <= another.Size() {
		for item := range set.table {
			if _, contains := another.table[item]; contains {
				result.Add(item)
			}
		}
	} else {
		for item := range another.table {
			if _, contains := set.table[item]; contains {
				result.Add(item)
			}
		}
	}

	return result
}

// Union returns the union of two sets.
// The new set consists of all elements that are in "set" or "another" (possibly both).
// Ref: https://en.wikipedia.org/wiki/Union_(set_theory)
func (set *Set[T]) Union(another *Set[T]) *Set[T] {
	result := New[T]()

	for item := range set.table {
		result.Add(item)
	}

	for item := range another.table {
		result.Add(item)
	}

	return result
}

// Difference returns the difference between two sets.
// The new set consists of all elements that are in "set" but not in "another".
// Ref: https://proofwiki.org/wiki/Definition:Set_Difference
func (set *Set[T]) Difference(another *Set[T]) *Set[T] {
	result := New[T]()

	for item := range set.table {
		if _, contains := another.table[item]; !contains {
			result.Add(item)
		}
	}

	return result
}
