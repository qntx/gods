package btreebimap_test

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"

	"github.com/qntx/gods/btreebimap"
)

func TestMapPut(t *testing.T) {
	m := btreebimap.New[int, string]()
	m.Put(5, "e")
	m.Put(6, "f")
	m.Put(7, "g")
	m.Put(3, "c")
	m.Put(4, "d")
	m.Put(1, "x")
	m.Put(2, "b")
	m.Put(1, "a") //overwrite

	if actualValue := m.Len(); actualValue != 7 {
		t.Errorf("Got %v expected %v", actualValue, 7)
	}

	slices.Equal(m.Keys(), []int{1, 2, 3, 4, 5, 6, 7})
	slices.Equal(m.Values(), []string{"a", "b", "c", "d", "e", "f", "g"})

	// key,expectedValue,expectedFound
	tests1 := [][]interface{}{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "e", true},
		{6, "f", true},
		{7, "g", true},
		{8, "", false},
	}

	for _, test := range tests1 {
		// retrievals
		actualValue, actualFound := m.Get(test[0].(int))
		if actualValue != test[1] || actualFound != test[2] {
			t.Errorf("Got %v expected %v", actualValue, test[1])
		}
	}
}

func TestMapDelete(t *testing.T) {
	m := btreebimap.New[int, string]()
	m.Put(5, "e")
	m.Put(6, "f")
	m.Put(7, "g")
	m.Put(3, "c")
	m.Put(4, "d")
	m.Put(1, "x")
	m.Put(2, "b")
	m.Put(1, "a") //overwrite

	m.Delete(5)
	m.Delete(6)
	m.Delete(7)
	m.Delete(8)
	m.Delete(5)

	slices.Equal(m.Keys(), []int{1, 2, 3, 4})
	slices.Equal(m.Values(), []string{"a", "b", "c", "d"})

	if actualValue := m.Len(); actualValue != 4 {
		t.Errorf("Got %v expected %v", actualValue, 4)
	}

	tests2 := [][]interface{}{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "", false},
		{6, "", false},
		{7, "", false},
		{8, "", false},
	}

	for _, test := range tests2 {
		actualValue, actualFound := m.Get(test[0].(int))
		if actualValue != test[1] || actualFound != test[2] {
			t.Errorf("Got %v expected %v", actualValue, test[1])
		}
	}

	m.Delete(1)
	m.Delete(4)
	m.Delete(2)
	m.Delete(3)
	m.Delete(2)
	m.Delete(2)

	slices.Equal(m.Keys(), nil)
	slices.Equal(m.Values(), nil)

	if actualValue := m.Len(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}

	if actualValue := m.IsEmpty(); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
}

func TestMapGetKey(t *testing.T) {
	m := btreebimap.New[int, string]()
	m.Put(5, "e")
	m.Put(6, "f")
	m.Put(7, "g")
	m.Put(3, "c")
	m.Put(4, "d")
	m.Put(1, "x")
	m.Put(2, "b")
	m.Put(1, "a") //overwrite

	// key,expectedValue,expectedFound
	tests1 := [][]interface{}{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "e", true},
		{6, "f", true},
		{7, "g", true},
		{0, "x", false},
	}

	for _, test := range tests1 {
		// retrievals
		actualValue, actualFound := m.GetKey(test[1].(string))
		if actualValue != test[0] || actualFound != test[2] {
			t.Errorf("Got %v expected %v", actualValue, test[0])
		}
	}
}
func TestMapSerialization(t *testing.T) {
	for range 10 {
		original := btreebimap.New[string, string]()
		original.Put("d", "4")
		original.Put("e", "5")
		original.Put("c", "3")
		original.Put("b", "2")
		original.Put("a", "1")

		serialized, err := original.MarshalJSON()
		if err != nil {
			t.Errorf("Got error %v", err)
		}

		deserialized := btreebimap.New[string, string]()

		err = deserialized.UnmarshalJSON(serialized)
		if err != nil {
			t.Errorf("Got error %v", err)
		}

		if original.Len() != deserialized.Len() {
			t.Errorf("Got map of size %d, expected %d", original.Len(), deserialized.Len())
		}
	}

	m := btreebimap.New[string, float64]()
	m.Put("a", 1.0)
	m.Put("b", 2.0)
	m.Put("c", 3.0)

	_, err := json.Marshal([]interface{}{"a", "b", "c", m})
	if err != nil {
		t.Errorf("Got error %v", err)
	}

	err = json.Unmarshal([]byte(`{"a":1,"b":2}`), &m)
	if err != nil {
		t.Errorf("Got error %v", err)
	}
}

func TestMapString(t *testing.T) {
	c := btreebimap.New[string, string]()
	c.Put("a", "a")

	if !strings.HasPrefix(c.String(), "TreeBidiMap") {
		t.Errorf("String should start with container name")
	}
}
