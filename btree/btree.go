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

var _ container.OrderedMap[int, int] = (*Tree[int, int])(nil)
var _ json.Marshaler = (*Tree[string, int])(nil)
var _ json.Unmarshaler = (*Tree[string, int])(nil)

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
// Returns the value and true if the key was found and removed, false otherwise.
// Time complexity: O(log n).
func (t *Tree[K, V]) Delete(key K) (v V, found bool) {
	node, index := t.lookup(key)
	if index == notFound {
		return v, false
	}

	v = node.entries[index].value
	t.delete(node, index)

	t.len--
	if t.len == 0 {
		t.root = nil
	}

	return v, true
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

func (t *Tree[K, V]) splitNonRoot(n *Node[K, V]) {
	mid := t.middle()
	p := n.parent

	// Promote middle entry to parent
	med := n.entries[mid]
	pi, _ := t.search(p, med.key)
	p.entries = slices.Insert(p.entries, pi, med)

	// Create new right sibling
	r := &Node[K, V]{
		parent:  p,
		entries: slices.Clone(n.entries[mid+1:]),
	}
	if !n.isLeaf() {
		r.children = slices.Clone(n.children[mid+1:])
		setParent(r.children, r)
	}

	// Update original node to be the left sibling
	n.entries = n.entries[:mid]
	if !n.isLeaf() {
		n.children = n.children[:mid+1]
	}

	// Insert right sibling into parent's children
	p.children = slices.Insert(p.children, pi+1, r)

	t.split(p)
}

func (t *Tree[K, V]) splitRoot() {
	mid := t.middle()
	med := t.root.entries[mid]

	l := &Node[K, V]{entries: slices.Clone(t.root.entries[:mid])}
	r := &Node[K, V]{entries: slices.Clone(t.root.entries[mid+1:])}

	if !t.root.isLeaf() {
		l.children = slices.Clone(t.root.children[:mid+1])
		r.children = slices.Clone(t.root.children[mid+1:])
		setParent(l.children, l)
		setParent(r.children, r)
	}

	nr := &Node[K, V]{
		entries:  []*entry[K, V]{med},
		children: []*Node[K, V]{l, r},
	}

	l.parent = nr
	r.parent = nr

	t.root = nr
}

// delete handles the core deletion logic.
func (t *Tree[K, V]) delete(n *Node[K, V], i int) {
	// If node is internal, swap with successor to move deletion to a leaf.
	if !n.isLeaf() {
		s := getMinNode(n.children[i+1])
		n.entries[i] = s.entries[0]
		n, i = s, 0 // Target the successor for actual deletion
	}

	// Delete entry from the (now guaranteed to be leaf) node.
	n.entries = slices.Delete(n.entries, i, i+1)
	t.rebalance(n)
}

// rebalance ensures B-tree properties post-deletion.
func (t *Tree[K, V]) rebalance(n *Node[K, V]) {
	// If node has enough entries, no rebalancing needed.
	if len(n.entries) >= t.minEntries() {
		return
	}

	// Root can have fewer entries or be empty (handled by merge).
	if n == t.root {
		return
	}

	// Non-root node with too few entries must have a parent.
	p := n.parent
	if p == nil {
		// Orphaned former root; no further rebalancing.
		return
	}

	i := findChildIndex(p, n)

	// Try borrowing from left sibling.
	if i > 0 {
		l := p.children[i-1]
		if len(l.entries) > t.minEntries() {
			t.borrowFromSibling(n, l, i)

			return
		}
	}

	// Try borrowing from right sibling.
	if i < len(p.children)-1 {
		r := p.children[i+1]
		if len(r.entries) > t.minEntries() {
			t.borrowFromSibling(n, r, i)

			return
		}
	}

	// Merge with a sibling.
	if i > 0 {
		t.mergeWithSibling(p.children[i-1], n, i-1) // Merge with left.
	} else {
		t.mergeWithSibling(n, p.children[i+1], i) // Merge with right.
	}

	t.rebalance(p)
}

// borrowFromSibling rotates to rebalance the tree.
func (t *Tree[K, V]) borrowFromSibling(n, s *Node[K, V], ni int) {
	p := n.parent
	// Find sibling's index in parent's children.
	si := findChildIndex(p, s)

	if si < ni { // Left sibling
		pi := ni - 1
		// Rotate right
		n.entries = slices.Insert(n.entries, 0, p.entries[pi])
		p.entries[pi] = s.entries[len(s.entries)-1]

		s.entries = slices.Delete(s.entries, len(s.entries)-1, len(s.entries))
		if !s.isLeaf() {
			c := s.children[len(s.children)-1]
			c.parent = n
			n.children = slices.Insert(n.children, 0, c)
			s.children = slices.Delete(s.children, len(s.children)-1, len(s.children))
		}
	} else { // Right sibling
		pi := ni
		// Rotate left
		n.entries = append(n.entries, p.entries[pi])
		p.entries[pi] = s.entries[0]

		s.entries = slices.Delete(s.entries, 0, 1)
		if !s.isLeaf() {
			c := s.children[0]
			c.parent = n
			n.children = append(n.children, c)
			s.children = slices.Delete(s.children, 0, 1)
		}
	}
}

// mergeWithSibling merges two adjacent nodes.
func (t *Tree[K, V]) mergeWithSibling(l, r *Node[K, V], li int) {
	p := l.parent
	// Pull down separator from parent
	sep := p.entries[li]
	l.entries = append(l.entries, sep)
	l.entries = append(l.entries, r.entries...)

	if !l.isLeaf() {
		l.children = append(l.children, r.children...)
		setParent(l.children, l)
	}

	// Remove separator and right node from parent
	p.entries = slices.Delete(p.entries, li, li+1)
	p.children = slices.Delete(p.children, li+1, li+2)

	if len(p.entries) == 0 && p == t.root {
		t.root = l
		l.parent = nil
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
