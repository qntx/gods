// Package avltree implements a self-balancing AVL binary search tree for ordered key-value storage.
//
// The AVL tree maintains O(log n) time complexity for insertion, deletion, and search operations
// by ensuring the tree remains height-balanced. This implementation is not thread-safe.
//
// Reference: https://en.wikipedia.org/wiki/AVL_tree
package avltree

import (
	"encoding/json"
	"fmt"
	"iter"
	"maps"
	"strings"

	"github.com/qntx/gods/cmp"
	"github.com/qntx/gods/container"
)

// Node represents a single element in the AVL tree.
type Node[K comparable, V any] struct {
	key           K           // The key used for ordering
	value         V           // The value associated with the key
	parent        *Node[K, V] // Parent node
	left          *Node[K, V] // Left child node
	right         *Node[K, V] // Right child node
	balanceFactor int8        // Balance factor: height(right) - height(left)
}

// Key returns the key stored in the node.
// Time complexity: O(1).
func (n *Node[K, V]) Key() K {
	return n.key
}

// Value returns the value associated with the node's key.
// Time complexity: O(1).
func (n *Node[K, V]) Value() V {
	return n.value
}

// Left returns the left child of the node, or nil if none exists.
// Time complexity: O(1).
func (n *Node[K, V]) Left() *Node[K, V] {
	return n.left
}

// Right returns the right child of the node, or nil if none exists.
// Time complexity: O(1).
func (n *Node[K, V]) Right() *Node[K, V] {
	return n.right
}

// Parent returns the parent of the node, or nil if the node is the root.
// Time complexity: O(1).
func (n *Node[K, V]) Parent() *Node[K, V] {
	return n.parent
}

// Size returns the number of nodes in the subtree rooted at this node.
// Computed dynamically by traversing the subtree. Time complexity: O(n).
func (n *Node[K, V]) Size() int {
	if n == nil {
		return 0
	}

	return 1 + n.left.Size() + n.right.Size()
}

// String returns a string representation of the node.
// Time complexity: O(1).
func (n *Node[K, V]) String() string {
	return fmt.Sprintf("%v", n.key)
}

var _ container.Map[int, int] = (*Tree[int, int])(nil)

// Tree manages an AVL tree storing key-value pairs.
//
// K must be comparable and compatible with the provided comparator.
// V can be any type.
type Tree[K comparable, V any] struct {
	root       *Node[K, V]       // Root node of the tree
	len        int               // Number of nodes in the tree
	comparator cmp.Comparator[K] // Comparator for ordering keys
}

// New creates a new AVL tree with a default comparator for ordered types.
//
// K must implement cmp.Ordered (e.g., int, string). Time complexity: O(1).
func New[K cmp.Ordered, V any]() *Tree[K, V] {
	return &Tree[K, V]{comparator: cmp.GenericComparator[K]}
}

// NewWith creates a new AVL tree with a custom comparator.
//
// The comparator defines the key ordering. Time complexity: O(1).
func NewWith[K comparable, V any](comparator cmp.Comparator[K]) *Tree[K, V] {
	return &Tree[K, V]{comparator: comparator}
}

// Put inserts or updates a key-value pair in the tree.
//
// If the key exists, its value is updated; otherwise, a new node is inserted.
// Panics if the key type is incompatible with the comparator.
// Time complexity: O(log n).
func (t *Tree[K, V]) Put(key K, val V) {
	if t.root == nil {
		t.root = &Node[K, V]{key: key, value: val}
		t.len++

		return
	}

	node, parent := t.root, (*Node[K, V])(nil)

	var cmpResult int

	for node != nil {
		parent = node
		cmpResult = t.comparator(key, node.key)

		switch {
		case cmpResult < 0:
			node = node.left
		case cmpResult > 0:
			node = node.right
		default: // cmpResult == 0
			node.value = val

			return
		}
	}

	newNode := &Node[K, V]{key: key, value: val, parent: parent}
	if cmpResult < 0 {
		parent.left = newNode
	} else {
		parent.right = newNode
	}

	t.len++

	t.insertFixup(parent)
}

