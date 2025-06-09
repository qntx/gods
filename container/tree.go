package container

// Tree is a generic interface for tree-based key-value mappings.
// It embeds Map[K, V], requiring implementations (e.g., binary search trees,
// AVL trees) to support all Map operations, typically with sorted key ordering.
// Type parameter K must be comparable for key equality checks.
// Type parameter V has no constraints, allowing any value type.
type Tree[K comparable, V any] interface {
	Map[K, V]
}
