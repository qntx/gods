package hashbimap_test

import (
	"testing"

	"github.com/qntx/gods/hashbimap"
)

func benchmarkGet(b *testing.B, m *hashbimap.Map[int, int], size int) {
	for range b.N {
		for n := range size {
			m.Get(n)
		}
	}
}

func benchmarkPut(b *testing.B, m *hashbimap.Map[int, int], size int) {
	for range b.N {
		for n := range size {
			m.Put(n, n)
		}
	}
}

func benchmarkRemove(b *testing.B, m *hashbimap.Map[int, int], size int) {
	for range b.N {
		for n := range size {
			m.Remove(n)
		}
	}
}

func BenchmarkHashBidiMapGet100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := hashbimap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkGet(b, m, size)
}

func BenchmarkHashBidiMapGet1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := hashbimap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkGet(b, m, size)
}

func BenchmarkHashBidiMapGet10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := hashbimap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkGet(b, m, size)
}

func BenchmarkHashBidiMapGet100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := hashbimap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkGet(b, m, size)
}

func BenchmarkHashBidiMapPut100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := hashbimap.New[int, int]()

	b.StartTimer()
	benchmarkPut(b, m, size)
}

func BenchmarkHashBidiMapPut1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := hashbimap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkPut(b, m, size)
}

func BenchmarkHashBidiMapPut10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := hashbimap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkPut(b, m, size)
}

func BenchmarkHashBidiMapPut100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := hashbimap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkPut(b, m, size)
}

func BenchmarkHashBidiMapRemove100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := hashbimap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkRemove(b, m, size)
}

func BenchmarkHashBidiMapRemove1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := hashbimap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkRemove(b, m, size)
}

func BenchmarkHashBidiMapRemove10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := hashbimap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkRemove(b, m, size)
}

func BenchmarkHashBidiMapRemove100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := hashbimap.New[int, int]()

	for n := range size {
		m.Put(n, n)
	}

	b.StartTimer()
	benchmarkRemove(b, m, size)
}
