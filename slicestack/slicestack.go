// Package slicestack implements a stack backed by a Go slice.
//
// Structure is not thread safe.
//
// Reference: https://en.wikipedia.org/wiki/Stack_%28abstract_data_type%29#Array
package slicestack

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/qntx/gods/container"
)

var _ container.Stack[int] = (*Stack[int])(nil)
var _ json.Marshaler = (*Stack[int])(nil)
var _ json.Unmarshaler = (*Stack[int])(nil)

// Stack holds elements in a slice.
type Stack[T comparable] struct {
	elements []T
}

// New instantiates a new empty stack.
func New[T comparable]() *Stack[T] {
	return &Stack[T]{elements: make([]T, 0)}
}

// Push adds a value onto the top of the stack.
func (s *Stack[T]) Push(value T) {
	s.elements = append(s.elements, value)
}

// Pop removes the top element from the stack and returns it.
// The second return parameter is true if a value was popped, or false if the stack was empty.
func (s *Stack[T]) Pop() (value T, ok bool) {
	if s.IsEmpty() {
		return
	}

	index := len(s.elements) - 1
	value = s.elements[index]
	s.elements = s.elements[:index]

	return value, true
}

// Peek returns the top element on the stack without removing it.
// The second return parameter is true if a value was peeked, or false if the stack was empty.
func (s *Stack[T]) Peek() (value T, ok bool) {
	if s.IsEmpty() {
		return
	}

	return s.elements[len(s.elements)-1], true
}

// IsEmpty returns true if the stack does not contain any elements.
func (s *Stack[T]) IsEmpty() bool {
	return len(s.elements) == 0
}

// Len returns the number of elements within the stack.
func (s *Stack[T]) Len() int {
	return len(s.elements)
}

// Clear removes all elements from the stack.
func (s *Stack[T]) Clear() {
	s.elements = make([]T, 0)
}

// Values returns all elements in the stack in LIFO (Last-In, First-Out) order.
// The element at the top of the stack will be the first element in the returned slice.
func (s *Stack[T]) Values() []T {
	size := len(s.elements)
	if size == 0 {
		return []T{}
	}

	// Create a new slice and copy elements in reverse order for LIFO
	reversedValues := make([]T, size)
	for i := range size {
		reversedValues[i] = s.elements[size-1-i]
	}

	return reversedValues
}

// ToSlice returns all elements in the stack in LIFO (Last-In, First-Out) order.
// The element at the top of the stack will be the first element in the returned slice.
func (s *Stack[T]) ToSlice() []T {
	return s.Values()
}

// String returns a string representation of the container.
// Elements are listed from bottom to top of the stack.
func (s *Stack[T]) String() string {
	str := "SliceStack\n"
	values := make([]string, len(s.elements))

	for i, value := range s.elements {
		values[i] = fmt.Sprintf("%v", value)
	}

	str += strings.Join(values, ", ")

	return str
}

// MarshalJSON outputs the JSON representation of the stack.
func (s *Stack[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.elements)
}

// UnmarshalJSON populates the stack from the input JSON representation.
func (s *Stack[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &s.elements)
}
