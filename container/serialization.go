// Package container provides interfaces for working with container data structures.
// It includes support for JSON serialization and deserialization, enabling containers
// to convert their elements to and from JSON format in a standardized way.
package container

// --------------------------------------------------------------------------------
// JSON Serialization Interface

// JSONSerializer defines methods for serializing container elements into JSON format.
//
// This interface provides a standard way to convert a container's data into a JSON byte
// slice. It includes a custom method ToJSON() and implements the json.Marshaler interface.
//
// Example usage:
//
//	type IntSlice []int
//	func (s IntSlice) ToJSON() ([]byte, error) {
//	    return json.Marshal(s)
//	}
//	func (s IntSlice) MarshalJSON() ([]byte, error) {
//	    return json.Marshal(s)
//	}
type JSONSerializer interface {
	// ToJSON returns the JSON representation of the container's elements as a byte slice.
	// It serializes the container's current state into a valid JSON format.
	// Returns an error if serialization fails, such as when elements are not JSON-compatible.
	ToJSON() ([]byte, error)

	// MarshalJSON implements the json.Marshaler interface, enabling the container to be
	// serialized using json.Marshal(). It returns the JSON representation of the container's
	// elements or an error if serialization fails.
	MarshalJSON() ([]byte, error)
}

// --------------------------------------------------------------------------------
// JSON Deserialization Interface

// JSONDeserializer defines methods for deserializing JSON data into container elements.
//
// This interface provides a standard way to populate a container's data from a JSON byte
// slice. It includes a custom method FromJSON() and implements the json.Unmarshaler interface.
//
// Example usage:
//
//	type IntSlice []int
//	func (s *IntSlice) FromJSON(data []byte) error {
//	    return json.Unmarshal(data, s)
//	}
//	func (s *IntSlice) UnmarshalJSON(data []byte) error {
//	    return json.Unmarshal(data, s)
//	}
type JSONDeserializer interface {
	// FromJSON populates the container's elements from the provided JSON byte slice.
	// It modifies the container's state to reflect the deserialized data.
	// Returns an error if the JSON is invalid or incompatible with the container's structure.
	FromJSON(data []byte) error

	// UnmarshalJSON implements the json.Unmarshaler interface, enabling the container to be
	// deserialized using json.Unmarshal(). It populates the container's elements from the
	// JSON data or returns an error if deserialization fails.
	UnmarshalJSON(data []byte) error
}

// --------------------------------------------------------------------------------
// Combined Interface (Optional)

// JSONSerializable combines JSONSerializer and JSONDeserializer into a single interface.
// It can be used for containers that support both serialization and deserialization.
//
// This is an optional convenience interface and may be implemented as needed.
type JSONSerializable interface {
	JSONSerializer
	JSONDeserializer
}
