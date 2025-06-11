// Package btree implements a B-tree for ordered key-value storage.
//
// A B-tree is a self-balancing tree data structure that maintains sorted data
// and allows searches, sequential access, insertions, and deletions in
// logarithmic time.
//
// Properties of a B-Tree of order m:
// - Every node has at most m children.
// - Every non-leaf node (except root) has at least ⌈m/2⌉ children.
// - The root has at least two children if it is not a leaf node.
// - A non-leaf node with k children contains k−1 keys.
// - All leaves appear at the same level.
//
// This implementation is not thread-safe.
//
// References: https://en.wikipedia.org/wiki/B-tree
package btree

import (
	"encoding/json"
	"fmt"
	"iter"
	"maps"
	"slices"
	"strings"

	"github.com/qntx/gods/cmp"
	"github.com/qntx/gods/container"
)

// notFound is a sentinel value indicating a key was not found.
const notFound = -1

// entry represents an internal key-value pair within a B-tree node.
type entry[K comparable, V any] struct {
	key   K
	value V
}

// Key returns the key of the entry.
func (e *entry[K, V]) Key() K {
	return e.key
}

// Value returns the value of the entry.
func (e *entry[K, V]) Value() V {
	return e.value
}

// String returns a string representation of the entry.
func (e *entry[K, V]) String() string {
	return fmt.Sprintf("%v", e.key)
}

// Node is a single element within the tree, containing entries and children.
type Node[K comparable, V any] struct {
	parent   *Node[K, V]
	children []*Node[K, V]  // Children nodes
	entries  []*entry[K, V] // Contained entries in node
}

// Parent returns the parent of the node.
func (n *Node[K, V]) Parent() *Node[K, V] {
	return n.parent
}

// Entries returns the key-value pairs of the node.
func (n *Node[K, V]) Entries() []*entry[K, V] {
	return n.entries
}

// Children returns the children nodes of the node.
func (n *Node[K, V]) Children() []*Node[K, V] {
	return n.children
}

// Size returns the number of nodes in the subtree rooted at this node.
// Time complexity: O(n).
func (n *Node[K, V]) Size() int {
	if n == nil {
		return 0
	}

	size := 1
	for _, children := range n.children {
		size += children.Size()
	}

	return size
}

// height is an internal helper to calculate the height from a given node downwards.
func (n *Node[K, V]) height() int {
	h := 1
	// A leaf has height 1. For non-leaf nodes, traverse down the leftmost path.
	for !n.isLeaf() {
		h++
		n = n.children[0]
	}

	return h
}

// String returns a string representation of the node.
func (n *Node[K, V]) String() string {
	var sb strings.Builder

	sb.WriteString("[")

	for i, e := range n.entries {
		if i > 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(e.String())
	}

	sb.WriteString("]")

	return sb.String()
}

// Tree holds the elements and configuration of the B-tree.
type Tree[K comparable, V any] struct {
	root *Node[K, V]       // Root node of the tree.
	cmp  cmp.Comparator[K] // Key comparator.
	len  int               // Total number of key-value pairs in the tree.
	m    int               // Order (maximum number of children).
}

// Root returns the root node of the tree.
func (t *Tree[K, V]) Root() *Node[K, V] {
	return t.root
}

// Comparator returns the comparator used by the tree.
func (t *Tree[K, V]) Comparator() cmp.Comparator[K] {
	return t.cmp
}

// MaxChildren returns the maximum number of children allowed in a node.
func (t *Tree[K, V]) MaxChildren() int {
	return t.m
}

// Ensure Tree implements the required interfaces.
var _ container.OrderedMap[int, int] = (*Tree[int, int])(nil)
var _ json.Marshaler = (*Tree[string, int])(nil)
var _ json.Unmarshaler = (*Tree[string, int])(nil)

