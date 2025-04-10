// Package rbtree provides JSON serialization and deserialization for the red-black tree.
//
// This file extends the Tree type with methods to convert to and from JSON format,
// implementing the container.JSONSerializer and container.JSONDeserializer interfaces.
package rbtree

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
	ErrMarshalJSONFailure   = errors.New("failed to marshal tree to JSON")
	ErrUnmarshalJSONFailure = errors.New("failed to unmarshal JSON into tree")
)

// --------------------------------------------------------------------------------
// Interface Assertions

// Ensure Tree implements required interfaces at compile time.
var (
	_ container.JSONSerializer   = (*Tree[string, int])(nil)
	_ container.JSONDeserializer = (*Tree[string, int])(nil)
	_ json.Marshaler             = (*Tree[string, int])(nil)
	_ json.Unmarshaler           = (*Tree[string, int])(nil)
)

// --------------------------------------------------------------------------------
// JSON Serialization Methods

// ToJSON serializes the tree into a JSON object.
//
// Converts the tree's key-value pairs into a JSON object where keys are the tree's
// keys and values are their corresponding values. Returns the JSON-encoded byte
// slice or an error if marshaling fails.
//
// Example:
//
//	t := New[string, int]()
//	t.Put("a", 1)
//	data, err := t.ToJSON() // Returns []byte(`{"a":1}`), nil
//
// Time complexity: O(n), where n is the number of nodes in the tree.
func (t *Tree[K, V]) ToJSON() ([]byte, error) {
	elems := make(map[K]V, t.Len())
	it := t.Iterator()

	for it.Next() {
		elems[it.Key()] = it.Value()
	}

	data, err := json.Marshal(elems)
	if err != nil {
		return nil, fmt.Errorf("rbtree: %w: %w", ErrMarshalJSONFailure, err)
	}

	return data, nil
}

// FromJSON populates the tree from a JSON object.
//
// Expects a JSON object (e.g., `{"a":1, "b":2}`). Clears the tree before loading
// and inserts each key-value pair. Returns an error if the JSON is invalid or
// unmarshaling fails.
//
// Example:
//
//	t := New[string, int]()
//	err := t.FromJSON([]byte(`{"a":1, "b":2}`)) // Tree contains {a:1, b:2}
//
// Time complexity: O(n log n), where n is the number of key-value pairs in the JSON.
func (t *Tree[K, V]) FromJSON(data []byte) error {
	var elems map[K]V
	if err := json.Unmarshal(data, &elems); err != nil {
		return fmt.Errorf("rbtree: %w: %w", ErrUnmarshalJSONFailure, err)
	}

	t.Clear()

	for k, v := range elems {
		t.Put(k, v)
	}

	return nil
}

// MarshalJSON implements json.Marshaler for seamless JSON encoding.
//
// Delegates to ToJSON() for consistency. Returns the JSON byte slice or an error
// if serialization fails.
//
// Time complexity: O(n), where n is the number of nodes in the tree.
func (t *Tree[K, V]) MarshalJSON() ([]byte, error) {
	return t.ToJSON()
}

// UnmarshalJSON implements json.Unmarshaler for seamless JSON decoding.
//
// Delegates to FromJSON() to populate the tree. Returns an error if deserialization
// fails.
//
// Time complexity: O(n log n), where n is the number of key-value pairs in the JSON.
func (t *Tree[K, V]) UnmarshalJSON(data []byte) error {
	return t.FromJSON(data)
}
