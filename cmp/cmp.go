// Package cmp provides generic utility functions, including comparators for ordering values.
//
// Comparators are designed for use in data structures like trees or sorting algorithms.
package cmp

import (
	"cmp"
	"math"
	"time"
)

// Epsilon is the default tolerance for floating-point comparisons.
//
// Used in comparators to handle precision issues. Value is 1e-15.
const Epsilon = 1e-15

// Ordered is a constraint that permits any ordered type: any type
// that supports the operators < <= >= >.
// If future releases of Go add new ordered types,
// this constraint will be modified to include them.
//
// Note that floating-point types may contain NaN ("not-a-number") values.
// An operator such as == or < will always report false when
// comparing a NaN value with any other value, NaN or not.
// See the [Compare] function for a consistent way to compare NaN values.
type Ordered = cmp.Ordered

// Comparator defines a function for comparing two values of type T.
//
// Returns:
//   - -1 if x < y
//   - 0 if x == y
//   - +1 if x > y
type Comparator[T any] func(x, y T) int

// Less reports whether x is less than y.
// For floating-point types, a NaN is considered less than any non-NaN,
// and -0.0 is not less than (is equal to) 0.0.
func Less[T Ordered](x, y T) bool {
	return (IsNaN(x) && !IsNaN(y)) || x < y
}

// Compare returns
//
//	-1 if x is less than y,
//	 0 if x equals y,
//	+1 if x is greater than y.
//
// For floating-point types, a NaN is considered less than any non-NaN,
// a NaN is considered equal to a NaN, and -0.0 is equal to 0.0.
func Compare[T Ordered](x, y T) int {
	xNaN := IsNaN(x)
	yNaN := IsNaN(y)
	if xNaN {
		if yNaN {
			return 0
		}
		return -1
	}
	if yNaN {
		return +1
	}
	if x < y {
		return -1
	}
	if x > y {
		return +1
	}
	return 0
}

// IsNaN reports whether x is a NaN without requiring the math package.
// This will always return false if T is not floating-point.
func IsNaN[T Ordered](x T) bool {
	return x != x
}

// Or returns the first of its arguments that is not equal to the zero value.
// If no argument is non-zero, it returns the zero value.
func Or[T comparable](vals ...T) T {
	var zero T
	for _, val := range vals {
		if val != zero {
			return val
		}
	}
	return zero
}

// TimeComparator compares two time.Time values.
//
// Uses time.Time's After and Before methods for precise ordering.
// Returns:
//   - 1 if a > b
//   - 0 if a == b
//   - -1 if a < b
//
// Time complexity: O(1).
func TimeComparator(a, b time.Time) int {
	if a.After(b) {
		return 1
	}

	if a.Before(b) {
		return -1
	}

	return 0
}

// Float64Comparator compares two float64 values directly with an epsilon tolerance.
//
// Accounts for floating-point precision by considering values equal if their
// difference is within epsilon. Handles special cases like NaN and ±0 consistently
// with cmp.Compare. Returns:
//   - -1 if x < y
//   - 0 if x ≈ y (within epsilon)
//   - +1 if x > y
//
// Special cases:
//   - NaN < any non-NaN value
//   - NaN == NaN
//   - -0.0 == 0.0
//
// Parameters:
//   - x, y: Values to compare.
//   - epsilon: Tolerance for equality (e.g., 1e-10). If ≤ 0, defaults to 1e-15.
//
// Time complexity: O(1).
func Float64Comparator(x, y, epsilon float64) int {
	// Use default epsilon if provided value is invalid.
	if epsilon <= 0 {
		epsilon = Epsilon
	}

	// Handle NaN cases first.
	switch {
	case math.IsNaN(x) && math.IsNaN(y):
		return 0
	case math.IsNaN(x):
		return -1
	case math.IsNaN(y):
		return 1
	}

	// Check approximate equality within epsilon.
	if math.Abs(x-y) <= epsilon {
		return 0
	}

	// Standard comparison.
	if x < y {
		return -1
	}

	return 1
}