// New creates a new B-tree with the specified order and a built-in comparator.
// The order `m` must be 3 or greater. Panics if order is invalid.
// K must be an ordered type (e.g., int, string). Time complexity: O(1).
func New[K cmp.Ordered, V any](order int) *Tree[K, V] {
	return NewWith[K, V](order, cmp.Compare[K])
}

// NewWith creates a new B-tree with a custom comparator.
// The order `m` must be 3 or greater. Panics if order is invalid.
// Time complexity: O(1).
func NewWith[K comparable, V any](order int, cmp cmp.Comparator[K]) *Tree[K, V] {
	if order < 3 {
		panic("Invalid B-tree order: must be 3 or greater")
	}

	return &Tree[K, V]{m: order, cmp: cmp}
}

// Put inserts a key-value pair into the tree, updating the value if the key already exists.
// Time complexity: O(log n).
func (t *Tree[K, V]) Put(key K, value V) {
	e := &entry[K, V]{key: key, value: value}

	if t.root == nil {
		t.root = &Node[K, V]{entries: []*entry[K, V]{e}}
		t.len++

		return
	}

	if t.insert(t.root, e) {
		t.len++
	}
}

// Get retrieves the value for a given key.
// Returns the value and true if found, or the zero value and false otherwise.
// Time complexity: O(log n).
func (t *Tree[K, V]) Get(key K) (V, bool) {
	node, index := t.lookup(key)
	if index != notFound {
		return node.entries[index].value, true
	}

	var zeroV V

	return zeroV, false
}

// GetNode retrieves the node containing the given key.
// Returns the node if found, nil otherwise.
// Time complexity: O(log n).
func (t *Tree[K, V]) GetNode(key K) *Node[K, V] {
	node, index := t.lookup(key)
	if index != notFound {
		return node
	}

	return nil
}

// Has checks if a key exists in the tree.
// Time complexity: O(log n).
func (t *Tree[K, V]) Has(key K) bool {
	_, index := t.lookup(key)

	return index != notFound
}

// Delete removes a key-value pair from the tree.
// Returns true if the key was found and removed, false otherwise.
// Time complexity: O(log n).
func (t *Tree[K, V]) Delete(key K) bool {
	node, index := t.lookup(key)
	if index == notFound {
		return false
	}

	t.delete(node, index)

	t.len--
	if t.len == 0 {
		t.root = nil
	}

	return true
}

// Begin returns the minimum key-value pair.
// Time complexity: O(log n).
func (t *Tree[K, V]) Begin() (k K, v V, ok bool) {
	if t.IsEmpty() {
		return k, v, false
	}

	node := getMinNode(t.root)
	e := node.entries[0]

	return e.key, e.value, true
}

// End returns the maximum key-value pair.
// Time complexity: O(log n).
func (t *Tree[K, V]) End() (k K, v V, ok bool) {
	if t.IsEmpty() {
		return k, v, false
	}

	node := getMaxNode(t.root)
	e := node.entries[len(node.entries)-1]

	return e.key, e.value, true
}

// DeleteBegin removes the minimum key-value pair.
// Returns the removed pair and true, or zero values and false if the tree is empty.
// Time complexity: O(log n).
func (t *Tree[K, V]) DeleteBegin() (k K, v V, ok bool) {
	key, value, found := t.Begin()
	if found {
		t.Delete(key)

		return key, value, true
	}

	return k, v, false
}

// DeleteEnd removes the maximum key-value pair.
// Returns the removed pair and true, or zero values and false if the tree is empty.
// Time complexity: O(log n).
func (t *Tree[K, V]) DeleteEnd() (k K, v V, ok bool) {
	key, value, found := t.End()
	if found {
		t.Delete(key)

		return key, value, true
	}

	return k, v, false
}

// GetBeginNode returns the node with the minimum key.
// Returns nil if the tree is empty.
// Time complexity: O(log n).
func (tree *Tree[K, V]) GetBeginNode() *Node[K, V] {
	return getMinNode(tree.root)
}

