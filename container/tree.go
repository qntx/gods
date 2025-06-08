package container

// Tree interface that all trees implement
type Tree[V any] interface {
	Container[V]
	// Empty() bool
	// Size() int
	// Clear()
	// Values() []interface{}
	// String() string
}
