package btree

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"testing"
)

// TestBTreeGet tests the Get method of the BTree.
func TestBTreeGet(t *testing.T) {
	t.Parallel()

	// Use t.Run to organize test cases
	t.Run("basic operations", func(t *testing.T) {
		t.Parallel()

		tree := New[int, string](3)

		// Initialize test data using a map for clarity
		testData := map[int]string{
			1: "a",
			2: "b",
			3: "c",
			4: "d",
			5: "e",
			6: "f",
			7: "g",
		}

		// Populate the tree with test data
		for k, v := range testData {
			tree.Put(k, v)
		}

		// Define test cases with a struct for better type safety
		tests := []struct {
			key       int
			wantVal   string
			wantFound bool
		}{
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

		// Run each test case
		for _, tt := range tests {
			t.Run(fmt.Sprintf("key=%d", tt.key), func(t *testing.T) {
				t.Parallel()

				gotVal, gotFound := tree.Get(tt.key)
				if gotVal != tt.wantVal || gotFound != tt.wantFound {
					t.Errorf("Get(%d) = (%q, %v), want (%q, %v)",
						tt.key, gotVal, gotFound, tt.wantVal, tt.wantFound)
				}
			})
		}
	})
}

// TestBTreeGet2 tests the Get method of the BTree with various key insertions.
func TestBTreeGet2(t *testing.T) {
	t.Parallel() // Enable parallel execution for this test

	// Use t.Run to group test cases
	t.Run("mixed order insertions", func(t *testing.T) {
		t.Parallel()

		tree := New[int, string](3)

		// Initialize test data using a map for better readability
		testData := map[int]string{
			7:  "g",
			9:  "i",
			10: "j",
			6:  "f",
			3:  "c",
			4:  "d",
			5:  "e",
			8:  "h",
			2:  "b",
			1:  "a",
		}

		// Populate the tree with test data
		for k, v := range testData {
			tree.Put(k, v)
		}

		// Define test cases with a struct for type safety
		tests := []struct {
			key       int
			wantVal   string
			wantFound bool
		}{
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

		// Run each test case in parallel
		for _, tt := range tests {
			// Capture range variable for parallel execution
			t.Run(fmt.Sprintf("key=%d", tt.key), func(t *testing.T) {
				t.Parallel()

				gotVal, gotFound := tree.Get(tt.key)
				if gotVal != tt.wantVal || gotFound != tt.wantFound {
					t.Errorf("Get(%d) = (%q, %v), want (%q, %v)",
						tt.key, gotVal, gotFound, tt.wantVal, tt.wantFound)
				}
			})
		}
	})
}

// TestBTreeGet3 tests the Size and GetNode methods of the BTree.
func TestBTreeGet3(t *testing.T) {
	t.Parallel() // Enable parallel execution for this test

	// Use t.Run to organize test scenarios
	t.Run("tree operations", func(t *testing.T) {

		tree := New[int, string](3)

		// Test initial state
		t.Run("initial size", func(t *testing.T) {

			if got := tree.Size(); got != 0 {
				t.Errorf("Size() = %d, want 0", got)
			}
		})

		t.Run("initial node size", func(t *testing.T) {

			if got := tree.GetNode(2).Size(); got != 0 {
				t.Errorf("GetNode(2).Size() = %d, want 0", got)
			}
		})

		// Populate the tree with test data
		insertions := []struct {
			key   int
			value string
		}{
			{1, "x"}, // Initial insertion
			{2, "b"}, // Add second key
			{1, "a"}, // Replace value for key 1
			{3, "c"},
			{4, "d"},
			{5, "e"},
			{6, "f"},
			{7, "g"},
		}
		for _, ins := range insertions {
			tree.Put(ins.key, ins.value)
		}

		// Define test cases for size checks
		sizeTests := []struct {
			name     string
			getSize  func() int
			wantSize int
		}{
			{"total size", tree.Size, 7},
			{"node 2 size", func() int { return tree.GetNode(2).Size() }, 3},
			{"node 4 size", func() int { return tree.GetNode(4).Size() }, 7},
			{"node 8 size", func() int { return tree.GetNode(8).Size() }, 0},
		}

		// Run size checks in parallel
		for _, tt := range sizeTests {
			// Capture range variable for parallel execution
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				if got := tt.getSize(); got != tt.wantSize {
					t.Errorf("%s = %d, want %d", tt.name, got, tt.wantSize)
				}
			})
		}
	})
}

