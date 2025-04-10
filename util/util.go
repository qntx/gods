// Package util provides common utility functions for general use.
//
// Includes functionalities such as:
//   - Sorting and comparators
//   - Type conversions
package util

import (
	"fmt"
	"strconv"
)

// --------------------------------------------------------------------------------
// Constants

// Epsilon is the default tolerance for floating-point comparisons.
//
// Used in comparators to handle precision issues. Value is 1e-15.
const Epsilon = 1e-15

// --------------------------------------------------------------------------------
// Conversion Functions

// ToString converts a value of any type to its string representation.
//
// Supports common built-in types with specific formatting:
//   - Integers (int8, int16, int32, int64, uint8, uint16, uint32, uint64)
//   - Floats (float32, float64) with 'g' format for concise output
//   - Boolean values
//   - Strings (returned as-is)
//   - Other types via fmt.Sprintf with %+v verb
//
// Time complexity: O(1) for most types, O(n) for complex types handled by fmt.Sprintf.
//
// Example:
//
//	ToString(42)       // Returns "42"
//	ToString(3.14)     // Returns "3.14"
//	ToString(true)     // Returns "true"
//	ToString("hello")  // Returns "hello"
func ToString(val any) string {
	switch v := val.(type) {
	case string:
		return v
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("%+v", v)
	}
}
