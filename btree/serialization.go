// Package btree provides JSON serialization for a generic B-tree data structure.
package btree

// import (
// 	"encoding/json"
// 	"errors"
// 	"fmt"

// 	"github.com/qntx/gods/container"
// )

// // --------------------------------------------------------------------------------
// // Interface Assertions

// // Verify Tree implements required interfaces at compile time.
// var (
// 	_ container.JSONSerializer   = (*Tree[string, int])(nil)
// 	_ container.JSONDeserializer = (*Tree[string, int])(nil)
// 	_ json.Marshaler             = (*Tree[string, int])(nil)
// 	_ json.Unmarshaler           = (*Tree[string, int])(nil)
// )

// // --------------------------------------------------------------------------------
// // JSON Serialization Methods

// // ToJSON serializes the tree into a JSON object mapping keys to values.
// //
// // The output is a JSON object where each key-value pair corresponds to an entry
// // in the tree (e.g., {"a": 1, "b": 2}). All elements must be JSON-serializable.
// //
// // Returns:
// //   - The JSON-encoded byte slice.
// //   - An error if marshaling fails.
// func (t *Tree[K, V]) ToJSON() ([]byte, error) {
// 	elements := make(map[K]V)
// 	it := t.Iterator()

// 	for it.Next() {
// 		elements[it.Key()] = it.Value()
// 	}

// 	return json.Marshal(elements)
// }

// // FromJSON populates the tree from a JSON object.
// //
// // The input must be a valid JSON object (e.g., {"a": 1, "b": 2}). The tree is
// // cleared before loading the new data.
// //
// // Returns:
// //
// //	An error if the JSON is invalid or unmarshaling fails.
// func (t *Tree[K, V]) FromJSON(data []byte) error {
// 	if len(data) == 0 {
// 		return errors.New("btree: empty JSON data")
// 	}

// 	elements := make(map[K]V)
// 	if err := json.Unmarshal(data, &elements); err != nil {
// 		return fmt.Errorf("btree: invalid JSON: %w", err)
// 	}

// 	t.Clear()

// 	for k, v := range elements {
// 		t.Put(k, v)
// 	}

// 	return nil
// }

// // MarshalJSON implements json.Marshaler for JSON encoding.
// //
// // Returns:
// //   - The JSON-encoded byte slice.
// //   - An error if marshaling fails.
// func (t *Tree[K, V]) MarshalJSON() ([]byte, error) {
// 	return t.ToJSON()
// }

// // UnmarshalJSON implements json.Unmarshaler for JSON decoding.
// //
// // Returns:
// //
// //	An error if deserialization fails.
// func (t *Tree[K, V]) UnmarshalJSON(data []byte) error {
// 	return t.FromJSON(data)
// }
