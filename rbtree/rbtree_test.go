package rbtree_test

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/qntx/gods/rbtree"
)

func TestRedBlackTreeGet(t *testing.T) {
	t.Parallel()

	tree := rbtree.New[int, string]()

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

	fmt.Println(tree)
	//
	//  RedBlackTree
	//  │           ┌── 6
	//  │       ┌── 5
	//  │   ┌── 4
	//  │   │   └── 3
	//  └── 2
	//       └── 1

	if actualValue := tree.Len(); actualValue != 6 {
		t.Errorf("Got %v expected %v", actualValue, 6)
	}

	if actualValue := tree.GetNode(4).Size(); actualValue != 4 {
		t.Errorf("Got %v expected %v", actualValue, 4)
	}

	if actualValue := tree.GetNode(2).Size(); actualValue != 6 {
		t.Errorf("Got %v expected %v", actualValue, 6)
	}

	if actualValue := tree.GetNode(8).Size(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}
}

func TestRedBlackTreePut(t *testing.T) {
	t.Parallel()

	// Initialize and populate the tree
	tree := rbtree.New[int, string]()
	tree.Put(5, "e")
	tree.Put(6, "f")
	tree.Put(7, "g")
	tree.Put(3, "c")
	tree.Put(4, "d")
	tree.Put(1, "x")
	tree.Put(2, "b")
	tree.Put(1, "a") // overwrite

	// Test length
	if got := tree.Len(); got != 7 {
		t.Errorf("Len() = %d, want 7", got)
	}

	// Test keys
	wantKeys := []int{1, 2, 3, 4, 5, 6, 7}
	if got := tree.Keys(); !slices.Equal(got, wantKeys) {
		t.Errorf("Keys() = %v, want %v", got, wantKeys)
	}

	// Test values
	wantValues := []string{"a", "b", "c", "d", "e", "f", "g"}
	if got := tree.Values(); !slices.Equal(got, wantValues) {
		t.Errorf("Values() = %v, want %v", got, wantValues)
	}

	// Test individual retrievals with structured data
	tests := []struct {
		key       int
		wantVal   string
		wantFound bool
	}{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "e", true},
		{6, "f", true},
		{7, "g", true},
		{8, "", false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Get(%d)", tt.key), func(t *testing.T) {
			t.Parallel()

			gotVal, gotFound := tree.Get(tt.key)
			if gotVal != tt.wantVal || gotFound != tt.wantFound {
				t.Errorf("Get(%d) = (%q, %v), want (%q, %v)", tt.key, gotVal, gotFound, tt.wantVal, tt.wantFound)
			}
		})
	}
}

func TestRedBlackTreeDelete(t *testing.T) {
	t.Parallel()

	// Initialize tree with data
	tree := rbtree.New[int, string]()
	tree.Put(5, "e")
	tree.Put(6, "f")
	tree.Put(7, "g")
	tree.Put(3, "c")
	tree.Put(4, "d")
	tree.Put(1, "x")
	tree.Put(2, "b")
	tree.Put(1, "a") // overwrite

	// Test partial removal
	t.Run("PartialRemoval", func(t *testing.T) {
		tree.Delete(5)
		tree.Delete(6)
		tree.Delete(7)
		tree.Delete(8) // Non-existent key
		tree.Delete(5) // Already removed

		wantKeys := []int{1, 2, 3, 4}
		if got := tree.Keys(); !slices.Equal(got, wantKeys) {
			t.Errorf("Keys() = %v, want %v", got, wantKeys)
		}

		wantValues := []string{"a", "b", "c", "d"}
		if got := tree.Values(); !slices.Equal(got, wantValues) {
			t.Errorf("Values() = %v, want %v", got, wantValues)
		}

		if got := tree.Len(); got != 4 {
			t.Errorf("Len() = %d, want 4", got)
		}

		// Structured test cases for Get
		tests := []struct {
			key       int
			wantVal   string
			wantFound bool
		}{
			{1, "a", true},
			{2, "b", true},
			{3, "c", true},
			{4, "d", true},
			{5, "", false},
			{6, "", false},
			{7, "", false},
			{8, "", false},
		}

		for _, tt := range tests {
			t.Run(fmt.Sprintf("Get(%d)", tt.key), func(t *testing.T) {
				gotVal, gotFound := tree.Get(tt.key)
				if gotVal != tt.wantVal || gotFound != tt.wantFound {
					t.Errorf("Get(%d) = (%q, %v), want (%q, %v)", tt.key, gotVal, gotFound, tt.wantVal, tt.wantFound)
				}
			})
		}
	})

	// Test full removal
	t.Run("FullRemoval", func(t *testing.T) {
		t.Parallel()

		tree.Delete(1)
		tree.Delete(4)
		tree.Delete(2)
		tree.Delete(3)
		tree.Delete(2) // Already removed
		tree.Delete(2) // Already removed

		wantKeys := []int{}
		if got := tree.Keys(); !slices.Equal(got, wantKeys) {
			t.Errorf("Keys() = %v, want %v", got, wantKeys)
		}

		wantValues := []string{}
		if got := tree.Values(); !slices.Equal(got, wantValues) {
			t.Errorf("Values() = %v, want %v", got, wantValues)
		}

		if gotEmpty, gotSize := tree.IsEmpty(), tree.Len(); !gotEmpty || gotSize != 0 {
			t.Errorf("Empty() = %v, Len() = %d, want true, 0", gotEmpty, gotSize)
		}
	})
}

func TestRedBlackTreeLeftAndRight(t *testing.T) {
	t.Parallel()

	tree := rbtree.New[int, string]()

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

	if actualValue, expectedValue := tree.GetBeginNode().Key(), 1; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, expectedValue := tree.GetBeginNode().Value(), "x"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, expectedValue := tree.GetEndNode().Key(), 7; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, expectedValue := tree.GetEndNode().Value(), "g"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestRedBlackTreeCeilingAndFloor(t *testing.T) {
	t.Parallel()

	tree := rbtree.New[int, string]()

	if node, found := tree.Floor(0); node != nil || found {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}

	if node, found := tree.Ceiling(0); node != nil || found {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}

	tree.Put(5, "e")
	tree.Put(6, "f")
	tree.Put(7, "g")
	tree.Put(3, "c")
	tree.Put(4, "d")
	tree.Put(1, "x")
	tree.Put(2, "b")

	if node, found := tree.Floor(4); node.Key() != 4 || !found {
		t.Errorf("Got %v expected %v", node.Key(), 4)
	}

	if node, found := tree.Floor(0); node != nil || found {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}

	if node, found := tree.Ceiling(4); node.Key() != 4 || !found {
		t.Errorf("Got %v expected %v", node.Key(), 4)
	}

	if node, found := tree.Ceiling(8); node != nil || found {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}
}

func TestRedBlackTreeSerialization(t *testing.T) {
	t.Parallel()

	tree := rbtree.New[string, string]()
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

	_, err = json.Marshal([]any{"a", "b", "c", tree})
	if err != nil {
		t.Errorf("Got error %v", err)
	}

	intTree := rbtree.New[string, int]()

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

func TestRedBlackTreeString(t *testing.T) {
	t.Parallel()

	c := rbtree.New[string, int]()
	c.Put("a", 1)

	if !strings.HasPrefix(c.String(), "RedBlackTree") {
		t.Errorf("String should start with container name")
	}
}
