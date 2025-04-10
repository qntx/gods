// Package btree provides an iterator for traversing a generic B-tree data structure.
package btree

import "github.com/qntx/gods/container"

// --------------------------------------------------------------------------------
// Types and Constants

// Iterator represents a stateful iterator for a B-tree, traversing key-value pairs.
//
// Type parameters:
//   - K: the key type, must be comparable
//   - V: the value type
type Iterator[K comparable, V any] struct {
	tree     *Tree[K, V]
	node     *Node[K, V]
	entry    *Entry[K, V]
	position position
}

// position defines the iterator's current state.
type position byte

const (
	// begin represents the position before the first element.
	begin position = iota
	// between represents the position at a valid element.
	between
	// end represents the position after the last element.
	end
)

// Verify Iterator implements required interface at compile time.
var _ container.ReverseIteratorWithKey[string, int] = (*Iterator[string, int])(nil)

// --------------------------------------------------------------------------------
// Iterator Creation

// Iterator returns a new iterator for traversing the tree's key-value pairs.
//
// The iterator starts at the "begin" position; call Next() to move to the first element.
//
// Returns:
//
//	A pointer to an initialized Iterator.
func (t *Tree[K, V]) Iterator() *Iterator[K, V] {
	return &Iterator[K, V]{tree: t, position: begin}
}

// --------------------------------------------------------------------------------
// Navigation Methods

// Next advances the iterator to the next element.
//
// If there is a next element, returns true and updates the iterator to point to it.
// Otherwise, moves to the "end" position and returns false. The current key and value
// can be retrieved with Key() and Value() when true is returned.
func (it *Iterator[K, V]) Next() bool {
	if it.position == end {
		it.End()

		return false
	}

	if it.position == begin {
		left := it.tree.Left()
		if left == nil {
			it.End()

			return false
		}

		it.node = left
		it.entry = left.Entries[0]
		it.position = between

		return true
	}

	e, _ := it.tree.search(it.node, it.entry.Key)
	if e+1 < len(it.node.Children) {
		it.node = it.node.Children[e+1]
		for len(it.node.Children) > 0 {
			it.node = it.node.Children[0]
		}

		it.entry = it.node.Entries[0]
		it.position = between

		return true
	}

	if e+1 < len(it.node.Entries) {
		it.entry = it.node.Entries[e+1]
		it.position = between

		return true
	}

	for it.node.Parent != nil {
		it.node = it.node.Parent
		e, _ := it.tree.search(it.node, it.entry.Key)

		if e < len(it.node.Entries) {
			it.entry = it.node.Entries[e]
			it.position = between

			return true
		}
	}

	it.End()

	return false
}

// Prev moves the iterator to the previous element.
//
// If there is a previous element, returns true and updates the iterator to point to it.
// Otherwise, moves to the "begin" position and returns false. The current key and value
// can be retrieved with Key() and Value() when true is returned.
func (it *Iterator[K, V]) Prev() bool {
	if it.position == begin {
		it.Begin()

		return false
	}

	if it.position == end {
		right := it.tree.Right()
		if right == nil {
			it.Begin()

			return false
		}

		it.node = right
		it.entry = right.Entries[len(right.Entries)-1]
		it.position = between

		return true
	}

	e, _ := it.tree.search(it.node, it.entry.Key)
	if e < len(it.node.Children) {
		it.node = it.node.Children[e]
		for len(it.node.Children) > 0 {
			it.node = it.node.Children[len(it.node.Children)-1]
		}

		it.entry = it.node.Entries[len(it.node.Entries)-1]
		it.position = between

		return true
	}

	if e-1 >= 0 {
		it.entry = it.node.Entries[e-1]
		it.position = between

		return true
	}

	for it.node.Parent != nil {
		it.node = it.node.Parent
		e, _ := it.tree.search(it.node, it.entry.Key)

		if e-1 >= 0 {
			it.entry = it.node.Entries[e-1]
			it.position = between

			return true
		}
	}

	it.Begin()

	return false
}

// --------------------------------------------------------------------------------
// Accessor Methods

// Key returns the current element's key.
//
// Panics if the iterator is not at a valid position (between elements).
func (it *Iterator[K, V]) Key() K {
	if it.position != between || it.entry == nil {
		panic("btree: iterator not at valid position")
	}

	return it.entry.Key
}

// Value returns the current element's value.
//
// Panics if the iterator is not at a valid position (between elements).
func (it *Iterator[K, V]) Value() V {
	if it.position != between || it.entry == nil {
		panic("btree: iterator not at valid position")
	}

	return it.entry.Value
}

// Node returns the current element's node.
//
// Returns nil if the iterator is not at a valid position.
func (it *Iterator[K, V]) Node() *Node[K, V] {
	if it.position != between {
		return nil
	}

	return it.node
}

// --------------------------------------------------------------------------------
// Positioning Methods

// Begin resets the iterator to the position before the first element.
//
// Call Next() to move to the first element.
func (it *Iterator[K, V]) Begin() {
	it.node = nil
	it.entry = nil
	it.position = begin
}

// End moves the iterator to the position after the last element.
//
// Call Prev() to move to the last element.
func (it *Iterator[K, V]) End() {
	it.node = nil
	it.entry = nil
	it.position = end
}

// First moves the iterator to the first element.
//
// Returns true if the tree is not empty, false otherwise.
func (it *Iterator[K, V]) First() bool {
	it.Begin()

	return it.Next()
}

// Last moves the iterator to the last element.
//
// Returns true if the tree is not empty, false otherwise.
func (it *Iterator[K, V]) Last() bool {
	it.End()

	return it.Prev()
}

// NextTo advances the iterator to the next element satisfying the given predicate.
//
// Returns true if such an element is found, false if the end is reached.
// The predicate function takes the key and value as arguments.
func (it *Iterator[K, V]) NextTo(f func(key K, value V) bool) bool {
	for it.Next() {
		if f(it.Key(), it.Value()) {
			return true
		}
	}

	return false
}

// PrevTo moves the iterator to the previous element satisfying the given predicate.
//
// Returns true if such an element is found, false if the beginning is reached.
// The predicate function takes the key and value as arguments.
func (it *Iterator[K, V]) PrevTo(f func(key K, value V) bool) bool {
	for it.Prev() {
		if f(it.Key(), it.Value()) {
			return true
		}
	}

	return false
}
