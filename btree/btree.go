// Package btree implements a B tree.
//
// According to Knuth's definition, a B-tree of order m is a tree which satisfies the following properties:
// - Every node has at most m children.
// - Every non-leaf node (except root) has at least ⌈m/2⌉ children.
// - The root has at least two children if it is not a leaf node.
// - A non-leaf node with k children contains k−1 keys.
// - All leaves appear in the same level
//
// Structure is not thread safe.
//
// References: https://en.wikipedia.org/wiki/B-tree
package btree

import (
	"bytes"
	"encoding/json"
	"fmt"
	"iter"
	"strings"

	"slices"

	"github.com/qntx/gods/cmp"
	"github.com/qntx/gods/container"
)

// Entry represents the key-value pair contained within nodes.
type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

func (entry *Entry[K, V]) String() string {
	return fmt.Sprintf("%v", entry.Key)
}

// Node is a single element within the tree.
type Node[K comparable, V any] struct {
	Parent   *Node[K, V]
	Entries  []*Entry[K, V] // Contained keys in node
	Children []*Node[K, V]  // Children nodes
}

// Height returns the height of the node.
func (node *Node[K, V]) height() int {
	height := 0
	for ; node != nil; node = node.Children[0] {
		height++

		if len(node.Children) == 0 {
			break
		}
	}

	return height
}

// Size returns the number of elements stored in the subtree.
// Computed dynamically on each call, i.e. the subtree is traversed to count the number of the nodes.
func (node *Node[K, V]) Size() int {
	if node == nil {
		return 0
	}

	size := 1
	for _, child := range node.Children {
		size += child.Size()
	}

	return size
}

var _ container.OrderedMap[int, int] = (*Tree[int, int])(nil)
var _ json.Marshaler = (*Tree[string, int])(nil)
var _ json.Unmarshaler = (*Tree[string, int])(nil)

// Tree holds elements of the B-tree.
type Tree[K comparable, V any] struct {
	Root       *Node[K, V]       // Root node
	Comparator cmp.Comparator[K] // Key comparator
	size       int               // Total number of keys in the tree
	m          int               // order (maximum number of children)
}

// New instantiates a B-tree with the order (maximum number of children) and the built-in comparator for K.
func New[K cmp.Ordered, V any](order int) *Tree[K, V] {
	return NewWith[K, V](order, cmp.Compare[K])
}

// NewWith instantiates a B-tree with the order (maximum number of children) and a custom key comparator.
func NewWith[K comparable, V any](order int, comparator cmp.Comparator[K]) *Tree[K, V] {
	if order < 3 {
		panic("Invalid order, should be at least 3")
	}

	return &Tree[K, V]{m: order, Comparator: comparator}
}

// Put inserts key-value pair node into the tree.
// If key already exists, then its value is updated with the new value.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *Tree[K, V]) Put(key K, value V) {
	entry := &Entry[K, V]{Key: key, Value: value}

	if tree.Root == nil {
		tree.Root = &Node[K, V]{Entries: []*Entry[K, V]{entry}, Children: []*Node[K, V]{}}
		tree.size++

		return
	}

	if tree.insert(tree.Root, entry) {
		tree.size++
	}
}

// Get searches the node in the tree by key and returns its value or nil if key is not found in tree.
// Second return parameter is true if key was found, otherwise false.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *Tree[K, V]) Get(key K) (value V, found bool) {
	node, index, found := tree.searchRecursively(tree.Root, key)
	if found {
		return node.Entries[index].Value, true
	}

	return value, false
}

// GetNode searches the node in the tree by key and returns its node or nil if key is not found in tree.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *Tree[K, V]) GetNode(key K) *Node[K, V] {
	node, _, _ := tree.searchRecursively(tree.Root, key)

	return node
}

// Has checks if the given key exists in the tree.
// Returns true if the key is found, false otherwise.
// Time complexity: O(log n).
func (tree *Tree[K, V]) Has(key K) bool {
	_, _, found := tree.searchRecursively(tree.Root, key)

	return found
}

