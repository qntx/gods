package util_test

import (
	"strings"
	"testing"

	"github.com/qntx/gods/util"
)

func TestToStringInts(t *testing.T) {
	t.Parallel()

	var value any

	value = int8(1)
	if actualValue, expectedValue := util.ToString(value), "1"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	value = int16(1)
	if actualValue, expectedValue := util.ToString(value), "1"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	value = int32(1)
	if actualValue, expectedValue := util.ToString(value), "1"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	value = int64(1)
	if actualValue, expectedValue := util.ToString(value), "1"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	value = rune(1)
	if actualValue, expectedValue := util.ToString(value), "1"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestToStringUInts(t *testing.T) {
	t.Parallel()

	var value any

	value = uint8(1)
	if actualValue, expectedValue := util.ToString(value), "1"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	value = uint16(1)
	if actualValue, expectedValue := util.ToString(value), "1"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	value = uint32(1)
	if actualValue, expectedValue := util.ToString(value), "1"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	value = uint64(1)
	if actualValue, expectedValue := util.ToString(value), "1"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	value = byte(1)
	if actualValue, expectedValue := util.ToString(value), "1"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestToStringFloats(t *testing.T) {
	t.Parallel()

	var value any

	value = float32(1.123456)
	if actualValue, expectedValue := util.ToString(value), "1.123456"; !strings.HasPrefix(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	value = float64(1.123456)
	if actualValue, expectedValue := util.ToString(value), "1.123456"; !strings.HasPrefix(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestToStringOther(t *testing.T) {
	t.Parallel()

	var value any

	value = "abc"
	if actualValue, expectedValue := util.ToString(value), "abc"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	value = true
	if actualValue, expectedValue := util.ToString(value), "true"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	type T struct {
		id   int
		name string
	}

	if actualValue, expectedValue := util.ToString(T{1, "abc"}), "{id:1 name:abc}"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}
