package rbtree_test

import (
	"testing"

	"github.com/qntx/gods/internal/testutil"
	"github.com/qntx/gods/rbtree"
)

func benchmarkGet(b *testing.B, tree *rbtree.Tree[int, struct{}], keys []int) {
	b.Helper()

	for range b.N {
		for key := range keys {
			tree.Get(key)
		}
	}
}

func benchmarkPut(b *testing.B, tree *rbtree.Tree[int, struct{}], keys []int) {
	b.Helper()

	for range b.N {
		for key := range keys {
			tree.Put(key, struct{}{})
		}
	}
}

func benchmarkDelete(b *testing.B, tree *rbtree.Tree[int, struct{}], keys []int) {
	b.Helper()

	for range b.N {
		for key := range keys {
			tree.Delete(key)
		}
	}
}

func BenchmarkRedBlackTreeGet100(b *testing.B) {
	b.StopTimer()

	size := 100
	tree := rbtree.New[int, struct{}]()
	keys := testutil.GeneratePermutedInts(size)

	for n := range keys {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkGet(b, tree, keys)
}

func BenchmarkRedBlackTreeGet1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	tree := rbtree.New[int, struct{}]()
	keys := testutil.GeneratePermutedInts(size)

	for n := range keys {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkGet(b, tree, keys)
}

func BenchmarkRedBlackTreeGet10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	tree := rbtree.New[int, struct{}]()
	keys := testutil.GeneratePermutedInts(size)

	for n := range keys {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkGet(b, tree, keys)
}

func BenchmarkRedBlackTreeGet100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	tree := rbtree.New[int, struct{}]()
	keys := testutil.GeneratePermutedInts(size)

	for n := range keys {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkGet(b, tree, keys)
}

func BenchmarkRedBlackTreePut100(b *testing.B) {
	b.StopTimer()

	size := 100
	tree := rbtree.New[int, struct{}]()
	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, tree, keys)
}

func BenchmarkRedBlackTreePut1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	tree := rbtree.New[int, struct{}]()
	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, tree, keys)
}

func BenchmarkRedBlackTreePut10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	tree := rbtree.New[int, struct{}]()
	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, tree, keys)
}

func BenchmarkRedBlackTreePut100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	tree := rbtree.New[int, struct{}]()
	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, tree, keys)
}

func BenchmarkRedBlackTreeDelete100(b *testing.B) {
	b.StopTimer()

	size := 100
	tree := rbtree.New[int, struct{}]()
	keys := testutil.GeneratePermutedInts(size)

	for n := range keys {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkDelete(b, tree, keys)
}

func BenchmarkRedBlackTreeDelete1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	tree := rbtree.New[int, struct{}]()
	keys := testutil.GeneratePermutedInts(size)

	for n := range keys {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkDelete(b, tree, keys)
}

func BenchmarkRedBlackTreeDelete10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	tree := rbtree.New[int, struct{}]()
	keys := testutil.GeneratePermutedInts(size)

	for n := range keys {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkDelete(b, tree, keys)
}

func BenchmarkRedBlackTreeDelete100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	tree := rbtree.New[int, struct{}]()
	keys := testutil.GeneratePermutedInts(size)

	for n := range keys {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkDelete(b, tree, keys)
}
