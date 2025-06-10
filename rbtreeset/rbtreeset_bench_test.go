package rbtreeset_test

import (
	"testing"

	"github.com/qntx/gods/rbtreeset"
)

func benchmarkContains(b *testing.B, set *rbtreeset.Set[int], size int) {
	b.Helper()

	for b.N > 0 {
		for n := range size {
			set.Contains(n)
		}
	}
}

func benchmarkAdd(b *testing.B, set *rbtreeset.Set[int], size int) {
	b.Helper()

	for b.N > 0 {
		for n := range size {
			set.Add(n)
		}
	}
}

func benchmarkRemove(b *testing.B, set *rbtreeset.Set[int], size int) {
	b.Helper()

	for b.N > 0 {
		for n := range size {
			set.Remove(n)
		}
	}
}

func BenchmarkTreeSetContains100(b *testing.B) {
	b.StopTimer()

	size := 100
	set := rbtreeset.New[int]()

	for n := range size {
		set.Add(n)
	}

	b.StartTimer()
	benchmarkContains(b, set, size)
}

func BenchmarkTreeSetContains1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	set := rbtreeset.New[int]()

	for n := range size {
		set.Add(n)
	}

	b.StartTimer()
	benchmarkContains(b, set, size)
}

func BenchmarkTreeSetContains10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	set := rbtreeset.New[int]()

	for n := range size {
		set.Add(n)
	}

	b.StartTimer()
	benchmarkContains(b, set, size)
}

func BenchmarkTreeSetContains100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	set := rbtreeset.New[int]()

	for n := range size {
		set.Add(n)
	}

	b.StartTimer()
	benchmarkContains(b, set, size)
}

func BenchmarkTreeSetAdd100(b *testing.B) {
	b.StopTimer()

	size := 100
	set := rbtreeset.New[int]()

	b.StartTimer()
	benchmarkAdd(b, set, size)
}

func BenchmarkTreeSetAdd1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	set := rbtreeset.New[int]()

	for n := range size {
		set.Add(n)
	}

	b.StartTimer()
	benchmarkAdd(b, set, size)
}

func BenchmarkTreeSetAdd10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	set := rbtreeset.New[int]()

	for n := range size {
		set.Add(n)
	}

	b.StartTimer()
	benchmarkAdd(b, set, size)
}

func BenchmarkTreeSetAdd100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	set := rbtreeset.New[int]()

	for n := range size {
		set.Add(n)
	}

	b.StartTimer()
	benchmarkAdd(b, set, size)
}

func BenchmarkTreeSetRemove100(b *testing.B) {
	b.StopTimer()

	size := 100
	set := rbtreeset.New[int]()

	for n := range size {
		set.Add(n)
	}

	b.StartTimer()
	benchmarkRemove(b, set, size)
}

func BenchmarkTreeSetRemove1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	set := rbtreeset.New[int]()

	for n := range size {
		set.Add(n)
	}

	b.StartTimer()
	benchmarkRemove(b, set, size)
}

func BenchmarkTreeSetRemove10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	set := rbtreeset.New[int]()

	for n := range size {
		set.Add(n)
	}

	b.StartTimer()
	benchmarkRemove(b, set, size)
}

func BenchmarkTreeSetRemove100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	set := rbtreeset.New[int]()

	for n := range size {
		set.Add(n)
	}

	b.StartTimer()
	benchmarkRemove(b, set, size)
}
