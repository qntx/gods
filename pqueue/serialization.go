// Package pqueue provides JSON serialization and deserialization for priority queues.
//
// This file extends the PriorityQueue type with methods to convert to and from JSON format,
// implementing the container.JSONSerializer and container.JSONDeserializer interfaces.
package pqueue

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/qntx/gods/container"
)

// --------------------------------------------------------------------------------
// Constants and Errors
// --------------------------------------------------------------------------------

// Predefined errors for JSON operations.
var (
	ErrMarshalJSONFailure   = errors.New("failed to marshal priority queue to JSON")
	ErrUnmarshalJSONFailure = errors.New("failed to unmarshal JSON into priority queue")
)

// --------------------------------------------------------------------------------
// Interface Assertions
// --------------------------------------------------------------------------------

// Ensure PriorityQueue implements required interfaces at compile time.
var (
	_ container.JSONSerializer   = (*PriorityQueue[int, int])(nil)
	_ container.JSONDeserializer = (*PriorityQueue[int, int])(nil)
	_ json.Marshaler             = (*PriorityQueue[int, int])(nil)
	_ json.Unmarshaler           = (*PriorityQueue[int, int])(nil)
)

// --------------------------------------------------------------------------------
// JSON Serialization Helper Types
// --------------------------------------------------------------------------------

// pqJSON is a helper struct for serializing the priority queue to JSON.
type pqJSON[T comparable, V any] struct {
	Items []itemJSON[T, V] `json:"items"`
	Kind  HeapKind         `json:"kind"`
}

// itemJSON represents a serialized priority queue item.
type itemJSON[T comparable, V any] struct {
	Value    T `json:"value"`
	Priority V `json:"priority"`
}

// --------------------------------------------------------------------------------
// JSON Serialization Methods
// --------------------------------------------------------------------------------

// ToJSON serializes the priority queue into a JSON object.
//
// Converts the queue's elements into a JSON structure that preserves both values
// and priorities. Returns the JSON-encoded byte slice or an error if marshaling fails.
//
// Example:
//
//	pq := New[string, int](MinHeap)
//	pq.Put("task1", 5)
//	data, err := pq.ToJSON()
//
// Time complexity: O(n), where n is the number of items in the queue.
func (pq *PriorityQueue[T, V]) ToJSON() ([]byte, error) {
	items := make([]itemJSON[T, V], 0, len(pq.heap))
	for _, item := range pq.heap {
		items = append(items, itemJSON[T, V]{
			Value:    item.Value,
			Priority: item.Priority,
		})
	}

	data, err := json.Marshal(pqJSON[T, V]{
		Items: items,
		Kind:  pq.kind,
	})
	if err != nil {
		return nil, fmt.Errorf("pqueue: %w: %w", ErrMarshalJSONFailure, err)
	}

	return data, nil
}

// FromJSON populates the priority queue from JSON data.
//
// Expects a JSON structure with items array and kind. Clears the queue before loading
// and inserts each value-priority pair. Returns an error if the JSON is invalid or
// unmarshaling fails.
//
// Example:
//
//	pq := New[string, int](MinHeap)
//	err := pq.FromJSON([]byte(`{"items":[{"value":"task1","priority":5}],"kind":0}`))
//
// Time complexity: O(n log n), where n is the number of items in the JSON.
func (pq *PriorityQueue[T, V]) FromJSON(data []byte) error {
	var jsonQueue pqJSON[T, V]
	if err := json.Unmarshal(data, &jsonQueue); err != nil {
		return fmt.Errorf("pqueue: %w: %w", ErrUnmarshalJSONFailure, err)
	}

	pq.Clear() // Clear existing data
	pq.kind = jsonQueue.Kind

	// Add items one by one
	for _, item := range jsonQueue.Items {
		pq.Put(item.Value, item.Priority)
	}
	return nil
}

// MarshalJSON implements json.Marshaler for seamless JSON encoding.
//
// Delegates to ToJSON() for consistency. Returns the JSON byte slice or an error
// if serialization fails.
//
// Time complexity: O(n), where n is the number of items in the queue.
func (pq *PriorityQueue[T, V]) MarshalJSON() ([]byte, error) {
	return pq.ToJSON()
}

// UnmarshalJSON implements json.Unmarshaler for seamless JSON decoding.
//
// Delegates to FromJSON() to populate the queue. Returns an error if deserialization
// fails.
//
// Time complexity: O(n log n), where n is the number of items in the JSON.
func (pq *PriorityQueue[T, V]) UnmarshalJSON(data []byte) error {
	return pq.FromJSON(data)
}
