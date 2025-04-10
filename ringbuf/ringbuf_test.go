package ringbuf_test

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"

	"github.com/qntx/gods/ringbuf"
)

func TestQueuePushFront(t *testing.T) {
	t.Parallel()

	queue := ringbuf.New[int](3)
	if actualValue := queue.Empty(); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}

	queue.PushFront(1)
	queue.PushFront(2)
	queue.PushFront(3)

	if actualValue := queue.Values(); actualValue[0] != 3 || actualValue[1] != 2 || actualValue[2] != 1 {
		t.Errorf("Got %v expected %v", actualValue, "[3,2,1]")
	}

	if actualValue := queue.Empty(); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}

	if actualValue := queue.Len(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}

	if actualValue, ok := queue.Front(); actualValue != 3 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}

	if actualValue, ok := queue.Back(); actualValue != 1 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}

	queue.PushFront(4)

	if actualValue := queue.Values(); actualValue[0] != 4 || actualValue[1] != 3 || actualValue[2] != 2 {
		t.Errorf("Got %v expected %v", actualValue, "[4,3,2]")
	}
}

func TestQueuePushBack(t *testing.T) {
	t.Parallel()

	queue := ringbuf.New[int](3)
	if actualValue := queue.Empty(); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}

	queue.PushBack(1)
	queue.PushBack(2)
	queue.PushBack(3)

	if actualValue := queue.Values(); actualValue[0] != 1 || actualValue[1] != 2 || actualValue[2] != 3 {
		t.Errorf("Got %v expected %v", actualValue, "[1,2,3]")
	}

	if actualValue := queue.Empty(); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}

	if actualValue := queue.Len(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}

	if actualValue, ok := queue.Front(); actualValue != 1 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}

	if actualValue, ok := queue.Back(); actualValue != 3 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
}

func TestQueueFront(t *testing.T) {
	t.Parallel()

	queue := ringbuf.New[int](3)
	if actualValue, ok := queue.Front(); actualValue != 0 || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}

	queue.PushBack(1)
	queue.PushBack(2)
	queue.PushBack(3)

	if actualValue, ok := queue.Front(); actualValue != 1 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}
}
func TestQueueBack(t *testing.T) {
	t.Parallel()

	queue := ringbuf.New[int](3)
	if actualValue, ok := queue.Back(); actualValue != 0 || ok {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}

	queue.PushBack(1)
	queue.PushBack(2)
	queue.PushBack(3)

	if actualValue, ok := queue.Back(); actualValue != 3 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
}

func TestQueuePeek(t *testing.T) {
	t.Parallel()

	queue := ringbuf.New[int](3)
	if actualValue, ok := queue.Peek(0); actualValue != 0 || ok {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}

	queue.PushBack(1)
	queue.PushBack(2)
	queue.PushBack(3)

	if actualValue, ok := queue.Peek(0); actualValue != 1 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}

	if actualValue, ok := queue.Peek(1); actualValue != 2 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}

	if actualValue, ok := queue.Peek(2); actualValue != 3 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}

	if actualValue, ok := queue.Peek(3); actualValue != 0 || ok {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}
}

func TestQueuePopFront(t *testing.T) {
	t.Parallel()

	assert := func(actualValue any, expectedValue any) {
		if actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
	}

	queue := ringbuf.New[int](3)
	assert(queue.Empty(), true)
	assert(queue.Empty(), true)
	assert(queue.Full(), false)
	assert(queue.Len(), 0)
	queue.PushBack(1)
	assert(queue.Len(), 1)
	queue.PushBack(2)
	assert(queue.Len(), 2)

	queue.PushBack(3)
	assert(queue.Len(), 3)
	assert(queue.Empty(), false)
	assert(queue.Full(), true)

	queue.PopFront()
	assert(queue.Len(), 2)

	if actualValue, ok := queue.Front(); actualValue != 2 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}

	assert(queue.Len(), 2)

	if actualValue, ok := queue.PopFront(); actualValue != 2 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}

	assert(queue.Len(), 1)

	if actualValue, ok := queue.PopFront(); actualValue != 3 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}

	assert(queue.Len(), 0)
	assert(queue.Empty(), true)
	assert(queue.Full(), false)

	if actualValue, ok := queue.PopFront(); actualValue != 0 || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}

	assert(queue.Len(), 0)

	assert(queue.Empty(), true)
	assert(queue.Full(), false)
	assert(len(queue.Values()), 0)
}

