package errors

import (
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type nonComparableError struct {
	items []string
}

func (e nonComparableError) Error() string {
	return "bad"
}

type uiErr string

func (e uiErr) Error() string {
	return "internal error"
}

func (e uiErr) UIError() string {
	return string(e)
}

func TestNew(t *testing.T) {
	t.Run("new error", func(t *testing.T) {
		err := New("hello")
		require.Equal(t, "hello", err.Message)
		require.Equal(t, "errors_test.go:32", err.Line)
		t.Log(err)
	})

	t.Run("new error with format", func(t *testing.T) {
		err := Newf("hello %v", 42)
		require.Equal(t, "hello 42", err.Message)
		require.Equal(t, "errors_test.go:39", err.Line)
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

	t.Run("is does not panic with non comparable wrapped errors", func(t *testing.T) {
		a := New("err").Wrap(nonComparableError{items: []string{"a"}})
		b := New("err").Wrap(nonComparableError{items: []string{"a"}})

		require.NotPanics(t, func() {
			require.False(t, Is(a, b))
		})
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

	t.Run("pointer enriched error returns the tag value", func(t *testing.T) {
		err := New("test").WithTag("foo", "bar")
		require.Equal(t, "bar", Tag(&err, "foo"))
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
			WithUIError("something went wrong").
			Wrap(New("werr").WithType("boo")).
			Error()
		require.Contains(t, err, "err")
		require.Contains(t, err, "werr")
		require.Contains(t, err, "boo")
		require.Contains(t, err, "something went wrong")
		t.Log(err)
	})

	t.Run("stringify recursively wrapped enriched errors", func(t *testing.T) {
		err := New("err").
			Wrap(New("werr").
				WithType("boo").
				Wrap(New("deep").WithUIError("display message"))).
			Error()

		require.Contains(t, err, `"message": "err"`)
		require.Contains(t, err, `"message": "werr"`)
		require.Contains(t, err, `"message": "deep"`)
		require.Contains(t, err, `"ui": "display message"`)
		t.Log(err)
	})

	t.Run("stringify enriched error wrapped by pointer", func(t *testing.T) {
		werr := New("werr").WithType("boo")
		err := New("err").Wrap(&werr).Error()

		require.Contains(t, err, `"message": "werr"`)
		require.Contains(t, err, `"type": "boo"`)
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

func TestSetEncoder(t *testing.T) {
	t.Run("custom encoder is used", func(t *testing.T) {
		SetEncoder(func(any) ([]byte, error) {
			return []byte("custom"), nil
		})
		defer SetInlineEncoder()

		require.Equal(t, "custom", New("err").Error())
	})

	t.Run("nil encoder panics", func(t *testing.T) {
		require.PanicsWithValue(t, "errors: nil encoder", func() {
			SetEncoder(nil)
		})
	})
}

func TestUIError(t *testing.T) {
	t.Run("returns ui message from enriched error", func(t *testing.T) {
		err := New("internal error").WithUIError("something went wrong")
		require.Equal(t, "something went wrong", UIError(err))
	})

	t.Run("returns ui message from wrapped enriched error", func(t *testing.T) {
		err := fmt.Errorf("request failed: %w", New("internal error").WithUIError("something went wrong"))
		require.Equal(t, "something went wrong", UIError(err))
	})

	t.Run("returns ui message from any error implementing the interface", func(t *testing.T) {
		require.Equal(t, "something went wrong", UIError(uiErr("something went wrong")))
	})

	t.Run("falls back to error string when enriched error ui message is not set", func(t *testing.T) {
		err := New("internal error")
		require.Equal(t, err.Error(), UIError(err))
	})

	t.Run("falls back to wrapped error string when enriched error ui message is not set", func(t *testing.T) {
		wrapped := New("internal error")
		err := fmt.Errorf("request failed: %w", wrapped)
		require.Equal(t, wrapped.Error(), UIError(err))
	})

	t.Run("falls back to error string for non enriched errors", func(t *testing.T) {
		err := fmt.Errorf("internal error")
		require.Equal(t, err.Error(), UIError(err))
	})

	t.Run("returns empty string for nil", func(t *testing.T) {
		require.Empty(t, UIError(nil))
	})
}

func BenchmarkIs(b *testing.B) {
	b.Run("plain enriched error", func(b *testing.B) {
		err := New("err")
		target := err

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if !Is(err, target) {
				b.Fatal("expected match")
			}
		}
	})

	b.Run("tagged enriched error", func(b *testing.B) {
		err := New("err").
			WithTag("method", "GET").
			WithTag("path", "/cookies").
			WithTag("code", 401)
		target := err

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if !Is(err, target) {
				b.Fatal("expected match")
			}
		}
	})

	b.Run("tagged enriched error with reflect type", func(b *testing.B) {
		err := New("err").
			WithTag("method", "GET").
			WithTag("receiver-type", reflect.TypeOf(Error{})).
			WithTag("code", 401)
		target := err

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if !Is(err, target) {
				b.Fatal("expected match")
			}
		}
	})

	b.Run("wrapped enriched error", func(b *testing.B) {
		wrapped := New("wrapped").WithType("wrapped-code")
		err := New("err").Wrap(wrapped)
		target := wrapped

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if !Is(err, target) {
				b.Fatal("expected match")
			}
		}
	})
}

func BenchmarkErrorIs(b *testing.B) {
	b.Run("plain enriched error", func(b *testing.B) {
		err := New("err")
		target := err

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if !err.Is(target) {
				b.Fatal("expected match")
			}
		}
	})

	b.Run("tagged enriched error", func(b *testing.B) {
		err := New("err").
			WithTag("method", "GET").
			WithTag("path", "/cookies").
			WithTag("code", 401)
		target := err

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if !err.Is(target) {
				b.Fatal("expected match")
			}
		}
	})

	b.Run("tagged enriched error with reflect type", func(b *testing.B) {
		err := New("err").
			WithTag("method", "GET").
			WithTag("receiver-type", reflect.TypeOf(Error{})).
			WithTag("code", 401)
		target := err

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if !err.Is(target) {
				b.Fatal("expected match")
			}
		}
	})

	b.Run("wrapped enriched error", func(b *testing.B) {
		wrapped := io.EOF
		err := New("err").Wrap(wrapped)
		target := err

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if !err.Is(target) {
				b.Fatal("expected match")
			}
		}
	})
}
