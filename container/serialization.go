// Package container provides interfaces for managing container data structures.
// It supports JSON serialization and deserialization, allowing containers to
// convert their elements to and from JSON in a standardized manner.
package container

import "encoding/json"

// JSONCodec defines an interface for containers that support both JSON
// serialization and deserialization. It combines the Marshaler and Unmarshaler
// interfaces for convenience.
//
// This interface is optional and may be implemented as needed.
type JSONCodec interface {
	json.Marshaler
	json.Unmarshaler
}
