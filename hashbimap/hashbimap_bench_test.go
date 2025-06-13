package hashbimap_test

import (
	"testing"

	"github.com/qntx/gods/hashbimap"
	"github.com/qntx/gods/internal/testutil"
)

func benchmarkGet(b *testing.B, m *hashbimap.Map[int, int], keys []int) {
	for range b.N {
		for key := range keys {
			m.Get(key)
		}
	}
}

func benchmarkPut(b *testing.B, m *hashbimap.Map[int, int], keys []int) {
	for range b.N {
		for key := range keys {
			m.Put(key, key)
		}
	}
}

func benchmarkDelete(b *testing.B, m *hashbimap.Map[int, int], keys []int) {
	for range b.N {
		for key := range keys {
			m.Delete(key)
		}
	}
}

func BenchmarkHashBidiMapGet100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := hashbimap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkGet(b, m, keys)
}

func BenchmarkHashBidiMapGet1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := hashbimap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkGet(b, m, keys)
}

func BenchmarkHashBidiMapGet10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := hashbimap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkGet(b, m, keys)
}

func BenchmarkHashBidiMapGet100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := hashbimap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkGet(b, m, keys)
}

func BenchmarkHashBidiMapPut100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := hashbimap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, m, keys)
}

func BenchmarkHashBidiMapPut1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := hashbimap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, m, keys)
}

func BenchmarkHashBidiMapPut10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := hashbimap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, m, keys)
}

func BenchmarkHashBidiMapPut100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := hashbimap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, m, keys)
}

func BenchmarkHashBidiMapDelete100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := hashbimap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkDelete(b, m, keys)
}

func BenchmarkHashBidiMapDelete1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := hashbimap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkDelete(b, m, keys)
}

func BenchmarkHashBidiMapDelete10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := hashbimap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkDelete(b, m, keys)
}

func BenchmarkHashBidiMapDelete100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := hashbimap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkDelete(b, m, keys)
}
