package hashmap_test

import (
	"testing"

	"github.com/qntx/gods/hashmap"
)

func benchmarkGet(b *testing.B, m *hashmap.Map[int, int], size int) {
	for range b.N {
		for n := range size {
			m.Get(n)
		}
	}
}

func benchmarkPut(b *testing.B, m *hashmap.Map[int, int], size int) {
	for range b.N {
		for n := range size {
			m.Put(n, n)
		}
	}
}

func benchmarkRemove(b *testing.B, m *hashmap.Map[int, int], size int) {
	for range b.N {
		for n := range size {
			m.Delete(n)
		}
	}
}

func BenchmarkHashMapGet100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := hashmap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkGet(b, m, size)
}

func BenchmarkHashMapGet1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := hashmap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkGet(b, m, size)
}

func BenchmarkHashMapGet10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := hashmap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkGet(b, m, size)
}

func BenchmarkHashMapGet100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := hashmap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkGet(b, m, size)
}

func BenchmarkHashMapPut100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := hashmap.New[int, int]()

	b.StartTimer()
	benchmarkPut(b, m, size)
}

func BenchmarkHashMapPut1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := hashmap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkPut(b, m, size)
}

func BenchmarkHashMapPut10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := hashmap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkPut(b, m, size)
}

func BenchmarkHashMapPut100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := hashmap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkPut(b, m, size)
}

func BenchmarkHashMapRemove100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := hashmap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkRemove(b, m, size)
}

func BenchmarkHashMapRemove1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := hashmap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkRemove(b, m, size)
}

func BenchmarkHashMapRemove10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := hashmap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkRemove(b, m, size)
}

func BenchmarkHashMapRemove100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := hashmap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkRemove(b, m, size)
}
