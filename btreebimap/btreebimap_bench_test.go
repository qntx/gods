package btreebimap_test

import (
	"testing"

	"github.com/qntx/gods/btreebimap"
	"github.com/qntx/gods/internal/testutil"
)

func benchmarkGet(b *testing.B, m *btreebimap.Map[int, int], keys []int) {
	for range b.N {
		for key := range keys {
			m.Get(key)
		}
	}
}

func benchmarkPut(b *testing.B, m *btreebimap.Map[int, int], keys []int) {
	for range b.N {
		for key := range keys {
			m.Put(key, key)
		}
	}
}

func benchmarkDelete(b *testing.B, m *btreebimap.Map[int, int], keys []int) {
	for range b.N {
		for key := range keys {
			m.Delete(key)
		}
	}
}

func BenchmarkBTreeBiMapGet100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := btreebimap.New[int, int]()

	keys := testutil.GeneratePermutedInts(size)
	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkGet(b, m, keys)
}

func BenchmarkBTreeBiMapGet1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := btreebimap.New[int, int]()

	keys := testutil.GeneratePermutedInts(size)
	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkGet(b, m, keys)
}

func BenchmarkBTreeBiMapGet10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := btreebimap.New[int, int]()

	keys := testutil.GeneratePermutedInts(size)
	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkGet(b, m, keys)
}

func BenchmarkBTreeBiMapGet100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := btreebimap.New[int, int]()

	keys := testutil.GeneratePermutedInts(size)
	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkGet(b, m, keys)
}

func BenchmarkBTreeBiMapPut100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := btreebimap.New[int, int]()

	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, m, keys)
}

func BenchmarkBTreeBiMapPut1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := btreebimap.New[int, int]()

	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, m, keys)
}

func BenchmarkBTreeBiMapPut10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := btreebimap.New[int, int]()

	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, m, keys)
}

func BenchmarkBTreeBiMapPut100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := btreebimap.New[int, int]()

	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, m, keys)
}

func BenchmarkBTreeBiMapDelete100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := btreebimap.New[int, int]()

	keys := testutil.GeneratePermutedInts(size)
	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkDelete(b, m, keys)
}

func BenchmarkBTreeBiMapDelete1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := btreebimap.New[int, int]()

	keys := testutil.GeneratePermutedInts(size)
	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkDelete(b, m, keys)
}

func BenchmarkBTreeBiMapDelete10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := btreebimap.New[int, int]()

	keys := testutil.GeneratePermutedInts(size)
	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkDelete(b, m, keys)
}

func BenchmarkBTreeBiMapDelete100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := btreebimap.New[int, int]()

	keys := testutil.GeneratePermutedInts(size)
	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkDelete(b, m, keys)
}
