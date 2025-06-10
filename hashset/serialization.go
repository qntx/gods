package hashset

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/qntx/gods/container"
)

var _ container.JSONCodec = (*Set[int])(nil)

// MarshalJSON creates a JSON array from the set, it marshals all elements.
func (s Set[T]) MarshalJSON() ([]byte, error) {
	items := make([]string, 0, s.Len())

	for elem := range s {
		b, err := json.Marshal(elem)
		if err != nil {
			return nil, err
		}

		items = append(items, string(b))
	}

	return []byte(fmt.Sprintf("[%s]", strings.Join(items, ","))), nil
}

// UnmarshalJSON recreates a set from a JSON array, it only decodes
// primitive types. Numbers are decoded as json.Number.
func (s *Set[T]) UnmarshalJSON(b []byte) error {
	var i []T

	err := json.Unmarshal(b, &i)
	if err != nil {
		return err
	}

	s.Append(i...)

	return nil
}
