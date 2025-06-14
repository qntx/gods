package btree

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"
)

func assertValidTree[K comparable, V any](t *testing.T, tree *Tree[K, V], expectedSize int) {
	if actualValue, expectedValue := tree.len, expectedSize; actualValue != expectedValue {
		t.Errorf("Got %v expected %v for tree size", actualValue, expectedValue)
	}
}

func assertValidTreeNode[K comparable, V any](t *testing.T, node *Node[K, V], expectedEntries int, expectedChildren int, keys []K, hasParent bool) {
	if actualValue, expectedValue := node.Parent() != nil, hasParent; actualValue != expectedValue {
		t.Errorf("Got %v expected %v for hasParent", actualValue, expectedValue)
	}

	if actualValue, expectedValue := len(node.Entries()), expectedEntries; actualValue != expectedValue {
		t.Errorf("Got %v expected %v for entries size", actualValue, expectedValue)
	}

	if actualValue, expectedValue := len(node.Children()), expectedChildren; actualValue != expectedValue {
		t.Errorf("Got %v expected %v for children size", actualValue, expectedValue)
	}

	for i, key := range keys {
		if actualValue, expectedValue := node.Entries()[i].Key(), key; actualValue != expectedValue {
			t.Errorf("Got %v expected %v for key", actualValue, expectedValue)
		}
	}
}

func TestBTreeGet1(t *testing.T) {
	tree := New[int, string](3)
	tree.Put(1, "a")
	tree.Put(2, "b")
	tree.Put(3, "c")
	tree.Put(4, "d")
	tree.Put(5, "e")
	tree.Put(6, "f")
	tree.Put(7, "g")

	tests := [][]interface{}{
		{0, "", false},
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "e", true},
		{6, "f", true},
		{7, "g", true},
		{8, "", false},
	}

	for _, test := range tests {
		if value, found := tree.Get(test[0].(int)); value != test[1] || found != test[2] {
			t.Errorf("Got %v,%v expected %v,%v", value, found, test[1], test[2])
		}
	}
}

func TestBTreeGet2(t *testing.T) {
	tree := New[int, string](3)
	tree.Put(7, "g")
	tree.Put(9, "i")
	tree.Put(10, "j")
	tree.Put(6, "f")
	tree.Put(3, "c")
	tree.Put(4, "d")
	tree.Put(5, "e")
	tree.Put(8, "h")
	tree.Put(2, "b")
	tree.Put(1, "a")

	tests := [][]interface{}{
		{0, "", false},
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "e", true},
		{6, "f", true},
		{7, "g", true},
		{8, "h", true},
		{9, "i", true},
		{10, "j", true},
		{11, "", false},
	}

	for _, test := range tests {
		if value, found := tree.Get(test[0].(int)); value != test[1] || found != test[2] {
			t.Errorf("Got %v,%v expected %v,%v", value, found, test[1], test[2])
		}
	}
}

func TestBTreeGet3(t *testing.T) {
	tree := New[int, string](3)

	if actualValue := tree.Len(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}

	if actualValue := tree.GetNode(2).Size(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}

	tree.Put(1, "x") // 1->x
	tree.Put(2, "b") // 1->x, 2->b (in order)
	tree.Put(1, "a") // 1->a, 2->b (in order, replacement)
	tree.Put(3, "c") // 1->a, 2->b, 3->c (in order)
	tree.Put(4, "d") // 1->a, 2->b, 3->c, 4->d (in order)
	tree.Put(5, "e") // 1->a, 2->b, 3->c, 4->d, 5->e (in order)
	tree.Put(6, "f") // 1->a, 2->b, 3->c, 4->d, 5->e, 6->f (in order)
	tree.Put(7, "g") // 1->a, 2->b, 3->c, 4->d, 5->e, 6->f, 7->g (in order)

	// BTree
	//         1
	//     2
	//         3
	// 4
	//         5
	//     6
	//         7

	if actualValue := tree.Len(); actualValue != 7 {
		t.Errorf("Got %v expected %v", actualValue, 7)
	}

	if actualValue := tree.GetNode(2).Size(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}

	if actualValue := tree.GetNode(4).Size(); actualValue != 7 {
		t.Errorf("Got %v expected %v", actualValue, 7)
	}

	if actualValue := tree.GetNode(8).Size(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}
}

