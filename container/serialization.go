// Package container provides interfaces for working with container data structures.
// It includes support for JSON serialization and deserialization, enabling containers
// to convert their elements to and from JSON format in a standardized way.
package container

import "encoding/json"

// JSONSerializable combines JSONSerializer and JSONDeserializer into a single interface.
// It can be used for containers that support both serialization and deserialization.
//
// This is an optional convenience interface and may be implemented as needed.
type JSONSerializable interface {
	json.Marshaler
	json.Unmarshaler
}
