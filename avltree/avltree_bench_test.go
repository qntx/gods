package avltree_test

import (
	"testing"

	"github.com/qntx/gods/avltree"
)

func benchmarkGet(b *testing.B, tree *avltree.Tree[int, struct{}], size int) {
	for range b.N {
		for n := range size {
			tree.Get(n)
		}
	}
}

func benchmarkPut(b *testing.B, tree *avltree.Tree[int, struct{}], size int) {
	for range b.N {
		for n := range size {
			tree.Put(n, struct{}{})
		}
	}
}

func benchmarkRemove(b *testing.B, tree *avltree.Tree[int, struct{}], size int) {
	for range b.N {
		for n := range size {
			tree.Delete(n)
		}
	}
}

func BenchmarkAVLTreeGet100(b *testing.B) {
	b.StopTimer()

	size := 100
	tree := avltree.New[int, struct{}]()

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkGet(b, tree, size)
}

func BenchmarkAVLTreeGet1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	tree := avltree.New[int, struct{}]()

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkGet(b, tree, size)
}

func BenchmarkAVLTreeGet10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	tree := avltree.New[int, struct{}]()

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkGet(b, tree, size)
}

func BenchmarkAVLTreeGet100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	tree := avltree.New[int, struct{}]()

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkGet(b, tree, size)
}

func BenchmarkAVLTreePut100(b *testing.B) {
	b.StopTimer()

	size := 100
	tree := avltree.New[int, struct{}]()

	b.StartTimer()
	benchmarkPut(b, tree, size)
}

func BenchmarkAVLTreePut1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	tree := avltree.New[int, struct{}]()

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkPut(b, tree, size)
}

func BenchmarkAVLTreePut10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	tree := avltree.New[int, struct{}]()

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkPut(b, tree, size)
}

func BenchmarkAVLTreePut100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	tree := avltree.New[int, struct{}]()

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkPut(b, tree, size)
}

func BenchmarkAVLTreeRemove100(b *testing.B) {
	b.StopTimer()

	size := 100
	tree := avltree.New[int, struct{}]()

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkRemove(b, tree, size)
}

func BenchmarkAVLTreeRemove1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	tree := avltree.New[int, struct{}]()

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkRemove(b, tree, size)
}

func BenchmarkAVLTreeRemove10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	tree := avltree.New[int, struct{}]()

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkRemove(b, tree, size)
}

func BenchmarkAVLTreeRemove100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	tree := avltree.New[int, struct{}]()

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkRemove(b, tree, size)
}
