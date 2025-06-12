package pqueue_test

import (
	"math/rand"
	"testing"

	"github.com/qntx/gods/pqueue"
)

type Element struct {
	name     string
	priority int
}

func TestPriorityQueueEnqueue(t *testing.T) {
	queue := pqueue.New[Element, int](pqueue.MaxHeap)

	if actualValue := queue.IsEmpty(); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}

	a := Element{name: "a", priority: 1}
	c := Element{name: "c", priority: 3}
	b := Element{name: "b", priority: 2}

	queue.Enqueue(a, a.priority)
	queue.Enqueue(c, c.priority)
	queue.Enqueue(b, b.priority)
	items := queue.Items()
	if len(items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(items))
	}

	v, p, ok := queue.Peek()
	if !ok || v.name != "c" || p != 3 {
		t.Errorf("Expected peek to be 'c', got %v", v)
	}
	v1, p1, ok1 := queue.Dequeue()
	if !ok1 || v1.name != "c" || p1 != 3 {
		t.Errorf("Expected first item to be 'c' with priority 3, got %v", v1)
	}

	v2, p2, ok2 := queue.Dequeue()
	if !ok2 || v2.name != "b" || p2 != 2 {
		t.Errorf("Expected second item to be 'b' with priority 2, got %v", v2)
	}

	v3, p3, ok3 := queue.Dequeue()
	if !ok3 || v3.name != "a" || p3 != 1 {
		t.Errorf("Expected third item to be 'a' with priority 1, got %v", v3)
	}

	if !queue.IsEmpty() {
		t.Error("Queue should be empty")
	}
}

func TestPriorityQueueEnqueueBulk(t *testing.T) {
	// Use MinHeap for natural number ordering
	queue := pqueue.New[int, int](pqueue.MinHeap)

	// Add elements with their values as priorities
	queue.Enqueue(15, 15)
	queue.Enqueue(20, 20)
	queue.Enqueue(3, 3)
	queue.Enqueue(1, 1)
	queue.Enqueue(2, 2)

	// Dequeue items one by one and verify order (should be ascending for MinHeap)
	v1, p1, ok1 := queue.Dequeue()
	if !ok1 || v1 != 1 || p1 != 1 {
		t.Errorf("Expected 1, got %v", v1)
	}

	v2, p2, ok2 := queue.Dequeue()
	if !ok2 || v2 != 2 || p2 != 2 {
		t.Errorf("Expected 2, got %v", v2)
	}

	v3, p3, ok3 := queue.Dequeue()
	if !ok3 || v3 != 3 || p3 != 3 {
		t.Errorf("Expected 3, got %v", v3)
	}

	v4, p4, ok4 := queue.Dequeue()
	if !ok4 || v4 != 15 || p4 != 15 {
		t.Errorf("Expected 15, got %v", v4)
	}

	v5, p5, ok5 := queue.Dequeue()
	if !ok5 || v5 != 20 || p5 != 20 {
		t.Errorf("Expected 20, got %v", v5)
	}

	// Queue should be empty after getting all items
	if !queue.IsEmpty() {
		t.Error("Queue should be empty")
	}

	// Test Clear() method
	queue.Enqueue(1, 1)
	queue.Enqueue(2, 2)
	queue.Clear()
	if !queue.IsEmpty() {
		t.Errorf("Queue should be empty after Clear()")
	}
}

func TestPriorityQueueDequeue(t *testing.T) {
	// Create a queue with MinHeap for natural ordering
	queue := pqueue.New[int, int](pqueue.MinHeap)

	// Check if queue is empty initially
	if !queue.IsEmpty() {
		t.Errorf("Queue should be empty initially")
	}

	// Add some items with values as priorities
	queue.Enqueue(3, 3)
	queue.Enqueue(2, 2)
	queue.Enqueue(1, 1)

	// Dequeue first item (should be 1 with MinHeap)
	v1, p1, ok1 := queue.Dequeue()
	if !ok1 || v1 != 1 || p1 != 1 {
		t.Errorf("Expected first item to be 1, got %v", v1)
	}

	// Dequeue second item (should be 2)
	v2, p2, ok2 := queue.Dequeue()
	if !ok2 || v2 != 2 || p2 != 2 {
		t.Errorf("Expected second item to be 2, got %v", v2)
	}

	// Dequeue third item (should be 3)
	v3, p3, ok3 := queue.Dequeue()
	if !ok3 || v3 != 3 || p3 != 3 {
		t.Errorf("Expected third item to be 3, got %v", v3)
	}

	// Queue should be empty now
	_, _, ok4 := queue.Dequeue()
	if ok4 {
		t.Errorf("Queue should be empty, but got %v", ok4)
	}

	// Check if queue is empty
	if !queue.IsEmpty() {
		t.Errorf("Queue should be empty after all items are removed")
	}

	// Check items length
	if items := queue.Items(); len(items) != 0 {
		t.Errorf("Expected 0 items, got %d", len(items))
	}
}

