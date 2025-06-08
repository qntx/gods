package pqueue

import (
	"encoding/json"
	"math/rand"
	"strings"
	"testing"
)

type Element struct {
	priority int
	name     string
}

func TestBinaryQueueEnqueue(t *testing.T) {
	queue := New[Element, int](MaxHeap)

	if actualValue := queue.Empty(); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}

	a := Element{name: "a", priority: 1}
	c := Element{name: "c", priority: 3}
	b := Element{name: "b", priority: 2}

	queue.Put(a, a.priority)
	queue.Put(c, c.priority)
	queue.Put(b, b.priority)
	items := queue.Items()
	if len(items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(items))
	}

	peek, ok := queue.Peek()
	if !ok || peek.Value.name != "c" {
		t.Errorf("Expected peek to be 'c', got %v", peek)
	}
	item1, ok1 := queue.Get()
	if !ok1 || item1.Value.name != "c" || item1.Priority != 3 {
		t.Errorf("Expected first item to be 'c' with priority 3, got %v", item1)
	}

	item2, ok2 := queue.Get()
	if !ok2 || item2.Value.name != "b" || item2.Priority != 2 {
		t.Errorf("Expected second item to be 'b' with priority 2, got %v", item2)
	}

	item3, ok3 := queue.Get()
	if !ok3 || item3.Value.name != "a" || item3.Priority != 1 {
		t.Errorf("Expected third item to be 'a' with priority 1, got %v", item3)
	}

	if !queue.Empty() {
		t.Error("Queue should be empty")
	}
}

func TestBinaryQueueEnqueueBulk(t *testing.T) {
	// Use MinHeap for natural number ordering
	queue := New[int, int](MinHeap)

	// Add elements with their values as priorities
	queue.Put(15, 15)
	queue.Put(20, 20)
	queue.Put(3, 3)
	queue.Put(1, 1)
	queue.Put(2, 2)

	// Get items one by one and verify order (should be ascending for MinHeap)
	item1, ok1 := queue.Get()
	if !ok1 || item1.Value != 1 || item1.Priority != 1 {
		t.Errorf("Expected 1, got %v", item1)
	}

	item2, ok2 := queue.Get()
	if !ok2 || item2.Value != 2 || item2.Priority != 2 {
		t.Errorf("Expected 2, got %v", item2)
	}

	item3, ok3 := queue.Get()
	if !ok3 || item3.Value != 3 || item3.Priority != 3 {
		t.Errorf("Expected 3, got %v", item3)
	}

	item4, ok4 := queue.Get()
	if !ok4 || item4.Value != 15 || item4.Priority != 15 {
		t.Errorf("Expected 15, got %v", item4)
	}

	item5, ok5 := queue.Get()
	if !ok5 || item5.Value != 20 || item5.Priority != 20 {
		t.Errorf("Expected 20, got %v", item5)
	}

	// Queue should be empty after getting all items
	if !queue.Empty() {
		t.Error("Queue should be empty")
	}

	// Test Clear() method
	queue.Put(1, 1)
	queue.Put(2, 2)
	queue.Clear()
	if !queue.Empty() {
		t.Errorf("Queue should be empty after Clear()")
	}
}

func TestBinaryQueueDequeue(t *testing.T) {
	// Create a queue with MinHeap for natural ordering
	queue := New[int, int](MinHeap)

	// Check if queue is empty initially
	if !queue.Empty() {
		t.Errorf("Queue should be empty initially")
	}

	// Add some items with values as priorities
	queue.Put(3, 3)
	queue.Put(2, 2)
	queue.Put(1, 1)

	// Get first item (should be 1 with MinHeap)
	item1, ok1 := queue.Get()
	if !ok1 || item1.Value != 1 || item1.Priority != 1 {
		t.Errorf("Expected first item to be 1, got %v", item1)
	}

	// Get second item (should be 2)
	item2, ok2 := queue.Get()
	if !ok2 || item2.Value != 2 || item2.Priority != 2 {
		t.Errorf("Expected second item to be 2, got %v", item2)
	}

	// Get third item (should be 3)
	item3, ok3 := queue.Get()
	if !ok3 || item3.Value != 3 || item3.Priority != 3 {
		t.Errorf("Expected third item to be 3, got %v", item3)
	}

	// Queue should be empty now
	item4, ok4 := queue.Get()
	if ok4 {
		t.Errorf("Queue should be empty, but got %v", item4)
	}

	// Check if queue is empty
	if !queue.Empty() {
		t.Errorf("Queue should be empty after all items are removed")
	}

	// Check items length
	if items := queue.Items(); len(items) != 0 {
		t.Errorf("Expected 0 items, got %d", len(items))
	}
}

