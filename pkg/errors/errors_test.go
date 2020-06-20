package errors

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	err := New("a simple error")
	t.Log(err)
}

func TestErrorWithTags(t *testing.T) {
	err := New("an error with tags").
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
	t.Log("\n", err)
}

func TestErrorWithWrap(t *testing.T) {
	b := New("b").Wrap(errors.New("c"))
	a := New("a").
		Wrap(b).
		Wrap(nil)
	require.True(t, b.Is(a.wrap))

	d := New("d").
		Wrap(a).
		Wrap(b).
		Wrap(errors.New("f"))
	t.Log("\n", d)
}

func TestErrorWithTagsAndWrap(t *testing.T) {
	err := New("an error with tags").
		Tag("uint64", uint64(64)).
		Tag("float64", float64(64.42)).
		Tag("slice", []string{"hello", "world"}).
		Wrap(errors.New("another error"))
	t.Log("\n", err)
}

func TestLookup(t *testing.T) {
	err := New("error").Tag("foo", 42)

	v, found := err.Lookup("foo")
	require.True(t, found)
	require.Equal(t, "42", v)

	v, found = err.Lookup("bar")
	require.False(t, found)
	require.Empty(t, v)
}

func TestErrorUwrap(t *testing.T) {
	a := New("a").Tag("wrap", true)
	b := New("b").Wrap(a)

	err := b.Unwrap()
	require.Equal(t, a, err)

	err = Unwrap(b)
	require.Equal(t, a, err)
}

func TestAs(t *testing.T) {
	a := New("a").Tag("wrap", true)
	b := New("b").Wrap(a)
	c := New("c").Wrap(b)
	d := New("d")

	require.True(t, As(c, &a))
	require.True(t, As(c, &d))
}

func TestIs(t *testing.T) {
	a := New("a").Tag("wrap", true)
	b := New("b").Wrap(a)
	c := New("c").Wrap(b)
	d := New("d")
	e := errors.New("e")

	require.True(t, Is(c, a))
	require.False(t, Is(c, d))
	require.False(t, Is(d, e))
}

func TestTag(t *testing.T) {
	a := errors.New("a")
	v, isTagged := Tag(a, "test")
	require.Empty(t, v)
	require.False(t, isTagged)

	b := New("b").Tag("test", "true")
	v, isTagged = Tag(b, "test")
	require.Equal(t, "true", v)
	require.True(t, isTagged)
}

func TestNewf(t *testing.T) {
	err := Newf("hello %q", "world")
	t.Log(err)
}

type goStringer struct{}

func (s goStringer) GoString() string {
	return "go stringer !"
}

func BenchmarkNew(b *testing.B) {
	for n := 0; n < b.N; n++ {
		New("an error with tags").
			Tag("string", "hello world").
			Tag("int8", int8(8)).
			Tag("int16", int16(16)).
			Tag("int32", int32(32)).
			Tag("int64", int64(64))
	}
}

func BenchmarkError(b *testing.B) {
	var s string

	for n := 0; n < b.N; n++ {
		s = New("an error with tags").
			Tag("string", "hello world").
			Tag("int8", int8(8)).
			Tag("int16", int16(16)).
			Tag("int32", int32(32)).
			Tag("int64", int64(64)).
			Error()
	}

	b.Log(s)
}
