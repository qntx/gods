// Package rbtree implements a red-black tree for ordered key-value storage.
//
// It provides a self-balancing binary search tree with O(log n) operations for
// insertion, deletion, and lookup. Used by TreeSet and TreeMap. Not thread-safe.
//
// Reference: https://en.wikipedia.org/wiki/Red%E2%80%93black_tree
package rbtree

import (
	"cmp"
	"errors"
	"fmt"
	"strings"

	"github.com/qntx/gods/util"
)

// --------------------------------------------------------------------------------
// Constants and Errors

// Color constants for red-black tree nodes.
const (
	black color = true  // Represents a black node.
	red   color = false // Represents a red node.
)

// Predefined errors for tree operations.
var (
	ErrInvalidKeyType = errors.New("key type does not match comparator")
)

// --------------------------------------------------------------------------------
// Types

// color represents the color of a red-black tree node (red or black).
type color bool

// Tree manages a red-black tree with key-value pairs.
//
// K must be comparable and compatible with the provided comparator.
// V can be any type.
type Tree[K comparable, V any] struct {
	Root       *Node[K, V]        // Root node of the tree.
	len        int                // Number of nodes in the tree.
	Comparator util.Comparator[K] // Comparator for ordering keys.
}

// Node represents a single element in the red-black tree.
type Node[K comparable, V any] struct {
	Key    K           // Key for ordering.
	Value  V           // Associated value.
	color  color       // Node color (red or black).
	Left   *Node[K, V] // Left child.
	Right  *Node[K, V] // Right child.
	Parent *Node[K, V] // Parent node.
}

// --------------------------------------------------------------------------------
// Constructors

// New creates a new red-black tree with the built-in comparator for ordered types.
//
// K must implement cmp.Ordered (e.g., int, string). Time complexity: O(1).
func New[K cmp.Ordered, V any]() *Tree[K, V] {
	return &Tree[K, V]{Comparator: cmp.Compare[K]}
}

// NewWith creates a new red-black tree with a custom comparator.
//
// The comparator defines the ordering of keys. Time complexity: O(1).
func NewWith[K comparable, V any](comparator util.Comparator[K]) *Tree[K, V] {
	return &Tree[K, V]{Comparator: comparator}
}

// --------------------------------------------------------------------------------
// Public Methods

// Put inserts or updates a key-value pair in the tree.
//
// If the key exists, its value is updated; otherwise, a new node is inserted.
// Panics if the key type is incompatible with the comparator.
// Time complexity: O(log n).
func (t *Tree[K, V]) Put(key K, val V) {
	t.validateKey(key)

	if t.Root == nil {
		t.Root = &Node[K, V]{Key: key, Value: val, color: black}
		t.len++

		return
	}

	node, parent := t.Root, (*Node[K, V])(nil)
	for node != nil {
		parent = node

		switch cmp := t.Comparator(key, node.Key); {
		case cmp == 0:
			node.Value = val

			return
		case cmp < 0:
			node = node.Left
		default:
			node = node.Right
		}
	}

	newNode := &Node[K, V]{Key: key, Value: val, color: red, Parent: parent}
	if t.Comparator(key, parent.Key) < 0 {
		parent.Left = newNode
	} else {
		parent.Right = newNode
	}

	t.insertFixup(newNode)

	t.len++
}

// Get retrieves the value associated with the given key.
//
// Returns the value and true if found, zero value and false otherwise.
// Panics if the key type is incompatible with the comparator.
// Time complexity: O(log n).
func (t *Tree[K, V]) Get(key K) (val V, found bool) {
	if node := t.lookup(key); node != nil {
		return node.Value, true
	}

	return val, false
}

// GetNode retrieves the node associated with the given key.
//
// Returns the node if found, nil otherwise. Panics if the key type is
// incompatible with the comparator. Time complexity: O(log n).
func (t *Tree[K, V]) GetNode(key K) *Node[K, V] {
	return t.lookup(key)
}

// Remove deletes the node with the given key from the tree.
//
// Does nothing if the key is not found. Panics if the key type is incompatible
// with the comparator. Time complexity: O(log n).
func (t *Tree[K, V]) Remove(key K) {
	node := t.lookup(key)
	if node == nil {
		return
	}

	t.deleteNode(node)

	t.len--
}

