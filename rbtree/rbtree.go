// Package rbtree implements a red-black tree for ordered key-value storage.
//
// It provides a self-balancing binary search tree with O(log n) operations for
// insertion, deletion, and lookup. Used by TreeSet and TreeMap. Not thread-safe.
//
// Reference: https://en.wikipedia.org/wiki/Red%E2%80%93black_tree
package rbtree

import (
	"fmt"
	"strings"

	"github.com/qntx/gods/cmp"
)

// Color represents the color of a red-black tree node (red or black).
type Color bool

const (
	black Color = true  // Represents a black node.
	red   Color = false // Represents a red node.
)

// Node represents a single element in the red-black tree.
type Node[K comparable, V any] struct {
	color  Color       // Node color (red or black).
	key    K           // Key for ordering.
	value  V           // Associated value.
	parent *Node[K, V] // Parent node.
	left   *Node[K, V] // Left child.
	right  *Node[K, V] // Right child.
}

// Color returns the color of the node (red or black).
// This method provides read-only access to the node's color.
// Time complexity: O(1).
func (n *Node[K, V]) Color() Color {
	return n.color
}

// Key returns the key stored in the node.
// This key is used for ordering within the red-black tree.
// This method provides read-only access to the node's key.
// Time complexity: O(1).
func (n *Node[K, V]) Key() K {
	return n.key
}

// Value returns the value associated with the key in the node.
// This method provides read-only access to the node's value.
// Time complexity: O(1).
func (n *Node[K, V]) Value() V {
	return n.value
}

// Left returns the left child of the node.
// It returns nil if the node has no left child.
// This method provides read-only access to the node's left child.
// Time complexity: O(1).
func (n *Node[K, V]) Left() *Node[K, V] {
	return n.left
}

// Right returns the right child of the node.
// It returns nil if the node has no right child.
// This method provides read-only access to the node's right child.
// Time complexity: O(1).
func (n *Node[K, V]) Right() *Node[K, V] {
	return n.right
}

// Parent returns the parent of the node.
// It returns nil if the node is the root or has no parent.
// This method provides read-only access to the node's parent.
// Time complexity: O(1).
func (n *Node[K, V]) Parent() *Node[K, V] {
	return n.parent
}

// Size returns the number of nodes in the subtree rooted at this node.
//
// Computed dynamically by traversing the subtree. Time complexity: O(n).
func (n *Node[K, V]) Size() int {
	if n == nil {
		return 0
	}

	return 1 + n.left.Size() + n.right.Size()
}

// String returns a string representation of the node.
//
// Time complexity: O(1).
func (n *Node[K, V]) String() string {
	return fmt.Sprintf("%v", n.key)
}

// uncle returns the uncle of the node.
//
// Returns nil if not applicable. Time complexity: O(1).
func (n *Node[K, V]) uncle() *Node[K, V] {
	if gp := n.grandparent(); gp != nil {
		if n.parent == gp.left {
			return gp.right
		}

		return gp.left
	}

	return nil
}

// sibling returns the sibling of the node.
//
// Returns nil if not applicable. Time complexity: O(1).
func (n *Node[K, V]) sibling() *Node[K, V] {
	if n != nil && n.parent != nil {
		if n == n.parent.left {
			return n.parent.right
		}

		return n.parent.left
	}

	return nil
}

// grandparent returns the grandparent of the node.
//
// Returns nil if not applicable. Time complexity: O(1).
func (n *Node[K, V]) grandparent() *Node[K, V] {
	if n != nil && n.parent != nil {
		return n.parent.parent
	}

	return nil
}

// Tree manages a red-black tree with key-value pairs.
//
// K must be comparable and compatible with the provided comparator.
// V can be any type.
type Tree[K comparable, V any] struct {
	root       *Node[K, V]       // Root node of the tree.
	len        int               // Number of nodes in the tree.
	comparator cmp.Comparator[K] // Comparator for ordering keys.
}

