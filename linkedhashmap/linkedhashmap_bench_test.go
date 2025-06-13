package linkedhashmap_test

import (
	"testing"

	"github.com/qntx/gods/internal/testutil"
	"github.com/qntx/gods/linkedhashmap"
)

func benchmarkGet(b *testing.B, m *linkedhashmap.Map[int, int], keys []int) {
	b.Helper()

	for range b.N {
		for key := range keys {
			m.Get(key)
		}
	}
}

func benchmarkPut(b *testing.B, m *linkedhashmap.Map[int, int], keys []int) {
	b.Helper()

	for range b.N {
		for key := range keys {
			m.Put(key, key)
		}
	}
}

func benchmarkDelete(b *testing.B, m *linkedhashmap.Map[int, int], keys []int) {
	b.Helper()

	for range b.N {
		for key := range keys {
			m.Delete(key)
		}
	}
}

func BenchmarkTreeMapGet100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := linkedhashmap.New[int, int]()

	keys := testutil.GeneratePermutedInts(size)
	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkGet(b, m, keys)
}

func BenchmarkTreeMapGet1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := linkedhashmap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkGet(b, m, keys)
}

func BenchmarkTreeMapGet10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := linkedhashmap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkGet(b, m, keys)
}

func BenchmarkTreeMapGet100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := linkedhashmap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkGet(b, m, keys)
}

func BenchmarkTreeMapPut100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := linkedhashmap.New[int, int]()

	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, m, keys)
}

func BenchmarkTreeMapPut1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := linkedhashmap.New[int, int]()

	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, m, keys)
}

func BenchmarkTreeMapPut10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := linkedhashmap.New[int, int]()

	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, m, keys)
}

func BenchmarkTreeMapPut100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := linkedhashmap.New[int, int]()

	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, m, keys)
}

func BenchmarkTreeMapDelete100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := linkedhashmap.New[int, int]()

	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkDelete(b, m, keys)
}

func BenchmarkTreeMapDelete1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := linkedhashmap.New[int, int]()

	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkDelete(b, m, keys)
}

func BenchmarkTreeMapDelete10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := linkedhashmap.New[int, int]()

	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkDelete(b, m, keys)
}

func BenchmarkTreeMapDelete100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := linkedhashmap.New[int, int]()

	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkDelete(b, m, keys)
}
