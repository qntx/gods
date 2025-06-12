package slicestack_test

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/qntx/gods/slicestack"
)

func TestStackPush(t *testing.T) {
	stack := slicestack.New[int]()
	if actualValue := stack.IsEmpty(); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}

	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	if actualValue := stack.Values(); actualValue[0] != 3 || actualValue[1] != 2 || actualValue[2] != 1 {
		t.Errorf("Got %v expected %v", actualValue, "[3,2,1]")
	}

	if actualValue := stack.IsEmpty(); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}

	if actualValue := stack.Len(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}

	if actualValue, ok := stack.Peek(); actualValue != 3 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
}

func TestStackPeek(t *testing.T) {
	stack := slicestack.New[int]()
	if actualValue, ok := stack.Peek(); actualValue != 0 || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}

	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	if actualValue, ok := stack.Peek(); actualValue != 3 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
}

func TestStackPop(t *testing.T) {
	stack := slicestack.New[int]()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	stack.Pop()

	if actualValue, ok := stack.Peek(); actualValue != 2 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}

	if actualValue, ok := stack.Pop(); actualValue != 2 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}

	if actualValue, ok := stack.Pop(); actualValue != 1 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}

	if actualValue, ok := stack.Pop(); actualValue != 0 || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}

	if actualValue := stack.IsEmpty(); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}

	if actualValue := stack.Values(); len(actualValue) != 0 {
		t.Errorf("Got %v expected %v", actualValue, "[]")
	}
}

func TestStackSerialization(t *testing.T) {
	stack := slicestack.New[string]()
	stack.Push("a")
	stack.Push("b")
	stack.Push("c")

	var err error

	assert := func() {
		slices.Equal(stack.Values(), []string{"c", "b", "a"})

		if actualValue, expectedValue := stack.Len(), 3; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}

		if err != nil {
			t.Errorf("Got error %v", err)
		}
	}

	assert()

	bytes, err := stack.MarshalJSON()

	assert()

	err = stack.UnmarshalJSON(bytes)

	assert()

	bytes, err = json.Marshal([]any{"a", "b", "c", stack})
	if err != nil {
		t.Errorf("Got error %v", err)
	}

	err = json.Unmarshal([]byte(`["a","b","c"]`), &stack)
	if err != nil {
		t.Errorf("Got error %v", err)
	}

	assert()
}