// GetEndNode returns the node with the maximum key.
// Returns nil if the tree is empty.
// Time complexity: O(log n).
func (tree *Tree[K, V]) GetEndNode() *Node[K, V] {
	return getMaxNode(tree.root)
}

// Height returns the height of the tree. A tree with a single node has a height of 1.
// Returns 0 if the tree is empty.
// Time complexity: O(log n).
func (t *Tree[K, V]) Height() int {
	if t.root == nil {
		return 0
	}

	return t.root.height()
}

// Len returns the number of items in the tree. Time complexity: O(1).
func (t *Tree[K, V]) Len() int { return t.len }

// IsEmpty returns true if the tree has no items. Time complexity: O(1).
func (t *Tree[K, V]) IsEmpty() bool { return t.len == 0 }

// Clear removes all items from the tree. Time complexity: O(1).
func (t *Tree[K, V]) Clear() {
	t.root = nil
	t.len = 0
}

// Keys returns a slice of all keys in sorted order. Time complexity: O(n).
func (t *Tree[K, V]) Keys() []K {
	keys := make([]K, 0, t.len)
	for k := range t.Iter() {
		keys = append(keys, k)
	}

	return keys
}

// Values returns a slice of all values in sorted key order. Time complexity: O(n).
func (t *Tree[K, V]) Values() []V {
	values := make([]V, 0, t.len)
	for _, v := range t.Iter() {
		values = append(values, v)
	}

	return values
}

// ToSlice returns a slice of all key-value pairs in sorted order. Time complexity: O(n).
func (t *Tree[K, V]) ToSlice() []V {
	return t.Values()
}

// Entries returns slices of all keys and values in sorted order. Time complexity: O(n).
func (t *Tree[K, V]) Entries() ([]K, []V) {
	keys := make([]K, 0, t.len)
	values := make([]V, 0, t.len)

	for k, v := range t.Iter() {
		keys = append(keys, k)
		values = append(values, v)
	}

	return keys, values
}

// Clone creates a deep copy of the tree. Time complexity: O(n).
func (t *Tree[K, V]) Clone() container.Map[K, V] {
	newTree := &Tree[K, V]{m: t.m, cmp: t.cmp, len: t.len}
	if t.root != nil {
		newTree.root = cloneNode(t.root, nil)
	}

	return newTree
}

// Iter returns an iterator for in-order traversal.
func (t *Tree[K, V]) Iter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		inorder(t.root, yield)
	}
}

// RIter returns an iterator for reverse-order traversal.
func (t *Tree[K, V]) RIter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		inorderReverse(t.root, yield)
	}
}

// String returns a string representation of the tree for debugging.
func (t *Tree[K, V]) String() string {
	if t.IsEmpty() {
		return "BTree[]"
	}

	var sb strings.Builder

	sb.WriteString("BTree\n")
	t.output(&sb, t.root, "", true)

	return sb.String()
}