// Delete remove the node from the tree by key.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *Tree[K, V]) Delete(key K) bool {
	node, index, found := tree.searchRecursively(tree.Root, key)
	if found {
		tree.delete(node, index)

		tree.size--

		return true
	}

	return false
}

// GetBeginNode returns the node with the minimum key.
// Returns nil if the tree is empty.
// Time complexity: O(log n).
func (tree *Tree[K, V]) GetBeginNode() *Node[K, V] {
	return tree.left(tree.Root)
}

// GetEndNode returns the node with the maximum key.
// Returns nil if the tree is empty.
// Time complexity: O(log n).
func (tree *Tree[K, V]) GetEndNode() *Node[K, V] {
	return tree.right(tree.Root)
}

// Begin returns the minimum key-value pair in the tree.
// Returns the key, value, and true if found, zero values and false otherwise.
// Time complexity: O(log n).
func (tree *Tree[K, V]) Begin() (key K, value V, found bool) {
	node := tree.GetBeginNode()
	if node != nil && len(node.Entries) > 0 {
		entry := node.Entries[0]

		return entry.Key, entry.Value, true
	}

	var zeroK K

	var zeroV V

	return zeroK, zeroV, false
}

// End returns the maximum key-value pair in the tree.
// Returns the key, value, and true if found, zero values and false otherwise.
// Time complexity: O(log n).
func (tree *Tree[K, V]) End() (key K, value V, found bool) {
	node := tree.GetEndNode()
	if node != nil && len(node.Entries) > 0 {
		entry := node.Entries[len(node.Entries)-1]

		return entry.Key, entry.Value, true
	}

	var zeroK K

	var zeroV V

	return zeroK, zeroV, false
}

// DeleteBegin removes the minimum key-value pair from the tree.
// Returns the removed key, value, and true if successful, zero values and false otherwise.
// Time complexity: O(log n).
func (tree *Tree[K, V]) DeleteBegin() (key K, value V, removed bool) {
	node := tree.GetBeginNode()
	if node != nil && len(node.Entries) > 0 {
		entry := node.Entries[0]
		key, value = entry.Key, entry.Value
		tree.Delete(key)

		return key, value, true
	}

	var zeroK K

	var zeroV V

	return zeroK, zeroV, false
}

// DeleteEnd removes the maximum key-value pair from the tree.
// Returns the removed key, value, and true if successful, zero values and false otherwise.
// Time complexity: O(log n).
func (tree *Tree[K, V]) DeleteEnd() (key K, value V, removed bool) {
	node := tree.GetEndNode()
	if node != nil && len(node.Entries) > 0 {
		entry := node.Entries[len(node.Entries)-1]
		key, value = entry.Key, entry.Value
		tree.Delete(key)

		return key, value, true
	}

	var zeroK K

	var zeroV V

	return zeroK, zeroV, false
}

// Len returns the number of key-value pairs in the tree.
// Time complexity: O(1).
func (tree *Tree[K, V]) Len() int {
	return tree.size
}

// IsEmpty returns true if the tree contains no key-value pairs.
// Time complexity: O(1).
func (tree *Tree[K, V]) IsEmpty() bool {
	return tree.size == 0
}

// Clear removes all key-value pairs from the tree.
// Time complexity: O(1).
func (tree *Tree[K, V]) Clear() {
	tree.Root = nil
	tree.size = 0
}

// Clone creates a deep copy of the tree.
// Time complexity: O(n).
func (tree *Tree[K, V]) Clone() container.Map[K, V] {
	newTree := &Tree[K, V]{
		Comparator: tree.Comparator,
		m:          tree.m,
		size:       tree.size,
	}
	if tree.Root != nil {
		newTree.Root = cloneNode(tree.Root, nil)
	}

	return newTree
}

func (tree *Tree[K, V]) MaxChildren() int {
	return tree.m
}

// Keys returns all keys in-order.
func (tree *Tree[K, V]) Keys() []K {
	keys := make([]K, 0, tree.size)

	for k := range tree.Iter() {
		keys = append(keys, k)
	}

	return keys
}

