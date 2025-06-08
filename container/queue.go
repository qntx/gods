package container

// Queue interface that all queues implement
type Queue[T comparable] interface {
	Enqueue(value T)
	Dequeue() (value T, ok bool)
	Peek() (value T, ok bool)

	Container[T]
	// Empty() bool
	// Size() int
	// Clear()
	// Values() []interface{}
	// String() string
}
