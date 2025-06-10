package linkedhashset_test

import (
	"testing"

	"github.com/qntx/gods/linkedhashset"
)

func benchmarkContains(b *testing.B, set *linkedhashset.Set[int], size int) {
	for range b.N {
		for n := range size {
			set.Contains(n)
		}
	}
}

func benchmarkAdd(b *testing.B, set *linkedhashset.Set[int], size int) {
	for range b.N {
		for n := range size {
			set.Add(n)
		}
	}
}

func benchmarkRemove(b *testing.B, set *linkedhashset.Set[int], size int) {
	for range b.N {
		for n := range size {
			set.Remove(n)
		}
	}
}

func BenchmarkHashSetContains100(b *testing.B) {
	b.StopTimer()

	size := 100
	set := linkedhashset.New[int]()

	for n := range size {
		set.Add(n)
	}

	b.StartTimer()
	benchmarkContains(b, set, size)
}

func BenchmarkHashSetContains1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	set := linkedhashset.New[int]()

	for n := range size {
		set.Add(n)
	}

	b.StartTimer()
	benchmarkContains(b, set, size)
}

func BenchmarkHashSetContains10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	set := linkedhashset.New[int]()

	for n := range size {
		set.Add(n)
	}

	b.StartTimer()
	benchmarkContains(b, set, size)
}

func BenchmarkHashSetContains100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	set := linkedhashset.New[int]()

	for n := range size {
		set.Add(n)
	}

	b.StartTimer()
	benchmarkContains(b, set, size)
}

func BenchmarkHashSetAdd100(b *testing.B) {
	b.StopTimer()

	size := 100
	set := linkedhashset.New[int]()

	b.StartTimer()
	benchmarkAdd(b, set, size)
}

func BenchmarkHashSetAdd1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	set := linkedhashset.New[int]()

	for n := range size {
		set.Add(n)
	}

	b.StartTimer()
	benchmarkAdd(b, set, size)
}

func BenchmarkHashSetAdd10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	set := linkedhashset.New[int]()

	for n := range size {
		set.Add(n)
	}

	b.StartTimer()
	benchmarkAdd(b, set, size)
}

func BenchmarkHashSetAdd100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	set := linkedhashset.New[int]()

	for n := range size {
		set.Add(n)
	}

	b.StartTimer()
	benchmarkAdd(b, set, size)
}

func BenchmarkHashSetRemove100(b *testing.B) {
	b.StopTimer()

	size := 100
	set := linkedhashset.New[int]()

	for n := range size {
		set.Add(n)
	}

	b.StartTimer()
	benchmarkRemove(b, set, size)
}

func BenchmarkHashSetRemove1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	set := linkedhashset.New[int]()

	for n := range size {
		set.Add(n)
	}

	b.StartTimer()
	benchmarkRemove(b, set, size)
}

func BenchmarkHashSetRemove10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	set := linkedhashset.New[int]()

	for n := range size {
		set.Add(n)
	}

	b.StartTimer()
	benchmarkRemove(b, set, size)
}

func BenchmarkHashSetRemove100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	set := linkedhashset.New[int]()

	for n := range size {
		set.Add(n)
	}

	b.StartTimer()
	benchmarkRemove(b, set, size)
}