// Values returns all values in-order based on the key.
func (tree *Tree[K, V]) Values() []V {
	values := make([]V, 0, tree.size)

	for _, v := range tree.Iter() {
		values = append(values, v)
	}

	return values
}

// ToSlice returns all key-value pairs in-order based on the key.
func (tree *Tree[K, V]) ToSlice() []V {
	return tree.Values()
}

// Entries returns all key-value pairs in-order based on the key.
func (tree *Tree[K, V]) Entries() ([]K, []V) {
	keys := make([]K, 0, tree.size)
	values := make([]V, 0, tree.size)

	for k, v := range tree.Iter() {
		keys = append(keys, k)
		values = append(values, v)
	}

	return keys, values
}

// Iter returns an iterator over key-value pairs from the tree.
// The iteration order is in-order based on the key.
func (tree *Tree[K, V]) Iter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		it := tree.Iterator()
		for it.Next() {
			if !yield(it.Key(), it.Value()) {
				return
			}
		}
	}
}

// RIter returns an iterator over key-value pairs in descending order.
// Time complexity: O(n) for full iteration.
func (tree *Tree[K, V]) RIter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		it := tree.Iterator()
		it.End()

		for it.Prev() {
			if !yield(it.Key(), it.Value()) {
				return
			}
		}
	}
}

// Height returns the height of the tree.
func (tree *Tree[K, V]) Height() int {
	return tree.Root.height()
}

// String returns a string representation of container (for debugging purposes).
func (tree *Tree[K, V]) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("BTree\n")

	if !tree.IsEmpty() {
		tree.output(&buffer, tree.Root, 0)
	}

	return buffer.String()
}

// ToJSON outputs the JSON representation of the tree.
func (tree *Tree[K, V]) MarshalJSON() ([]byte, error) {
	elements := make(map[K]V)
	it := tree.Iterator()

	for it.Next() {
		elements[it.Key()] = it.Value()
	}

	return json.Marshal(&elements)
}

// FromJSON populates the tree from the input JSON representation.
func (tree *Tree[K, V]) UnmarshalJSON(data []byte) error {
	elements := make(map[K]V)

	err := json.Unmarshal(data, &elements)
	if err != nil {
		return err
	}

	tree.Clear()

	for key, value := range elements {
		tree.Put(key, value)
	}

	return err
}

func (tree *Tree[K, V]) output(buffer *bytes.Buffer, node *Node[K, V], level int) {
	for e := range len(node.Entries) + 1 {
		if e < len(node.Children) {
			tree.output(buffer, node.Children[e], level+1)
		}

		if e < len(node.Entries) {
			buffer.WriteString(strings.Repeat("    ", level))
			buffer.WriteString(fmt.Sprintf("%v", node.Entries[e].Key) + "\n")
		}
	}
}

func (tree *Tree[K, V]) isLeaf(node *Node[K, V]) bool {
	return len(node.Children) == 0
}

func (tree *Tree[K, V]) shouldSplit(node *Node[K, V]) bool {
	return len(node.Entries) > tree.maxEntries()
}

func (tree *Tree[K, V]) maxChildren() int {
	return tree.m
}

func (tree *Tree[K, V]) minChildren() int {
	return (tree.m + 1) / 2 // ceil(m/2)
}

func (tree *Tree[K, V]) maxEntries() int {
	return tree.maxChildren() - 1
}

func (tree *Tree[K, V]) minEntries() int {
	return tree.minChildren() - 1
}

func (tree *Tree[K, V]) middle() int {
	return (tree.m - 1) / 2 // "-1" to favor right nodes to have more keys when splitting
}

// search searches only within the single node among its entries.
func (tree *Tree[K, V]) search(node *Node[K, V], key K) (index int, found bool) {
	low, high := 0, len(node.Entries)-1

	var mid int

	for low <= high {
		mid = (high + low) / 2
		compare := tree.Comparator(key, node.Entries[mid].Key)

		switch {
		case compare > 0:
			low = mid + 1
		case compare < 0:
			high = mid - 1
		case compare == 0:
			return mid, true
		}
	}

	return low, false
}

