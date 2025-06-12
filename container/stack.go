package container

// Stack is a generic interface for a last-in, first-out (LIFO) data structure.
// It supports adding elements to the top (push) and removing them from the top (pop).
// Implementations (e.g., array-based or linked-list stacks) must provide all
// operations defined here, including those inherited from Container[T] (e.g., Len, IsEmpty, Clear).
// Type parameter T must be comparable to enable equality checks for elements, if needed by specific implementations,
// though stack operations themselves (Push, Pop, Peek) do not inherently require comparability.
type Stack[T comparable] interface {
	Container[T]

	// Push adds an element to the top of the stack.
	Push(value T)

	// Pop removes and returns the top element of the stack.
	// Returns the element and true if the stack is non-empty,
	// or the zero value of T and false if the stack is empty.
	Pop() (value T, ok bool)

	// Peek returns the top element of the stack without removing it.
	// Returns the element and true if the stack is non-empty,
	// or the zero value of T and false if the stack is empty.
	Peek() (value T, ok bool)
}
