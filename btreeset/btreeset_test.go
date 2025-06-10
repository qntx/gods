package btreeset_test

import (
	"encoding/json"
	"fmt"
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

func TestSetAdd(t *testing.T) {
	set := btreeset.New[int]()
	set.Add()
	set.Add(1)
	set.Add(2)
	set.Add(2, 3)
	set.Add()

	if actualValue := set.Empty(); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}

	if actualValue := set.Len(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
}

func TestSetContains(t *testing.T) {
	set := btreeset.New[int]()
	set.Add(3, 1, 2)

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
	set.Add(3, 1, 2)
	set.Remove()

	if actualValue := set.Len(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}

	set.Remove(1)

	if actualValue := set.Len(); actualValue != 2 {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}

	set.Remove(3)
	set.Remove(3)
	set.Remove()
	set.Remove(2)

	if actualValue := set.Len(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}
}

func TestSetEach(t *testing.T) {
	set := btreeset.New[string]()
	set.Add("c", "a", "b")
	set.Each(func(index int, value string) {
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
	})
}

func TestSetMap(t *testing.T) {
	set := btreeset.New[string]()
	set.Add("c", "a", "b")

	mappedSet := set.Map(func(index int, value string) string {
		return "mapped: " + value
	})
	if actualValue, expectedValue := mappedSet.Contains("mapped: a", "mapped: b", "mapped: c"), true; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, expectedValue := mappedSet.Contains("mapped: a", "mapped: b", "mapped: x"), false; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if mappedSet.Len() != 3 {
		t.Errorf("Got %v expected %v", mappedSet.Len(), 3)
	}
}

func TestSetSelect(t *testing.T) {
	set := btreeset.New[string]()
	set.Add("c", "a", "b")

	selectedSet := set.Select(func(index int, value string) bool {
		return value >= "a" && value <= "b"
	})
	if actualValue, expectedValue := selectedSet.Contains("a", "b"), true; actualValue != expectedValue {
		fmt.Println("A: ", selectedSet.Contains("b"))
		t.Errorf("Got %v (%v) expected %v (%v)", actualValue, selectedSet.Values(), expectedValue, "[a b]")
	}

	if actualValue, expectedValue := selectedSet.Contains("a", "b", "c"), false; actualValue != expectedValue {
		t.Errorf("Got %v (%v) expected %v (%v)", actualValue, selectedSet.Values(), expectedValue, "[a b]")
	}

	if selectedSet.Len() != 2 {
		t.Errorf("Got %v expected %v", selectedSet.Len(), 3)
	}
}

func TestSetAny(t *testing.T) {
	set := btreeset.New[string]()
	set.Add("c", "a", "b")

	any := set.Any(func(index int, value string) bool {
		return value == "c"
	})
	if any != true {
		t.Errorf("Got %v expected %v", any, true)
	}

	any = set.Any(func(index int, value string) bool {
		return value == "x"
	})
	if any != false {
		t.Errorf("Got %v expected %v", any, false)
	}
}

func TestSetAll(t *testing.T) {
	set := btreeset.New[string]()
	set.Add("c", "a", "b")

	all := set.All(func(index int, value string) bool {
		return value >= "a" && value <= "c"
	})
	if all != true {
		t.Errorf("Got %v expected %v", all, true)
	}

	all = set.All(func(index int, value string) bool {
		return value >= "a" && value <= "b"
	})
	if all != false {
		t.Errorf("Got %v expected %v", all, false)
	}
}

func TestSetFind(t *testing.T) {
	set := btreeset.New[string]()
	set.Add("c", "a", "b")

	foundIndex, foundValue := set.Find(func(index int, value string) bool {
		return value == "c"
	})
	if foundValue != "c" || foundIndex != 2 {
		t.Errorf("Got %v at %v expected %v at %v", foundValue, foundIndex, "c", 2)
	}

	foundIndex, foundValue = set.Find(func(index int, value string) bool {
		return value == "x"
	})
	if foundValue != "" || foundIndex != -1 {
		t.Errorf("Got %v at %v expected %v at %v", foundValue, foundIndex, nil, nil)
	}
}

func TestSetChaining(t *testing.T) {
	set := btreeset.New[string]()
	set.Add("c", "a", "b")
}

func TestSetSerialization(t *testing.T) {
	set := btreeset.New[string]()
	set.Add("a", "b", "c")

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
	c.Add(1)

	if !strings.HasPrefix(c.String(), "TreeSet") {
		t.Errorf("String should start with container name")
	}
}

func TestSetIntersection(t *testing.T) {
	set := btreeset.New[string]()
	another := btreeset.New[string]()

	intersection := set.Intersection(another)
	if actualValue, expectedValue := intersection.Len(), 0; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	set.Add("a", "b", "c", "d")
	another.Add("c", "d", "e", "f")

	intersection = set.Intersection(another)

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

	set.Add("a", "b", "c", "d")
	another.Add("c", "d", "e", "f")

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

	set.Add("a", "b", "c", "d")
	another.Add("c", "d", "e", "f")

	difference = set.Difference(another)

	if actualValue, expectedValue := difference.Len(), 2; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue := difference.Contains("a", "b"); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
}