// Empty checks if the tree contains no nodes.
//
// Time complexity: O(1).
func (t *Tree[K, V]) Empty() bool {
	return t.len == 0
}

// Len returns the number of nodes in the tree.
//
// Time complexity: O(1).
func (t *Tree[K, V]) Len() int {
	return t.len
}

// Size returns the number of nodes in the subtree rooted at this node.
//
// Computed dynamically by traversing the subtree. Time complexity: O(n).
func (n *Node[K, V]) Size() int {
	if n == nil {
		return 0
	}

	return 1 + n.Left.Size() + n.Right.Size()
}

// Keys returns all keys in in-order traversal.
//
// Time complexity: O(n).
func (t *Tree[K, V]) Keys() []K {
	keys := make([]K, t.len)
	it := t.Iterator()

	for i := 0; it.Next(); i++ {
		keys[i] = it.Key()
	}

	return keys
}

// Values returns all values in in-order traversal based on keys.
//
// Time complexity: O(n).
func (t *Tree[K, V]) Values() []V {
	vals := make([]V, t.len)
	it := t.Iterator()

	for i := 0; it.Next(); i++ {
		vals[i] = it.Value()
	}

	return vals
}

// KeysAndValues returns all keys and values in in-order traversal.
//
// More efficient than calling Keys() and Values() separately as it traverses
// the tree only once. Time complexity: O(n).
func (t *Tree[K, V]) KeysAndValues() ([]K, []V) {
	keys := make([]K, t.len)
	vals := make([]V, t.len)
	it := t.Iterator()

	for i := 0; it.Next(); i++ {
		keys[i], vals[i] = it.Key(), it.Value()
	}

	return keys, vals
}

// Left returns the leftmost (minimum) node or nil if the tree is empty.
//
// Time complexity: O(log n).
func (t *Tree[K, V]) Left() *Node[K, V] {
	return t.minNode(t.Root)
}

// Right returns the rightmost (maximum) node or nil if the tree is empty.
//
// Time complexity: O(log n).
func (t *Tree[K, V]) Right() *Node[K, V] {
	return t.maxNode(t.Root)
}

// Floor finds the largest node less than or equal to the given key.
//
// Returns the node and true if found, nil and false otherwise. Panics if the
// key type is incompatible with the comparator. Time complexity: O(log n).
func (t *Tree[K, V]) Floor(key K) (*Node[K, V], bool) {
	t.validateKey(key)

	var floor *Node[K, V]

	node := t.Root
	for node != nil {
		switch cmp := t.Comparator(key, node.Key); {
		case cmp == 0:
			return node, true
		case cmp > 0:
			floor = node
			node = node.Right
		default:
			node = node.Left
		}
	}

	return floor, floor != nil
}

// Ceiling finds the smallest node greater than or equal to the given key.
//
// Returns the node and true if found, nil and false otherwise. Panics if the
// key type is incompatible with the comparator. Time complexity: O(log n).
func (t *Tree[K, V]) Ceiling(key K) (*Node[K, V], bool) {
	t.validateKey(key)

	var ceil *Node[K, V]

	node := t.Root
	for node != nil {
		switch cmp := t.Comparator(key, node.Key); {
		case cmp == 0:
			return node, true
		case cmp < 0:
			ceil = node
			node = node.Left
		default:
			node = node.Right
		}
	}

	return ceil, ceil != nil
}

// Clear removes all nodes from the tree.
//
// Time complexity: O(1).
func (t *Tree[K, V]) Clear() {
	t.Root = nil
	t.len = 0
}

// String returns a string representation of the tree.
//
// Time complexity: O(n).
func (t *Tree[K, V]) String() string {
	if t.Empty() {
		return "RedBlackTree[]"
	}

	var sb strings.Builder

	sb.WriteString("RedBlackTree\n")
	t.output(t.Root, "", true, &sb)

	return sb.String()
}

// String returns a string representation of the node.
//
// Time complexity: O(1).
func (n *Node[K, V]) String() string {
	return fmt.Sprintf("%v", n.Key)
}

// --------------------------------------------------------------------------------
// Private Methods

