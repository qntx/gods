package btree_test

import (
	"testing"

	"github.com/qntx/gods/btree"
	"github.com/qntx/gods/internal/testutil"
)

func benchmarkGet(b *testing.B, tree *btree.Tree[int, struct{}], keys []int) {
	for range b.N {
		for n := range keys {
			tree.Get(n)
		}
	}
}

func benchmarkPut(b *testing.B, tree *btree.Tree[int, struct{}], keys []int) {
	for range b.N {
		for n := range keys {
			tree.Put(n, struct{}{})
		}
	}
}

func benchmarkDelete(b *testing.B, tree *btree.Tree[int, struct{}], keys []int) {
	for range b.N {
		for n := range keys {
			tree.Delete(n)
		}
	}
}

func BenchmarkBTreeGet100(b *testing.B) {
	b.StopTimer()

	size := 100
	tree := btree.New[int, struct{}](128)

	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		tree.Put(key, struct{}{})
	}

	b.StartTimer()
	benchmarkGet(b, tree, keys)
}

func BenchmarkBTreeGet1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	tree := btree.New[int, struct{}](128)

	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		tree.Put(key, struct{}{})
	}

	b.StartTimer()
	benchmarkGet(b, tree, keys)
}

func BenchmarkBTreeGet10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	tree := btree.New[int, struct{}](128)

	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		tree.Put(key, struct{}{})
	}

	b.StartTimer()
	benchmarkGet(b, tree, keys)
}

func BenchmarkBTreeGet100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	tree := btree.New[int, struct{}](128)

	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		tree.Put(key, struct{}{})
	}

	b.StartTimer()
	benchmarkGet(b, tree, keys)
}

func BenchmarkBTreePut100(b *testing.B) {
	b.StopTimer()

	size := 100
	tree := btree.New[int, struct{}](128)

	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, tree, keys)
}

func BenchmarkBTreePut1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	tree := btree.New[int, struct{}](128)

	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, tree, keys)
}

func BenchmarkBTreePut10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	tree := btree.New[int, struct{}](128)

	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, tree, keys)
}

func BenchmarkBTreePut100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	tree := btree.New[int, struct{}](128)

	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, tree, keys)
}

func BenchmarkBTreeDelete100(b *testing.B) {
	b.StopTimer()

	size := 100
	tree := btree.New[int, struct{}](128)

	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		tree.Put(key, struct{}{})
	}

	b.StartTimer()
	benchmarkDelete(b, tree, keys)
}

func BenchmarkBTreeDelete1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	tree := btree.New[int, struct{}](128)

	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		tree.Put(key, struct{}{})
	}

	b.StartTimer()
	benchmarkDelete(b, tree, keys)
}

func BenchmarkBTreeDelete10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	tree := btree.New[int, struct{}](128)

	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		tree.Put(key, struct{}{})
	}

	b.StartTimer()
	benchmarkDelete(b, tree, keys)
}

func BenchmarkBTreeDelete100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	tree := btree.New[int, struct{}](128)

	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		tree.Put(key, struct{}{})
	}

	b.StartTimer()
	benchmarkDelete(b, tree, keys)
}
