package hashbidimap

import (
	"encoding/json"
	"strings"
	"testing"
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
	m := New[int, string]()
	m.Put(5, "e")
	m.Put(6, "f")
	m.Put(7, "g")
	m.Put(3, "c")
	m.Put(4, "d")
	m.Put(1, "x")
	m.Put(2, "b")
	m.Put(1, "a") //overwrite

	if actualValue := m.Size(); actualValue != 7 {
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
	m := New[int, string]()
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

	if actualValue := m.Size(); actualValue != 4 {
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

	if actualValue := m.Size(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}

	if actualValue := m.Empty(); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
}

func TestMapGetKey(t *testing.T) {
	m := New[int, string]()
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
	m := New[string, float64]()
	m.Put("a", 1.0)
	m.Put("b", 2.0)
	m.Put("c", 3.0)

	var err error

	assert := func() {
		SameElements(t, m.Keys(), []string{"a", "b", "c"})
		SameElements(t, m.Values(), []float64{1.0, 2.0, 3.0})

		if actualValue, expectedValue := m.Size(), 3; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}

		if err != nil {
			t.Errorf("Got error %v", err)
		}
	}

	assert()

	bytes, err := m.MarshalJSON()

	assert()

	err = m.UnmarshalJSON(bytes)

	assert()

	_, err = json.Marshal([]any{"a", "b", "c", m})
	if err != nil {
		t.Errorf("Got error %v", err)
	}

	err = json.Unmarshal([]byte(`{"a":1,"b":2}`), &m)
	if err != nil {
		t.Errorf("Got error %v", err)
	}
}

func TestMapString(t *testing.T) {
	c := New[string, int]()
	c.Put("a", 1)

	if !strings.HasPrefix(c.String(), "HashBidiMap") {
		t.Errorf("String should start with container name")
	}
}

func benchmarkGet(b *testing.B, m *Map[int, int], size int) {
	for range b.N {
		for n := range size {
			m.Get(n)
		}
	}
}

func benchmarkPut(b *testing.B, m *Map[int, int], size int) {
	for range b.N {
		for n := range size {
			m.Put(n, n)
		}
	}
}

func benchmarkRemove(b *testing.B, m *Map[int, int], size int) {
	for range b.N {
		for n := range size {
			m.Remove(n)
		}
	}
}

func BenchmarkHashBidiMapGet100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkGet(b, m, size)
}

func BenchmarkHashBidiMapGet1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkGet(b, m, size)
}

func BenchmarkHashBidiMapGet10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkGet(b, m, size)
}

func BenchmarkHashBidiMapGet100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkGet(b, m, size)
}

func BenchmarkHashBidiMapPut100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := New[int, int]()

	b.StartTimer()
	benchmarkPut(b, m, size)
}

func BenchmarkHashBidiMapPut1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkPut(b, m, size)
}

func BenchmarkHashBidiMapPut10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkPut(b, m, size)
}

func BenchmarkHashBidiMapPut100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkPut(b, m, size)
}

func BenchmarkHashBidiMapRemove100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkRemove(b, m, size)
}

func BenchmarkHashBidiMapRemove1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkRemove(b, m, size)
}

func BenchmarkHashBidiMapRemove10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkRemove(b, m, size)
}

func BenchmarkHashBidiMapRemove100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkRemove(b, m, size)
}
