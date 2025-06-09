// Package rbtree provides an iterator for traversing the red-black tree.
//
// This file implements a stateful iterator for the Tree type, supporting both
// forward and reverse iteration over key-value pairs with O(log n) access time.
package rbtree

import (
	"errors"

	"github.com/qntx/gods/container"
)

// Position constants for iterator state.
type position byte

const (
	begin   position = iota // Before the first element.
	between                 // Between elements (valid position).
	end                     // Past the last element.
)

// Predefined errors for iterator operations.
var (
	ErrInvalidIteratorPosition = errors.New("iterator accessed at invalid position")
)

// Ensure Iterator implements container.ReverseIteratorWithKey at compile time.
var _ container.ReverseIteratorWithKey[string, int] = (*Iterator[string, int])(nil)

// Iterator provides forward and reverse traversal over a Tree's key-value pairs.
//
// It maintains a position in the tree, allowing O(log n) navigation between elements.
// The iterator is read-only and does not modify the underlying tree. Most operations
// are O(log n) due to tree traversal, except Begin and End which are O(1).
type Iterator[K comparable, V any] struct {
	tree     *Tree[K, V] // Reference to the tree being iterated.
	node     *Node[K, V] // Current node, nil if at begin or end.
	position position    // Current state: begin, between, or end.
}

// Iterator creates a new iterator for the tree.
//
// Starts before the first element (begin state). Use Next() to reach the first
// element or End() followed by Prev() for the last. Time complexity: O(1).
func (t *Tree[K, V]) Iterator() *Iterator[K, V] {
	return &Iterator[K, V]{tree: t, node: nil, position: begin}
}

// IteratorAt creates a new iterator starting at a specific node.
//
// Starts in the between state at the given node. Time complexity: O(1).
func (t *Tree[K, V]) IteratorAt(node *Node[K, V]) *Iterator[K, V] {
	return &Iterator[K, V]{tree: t, node: node, position: between}
}

// Next advances the iterator to the next element in in-order traversal.
//
// Returns true if the iterator is at a valid element after moving, false if it
// reaches the end. Updates the iterator's state. Time complexity: O(log n).
func (it *Iterator[K, V]) Next() bool {
	switch it.position {
	case end:
		return false
	case begin:
		if left := it.tree.GetLeftNode(); left != nil {
			it.node = left
			it.position = between

			return true
		}

		it.position = end

		return false
	case between:
		if it.node.Right != nil {
			it.node = it.node.Right
			for it.node.Left != nil {
				it.node = it.node.Left
			}

			return true
		}

		for it.node.Parent != nil {
			node := it.node
			it.node = it.node.Parent

			if node == it.node.Left {
				return true
			}
		}
	}

	it.node = nil
	it.position = end

	return false
}

// Prev moves the iterator to the previous element in in-order traversal.
//
// Returns true if the iterator is at a valid element after moving, false if it
// reaches the beginning. Updates the iterator's state. Time complexity: O(log n).
func (it *Iterator[K, V]) Prev() bool {
	switch it.position {
	case begin:
		return false
	case end:
		if right := it.tree.GetRightNode(); right != nil {
			it.node = right
			it.position = between

			return true
		}

		it.position = begin

		return false
	case between:
		if it.node.Left != nil {
			it.node = it.node.Left
			for it.node.Right != nil {
				it.node = it.node.Right
			}

			return true
		}

		for it.node.Parent != nil {
			node := it.node
			it.node = it.node.Parent

			if node == it.node.Right {
				return true
			}
		}
	}

	it.node = nil
	it.position = begin

	return false
}

// Key returns the current element's key.
//
// Panics if the iterator is not at a valid position (begin or end state).
// Time complexity: O(1).
func (it *Iterator[K, V]) Key() K {
	if !it.valid() {
		panic("rbtree: " + ErrInvalidIteratorPosition.Error())
	}

	return it.node.Key
}

// Value returns the current element's value.
//
// Panics if the iterator is not at a valid position (begin or end state).
// Time complexity: O(1).
func (it *Iterator[K, V]) Value() V {
	if !it.valid() {
		panic("rbtree: " + ErrInvalidIteratorPosition.Error())
	}

	return it.node.Value
}

// Node returns the current node.
//
// Returns nil if the iterator is at begin or end. Time complexity: O(1).
func (it *Iterator[K, V]) Node() *Node[K, V] {
	return it.node
}

// Begin resets the iterator to before the first element.
//
// Use Next() to move to the first element. Time complexity: O(1).
func (it *Iterator[K, V]) Begin() {
	it.node = nil
	it.position = begin
}

// End moves the iterator past the last element.
//
// Use Prev() to move to the last element. Time complexity: O(1).
func (it *Iterator[K, V]) End() {
	it.node = nil
	it.position = end
}

// First moves the iterator to the first element.
//
// Returns true if the tree is non-empty, false otherwise. Time complexity: O(log n).
func (it *Iterator[K, V]) First() bool {
	it.Begin()

	return it.Next()
}

// Last moves the iterator to the last element.
//
// Returns true if the tree is non-empty, false otherwise. Time complexity: O(log n).
func (it *Iterator[K, V]) Last() bool {
	it.End()

	return it.Prev()
}

// NextTo advances to the next element satisfying the given condition.
//
// Moves forward from the current position until an element matches the predicate
// or the end is reached. Returns true if a match is found. Time complexity: O(n)
// in the worst case.
func (it *Iterator[K, V]) NextTo(f func(key K, value V) bool) bool {
	for it.Next() {
		if f(it.Key(), it.Value()) {
			return true
		}
	}

	return false
}

// PrevTo moves to the previous element satisfying the given condition.
//
// Moves backward from the current position until an element matches the predicate
// or the beginning is reached. Returns true if a match is found. Time complexity: O(n)
// in the worst case.
func (it *Iterator[K, V]) PrevTo(f func(key K, value V) bool) bool {
	for it.Prev() {
		if f(it.Key(), it.Value()) {
			return true
		}
	}

	return false
}

// valid checks if the iterator is at a valid element position.
func (it *Iterator[K, V]) valid() bool {
	return it.position == between && it.node != nil
}
