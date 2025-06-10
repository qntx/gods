package avltree

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"
)

func TestAVLTreeGet(t *testing.T) {
	tree := New[int, string]()

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
	//
	//  AVLTree
	//  │       ┌── 6
	//  │   ┌── 5
	//  └── 4
	//      │   ┌── 3
	//      └── 2
	//          └── 1

	if actualValue := tree.Len(); actualValue != 6 {
		t.Errorf("Got %v expected %v", actualValue, 6)
	}

	if actualValue := tree.GetNode(2).Size(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}

	if actualValue := tree.GetNode(4).Size(); actualValue != 6 {
		t.Errorf("Got %v expected %v", actualValue, 6)
	}

	if actualValue := tree.GetNode(7).Size(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}
}

func TestAVLTreePut(t *testing.T) {
	tree := New[int, string]()
	tree.Put(5, "e")
	tree.Put(6, "f")
	tree.Put(7, "g")
	tree.Put(3, "c")
	tree.Put(4, "d")
	tree.Put(1, "x")
	tree.Put(2, "b")
	tree.Put(1, "a") //overwrite

	if actualValue := tree.Len(); actualValue != 7 {
		t.Errorf("Got %v expected %v", actualValue, 7)
	}

	if actualValue, expectedValue := tree.Keys(), []int{1, 2, 3, 4, 5, 6, 7}; !slices.Equal(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, expectedValue := tree.Values(), []string{"a", "b", "c", "d", "e", "f", "g"}; !slices.Equal(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	tests1 := [][]interface{}{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "e", true},
		{6, "f", true},
		{7, "g", true},
		{8, "", false},
	}

	for _, test := range tests1 {
		// retrievals
		actualValue, actualFound := tree.Get(test[0].(int))
		if actualValue != test[1] || actualFound != test[2] {
			t.Errorf("Got %v expected %v", actualValue, test[1])
		}
	}
}

func TestAVLTreeRemove(t *testing.T) {
	tree := New[int, string]()
	tree.Put(5, "e")
	tree.Put(6, "f")
	tree.Put(7, "g")
	tree.Put(3, "c")
	tree.Put(4, "d")
	tree.Put(1, "x")
	tree.Put(2, "b")
	tree.Put(1, "a") //overwrite

	tree.Delete(5)
	tree.Delete(6)
	tree.Delete(7)
	tree.Delete(8)
	tree.Delete(5)

	if actualValue, expectedValue := tree.Keys(), []int{1, 2, 3, 4}; !slices.Equal(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, expectedValue := tree.Values(), []string{"a", "b", "c", "d"}; !slices.Equal(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue := tree.Len(); actualValue != 4 {
		t.Errorf("Got %v expected %v", actualValue, 7)
	}

	tests2 := [][]interface{}{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "", false},
		{6, "", false},
		{7, "", false},
		{8, "", false},
	}

	for _, test := range tests2 {
		actualValue, actualFound := tree.Get(test[0].(int))
		if actualValue != test[1] || actualFound != test[2] {
			t.Errorf("Got %v expected %v", actualValue, test[1])
		}
	}

	tree.Delete(1)
	tree.Delete(4)
	tree.Delete(2)
	tree.Delete(3)
	tree.Delete(2)
	tree.Delete(2)

	if actualValue, expectedValue := tree.Keys(), []int{}; !slices.Equal(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, expectedValue := tree.Values(), []string{}; !slices.Equal(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if empty, size := tree.IsEmpty(), tree.Len(); empty != true || size != -0 {
		t.Errorf("Got %v expected %v", empty, true)
	}
}

func TestAVLTreeLeftAndRight(t *testing.T) {
	tree := New[int, string]()

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

func TestAVLTreeCeilingAndFloor(t *testing.T) {
	tree := New[int, string]()

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

func TestAVLTreeIterEmpty(t *testing.T) {
	tree := New[int, string]()
	count := 0

	for range tree.Iter() {
		count++
	}

	if count != 0 {
		t.Errorf("should not iterate on an empty tree, but counted %d elements", count)
	}
}

func TestTreeIterForward(t *testing.T) {
	tree := New[int, string]()
	tree.Put(5, "e")
	tree.Put(6, "f")
	tree.Put(7, "g")
	tree.Put(3, "c")
	tree.Put(4, "d")
	tree.Put(1, "a")
	tree.Put(2, "b")

	expectedKeys := []int{1, 2, 3, 4, 5, 6, 7}
	expectedValues := []string{"a", "b", "c", "d", "e", "f", "g"}

	actualKeys := make([]int, 0, tree.Len())
	actualValues := make([]string, 0, tree.Len())

	for key, value := range tree.Iter() {
		actualKeys = append(actualKeys, key)
		actualValues = append(actualValues, value)
	}

	if !slices.Equal(actualKeys, expectedKeys) {
		t.Errorf("forward iteration keys mismatch:\ngot:  %v\nwant: %v", actualKeys, expectedKeys)
	}

	if !slices.Equal(actualValues, expectedValues) {
		t.Errorf("forward iteration values mismatch:\ngot:  %v\nwant: %v", actualValues, expectedValues)
	}
}

func TestTreeReverseIterOnEmpty(t *testing.T) {
	tree := New[int, string]()

	var count int
	for range tree.ReverseIter() {
		count++
	}

	if count != 0 {
		t.Errorf("ReverseIter() on an empty tree should produce 0 elements, got %d", count)
	}
}

func TestTreeReverseIter(t *testing.T) {
	tree := New[int, string]()
	tree.Put(3, "c")
	tree.Put(1, "a")
	tree.Put(2, "b")

	expectedKeys := []int{3, 2, 1}
	expectedValues := []string{"c", "b", "a"}

	actualKeys := make([]int, 0, tree.Len())
	actualValues := make([]string, 0, tree.Len())

	for key, value := range tree.ReverseIter() {
		actualKeys = append(actualKeys, key)
		actualValues = append(actualValues, value)
	}

	if !slices.Equal(actualKeys, expectedKeys) {
		t.Errorf("reverse iteration keys mismatch:\ngot:  %v\nwant: %v", actualKeys, expectedKeys)
	}

	if !slices.Equal(actualValues, expectedValues) {
		t.Errorf("reverse iteration values mismatch:\ngot:  %v\nwant: %v", actualValues, expectedValues)
	}
}

func TestTreeIterSkipsDeleted(t *testing.T) {
	tree := New[int, string]()
	tree.Put(10, "a")
	tree.Put(20, "b")
	tree.Put(30, "c")
	tree.Put(5, "d")

	tree.Delete(20) // Delete a middle element
	tree.Delete(5)  // Delete a leaf element

	expectedKeys := []int{10, 30}

	var actualKeys []int
	for k := range tree.Iter() {
		actualKeys = append(actualKeys, k)
	}

	if !slices.Equal(actualKeys, expectedKeys) {
		t.Errorf("iterator did not skip deleted elements correctly:\ngot:  %v\nwant: %v", actualKeys, expectedKeys)
	}
}

func TestAVLTreeSerialization(t *testing.T) {
	tree := New[string, string]()
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

	intTree := New[string, int]()

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

func TestAVLTreeString(t *testing.T) {
	c := New[int, int]()
	c.Put(1, 1)
	c.Put(2, 1)
	c.Put(3, 1)
	c.Put(4, 1)
	c.Put(5, 1)
	c.Put(6, 1)
	c.Put(7, 1)
	c.Put(8, 1)

	if !strings.HasPrefix(c.String(), "AVLTree") {
		t.Errorf("String should start with container name")
	}
}