func TestBinaryQueueRandom(t *testing.T) {
	queue := New[int, int](MinHeap)

	r := rand.New(rand.NewSource(3))
	for i := 0; i < 10000; i++ {
		val := int(r.Int31n(30))
		queue.Put(val, val)
	}

	prev, ok := queue.Get()
	if !ok {
		t.Fatal("Failed to get first item from queue")
	}

	for !queue.Empty() {
		curr, _ := queue.Get()
		if prev.Priority > curr.Priority {
			t.Errorf("Queue property invalidated. prev: %v current: %v", prev, curr)
		}
		prev = curr
	}
}

func TestPriorityQueueEmpty(t *testing.T) {
	queue := New[int, int](MinHeap)

	if !queue.Empty() {
		t.Error("Queue should be empty on initialization")
	}

	queue.Put(1, 1)
	if queue.Empty() {
		t.Error("Queue should not be empty after adding item")
	}

	_, _ = queue.Get()
	if !queue.Empty() {
		t.Error("Queue should be empty after removing only item")
	}
}

func TestPriorityQueuePeek(t *testing.T) {
	queue := New[int, int](MinHeap)

	peek, ok := queue.Peek()
	if ok {
		t.Errorf("Should not be able to peek empty queue, got %v", peek)
	}

	queue.Put(3, 3)
	queue.Put(2, 2)
	queue.Put(1, 1)
	peek, ok = queue.Peek()
	if !ok || peek.Value != 1 || peek.Priority != 1 {
		t.Errorf("Expected peek to return 1, got %v", peek)
	}

	peek2, ok2 := queue.Peek()
	if !ok2 || peek2.Value != 1 || peek2.Priority != 1 {
		t.Errorf("Expected second peek to return 1, got %v", peek2)
	}

	if len(queue.Items()) != 3 {
		t.Errorf("Expected 3 items in queue after peek, got %d", len(queue.Items()))
	}
}