// New creates a new red-black tree with the built-in comparator for ordered types.
//
// K must implement cmp.Ordered (e.g., int, string). Time complexity: O(1).
func New[K cmp.Ordered, V any]() *Tree[K, V] {
	return &Tree[K, V]{comparator: cmp.GenericComparator[K]}
}

// NewWith creates a new red-black tree with a custom comparator.
//
// The comparator defines the ordering of keys. Time complexity: O(1).
func NewWith[K comparable, V any](comparator cmp.Comparator[K]) *Tree[K, V] {
	return &Tree[K, V]{comparator: comparator}
}

// Put inserts or updates a key-value pair in the tree.
//
// If the key exists, its value is updated; otherwise, a new node is inserted.
// Panics if the key type is incompatible with the comparator.
// Time complexity: O(log n).
func (t *Tree[K, V]) Put(key K, val V) {
	// Case 1: Tree is empty.
	// The new node becomes the root and is colored black (Property 2).
	if t.root == nil {
		t.root = &Node[K, V]{key: key, value: val, color: black}
		t.len++

		return
	}

	// Case 2: Tree is not empty.
	// Traverse the tree to find the insertion point or an existing node with the same key.
	node, parent := t.root, (*Node[K, V])(nil) // `node` is current, `parent` trails `node`.
	for node != nil {
		parent = node // `parent` will be the parent of the new node if key is not found.

		switch cmp := t.comparator(key, node.key); {
		case cmp == 0:
			// Key already exists, update its value.
			node.value = val

			return
		case cmp < 0:
			// Key is less than current node's key, go left.
			node = node.left
		default: // cmp > 0
			// Key is greater than current node's key, go right.
			node = node.right
		}
	}

	// Key not found, insert a new node.
	// New nodes are initially colored red to simplify maintaining Red-Black properties.
	// The `parent` variable now holds the parent of the new node.
	n := &Node[K, V]{key: key, value: val, color: red, parent: parent}

	// Link the new node to its parent.
	if t.comparator(key, parent.key) < 0 {
		parent.left = n
	} else {
		parent.right = n
	}

	// Rebalance the tree to maintain Red-Black properties after insertion.
	t.insertFixup(n)

	t.len++ // Increment the tree size.
}

// Remove deletes the node with the given key from the tree.
//
// Does nothing if key not found. Panics if key type is incompatible with comparator.
// Time complexity: O(log n).
func (t *Tree[K, V]) Remove(key K) {
	// Step 1: Find node to remove.
	n := t.lookup(key)
	if n == nil {
		return // Not found.
	}

	// unlink: node to be unlinked.
	// child: node that replaces unlink.
	var child *Node[K, V]
	unlink := n

	// Step 2: If unlink has two children, copy predecessor's data to unlink,
	// then target predecessor for unlinking. Predecessor has at most one child.
	if unlink.left != nil && unlink.right != nil {
		pred := t.getRightNode(unlink.left) // In-order predecessor.
		unlink.key, unlink.value = pred.key, pred.value
		unlink = pred // Unlink predecessor.
	}

	// Step 3: unlink now has 0 or 1 child.
	// Set child to unlink's only child (right if left nil, else left) or nil if leaf.
	child = ternary(unlink.left == nil, unlink.right, unlink.left)

	// Step 4: If unlink was black, fix Red-Black properties (black-height, red node rules).
	if unlink.color == black {
		// Set unlink's color like child's for fixup:
		// - RED child: unlink RED, fixup makes BLACK, absorbs extra black.
		// - BLACK/nil child: unlink BLACK, fixup handles black-height deficit.
		unlink.color = color(child)
		t.deleteFixup(unlink) // Fix Red-Black balance at unlink's position.
	}

	// Step 5: Replace unlink with child in the tree.
	t.replaceNode(unlink, child)

	// Step 6: Ensure new root (if any) is black.
	// deleteFixup usually handles this, but covers red node deletion case.
	if unlink.parent == nil && child != nil {
		child.color = black
	}

	// Step 7: Decrement tree size.
	t.len--
}

