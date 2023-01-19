package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("new error", func(t *testing.T) {
		err := New("hello")
		require.Equal(t, "hello", err.Message)
		require.Equal(t, "errors_test.go:12", err.Line)
		t.Log(err)
	})

	t.Run("new error with format", func(t *testing.T) {
		err := Newf("hello %v", 42)
		require.Equal(t, "hello 42", err.Message)
		require.Equal(t, "errors_test.go:19", err.Line)
		t.Log(err)
	})
}

func TestUnwrap(t *testing.T) {
	t.Run("enriched error is unwraped", func(t *testing.T) {
		werr := fmt.Errorf("werr")
		err := New("err").Wrap(werr)
		require.Equal(t, werr, Unwrap(err))
	})

	t.Run("enriched error is not unwraped", func(t *testing.T) {
		err := New("err")
		require.Nil(t, Unwrap(err))
	})
}

func TestIs(t *testing.T) {
	t.Run("is enriched error is true", func(t *testing.T) {
		err := New("test")
		require.True(t, Is(err, err))
	})

	t.Run("is enriched error is false", func(t *testing.T) {
		err := New("test")
		require.False(t, Is(err, New("test b")))
	})

	t.Run("is nested enriched error is true", func(t *testing.T) {
		werr := New("werr")
		err := fmt.Errorf("err: %w", werr)
		require.True(t, Is(err, werr))
	})

	t.Run("is nested enriched error is false", func(t *testing.T) {
		werr := New("werr")
		err := fmt.Errorf("err: %w", New("werr"))
		require.False(t, Is(err, werr))
	})

	t.Run("is not enriched nested error is true", func(t *testing.T) {
		werr := fmt.Errorf("werr")
		err := New("err").Wrap(werr)
		require.True(t, Is(err, werr))
	})

	t.Run("is not enriched nested error is false", func(t *testing.T) {
		werr := fmt.Errorf("werr")
		err := New("err").Wrap(fmt.Errorf("werr"))
		require.False(t, Is(err, werr))
	})
}

func TestAs(t *testing.T) {
	t.Run("has enriched error is true", func(t *testing.T) {
		var ierr Error
		err := New("err")
		require.True(t, As(err, &ierr))
	})

	t.Run("has not enriched error is false", func(t *testing.T) {
		var ierr Error
		err := fmt.Errorf("err")
		require.False(t, As(err, &ierr))
	})

	t.Run("has nested enriched error is true", func(t *testing.T) {
		var ierr Error
		err := fmt.Errorf("err: %w", New("werr"))
		require.True(t, As(err, &ierr))
	})
}

func TestHasType(t *testing.T) {
	t.Run("nil error is empty", func(t *testing.T) {
		require.True(t, HasType(nil, ""))
	})

	t.Run("enriched error is of the default type", func(t *testing.T) {
		err := New("err")
		require.True(t, HasType(err, "errors.Error"))
	})

	t.Run("enriched error is of the defined type", func(t *testing.T) {
		err := New("err").WithType("foo")
		require.True(t, HasType(err, "foo"))
	})

	t.Run("enriched error is not of the requested type", func(t *testing.T) {
		err := New("err").WithType("foo")
		require.False(t, HasType(err, "bar"))
	})

	t.Run("non enriched error is of the default type", func(t *testing.T) {
		err := fmt.Errorf("err")
		require.True(t, HasType(err, "*errors.errorString"))
	})

	t.Run("non enriched error is not of the requested type", func(t *testing.T) {
		err := fmt.Errorf("err")
		require.False(t, HasType(err, "foo"))
	})

	t.Run("enriched error is of the nested enriched type", func(t *testing.T) {
		err := New("err").Wrap(New("werr").WithType("foo"))
		require.True(t, HasType(err, "foo"))
	})

	t.Run("enriched error is of the nested non enriched type", func(t *testing.T) {
		err := New("err").Wrap(fmt.Errorf("werr"))
		require.True(t, HasType(err, "*errors.errorString"))
	})

	t.Run("non enriched error is of the nested enriched type", func(t *testing.T) {
		err := fmt.Errorf("err: %w", New("werr").WithType("foo"))
		require.True(t, HasType(err, "foo"))
	})
}

func TestTag(t *testing.T) {
	t.Run("enriched error returns the tag value", func(t *testing.T) {
		err := New("test").WithTag("foo", "bar")
		require.Equal(t, "bar", Tag(err, "foo"))
	})

	t.Run("enriched error does not returns the tag value", func(t *testing.T) {
		err := New("test")
		require.Empty(t, Tag(err, "foo"))
	})

	t.Run("nested enriched error in enriched error returns the tag value", func(t *testing.T) {
		err := New("err").Wrap(New("werr").WithTag("foo", "bar"))
		require.Equal(t, "bar", Tag(err, "foo"))
	})

	t.Run("nested enriched error in non enriched error returns the tag value", func(t *testing.T) {
		err := fmt.Errorf("err: %w", New("werr").WithTag("foo", "bar"))
		require.Equal(t, "bar", Tag(err, "foo"))
	})

	t.Run("non enriched error does not returns the tag value", func(t *testing.T) {
		err := fmt.Errorf("err")
		require.Empty(t, Tag(err, "foo"))
	})
}

func TestError(t *testing.T) {
	SetIndentEncoder()
	defer SetInlineEncoder()

	t.Run("stringify an enriched error", func(t *testing.T) {
		err := New("err").
			WithTag("foo", "bar").
			Error()
		require.Contains(t, err, "err")
		t.Log(err)
	})

	t.Run("stringify an enriched error wrapped in an enriched error", func(t *testing.T) {
		err := New("err").
			WithTag("foo", "bar").
			Wrap(New("werr").WithType("boo")).
			Error()
		require.Contains(t, err, "err")
		require.Contains(t, err, "werr")
		require.Contains(t, err, "boo")
		t.Log(err)
	})

	t.Run("stringify a non enriched error wrapped in an enriched error", func(t *testing.T) {
		err := New("err").
			WithTag("foo", "bar").
			Wrap(fmt.Errorf("werr")).
			Error()

		require.Contains(t, err, "err")
		require.Contains(t, err, "werr")
		t.Log(err)
	})

	t.Run("stringify a non enriched error wrapped in an enriched error", func(t *testing.T) {
		err := fmt.Errorf("err: %w", New("werr")).Error()
		require.Contains(t, err, "werr")
		require.NotContains(t, err, "*errors.errorString")
		require.Contains(t, err, "err")
		t.Log(err)
	})

	t.Run("stringify an enriched error with a bad tag", func(t *testing.T) {
		err := New("err").WithTag("func", func() {}).Error()
		require.Contains(t, err, "encoding error failed")
		t.Log(err)
	})
}