func TestBinaryQueueIteratorPrev(t *testing.T) {
	queue := New[int, int](MinHeap)
	queue.Put(3, 3)
	queue.Put(2, 2)
	queue.Put(1, 1)

	it := queue.Iterator()
	for it.Next() {
	}
	count := 0
	for it.Prev() {
		count++
		index := it.Index()
		value := it.Value()
		switch index {
		case 0:
			if actualValue, expectedValue := value, 1; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 1:
			if actualValue, expectedValue := value, 2; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 2:
			if actualValue, expectedValue := value, 3; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			t.Errorf("Too many")
		}
		if actualValue, expectedValue := index, 3-count; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
	}
	if actualValue, expectedValue := count, 3; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBinaryQueueIteratorBegin(t *testing.T) {
	queue := New[int, int](MinHeap)
	it := queue.Iterator()
	it.Begin()
	queue.Put(2, 2)
	queue.Put(3, 3)
	queue.Put(1, 1)
	for it.Next() {
	}
	it.Begin()
	it.Next()
	if index, value := it.Index(), it.Value(); index != 0 || value != 1 {
		t.Errorf("Got %v,%v expected %v,%v", index, value, 0, 1)
	}
}

func TestBinaryQueueIteratorEnd(t *testing.T) {
	queue := New[int, int](MinHeap)
	it := queue.Iterator()

	if index := it.Index(); index != -1 {
		t.Errorf("Got %v expected %v", index, -1)
	}

	it.End()
	if index := it.Index(); index != 0 {
		t.Errorf("Got %v expected %v", index, 0)
	}

	queue.Put(3, 3)
	queue.Put(2, 2)
	queue.Put(1, 1)
	it.End()
	if index := it.Index(); index != queue.Len() {
		t.Errorf("Got %v expected %v", index, queue.Len())
	}

	it.Prev()
	if index, value := it.Index(), it.Value(); index != queue.Len()-1 || value != 3 {
		t.Errorf("Got %v,%v expected %v,%v", index, value, queue.Len()-1, 3)
	}
}

func TestBinaryQueueIteratorFirst(t *testing.T) {
	queue := New[int, int](MinHeap)
	it := queue.Iterator()
	if actualValue, expectedValue := it.First(), false; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	queue.Put(3, 3) // [3]
	queue.Put(2, 2) // [2,3]
	queue.Put(1, 1) // [1,3,2](2 swapped with 1, hence last)
	if actualValue, expectedValue := it.First(), true; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if index, value := it.Index(), it.Value(); index != 0 || value != 1 {
		t.Errorf("Got %v,%v expected %v,%v", index, value, 0, 1)
	}
}

func TestBinaryQueueIteratorLast(t *testing.T) {
	queue := New[int, int](MinHeap)
	it := queue.Iterator()
	if actualValue, expectedValue := it.Last(), false; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	queue.Put(2, 2)
	queue.Put(3, 3)
	queue.Put(1, 1)
	if actualValue, expectedValue := it.Last(), true; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if index, value := it.Index(), it.Value(); index != 2 || value != 3 {
		t.Errorf("Got %v,%v expected %v,%v", index, value, 2, 3)
	}
}

func TestBinaryQueueIteratorNextTo(t *testing.T) {
	// Sample seek function, i.e. string starting with "b"
	seek := func(index int, value string) bool {
		return strings.HasSuffix(value, "b")
	}

	// NextTo (empty)
	{
		queue := New[string, string](MinHeap)
		it := queue.Iterator()
		for it.NextTo(seek) {
			t.Errorf("Shouldn't iterate on empty list")
		}
	}

	// NextTo (not found)
	{
		queue := New[string, string](MinHeap)
		queue.Put("xx", "xx")
		queue.Put("yy", "yy")
		it := queue.Iterator()
		for it.NextTo(seek) {
			t.Errorf("Shouldn't iterate on empty list")
		}
	}

	// NextTo (found)
	{
		queue := New[string, string](MinHeap)
		queue.Put("aa", "aa")
		queue.Put("bb", "bb")
		queue.Put("cc", "cc")
		it := queue.Iterator()
		it.Begin()
		if !it.NextTo(seek) {
			t.Errorf("Shouldn't iterate on empty list")
		}
		if index, value := it.Index(), it.Value(); index != 1 || value != "bb" {
			t.Errorf("Got %v,%v expected %v,%v", index, value, 1, "bb")
		}
		if !it.Next() {
			t.Errorf("Should go to first element")
		}
		if index, value := it.Index(), it.Value(); index != 2 || value != "cc" {
			t.Errorf("Got %v,%v expected %v,%v", index, value, 2, "cc")
		}
		if it.Next() {
			t.Errorf("Should not go past last element")
		}
	}
}

func TestBinaryQueueIteratorPrevTo(t *testing.T) {
	// Sample seek function, i.e. string starting with "b"
	seek := func(index int, value string) bool {
		return strings.HasSuffix(value, "b")
	}

	// PrevTo (empty)
	{
		queue := New[string, string](MinHeap)
		it := queue.Iterator()
		it.End()
		for it.PrevTo(seek) {
			t.Errorf("Shouldn't iterate on empty list")
		}
	}

	// PrevTo (not found)
	{
		queue := New[string, string](MinHeap)
		queue.Put("xx", "xx")
		queue.Put("yy", "yy")
		it := queue.Iterator()
		it.End()
		for it.PrevTo(seek) {
			t.Errorf("Shouldn't iterate on empty list")
		}
	}

	// PrevTo (found)
	{
		queue := New[string, string](MinHeap)
		queue.Put("aa", "aa")
		queue.Put("bb", "bb")
		queue.Put("cc", "cc")
		it := queue.Iterator()
		it.End()
		if !it.PrevTo(seek) {
			t.Errorf("Shouldn't iterate on empty list")
		}
		if index, value := it.Index(), it.Value(); index != 1 || value != "bb" {
			t.Errorf("Got %v,%v expected %v,%v", index, value, 1, "bb")
		}
		if !it.Prev() {
			t.Errorf("Should go to first element")
		}
		if index, value := it.Index(), it.Value(); index != 0 || value != "aa" {
			t.Errorf("Got %v,%v expected %v,%v", index, value, 0, "aa")
		}
		if it.Prev() {
			t.Errorf("Should not go before first element")
		}
	}
}

func TestBinaryQueueSerialization(t *testing.T) {
	queue := New[string, string](MinHeap)

	queue.Put("c", "c")
	queue.Put("b", "b")
	queue.Put("a", "a")

	var err error
	assert := func() {
		if actualValue := queue.Items(); actualValue[0].Value != "a" || actualValue[1].Value != "b" || actualValue[2].Value != "c" {
			t.Errorf("Got %v expected %v", actualValue, "[1,3,2]")
		}
		if actualValue := queue.Len(); actualValue != 3 {
			t.Errorf("Got %v expected %v", actualValue, 3)
		}
		if actualValue, ok := queue.Peek(); actualValue.Value != "a" || !ok {
			t.Errorf("Got %v expected %v", actualValue, "a")
		}
		if err != nil {
			t.Errorf("Got error %v", err)
		}
	}

	assert()

	bytes, err := queue.MarshalJSON()
	assert()

	err = queue.UnmarshalJSON(bytes)
	assert()

	bytes, err = json.Marshal([]interface{}{"a", "b", "c", queue})
	if err != nil {
		t.Errorf("Got error %v", err)
	}

	err = json.Unmarshal([]byte(`["a","b","c"]`), &queue)
	if err != nil {
		t.Errorf("Got error %v", err)
	}
	assert()
}

func benchmarkEnqueue(b *testing.B, queue *PriorityQueue[Element, int], size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			queue.Put(Element{}, n)
		}
	}
}

