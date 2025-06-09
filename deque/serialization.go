// Package deque provides JSON serialization and deserialization for a circular buffer queue.
//
// This package extends the Queue type with methods to convert to and from JSON format,
// implementing the container.JSONSerializer and container.JSONDeserializer interfaces.
package deque

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/qntx/gods/container"
)

var (
	ErrMarshalJSON = errors.New("failed to marshal queue to JSON")
	ErrInvalidJSON = errors.New("invalid JSON data")
)

// Verify Queue satisfies required interfaces at compile time.
var _ container.JSONCodec = (*Deque[int])(nil)

// MarshalJSON serializes the queue's elements into a JSON array in FIFO order.
//
// Elements are marshaled as a JSON array (e.g., "[1,2]"). The method returns an
// error if the elements are not JSON-serializable.
//
// Example:
//
//	q := New[int](3)
//	q.PushBack(1)
//	q.PushBack(2)
//	data, err := q.MarshalJSON() // Returns []byte("[1,2]"), nil
//
// Returns:
//   - The JSON-encoded byte slice.
//   - An error if marshaling fails.
//
// Time complexity: O(n), where n is the number of elements.
func (q *Deque[T]) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(q.Values())
	if err != nil {
		return nil, fmt.Errorf("ringbuf: %w: %w", ErrMarshalJSON, err)
	}

	return data, nil
}

// UnmarshalJSON populates the queue from a JSON array, appending elements to the back.
//
// The input must be a valid JSON array (e.g., "[1,2,3]"). The queue is cleared
// before loading. If the capacity is exceeded, older elements are overwritten.
//
// Example:
//
//	q := New[int](2)
//	err := q.UnmarshalJSON([]byte("[1,2,3]")) // Queue contains [2,3] after overflow
//
// Returns:
//
//	An error if the JSON is invalid or elements cannot be unmarshaled into type T.
//
// Time complexity: O(n), where n is the number of elements in the JSON array.
func (q *Deque[T]) UnmarshalJSON(data []byte) error {
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
