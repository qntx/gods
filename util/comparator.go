// Package util provides generic utility functions, including comparators for ordering values.
//
// Comparators are designed for use in data structures like trees or sorting algorithms.
package util

import (
	"math"
	"time"
)

// --------------------------------------------------------------------------------
// Types

// Comparator defines a function for comparing two values of type T.
//
// Returns:
//   - -1 if x < y
//   - 0 if x == y
//   - +1 if x > y
type Comparator[T any] func(x, y T) int

// --------------------------------------------------------------------------------
// Time Comparator

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

// --------------------------------------------------------------------------------
// Float64 Comparators

// Float64Comparator compares two float64 values with an epsilon tolerance.
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
func Float64SimpleReverseComparator(x, y float64) int {
	return Float64ReverseComparator(x, y, Epsilon)
}
