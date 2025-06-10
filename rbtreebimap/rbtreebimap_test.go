package rbtreebimap_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/qntx/gods/rbtreebimap"
)

func SameElements[T comparable](t *testing.T, actual, expected []T) {
	if len(actual) != len(expected) {
		t.Errorf("Got %d expected %d", len(actual), len(expected))
	}
outer:
	for _, e := range expected {
		for _, a := range actual {
			if e == a {
				continue outer
			}
		}
		t.Errorf("Did not find expected element %v in %v", e, actual)
	}
}

func TestMapPut(t *testing.T) {
	m := rbtreebimap.New[int, string]()
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

	SameElements(t, m.Keys(), []int{1, 2, 3, 4, 5, 6, 7})
	SameElements(t, m.Values(), []string{"a", "b", "c", "d", "e", "f", "g"})

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

func TestMapRemove(t *testing.T) {
	m := rbtreebimap.New[int, string]()
	m.Put(5, "e")
	m.Put(6, "f")
	m.Put(7, "g")
	m.Put(3, "c")
	m.Put(4, "d")
	m.Put(1, "x")
	m.Put(2, "b")
	m.Put(1, "a") //overwrite

	m.Remove(5)
	m.Remove(6)
	m.Remove(7)
	m.Remove(8)
	m.Remove(5)

	SameElements(t, m.Keys(), []int{1, 2, 3, 4})
	SameElements(t, m.Values(), []string{"a", "b", "c", "d"})

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

	m.Remove(1)
	m.Remove(4)
	m.Remove(2)
	m.Remove(3)
	m.Remove(2)
	m.Remove(2)

	SameElements(t, m.Keys(), nil)
	SameElements(t, m.Values(), nil)

	if actualValue := m.Len(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}

	if actualValue := m.Empty(); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
}

func TestMapGetKey(t *testing.T) {
	m := rbtreebimap.New[int, string]()
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

func TestMapEach(t *testing.T) {
	m := rbtreebimap.New[string, int]()
	m.Put("c", 3)
	m.Put("a", 1)
	m.Put("b", 2)

	count := 0

	m.Each(func(key string, value int) {
		count++
		if actualValue, expectedValue := count, value; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}

		switch value {
		case 1:
			if actualValue, expectedValue := key, "a"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 2:
			if actualValue, expectedValue := key, "b"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 3:
			if actualValue, expectedValue := key, "c"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			t.Errorf("Too many")
		}
	})
}

func TestMapMap(t *testing.T) {
	m := rbtreebimap.New[string, int]()
	m.Put("c", 3)
	m.Put("a", 1)
	m.Put("b", 2)

	mappedMap := m.Map(func(key1 string, value1 int) (key2 string, value2 int) {
		return key1, value1 * value1
	})
	if actualValue, _ := mappedMap.Get("a"); actualValue != 1 {
		t.Errorf("Got %v expected %v", actualValue, "mapped: a")
	}

	if actualValue, _ := mappedMap.Get("b"); actualValue != 4 {
		t.Errorf("Got %v expected %v", actualValue, "mapped: b")
	}

	if actualValue, _ := mappedMap.Get("c"); actualValue != 9 {
		t.Errorf("Got %v expected %v", actualValue, "mapped: c")
	}

	if mappedMap.Len() != 3 {
		t.Errorf("Got %v expected %v", mappedMap.Len(), 3)
	}
}

func TestMapSelect(t *testing.T) {
	m := rbtreebimap.New[string, int]()
	m.Put("c", 3)
	m.Put("a", 1)
	m.Put("b", 2)

	selectedMap := m.Select(func(key string, value int) bool {
		return key >= "a" && key <= "b"
	})
	if actualValue, _ := selectedMap.Get("a"); actualValue != 1 {
		t.Errorf("Got %v expected %v", actualValue, "value: a")
	}

	if actualValue, _ := selectedMap.Get("b"); actualValue != 2 {
		t.Errorf("Got %v expected %v", actualValue, "value: b")
	}

	if selectedMap.Len() != 2 {
		t.Errorf("Got %v expected %v", selectedMap.Len(), 2)
	}
}

func TestMapAny(t *testing.T) {
	m := rbtreebimap.New[string, int]()
	m.Put("c", 3)
	m.Put("a", 1)
	m.Put("b", 2)

	any := m.Any(func(key string, value int) bool {
		return value == 3
	})
	if any != true {
		t.Errorf("Got %v expected %v", any, true)
	}

	any = m.Any(func(key string, value int) bool {
		return value == 4
	})
	if any != false {
		t.Errorf("Got %v expected %v", any, false)
	}
}

func TestMapAll(t *testing.T) {
	m := rbtreebimap.New[string, int]()
	m.Put("c", 3)
	m.Put("a", 1)
	m.Put("b", 2)

	all := m.All(func(key string, value int) bool {
		return key >= "a" && key <= "c"
	})
	if all != true {
		t.Errorf("Got %v expected %v", all, true)
	}

	all = m.All(func(key string, value int) bool {
		return key >= "a" && key <= "b"
	})
	if all != false {
		t.Errorf("Got %v expected %v", all, false)
	}
}

func TestMapFind(t *testing.T) {
	m := rbtreebimap.New[string, int]()
	m.Put("c", 3)
	m.Put("a", 1)
	m.Put("b", 2)

	foundKey, foundValue := m.Find(func(key string, value int) bool {
		return key == "c"
	})
	if foundKey != "c" || foundValue != 3 {
		t.Errorf("Got %v -> %v expected %v -> %v", foundKey, foundValue, "c", 3)
	}

	foundKey, foundValue = m.Find(func(key string, value int) bool {
		return key == "x"
	})
	if foundKey != "" || foundValue != 0 {
		t.Errorf("Got %v at %v expected %v at %v", foundValue, foundKey, nil, nil)
	}
}

func TestMapChaining(t *testing.T) {
	m := rbtreebimap.New[string, int]()
	m.Put("c", 3)
	m.Put("a", 1)
	m.Put("b", 2)

	chainedMap := m.Select(func(key string, value int) bool {
		return value > 1
	}).Map(func(key string, value int) (string, int) {
		return key + key, value * value
	})
	if actualValue := chainedMap.Len(); actualValue != 2 {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}

	if actualValue, found := chainedMap.Get("aa"); actualValue != 0 || found {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}

	if actualValue, found := chainedMap.Get("bb"); actualValue != 4 || !found {
		t.Errorf("Got %v expected %v", actualValue, 4)
	}

	if actualValue, found := chainedMap.Get("cc"); actualValue != 9 || !found {
		t.Errorf("Got %v expected %v", actualValue, 9)
	}
}

func TestMapSerialization(t *testing.T) {
	for range 10 {
		original := rbtreebimap.New[string, string]()
		original.Put("d", "4")
		original.Put("e", "5")
		original.Put("c", "3")
		original.Put("b", "2")
		original.Put("a", "1")

		serialized, err := original.MarshalJSON()
		if err != nil {
			t.Errorf("Got error %v", err)
		}

		deserialized := rbtreebimap.New[string, string]()

		err = deserialized.UnmarshalJSON(serialized)
		if err != nil {
			t.Errorf("Got error %v", err)
		}

		if original.Len() != deserialized.Len() {
			t.Errorf("Got map of size %d, expected %d", original.Len(), deserialized.Len())
		}

		original.Each(func(key string, expected string) {
			actual, ok := deserialized.Get(key)
			if !ok || actual != expected {
				t.Errorf("Did not find expected value %v for key %v in deserialied map (got %q)", expected, key, actual)
			}
		})
	}

	m := rbtreebimap.New[string, float64]()
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
	c := rbtreebimap.New[string, string]()
	c.Put("a", "a")

	if !strings.HasPrefix(c.String(), "TreeBidiMap") {
		t.Errorf("String should start with container name")
	}
}
