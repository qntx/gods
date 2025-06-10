package container

import "iter"

// Map is a generic interface for key-value mappings, where keys are unique and
// map to a single value. Implementations (e.g., hash maps, tree maps) must
// support all operations defined here. This interface does not assume any
// specific ordering of keys, making it suitable for both ordered and unordered
// maps.
//
// Type parameter K must be comparable to ensure keys can be used in equality
// checks. Type parameter V represents the value type and has no constraints,
// allowing any type to be stored as a value.
//
// The Map interface extends Container[V], inheriting operations like Len,
// IsEmpty, and Clear, where the element type is the value type V.
type Map[K comparable, V any] interface {
	Container[V]

	// Put associates the specified value with the given key in the map.
	// If the key already exists, its value is updated with the new value.
	// If the key does not exist, a new key-value pair is added.
	Put(key K, value V)

	// Get retrieves the value associated with the specified key.
	// Returns the value and true if the key was found, or the zero value of V
	// and false if the key was not present.
	Get(key K) (V, bool)

	// Has returns true if the specified key is present in the map, false otherwise.
	Has(key K) bool

	// Delete removes the key-value pair associated with the specified key.
	// Returns true if the key was found and removed, false if the key was not present.
	Delete(key K) bool

	// Keys returns a slice containing all keys in the map.
	// The order of keys matches the map's ordering (e.g., sorted for tree maps).
	// Returns an empty slice if the map is empty.
	Keys() []K

	// Values returns a slice containing all values in the map.
	// The order of values corresponds to the order of keys in the map's ordering.
	// Returns an empty slice if the map is empty.
	Values() []V

	// Entries returns two slices containing all keys and values in the map,
	// respectively, in the map's ordering (e.g., sorted for tree maps).
	// The returned slices are of equal length, where the i-th element of the keys
	// slice corresponds to the i-th element of the values slice.
	// Returns empty slices if the map is empty.
	//
	// Example:
	//   keys, values := m.Entries()
	//   for i := 0; i < len(keys); i++ {
	//       fmt.Printf("Key: %v, Value: %v\n", keys[i], values[i])
	//   }
	Entries() (keys []K, values []V)

	// Clone returns a clone of the map using the same implementation,
	// duplicating all keys and values.
	Clone() Map[K, V]
}

// OrderedMap is a generic interface for key-value mappings that maintain a
// specific ordering of keys, such as sorted order in tree maps. It extends the
// Map interface with operations that depend on this ordering, such as accessing
// the first or last key-value pair, iterating in order, or retrieving keys and
// values in the map's defined order.
//
// Implementations must define a consistent ordering (e.g., sorted by key for
// tree maps). Type parameters K and V follow the same constraints as in Map.
type OrderedMap[K comparable, V any] interface {
	Map[K, V]

	// Begin returns the first key-value pair in the map's ordering.
	// For ordered implementations (e.g., tree maps), this is the smallest key
	// according to the map's comparison function.
	// Returns the key, value, and true if the map is non-empty, or zero values
	// and false if the map is empty.
	Begin() (K, V, bool)

	// End returns the last key-value pair in the map's ordering.
	// For ordered implementations, this is the largest key according to the
	// map's comparison function.
	// Returns the key, value, and true if the map is non-empty, or zero values
	// and false if the map is empty.
	End() (K, V, bool)

	// DeleteBegin removes and returns the first key-value pair in the map's ordering.
	// For ordered implementations, this removes the smallest key.
	// Returns the key, value, and true if the map was non-empty, or zero values
	// and false if the map was empty.
	DeleteBegin() (K, V, bool)

	// DeleteEnd removes and returns the last key-value pair in the map's ordering.
	// For ordered implementations, this removes the largest key.
	// Returns the key, value, and true if the map was non-empty, or zero values
	// and false if the map was empty.
	DeleteEnd() (K, V, bool)

	// Iter returns an iterator over the key-value pairs in the map.
	// The iterator yields pairs in the map's ordering (e.g., sorted for tree maps).
	// Suitable for range loops.
	//
	// Example:
	//   for key, value := range m.Iter() {
	//       fmt.Printf("Key: %v, Value: %v\n", key, value)
	//   }
	//
	// The iterator is safe for concurrent reads but not for concurrent writes
	// unless the implementation explicitly supports it.
	Iter() iter.Seq2[K, V]

	// RIter returns an iterator over the key-value pairs in the map.
	// The iterator yields pairs in the map's reverse ordering.
	// Suitable for range loops.
	//
	// Example:
	//   for key, value := range m.RIter() {
	//       fmt.Printf("Key: %v, Value: %v\n", key, value)
	//   }
	//
	// The iterator is safe for concurrent reads but not for concurrent writes
	// unless the implementation explicitly supports it.
	RIter() iter.Seq2[K, V]
}

// BidiMap is a generic interface for bidirectional maps, extending the Map
// interface to support lookups by both keys and values. In a bidirectional map,
// both keys and values are unique, allowing values to map back to their
// corresponding keys.
//
// Type parameter K must be comparable for key equality checks.
// Type parameter V must also be comparable to ensure value uniqueness and
// lookup capability.
//
// BidiMap inherits all Map operations and adds methods specific to
// bidirectional functionality. Implementations must maintain the invariant that
// each value is associated with exactly one key, just as each key is associated
// with exactly one value.
type BidiMap[K comparable, V comparable] interface {
	Map[K, V]

	// GetKey retrieves the key associated with the specified value.
	// Returns the key and true if the value was found, or the zero value of K
	// and false if the value was not present.
	GetKey(value V) (K, bool)

	// DeleteValue removes the key-value pair associated with the specified value.
	// Returns true if the value was found and removed, false if the value was not present.
	DeleteValue(value V) bool
}
