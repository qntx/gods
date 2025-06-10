package btreebimap_test

import (
	"testing"

	"github.com/qntx/gods/btreebimap"
)

func benchmarkGet(b *testing.B, m *btreebimap.Map[int, int], size int) {
	for range b.N {
		for n := range size {
			m.Get(n)
		}
	}
}

func benchmarkPut(b *testing.B, m *btreebimap.Map[int, int], size int) {
	for range b.N {
		for n := range size {
			m.Put(n, n)
		}
	}
}

func benchmarkRemove(b *testing.B, m *btreebimap.Map[int, int], size int) {
	for range b.N {
		for n := range size {
			m.Remove(n)
		}
	}
}

func BenchmarkTreeBidiMapGet100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := btreebimap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkGet(b, m, size)
}

func BenchmarkTreeBidiMapGet1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := btreebimap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkGet(b, m, size)
}

func BenchmarkTreeBidiMapGet10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := btreebimap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkGet(b, m, size)
}

func BenchmarkTreeBidiMapGet100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := btreebimap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkGet(b, m, size)
}

func BenchmarkTreeBidiMapPut100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := btreebimap.New[int, int]()

	b.StartTimer()
	benchmarkPut(b, m, size)
}

func BenchmarkTreeBidiMapPut1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := btreebimap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkPut(b, m, size)
}

func BenchmarkTreeBidiMapPut10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := btreebimap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkPut(b, m, size)
}

func BenchmarkTreeBidiMapPut100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := btreebimap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkPut(b, m, size)
}

func BenchmarkTreeBidiMapRemove100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := btreebimap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkRemove(b, m, size)
}

func BenchmarkTreeBidiMapRemove1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := btreebimap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkRemove(b, m, size)
}

func BenchmarkTreeBidiMapRemove10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := btreebimap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkRemove(b, m, size)
}

func BenchmarkTreeBidiMapRemove100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := btreebimap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkRemove(b, m, size)
}