// Delete removes the node with the specified key from the tree.
//
// Returns true if a node was removed, false if the key was not found.
// Panics if the key type is incompatible with the comparator.
// Time complexity: O(log n).
func (t *Tree[K, V]) Delete(key K) bool {
	node := t.lookup(key)
	if node == nil {
		return false
	}

	var fixupStartNode *Node[K, V]

	if node.left != nil && node.right != nil {
		// Node has two children: find the in-order successor (smallest node in right subtree)
		successor := t.getLeftNode(node.right)
		// Copy successor's data to the current node
		node.key = successor.key
		node.value = successor.value
		// Mark successor for deletion (simplifies the problem)
		node = successor
	}

	// At this point, 'node' has at most one child
	fixupStartNode = node.parent

	child := node.left
	if child == nil {
		child = node.right
	}

	t.replaceNode(node, child)

	t.len--

	if fixupStartNode != nil {
		t.deleteFixup(fixupStartNode)
	}

	return true
}

// Get retrieves the value associated with the specified key.
//
// Returns the value and true if found, or a zero value and false if not.
// Panics if the key type is incompatible with the comparator.
// Time complexity: O(log n).
func (t *Tree[K, V]) Get(key K) (val V, ok bool) {
	if node := t.lookup(key); node != nil {
		return node.value, true
	}

	var zeroVal V

	return zeroVal, false
}

// GetNode retrieves the node associated with the specified key.
//
// Returns the node if found, or nil if not.
// Panics if the key type is incompatible with the comparator.
// Time complexity: O(log n).
func (t *Tree[K, V]) GetNode(key K) *Node[K, V] {
	return t.lookup(key)
}

// Has checks if the specified key exists in the tree.
//
// Returns true if the key is found, false otherwise.
// Panics if the key type is incompatible with the comparator.
// Time complexity: O(log n).
func (t *Tree[K, V]) Has(key K) bool {
	return t.lookup(key) != nil
}

// GetBeginNode returns the leftmost node (minimum key), or nil if the tree is empty.
// Time complexity: O(log n).
func (t *Tree[K, V]) GetBeginNode() *Node[K, V] {
	return t.getLeftNode(t.root)
}

// GetEndNode returns the rightmost node (maximum key), or nil if the tree is empty.
// Time complexity: O(log n).
func (t *Tree[K, V]) GetEndNode() *Node[K, V] {
	return t.getRightNode(t.root)
}

// Begin returns the minimum key and value in the tree.
//
// Returns found as true if an element is found, false otherwise.
// Time complexity: O(log n).
func (t *Tree[K, V]) Begin() (key K, value V, found bool) {
	node := t.GetBeginNode()
	if node != nil {
		return node.key, node.value, true
	}

	var zeroKey K

	var zeroValue V

	return zeroKey, zeroValue, false
}

// End returns the maximum key and value in the tree.
//
// Returns found as true if an element is found, false otherwise.
// Time complexity: O(log n).
func (t *Tree[K, V]) End() (key K, value V, found bool) {
	node := t.GetEndNode()
	if node != nil {
		return node.key, node.value, true
	}

	var zeroKey K

	var zeroValue V

	return zeroKey, zeroValue, false
}

// DeleteBegin removes the minimum key-value pair from the tree.
//
// Returns the removed key, value, and true if an element was removed, false otherwise.
// Time complexity: O(log n).
func (t *Tree[K, V]) DeleteBegin() (key K, value V, removed bool) {
	node := t.GetBeginNode()
	if node != nil {
		k, v := node.key, node.value
		t.Delete(k)

		return k, v, true
	}

	var zeroKey K

	var zeroValue V

	return zeroKey, zeroValue, false
}

