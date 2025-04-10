package main

import (
	"log"

	"github.com/qntx/gods/ringbuf"
)

func main() {
	// Initialize a circular buffer with capacity 3
	queue := ringbuf.New[int](3)
	log.Printf("Initialized circular buffer: capacity=%d, values=%v, len=%d, full=%t",
		3, queue.Values(), queue.Len(), queue.Full())

	// Test operations on an empty queue
	log.Println("Testing empty queue operations:")
	if val, ok := queue.PopBack(); ok {
		log.Printf("PopBack: %d (should not happen)", val)
	} else {
		log.Println("PopBack: queue is empty, returned (0, false)")
	}
	if val, ok := queue.PopFront(); ok {
		log.Printf("PopFront: %d (should not happen)", val)
	} else {
		log.Println("PopFront: queue is empty, returned (0, false)")
	}

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
	if val, ok := queue.Front(); ok {
		log.Printf("Front: %d (oldest element)", val)
	}
	if val, ok := queue.Back(); ok {
		log.Printf("Back: %d (newest element)", val)
	}
	if val, ok := queue.Peek(1); ok {
		log.Printf("Peek(1): %d (middle element)", val)
	}

	// Test PopBack: Remove 3 from the back
	// [1, 2, 3] -> [1, 2, _]
	log.Println("\nTesting PopBack:")
	if val, ok := queue.PopBack(); ok {
		log.Printf("PopBack: %d, values=%v, len=%d, full=%t", val, queue.Values(), queue.Len(), queue.Full())
	}

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
	if val, ok := queue.PopFront(); ok {
		log.Printf("PopFront: %d, values=%v, len=%d, full=%t", val, queue.Values(), queue.Len(), queue.Full())
	}

	// Test PushFront: Insert 6 from the front
	// [2, 5, _] -> [6, 2, 5]
	queue.PushFront(6)
	log.Printf("PushFront(6): values=%v, len=%d, full=%t", queue.Values(), queue.Len(), queue.Full())

	// Test continuous PopBack to empty the queue
	// [6, 2, 5] -> [6, 2, _] -> [6, _, _] -> [_, _, _]
	log.Println("\nTesting continuous PopBack:")
	for range 3 {
		if val, ok := queue.PopBack(); ok {
			log.Printf("PopBack: %d, values=%v, len=%d, full=%t", val, queue.Values(), queue.Len(), queue.Full())
		} else {
			log.Println("PopBack: queue is empty")
		}
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
