package main

import (
	"log"

	"github.com/qntx/gods/deque"
)

func main() {
	// Initialize a circular buffer with capacity 3
	queue := deque.New[int](3)
	log.Printf("Initialized circular buffer: capacity=%d, values=%v, len=%d, full=%t",
		3, queue.Values(), queue.Len(), queue.Full())

	// Test PushBack: Insert 1, 2, 3 from the back
	// [1, _, _] -> [1, 2, _] -> [1, 2, 3]
	log.Println("\nTesting PushBack:")
	queue.PushBack(1)
	log.Printf("PushBack(1): values=%v, len=%d, full=%t", queue.Values(), queue.Len(), queue.Full())
	queue.PushBack(2)
	log.Printf("PushBack(2): values=%v, len=%d, full=%t", queue.Values(), queue.Len(), queue.Full())
	queue.PushBack(3)
	log.Printf("PushBack(3): values=%v, len=%d, full=%t", queue.Values(), queue.Len(), queue.Full())

	// Test Peek operations
	log.Println("\nTesting Peek:")
	val, ok := queue.Front()
	if !ok {
		log.Fatal("Front: expected value")
	}
	log.Printf("Front: %d (oldest element)", val)
	val, ok = queue.Back()
	if !ok {
		log.Fatal("Back: expected value")
	}
	log.Printf("Back: %d (newest element)", val)
	val = queue.Get(1)
	log.Printf("Peek(1): %d (middle element)", val)

	// Test PopBack: Remove 3 from the back
	// [1, 2, 3] -> [1, 2, _]
	log.Println("\nTesting PopBack:")
	val, ok = queue.PopBack()
	if !ok {
		log.Fatal("PopBack: expected value")
	}
	log.Printf("PopBack: %d, values=%v, len=%d, full=%t", val, queue.Values(), queue.Len(), queue.Full())

	// Test PushFront: Insert 4 from the front
	// [1, 2, _] -> [4, 1, 2]
	log.Println("\nTesting PushFront:")
	queue.PushFront(4)
	log.Printf("PushFront(4): values=%v, len=%d, full=%t", queue.Values(), queue.Len(), queue.Full())

	// Test PushBack when full: Insert 5, overwrites oldest (4)
	// [4, 1, 2] -> [1, 2, 5]
	log.Println("\nTesting PushBack when full:")
	queue.PushBack(5)
	log.Printf("PushBack(5): values=%v, len=%d, full=%t", queue.Values(), queue.Len(), queue.Full())

	// Test PopFront: Remove 1 from the front
	// [1, 2, 5] -> [2, 5, _]
	log.Println("\nTesting PopFront:")
	val, ok = queue.PopFront()
	if !ok {
		log.Fatal("PopFront: expected value")
	}
	log.Printf("PopFront: %d, values=%v, len=%d, full=%t", val, queue.Values(), queue.Len(), queue.Full())

	// Test PushFront: Insert 6 from the front
	// [2, 5, _] -> [6, 2, 5]
	queue.PushFront(6)
	log.Printf("PushFront(6): values=%v, len=%d, full=%t", queue.Values(), queue.Len(), queue.Full())

	// Test continuous PopBack to empty the queue
	// [6, 2, 5] -> [6, 2, _] -> [6, _, _] -> [_, _, _]
	log.Println("\nTesting continuous PopBack:")
	for range 3 {
		val, ok := queue.PopBack()
		if !ok {
			log.Fatal("PopBack: expected value")
		}
		log.Printf("PopBack: %d, values=%v, len=%d, full=%t", val, queue.Values(), queue.Len(), queue.Full())
	}

	// Test mixed PushFront and PushBack
	log.Println("\nTesting mixed PushFront and PushBack:")
	queue.PushFront(7) // [7, _, _]
	log.Printf("PushFront(7): values=%v, len=%d", queue.Values(), queue.Len())
	queue.PushBack(8) // [7, 8, _]
	log.Printf("PushBack(8): values=%v, len=%d", queue.Values(), queue.Len())

	// Test Clear
	log.Println("\nTesting Clear:")
	queue.Clear()
	log.Printf("Clear: values=%v, len=%d, empty=%t", queue.Values(), queue.Len(), queue.Empty())
}