// MarshalJSON implements the json.Marshaler interface. Time complexity: O(n).
func (t *Tree[K, V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(maps.Collect(t.Iter()))
}

// UnmarshalJSON implements the json.Unmarshaler interface. Time complexity: O(m log m).
func (t *Tree[K, V]) UnmarshalJSON(data []byte) error {
	var elems map[K]V
	if err := json.Unmarshal(data, &elems); err != nil {
		return err
	}

	t.Clear()

	for k, v := range elems {
		t.Put(k, v)
	}

	return nil
}

// isLeaf checks if a node is a leaf (has no children).
func (n *Node[K, V]) isLeaf() bool {
	return len(n.children) == 0
}

// lookup finds the node and entry index for a given key.
// Returns the node and index, or (nil, notFound) if the key doesn't exist.
func (t *Tree[K, V]) lookup(key K) (*Node[K, V], int) {
	if t.root == nil {
		return nil, notFound
	}

	node := t.root

	for {
		index, found := t.search(node, key)
		if found {
			return node, index
		}

		if node.isLeaf() {
			return nil, notFound
		}

		node = node.children[index]
	}
}

// search performs a binary search for a key within a single node's entries.
func (t *Tree[K, V]) search(node *Node[K, V], key K) (index int, found bool) {
	return slices.BinarySearchFunc(node.entries, key, func(e *entry[K, V], k K) int {
		return t.cmp(e.key, k)
	})
}

// insert handles the insertion of an entry, returning true if the tree size increased.
func (t *Tree[K, V]) insert(node *Node[K, V], e *entry[K, V]) bool {
	if node.isLeaf() {
		return t.insertIntoLeaf(node, e)
	}

	return t.insertIntoInternal(node, e)
}

func (t *Tree[K, V]) insertIntoLeaf(node *Node[K, V], e *entry[K, V]) bool {
	index, found := t.search(node, e.key)
	if found {
		node.entries[index] = e

		return false
	}

	node.entries = slices.Insert(node.entries, index, e)
	t.split(node)

	return true
}

func (t *Tree[K, V]) insertIntoInternal(node *Node[K, V], e *entry[K, V]) bool {
	index, found := t.search(node, e.key)
	if found {
		node.entries[index] = e

		return false
	}

	return t.insert(node.children[index], e)
}

// split divides a node if it has too many entries.
func (t *Tree[K, V]) split(node *Node[K, V]) {
	if len(node.entries) <= t.maxEntries() {
		return
	}

	if node == t.root {
		t.splitRoot()

		return
	}

	t.splitNonRoot(node)
}

func (t *Tree[K, V]) splitNonRoot(node *Node[K, V]) {
	middle := t.middle()
	parent := node.parent

	// Promote middle entry to parent
	medianEntry := node.entries[middle]
	parentIndex, _ := t.search(parent, medianEntry.key)
	parent.entries = slices.Insert(parent.entries, parentIndex, medianEntry)

	// Create new right sibling
	rightSibling := &Node[K, V]{
		parent:  parent,
		entries: slices.Clone(node.entries[middle+1:]),
	}
	if !node.isLeaf() {
		rightSibling.children = slices.Clone(node.children[middle+1:])
		setParent(rightSibling.children, rightSibling)
	}

	// Update original node to be the left sibling
	node.entries = node.entries[:middle]
	if !node.isLeaf() {
		node.children = node.children[:middle+1]
	}

	// Insert right sibling into parent's children
	parent.children = slices.Insert(parent.children, parentIndex+1, rightSibling)

	t.split(parent)
}

func (t *Tree[K, V]) splitRoot() {
	middle := t.middle()
	medianEntry := t.root.entries[middle]

	left := &Node[K, V]{entries: slices.Clone(t.root.entries[:middle])}
	right := &Node[K, V]{entries: slices.Clone(t.root.entries[middle+1:])}

	if !t.root.isLeaf() {
		left.children = slices.Clone(t.root.children[:middle+1])
		right.children = slices.Clone(t.root.children[middle+1:])
		setParent(left.children, left)
		setParent(right.children, right)
	}

	newRoot := &Node[K, V]{
		entries:  []*entry[K, V]{medianEntry},
		children: []*Node[K, V]{left, right},
	}
	left.parent = newRoot
	right.parent = newRoot
	t.root = newRoot
}

// delete handles the core deletion logic.
func (t *Tree[K, V]) delete(node *Node[K, V], index int) {
	// If node is an internal node, swap with successor to move deletion to a leaf.
	if !node.isLeaf() {
		successorNode := getMinNode(node.children[index+1])
		node.entries[index] = successorNode.entries[0]
		node, index = successorNode, 0 // Target the successor for actual deletion
	}

	// Delete entry from the (now guaranteed to be leaf) node.
	node.entries = slices.Delete(node.entries, index, index+1)
	t.rebalance(node)
}

// rebalance ensures the B-tree properties are maintained after a deletion.
func (t *Tree[K, V]) rebalance(node *Node[K, V]) {
	// If node has enough entries, it doesn't need rebalancing.
	if len(node.entries) >= t.minEntries() {
		return
	}

	// If node is the current root of the tree.
	// The root can have fewer than minEntries (e.g., if it's the only node or tree height is 1).
	// If it becomes empty, mergeWithSibling handles root replacement.
	if node == t.root {
		return
	}

	// At this point: node is NOT t.root AND len(node.entries) < t.minEntries().
	// Such a node must have a parent. If not, it's an orphaned former root.
	parent := node.parent
	if parent == nil {
		// This state (node != t.root && node.parent == nil) occurs if 'node'
		// was the root, got emptied during a merge, t.root was reassigned,
		// and rebalance is now called on this 'node' (the old, detached root).
		// No further parent-based rebalancing is applicable.
		return
	}

	nodeIndex := findChildIndex(parent, node)

	// Try to borrow from left sibling
	if nodeIndex > 0 {
		leftSibling := parent.children[nodeIndex-1]
		if len(leftSibling.entries) > t.minEntries() {
			t.borrowFromSibling(node, leftSibling, nodeIndex)

			return
		}
	}

	// Try to borrow from right sibling
	if nodeIndex < len(parent.children)-1 {
		rightSibling := parent.children[nodeIndex+1]
		if len(rightSibling.entries) > t.minEntries() {
			t.borrowFromSibling(node, rightSibling, nodeIndex)

			return
		}
	}

	// Merge with a sibling
	if nodeIndex > 0 {
		t.mergeWithSibling(parent.children[nodeIndex-1], node, nodeIndex-1) // Merge with left
	} else {
		t.mergeWithSibling(node, parent.children[nodeIndex+1], nodeIndex) // Merge with right
	}

	t.rebalance(parent)
}

// borrowFromSibling performs a rotation to rebalance the tree.
func (t *Tree[K, V]) borrowFromSibling(node, sibling *Node[K, V], nodeIndexInParent int) {
	parent := node.parent
	// Determine if the sibling is to the left or right of the node
	// by finding the sibling's index in the parent's children list.
	siblingIndexInParent := findChildIndex(parent, sibling)

	if siblingIndexInParent < nodeIndexInParent { // Sibling is to the left
		parentIndex := nodeIndexInParent - 1 // This is also siblingIndexInParent
		// Rotate right
		node.entries = slices.Insert(node.entries, 0, parent.entries[parentIndex])
		parent.entries[parentIndex] = sibling.entries[len(sibling.entries)-1]

		sibling.entries = slices.Delete(sibling.entries, len(sibling.entries)-1, len(sibling.entries))
		if !sibling.isLeaf() {
			childToMove := sibling.children[len(sibling.children)-1]
			childToMove.parent = node
			node.children = slices.Insert(node.children, 0, childToMove)
			sibling.children = slices.Delete(sibling.children, len(sibling.children)-1, len(sibling.children))
		}
	} else { // Sibling is to the right (siblingIndexInParent > nodeIndexInParent)
		parentIndex := nodeIndexInParent // Separator entry in parent is at node's original index
		// Rotate left
		node.entries = append(node.entries, parent.entries[parentIndex])
		parent.entries[parentIndex] = sibling.entries[0]

		sibling.entries = slices.Delete(sibling.entries, 0, 1)
		if !sibling.isLeaf() {
			childToMove := sibling.children[0]
			childToMove.parent = node
			node.children = append(node.children, childToMove)
			sibling.children = slices.Delete(sibling.children, 0, 1)
		}
	}
}

// mergeWithSibling merges two adjacent nodes.
func (t *Tree[K, V]) mergeWithSibling(left, right *Node[K, V], leftIndexInParent int) {
	parent := left.parent
	// Pull down separator from parent
	separator := parent.entries[leftIndexInParent]
	left.entries = append(left.entries, separator)
	left.entries = append(left.entries, right.entries...)

	if !left.isLeaf() {
		left.children = append(left.children, right.children...)
		setParent(left.children, left)
	}

	// Remove separator from parent and right children from parent's children
	parent.entries = slices.Delete(parent.entries, leftIndexInParent, leftIndexInParent+1)
	parent.children = slices.Delete(parent.children, leftIndexInParent+1, leftIndexInParent+2)

	if len(parent.entries) == 0 && parent == t.root {
		t.root = left
		left.parent = nil
	}
}

func (t *Tree[K, V]) maxEntries() int { return t.m - 1 }
func (t *Tree[K, V]) minEntries() int { return (t.m+1)/2 - 1 }
func (t *Tree[K, V]) middle() int     { return (t.m - 1) / 2 }

func setParent[K comparable, V any](nodes []*Node[K, V], parent *Node[K, V]) {
	for _, n := range nodes {
		n.parent = parent
	}
}

// findChildIndex finds the index of a children in its parent's children slice.
func findChildIndex[K comparable, V any](parent, children *Node[K, V]) int {
	for i, c := range parent.children {
		if c == children {
			return i
		}
	}

	return notFound
}

func getMinNode[K comparable, V any](node *Node[K, V]) *Node[K, V] {
	if node == nil {
		return nil
	}

	for !node.isLeaf() {
		node = node.children[0]
	}

	return node
}

func getMaxNode[K comparable, V any](node *Node[K, V]) *Node[K, V] {
	if node == nil {
		return nil
	}

	for !node.isLeaf() {
		node = node.children[len(node.children)-1]
	}

	return node
}

func cloneNode[K comparable, V any](node *Node[K, V], parent *Node[K, V]) *Node[K, V] {
	if node == nil {
		return nil
	}

	newNode := &Node[K, V]{parent: parent}
	newNode.entries = make([]*entry[K, V], len(node.entries))

	for i, e := range node.entries {
		newNode.entries[i] = &entry[K, V]{key: e.key, value: e.value}
	}

	if !node.isLeaf() {
		newNode.children = make([]*Node[K, V], len(node.children))
		for i, c := range node.children {
			newNode.children[i] = cloneNode(c, newNode)
		}
	}

	return newNode
}

// inorder traversal for the iterator.
func inorder[K comparable, V any](n *Node[K, V], yield func(K, V) bool) bool {
	if n == nil {
		return true
	}

	for i := range len(n.entries) {
		if !n.isLeaf() {
			if !inorder(n.children[i], yield) {
				return false
			}
		}

		if !yield(n.entries[i].key, n.entries[i].value) {
			return false
		}
	}

	if !n.isLeaf() {
		if !inorder(n.children[len(n.children)-1], yield) {
			return false
		}
	}

	return true
}

// inorderReverse traversal for the reverse iterator.
func inorderReverse[K comparable, V any](n *Node[K, V], yield func(K, V) bool) bool {
	if n == nil {
		return true
	}

	if !n.isLeaf() {
		if !inorderReverse(n.children[len(n.children)-1], yield) {
			return false
		}
	}

	for i := len(n.entries) - 1; i >= 0; i-- {
		if !yield(n.entries[i].key, n.entries[i].value) {
			return false
		}

		if !n.isLeaf() {
			if !inorderReverse(n.children[i], yield) {
				return false
			}
		}
	}

	return true
}

// output generates the string representation of the tree.
func (t *Tree[K, V]) output(sb *strings.Builder, n *Node[K, V], prefix string, isTail bool) {
	if n == nil {
		return
	}

	sb.WriteString(prefix)

	if isTail {
		sb.WriteString("└── ")

		prefix += "    "
	} else {
		sb.WriteString("├── ")

		prefix += "│   "
	}

	entryKeys := make([]string, len(n.entries))
	for i, e := range n.entries {
		entryKeys[i] = fmt.Sprintf("%v", e.key)
	}

	sb.WriteString(strings.Join(entryKeys, ", ") + "\n")

	for i, children := range n.children {
		t.output(sb, children, prefix, i == len(n.children)-1)
	}
}