func TestQueuePopFrontFull(t *testing.T) {
	t.Parallel()

	assert := func(actualValue any, expectedValue any) {
		if actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
	}

	queue := ringbuf.New[int](2)
	assert(queue.Empty(), true)
	assert(queue.Full(), false)
	assert(queue.Len(), 0)

	queue.PushBack(1)
	assert(queue.Len(), 1)

	queue.PushBack(2)
	assert(queue.Len(), 2)
	assert(queue.Full(), true)

	if actualValue, ok := queue.Front(); actualValue != 1 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}

	queue.PushBack(3) // overwrites 1
	assert(queue.Len(), 2)

	if actualValue, ok := queue.PopFront(); actualValue != 2 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}

	if actualValue, expectedValue := queue.Len(), 1; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, ok := queue.Front(); actualValue != 3 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}

	if actualValue, expectedValue := queue.Len(), 1; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, ok := queue.PopFront(); actualValue != 3 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}

	assert(queue.Len(), 0)

	if actualValue, ok := queue.PopFront(); actualValue != 0 || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}

	assert(queue.Empty(), true)
	assert(queue.Full(), false)
	assert(len(queue.Values()), 0)
}

func TestQueuePopBack(t *testing.T) {
	t.Parallel()

	assert := func(actualValue any, expectedValue any) {
		if actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
	}

	queue := ringbuf.New[int](3)
	assert(queue.Empty(), true)
	assert(queue.Full(), false)
	assert(queue.Len(), 0)

	if actualValue, ok := queue.PopBack(); actualValue != 0 || ok {
		t.Errorf("Got %v, %v expected 0, false", actualValue, ok)
	}

	queue.PushBack(1)
	assert(queue.Len(), 1)

	if actualValue, ok := queue.PopBack(); actualValue != 1 || !ok {
		t.Errorf("Got %v, %v expected 1, true", actualValue, ok)
	}

	assert(queue.Empty(), true)
	assert(queue.Len(), 0)

	queue.PushBack(1)
	queue.PushBack(2)
	queue.PushBack(3)
	assert(queue.Len(), 3)
	assert(queue.Full(), true)

	if actualValue, ok := queue.PopBack(); actualValue != 3 || !ok {
		t.Errorf("Got %v, %v expected 3, true", actualValue, ok)
	}

	assert(queue.Len(), 2)

	if actualValue, ok := queue.Back(); actualValue != 2 || !ok {
		t.Errorf("Got %v, %v expected 2, true", actualValue, ok)
	}

	if actualValue, ok := queue.Front(); actualValue != 1 || !ok {
		t.Errorf("Got %v, %v expected 1, true", actualValue, ok)
	}

	queue.PushBack(4)
	assert(queue.Len(), 3)
	assert(queue.Full(), true)

	if actualValue, ok := queue.PopBack(); actualValue != 4 || !ok {
		t.Errorf("Got %v, %v expected 4, true", actualValue, ok)
	}

	assert(queue.Len(), 2)

	if actualValue, ok := queue.Back(); actualValue != 2 || !ok {
		t.Errorf("Got %v, %v expected 2, true", actualValue, ok)
	}

	if actualValue, ok := queue.Front(); actualValue != 1 || !ok {
		t.Errorf("Got %v, %v expected 1, true", actualValue, ok)
	}

	if actualValue, ok := queue.PopBack(); actualValue != 2 || !ok {
		t.Errorf("Got %v, %v expected 2, true", actualValue, ok)
	}

	assert(queue.Len(), 1)

	if actualValue, ok := queue.PopBack(); actualValue != 1 || !ok {
		t.Errorf("Got %v, %v expected 1, true", actualValue, ok)
	}

	assert(queue.Empty(), true)
	assert(queue.Full(), false)
	assert(len(queue.Values()), 0)
}

