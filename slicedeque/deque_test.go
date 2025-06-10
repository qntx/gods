package slicedeque_test

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"

	"github.com/qntx/gods/slicedeque"
)

func TestQueuePushFront(t *testing.T) {
	t.Parallel()

	queue := slicedeque.New[int](3)
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

	queue := slicedeque.New[int](3)
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

	queue := slicedeque.New[int](3)
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

	queue := slicedeque.New[int](3)
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

func TestQueueGet(t *testing.T) {
	t.Parallel()

	queue := slicedeque.New[int](3)

	queue.PushBack(1)
	queue.PushBack(2)
	queue.PushBack(3)

	if actualValue := queue.Get(0); actualValue != 1 {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}

	if actualValue := queue.Get(1); actualValue != 2 {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}

	if actualValue := queue.Get(2); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
}

func TestQueuePopFront(t *testing.T) {
	t.Parallel()

	assert := func(actualValue any, expectedValue any) {
		if actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
	}

	queue := slicedeque.New[int](3)
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

	queue := slicedeque.New[int](2)
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

	queue := slicedeque.New[int](3)
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

	queue := slicedeque.New[int](2)
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

	queue := slicedeque.New[int](3)
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

func TestQueueSerialization(t *testing.T) {
	t.Parallel()

	queue := slicedeque.New[string](3)
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

	bytes, err := queue.MarshalJSON()

	assert()

	err = queue.UnmarshalJSON(bytes)

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

	c := slicedeque.New[int](3)
	c.PushBack(1)

	if !strings.HasPrefix(c.String(), "Deque") {
		t.Errorf("String should start with container name")
	}
}
