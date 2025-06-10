package rbtreebidimap

import (
	"encoding/json"

	"github.com/qntx/gods/container"
)

// Assert Serialization implementation.
var _ container.JSONCodec = (*Map[string, int])(nil)

// MarshalJSON @implements json.Marshaler.
func (m *Map[K, V]) MarshalJSON() ([]byte, error) {
	return m.forwardMap.MarshalJSON()
}

// UnmarshalJSON @implements json.Unmarshaler.
func (m *Map[K, V]) UnmarshalJSON(data []byte) error {
	var elements map[K]V

	err := json.Unmarshal(data, &elements)
	if err != nil {
		return err
	}

	m.Clear()

	for key, value := range elements {
		m.Put(key, value)
	}

	return nil
}