// DeleteEnd removes the maximum key-value pair from the tree.
//
// Returns the removed key, value, and true if an element was removed, false otherwise.
// Time complexity: O(log n).
func (t *Tree[K, V]) DeleteEnd() (key K, value V, removed bool) {
	node := t.GetEndNode()
	if node != nil {
		k, v := node.key, node.value
		t.Delete(k)

		return k, v, true
	}

	var zeroKey K

	var zeroValue V

	return zeroKey, zeroValue, false
}

// Floor finds the largest node with a key less than or equal to the given key.
//
// Returns the node and true if found, or nil and false if not.
// Panics if the key type is incompatible with the comparator.
// Time complexity: O(log n).
func (t *Tree[K, V]) Floor(key K) (*Node[K, V], bool) {
	var floor *Node[K, V]

	node := t.root
	for node != nil {
		switch cmp := t.comparator(key, node.key); {
		case cmp == 0:
			return node, true
		case cmp > 0:
			floor = node
			node = node.right
		default:
			node = node.left
		}
	}

	return floor, floor != nil
}

// Ceiling finds the smallest node with a key greater than or equal to the given key.
//
// Returns the node and true if found, or nil and false if not.
// Panics if the key type is incompatible with the comparator.
// Time complexity: O(log n).
func (t *Tree[K, V]) Ceiling(key K) (*Node[K, V], bool) {
	var ceil *Node[K, V]

	node := t.root
	for node != nil {
		switch cmp := t.comparator(key, node.key); {
		case cmp == 0:
			return node, true
		case cmp < 0:
			ceil = node
			node = node.left
		default:
			node = node.right
		}
	}

	return ceil, ceil != nil
}

// Keys returns all keys in in-order sequence.
// Time complexity: O(n).
func (t *Tree[K, V]) Keys() []K {
	keys := make([]K, 0, t.len)
	t.inOrderTraversal(func(n *Node[K, V]) {
		keys = append(keys, n.key)
	})

	return keys
}

// Values returns all values in in-order sequence based on their keys.
// Time complexity: O(n).
func (t *Tree[K, V]) Values() []V {
	values := make([]V, 0, t.len)
	t.inOrderTraversal(func(n *Node[K, V]) {
		values = append(values, n.value)
	})

	return values
}

// ToSlice returns all values in in-order sequence.
// Time complexity: O(n).
func (t *Tree[K, V]) ToSlice() []V {
	return t.Values()
}

// Entries returns all keys and values in in-order sequence.
//
// More efficient than calling Keys() and Values() separately as it traverses
// the tree only once. Time complexity: O(n).
func (t *Tree[K, V]) Entries() ([]K, []V) {
	keys := make([]K, 0, t.len)
	vals := make([]V, 0, t.len)
	t.inOrderTraversal(func(n *Node[K, V]) {
		keys = append(keys, n.key)
		vals = append(vals, n.value)
	})

	return keys, vals
}

// Len returns the number of nodes in the tree.
// Time complexity: O(1).
func (t *Tree[K, V]) Len() int {
	return t.len
}

// IsEmpty checks if the tree contains no nodes.
// Time complexity: O(1).
func (t *Tree[K, V]) IsEmpty() bool {
	return t.len == 0
}

// Clear removes all nodes from the tree.
// Time complexity: O(1).
func (t *Tree[K, V]) Clear() {
	t.root = nil
	t.len = 0
}

// Clone creates a deep copy of the tree.
//
// The new tree has independent nodes from the original.
// Time complexity: O(n).
func (t *Tree[K, V]) Clone() container.Map[K, V] {
	newTree := &Tree[K, V]{
		comparator: t.comparator,
		len:        t.len,
	}

	if t.root == nil {
		return newTree
	}

	newTree.root = cloneNode(t.root, nil)

	return newTree
}