func TestBTreePut1(t *testing.T) {
	// https://upload.wikimedia.org/wikipedia/commons/3/33/B_tree_insertion_example.png
	tree := New[int, int](3)
	assertValidTree(t, tree, 0)

	tree.Put(1, 0)
	assertValidTree(t, tree, 1)
	assertValidTreeNode(t, tree.Root, 1, 0, []int{1}, false)

	tree.Put(2, 1)
	assertValidTree(t, tree, 2)
	assertValidTreeNode(t, tree.Root, 2, 0, []int{1, 2}, false)

	tree.Put(3, 2)
	assertValidTree(t, tree, 3)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{2}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{3}, true)

	tree.Put(4, 2)
	assertValidTree(t, tree, 4)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{2}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 2, 0, []int{3, 4}, true)

	tree.Put(5, 2)
	assertValidTree(t, tree, 5)
	assertValidTreeNode(t, tree.Root, 2, 3, []int{2, 4}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{3}, true)
	assertValidTreeNode(t, tree.Root.Children[2], 1, 0, []int{5}, true)

	tree.Put(6, 2)
	assertValidTree(t, tree, 6)
	assertValidTreeNode(t, tree.Root, 2, 3, []int{2, 4}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{3}, true)
	assertValidTreeNode(t, tree.Root.Children[2], 2, 0, []int{5, 6}, true)

	tree.Put(7, 2)
	assertValidTree(t, tree, 7)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{4}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 2, []int{2}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 2, []int{6}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[1], 1, 0, []int{3}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[0], 1, 0, []int{5}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[1], 1, 0, []int{7}, true)
}

func TestBTreePut2(t *testing.T) {
	tree := New[int, int](4)
	assertValidTree(t, tree, 0)

	tree.Put(0, 0)
	assertValidTree(t, tree, 1)
	assertValidTreeNode(t, tree.Root, 1, 0, []int{0}, false)

	tree.Put(2, 2)
	assertValidTree(t, tree, 2)
	assertValidTreeNode(t, tree.Root, 2, 0, []int{0, 2}, false)

	tree.Put(1, 1)
	assertValidTree(t, tree, 3)
	assertValidTreeNode(t, tree.Root, 3, 0, []int{0, 1, 2}, false)

	tree.Put(1, 1)
	assertValidTree(t, tree, 3)
	assertValidTreeNode(t, tree.Root, 3, 0, []int{0, 1, 2}, false)

	tree.Put(3, 3)
	assertValidTree(t, tree, 4)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{1}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{0}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 2, 0, []int{2, 3}, true)

	tree.Put(4, 4)
	assertValidTree(t, tree, 5)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{1}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{0}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 3, 0, []int{2, 3, 4}, true)

	tree.Put(5, 5)
	assertValidTree(t, tree, 6)
	assertValidTreeNode(t, tree.Root, 2, 3, []int{1, 3}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{0}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{2}, true)
	assertValidTreeNode(t, tree.Root.Children[2], 2, 0, []int{4, 5}, true)
}