// validateKey ensures the key is compatible with the comparator.
//
// Panics if the key type does not match the comparator's expectations.
func (t *Tree[K, V]) validateKey(key K) {
	if _, err := safeCompare(t.Comparator, key, key); err != nil {
		panic(fmt.Sprintf("rbtree: %v", err))
	}
}

// lookup finds the node with the given key.
//
// Returns nil if not found. Time complexity: O(log n).
func (t *Tree[K, V]) lookup(key K) *Node[K, V] {
	t.validateKey(key)

	node := t.Root
	for node != nil {
		switch cmp := t.Comparator(key, node.Key); {
		case cmp == 0:
			return node
		case cmp < 0:
			node = node.Left
		default:
			node = node.Right
		}
	}

	return nil
}

// minNode finds the leftmost node in the subtree.
//
// Returns nil if the subtree is empty. Time complexity: O(log n).
func (t *Tree[K, V]) minNode(node *Node[K, V]) *Node[K, V] {
	for node != nil && node.Left != nil {
		node = node.Left
	}

	return node
}

// maxNode finds the rightmost node in the subtree.
//
// Returns nil if the subtree is empty. Time complexity: O(log n).
func (t *Tree[K, V]) maxNode(node *Node[K, V]) *Node[K, V] {
	for node != nil && node.Right != nil {
		node = node.Right
	}

	return node
}

// output builds a string representation of the tree recursively.
func (t *Tree[K, V]) output(node *Node[K, V], prefix string, isTail bool, sb *strings.Builder) {
	if node.Right != nil {
		newPrefix := prefix + ternary(isTail, "│   ", "    ")
		t.output(node.Right, newPrefix, false, sb)
	}

	sb.WriteString(prefix)
	sb.WriteString(ternary(isTail, "└── ", "┌── "))
	sb.WriteString(node.String() + "\n")

	if node.Left != nil {
		newPrefix := prefix + ternary(isTail, "    ", "│   ")
		t.output(node.Left, newPrefix, true, sb)
	}
}

// grandparent returns the grandparent of the node.
//
// Returns nil if not applicable. Time complexity: O(1).
func (n *Node[K, V]) grandparent() *Node[K, V] {
	if n != nil && n.Parent != nil {
		return n.Parent.Parent
	}

	return nil
}

// uncle returns the uncle of the node.
//
// Returns nil if not applicable. Time complexity: O(1).
func (n *Node[K, V]) uncle() *Node[K, V] {
	if gp := n.grandparent(); gp != nil {
		if n.Parent == gp.Left {
			return gp.Right
		}

		return gp.Left
	}

	return nil
}

// sibling returns the sibling of the node.
//
// Returns nil if not applicable. Time complexity: O(1).
func (n *Node[K, V]) sibling() *Node[K, V] {
	if n != nil && n.Parent != nil {
		if n == n.Parent.Left {
			return n.Parent.Right
		}

		return n.Parent.Left
	}

	return nil
}

// rotateLeft performs a left rotation around the node.
func (t *Tree[K, V]) rotateLeft(n *Node[K, V]) {
	r := n.Right
	t.replaceNode(n, r)

	n.Right = r.Left
	if r.Left != nil {
		r.Left.Parent = n
	}

	r.Left = n
	n.Parent = r
}

// rotateRight performs a right rotation around the node.
func (t *Tree[K, V]) rotateRight(n *Node[K, V]) {
	l := n.Left
	t.replaceNode(n, l)

	n.Left = l.Right
	if l.Right != nil {
		l.Right.Parent = n
	}

	l.Right = n
	n.Parent = l
}

// replaceNode replaces oldNode with newNode in the tree structure.
func (t *Tree[K, V]) replaceNode(oldNode, newNode *Node[K, V]) {
	if oldNode.Parent == nil {
		t.Root = newNode
	} else if oldNode == oldNode.Parent.Left {
		oldNode.Parent.Left = newNode
	} else {
		oldNode.Parent.Right = newNode
	}

	if newNode != nil {
		newNode.Parent = oldNode.Parent
	}
}

// insertFixup balances the tree after insertion.
func (t *Tree[K, V]) insertFixup(n *Node[K, V]) {
	if n.Parent == nil {
		n.color = black

		return
	}

	if nodeColor(n.Parent) == black {
		return
	}

	if uncle := n.uncle(); nodeColor(uncle) == red {
		n.Parent.color = black
		uncle.color = black
		gp := n.grandparent()
		gp.color = red
		t.insertFixup(gp)

		return
	}

	t.insertFixupStep(n)
}