// NewFloat64Comparator creates a Comparator for float64 values with a specified epsilon tolerance.
//
// Creates a closure that remembers the epsilon value and returns a function conforming to
// the Comparator[float64] type. This allows pre-configuring the epsilon tolerance
// without passing it in every comparison operation.
//
// Parameters:
//   - epsilon: Tolerance for equality (e.g., 1e-10). If ≤ 0, defaults to 1e-15.
//
// Returns:
//   - A Comparator[float64] function that compares with the specified epsilon.
//
// Time complexity: O(1) for creation, O(1) for each comparison.
func NewFloat64Comparator(epsilon float64) Comparator[float64] {
	// Create and return a closure that remembers the epsilon value
	return func(x, y float64) int {
		return Float64Comparator(x, y, epsilon)
	}
}

// Float64ReverseComparator compares two float64 values with an epsilon tolerance in reverse order.
//
// Similar to Float64Comparator but returns the opposite result for descending order.
// Handles special cases like NaN and ±0 consistently with cmp.Compare. Returns:
//   - -1 if x > y
//   - 0 if x ≈ y (within epsilon)
//   - +1 if x < y
//
// Special cases:
//   - NaN > any non-NaN value (reversed from Float64Comparator)
//   - NaN == NaN
//   - -0.0 == 0.0
//
// Parameters:
//   - x, y: Values to compare.
//   - epsilon: Tolerance for equality (e.g., 1e-10). If ≤ 0, defaults to 1e-15.
//
// Time complexity: O(1).
func Float64ReverseComparator(x, y, epsilon float64) int {
	// Use default epsilon if provided value is invalid.
	if epsilon <= 0 {
		epsilon = Epsilon
	}

	// Handle NaN cases first, reversing the order.
	switch {
	case math.IsNaN(x) && math.IsNaN(y):
		return 0
	case math.IsNaN(x):
		return 1 // NaN is "greater" in reverse order.
	case math.IsNaN(y):
		return -1
	}

	// Check approximate equality within epsilon.
	if math.Abs(x-y) <= epsilon {
		return 0
	}

	// Reverse the standard comparison.
	if x > y {
		return -1
	}

	return 1
}

// NewFloat64ReverseComparator creates a reverse Comparator for float64 values with a specified epsilon tolerance.
//
// Creates a closure that remembers the epsilon value and returns a function conforming to
// the Comparator[float64] type. This allows pre-configuring the epsilon tolerance
// for descending order comparisons without passing it in every operation.
//
// Parameters:
//   - epsilon: Tolerance for equality (e.g., 1e-10). If ≤ 0, defaults to 1e-15.
//
// Returns:
//   - A Comparator[float64] function that compares with the specified epsilon in reverse order.
//
// Time complexity: O(1) for creation, O(1) for each comparison.
func NewFloat64ReverseComparator(epsilon float64) Comparator[float64] {
	// Create and return a closure that remembers the epsilon value
	return func(x, y float64) int {
		return Float64ReverseComparator(x, y, epsilon)
	}
}

// Float64SimpleComparator is a simplified version of Float64Comparator that implements Comparator[float64].
//
// Uses the default Epsilon value for comparison tolerance.
//
// Returns:
//   - -1 if x < y
//   - 0 if x ≈ y (within Epsilon)
//   - +1 if x > y
//
// Time complexity: O(1).
//
// Note: This is equivalent to calling NewFloat64Comparator(Epsilon) but avoids the closure overhead.
func Float64SimpleComparator(x, y float64) int {
	return Float64Comparator(x, y, Epsilon)
}

// Float64SimpleReverseComparator is a simplified version of Float64ReverseComparator that implements Comparator[float64].
//
// Uses the default Epsilon value for comparison tolerance.
//
// Returns:
//   - -1 if x > y
//   - 0 if x ≈ y (within Epsilon)
//   - +1 if x < y
//
// Time complexity: O(1).
//
// Note: This is equivalent to calling NewFloat64ReverseComparator(Epsilon) but avoids the closure overhead.
func Float64SimpleReverseComparator(x, y float64) int {
	return Float64ReverseComparator(x, y, Epsilon)
}