// Iter returns an iterator over all key-value pairs in sorted order.
//
// Conforms to Go 1.22+ iterator design (iter.Seq2). Yields pairs via an efficient,
// non-recursive in-order traversal. First element retrieval is O(log n), subsequent
// steps are amortized O(1), with overall iteration complexity of O(n).
func (t *Tree[K, V]) Iter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		node := t.GetBeginNode()
		for node != nil {
			if !yield(node.key, node.value) {
				return
			}

			if node.right != nil {
				node = t.getLeftNode(node.right)
			} else {
				for node.parent != nil && node == node.parent.right {
					node = node.parent
				}

				node = node.parent
			}
		}
	}
}

// ReverseIter returns a reverse iterator over all key-value pairs (from largest to smallest).
//
// Conforms to Go 1.22+ iterator design (iter.Seq2). Yields pairs via an efficient,
// non-recursive reverse in-order traversal. First element retrieval is O(log n), subsequent
// steps are amortized O(1), with overall iteration complexity of O(n).
func (t *Tree[K, V]) ReverseIter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		node := t.GetEndNode()
		for node != nil {
			if !yield(node.key, node.value) {
				return
			}

			if node.left != nil {
				node = t.getRightNode(node.left)
			} else {
				for node.parent != nil && node == node.parent.left {
					node = node.parent
				}

				node = node.parent
			}
		}
	}
}

var _ json.Marshaler = (*Tree[string, int])(nil)
var _ json.Unmarshaler = (*Tree[string, int])(nil)

// ToJSON outputs the JSON representation of the tree.
func (tree *Tree[K, V]) MarshalJSON() ([]byte, error) {
	elems := maps.Collect(tree.Iter())

	return json.Marshal(&elems)
}

// FromJSON populates the tree from the input JSON representation.
func (tree *Tree[K, V]) UnmarshalJSON(data []byte) error {
	elems := make(map[K]V)

	err := json.Unmarshal(data, &elems)
	if err != nil {
		return err
	}

	tree.Clear()

	for key, value := range elems {
		tree.Put(key, value)
	}

	return nil
}

// String returns a string representation of the tree.
// Time complexity: O(n).
func (t *Tree[K, V]) String() string {
	if t.IsEmpty() {
		return "AVLTree[]"
	}

	var sb strings.Builder

	sb.WriteString("AVLTree\n")
	t.output(t.root, "", true, &sb)

	return sb.String()
}

// Comparator returns the comparator used by the tree.
// Time complexity: O(1).
func (t *Tree[K, V]) Comparator() cmp.Comparator[K] {
	return t.comparator
}

// --- Internal and Helper Methods ---

// lookup finds the node with the specified key, or nil if not found.
// Time complexity: O(log n).
func (t *Tree[K, V]) lookup(key K) *Node[K, V] {
	node := t.root
	for node != nil {
		switch cmp := t.comparator(key, node.key); {
		case cmp == 0:
			return node
		case cmp < 0:
			node = node.left
		default:
			node = node.right
		}
	}

	return nil
}

// getLeftNode finds the leftmost node in the subtree, or nil if empty.
// Time complexity: O(log n).
func (t *Tree[K, V]) getLeftNode(node *Node[K, V]) *Node[K, V] {
	for node != nil && node.left != nil {
		node = node.left
	}

	return node
}

// getRightNode finds the rightmost node in the subtree, or nil if empty.
// Time complexity: O(log n).
func (t *Tree[K, V]) getRightNode(node *Node[K, V]) *Node[K, V] {
	for node != nil && node.right != nil {
		node = node.right
	}

	return node
}

// replaceNode replaces the old node with the new node in the tree structure.
func (t *Tree[K, V]) replaceNode(old, new *Node[K, V]) {
	if old.parent == nil {
		t.root = new
	} else if old == old.parent.left {
		old.parent.left = new
	} else {
		old.parent.right = new
	}

	if new != nil {
		new.parent = old.parent
	}
}

// --- AVL Balance Logic ---

