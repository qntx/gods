// Package ringbuf provides JSON serialization and deserialization for the circular buffer queue.
//
// This file extends the Queue type with methods to convert to and from JSON format,
// implementing the container.JSONSerializer and container.JSONDeserializer interfaces.
package ringbuf

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/qntx/gods/container"
)

// --------------------------------------------------------------------------------
// Constants and Errors

// Predefined errors for JSON operations.
var (
	ErrMarshalJSON = errors.New("failed to marshal queue to JSON")
	ErrInvalidJSON = errors.New("invalid JSON data")
)

// --------------------------------------------------------------------------------
// Interface Assertions

// Ensure Queue implements required interfaces at compile time.
var (
	_ container.JSONSerializer   = (*Queue[int])(nil)
	_ container.JSONDeserializer = (*Queue[int])(nil)
	_ json.Marshaler             = (*Queue[int])(nil)
	_ json.Unmarshaler           = (*Queue[int])(nil)
)

// --------------------------------------------------------------------------------
// JSON Serialization Methods

// ToJSON serializes the queue's elements into a JSON array in FIFO order.
//
// Returns the JSON-encoded byte slice or an error if marshaling fails.
// Elements must be JSON-serializable; otherwise, an error is returned.
//
// Example:
//
//	q := New[int](3)
//	q.PushBack(1)
//	q.PushBack(2)
//	data, err := q.ToJSON() // Returns []byte("[1,2]"), nil
//
// Time complexity: O(n), where n is the number of elements.
func (q *Queue[T]) ToJSON() ([]byte, error) {
	data, err := json.Marshal(q.Values())
	if err != nil {
		return nil, fmt.Errorf("ringbuf: %w: %w", ErrMarshalJSON, err)
	}

	return data, nil
}

// FromJSON populates the queue from a JSON array, appending elements to the back.
//
// Expects a valid JSON array (e.g., "[1,2,3]"). Clears the queue before loading
// to ensure a clean state. If the capacity is exceeded, older elements (front) are
// overwritten. Returns an error if the JSON is invalid or elements cannot be
// unmarshaled into type T.
//
// Example:
//
//	q := New[int](2)
//	err := q.FromJSON([]byte("[1,2,3]")) // Queue contains [2,3] after overflow
//
// Time complexity: O(n), where n is the number of elements in the JSON array.
func (q *Queue[T]) FromJSON(data []byte) error {
	var vals []T
	if err := json.Unmarshal(data, &vals); err != nil {
		return fmt.Errorf("ringbuf: %w: %w", ErrInvalidJSON, err)
	}

	q.Clear()

	for _, v := range vals {
		q.PushBack(v)
	}

	return nil
}

// MarshalJSON implements json.Marshaler for seamless JSON encoding.
//
// Delegates to ToJSON() for consistency. Returns the JSON byte slice or an error
// if serialization fails.
//
// Time complexity: O(n), where n is the number of elements.
func (q *Queue[T]) MarshalJSON() ([]byte, error) {
	return q.ToJSON()
}

// UnmarshalJSON implements json.Unmarshaler for seamless JSON decoding.
//
// Delegates to FromJSON() to populate the queue. Returns an error if the JSON
// data is invalid or deserialization fails.
//
// Time complexity: O(n), where n is the number of elements in the JSON array.
func (q *Queue[T]) UnmarshalJSON(data []byte) error {
	return q.FromJSON(data)
}