func TestBTreePut3(t *testing.T) {
	// http://www.geeksforgeeks.org/b-tree-set-1-insert-2/
	tree := New[int, int](6)
	assertValidTree(t, tree, 0)

	tree.Put(10, 0)
	assertValidTree(t, tree, 1)
	assertValidTreeNode(t, tree.Root, 1, 0, []int{10}, false)

	tree.Put(20, 1)
	assertValidTree(t, tree, 2)
	assertValidTreeNode(t, tree.Root, 2, 0, []int{10, 20}, false)

	tree.Put(30, 2)
	assertValidTree(t, tree, 3)
	assertValidTreeNode(t, tree.Root, 3, 0, []int{10, 20, 30}, false)

	tree.Put(40, 3)
	assertValidTree(t, tree, 4)
	assertValidTreeNode(t, tree.Root, 4, 0, []int{10, 20, 30, 40}, false)

	tree.Put(50, 4)
	assertValidTree(t, tree, 5)
	assertValidTreeNode(t, tree.Root, 5, 0, []int{10, 20, 30, 40, 50}, false)

	tree.Put(60, 5)
	assertValidTree(t, tree, 6)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{30}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 2, 0, []int{10, 20}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 3, 0, []int{40, 50, 60}, true)

	tree.Put(70, 6)
	assertValidTree(t, tree, 7)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{30}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 2, 0, []int{10, 20}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 4, 0, []int{40, 50, 60, 70}, true)

	tree.Put(80, 7)
	assertValidTree(t, tree, 8)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{30}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 2, 0, []int{10, 20}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 5, 0, []int{40, 50, 60, 70, 80}, true)

	tree.Put(90, 8)
	assertValidTree(t, tree, 9)
	assertValidTreeNode(t, tree.Root, 2, 3, []int{30, 60}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 2, 0, []int{10, 20}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 2, 0, []int{40, 50}, true)
	assertValidTreeNode(t, tree.Root.Children[2], 3, 0, []int{70, 80, 90}, true)
}

func TestBTreePut4(t *testing.T) {
	tree := New[int, *struct{}](3)
	assertValidTree(t, tree, 0)

	tree.Put(6, nil)
	assertValidTree(t, tree, 1)
	assertValidTreeNode(t, tree.Root, 1, 0, []int{6}, false)

	tree.Put(5, nil)
	assertValidTree(t, tree, 2)
	assertValidTreeNode(t, tree.Root, 2, 0, []int{5, 6}, false)

	tree.Put(4, nil)
	assertValidTree(t, tree, 3)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{5}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{4}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{6}, true)

	tree.Put(3, nil)
	assertValidTree(t, tree, 4)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{5}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 2, 0, []int{3, 4}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{6}, true)

	tree.Put(2, nil)
	assertValidTree(t, tree, 5)
	assertValidTreeNode(t, tree.Root, 2, 3, []int{3, 5}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{2}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{4}, true)
	assertValidTreeNode(t, tree.Root.Children[2], 1, 0, []int{6}, true)

	tree.Put(1, nil)
	assertValidTree(t, tree, 6)
	assertValidTreeNode(t, tree.Root, 2, 3, []int{3, 5}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 2, 0, []int{1, 2}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{4}, true)
	assertValidTreeNode(t, tree.Root.Children[2], 1, 0, []int{6}, true)

	tree.Put(0, nil)
	assertValidTree(t, tree, 7)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{3}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 2, []int{1}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 2, []int{5}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[0], 1, 0, []int{0}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[1], 1, 0, []int{2}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[0], 1, 0, []int{4}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[1], 1, 0, []int{6}, true)

	tree.Put(-1, nil)
	assertValidTree(t, tree, 8)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{3}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 2, []int{1}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 2, []int{5}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[0], 2, 0, []int{-1, 0}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[1], 1, 0, []int{2}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[0], 1, 0, []int{4}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[1], 1, 0, []int{6}, true)

	tree.Put(-2, nil)
	assertValidTree(t, tree, 9)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{3}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 2, 3, []int{-1, 1}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 2, []int{5}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[0], 1, 0, []int{-2}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[1], 1, 0, []int{0}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[2], 1, 0, []int{2}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[0], 1, 0, []int{4}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[1], 1, 0, []int{6}, true)

	tree.Put(-3, nil)
	assertValidTree(t, tree, 10)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{3}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 2, 3, []int{-1, 1}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 2, []int{5}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[0], 2, 0, []int{-3, -2}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[1], 1, 0, []int{0}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[2], 1, 0, []int{2}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[0], 1, 0, []int{4}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[1], 1, 0, []int{6}, true)

	tree.Put(-4, nil)
	assertValidTree(t, tree, 11)
	assertValidTreeNode(t, tree.Root, 2, 3, []int{-1, 3}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 2, []int{-3}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 2, []int{1}, true)
	assertValidTreeNode(t, tree.Root.Children[2], 1, 2, []int{5}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[0], 1, 0, []int{-4}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[1], 1, 0, []int{-2}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[0], 1, 0, []int{0}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[1], 1, 0, []int{2}, true)
	assertValidTreeNode(t, tree.Root.Children[2].Children[0], 1, 0, []int{4}, true)
	assertValidTreeNode(t, tree.Root.Children[2].Children[1], 1, 0, []int{6}, true)
}

