package rbtreebimap_test

import (
	"testing"

	"github.com/qntx/gods/internal/testutil"
	"github.com/qntx/gods/rbtreebimap"
)

func benchmarkGet(b *testing.B, m *rbtreebimap.Map[int, int], keys []int) {
	b.Helper()

	for range b.N {
		for key := range keys {
			m.Get(key)
		}
	}
}

func benchmarkPut(b *testing.B, m *rbtreebimap.Map[int, int], keys []int) {
	b.Helper()

	for range b.N {
		for key := range keys {
			m.Put(key, key)
		}
	}
}

func benchmarkDelete(b *testing.B, m *rbtreebimap.Map[int, int], keys []int) {
	b.Helper()

	for range b.N {
		for key := range keys {
			m.Delete(key)
		}
	}
}

func BenchmarkTreeBidiMapGet100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := rbtreebimap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkGet(b, m, keys)
}

func BenchmarkTreeBidiMapGet1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := rbtreebimap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkGet(b, m, keys)
}

func BenchmarkTreeBidiMapGet10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := rbtreebimap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkGet(b, m, keys)
}

func BenchmarkTreeBidiMapGet100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := rbtreebimap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkGet(b, m, keys)
}

func BenchmarkTreeBidiMapPut100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := rbtreebimap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, m, keys)
}

func BenchmarkTreeBidiMapPut1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := rbtreebimap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, m, keys)
}

func BenchmarkTreeBidiMapPut10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := rbtreebimap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, m, keys)
}

func BenchmarkTreeBidiMapPut100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := rbtreebimap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, m, keys)
}

func BenchmarkTreeBidiMapDelete100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := rbtreebimap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkDelete(b, m, keys)
}

func BenchmarkTreeBidiMapDelete1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := rbtreebimap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkDelete(b, m, keys)
}

func BenchmarkTreeBidiMapDelete10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := rbtreebimap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkDelete(b, m, keys)
}

func BenchmarkTreeBidiMapDelete100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := rbtreebimap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkDelete(b, m, keys)
}
