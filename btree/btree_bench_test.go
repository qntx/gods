package btree_test

import (
	"testing"

	"github.com/qntx/gods/btree"
)

func benchmarkGet(b *testing.B, tree *btree.Tree[int, struct{}], size int) {
	for range b.N {
		for n := range size {
			tree.Get(n)
		}
	}
}

func benchmarkPut(b *testing.B, tree *btree.Tree[int, struct{}], size int) {
	for range b.N {
		for n := range size {
			tree.Put(n, struct{}{})
		}
	}
}

func benchmarkDelete(b *testing.B, tree *btree.Tree[int, struct{}], size int) {
	for range b.N {
		for n := range size {
			tree.Delete(n)
		}
	}
}

func BenchmarkBTreeGet100(b *testing.B) {
	b.StopTimer()

	size := 100
	tree := btree.New[int, struct{}](128)

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkGet(b, tree, size)
}

func BenchmarkBTreeGet1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	tree := btree.New[int, struct{}](128)

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkGet(b, tree, size)
}

func BenchmarkBTreeGet10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	tree := btree.New[int, struct{}](128)

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkGet(b, tree, size)
}

func BenchmarkBTreeGet100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	tree := btree.New[int, struct{}](128)

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkGet(b, tree, size)
}

func BenchmarkBTreePut100(b *testing.B) {
	b.StopTimer()

	size := 100
	tree := btree.New[int, struct{}](128)

	b.StartTimer()
	benchmarkPut(b, tree, size)
}

func BenchmarkBTreePut1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	tree := btree.New[int, struct{}](128)

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkPut(b, tree, size)
}

func BenchmarkBTreePut10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	tree := btree.New[int, struct{}](128)

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkPut(b, tree, size)
}

func BenchmarkBTreePut100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	tree := btree.New[int, struct{}](128)

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkPut(b, tree, size)
}

func BenchmarkBTreeDelete100(b *testing.B) {
	b.StopTimer()

	size := 100
	tree := btree.New[int, struct{}](128)

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkDelete(b, tree, size)
}

func BenchmarkBTreeDelete1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	tree := btree.New[int, struct{}](128)

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkDelete(b, tree, size)
}

func BenchmarkBTreeDelete10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	tree := btree.New[int, struct{}](128)

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkDelete(b, tree, size)
}

func BenchmarkBTreeDelete100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	tree := btree.New[int, struct{}](128)

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkDelete(b, tree, size)
}
