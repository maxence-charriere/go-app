package logs

import (
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	log := New("a simple log")
	t.Log(log)
}

func TestLogWithTags(t *testing.T) {
	log := New("an log with tags").
		Tag("string", "hello world").
		Tag("go-stringer", goStringer{}).
		Tag("duration", time.Duration(3600000000)).
		Tag("int", 42).
		Tag("int8", int8(8)).
		Tag("int16", int16(16)).
		Tag("int32", int32(32)).
		Tag("int64", int64(64)).
		Tag("uint", uint(42)).
		Tag("uint8", uint8(8)).
		Tag("uint16", uint16(16)).
		Tag("uint32", uint32(32)).
		Tag("uint64", uint64(64)).
		Tag("float32", float32(32.42)).
		Tag("float64", float64(64.42)).
		Tag("slice", []string{"hello", "world"})
	t.Log("\n", log)
}

func TestNewf(t *testing.T) {
	log := Newf("hello %q", "world")
	t.Log(log)
}

type goStringer struct{}

func (s goStringer) GoString() string {
	return "go stringer !"
}

func BenchmarkNew(b *testing.B) {
	for n := 0; n < b.N; n++ {
		New("a log with tags").
			Tag("string", "hello world").
			Tag("int8", int8(8)).
			Tag("int16", int16(16)).
			Tag("int32", int32(32)).
			Tag("int64", int64(64))
	}
}

func BenchmarkString(b *testing.B) {
	var s string

	for n := 0; n < b.N; n++ {
		s = New("a log with tags").
			Tag("string", "hello world").
			Tag("int8", int8(8)).
			Tag("int16", int16(16)).
			Tag("int32", int32(32)).
			Tag("int64", int64(64)).
			String()
	}

	b.Log(s)
}
