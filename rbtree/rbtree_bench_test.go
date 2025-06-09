package rbtree_test

import (
	"testing"

	"github.com/qntx/gods/rbtree"
)

func benchmarkGet(b *testing.B, tree *rbtree.Tree[int, struct{}], size int) {
	b.Helper()

	for b.N > 0 {
		for n := range size {
			tree.Get(n)
		}
	}
}

func benchmarkPut(b *testing.B, tree *rbtree.Tree[int, struct{}], size int) {
	b.Helper()

	for b.N > 0 {
		for n := range size {
			tree.Put(n, struct{}{})
		}
	}
}

func benchmarkRemove(b *testing.B, tree *rbtree.Tree[int, struct{}], size int) {
	b.Helper()

	for b.N > 0 {
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