func TestBTreePut1(t *testing.T) {
	// https://upload.wikimedia.org/wikipedia/commons/3/33/B_tree_insertion_example.png
	tree := New[int, int](3)
	assertValidTree(t, tree, 0)

	tree.Put(1, 0)
	assertValidTree(t, tree, 1)
	assertValidTreeNode(t, tree.Root(), 1, 0, []int{1}, false)

	tree.Put(2, 1)
	assertValidTree(t, tree, 2)
	assertValidTreeNode(t, tree.Root(), 2, 0, []int{1, 2}, false)

	tree.Put(3, 2)
	assertValidTree(t, tree, 3)
	assertValidTreeNode(t, tree.Root(), 1, 2, []int{2}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 1, 0, []int{3}, true)

	tree.Put(4, 2)
	assertValidTree(t, tree, 4)
	assertValidTreeNode(t, tree.Root(), 1, 2, []int{2}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 2, 0, []int{3, 4}, true)

	tree.Put(5, 2)
	assertValidTree(t, tree, 5)
	assertValidTreeNode(t, tree.Root(), 2, 3, []int{2, 4}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 1, 0, []int{3}, true)
	assertValidTreeNode(t, tree.Root().Children()[2], 1, 0, []int{5}, true)

	tree.Put(6, 2)
	assertValidTree(t, tree, 6)
	assertValidTreeNode(t, tree.Root(), 2, 3, []int{2, 4}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 1, 0, []int{3}, true)
	assertValidTreeNode(t, tree.Root().Children()[2], 2, 0, []int{5, 6}, true)

	tree.Put(7, 2)
	assertValidTree(t, tree, 7)
	assertValidTreeNode(t, tree.Root(), 1, 2, []int{4}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 1, 2, []int{2}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 1, 2, []int{6}, true)
	assertValidTreeNode(t, tree.Root().Children()[0].Children()[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root().Children()[0].Children()[1], 1, 0, []int{3}, true)
	assertValidTreeNode(t, tree.Root().Children()[1].Children()[0], 1, 0, []int{5}, true)
	assertValidTreeNode(t, tree.Root().Children()[1].Children()[1], 1, 0, []int{7}, true)
}

func TestBTreePut2(t *testing.T) {
	tree := New[int, int](4)
	assertValidTree(t, tree, 0)

	tree.Put(0, 0)
	assertValidTree(t, tree, 1)
	assertValidTreeNode(t, tree.Root(), 1, 0, []int{0}, false)

	tree.Put(2, 2)
	assertValidTree(t, tree, 2)
	assertValidTreeNode(t, tree.Root(), 2, 0, []int{0, 2}, false)

	tree.Put(1, 1)
	assertValidTree(t, tree, 3)
	assertValidTreeNode(t, tree.Root(), 3, 0, []int{0, 1, 2}, false)

	tree.Put(1, 1)
	assertValidTree(t, tree, 3)
	assertValidTreeNode(t, tree.Root(), 3, 0, []int{0, 1, 2}, false)

	tree.Put(3, 3)
	assertValidTree(t, tree, 4)
	assertValidTreeNode(t, tree.Root(), 1, 2, []int{1}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 1, 0, []int{0}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 2, 0, []int{2, 3}, true)

	tree.Put(4, 4)
	assertValidTree(t, tree, 5)
	assertValidTreeNode(t, tree.Root(), 1, 2, []int{1}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 1, 0, []int{0}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 3, 0, []int{2, 3, 4}, true)

	tree.Put(5, 5)
	assertValidTree(t, tree, 6)
	assertValidTreeNode(t, tree.Root(), 2, 3, []int{1, 3}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 1, 0, []int{0}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 1, 0, []int{2}, true)
	assertValidTreeNode(t, tree.Root().Children()[2], 2, 0, []int{4, 5}, true)
}

func TestBTreePut3(t *testing.T) {
	// http://www.geeksforgeeks.org/b-tree-set-1-insert-2/
	tree := New[int, int](6)
	assertValidTree(t, tree, 0)

	tree.Put(10, 0)
	assertValidTree(t, tree, 1)
	assertValidTreeNode(t, tree.Root(), 1, 0, []int{10}, false)

	tree.Put(20, 1)
	assertValidTree(t, tree, 2)
	assertValidTreeNode(t, tree.Root(), 2, 0, []int{10, 20}, false)

	tree.Put(30, 2)
	assertValidTree(t, tree, 3)
	assertValidTreeNode(t, tree.Root(), 3, 0, []int{10, 20, 30}, false)

	tree.Put(40, 3)
	assertValidTree(t, tree, 4)
	assertValidTreeNode(t, tree.Root(), 4, 0, []int{10, 20, 30, 40}, false)

	tree.Put(50, 4)
	assertValidTree(t, tree, 5)
	assertValidTreeNode(t, tree.Root(), 5, 0, []int{10, 20, 30, 40, 50}, false)

	tree.Put(60, 5)
	assertValidTree(t, tree, 6)
	assertValidTreeNode(t, tree.Root(), 1, 2, []int{30}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 2, 0, []int{10, 20}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 3, 0, []int{40, 50, 60}, true)

	tree.Put(70, 6)
	assertValidTree(t, tree, 7)
	assertValidTreeNode(t, tree.Root(), 1, 2, []int{30}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 2, 0, []int{10, 20}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 4, 0, []int{40, 50, 60, 70}, true)

	tree.Put(80, 7)
	assertValidTree(t, tree, 8)
	assertValidTreeNode(t, tree.Root(), 1, 2, []int{30}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 2, 0, []int{10, 20}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 5, 0, []int{40, 50, 60, 70, 80}, true)

	tree.Put(90, 8)
	assertValidTree(t, tree, 9)
	assertValidTreeNode(t, tree.Root(), 2, 3, []int{30, 60}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 2, 0, []int{10, 20}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 2, 0, []int{40, 50}, true)
	assertValidTreeNode(t, tree.Root().Children()[2], 3, 0, []int{70, 80, 90}, true)
}

func TestBTreePut4(t *testing.T) {
	tree := New[int, *struct{}](3)
	assertValidTree(t, tree, 0)

	tree.Put(6, nil)
	assertValidTree(t, tree, 1)
	assertValidTreeNode(t, tree.Root(), 1, 0, []int{6}, false)

	tree.Put(5, nil)
	assertValidTree(t, tree, 2)
	assertValidTreeNode(t, tree.Root(), 2, 0, []int{5, 6}, false)

	tree.Put(4, nil)
	assertValidTree(t, tree, 3)
	assertValidTreeNode(t, tree.Root(), 1, 2, []int{5}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 1, 0, []int{4}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 1, 0, []int{6}, true)

	tree.Put(3, nil)
	assertValidTree(t, tree, 4)
	assertValidTreeNode(t, tree.Root(), 1, 2, []int{5}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 2, 0, []int{3, 4}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 1, 0, []int{6}, true)

	tree.Put(2, nil)
	assertValidTree(t, tree, 5)
	assertValidTreeNode(t, tree.Root(), 2, 3, []int{3, 5}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 1, 0, []int{2}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 1, 0, []int{4}, true)
	assertValidTreeNode(t, tree.Root().Children()[2], 1, 0, []int{6}, true)

	tree.Put(1, nil)
	assertValidTree(t, tree, 6)
	assertValidTreeNode(t, tree.Root(), 2, 3, []int{3, 5}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 2, 0, []int{1, 2}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 1, 0, []int{4}, true)
	assertValidTreeNode(t, tree.Root().Children()[2], 1, 0, []int{6}, true)

	tree.Put(0, nil)
	assertValidTree(t, tree, 7)
	assertValidTreeNode(t, tree.Root(), 1, 2, []int{3}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 1, 2, []int{1}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 1, 2, []int{5}, true)
	assertValidTreeNode(t, tree.Root().Children()[0].Children()[0], 1, 0, []int{0}, true)
	assertValidTreeNode(t, tree.Root().Children()[0].Children()[1], 1, 0, []int{2}, true)
	assertValidTreeNode(t, tree.Root().Children()[1].Children()[0], 1, 0, []int{4}, true)
	assertValidTreeNode(t, tree.Root().Children()[1].Children()[1], 1, 0, []int{6}, true)

	tree.Put(-1, nil)
	assertValidTree(t, tree, 8)
	assertValidTreeNode(t, tree.Root(), 1, 2, []int{3}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 1, 2, []int{1}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 1, 2, []int{5}, true)
	assertValidTreeNode(t, tree.Root().Children()[0].Children()[0], 2, 0, []int{-1, 0}, true)
	assertValidTreeNode(t, tree.Root().Children()[0].Children()[1], 1, 0, []int{2}, true)
	assertValidTreeNode(t, tree.Root().Children()[1].Children()[0], 1, 0, []int{4}, true)
	assertValidTreeNode(t, tree.Root().Children()[1].Children()[1], 1, 0, []int{6}, true)

	tree.Put(-2, nil)
	assertValidTree(t, tree, 9)
	assertValidTreeNode(t, tree.Root(), 1, 2, []int{3}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 2, 3, []int{-1, 1}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 1, 2, []int{5}, true)
	assertValidTreeNode(t, tree.Root().Children()[0].Children()[0], 1, 0, []int{-2}, true)
	assertValidTreeNode(t, tree.Root().Children()[0].Children()[1], 1, 0, []int{0}, true)
	assertValidTreeNode(t, tree.Root().Children()[0].Children()[2], 1, 0, []int{2}, true)
	assertValidTreeNode(t, tree.Root().Children()[1].Children()[0], 1, 0, []int{4}, true)
	assertValidTreeNode(t, tree.Root().Children()[1].Children()[1], 1, 0, []int{6}, true)

	tree.Put(-3, nil)
	assertValidTree(t, tree, 10)
	assertValidTreeNode(t, tree.Root(), 1, 2, []int{3}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 2, 3, []int{-1, 1}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 1, 2, []int{5}, true)
	assertValidTreeNode(t, tree.Root().Children()[0].Children()[0], 2, 0, []int{-3, -2}, true)
	assertValidTreeNode(t, tree.Root().Children()[0].Children()[1], 1, 0, []int{0}, true)
	assertValidTreeNode(t, tree.Root().Children()[0].Children()[2], 1, 0, []int{2}, true)
	assertValidTreeNode(t, tree.Root().Children()[1].Children()[0], 1, 0, []int{4}, true)
	assertValidTreeNode(t, tree.Root().Children()[1].Children()[1], 1, 0, []int{6}, true)

	tree.Put(-4, nil)
	assertValidTree(t, tree, 11)
	assertValidTreeNode(t, tree.Root(), 2, 3, []int{-1, 3}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 1, 2, []int{-3}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 1, 2, []int{1}, true)
	assertValidTreeNode(t, tree.Root().Children()[2], 1, 2, []int{5}, true)
	assertValidTreeNode(t, tree.Root().Children()[0].Children()[0], 1, 0, []int{-4}, true)
	assertValidTreeNode(t, tree.Root().Children()[0].Children()[1], 1, 0, []int{-2}, true)
	assertValidTreeNode(t, tree.Root().Children()[1].Children()[0], 1, 0, []int{0}, true)
	assertValidTreeNode(t, tree.Root().Children()[1].Children()[1], 1, 0, []int{2}, true)
	assertValidTreeNode(t, tree.Root().Children()[2].Children()[0], 1, 0, []int{4}, true)
	assertValidTreeNode(t, tree.Root().Children()[2].Children()[1], 1, 0, []int{6}, true)
}

func TestBTreeDelete1(t *testing.T) {
	// empty
	tree := New[int, int](3)
	tree.Delete(1)
	assertValidTree(t, tree, 0)
}

func TestBTreeDelete2(t *testing.T) {
	// leaf node (no underflow)
	tree := New[int, *struct{}](3)
	tree.Put(1, nil)
	tree.Put(2, nil)

	tree.Delete(1)
	assertValidTree(t, tree, 1)
	assertValidTreeNode(t, tree.Root(), 1, 0, []int{2}, false)

	tree.Delete(2)
	assertValidTree(t, tree, 0)
}

func TestBTreeDelete3(t *testing.T) {
	// merge with right (underflow)
	{
		tree := New[int, *struct{}](3)
		tree.Put(1, nil)
		tree.Put(2, nil)
		tree.Put(3, nil)

		tree.Delete(1)
		assertValidTree(t, tree, 2)
		assertValidTreeNode(t, tree.Root(), 2, 0, []int{2, 3}, false)
	}
	// merge with left (underflow)
	{
		tree := New[int, *struct{}](3)
		tree.Put(1, nil)
		tree.Put(2, nil)
		tree.Put(3, nil)

		tree.Delete(3)
		assertValidTree(t, tree, 2)
		assertValidTreeNode(t, tree.Root(), 2, 0, []int{1, 2}, false)
	}
}

func TestBTreeDelete4(t *testing.T) {
	// rotate left (underflow)
	tree := New[int, *struct{}](3)
	tree.Put(1, nil)
	tree.Put(2, nil)
	tree.Put(3, nil)
	tree.Put(4, nil)

	assertValidTree(t, tree, 4)
	assertValidTreeNode(t, tree.Root(), 1, 2, []int{2}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 2, 0, []int{3, 4}, true)

	tree.Delete(1)
	assertValidTree(t, tree, 3)
	assertValidTreeNode(t, tree.Root(), 1, 2, []int{3}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 1, 0, []int{2}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 1, 0, []int{4}, true)
}

func TestBTreeDelete5(t *testing.T) {
	// rotate right (underflow)
	tree := New[int, *struct{}](3)
	tree.Put(1, nil)
	tree.Put(2, nil)
	tree.Put(3, nil)
	tree.Put(0, nil)

	assertValidTree(t, tree, 4)
	assertValidTreeNode(t, tree.Root(), 1, 2, []int{2}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 2, 0, []int{0, 1}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 1, 0, []int{3}, true)

	tree.Delete(3)
	assertValidTree(t, tree, 3)
	assertValidTreeNode(t, tree.Root(), 1, 2, []int{1}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 1, 0, []int{0}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 1, 0, []int{2}, true)
}

func TestBTreeDelete6(t *testing.T) {
	// root height reduction after a series of underflows on right side
	// use simulator: https://www.cs.usfca.edu/~galles/visualization/BTree.html
	tree := New[int, *struct{}](3)
	tree.Put(1, nil)
	tree.Put(2, nil)
	tree.Put(3, nil)
	tree.Put(4, nil)
	tree.Put(5, nil)
	tree.Put(6, nil)
	tree.Put(7, nil)

	assertValidTree(t, tree, 7)
	assertValidTreeNode(t, tree.Root(), 1, 2, []int{4}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 1, 2, []int{2}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 1, 2, []int{6}, true)
	assertValidTreeNode(t, tree.Root().Children()[0].Children()[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root().Children()[0].Children()[1], 1, 0, []int{3}, true)
	assertValidTreeNode(t, tree.Root().Children()[1].Children()[0], 1, 0, []int{5}, true)
	assertValidTreeNode(t, tree.Root().Children()[1].Children()[1], 1, 0, []int{7}, true)

	tree.Delete(7)
	assertValidTree(t, tree, 6)
	assertValidTreeNode(t, tree.Root(), 2, 3, []int{2, 4}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 1, 0, []int{3}, true)
	assertValidTreeNode(t, tree.Root().Children()[2], 2, 0, []int{5, 6}, true)
}

func TestBTreeDelete7(t *testing.T) {
	// root height reduction after a series of underflows on left side
	// use simulator: https://www.cs.usfca.edu/~galles/visualization/BTree.html
	tree := New[int, *struct{}](3)
	tree.Put(1, nil)
	tree.Put(2, nil)
	tree.Put(3, nil)
	tree.Put(4, nil)
	tree.Put(5, nil)
	tree.Put(6, nil)
	tree.Put(7, nil)

	assertValidTree(t, tree, 7)
	assertValidTreeNode(t, tree.Root(), 1, 2, []int{4}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 1, 2, []int{2}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 1, 2, []int{6}, true)
	assertValidTreeNode(t, tree.Root().Children()[0].Children()[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root().Children()[0].Children()[1], 1, 0, []int{3}, true)
	assertValidTreeNode(t, tree.Root().Children()[1].Children()[0], 1, 0, []int{5}, true)
	assertValidTreeNode(t, tree.Root().Children()[1].Children()[1], 1, 0, []int{7}, true)

	tree.Delete(1) // series of underflows
	assertValidTree(t, tree, 6)
	assertValidTreeNode(t, tree.Root(), 2, 3, []int{4, 6}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 2, 0, []int{2, 3}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 1, 0, []int{5}, true)
	assertValidTreeNode(t, tree.Root().Children()[2], 1, 0, []int{7}, true)

	// clear all remaining
	tree.Delete(2)
	assertValidTree(t, tree, 5)
	assertValidTreeNode(t, tree.Root(), 2, 3, []int{4, 6}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 1, 0, []int{3}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 1, 0, []int{5}, true)
	assertValidTreeNode(t, tree.Root().Children()[2], 1, 0, []int{7}, true)

	tree.Delete(3)
	assertValidTree(t, tree, 4)
	assertValidTreeNode(t, tree.Root(), 1, 2, []int{6}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 2, 0, []int{4, 5}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 1, 0, []int{7}, true)

	tree.Delete(4)
	assertValidTree(t, tree, 3)
	assertValidTreeNode(t, tree.Root(), 1, 2, []int{6}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 1, 0, []int{5}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 1, 0, []int{7}, true)

	tree.Delete(5)
	assertValidTree(t, tree, 2)
	assertValidTreeNode(t, tree.Root(), 2, 0, []int{6, 7}, false)

	tree.Delete(6)
	assertValidTree(t, tree, 1)
	assertValidTreeNode(t, tree.Root(), 1, 0, []int{7}, false)

	tree.Delete(7)
	assertValidTree(t, tree, 0)
}

func TestBTreeDelete8(t *testing.T) {
	// use simulator: https://www.cs.usfca.edu/~galles/visualization/BTree.html
	tree := New[int, *struct{}](3)
	tree.Put(1, nil)
	tree.Put(2, nil)
	tree.Put(3, nil)
	tree.Put(4, nil)
	tree.Put(5, nil)
	tree.Put(6, nil)
	tree.Put(7, nil)
	tree.Put(8, nil)
	tree.Put(9, nil)

	assertValidTree(t, tree, 9)
	assertValidTreeNode(t, tree.Root(), 1, 2, []int{4}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 1, 2, []int{2}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 2, 3, []int{6, 8}, true)
	assertValidTreeNode(t, tree.Root().Children()[0].Children()[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root().Children()[0].Children()[1], 1, 0, []int{3}, true)
	assertValidTreeNode(t, tree.Root().Children()[1].Children()[0], 1, 0, []int{5}, true)
	assertValidTreeNode(t, tree.Root().Children()[1].Children()[1], 1, 0, []int{7}, true)
	assertValidTreeNode(t, tree.Root().Children()[1].Children()[2], 1, 0, []int{9}, true)

	tree.Delete(1)
	assertValidTree(t, tree, 8)
	assertValidTreeNode(t, tree.Root(), 1, 2, []int{6}, false)
	assertValidTreeNode(t, tree.Root().Children()[0], 1, 2, []int{4}, true)
	assertValidTreeNode(t, tree.Root().Children()[1], 1, 2, []int{8}, true)
	assertValidTreeNode(t, tree.Root().Children()[0].Children()[0], 2, 0, []int{2, 3}, true)
	assertValidTreeNode(t, tree.Root().Children()[0].Children()[1], 1, 0, []int{5}, true)
	assertValidTreeNode(t, tree.Root().Children()[1].Children()[0], 1, 0, []int{7}, true)
	assertValidTreeNode(t, tree.Root().Children()[1].Children()[1], 1, 0, []int{9}, true)
}

func TestBTreeDelete9(t *testing.T) {
	const max = 1000

	orders := []int{3, 4, 5, 6, 7, 8, 9, 10, 20, 100, 500, 1000, 5000, 10000}
	for _, order := range orders {
		tree := New[int, int](order)

		{
			for i := 1; i <= max; i++ {
				tree.Put(i, i)
			}

			assertValidTree(t, tree, max)

			for i := 1; i <= max; i++ {
				if _, found := tree.Get(i); !found {
					t.Errorf("Not found %v", i)
				}
			}

			for i := 1; i <= max; i++ {
				tree.Delete(i)
			}

			assertValidTree(t, tree, 0)
		}

		{
			for i := max; i > 0; i-- {
				tree.Put(i, i)
			}

			assertValidTree(t, tree, max)

			for i := max; i > 0; i-- {
				if _, found := tree.Get(i); !found {
					t.Errorf("Not found %v", i)
				}
			}

			for i := max; i > 0; i-- {
				tree.Delete(i)
			}

			assertValidTree(t, tree, 0)
		}
	}
}

func TestBTreeHeight(t *testing.T) {
	tree := New[int, int](3)
	if actualValue, expectedValue := tree.Height(), 0; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	tree.Put(1, 0)

	if actualValue, expectedValue := tree.Height(), 1; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	tree.Put(2, 1)

	if actualValue, expectedValue := tree.Height(), 1; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	tree.Put(3, 2)

	if actualValue, expectedValue := tree.Height(), 2; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	tree.Put(4, 2)

	if actualValue, expectedValue := tree.Height(), 2; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	tree.Put(5, 2)

	if actualValue, expectedValue := tree.Height(), 2; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	tree.Put(6, 2)

	if actualValue, expectedValue := tree.Height(), 2; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	tree.Put(7, 2)

	if actualValue, expectedValue := tree.Height(), 3; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	tree.Delete(1)
	tree.Delete(2)
	tree.Delete(3)
	tree.Delete(4)
	tree.Delete(5)
	tree.Delete(6)
	tree.Delete(7)

	if actualValue, expectedValue := tree.Height(), 0; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBTreeLeftAndRight(t *testing.T) {
	tree := New[int, string](3)

	if actualValue := tree.GetBeginNode(); actualValue != nil {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}

	if actualValue := tree.GetEndNode(); actualValue != nil {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}

	tree.Put(1, "a")
	tree.Put(5, "e")
	tree.Put(6, "f")
	tree.Put(7, "g")
	tree.Put(3, "c")
	tree.Put(4, "d")
	tree.Put(1, "x") // overwrite
	tree.Put(2, "b")

	if actualValue, expectedValue := tree.GetBeginNode().Entries()[0].Key(), 1; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, expectedValue := tree.GetBeginNode().Entries()[0].Value(), "x"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, expectedValue := tree.GetEndNode().Entries()[len(tree.GetEndNode().Entries())-1].Key(), 7; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, expectedValue := tree.GetEndNode().Entries()[len(tree.GetEndNode().Entries())-1].Value(), "g"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBTreeIteratorValuesAndKeys(t *testing.T) {
	tree := New[int, string](4)
	tree.Put(4, "d")
	tree.Put(5, "e")
	tree.Put(6, "f")
	tree.Put(3, "c")
	tree.Put(1, "a")
	tree.Put(7, "g")
	tree.Put(2, "b")
	tree.Put(1, "x") // override

	if actualValue, expectedValue := tree.Keys(), []int{1, 2, 3, 4, 5, 6, 7}; !slices.Equal(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, expectedValue := tree.Values(), []string{"x", "b", "c", "d", "e", "f", "g"}; !slices.Equal(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue := tree.Len(); actualValue != 7 {
		t.Errorf("Got %v expected %v", actualValue, 7)
	}
}

func TestBTreeIteratorNextOnEmpty(t *testing.T) {
	tree := New[int, string](3)

	for range tree.Iter() {
		t.Errorf("Shouldn't iterate on empty tree")
	}
}

func TestBTreeIteratorPrevOnEmpty(t *testing.T) {
	tree := New[int, string](3)

	for range tree.RIter() {
		t.Errorf("Shouldn't iterate on empty tree")
	}
}

func TestBTreeIterator1Next(t *testing.T) {
	tree := New[int, string](3)
	tree.Put(5, "e")
	tree.Put(6, "f")
	tree.Put(7, "g")
	tree.Put(3, "c")
	tree.Put(4, "d")
	tree.Put(1, "x")
	tree.Put(2, "b")
	tree.Put(1, "a") //overwrite

	count := 0

	for k := range tree.Iter() {
		count++

		if actualValue, expectedValue := k, count; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
	}

	if actualValue, expectedValue := count, tree.Len(); actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBTreeIterator1Prev(t *testing.T) {
	tree := New[int, string](3)
	tree.Put(5, "e")
	tree.Put(6, "f")
	tree.Put(7, "g")
	tree.Put(3, "c")
	tree.Put(4, "d")
	tree.Put(1, "x")
	tree.Put(2, "b")
	tree.Put(1, "a") //overwrite

	for range tree.RIter() {
	}

	countDown := tree.len

	for k := range tree.RIter() {
		if actualValue, expectedValue := k, countDown; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}

		countDown--
	}

	if actualValue, expectedValue := countDown, 0; actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBTreeIterator2Next(t *testing.T) {
	tree := New[int, string](3)
	tree.Put(3, "c")
	tree.Put(1, "a")
	tree.Put(2, "b")

	count := 0

	for k := range tree.Iter() {
		count++

		if actualValue, expectedValue := k, count; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
	}

	if actualValue, expectedValue := count, tree.Len(); actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBTreeIterator2Prev(t *testing.T) {
	tree := New[int, string](3)
	tree.Put(3, "c")
	tree.Put(1, "a")
	tree.Put(2, "b")

	for range tree.Iter() {
	}

	countDown := tree.len

	for k := range tree.RIter() {
		if actualValue, expectedValue := k, countDown; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}

		countDown--
	}

	if actualValue, expectedValue := countDown, 0; actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBTreeIterator3Next(t *testing.T) {
	tree := New[int, string](3)
	tree.Put(1, "a")

	count := 0

	for k := range tree.Iter() {
		count++

		if actualValue, expectedValue := k, count; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
	}

	if actualValue, expectedValue := count, tree.Len(); actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBTreeIterator3Prev(t *testing.T) {
	tree := New[int, string](3)
	tree.Put(1, "a")

	for range tree.Iter() {
	}

	countDown := tree.len

	for k := range tree.RIter() {
		if actualValue, expectedValue := k, countDown; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}

		countDown--
	}

	if actualValue, expectedValue := countDown, 0; actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBTreeIterator4Next(t *testing.T) {
	tree := New[int, int](3)
	tree.Put(13, 5)
	tree.Put(8, 3)
	tree.Put(17, 7)
	tree.Put(1, 1)
	tree.Put(11, 4)
	tree.Put(15, 6)
	tree.Put(25, 9)
	tree.Put(6, 2)
	tree.Put(22, 8)
	tree.Put(27, 10)

	count := 0

	for _, v := range tree.Iter() {
		count++

		if actualValue, expectedValue := v, count; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
	}

	if actualValue, expectedValue := count, tree.Len(); actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBTreeIterator4Prev(t *testing.T) {
	tree := New[int, int](3)
	tree.Put(13, 5)
	tree.Put(8, 3)
	tree.Put(17, 7)
	tree.Put(1, 1)
	tree.Put(11, 4)
	tree.Put(15, 6)
	tree.Put(25, 9)
	tree.Put(6, 2)
	tree.Put(22, 8)
	tree.Put(27, 10)
	count := tree.Len()

	for _, v := range tree.RIter() {
		if actualValue, expectedValue := v, count; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}

		count--
	}

	if actualValue, expectedValue := count, 0; actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBTreeSearch(t *testing.T) {
	{
		tree := New[int, int](3)
		tree.root = &Node[int, int]{entries: []*entry[int, int]{}, children: make([]*Node[int, int], 0)}
		tests := [][]interface{}{
			{0, 0, false},
		}

		for _, test := range tests {
			index, found := tree.search(tree.Root(), test[0].(int))
			if actualValue, expectedValue := index, test[1]; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}

			if actualValue, expectedValue := found, test[2]; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}
	}
	{
		tree := New[int, int](3)
		tree.root = &Node[int, int]{entries: []*entry[int, int]{{2, 0}, {4, 1}, {6, 2}}, children: []*Node[int, int]{}}
		tests := [][]interface{}{
			{0, 0, false},
			{1, 0, false},
			{2, 0, true},
			{3, 1, false},
			{4, 1, true},
			{5, 2, false},
			{6, 2, true},
			{7, 3, false},
		}

		for _, test := range tests {
			index, found := tree.search(tree.Root(), test[0].(int))
			if actualValue, expectedValue := index, test[1]; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}

			if actualValue, expectedValue := found, test[2]; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}
	}
}

func TestBTreeSerialization(t *testing.T) {
	tree := New[string, string](3)
	tree.Put("c", "3")
	tree.Put("b", "2")
	tree.Put("a", "1")

	var err error

	assert := func() {
		if actualValue, expectedValue := tree.Len(), 3; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}

		if actualValue, expectedValue := tree.Keys(), []string{"a", "b", "c"}; !slices.Equal(actualValue, expectedValue) {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}

		if actualValue, expectedValue := tree.Values(), []string{"1", "2", "3"}; !slices.Equal(actualValue, expectedValue) {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}

		if err != nil {
			t.Errorf("Got error %v", err)
		}
	}

	assert()

	bytes, err := tree.MarshalJSON()

	assert()

	err = tree.UnmarshalJSON(bytes)

	assert()

	_, err = json.Marshal([]interface{}{"a", "b", "c", tree})
	if err != nil {
		t.Errorf("Got error %v", err)
	}

	intTree := New[string, int](3)

	err = json.Unmarshal([]byte(`{"a":1,"b":2}`), intTree)
	if err != nil {
		t.Errorf("Got error %v", err)
	}

	if actualValue, expectedValue := intTree.Len(), 2; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, expectedValue := intTree.Keys(), []string{"a", "b"}; !slices.Equal(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, expectedValue := intTree.Values(), []int{1, 2}; !slices.Equal(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBTreeString(t *testing.T) {
	c := New[string, int](3)
	c.Put("a", 1)

	if !strings.HasPrefix(c.String(), "BTree") {
		t.Errorf("String should start with container name")
	}
}
