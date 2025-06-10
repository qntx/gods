package hashbidimap

import (
	"encoding/json"

	"github.com/qntx/gods/container"
)

// Assert Serialization implementation.
var _ container.JSONCodec = (*Map[string, int])(nil)

// MarshalJSON outputs the JSON representation of the map.
func (m *Map[K, V]) MarshalJSON() ([]byte, error) {
	return m.forwardMap.MarshalJSON()
}

// UnmarshalJSON populates the map from the input JSON representation.
func (m *Map[K, V]) UnmarshalJSON(data []byte) error {
	var elements map[K]V

	err := json.Unmarshal(data, &elements)
	if err != nil {
		return err
	}

	m.Clear()

	for k, v := range elements {
		m.Put(k, v)
	}

	return nil
}
