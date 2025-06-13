package rbtreeset_test

import (
	"testing"

	"github.com/qntx/gods/internal/testutil"
	"github.com/qntx/gods/rbtreeset"
)

func benchmarkContains(b *testing.B, set *rbtreeset.Set[int], keys []int) {
	b.Helper()

	for range b.N {
		for key := range keys {
			set.Contains(key)
		}
	}
}

func benchmarkAdd(b *testing.B, set *rbtreeset.Set[int], keys []int) {
	b.Helper()

	for range b.N {
		for key := range keys {
			set.Add(key)
		}
	}
}

func benchmarkRemove(b *testing.B, set *rbtreeset.Set[int], keys []int) {
	b.Helper()

	for range b.N {
		for key := range keys {
			set.Remove(key)
		}
	}
}

func BenchmarkTreeSetContains100(b *testing.B) {
	b.StopTimer()

	size := 100
	set := rbtreeset.New[int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		set.Add(key)
	}

	b.StartTimer()
	benchmarkContains(b, set, keys)
}

func BenchmarkTreeSetContains1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	set := rbtreeset.New[int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		set.Add(key)
	}

	b.StartTimer()
	benchmarkContains(b, set, keys)
}

func BenchmarkTreeSetContains10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	set := rbtreeset.New[int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		set.Add(key)
	}

	b.StartTimer()
	benchmarkContains(b, set, keys)
}

func BenchmarkTreeSetContains100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	set := rbtreeset.New[int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		set.Add(key)
	}

	b.StartTimer()
	benchmarkContains(b, set, keys)
}

func BenchmarkTreeSetAdd100(b *testing.B) {
	b.StopTimer()

	size := 100
	set := rbtreeset.New[int]()
	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkAdd(b, set, keys)
}

func BenchmarkTreeSetAdd1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	set := rbtreeset.New[int]()
	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkAdd(b, set, keys)
}

func BenchmarkTreeSetAdd10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	set := rbtreeset.New[int]()
	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkAdd(b, set, keys)
}

func BenchmarkTreeSetAdd100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	set := rbtreeset.New[int]()
	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkAdd(b, set, keys)
}

func BenchmarkTreeSetRemove100(b *testing.B) {
	b.StopTimer()

	size := 100
	set := rbtreeset.New[int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		set.Add(key)
	}

	b.StartTimer()
	benchmarkRemove(b, set, keys)
}

func BenchmarkTreeSetRemove1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	set := rbtreeset.New[int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		set.Add(key)
	}

	b.StartTimer()
	benchmarkRemove(b, set, keys)
}

func BenchmarkTreeSetRemove10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	set := rbtreeset.New[int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		set.Add(key)
	}

	b.StartTimer()
	benchmarkRemove(b, set, keys)
}

func BenchmarkTreeSetRemove100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	set := rbtreeset.New[int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		set.Add(key)
	}

	b.StartTimer()
	benchmarkRemove(b, set, keys)
}