func TestQueuePopBackFull(t *testing.T) {
	t.Parallel()

	assert := func(actualValue any, expectedValue any) {
		if actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
	}

	queue := ringbuf.New[int](2)
	assert(queue.Empty(), true)
	assert(queue.Full(), false)
	assert(queue.Len(), 0)

	queue.PushBack(1)
	queue.PushBack(2)
	assert(queue.Len(), 2)
	assert(queue.Full(), true)

	if actualValue, ok := queue.Back(); actualValue != 2 || !ok {
		t.Errorf("Got %v, %v expected 2, true", actualValue, ok)
	}

	queue.PushBack(3) // Overwrites 1
	assert(queue.Len(), 2)

	if actualValue, ok := queue.Back(); actualValue != 3 || !ok {
		t.Errorf("Got %v, %v expected 3, true", actualValue, ok)
	}

	if actualValue, ok := queue.Front(); actualValue != 2 || !ok {
		t.Errorf("Got %v, %v expected 2, true", actualValue, ok)
	}

	if actualValue, ok := queue.PopBack(); actualValue != 3 || !ok {
		t.Errorf("Got %v, %v expected 3, true", actualValue, ok)
	}

	assert(queue.Len(), 1)

	if actualValue, ok := queue.Back(); actualValue != 2 || !ok {
		t.Errorf("Got %v, %v expected 2, true", actualValue, ok)
	}

	if actualValue, ok := queue.PopBack(); actualValue != 2 || !ok {
		t.Errorf("Got %v, %v expected 2, true", actualValue, ok)
	}

	assert(queue.Empty(), true)
	assert(queue.Full(), false)
	assert(len(queue.Values()), 0)
}

func TestQueueClear(t *testing.T) {
	t.Parallel()

	queue := ringbuf.New[int](3)
	queue.PushBack(1)
	queue.PushBack(2)
	queue.PushBack(3)

	queue.Clear()

	if actualValue := queue.Empty(); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}

	if actualValue := queue.Full(); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}

	if actualValue := queue.Len(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}

	if actualValue := queue.Values(); len(actualValue) != 0 {
		t.Errorf("Got %v expected %v", actualValue, "[]")
	}
}

func TestQueueIteratorOnEmpty(t *testing.T) {
	t.Parallel()

	queue := ringbuf.New[int](3)

	it := queue.Iterator()
	for it.Next() {
		t.Errorf("Shouldn't iterate on empty queue")
	}
}

func TestQueueIteratorNext(t *testing.T) {
	t.Parallel()

	queue := ringbuf.New[string](3)
	queue.PushBack("a")
	queue.PushBack("b")
	queue.PushBack("c")

	it := queue.Iterator()
	count := 0

	for it.Next() {
		count++
		index := it.Index()
		value := it.Value()

		switch index {
		case 0:
			if actualValue, expectedValue := value, "a"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 1:
			if actualValue, expectedValue := value, "b"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 2:
			if actualValue, expectedValue := value, "c"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			t.Errorf("Too many")
		}

		if actualValue, expectedValue := index, count-1; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
	}

	if actualValue, expectedValue := count, 3; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	queue.Clear()

	it = queue.Iterator()
	for it.Next() {
		t.Errorf("Shouldn't iterate on empty queue")
	}
}

