package hashmap

import (
	"encoding/json"

	"github.com/qntx/gods/container"
)

// Assert Serialization implementation.
var _ container.JSONCodec = (*Map[string, int])(nil)

// MarshalJSON outputs the JSON representation of the map.
func (m *Map[K, V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.m)
}

// UnmarshalJSON populates the map from the input JSON representation.
func (m *Map[K, V]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &m.m)
}