// insertFixupStep handles rotation cases for insertion balancing.
func (t *Tree[K, V]) insertFixupStep(n *Node[K, V]) {
	gp := n.grandparent()
	if n == n.Parent.Right && n.Parent == gp.Left {
		t.rotateLeft(n.Parent)
		n = n.Left
	} else if n == n.Parent.Left && n.Parent == gp.Right {
		t.rotateRight(n.Parent)
		n = n.Right
	}

	n.Parent.color = black
	gp.color = red

	if n == n.Parent.Left {
		t.rotateRight(gp)
	} else {
		t.rotateLeft(gp)
	}
}

// deleteNode removes a node from the tree and rebalances.
func (t *Tree[K, V]) deleteNode(n *Node[K, V]) {
	var child *Node[K, V]

	if n.Left != nil && n.Right != nil {
		pred := t.maxNode(n.Left)
		n.Key, n.Value = pred.Key, pred.Value
		n = pred
	}

	child = ternary(n.Left == nil, n.Right, n.Left)
	if n.color == black {
		n.color = nodeColor(child)
		t.deleteFixup(n)
	}

	t.replaceNode(n, child)

	if n.Parent == nil && child != nil {
		child.color = black
	}
}

// deleteFixup balances the tree after deletion.
func (t *Tree[K, V]) deleteFixup(n *Node[K, V]) {
	if n.Parent == nil {
		return
	}

	s := n.sibling()
	if nodeColor(s) == red {
		n.Parent.color = red
		s.color = black

		if n == n.Parent.Left {
			t.rotateLeft(n.Parent)
		} else {
			t.rotateRight(n.Parent)
		}

		s = n.sibling()
	}

	t.deleteFixupCases(n, s)
}

// deleteFixupCases handles specific deletion balancing cases.
func (t *Tree[K, V]) deleteFixupCases(n, s *Node[K, V]) {
	if nodeColor(n.Parent) == black && nodeColor(s) == black &&
		nodeColor(s.Left) == black && nodeColor(s.Right) == black {
		s.color = red

		t.deleteFixup(n.Parent)

		return
	}

	if nodeColor(n.Parent) == red && nodeColor(s) == black &&
		nodeColor(s.Left) == black && nodeColor(s.Right) == black {
		s.color = red
		n.Parent.color = black

		return
	}

	t.deleteFixupRotations(n, s)
}

// deleteFixupRotations handles rotation cases for deletion balancing.
func (t *Tree[K, V]) deleteFixupRotations(n, s *Node[K, V]) {
	if n == n.Parent.Left && nodeColor(s) == black &&
		nodeColor(s.Left) == red && nodeColor(s.Right) == black {
		s.color = red
		s.Left.color = black
		t.rotateRight(s)
		s = n.sibling()
	} else if n == n.Parent.Right && nodeColor(s) == black &&
		nodeColor(s.Right) == red && nodeColor(s.Left) == black {
		s.color = red
		s.Right.color = black
		t.rotateLeft(s)
		s = n.sibling()
	}

	s.color = nodeColor(n.Parent)
	n.Parent.color = black

	if n == n.Parent.Left {
		s.Right.color = black

		t.rotateLeft(n.Parent)
	} else {
		s.Left.color = black

		t.rotateRight(n.Parent)
	}
}

// nodeColor returns the color of a node, black if nil.
func nodeColor[K comparable, V any](n *Node[K, V]) color {
	if n == nil {
		return black
	}

	return n.color
}

// safeCompare wraps a comparator call with error handling.
//
// Returns the comparison result and any error from a panic.
func safeCompare[K comparable](cmp util.Comparator[K], a, b K) (int, error) {
	defer func() {
		if r := recover(); r != nil {
			err, _ := r.(error)
			panic(fmt.Errorf("%w: %w", ErrInvalidKeyType, err))
		}
	}()

	return cmp(a, b), nil
}

// ternary is a helper for conditional expressions.
func ternary[T any](cond bool, trueVal, falseVal T) T {
	if cond {
		return trueVal
	}

	return falseVal
}
