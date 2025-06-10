package hashset

import (
	"encoding/json"
	"fmt"
	"iter"
	"strings"

	"github.com/qntx/gods/container"
)

type Set[T comparable] map[T]struct{}

var _ container.Set[string] = (*Set[string])(nil)
var _ json.Marshaler = (*Set[string])(nil)
var _ json.Unmarshaler = (*Set[string])(nil)

func New[T comparable](vals ...T) *Set[T] {
	t := make(Set[T])
	for _, v := range vals {
		t.Add(v)
	}

	return &t
}

func NewWith[T comparable](cardinality int, vals ...T) *Set[T] {
	t := make(Set[T], cardinality)
	for _, v := range vals {
		t.Add(v)
	}

	return &t
}

func (s Set[T]) Add(v T) bool {
	prevLen := len(s)
	s[v] = struct{}{}

	return prevLen != len(s)
}

func (s *Set[T]) Append(v ...T) int {
	prevLen := len(*s)

	for _, val := range v {
		(*s)[val] = struct{}{}
	}

	return len(*s) - prevLen
}

// private version of Add which doesn't return a value.
func (s *Set[T]) add(v T) {
	(*s)[v] = struct{}{}
}

func (s *Set[T]) Len() int {
	return len(*s)
}

func (s *Set[T]) Clear() {
	// Constructions like this are optimised by compiler, and replaced by
	// mapclear() function, defined in
	// https://github.com/golang/go/blob/29bbca5c2c1ad41b2a9747890d183b6dd3a4ace4/src/runtime/map.go#L993)
	for key := range *s {
		delete(*s, key)
	}
}

func (s *Set[T]) Clone() container.Set[T] {
	clonedSet := NewWith[T](s.Len())
	for elem := range *s {
		clonedSet.add(elem)
	}

	return clonedSet
}

func (s *Set[T]) Contains(v ...T) bool {
	for _, val := range v {
		if _, ok := (*s)[val]; !ok {
			return false
		}
	}

	return true
}

func (s *Set[T]) ContainsOne(v T) bool {
	_, ok := (*s)[v]

	return ok
}

func (s *Set[T]) ContainsAny(v ...T) bool {
	for _, val := range v {
		if _, ok := (*s)[val]; ok {
			return true
		}
	}

	return false
}

func (s *Set[T]) ContainsAnyElement(other container.Set[T]) bool {
	o := other.(*Set[T])

	// loop over smaller set
	if s.Len() < other.Len() {
		for elem := range *s {
			if o.contains(elem) {
				return true
			}
		}
	} else {
		for elem := range *o {
			if s.contains(elem) {
				return true
			}
		}
	}

	return false
}

// private version of Contains for a single element v.
func (s *Set[T]) contains(v T) (ok bool) {
	_, ok = (*s)[v]

	return ok
}

func (s *Set[T]) Difference(other container.Set[T]) container.Set[T] {
	o := other.(*Set[T])

	diff := New[T]()

	for elem := range *s {
		if !o.contains(elem) {
			diff.add(elem)
		}
	}

	return diff
}

func (s *Set[T]) Each(cb func(T) bool) {
	for elem := range *s {
		if cb(elem) {
			break
		}
	}
}

func (s *Set[T]) Equal(other container.Set[T]) bool {
	o := other.(*Set[T])

	if s.Len() != other.Len() {
		return false
	}

	for elem := range *s {
		if !o.contains(elem) {
			return false
		}
	}

	return true
}

func (s *Set[T]) Intersect(other container.Set[T]) container.Set[T] {
	o := other.(*Set[T])

	intersection := New[T]()
	// loop over smaller set
	if s.Len() < other.Len() {
		for elem := range *s {
			if o.contains(elem) {
				intersection.add(elem)
			}
		}
	} else {
		for elem := range *o {
			if s.contains(elem) {
				intersection.add(elem)
			}
		}
	}

	return intersection
}

func (s *Set[T]) IsEmpty() bool {
	return s.Len() == 0
}

func (s *Set[T]) IsProperSubset(other container.Set[T]) bool {
	return s.Len() < other.Len() && s.IsSubset(other)
}

func (s *Set[T]) IsProperSuperset(other container.Set[T]) bool {
	return s.Len() > other.Len() && s.IsSuperset(other)
}

func (s *Set[T]) IsSubset(other container.Set[T]) bool {
	o := other.(*Set[T])

	if s.Len() > other.Len() {
		return false
	}

	for elem := range *s {
		if !o.contains(elem) {
			return false
		}
	}

	return true
}

func (s *Set[T]) IsSuperset(other container.Set[T]) bool {
	return other.IsSubset(s)
}

// Iter returns an iterator over all elements in the set in sorted order.
func (s *Set[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		s.Each(func(t T) bool {
			return !yield(t)
		})
	}
}

// if set is already empty.
func (s *Set[T]) Pop() (v T, ok bool) {
	for item := range *s {
		delete(*s, item)

		return item, true
	}

	return v, false
}

func (s Set[T]) Remove(v T) {
	delete(s, v)
}

func (s Set[T]) RemoveAll(i ...T) {
	for _, elem := range i {
		delete(s, elem)
	}
}

func (s Set[T]) MarshalJSON() ([]byte, error) {
	items := make([]string, 0, s.Len())

	for elem := range s {
		b, err := json.Marshal(elem)
		if err != nil {
			return nil, err
		}

		items = append(items, string(b))
	}

	return []byte(fmt.Sprintf("[%s]", strings.Join(items, ","))), nil
}

func (s *Set[T]) UnmarshalJSON(b []byte) error {
	var i []T

	err := json.Unmarshal(b, &i)
	if err != nil {
		return err
	}

	s.Append(i...)

	return nil
}

func (s Set[T]) String() string {
	items := make([]string, 0, len(s))

	for elem := range s {
		items = append(items, fmt.Sprintf("%v", elem))
	}

	return fmt.Sprintf("Set{%s}", strings.Join(items, ", "))
}

func (s *Set[T]) SymmetricDifference(other container.Set[T]) container.Set[T] {
	o := other.(*Set[T])

	sd := New[T]()

	for elem := range *s {
		if !o.contains(elem) {
			sd.add(elem)
		}
	}

	for elem := range *o {
		if !s.contains(elem) {
			sd.add(elem)
		}
	}

	return sd
}

func (s Set[T]) ToSlice() []T {
	keys := make([]T, 0, s.Len())
	for elem := range s {
		keys = append(keys, elem)
	}

	return keys
}

func (s Set[T]) Union(other container.Set[T]) container.Set[T] {
	o := other.(*Set[T])

	n := s.Len()
	if o.Len() > n {
		n = o.Len()
	}

	unionedSet := make(Set[T], n)

	for elem := range s {
		unionedSet.add(elem)
	}

	for elem := range *o {
		unionedSet.add(elem)
	}

	return &unionedSet
}
