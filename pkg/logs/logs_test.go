package logs

import (
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	SetIndentEncoder()
	defer SetInlineEncoder()

	log := New("a simple log")
	t.Log(log)
}

func TestLogWithTags(t *testing.T) {
	log := New("an log with tags").
		WithTag("string", "hello world").
		WithTag("go-stringer", goStringer{}).
		WithTag("duration", time.Duration(3600000000)).
		WithTag("int", 42).
		WithTag("int8", int8(8)).
		WithTag("int16", int16(16)).
		WithTag("int32", int32(32)).
		WithTag("int64", int64(64)).
		WithTag("uint", uint(42)).
		WithTag("uint8", uint8(8)).
		WithTag("uint16", uint16(16)).
		WithTag("uint32", uint32(32)).
		WithTag("uint64", uint64(64)).
		WithTag("float32", float32(32.42)).
		WithTag("float64", float64(64.42)).
		WithTag("slice", []string{"hello", "world"})
	t.Log("\n", log)
}

func TestLogWithBadTag(t *testing.T) {
	log := New("an log with tags").
		WithTag("func", func() {})
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
			WithTag("string", "hello world").
			WithTag("int8", int8(8)).
			WithTag("int16", int16(16)).
			WithTag("int32", int32(32)).
			WithTag("int64", int64(64))
	}
}

func BenchmarkString(b *testing.B) {
	var s string

	for n := 0; n < b.N; n++ {
		s = New("a log with tags").
			WithTag("string", "hello world").
			WithTag("int8", int8(8)).
			WithTag("int16", int16(16)).
			WithTag("int32", int32(32)).
			WithTag("int64", int64(64)).
			String()
	}

	b.Log(s)
}