func benchmarkDequeue(b *testing.B, queue *PriorityQueue[Element, int], size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			queue.Get()
		}
	}
}

func BenchmarkBinaryQueueDequeue100(b *testing.B) {
	b.StopTimer()
	size := 100
	queue := New[Element, int](MinHeap)
	for n := 0; n < size; n++ {
		queue.Put(Element{}, n)
	}
	b.StartTimer()
	benchmarkDequeue(b, queue, size)
}

func BenchmarkBinaryQueueDequeue1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	queue := New[Element, int](MinHeap)
	for n := 0; n < size; n++ {
		queue.Put(Element{}, n)
	}
	b.StartTimer()
	benchmarkDequeue(b, queue, size)
}

func BenchmarkBinaryQueueDequeue10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	queue := New[Element, int](MinHeap)
	for n := 0; n < size; n++ {
		queue.Put(Element{}, n)
	}
	b.StartTimer()
	benchmarkDequeue(b, queue, size)
}

func BenchmarkBinaryQueueDequeue100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	queue := New[Element, int](MinHeap)
	for n := 0; n < size; n++ {
		queue.Put(Element{}, n)
	}
	b.StartTimer()
	benchmarkDequeue(b, queue, size)
}

func BenchmarkBinaryQueueEnqueue100(b *testing.B) {
	b.StopTimer()
	size := 100
	queue := New[Element, int](MinHeap)
	for n := 0; n < size; n++ {
		queue.Put(Element{}, n)
	}
	b.StartTimer()
	benchmarkEnqueue(b, queue, size)
}

func BenchmarkBinaryQueueEnqueue1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	queue := New[Element, int](MinHeap)
	for n := 0; n < size; n++ {
		queue.Put(Element{}, n)
	}
	b.StartTimer()
	benchmarkEnqueue(b, queue, size)
}

func BenchmarkBinaryQueueEnqueue10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	queue := New[Element, int](MinHeap)
	for n := 0; n < size; n++ {
		queue.Put(Element{}, n)
	}
	b.StartTimer()
	benchmarkEnqueue(b, queue, size)
}

func BenchmarkBinaryQueueEnqueue100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	queue := New[Element, int](MinHeap)
	for n := 0; n < size; n++ {
		queue.Put(Element{}, n)
	}
	b.StartTimer()
	benchmarkEnqueue(b, queue, size)
}
