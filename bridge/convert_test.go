package bridge

import (
	"reflect"
	"testing"
)

func TestStrings(t *testing.T) {
	v := []interface{}{
		"hello",
		"world",
	}

	expected := []string{
		"hello",
		"world",
	}

	if c := Strings(v); !reflect.DeepEqual(c, expected) {
		t.Error("expected:", expected)
		t.Error("value:   ", c)
	}
}

func TestStringsPanic(t *testing.T) {
	defer func() { recover() }()

	v := []interface{}{
		4,
		2,
	}

	Strings(v)
	t.Error("no panic")
}