// Get retrieves the value associated with the given key.
//
// Returns the value and true if found, zero value and false otherwise.
// Panics if the key type is incompatible with the comparator.
// Time complexity: O(log n).
func (t *Tree[K, V]) Get(key K) (val V, ok bool) {
	if node := t.lookup(key); node != nil {
		return node.value, true
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

// GetLeftNode returns the leftmost (minimum key) node or nil if the tree is empty.
// Renamed from MinNode for clarity.
// Time complexity: O(log n).
func (t *Tree[K, V]) GetLeftNode() *Node[K, V] {
	return t.getLeftNode(t.root)
}

// GetRightNode returns the rightmost (maximum key) node or nil if the tree is empty.
// Renamed from MaxNode for clarity.
// Time complexity: O(log n).
func (t *Tree[K, V]) GetRightNode() *Node[K, V] {
	return t.getRightNode(t.root)
}

// Left returns the minimum key and value in the tree.
// Returns found as true if an element is found, otherwise false.
// Time complexity: O(log n).
func (t *Tree[K, V]) Left() (key K, value V, found bool) {
	node := t.GetLeftNode()
	if node != nil {
		return node.key, node.value, true
	}
	var zeroKey K
	var zeroValue V
	return zeroKey, zeroValue, false
}

// Right returns the maximum key and value in the tree.
// Returns found as true if an element is found, otherwise false.
// Time complexity: O(log n).
func (t *Tree[K, V]) Right() (key K, value V, found bool) {
	node := t.GetRightNode()
	if node != nil {
		return node.key, node.value, true
	}
	var zeroKey K
	var zeroValue V
	return zeroKey, zeroValue, false
}

// RemoveLeft removes the minimum key-value pair from the tree.
// Returns the removed key, value, and true if an element was removed, otherwise false.
// Time complexity: O(log n).
func (t *Tree[K, V]) RemoveLeft() (key K, value V, removed bool) {
	node := t.GetLeftNode()
	if node != nil {
		k, v := node.key, node.value
		t.Remove(k)
		return k, v, true
	}
	var zeroKey K
	var zeroValue V
	return zeroKey, zeroValue, false
}

// RemoveRight removes the maximum key-value pair from the tree.
// Returns the removed key, value, and true if an element was removed, otherwise false.
// Time complexity: O(log n).
func (t *Tree[K, V]) RemoveRight() (key K, value V, removed bool) {
	node := t.GetRightNode()
	if node != nil {
		k, v := node.key, node.value
		t.Remove(k)
		return k, v, true
	}
	var zeroKey K
	var zeroValue V
	return zeroKey, zeroValue, false
}

// Floor finds the largest node less than or equal to the given key.
//
// Returns the node and true if found, nil and false otherwise. Panics if the
// key type is incompatible with the comparator. Time complexity: O(log n).
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

// Ceiling finds the smallest node greater than or equal to the given key.
//
// Returns the node and true if found, nil and false otherwise. Panics if the
// key type is incompatible with the comparator. Time complexity: O(log n).
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

// Entries returns all keys and values in in-order traversal.
//
// More efficient than calling Keys() and Values() separately as it traverses
// the tree only once. Time complexity: O(n).
func (t *Tree[K, V]) Entries() ([]K, []V) {
	keys := make([]K, t.len)
	vals := make([]V, t.len)
	it := t.Iterator()

	for i := 0; it.Next(); i++ {
		keys[i], vals[i] = it.Key(), it.Value()
	}

	return keys, vals
}

// Len returns the number of nodes in the tree.
//
// Time complexity: O(1).
func (t *Tree[K, V]) Len() int {
	return t.len
}

// Empty checks if the tree contains no nodes.
// Time complexity: O(1).
func (t *Tree[K, V]) Empty() bool {
	return t.len == 0
}

// Clear removes all nodes from the tree.
//
// Time complexity: O(1).
func (t *Tree[K, V]) Clear() {
	t.root = nil
	t.len = 0
}

// Comparator returns the comparator used by the tree.
//
// Time complexity: O(1).
func (t *Tree[K, V]) Comparator() cmp.Comparator[K] {
	return t.comparator
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
	t.output(t.root, "", true, &sb)

	return sb.String()
}

// lookup finds the node with the given key.
//
// Returns nil if not found. Time complexity: O(log n).
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

// getLeftNode finds the leftmost node in the subtree.
// Renamed from minNode for clarity as an internal helper.
// Returns nil if the subtree is empty. Time complexity: O(log n).
func (t *Tree[K, V]) getLeftNode(node *Node[K, V]) *Node[K, V] {
	for node != nil && node.left != nil {
		node = node.left
	}

	return node
}

// getRightNode finds the rightmost node in the subtree.
// Renamed from maxNode for clarity as an internal helper.
// Returns nil if the subtree is empty. Time complexity: O(log n).
func (t *Tree[K, V]) getRightNode(node *Node[K, V]) *Node[K, V] {
	for node != nil && node.right != nil {
		node = node.right
	}

	return node
}

// replaceNode substitutes the `old` node with the `new` node in the tree structure.
// It updates the parent of `old` to point to `new`, and sets `new`'s parent
// to be `old`'s parent. This function does not modify children of `old` or `new`.
//
// - old: The node to be replaced.
// - new: The node to take `old`'s place. Can be nil (e.g., when deleting a leaf).
func (t *Tree[K, V]) replaceNode(old, new *Node[K, V]) {
	// Case 1: `old` is the root of the tree.
	if old.parent == nil {
		t.root = new
	} else if old == old.parent.left {
		// Case 2: `old` is a left child.
		old.parent.left = new
	} else {
		// Case 3: `old` is a right child.
		old.parent.right = new
	}

	// Update `new` node's parent pointer, if `new` is not nil.
	if new != nil {
		new.parent = old.parent
	}
}

// rotateLeft performs a left rotation around the given node n.
// This operation is a fundamental tree restructuring maneuver used in balancing.
//
// Diagram (n is the pivot):
//
//	  P             P
//	  |             |
//	  n             r        (r becomes the new child of P)
//	 / \           / \
//	L   r   ==>   n   RR     (n becomes the left child of r)
//	   / \       / \
//	  RL  RR    L   RL       (r's original left child (RL) becomes n's new right child)
//
// - n: The node to rotate around. n.right (r) must not be nil.
// - r: The right child of n, which will become the new root of this subtree.
func (t *Tree[K, V]) rotateLeft(n *Node[K, V]) {
	r := n.right // r is n's right child, which will move up.

	// Step 1: Replace n with r in n's parent's child link.
	// r's parent becomes n's original parent.
	t.replaceNode(n, r)

	// Step 2: r's original left child (RL) becomes n's new right child.
	n.right = r.left
	if r.left != nil {
		r.left.parent = n // Update RL's parent to n.
	}

	// Step 3: n becomes the left child of r.
	r.left = n
	n.parent = r // Update n's parent to r.
}

// rotateRight performs a right rotation around the given node n.
// This operation is the mirror image of a left rotation.
//
// Diagram (n is the pivot):
//
//	    P             P
//	    |             |
//	    n             l        (l becomes the new child of P)
//	   / \           / \
//	  l   R   ==>   LL  n      (n becomes the right child of l)
//	 / \               / \
//	LL  LR            LR  R      (l's original right child (LR) becomes n's new left child)
//
// - n: The node to rotate around. n.left (l) must not be nil.
// - l: The left child of n, which will become the new root of this subtree.
func (t *Tree[K, V]) rotateRight(n *Node[K, V]) {
	l := n.left // l is n's left child, which will move up.

	// Step 1: Replace n with l in n's parent's child link.
	// l's parent becomes n's original parent.
	t.replaceNode(n, l)

	// Step 2: l's original right child (LR) becomes n's new left child.
	n.left = l.right
	if l.right != nil {
		l.right.parent = n // Update LR's parent to n.
	}

	// Step 3: n becomes the right child of l.
	l.right = n
	n.parent = l // Update n's parent to l.
}

// insertFixup restores Red-Black properties after a new node `n` (which is red) is inserted.
// This function is called when `n` and its parent `n.parent` are both red (a "double red" violation).
// It iteratively moves up the tree, resolving violations until properties are restored.
//
// Legend for understanding cases:
//
//	N: Current node being fixed (initially the newly inserted node, always red at the start of a fixup iteration within the loop).
//	P: parent of N.
//	G: Grandparent of N.
//	U: Uncle of N (sibling of P).
func (t *Tree[K, V]) insertFixup(n *Node[K, V]) {
	parent := n.parent

	// Case 0: N is the root. Color it black.
	if parent == nil {
		n.color = black // Property 2: Root is black.

		return
	}

	// Case 1: parent P is black.
	// If P is black, and N is red, no Red-Black properties are violated.
	if parent.color == black {
		return // Tree is already balanced with respect to N.
	}

	// At this point, N is red and P is red (double red violation).
	// G (grandparent) must exist because P is red, and the root cannot be red if it has children.
	// The main loop of fixup starts when P is red.
	// The `insertFixupStep` handles the core logic when P is red and U (uncle) is black.
	// If U is red, colors are flipped, and fixup continues from G.
	uncle := n.uncle()             // Maybe nil.
	grandparent := n.grandparent() // Must exist as P is red.

	// Case 2: parent P is red, Uncle U is red.
	// This is the "color flip" case.
	//
	//       G(black)           G(red) <-- Recolor G, P, U
	//      /    \             /    \
	//   P(red)  U(red) ==> P(black) U(black)
	//   /                  /
	// N(red)             N(red) (N is the new node for next iteration)
	//
	// Recolor P and U to black, G to red.
	// Then, recursively fixup G as it might now violate properties (e.g., G is root or G's parent is red).
	if color(uncle) == red {
		parent.color = black
		uncle.color = black
		grandparent.color = red

		t.insertFixup(grandparent) // Recursively fixup from G.

		return
	}

	// Case 3: parent P is red, Uncle U is BLACK (or nil).
	// This requires rotations. This logic is handled by insertFixupStep.
	// `n` is red, `parent` is red, `uncle` is black, `grandparent` is black (initially).
	// Determine if P is a left or right child to find its sibling `s`.
	if parent == grandparent.left {
		// P is G's left child. U (G.right) is black.
		// Case 3a: N is P's right child (forms a "triangle" G-P-N: G <-- P(L) --> N(R)).
		// This requires a left rotation on P to make it a "line".
		//
		//       G(B)                G(B)
		//      /    \              /    \
		//   P(R)   U(B)  ==>    N(R)   U(B)  (N becomes new P for next step)
		//     \                  /
		//    N(R)              P(R) (Original P)
		//
		if n == parent.right {
			t.rotateLeft(parent)
			// After rotation, the original `parent` is now the left child of `n` (original N).
			// `n` (the original N) is now where P was, and becomes the new `parent` for the line case.
			// The node that was `parent` (original P) is now `n.left`.
			// We update `parent` to be the new parent (original N) and `n` to be its child (original P).
			parent = n      // The original N is now the parent in the G-N-P line.
			n = parent.left // The original P is now the child n.
		}

		// Case 3b: N is P's left child (forms a "line" G-P-N: G <-- P(L) <-- N(L)).
		// This also covers the case after Case 3a's rotation.
		//
		//       G(B)                  P(B) <-- Recolor P, G
		//      /    \                /    \
		//   P(R)   U(B)   ==>    N(R)   G(R)
		//   /                               \
		// N(R)                               U(B)
		//
		// Recolor P to black, G to red. Perform a right rotation on G.
		parent.color = black
		grandparent.color = red
		t.rotateRight(grandparent)
	} else { // Symmetric case: P is G's right child. U (G.left) is black.
		// Case 3c: N is P's left child (forms a "triangle" G-P-N: G --> P(R) <-- N(L)).
		// Requires a right rotation on P.
		//
		//       G(B)                G(B)
		//      /    \              /    \
		//   U(B)   P(R)   ==>   U(B)  N(R)  (N becomes new P)
		//          /                        \
		//       N(R)                        P(R) (Original P)
		//
		if n == parent.left {
			t.rotateRight(parent)
			// Similar to Case 3a, adjust `n` and `parent` for the "line" case.
			// Original N (`n`) moves up. Original P (`parent`) becomes N's right child.
			// `n` (the original N) is now where P was, and becomes the new `parent` for the line case.
			// The node that was `parent` (original P) is now `n.right`.
			// We update `parent` to be the new parent (original N) and `n` to be its child (original P).
			parent = n       // The original N is now the parent in the G-N-P line.
			n = parent.right // The original P is now the child n.
		}

		// Case 3d: N is P's right child (forms a "line" G-P-N: G --> P(R) --> N(R)).
		//
		//       G(B)                  P(B) <-- Recolor P, G
		//      /    \                /    \
		//   U(B)   P(R)   ==>    G(R)   N(R)
		//            \            /
		//           N(R)        U(B)
		//
		// Recolor P to black, G to red. Perform a left rotation on G.
		parent.color = black
		grandparent.color = red
		t.rotateLeft(grandparent)
	}
}

// deleteFixup restores Red-Black properties after a node deletion.
//
// The parameter `x` is the node that has an "extra black" or whose removal
// might have caused a violation. If `x` is nil, it represents a nil child
// that was supposed to be black. The fixup proceeds by examining `x`'s sibling.
// This implementation follows the logic similar to CLRS, handling cases iteratively.
func (t *Tree[K, V]) deleteFixup(x *Node[K, V]) {
	// Loop as long as `x` is not the root and `x` is black (or represents a missing black node).
	// If `x` becomes red, it can absorb the "extra black", and the loop terminates.
	for x != t.root && color(x) == black {
		// Determine if `x` is a left or right child to find its sibling `s`.
		if x == x.parent.left {
			// s := x.parent.right // Original way to get sibling
			s := x.sibling() // `s` is `x`'s sibling.
			// `s` cannot be nil here if the tree properties were maintained before deletion,
			// because if `x` is black, its sibling must exist to maintain black-heights
			// unless `x`'s parent is red and `s` was removed (which is not this scenario).
			// Note: x.sibling() might return nil if x.parent is nil, but x != t.root ensures x.parent exists.

			// Case 1: `x`'s sibling `s` is red.
			// (P must be black; SL and SR must be black as children of red S)
			//
			//          P(B)                 S(B)  <-- s recolored black
			//         /    \               /    \
			//      X(B*)  S(R)  --L(P)--> P(R)  SR(B) <-- P recolored red
			//             /   \           /   \
			//         SL(B) SR(B)      X(B*) SL(B)
			//                                 (New S for X is original SL, which is black)
			//
			// Action: Recolor `s` to black, `x.parent` to red. Rotate left at `x.parent`.
			//         Update `s` to `x`'s new sibling (original SL).
			// Effect: Transforms Case 1 into Case 2, 3, or 4, where `x` has a black sibling.
			if color(s) == red {
				if s != nil { // s should not be nil if it's red
					s.color = black
				}
				x.parent.color = red
				t.rotateLeft(x.parent)
				s = x.parent.right // Update `s` to the new sibling. Crucial after rotation.
			}

			// At this point, `s` must be black (due to Case 1 transformation or initially).

			// Case 2: `x`'s sibling `s` is black, and both of `s`'s children (SL, SR) are black.
			//
			//      P(c)                P(c) <-- Becomes X' if it was black
			//     /    \              /    \
			//  X(B*)  S(B)  ----->  X(B*)  S(R) <-- s recolored red
			//         /   \                /   \
			//      SL(B) SR(B)           SL(B) SR(B)
			//
			// Action: Recolor `s` to red. Move `x` up to `x.parent`.
			// Effect: The "extra black" is passed up the tree. Loop continues from `x.parent`.
			//         If `x.parent` was red, it becomes black (absorbing extra black), loop terminates.
			if color(s.left) == black && color(s.right) == black {
				if s != nil {
					s.color = red
				} // s could be nil if tree is malformed, but typically not.
				x = x.parent // Move `x` up.
				continue     // Restart loop with new x; its sibling will be re-evaluated.
			}

			// If we reach here, s is black and at least one of s's children is red.

			// Case 3: `x`'s sibling `s` is black, `s.left` (SL) is red, `s.right` (SR) is black. (Triangle)
			//
			//      P(c)                   P(c)
			//     /    \                 /    \
			//  X(B*)  S(B)  --R(S) --> X(B*)  SL(B) <-- New S for X (orig SL, recolored black)
			//         /   \                    \
			//      SL(R) SR(B)                  S(R) <-- Orig S, recolored red
			//                                    \
			//                                    SR(B)
			//
			// Action: Recolor `s.left` to black, `s` to red. Rotate right at `s`.
			//         Update `s` to `x`'s new sibling (original SL).
			// Effect: Transforms Case 3 into Case 4. `x`'s new sibling `s` is black and has a red right child.
			if color(s.right) == black { // s.left must be red here.
				if s.left != nil {
					s.left.color = black
				}
				if s != nil {
					s.color = red
				}
				t.rotateRight(s)
				s = x.parent.right // Update `s` to the new sibling. Crucial after rotation.
			}

			// Case 4: `x`'s sibling `s` is black, and `s.right` (SR) is red. (Line)
			//
			//      P(c1)                  S(c1) <-- s takes P's original color
			//     /    \                 /     \
			//  X(B*)  S(B)  --L(P) --> P(B)    SR(B) <-- P becomes black, SR becomes black
			//         /   \            /   \
			//      SL(c2)SR(R)      X(B*) SL(c2)
			//
			// Action: Recolor `s` with `x.parent`'s color. Recolor `x.parent` to black.
			//         Recolor `s.right` to black. Rotate left at `x.parent`.
			//         Set `x` to `t.root` to terminate the loop.
			// Effect: Fixes the Red-Black properties. The "extra black" is absorbed.
			// This is reached if Case 2 was false, and Case 3 was false (meaning s.right was red).
			if s != nil {
				s.color = color(x.parent)
			}
			x.parent.color = black
			if s.right != nil {
				s.right.color = black
			}
			t.rotateLeft(x.parent)
			x = t.root // Terminate loop.

		} else { // Symmetric cases: `x` is a right child.
			// s := x.parent.left // Original way to get sibling
			s := x.sibling() // `s` is `x`'s sibling.

			// Case 1 (symmetric): `x`'s sibling `s` is red.
			// (P must be black; SL and SR must be black as children of red S)
			//
			//        P(B)                   S(B)  <-- s recolored black
			//       /    \                 /    \
			//    S(R)    X(B*) --R(P)--> SL(B)  P(R) <-- P recolored red
			//   /   \                           /   \
			// SL(B) SR(B)                    SR(B)  X(B*)
			//                                (New S for X is original SR, which is black)
			//
			// Action: Recolor `s` to black, `x.parent` to red. Rotate right at `x.parent`.
			//         Update `s` to `x`'s new sibling (original SR).
			// Effect: Transforms Case 1 into Case 2, 3, or 4 (symmetric) where `x` has a black sibling.
			if color(s) == red {
				if s != nil { // s should not be nil if it's red
					s.color = black
				}
				x.parent.color = red
				t.rotateRight(x.parent)
				s = x.parent.left // Update `s`. Crucial after rotation.
			}

			// Case 2 (symmetric): `s` is black, and both of `s`'s children (SL, SR) are black.
			//
			//      P(c)                   P(c) <-- Becomes X' if it was black
			//     /    \                 /    \
			//  S(B)    X(B*) -------->  S(R)   X(B*) <-- s recolored red
			//  /   \                    /   \
			// SL(B) SR(B)              SL(B) SR(B)
			//
			// Action: Recolor `s` to red. Move `x` up to `x.parent`.
			// Effect: The "extra black" is passed up the tree. Loop continues from `x.parent`.
			if color(s.right) == black && color(s.left) == black {
				if s != nil {
					s.color = red
				}
				x = x.parent
				continue // Restart loop with new x; its sibling will be re-evaluated.
			}

			// If we reach here, s is black and at least one of s's children is red.

			// Case 3 (symmetric): `s` is black, `s.right` (SR) is red, `s.left` (SL) is black. (Triangle)
			//
			//      P(c)                      P(c)
			//     /    \                    /    \
			//  S(B)    X(B*)  --L(S)-->  SR(B)  X(B*) <-- New S for X (orig SR, recolored black)
			//  /   \                    /
			// SL(B) SR(R)             S(R) <-- Orig S, recolored red
			//                         /
			//                       SL(B)
			//
			// Action: Recolor `s.right` to black, `s` to red. Rotate left at `s`.
			//         Update `s` to `x`'s new sibling (original SR).
			// Effect: Transforms Case 3 into Case 4 (symmetric). `x`'s new sibling `s` is black and has a red left child.
			if color(s.left) == black { // s.right must be red here.
				if s.right != nil {
					s.right.color = black
				}
				if s != nil {
					s.color = red
				}
				t.rotateLeft(s)
				s = x.parent.left // Update `s`. Crucial after rotation.
			}

			// Case 4 (symmetric): `s` is black, and `s.left` (SL) is red. (Line)
			//
			//      P(c1)                     S(c1) <-- s takes P's original color
			//     /    \                    /     \
			//   S(B)    X(B*)  --R(P)-->  SL(B)   P(B) <-- P becomes black, SL becomes black
			//  /   \                             /   \
			// SL(R) SR(c2)                     SR(c2) X(B*)
			//
			// Action: Recolor `s` with `x.parent`'s color. Recolor `x.parent` to black.
			//         Recolor `s.left` to black. Rotate right at `x.parent`.
			//         Set `x` to `t.root` to terminate the loop.
			// Effect: Fixes the Red-Black properties. The "extra black" is absorbed.
			// This is reached if Case 2 was false, and Case 3 was false (meaning s.left was red).
			if s != nil {
				s.color = color(x.parent)
			}
			x.parent.color = black
			if s.left != nil {
				s.left.color = black
			}
			t.rotateRight(x.parent)
			x = t.root // Terminate loop.
		}
	}
	// Ensure `x` (which could be the original `x` if it became red, or the root) is black.
	// This handles the case where `x` was red and absorbed the extra black, or `x` is the root.
	if x != nil {
		x.color = black
	}
}

// output builds a string representation of the tree recursively.
func (t *Tree[K, V]) output(node *Node[K, V], prefix string, isTail bool, sb *strings.Builder) {
	if node.right != nil {
		newPrefix := prefix + ternary(isTail, "│   ", "    ")
		t.output(node.right, newPrefix, false, sb)
	}

	sb.WriteString(prefix)
	sb.WriteString(ternary(isTail, "└── ", "┌── "))
	sb.WriteString(node.String() + "\n")

	if node.left != nil {
		newPrefix := prefix + ternary(isTail, "    ", "│   ")
		t.output(node.left, newPrefix, true, sb)
	}
}

// color returns the color of a node, black if nil.
func color[K comparable, V any](n *Node[K, V]) Color {
	if n == nil {
		return black
	}

	return n.color
}

// ternary is a helper for conditional expressions.
func ternary[T any](cond bool, trueVal, falseVal T) T {
	if cond {
		return trueVal
	}

	return falseVal
}