// searchRecursively searches recursively down the tree starting at the startNode.
func (tree *Tree[K, V]) searchRecursively(startNode *Node[K, V], key K) (node *Node[K, V], index int, found bool) {
	if tree.IsEmpty() {
		return nil, -1, false
	}

	node = startNode

	for {
		index, found = tree.search(node, key)
		if found {
			return node, index, true
		}

		if tree.isLeaf(node) {
			return nil, -1, false
		}

		node = node.Children[index]
	}
}

func (tree *Tree[K, V]) insert(node *Node[K, V], entry *Entry[K, V]) (inserted bool) {
	if tree.isLeaf(node) {
		return tree.insertIntoLeaf(node, entry)
	}

	return tree.insertIntoInternal(node, entry)
}

func (tree *Tree[K, V]) insertIntoLeaf(node *Node[K, V], entry *Entry[K, V]) (inserted bool) {
	insertPosition, found := tree.search(node, entry.Key)
	if found {
		node.Entries[insertPosition] = entry

		return false
	}
	// Insert entry's key in the middle of the node
	node.Entries = append(node.Entries, nil)
	copy(node.Entries[insertPosition+1:], node.Entries[insertPosition:])
	node.Entries[insertPosition] = entry
	tree.split(node)

	return true
}

func (tree *Tree[K, V]) insertIntoInternal(node *Node[K, V], entry *Entry[K, V]) (inserted bool) {
	insertPosition, found := tree.search(node, entry.Key)
	if found {
		node.Entries[insertPosition] = entry

		return false
	}

	return tree.insert(node.Children[insertPosition], entry)
}

func (tree *Tree[K, V]) split(node *Node[K, V]) {
	if !tree.shouldSplit(node) {
		return
	}

	if node == tree.Root {
		tree.splitRoot()

		return
	}

	tree.splitNonRoot(node)
}

func (tree *Tree[K, V]) splitNonRoot(node *Node[K, V]) {
	middle := tree.middle()
	parent := node.Parent

	left := &Node[K, V]{Entries: slices.Clone(node.Entries[:middle]), Parent: parent}
	right := &Node[K, V]{Entries: slices.Clone(node.Entries[middle+1:]), Parent: parent}

	// Move children from the node to be split into left and right nodes
	if !tree.isLeaf(node) {
		left.Children = slices.Clone(node.Children[:middle+1])
		right.Children = slices.Clone(node.Children[middle+1:])
		setParent(left.Children, left)
		setParent(right.Children, right)
	}

	insertPosition, _ := tree.search(parent, node.Entries[middle].Key)

	// Insert middle key into parent
	parent.Entries = append(parent.Entries, nil)
	copy(parent.Entries[insertPosition+1:], parent.Entries[insertPosition:])
	parent.Entries[insertPosition] = node.Entries[middle]

	// Set child left of inserted key in parent to the created left node
	parent.Children[insertPosition] = left

	// Set child right of inserted key in parent to the created right node
	parent.Children = append(parent.Children, nil)
	copy(parent.Children[insertPosition+2:], parent.Children[insertPosition+1:])
	parent.Children[insertPosition+1] = right

	tree.split(parent)
}

func (tree *Tree[K, V]) splitRoot() {
	middle := tree.middle()

	left := &Node[K, V]{Entries: slices.Clone(tree.Root.Entries[:middle])}
	right := &Node[K, V]{Entries: slices.Clone(tree.Root.Entries[middle+1:])}

	// Move children from the node to be split into left and right nodes
	if !tree.isLeaf(tree.Root) {
		left.Children = slices.Clone(tree.Root.Children[:middle+1])
		right.Children = slices.Clone(tree.Root.Children[middle+1:])
		setParent(left.Children, left)
		setParent(right.Children, right)
	}

	// Root is a node with one entry and two children (left and right)
	newRoot := &Node[K, V]{
		Entries:  []*Entry[K, V]{tree.Root.Entries[middle]},
		Children: []*Node[K, V]{left, right},
	}

	left.Parent = newRoot
	right.Parent = newRoot
	tree.Root = newRoot
}

