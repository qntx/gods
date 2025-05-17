// Package container_test contains tests for the container package.
//
// It verifies the behavior of the Container interface and its utility functions,
// ensuring they work correctly with various data types and edge cases.
package container_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/qntx/gods/container"
)

// --------------------------------------------------------------------------------
// Test Implementation of Container
// --------------------------------------------------------------------------------

// containerTest is a test implementation of the container.Container interface.
type containerTest[T any] struct {
	values []T
}

// newContainerTest creates a new containerTest instance with the given values.
func newContainerTest[T any](values ...T) *containerTest[T] {
	return &containerTest[T]{values: values}
}

// Empty returns true if the container has no elements.
func (c *containerTest[T]) Empty() bool {
	return len(c.values) == 0
}

// Size returns the number of elements in the container.
func (c *containerTest[T]) Size() int {
	return len(c.values)
}

// Clear removes all elements from the container.
func (c *containerTest[T]) Clear() {
	c.values = nil // Use nil instead of empty slice for cleaner reset
}

// Values returns a slice of all elements in the container.
func (c *containerTest[T]) Values() []T {
	return c.values
}

// String returns a string representation of the container.
func (c *containerTest[T]) String() string {
	var sb strings.Builder

	sb.WriteString("containerTest{")

	for i, v := range c.values {
		if i > 0 {
			sb.WriteString(", ")
		}

		fmt.Fprintf(&sb, "%v", v)
	}

	sb.WriteString("}")

	return sb.String()
}

// --------------------------------------------------------------------------------
// Test Cases
// --------------------------------------------------------------------------------

func TestContainerMethods(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		init      []int
		wantSize  int
		wantEmpty bool
		wantStr   string
	}{
		{name: "empty", init: nil, wantSize: 0, wantEmpty: true, wantStr: "containerTest{}"},
		{name: "single", init: []int{42}, wantSize: 1, wantEmpty: false, wantStr: "containerTest{42}"},
		{name: "multiple", init: []int{1, 2, 3}, wantSize: 3, wantEmpty: false, wantStr: "containerTest{1, 2, 3}"},
	}

	for _, tt := range tests {
		// Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := newContainerTest(tt.init...)

			// Test Empty
			if got := c.Empty(); got != tt.wantEmpty {
				t.Errorf("Empty() = %v, want %v", got, tt.wantEmpty)
			}

			// Test Size
			if got := c.Size(); got != tt.wantSize {
				t.Errorf("Size() = %d, want %d", got, tt.wantSize)
			}

			// Test Values
			if got := c.Values(); len(got) != tt.wantSize {
				t.Errorf("Values() length = %d, want %d", len(got), tt.wantSize)
			}

			// Test String
			if got := c.String(); got != tt.wantStr {
				t.Errorf("String() = %q, want %q", got, tt.wantStr)
			}

			// Test Clear
			c.Clear()

			if !c.Empty() || c.Size() != 0 {
				t.Errorf("Clear() failed: Empty() = %v, Size() = %d", c.Empty(), c.Size())
			}
		})
	}
}

func TestGetSortedValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{name: "empty", input: nil, want: nil},
		{name: "single", input: []int{5}, want: []int{5}},
		{name: "unsorted", input: []int{5, 1, 3, 2, 4}, want: []int{1, 2, 3, 4, 5}},
		{name: "sorted", input: []int{1, 2, 3}, want: []int{1, 2, 3}},
	}

	for _, tt := range tests {
		// Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := newContainerTest(tt.input...)
			got := container.GetSortedValues(c)

			// Check length
			if len(got) != len(tt.want) {
				t.Errorf("GetSortedValues() length = %d, want %d", len(got), len(tt.want))
			}

			// Check sorted order and values
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("GetSortedValues() = %v, want %v", got, tt.want)

					break
				}
			}

			// Verify original container is unchanged
			orig := c.Values()
			if len(orig) != len(tt.input) {
				t.Errorf("Original values modified: got %v, want %v", orig, tt.input)
			}
		})
	}
}

// notInt is a custom type for testing non-ordered values.
type notInt struct {
	i int
}

func TestGetSortedValuesFunc(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input []notInt
		want  []notInt
	}{
		{name: "empty", input: nil, want: nil},
		{name: "single", input: []notInt{{5}}, want: []notInt{{5}}},
		{name: "unsorted", input: []notInt{{5}, {1}, {3}, {2}, {4}}, want: []notInt{{1}, {2}, {3}, {4}, {5}}},
	}

	cmpFunc := func(a, b notInt) int {
		return a.i - b.i // Simple comparison for testing
	}

	for _, tt := range tests {
		// Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := newContainerTest(tt.input...)
			got := container.GetSortedValuesFunc(c, cmpFunc)

			// Check length
			if len(got) != len(tt.want) {
				t.Errorf("GetSortedValuesFunc() length = %d, want %d", len(got), len(tt.want))
			}

			// Check sorted order and values
			for i := range got {
				if got[i].i != tt.want[i].i {
					t.Errorf("GetSortedValuesFunc() = %v, want %v", got, tt.want)

					break
				}
			}

			// Verify original container is unchanged
			orig := c.Values()
			if len(orig) != len(tt.input) {
				t.Errorf("Original values modified: got %v, want %v", orig, tt.input)
			}
		})
	}
}