func TestPriorityQueueRandom(t *testing.T) {
	queue := pqueue.New[int, int](pqueue.MinHeap)

	r := rand.New(rand.NewSource(3))
	for i := 0; i < 10000; i++ {
		val := int(r.Int31n(30))
		queue.Enqueue(val, val)
	}

	prev, _, ok := queue.Dequeue()
	if !ok {
		t.Fatal("Failed to get first item from queue")
	}

	for !queue.IsEmpty() {
		curr, _, ok := queue.Dequeue()
		if !ok {
			t.Fatal("Failed to get item from queue")
		}
		if prev > curr {
			t.Errorf("Queue property invalidated. prev: %v current: %v", prev, curr)
		}
		prev = curr
	}
}

func TestPriorityQueueIsEmpty(t *testing.T) {
	queue := pqueue.New[int, int](pqueue.MinHeap)

	if !queue.IsEmpty() {
		t.Error("Queue should be empty on initialization")
	}

	queue.Enqueue(1, 1)
	if queue.IsEmpty() {
		t.Error("Queue should not be empty after adding item")
	}

	_, _, _ = queue.Dequeue()
	if !queue.IsEmpty() {
		t.Error("Queue should be empty after removing only item")
	}
}

func TestPriorityQueuePeek(t *testing.T) {
	queue := pqueue.New[int, int](pqueue.MinHeap)

	_, _, ok := queue.Peek()
	if ok {
		t.Errorf("Should not be able to peek empty queue, got %v", ok)
	}

	queue.Enqueue(3, 3)
	queue.Enqueue(2, 2)
	queue.Enqueue(1, 1)
	v, p, ok := queue.Peek()
	if !ok || v != 1 || p != 1 {
		t.Errorf("Expected peek to return 1, got %v", v)
	}

	v2, p2, ok2 := queue.Peek()
	if !ok2 || v2 != 1 || p2 != 1 {
		t.Errorf("Expected second peek to return 1, got %v", v2)
	}

	if len(queue.Items()) != 3 {
		t.Errorf("Expected 3 items in queue after peek, got %d", len(queue.Items()))
	}
}

// func TestPriorityQueueSerialization(t *testing.T) {
// 	queue := pqueue.New[string, string](pqueue.MinHeap)

// 	queue.Enqueue("c", "c")
// 	queue.Enqueue("b", "b")
// 	queue.Enqueue("a", "a")

// 	var err error
// 	assert := func() {
// 		if actualValue := queue.Items(); actualValue[0].Value != "a" || actualValue[1].Value != "b" || actualValue[2].Value != "c" {
// 			t.Errorf("Got %v expected %v", actualValue, "[1,3,2]")
// 		}
// 		if actualValue := queue.Len(); actualValue != 3 {
// 			t.Errorf("Got %v expected %v", actualValue, 3)
// 		}
// 		if actualValue, _, ok := queue.Peek(); actualValue != "a" || !ok {
// 			t.Errorf("Got %v expected %v", actualValue, "a")
// 		}
// 		if err != nil {
// 			t.Errorf("Got error %v", err)
// 		}
// 	}

// 	assert()

// 	bytes, err := queue.MarshalJSON()
// 	assert()

// 	err = queue.UnmarshalJSON(bytes)
// 	assert()

// 	bytes, err = json.Marshal([]any{"a", "b", "c", queue})
// 	if err != nil {
// 		t.Errorf("Got error %v", err)
// 	}

// 	err = json.Unmarshal([]byte(`["a","b","c"]`), &queue)
// 	if err != nil {
// 		t.Errorf("Got error %v", err)
// 	}
// 	assert()
// }