func setParent[K comparable, V any](nodes []*Node[K, V], parent *Node[K, V]) {
	for _, node := range nodes {
		node.Parent = parent
	}
}

func (tree *Tree[K, V]) left(node *Node[K, V]) *Node[K, V] {
	if tree.IsEmpty() {
		return nil
	}

	current := node

	for {
		if tree.isLeaf(current) {
			return current
		}

		current = current.Children[0]
	}
}

func (tree *Tree[K, V]) right(node *Node[K, V]) *Node[K, V] {
	if tree.IsEmpty() {
		return nil
	}

	current := node

	for {
		if tree.isLeaf(current) {
			return current
		}

		current = current.Children[len(current.Children)-1]
	}
}

// leftSibling returns the node's left sibling and child index (in parent) if it exists, otherwise (nil,-1)
// key is any of keys in node (could even be deleted).
func (tree *Tree[K, V]) leftSibling(node *Node[K, V], key K) (*Node[K, V], int) {
	if node.Parent != nil {
		index, _ := tree.search(node.Parent, key)
		index--

		if index >= 0 && index < len(node.Parent.Children) {
			return node.Parent.Children[index], index
		}
	}

	return nil, -1
}

// rightSibling returns the node's right sibling and child index (in parent) if it exists, otherwise (nil,-1)
// key is any of keys in node (could even be deleted).
func (tree *Tree[K, V]) rightSibling(node *Node[K, V], key K) (*Node[K, V], int) {
	if node.Parent != nil {
		index, _ := tree.search(node.Parent, key)
		index++

		if index < len(node.Parent.Children) {
			return node.Parent.Children[index], index
		}
	}

	return nil, -1
}

// delete deletes an entry in node at entries' index
// ref.: https://en.wikipedia.org/wiki/B-tree#Deletion
func (tree *Tree[K, V]) delete(node *Node[K, V], index int) {
	// deleting from a leaf node
	if tree.isLeaf(node) {
		deletedKey := node.Entries[index].Key
		tree.deleteEntry(node, index)
		tree.rebalance(node, deletedKey)

		if len(tree.Root.Entries) == 0 {
			tree.Root = nil
		}

		return
	}

	// deleting from an internal node
	leftLargestNode := tree.right(node.Children[index]) // largest node in the left sub-tree (assumed to exist)
	leftLargestEntryIndex := len(leftLargestNode.Entries) - 1
	node.Entries[index] = leftLargestNode.Entries[leftLargestEntryIndex]
	deletedKey := leftLargestNode.Entries[leftLargestEntryIndex].Key
	tree.deleteEntry(leftLargestNode, leftLargestEntryIndex)
	tree.rebalance(leftLargestNode, deletedKey)
}