func TestBTreeRemove1(t *testing.T) {
	// empty
	tree := New[int, int](3)
	tree.Remove(1)
	assertValidTree(t, tree, 0)
}

func TestBTreeRemove2(t *testing.T) {
	// leaf node (no underflow)
	tree := New[int, *struct{}](3)
	tree.Put(1, nil)
	tree.Put(2, nil)

	tree.Remove(1)
	assertValidTree(t, tree, 1)
	assertValidTreeNode(t, tree.Root, 1, 0, []int{2}, false)

	tree.Remove(2)
	assertValidTree(t, tree, 0)
}

func TestBTreeRemove3(t *testing.T) {
	// merge with right (underflow)
	{
		tree := New[int, *struct{}](3)
		tree.Put(1, nil)
		tree.Put(2, nil)
		tree.Put(3, nil)

		tree.Remove(1)
		assertValidTree(t, tree, 2)
		assertValidTreeNode(t, tree.Root, 2, 0, []int{2, 3}, false)
	}
	// merge with left (underflow)
	{
		tree := New[int, *struct{}](3)
		tree.Put(1, nil)
		tree.Put(2, nil)
		tree.Put(3, nil)

		tree.Remove(3)
		assertValidTree(t, tree, 2)
		assertValidTreeNode(t, tree.Root, 2, 0, []int{1, 2}, false)
	}
}

func TestBTreeRemove4(t *testing.T) {
	// rotate left (underflow)
	tree := New[int, *struct{}](3)
	tree.Put(1, nil)
	tree.Put(2, nil)
	tree.Put(3, nil)
	tree.Put(4, nil)

	assertValidTree(t, tree, 4)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{2}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 2, 0, []int{3, 4}, true)

	tree.Remove(1)
	assertValidTree(t, tree, 3)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{3}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{2}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{4}, true)
}

func TestBTreeRemove5(t *testing.T) {
	// rotate right (underflow)
	tree := New[int, *struct{}](3)
	tree.Put(1, nil)
	tree.Put(2, nil)
	tree.Put(3, nil)
	tree.Put(0, nil)

	assertValidTree(t, tree, 4)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{2}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 2, 0, []int{0, 1}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{3}, true)

	tree.Remove(3)
	assertValidTree(t, tree, 3)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{1}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{0}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{2}, true)
}

func TestBTreeRemove6(t *testing.T) {
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
	assertValidTreeNode(t, tree.Root, 1, 2, []int{4}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 2, []int{2}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 2, []int{6}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[1], 1, 0, []int{3}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[0], 1, 0, []int{5}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[1], 1, 0, []int{7}, true)

	tree.Remove(7)
	assertValidTree(t, tree, 6)
	assertValidTreeNode(t, tree.Root, 2, 3, []int{2, 4}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{3}, true)
	assertValidTreeNode(t, tree.Root.Children[2], 2, 0, []int{5, 6}, true)
}

func TestBTreeRemove7(t *testing.T) {
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
	assertValidTreeNode(t, tree.Root, 1, 2, []int{4}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 2, []int{2}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 2, []int{6}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[1], 1, 0, []int{3}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[0], 1, 0, []int{5}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[1], 1, 0, []int{7}, true)

	tree.Remove(1) // series of underflows
	assertValidTree(t, tree, 6)
	assertValidTreeNode(t, tree.Root, 2, 3, []int{4, 6}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 2, 0, []int{2, 3}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{5}, true)
	assertValidTreeNode(t, tree.Root.Children[2], 1, 0, []int{7}, true)

	// clear all remaining
	tree.Remove(2)
	assertValidTree(t, tree, 5)
	assertValidTreeNode(t, tree.Root, 2, 3, []int{4, 6}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{3}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{5}, true)
	assertValidTreeNode(t, tree.Root.Children[2], 1, 0, []int{7}, true)

	tree.Remove(3)
	assertValidTree(t, tree, 4)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{6}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 2, 0, []int{4, 5}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{7}, true)

	tree.Remove(4)
	assertValidTree(t, tree, 3)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{6}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{5}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{7}, true)

	tree.Remove(5)
	assertValidTree(t, tree, 2)
	assertValidTreeNode(t, tree.Root, 2, 0, []int{6, 7}, false)

	tree.Remove(6)
	assertValidTree(t, tree, 1)
	assertValidTreeNode(t, tree.Root, 1, 0, []int{7}, false)

	tree.Remove(7)
	assertValidTree(t, tree, 0)
}

