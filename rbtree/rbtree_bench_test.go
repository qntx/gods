package rbtree_test

import (
	"sort"
	"testing"

	"github.com/qntx/gods/rbtree"
)

const (
	defaultSize = 5000 // Default benchmark size for consistent testing
)

// BenchmarkRedBlackTree measures the performance of red-black tree operations.
// It tests insertion and key retrieval separately for clarity.
func BenchmarkRedBlackTree(b *testing.B) {
	b.Run("Insert", func(b *testing.B) {
		for b.Loop() {
			t := rbtree.New[int, struct{}]()
			for i := range defaultSize {
				t.Put(i, struct{}{})
			}
		}
	})

	t := rbtree.New[int, struct{}]()
	for i := range defaultSize {
		t.Put(i, struct{}{})
	}

	b.Run("Keys", func(b *testing.B) {
		b.ResetTimer()

		for b.Loop() {
			_ = t.Keys()
		}
	})
}

// BenchmarkMap measures the performance of Go map operations with sorted keys.
// It tests insertion and sorted key retrieval separately for clarity.
func BenchmarkMap(b *testing.B) {
	b.Run("Insert", func(b *testing.B) {
		for b.Loop() {
			m := make(map[int]struct{}, defaultSize)
			for i := range defaultSize {
				m[i] = struct{}{}
			}
		}
	})

	m := make(map[int]struct{}, defaultSize)
	for i := range defaultSize {
		m[i] = struct{}{}
	}

	b.Run("SortedKeys", func(b *testing.B) {
		b.ResetTimer()

		for b.Loop() {
			keys := make([]int, 0, defaultSize)
			for k := range m {
				keys = append(keys, k)
			}

			sort.Ints(keys)
		}
	})
}

func benchmarkGet(b *testing.B, tree *rbtree.Tree[int, struct{}], size int) {
	b.Helper()

	for b.Loop() {
		for n := range size {
			tree.Get(n)
		}
	}
}

func benchmarkPut(b *testing.B, tree *rbtree.Tree[int, struct{}], size int) {
	b.Helper()

	for b.Loop() {
		for n := range size {
			tree.Put(n, struct{}{})
		}
	}
}

func benchmarkRemove(b *testing.B, tree *rbtree.Tree[int, struct{}], size int) {
	b.Helper()

	for b.Loop() {
		for n := range size {
			tree.Remove(n)
		}
	}
}

func BenchmarkRedBlackTreeGet100(b *testing.B) {
	b.StopTimer()

	size := 100
	tree := rbtree.New[int, struct{}]()

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkGet(b, tree, size)
}

func BenchmarkRedBlackTreeGet1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	tree := rbtree.New[int, struct{}]()

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkGet(b, tree, size)
}

func BenchmarkRedBlackTreeGet10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	tree := rbtree.New[int, struct{}]()

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkGet(b, tree, size)
}

func BenchmarkRedBlackTreeGet100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	tree := rbtree.New[int, struct{}]()

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkGet(b, tree, size)
}

func BenchmarkRedBlackTreePut100(b *testing.B) {
	b.StopTimer()

	size := 100
	tree := rbtree.New[int, struct{}]()

	b.StartTimer()
	benchmarkPut(b, tree, size)
}

func BenchmarkRedBlackTreePut1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	tree := rbtree.New[int, struct{}]()

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkPut(b, tree, size)
}

func BenchmarkRedBlackTreePut10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	tree := rbtree.New[int, struct{}]()

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkPut(b, tree, size)
}

func BenchmarkRedBlackTreePut100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	tree := rbtree.New[int, struct{}]()

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkPut(b, tree, size)
}

func BenchmarkRedBlackTreeRemove100(b *testing.B) {
	b.StopTimer()

	size := 100
	tree := rbtree.New[int, struct{}]()

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkRemove(b, tree, size)
}

func BenchmarkRedBlackTreeRemove1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	tree := rbtree.New[int, struct{}]()

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkRemove(b, tree, size)
}

func BenchmarkRedBlackTreeRemove10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	tree := rbtree.New[int, struct{}]()

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkRemove(b, tree, size)
}

func BenchmarkRedBlackTreeRemove100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	tree := rbtree.New[int, struct{}]()

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkRemove(b, tree, size)
}