// insertFixup rebalances the tree upward from the given node after insertion.
func (t *Tree[K, V]) insertFixup(node *Node[K, V]) {
	for node != nil {
		t.updateBalanceFactor(node)

		bf := node.balanceFactor
		if bf < -1 || bf > 1 {
			t.rebalance(node)

			break
		}

		if bf == 0 {
			break
		}

		node = node.parent
	}
}

// deleteFixup rebalances the tree upward from the given node after deletion.
func (t *Tree[K, V]) deleteFixup(node *Node[K, V]) {
	for node != nil {
		t.updateBalanceFactor(node)

		bf := node.balanceFactor
		if bf < -1 || bf > 1 {
			t.rebalance(node)
		}

		if node.balanceFactor != 0 {
			break
		}

		node = node.parent
	}
}

// rebalance balances the subtree rooted at the unbalanced node.
// Assumes the node has a balance factor > 1 or < -1.
func (t *Tree[K, V]) rebalance(node *Node[K, V]) {
	if node.balanceFactor < -1 {
		if node.left.balanceFactor > 0 {
			t.rotateLeft(node.left)
		}

		t.rotateRight(node)
	} else {
		if node.right.balanceFactor < 0 {
			t.rotateRight(node.right)
		}

		t.rotateLeft(node)
	}
}

// rotateLeft performs a left rotation around the pivot node.
func (t *Tree[K, V]) rotateLeft(pivot *Node[K, V]) {
	r := pivot.right
	t.replaceNode(pivot, r)

	pivot.right = r.left
	if pivot.right != nil {
		pivot.right.parent = pivot
	}

	r.left = pivot
	pivot.parent = r

	t.updateBalanceFactor(pivot)
	t.updateBalanceFactor(r)
}

// rotateRight performs a right rotation around the pivot node.
func (t *Tree[K, V]) rotateRight(pivot *Node[K, V]) {
	l := pivot.left
	t.replaceNode(pivot, l)

	pivot.left = l.right
	if pivot.left != nil {
		pivot.left.parent = pivot
	}

	l.right = pivot
	pivot.parent = l

	t.updateBalanceFactor(pivot)
	t.updateBalanceFactor(l)
}

// height returns the height of a node. A nil node has height -1.
func (t *Tree[K, V]) height(n *Node[K, V]) int {
	if n == nil {
		return -1
	}

	return 1 + max(t.height(n.left), t.height(n.right))
}

// updateBalanceFactor recalculates and updates the balance factor of a node.
func (t *Tree[K, V]) updateBalanceFactor(n *Node[K, V]) {
	if n == nil {
		return
	}

	n.balanceFactor = int8(t.height(n.right) - t.height(n.left))
}

// --- Utility Functions ---

// inOrderTraversal traverses the tree in-order, applying function f to each node.
func (t *Tree[K, V]) inOrderTraversal(f func(*Node[K, V])) {
	var visit func(*Node[K, V])
	visit = func(n *Node[K, V]) {
		if n == nil {
			return
		}

		visit(n.left)
		f(n)
		visit(n.right)
	}
	visit(t.root)
}

// output recursively builds a string representation of the tree for printing.
func (t *Tree[K, V]) output(node *Node[K, V], prefix string, isTail bool, sb *strings.Builder) {
	if node.right != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}

		t.output(node.right, newPrefix, false, sb)
	}

	sb.WriteString(prefix)

	if isTail {
		sb.WriteString("└── ")
	} else {
		sb.WriteString("┌── ")
	}

	sb.WriteString(node.String() + "\n")

	if node.left != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}

		t.output(node.left, newPrefix, true, sb)
	}
}

// cloneNode creates a deep copy of a node and its subtree, setting the parent for the new node.
func cloneNode[K comparable, V any](node *Node[K, V], parent *Node[K, V]) *Node[K, V] {
	if node == nil {
		return nil
	}

	n := &Node[K, V]{
		key:           node.key,
		value:         node.value,
		balanceFactor: node.balanceFactor,
		parent:        parent,
	}

	n.left = cloneNode(node.left, n)
	n.right = cloneNode(node.right, n)

	return n
}
