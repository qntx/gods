package avltree_test

import (
	"testing"

	"github.com/qntx/gods/avltree"
	"github.com/qntx/gods/internal/testutil"
)

func benchmarkGet(b *testing.B, tree *avltree.Tree[int, struct{}], keys []int) {
	b.Helper()

	for range b.N {
		for key := range keys {
			tree.Get(key)
		}
	}
}

func benchmarkPut(b *testing.B, tree *avltree.Tree[int, struct{}], keys []int) {
	b.Helper()

	for range b.N {
		for key := range keys {
			tree.Put(key, struct{}{})
		}
	}
}

func benchmarkDelete(b *testing.B, tree *avltree.Tree[int, struct{}], keys []int) {
	b.Helper()

	for range b.N {
		for key := range keys {
			tree.Delete(key)
		}
	}
}

func BenchmarkAVLTreeGet100(b *testing.B) {
	b.StopTimer()

	size := 100
	tree := avltree.New[int, struct{}]()

	keys := testutil.GeneratePermutedInts(size)
	for key := range keys {
		tree.Put(key, struct{}{})
	}

	b.StartTimer()
	benchmarkGet(b, tree, keys)
}

func BenchmarkAVLTreeGet1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	tree := avltree.New[int, struct{}]()

	keys := testutil.GeneratePermutedInts(size)
	for key := range keys {
		tree.Put(key, struct{}{})
	}

	b.StartTimer()
	benchmarkGet(b, tree, keys)
}

func BenchmarkAVLTreeGet10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	tree := avltree.New[int, struct{}]()

	keys := testutil.GeneratePermutedInts(size)
	for key := range keys {
		tree.Put(key, struct{}{})
	}

	b.StartTimer()
	benchmarkGet(b, tree, keys)
}

func BenchmarkAVLTreeGet100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	tree := avltree.New[int, struct{}]()

	keys := testutil.GeneratePermutedInts(size)
	for key := range keys {
		tree.Put(key, struct{}{})
	}

	b.StartTimer()
	benchmarkGet(b, tree, keys)
}

func BenchmarkAVLTreePut100(b *testing.B) {
	b.StopTimer()

	size := 100
	tree := avltree.New[int, struct{}]()

	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, tree, keys)
}

func BenchmarkAVLTreePut1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	tree := avltree.New[int, struct{}]()

	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, tree, keys)
}

func BenchmarkAVLTreePut10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	tree := avltree.New[int, struct{}]()

	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, tree, keys)
}

func BenchmarkAVLTreePut100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	tree := avltree.New[int, struct{}]()

	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPut(b, tree, keys)
}

func BenchmarkAVLTreeDelete100(b *testing.B) {
	b.StopTimer()

	size := 100
	tree := avltree.New[int, struct{}]()

	keys := testutil.GeneratePermutedInts(size)
	for key := range keys {
		tree.Put(key, struct{}{})
	}

	b.StartTimer()
	benchmarkDelete(b, tree, keys)
}

func BenchmarkAVLTreeDelete1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	tree := avltree.New[int, struct{}]()

	keys := testutil.GeneratePermutedInts(size)
	for key := range keys {
		tree.Put(key, struct{}{})
	}

	b.StartTimer()
	benchmarkDelete(b, tree, keys)
}

func BenchmarkAVLTreeDelete10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	tree := avltree.New[int, struct{}]()

	keys := testutil.GeneratePermutedInts(size)
	for key := range keys {
		tree.Put(key, struct{}{})
	}

	b.StartTimer()
	benchmarkDelete(b, tree, keys)
}

func BenchmarkAVLTreeDelete100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	tree := avltree.New[int, struct{}]()

	keys := testutil.GeneratePermutedInts(size)
	for key := range keys {
		tree.Put(key, struct{}{})
	}

	b.StartTimer()
	benchmarkDelete(b, tree, keys)
}