// rebalance rebalances the tree after deletion if necessary and returns true, otherwise false.
// Note that we first delete the entry and then call rebalance, thus the passed deleted key as reference.
func (tree *Tree[K, V]) rebalance(node *Node[K, V], deletedKey K) {
	// check if rebalancing is needed
	if node == nil || len(node.Entries) >= tree.minEntries() {
		return
	}

	// try to borrow from left sibling
	leftSibling, leftSiblingIndex := tree.leftSibling(node, deletedKey)
	if leftSibling != nil && len(leftSibling.Entries) > tree.minEntries() {
		// rotate right
		node.Entries = append([]*Entry[K, V]{node.Parent.Entries[leftSiblingIndex]}, node.Entries...) // prepend parent's separator entry to node's entries
		node.Parent.Entries[leftSiblingIndex] = leftSibling.Entries[len(leftSibling.Entries)-1]
		tree.deleteEntry(leftSibling, len(leftSibling.Entries)-1)

		if !tree.isLeaf(leftSibling) {
			leftSiblingRightMostChild := leftSibling.Children[len(leftSibling.Children)-1]
			leftSiblingRightMostChild.Parent = node
			node.Children = append([]*Node[K, V]{leftSiblingRightMostChild}, node.Children...)
			tree.deleteChild(leftSibling, len(leftSibling.Children)-1)
		}

		return
	}

	// try to borrow from right sibling
	rightSibling, rightSiblingIndex := tree.rightSibling(node, deletedKey)
	if rightSibling != nil && len(rightSibling.Entries) > tree.minEntries() {
		// rotate left
		node.Entries = append(node.Entries, node.Parent.Entries[rightSiblingIndex-1]) // append parent's separator entry to node's entries
		node.Parent.Entries[rightSiblingIndex-1] = rightSibling.Entries[0]
		tree.deleteEntry(rightSibling, 0)

		if !tree.isLeaf(rightSibling) {
			rightSiblingLeftMostChild := rightSibling.Children[0]
			rightSiblingLeftMostChild.Parent = node
			node.Children = append(node.Children, rightSiblingLeftMostChild)

			tree.deleteChild(rightSibling, 0)
		}

		return
	}

	// merge with siblings
	if rightSibling != nil {
		// merge with right sibling
		node.Entries = append(node.Entries, node.Parent.Entries[rightSiblingIndex-1])
		node.Entries = append(node.Entries, rightSibling.Entries...)
		deletedKey = node.Parent.Entries[rightSiblingIndex-1].Key
		tree.deleteEntry(node.Parent, rightSiblingIndex-1)
		tree.appendChildren(node.Parent.Children[rightSiblingIndex], node)
		tree.deleteChild(node.Parent, rightSiblingIndex)
	} else if leftSibling != nil {
		// merge with left sibling
		entries := slices.Clone(leftSibling.Entries)
		entries = append(entries, node.Parent.Entries[leftSiblingIndex])
		node.Entries = append(entries, node.Entries...)
		deletedKey = node.Parent.Entries[leftSiblingIndex].Key
		tree.deleteEntry(node.Parent, leftSiblingIndex)
		tree.prependChildren(node.Parent.Children[leftSiblingIndex], node)
		tree.deleteChild(node.Parent, leftSiblingIndex)
	}

	// make the merged node the root if its parent was the root and the root is empty
	if node.Parent == tree.Root && len(tree.Root.Entries) == 0 {
		tree.Root = node
		node.Parent = nil

		return
	}

	// parent might underflow, so try to rebalance if necessary
	tree.rebalance(node.Parent, deletedKey)
}

func (tree *Tree[K, V]) prependChildren(fromNode *Node[K, V], toNode *Node[K, V]) {
	children := append([]*Node[K, V](nil), fromNode.Children...)
	toNode.Children = append(children, toNode.Children...)
	setParent(fromNode.Children, toNode)
}

func (tree *Tree[K, V]) appendChildren(fromNode *Node[K, V], toNode *Node[K, V]) {
	toNode.Children = append(toNode.Children, fromNode.Children...)
	setParent(fromNode.Children, toNode)
}

func (tree *Tree[K, V]) deleteEntry(node *Node[K, V], index int) {
	copy(node.Entries[index:], node.Entries[index+1:])
	node.Entries[len(node.Entries)-1] = nil
	node.Entries = node.Entries[:len(node.Entries)-1]
}

func (tree *Tree[K, V]) deleteChild(node *Node[K, V], index int) {
	if index >= len(node.Children) {
		return
	}

	copy(node.Children[index:], node.Children[index+1:])
	node.Children[len(node.Children)-1] = nil
	node.Children = node.Children[:len(node.Children)-1]
}

func cloneNode[K comparable, V any](node *Node[K, V], parent *Node[K, V]) *Node[K, V] {
	if node == nil {
		return nil
	}

	newNode := &Node[K, V]{
		Parent:   parent,
		Entries:  make([]*Entry[K, V], len(node.Entries)),
		Children: make([]*Node[K, V], len(node.Children)),
	}
	for i, entry := range node.Entries {
		newNode.Entries[i] = &Entry[K, V]{Key: entry.Key, Value: entry.Value}
	}

	for i, child := range node.Children {
		newNode.Children[i] = cloneNode(child, newNode)
	}

	return newNode
}
