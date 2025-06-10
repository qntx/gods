package linkedhashset

import (
	"encoding/json"

	"github.com/qntx/gods/container"
)

// Assert Serialization implementation.
var _ container.JSONCodec = (*Set[int])(nil)

// MarshalJSON outputs the JSON representation of the set.
func (set *Set[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(set.Values())
}

// UnmarshalJSON populates the set from the input JSON representation.
func (set *Set[T]) UnmarshalJSON(data []byte) error {
	var elements []T

	err := json.Unmarshal(data, &elements)
	if err == nil {
		set.Clear()
		set.Add(elements...)
	}

	return err
}
