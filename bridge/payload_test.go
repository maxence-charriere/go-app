package bridge

import (
	"reflect"
	"testing"
)

func TestPayloadFromString(t *testing.T) {
	PayloadFromString("42")
}

func TestNewPayloadPanic(t *testing.T) {
	defer func() { recover() }()
	NewPayload(make(chan bool))
	t.Fatal("should have panic")
}

func TestPayloadLen(t *testing.T) {
	p := NewPayload(42)

	if l := p.Len(); l != 2 {
		t.Fatal("payload len should be 2:", l)
	}
}

func TestPayloadBytes(t *testing.T) {
	p := NewPayload(42)

	expect := []byte("42")
	if b := p.Bytes(); !reflect.DeepEqual(b, expect) {
		t.Fatalf("payload bytes should be %v: %v", expect, b)
	}
}

func TestPayloadString(t *testing.T) {
	p := NewPayload(42)

	if str := p.String(); str != "42" {
		t.Fatalf(`payload string should be "42": "%v"`, str)
	}
}

func TestPayloadUnmarshal(t *testing.T) {
	p := NewPayload("hello world")

	var res string
	p.Unmarshal(&res)

	if res != "hello world" {
		t.Fatalf(`unmarshaled result should be "hello wordl": "%v"`, res)
	}
}

func TestPayloadUnmarshalNotPtr(t *testing.T) {
	defer func() { recover() }()

	p := NewPayload("hello world")

	var res string
	p.Unmarshal(res)
	t.Fatal("should have panic")

}

func TestPayloadUnmarshalBadPayload(t *testing.T) {
	defer func() { recover() }()

	p := PayloadFromBytes([]byte("}dsfa{"))

	var res string
	p.Unmarshal(&res)
	t.Fatal("should have panic")
}