func TestBTreeRemove8(t *testing.T) {
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
	assertValidTreeNode(t, tree.Root, 1, 2, []int{4}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 2, []int{2}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 2, 3, []int{6, 8}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[1], 1, 0, []int{3}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[0], 1, 0, []int{5}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[1], 1, 0, []int{7}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[2], 1, 0, []int{9}, true)

	tree.Remove(1)
	assertValidTree(t, tree, 8)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{6}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 2, []int{4}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 2, []int{8}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[0], 2, 0, []int{2, 3}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[1], 1, 0, []int{5}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[0], 1, 0, []int{7}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[1], 1, 0, []int{9}, true)
}

func TestBTreeRemove9(t *testing.T) {
	const maxSize = 1000

	orders := []int{3, 4, 5, 6, 7, 8, 9, 10, 20, 100, 500, 1000, 5000, 10000}
	for _, order := range orders {
		tree := New[int, int](order)

		{
			for i := 1; i <= maxSize; i++ {
				tree.Put(i, i)
			}

			assertValidTree(t, tree, maxSize)

			for i := 1; i <= maxSize; i++ {
				if _, found := tree.Get(i); !found {
					t.Errorf("Not found %v", i)
				}
			}

			for i := 1; i <= maxSize; i++ {
				tree.Remove(i)
			}

			assertValidTree(t, tree, 0)
		}

		{
			for i := maxSize; i > 0; i-- {
				tree.Put(i, i)
			}

			assertValidTree(t, tree, maxSize)

			for i := maxSize; i > 0; i-- {
				if _, found := tree.Get(i); !found {
					t.Errorf("Not found %v", i)
				}
			}

			for i := maxSize; i > 0; i-- {
				tree.Remove(i)
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

	tree.Remove(1)
	tree.Remove(2)
	tree.Remove(3)
	tree.Remove(4)
	tree.Remove(5)
	tree.Remove(6)
	tree.Remove(7)

	if actualValue, expectedValue := tree.Height(), 0; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBTreeLeftAndRight(t *testing.T) {
	tree := New[int, string](3)

	if actualValue := tree.Left(); actualValue != nil {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}

	if actualValue := tree.Right(); actualValue != nil {
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

	if actualValue, expectedValue := tree.LeftKey(), 1; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, expectedValue := tree.LeftValue(), "x"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, expectedValue := tree.RightKey(), 7; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, expectedValue := tree.RightValue(), "g"; actualValue != expectedValue {
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

	if actualValue := tree.Size(); actualValue != 7 {
		t.Errorf("Got %v expected %v", actualValue, 7)
	}
}

func TestBTreeIteratorNextOnEmpty(t *testing.T) {
	tree := New[int, string](3)

	it := tree.Iterator()
	for it.Next() {
		t.Errorf("Shouldn't iterate on empty tree")
	}
}

func TestBTreeIteratorPrevOnEmpty(t *testing.T) {
	tree := New[int, string](3)

	it := tree.Iterator()
	for it.Prev() {
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
	it := tree.Iterator()
	count := 0

	for it.Next() {
		count++

		key := it.Key()
		if actualValue, expectedValue := key, count; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
	}

	if actualValue, expectedValue := count, tree.Size(); actualValue != expectedValue {
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

	it := tree.Iterator()
	for it.Next() {
	}

	countDown := tree.size

	for it.Prev() {
		key := it.Key()
		if actualValue, expectedValue := key, countDown; actualValue != expectedValue {
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
	it := tree.Iterator()
	count := 0

	for it.Next() {
		count++

		key := it.Key()
		if actualValue, expectedValue := key, count; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
	}

	if actualValue, expectedValue := count, tree.Size(); actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBTreeIterator2Prev(t *testing.T) {
	tree := New[int, string](3)
	tree.Put(3, "c")
	tree.Put(1, "a")
	tree.Put(2, "b")

	it := tree.Iterator()
	for it.Next() {
	}

	countDown := tree.size

	for it.Prev() {
		key := it.Key()
		if actualValue, expectedValue := key, countDown; actualValue != expectedValue {
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
	it := tree.Iterator()
	count := 0

	for it.Next() {
		count++

		key := it.Key()
		if actualValue, expectedValue := key, count; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
	}

	if actualValue, expectedValue := count, tree.Size(); actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBTreeIterator3Prev(t *testing.T) {
	tree := New[int, string](3)
	tree.Put(1, "a")

	it := tree.Iterator()
	for it.Next() {
	}

	countDown := tree.size

	for it.Prev() {
		key := it.Key()
		if actualValue, expectedValue := key, countDown; actualValue != expectedValue {
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
	it := tree.Iterator()
	count := 0

	for it.Next() {
		count++

		value := it.Value()
		if actualValue, expectedValue := value, count; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
	}

	if actualValue, expectedValue := count, tree.Size(); actualValue != expectedValue {
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
	it := tree.Iterator()
	count := tree.Size()

	for it.Next() {
	}

	for it.Prev() {
		value := it.Value()
		if actualValue, expectedValue := value, count; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}

		count--
	}

	if actualValue, expectedValue := count, 0; actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBTreeIteratorBegin(t *testing.T) {
	tree := New[int, string](3)
	tree.Put(3, "c")
	tree.Put(1, "a")
	tree.Put(2, "b")
	it := tree.Iterator()

	if it.node != nil {
		t.Errorf("Got %v expected %v", it.node, nil)
	}

	it.Begin()

	if it.node != nil {
		t.Errorf("Got %v expected %v", it.node, nil)
	}

	for it.Next() {
	}

	it.Begin()

	if it.node != nil {
		t.Errorf("Got %v expected %v", it.node, nil)
	}

	it.Next()

	if key, value := it.Key(), it.Value(); key != 1 || value != "a" {
		t.Errorf("Got %v,%v expected %v,%v", key, value, 1, "a")
	}
}

func TestBTreeIteratorEnd(t *testing.T) {
	tree := New[int, string](3)
	it := tree.Iterator()

	if it.node != nil {
		t.Errorf("Got %v expected %v", it.node, nil)
	}

	it.End()

	if it.node != nil {
		t.Errorf("Got %v expected %v", it.node, nil)
	}

	tree.Put(3, "c")
	tree.Put(1, "a")
	tree.Put(2, "b")
	it.End()

	if it.node != nil {
		t.Errorf("Got %v expected %v", it.node, nil)
	}

	it.Prev()

	if key, value := it.Key(), it.Value(); key != 3 || value != "c" {
		t.Errorf("Got %v,%v expected %v,%v", key, value, 3, "c")
	}
}

func TestBTreeIteratorFirst(t *testing.T) {
	tree := New[int, string](3)
	tree.Put(3, "c")
	tree.Put(1, "a")
	tree.Put(2, "b")

	it := tree.Iterator()
	if actualValue, expectedValue := it.First(), true; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if key, value := it.Key(), it.Value(); key != 1 || value != "a" {
		t.Errorf("Got %v,%v expected %v,%v", key, value, 1, "a")
	}
}

func TestBTreeIteratorLast(t *testing.T) {
	tree := New[int, string](3)
	tree.Put(3, "c")
	tree.Put(1, "a")
	tree.Put(2, "b")

	it := tree.Iterator()
	if actualValue, expectedValue := it.Last(), true; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if key, value := it.Key(), it.Value(); key != 3 || value != "c" {
		t.Errorf("Got %v,%v expected %v,%v", key, value, 3, "c")
	}
}

func TestBTreeSearch(t *testing.T) {
	{
		tree := New[int, int](3)
		tree.Root = &Node[int, int]{Entries: []*Entry[int, int]{}, Children: make([]*Node[int, int], 0)}
		tests := [][]any{
			{0, 0, false},
		}

		for _, test := range tests {
			key := test[0].(int)

			index, found := tree.search(tree.Root, key)
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
		tree.Root = &Node[int, int]{Entries: []*Entry[int, int]{{2, 0}, {4, 1}, {6, 2}}, Children: []*Node[int, int]{}}
		tests := [][]any{
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
			key := test[0].(int)

			index, found := tree.search(tree.Root, key)
			if actualValue, expectedValue := index, test[1]; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}

			if actualValue, expectedValue := found, test[2]; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}
	}
}

func assertValidTree[K comparable, V any](t *testing.T, tree *Tree[K, V], expectedSize int) {
	t.Helper()

	if actualValue, expectedValue := tree.size, expectedSize; actualValue != expectedValue {
		t.Errorf("Got %v expected %v for tree size", actualValue, expectedValue)
	}
}

func assertValidTreeNode[K comparable, V any](t *testing.T, node *Node[K, V], expectedEntries int, expectedChildren int, keys []K, hasParent bool) {
	t.Helper()

	if actualValue, expectedValue := node.Parent != nil, hasParent; actualValue != expectedValue {
		t.Errorf("Got %v expected %v for hasParent", actualValue, expectedValue)
	}

	if actualValue, expectedValue := len(node.Entries), expectedEntries; actualValue != expectedValue {
		t.Errorf("Got %v expected %v for entries size", actualValue, expectedValue)
	}

	if actualValue, expectedValue := len(node.Children), expectedChildren; actualValue != expectedValue {
		t.Errorf("Got %v expected %v for children size", actualValue, expectedValue)
	}

	for i, key := range keys {
		if actualValue, expectedValue := node.Entries[i].Key, key; actualValue != expectedValue {
			t.Errorf("Got %v expected %v for key", actualValue, expectedValue)
		}
	}
}

func TestBTreeIteratorNextTo(t *testing.T) {
	// Sample seek function, i.e. string starting with "b"
	seek := func(_ int, value string) bool {
		return strings.HasSuffix(value, "b")
	}

	// NextTo (empty)
	{
		tree := New[int, string](3)

		it := tree.Iterator()
		for it.NextTo(seek) {
			t.Errorf("Shouldn't iterate on empty tree")
		}
	}

	// NextTo (not found)
	{
		tree := New[int, string](3)
		tree.Put(0, "xx")
		tree.Put(1, "yy")

		it := tree.Iterator()
		for it.NextTo(seek) {
			t.Errorf("Shouldn't iterate on empty tree")
		}
	}

	// NextTo (found)
	{
		tree := New[int, string](3)
		tree.Put(2, "cc")
		tree.Put(0, "aa")
		tree.Put(1, "bb")
		it := tree.Iterator()
		it.Begin()

		if !it.NextTo(seek) {
			t.Errorf("Shouldn't iterate on empty tree")
		}

		if index, value := it.Key(), it.Value(); index != 1 || value != "bb" {
			t.Errorf("Got %v,%v expected %v,%v", index, value, 1, "bb")
		}

		if !it.Next() {
			t.Errorf("Should go to first element")
		}

		if index, value := it.Key(), it.Value(); index != 2 || value != "cc" {
			t.Errorf("Got %v,%v expected %v,%v", index, value, 2, "cc")
		}

		if it.Next() {
			t.Errorf("Should not go past last element")
		}
	}
}

func TestBTreeIteratorPrevTo(t *testing.T) {
	// Sample seek function, i.e. string starting with "b"
	seek := func(_ int, value string) bool {
		return strings.HasSuffix(value, "b")
	}

	// PrevTo (empty)
	{
		tree := New[int, string](3)
		it := tree.Iterator()
		it.End()

		for it.PrevTo(seek) {
			t.Errorf("Shouldn't iterate on empty tree")
		}
	}

	// PrevTo (not found)
	{
		tree := New[int, string](3)
		tree.Put(0, "xx")
		tree.Put(1, "yy")
		it := tree.Iterator()
		it.End()

		for it.PrevTo(seek) {
			t.Errorf("Shouldn't iterate on empty tree")
		}
	}

	// PrevTo (found)
	{
		tree := New[int, string](3)
		tree.Put(2, "cc")
		tree.Put(0, "aa")
		tree.Put(1, "bb")
		it := tree.Iterator()
		it.End()

		if !it.PrevTo(seek) {
			t.Errorf("Shouldn't iterate on empty tree")
		}

		if index, value := it.Key(), it.Value(); index != 1 || value != "bb" {
			t.Errorf("Got %v,%v expected %v,%v", index, value, 1, "bb")
		}

		if !it.Prev() {
			t.Errorf("Should go to first element")
		}

		if index, value := it.Key(), it.Value(); index != 0 || value != "aa" {
			t.Errorf("Got %v,%v expected %v,%v", index, value, 0, "aa")
		}

		if it.Prev() {
			t.Errorf("Should not go before first element")
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
		if actualValue, expectedValue := tree.Size(), 3; actualValue != expectedValue {
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

	bytes, err := tree.ToJSON()

	assert()

	err = tree.FromJSON(bytes)

	assert()

	_, err = json.Marshal([]any{"a", "b", "c", tree})
	if err != nil {
		t.Errorf("Got error %v", err)
	}

	intTree := New[string, int](3)

	err = json.Unmarshal([]byte(`{"a":1,"b":2}`), intTree)
	if err != nil {
		t.Errorf("Got error %v", err)
	}

	if actualValue, expectedValue := intTree.Size(), 2; actualValue != expectedValue {
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

func benchmarkGet(b *testing.B, tree *Tree[int, struct{}], size int) {
	b.Helper()

	for b.Loop() {
		for n := range size {
			tree.Get(n)
		}
	}
}

func benchmarkPut(b *testing.B, tree *Tree[int, struct{}], size int) {
	b.Helper()

	for b.Loop() {
		for n := range size {
			tree.Put(n, struct{}{})
		}
	}
}

func benchmarkRemove(b *testing.B, tree *Tree[int, struct{}], size int) {
	b.Helper()

	for b.Loop() {
		for n := range size {
			tree.Remove(n)
		}
	}
}

func BenchmarkBTreeGet100(b *testing.B) {
	b.StopTimer()

	size := 100
	tree := New[int, struct{}](128)

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkGet(b, tree, size)
}

func BenchmarkBTreeGet1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	tree := New[int, struct{}](128)

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkGet(b, tree, size)
}

func BenchmarkBTreeGet10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	tree := New[int, struct{}](128)

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkGet(b, tree, size)
}

func BenchmarkBTreeGet100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	tree := New[int, struct{}](128)

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkGet(b, tree, size)
}

func BenchmarkBTreePut100(b *testing.B) {
	b.StopTimer()

	size := 100
	tree := New[int, struct{}](128)

	b.StartTimer()
	benchmarkPut(b, tree, size)
}

func BenchmarkBTreePut1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	tree := New[int, struct{}](128)

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkPut(b, tree, size)
}

func BenchmarkBTreePut10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	tree := New[int, struct{}](128)

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkPut(b, tree, size)
}

func BenchmarkBTreePut100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	tree := New[int, struct{}](128)

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkPut(b, tree, size)
}

func BenchmarkBTreeRemove100(b *testing.B) {
	b.StopTimer()

	size := 100
	tree := New[int, struct{}](128)

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkRemove(b, tree, size)
}

func BenchmarkBTreeRemove1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	tree := New[int, struct{}](128)

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkRemove(b, tree, size)
}

func BenchmarkBTreeRemove10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	tree := New[int, struct{}](128)

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkRemove(b, tree, size)
}

func BenchmarkBTreeRemove100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	tree := New[int, struct{}](128)

	for n := range size {
		tree.Put(n, struct{}{})
	}

	b.StartTimer()
	benchmarkRemove(b, tree, size)
}
