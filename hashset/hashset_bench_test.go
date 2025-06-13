package hashset_test

import (
	"testing"

	"github.com/qntx/gods/hashset"
	"github.com/qntx/gods/internal/testutil"
)

func benchmarkContains(b *testing.B, set *hashset.Set[int], keys []int) {
	b.Helper()

	for range b.N {
		for key := range keys {
			set.Contains(key)
		}
	}
}

func benchmarkAdd(b *testing.B, set *hashset.Set[int], keys []int) {
	b.Helper()

	for range b.N {
		for key := range keys {
			set.Add(key)
		}
	}
}

func benchmarkRemove(b *testing.B, set *hashset.Set[int], keys []int) {
	b.Helper()

	for range b.N {
		for key := range keys {
			set.Remove(key)
		}
	}
}

func BenchmarkHashSetContains100(b *testing.B) {
	b.StopTimer()

	size := 100
	set := hashset.New[int]()

	keys := testutil.GeneratePermutedInts(size)
	for key := range keys {
		set.Add(key)
	}

	b.StartTimer()
	benchmarkContains(b, set, keys)
}

func BenchmarkHashSetContains1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	set := hashset.New[int]()

	keys := testutil.GeneratePermutedInts(size)
	for key := range keys {
		set.Add(key)
	}

	b.StartTimer()
	benchmarkContains(b, set, keys)
}

func BenchmarkHashSetContains10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	set := hashset.New[int]()

	keys := testutil.GeneratePermutedInts(size)
	for key := range keys {
		set.Add(key)
	}

	b.StartTimer()
	benchmarkContains(b, set, keys)
}

func BenchmarkHashSetContains100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	set := hashset.New[int]()

	keys := testutil.GeneratePermutedInts(size)
	for key := range keys {
		set.Add(key)
	}

	b.StartTimer()
	benchmarkContains(b, set, keys)
}

func BenchmarkHashSetAdd100(b *testing.B) {
	b.StopTimer()

	size := 100
	set := hashset.New[int]()
	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkAdd(b, set, keys)
}

func BenchmarkHashSetAdd1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	set := hashset.New[int]()

	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkAdd(b, set, keys)
}

func BenchmarkHashSetAdd10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	set := hashset.New[int]()

	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkAdd(b, set, keys)
}

func BenchmarkHashSetAdd100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	set := hashset.New[int]()

	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkAdd(b, set, keys)
}

func BenchmarkHashSetRemove100(b *testing.B) {
	b.StopTimer()

	size := 100
	set := hashset.New[int]()

	keys := testutil.GeneratePermutedInts(size)
	for key := range keys {
		set.Add(key)
	}

	b.StartTimer()
	benchmarkRemove(b, set, keys)
}

func BenchmarkHashSetRemove1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	set := hashset.New[int]()

	keys := testutil.GeneratePermutedInts(size)
	for key := range keys {
		set.Add(key)
	}

	b.StartTimer()
	benchmarkRemove(b, set, keys)
}

func BenchmarkHashSetRemove10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	set := hashset.New[int]()

	keys := testutil.GeneratePermutedInts(size)
	for key := range keys {
		set.Add(key)
	}

	b.StartTimer()
	benchmarkRemove(b, set, keys)
}

func BenchmarkHashSetRemove100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	set := hashset.New[int]()

	keys := testutil.GeneratePermutedInts(size)
	for key := range keys {
		set.Add(key)
	}

	b.StartTimer()
	benchmarkRemove(b, set, keys)
}