func TestQueueIteratorPrev(t *testing.T) {
	t.Parallel()

	queue := ringbuf.New[string](3)
	queue.PushBack("a")
	queue.PushBack("b")
	queue.PushBack("c")

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
			if actualValue, expectedValue := value, "a"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 1:
			if actualValue, expectedValue := value, "b"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 2:
			if actualValue, expectedValue := value, "c"; actualValue != expectedValue {
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

func TestQueueIteratorBegin(t *testing.T) {
	t.Parallel()

	queue := ringbuf.New[string](3)
	it := queue.Iterator()
	it.Begin()
	queue.PushBack("a")
	queue.PushBack("b")
	queue.PushBack("c")

	for it.Next() {
	}

	it.Begin()
	it.Next()

	if index, value := it.Index(), it.Value(); index != 0 || value != "a" {
		t.Errorf("Got %v,%v expected %v,%v", index, value, 0, "a")
	}
}

func TestQueueIteratorEnd(t *testing.T) {
	t.Parallel()

	queue := ringbuf.New[string](3)
	it := queue.Iterator()

	if index := it.Index(); index != -1 {
		t.Errorf("Got %v expected %v", index, -1)
	}

	it.End()

	if index := it.Index(); index != 0 {
		t.Errorf("Got %v expected %v", index, 0)
	}

	queue.PushBack("a")
	queue.PushBack("b")
	queue.PushBack("c")
	it.End()

	if index := it.Index(); index != queue.Len() {
		t.Errorf("Got %v expected %v", index, queue.Len())
	}

	it.Prev()

	if index, value := it.Index(), it.Value(); index != queue.Len()-1 || value != "c" {
		t.Errorf("Got %v,%v expected %v,%v", index, value, queue.Len()-1, "c")
	}
}

func TestQueueIteratorFirst(t *testing.T) {
	t.Parallel()

	queue := ringbuf.New[string](3)

	it := queue.Iterator()
	if actualValue, expectedValue := it.First(), false; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	queue.PushBack("a")
	queue.PushBack("b")
	queue.PushBack("c")

	if actualValue, expectedValue := it.First(), true; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if index, value := it.Index(), it.Value(); index != 0 || value != "a" {
		t.Errorf("Got %v,%v expected %v,%v", index, value, 0, "a")
	}
}

func TestQueueIteratorLast(t *testing.T) {
	t.Parallel()

	queue := ringbuf.New[string](3)

	it := queue.Iterator()
	if actualValue, expectedValue := it.Last(), false; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	queue.PushBack("a")
	queue.PushBack("b")
	queue.PushBack("c")

	if actualValue, expectedValue := it.Last(), true; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if index, value := it.Index(), it.Value(); index != 2 || value != "c" {
		t.Errorf("Got %v,%v expected %v,%v", index, value, 2, "c")
	}
}

func TestQueueIteratorNextTo(t *testing.T) {
	t.Parallel()

	// Sample seek function, i.e. string starting with "b"
	seek := func(_ int, value string) bool {
		return strings.HasSuffix(value, "b")
	}

	// NextTo (empty)
	{
		queue := ringbuf.New[string](3)

		it := queue.Iterator()
		for it.NextTo(seek) {
			t.Errorf("Shouldn't iterate on empty queue")
		}
	}

	// NextTo (not found)
	{
		queue := ringbuf.New[string](3)
		queue.PushBack("xx")
		queue.PushBack("yy")

		it := queue.Iterator()
		for it.NextTo(seek) {
			t.Errorf("Shouldn't iterate on empty queue")
		}
	}

	// NextTo (found)
	{
		queue := ringbuf.New[string](3)
		queue.PushBack("aa")
		queue.PushBack("bb")
		queue.PushBack("cc")
		it := queue.Iterator()
		it.Begin()

		if !it.NextTo(seek) {
			t.Errorf("Shouldn't iterate on empty queue")
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

func TestQueueIteratorPrevTo(t *testing.T) {
	t.Parallel()

	// Sample seek function, i.e. string starting with "b"
	seek := func(_ int, value string) bool {
		return strings.HasSuffix(value, "b")
	}

	// PrevTo (empty)
	{
		queue := ringbuf.New[string](3)
		it := queue.Iterator()
		it.End()

		for it.PrevTo(seek) {
			t.Errorf("Shouldn't iterate on empty queue")
		}
	}

	// PrevTo (not found)
	{
		queue := ringbuf.New[string](3)
		queue.PushBack("xx")
		queue.PushBack("yy")
		it := queue.Iterator()
		it.End()

		for it.PrevTo(seek) {
			t.Errorf("Shouldn't iterate on empty queue")
		}
	}

	// PrevTo (found)
	{
		queue := ringbuf.New[string](3)
		queue.PushBack("aa")
		queue.PushBack("bb")
		queue.PushBack("cc")
		it := queue.Iterator()
		it.End()

		if !it.PrevTo(seek) {
			t.Errorf("Shouldn't iterate on empty queue")
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

func TestQueueIterator(t *testing.T) {
	t.Parallel()

	assert := func(actualValue any, expectedValue any) {
		if actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
	}

	queue := ringbuf.New[string](2)

	queue.PushBack("a")
	queue.PushBack("b")
	queue.PushBack("c") // overwrites "a"

	it := queue.Iterator()

	if actualIndex, expectedIndex := it.Index(), -1; actualIndex != expectedIndex {
		t.Errorf("Got %v expected %v", actualIndex, expectedIndex)
	}

	assert(it.Next(), true)

	if actualValue, actualIndex, expectedValue, expectedIndex := it.Value(), it.Index(), "b", 0; actualValue != expectedValue || actualIndex != expectedIndex {
		t.Errorf("Got %v expected %v, Got %v expected %v", actualValue, expectedValue, actualIndex, expectedIndex)
	}

	assert(it.Next(), true)

	if actualValue, actualIndex, expectedValue, expectedIndex := it.Value(), it.Index(), "c", 1; actualValue != expectedValue || actualIndex != expectedIndex {
		t.Errorf("Got %v expected %v, Got %v expected %v", actualValue, expectedValue, actualIndex, expectedIndex)
	}

	assert(it.Next(), false)

	if actualIndex, expectedIndex := it.Index(), 2; actualIndex != expectedIndex {
		t.Errorf("Got %v expected %v", actualIndex, expectedIndex)
	}

	assert(it.Next(), false)

	assert(it.Prev(), true)

	if actualValue, actualIndex, expectedValue, expectedIndex := it.Value(), it.Index(), "c", 1; actualValue != expectedValue || actualIndex != expectedIndex {
		t.Errorf("Got %v expected %v, Got %v expected %v", actualValue, expectedValue, actualIndex, expectedIndex)
	}

	assert(it.Prev(), true)

	if actualValue, actualIndex, expectedValue, expectedIndex := it.Value(), it.Index(), "b", 0; actualValue != expectedValue || actualIndex != expectedIndex {
		t.Errorf("Got %v expected %v, Got %v expected %v", actualValue, expectedValue, actualIndex, expectedIndex)
	}

	assert(it.Prev(), false)

	if actualIndex, expectedIndex := it.Index(), -1; actualIndex != expectedIndex {
		t.Errorf("Got %v expected %v", actualIndex, expectedIndex)
	}
}

func TestQueueSerialization(t *testing.T) {
	t.Parallel()

	queue := ringbuf.New[string](3)
	queue.PushBack("a")
	queue.PushBack("b")
	queue.PushBack("c")

	var err error

	assert := func() {
		if !slices.Equal(queue.Values(), []string{"a", "b", "c"}) {
			t.Errorf("Got %v expected %v", queue.Values(), []string{"a", "b", "c"})
		}

		if actualValue, expectedValue := queue.Len(), 3; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}

		if err != nil {
			t.Errorf("Got error %v", err)
		}
	}

	assert()

	bytes, err := queue.ToJSON()

	assert()

	err = queue.FromJSON(bytes)

	assert()

	_, err = json.Marshal([]any{"a", "b", "c", queue})
	if err != nil {
		t.Errorf("Got error %v", err)
	}

	err = json.Unmarshal([]byte(`["a","b","c"]`), &queue)
	if err != nil {
		t.Errorf("Got error %v", err)
	}

	assert()
}

func TestQueueString(t *testing.T) {
	t.Parallel()

	c := ringbuf.New[int](3)
	c.PushBack(1)

	if !strings.HasPrefix(c.String(), "Queue") {
		t.Errorf("String should start with container name")
	}
}

func benchmarkPushBack(b *testing.B, queue *ringbuf.Queue[int], size int) {
	b.Helper()

	for b.Loop() {
		for n := range size {
			queue.PushBack(n)
		}
	}
}

func benchmarkPopFront(b *testing.B, queue *ringbuf.Queue[int], size int) {
	b.Helper()

	for b.Loop() {
		for range size {
			queue.PopFront()
		}
	}
}

func BenchmarkArrayQueuePopFront100(b *testing.B) {
	b.StopTimer()

	size := 100
	queue := ringbuf.New[int](3)

	for n := range size {
		queue.PushBack(n)
	}

	b.StartTimer()
	benchmarkPopFront(b, queue, size)
}

func BenchmarkArrayQueuePopFront1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	queue := ringbuf.New[int](3)

	for n := range size {
		queue.PushBack(n)
	}

	b.StartTimer()
	benchmarkPopFront(b, queue, size)
}

func BenchmarkArrayQueuePopFront10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	queue := ringbuf.New[int](3)

	for n := range size {
		queue.PushBack(n)
	}

	b.StartTimer()
	benchmarkPopFront(b, queue, size)
}

func BenchmarkArrayQueuePopFront100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	queue := ringbuf.New[int](3)

	for n := range size {
		queue.PushBack(n)
	}

	b.StartTimer()
	benchmarkPopFront(b, queue, size)
}

func BenchmarkArrayQueuePushBack100(b *testing.B) {
	b.StopTimer()

	size := 100
	queue := ringbuf.New[int](3)

	b.StartTimer()
	benchmarkPushBack(b, queue, size)
}

func BenchmarkArrayQueuePushBack1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	queue := ringbuf.New[int](3)

	for n := range size {
		queue.PushBack(n)
	}

	b.StartTimer()
	benchmarkPushBack(b, queue, size)
}

func BenchmarkArrayQueuePushBack10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	queue := ringbuf.New[int](3)

	for n := range size {
		queue.PushBack(n)
	}

	b.StartTimer()
	benchmarkPushBack(b, queue, size)
}

func BenchmarkArrayQueuePushBack100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	queue := ringbuf.New[int](3)

	for n := range size {
		queue.PushBack(n)
	}

	b.StartTimer()
	benchmarkPushBack(b, queue, size)
}
