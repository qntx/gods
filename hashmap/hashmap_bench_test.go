package hashmap_test

import (
	"testing"

	"github.com/qntx/gods/hashmap"
	"github.com/qntx/gods/internal/testutil"
)

func benchmarkGet(b *testing.B, m *hashmap.Map[int, int], keys []int) {
	for range b.N {
		for n := range keys {
			m.Get(n)
		}
	}
}

func benchmarkPut(b *testing.B, m *hashmap.Map[int, int], keys []int) {
	for range b.N {
		for n := range keys {
			m.Put(n, n)
		}
	}
}

func benchmarkDelete(b *testing.B, m *hashmap.Map[int, int], keys []int) {
	for range b.N {
		for n := range keys {
			m.Delete(n)
		}
	}
}

func BenchmarkHashMapGet100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := hashmap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkGet(b, m, keys)
}

func BenchmarkHashMapGet1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := hashmap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkGet(b, m, keys)
}

func BenchmarkHashMapGet10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := hashmap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkGet(b, m, keys)
}

func BenchmarkHashMapGet100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := hashmap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkGet(b, m, keys)
}

func BenchmarkHashMapPut100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := hashmap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, m, keys)
}

func BenchmarkHashMapPut1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := hashmap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, m, keys)
}

func BenchmarkHashMapPut10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := hashmap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, m, keys)
}

func BenchmarkHashMapPut100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := hashmap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, m, keys)
}

func BenchmarkHashMapDelete100(b *testing.B) {
	b.StopTimer()

	size := 100
	m := hashmap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkDelete(b, m, keys)
}

func BenchmarkHashMapDelete1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	m := hashmap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkDelete(b, m, keys)
}

func BenchmarkHashMapDelete10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	m := hashmap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkDelete(b, m, keys)
}

func BenchmarkHashMapDelete100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	m := hashmap.New[int, int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		m.Put(key, key)
	}

	b.StartTimer()
	benchmarkDelete(b, m, keys)
}
