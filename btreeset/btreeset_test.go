package btreeset_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/qntx/gods/btreeset"
)

func TestSetNew(t *testing.T) {
	set := btreeset.New(2, 1)
	if actualValue := set.Len(); actualValue != 2 {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}

	values := set.Values()
	if actualValue := values[0]; actualValue != 1 {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}

	if actualValue := values[1]; actualValue != 2 {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
}

func TestSetAppend(t *testing.T) {
	set := btreeset.New[int]()
	set.Append()
	set.Append(1)
	set.Append(2)
	set.Append(2, 3)
	set.Append()

	if actualValue := set.IsEmpty(); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}

	if actualValue := set.Len(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
}

func TestSetContains(t *testing.T) {
	set := btreeset.New[int]()
	set.Append(3, 1, 2)

	if actualValue := set.Contains(); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}

	if actualValue := set.Contains(1); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}

	if actualValue := set.Contains(1, 2, 3); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}

	if actualValue := set.Contains(1, 2, 3, 4); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
}

func TestSetRemove(t *testing.T) {
	set := btreeset.New[int]()
	set.Append(3, 1, 2)
	set.RemoveAll()

	if actualValue := set.Len(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}

	set.RemoveAll(1)

	if actualValue := set.Len(); actualValue != 2 {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}

	set.RemoveAll(3)
	set.RemoveAll(3)
	set.RemoveAll()
	set.RemoveAll(2)

	if actualValue := set.Len(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}
}

func TestSetChaining(t *testing.T) {
	set := btreeset.New[string]()
	set.Append("c", "a", "b")
}

func TestSetSerialization(t *testing.T) {
	set := btreeset.New[string]()
	set.Append("a", "b", "c")

	var err error

	assert := func() {
		if actualValue, expectedValue := set.Len(), 3; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}

		if actualValue := set.Contains("a", "b", "c"); actualValue != true {
			t.Errorf("Got %v expected %v", actualValue, true)
		}

		if err != nil {
			t.Errorf("Got error %v", err)
		}
	}

	assert()

	bytes, err := set.MarshalJSON()

	assert()

	err = set.UnmarshalJSON(bytes)

	assert()

	_, err = json.Marshal([]any{"a", "b", "c", set})
	if err != nil {
		t.Errorf("Got error %v", err)
	}

	err = json.Unmarshal([]byte(`["1","2","3"]`), &set)
	if err != nil {
		t.Errorf("Got error %v", err)
	}
}

func TestSetString(t *testing.T) {
	c := btreeset.New[int]()
	c.Append(1)

	if !strings.HasPrefix(c.String(), "TreeSet") {
		t.Errorf("String should start with container name")
	}
}

func TestSetIntersection(t *testing.T) {
	set := btreeset.New[string]()
	another := btreeset.New[string]()

	intersection := set.Intersect(another)
	if actualValue, expectedValue := intersection.Len(), 0; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	set.Append("a", "b", "c", "d")
	another.Append("c", "d", "e", "f")

	intersection = set.Intersect(another)

	if actualValue, expectedValue := intersection.Len(), 2; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue := intersection.Contains("c", "d"); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
}

func TestSetUnion(t *testing.T) {
	set := btreeset.New[string]()
	another := btreeset.New[string]()

	union := set.Union(another)
	if actualValue, expectedValue := union.Len(), 0; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	set.Append("a", "b", "c", "d")
	another.Append("c", "d", "e", "f")

	union = set.Union(another)

	if actualValue, expectedValue := union.Len(), 6; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue := union.Contains("a", "b", "c", "d", "e", "f"); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
}

func TestSetDifference(t *testing.T) {
	set := btreeset.New[string]()
	another := btreeset.New[string]()

	difference := set.Difference(another)
	if actualValue, expectedValue := difference.Len(), 0; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	set.Append("a", "b", "c", "d")
	another.Append("c", "d", "e", "f")

	difference = set.Difference(another)

	if actualValue, expectedValue := difference.Len(), 2; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue := difference.Contains("a", "b"); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
}
