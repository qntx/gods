package linkedhashmap_test

import (
	"testing"

	"github.com/qntx/gods/linkedhashmap"
)

func benchmarkGet(b *testing.B, m *linkedhashmap.Map[int, int], size int) {
	for range b.N {
		for n := range size {
			m.Get(n)
		}
	}
}

func benchmarkPut(b *testing.B, m *linkedhashmap.Map[int, int], size int) {
	for range b.N {
		for n := range size {
			m.Put(n, n)
		}
	}
}

func benchmarkRemove(b *testing.B, m *linkedhashmap.Map[int, int], size int) {
	for range b.N {
		for n := range size {
			m.Remove(n)
		}
	}
}

func BenchmarkTreeMapGet100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := linkedhashmap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkGet(b, m, size)
}

func BenchmarkTreeMapGet1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := linkedhashmap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkGet(b, m, size)
}

func BenchmarkTreeMapGet10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := linkedhashmap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkGet(b, m, size)
}

func BenchmarkTreeMapGet100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := linkedhashmap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkGet(b, m, size)
}

func BenchmarkTreeMapPut100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := linkedhashmap.New[int, int]()

	b.StartTimer()
	benchmarkPut(b, m, size)
}

func BenchmarkTreeMapPut1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := linkedhashmap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkPut(b, m, size)
}

func BenchmarkTreeMapPut10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := linkedhashmap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkPut(b, m, size)
}

func BenchmarkTreeMapPut100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := linkedhashmap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkPut(b, m, size)
}

func BenchmarkTreeMapRemove100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := linkedhashmap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkRemove(b, m, size)
}

func BenchmarkTreeMapRemove1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := linkedhashmap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkRemove(b, m, size)
}

func BenchmarkTreeMapRemove10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := linkedhashmap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkRemove(b, m, size)
}

func BenchmarkTreeMapRemove100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := linkedhashmap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkRemove(b, m, size)
}
