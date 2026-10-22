package errors

import (
	stderr "errors"
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
		require.Equal(t, "errors_test.go:33", err.Line)
		t.Log(err)
	})

	t.Run("new error with format", func(t *testing.T) {
		err := Newf("hello %v", 42)
		require.Equal(t, "hello 42", err.Message)
		require.Equal(t, "errors_test.go:40", err.Line)
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

// Unlike errors.Join, this wrapper preserves nil children for traversal tests.
type multiError []error

func (e multiError) Error() string   { return "multiple errors" }
func (e multiError) Unwrap() []error { return e }

type unvisitedError struct{}

func (unvisitedError) Error() string { return "must not unwrap" }
func (unvisitedError) Unwrap() error { panic("visited a later branch after a match") }

func TestMetadataErrorTrees(t *testing.T) {
	leaf := New("database connection refused").
		WithType("database").
		WithTag("operation", "save").
		WithUIError("Unable to save your changes.")
	plain := stderr.New("cleanup failed")
	for _, tt := range []struct {
		name string
		err  error
	}{
		{"join", stderr.Join(plain, leaf)},
		{"multiple percent w", fmt.Errorf("cleanup: %w; save: %w", plain, leaf)},
		{"nested joins", stderr.Join(plain, stderr.Join(plain, leaf))},
		{"single wrapper around join", fmt.Errorf("request: %w", stderr.Join(plain, leaf))},
		{"join around single wrapper", stderr.Join(plain, fmt.Errorf("save: %w", leaf))},
		{"pointer enriched error", stderr.Join(plain, &leaf)},
		{"nil children", multiError{nil, plain, nil, leaf, nil}},
		{"empty branch before match", multiError{multiError{}, leaf}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.True(t, HasType(tt.err, "database"))
			require.Equal(t, "save", Tag(tt.err, "operation"))
			require.Equal(t, "Unable to save your changes.", UIError(tt.err))
		})
	}
}

func TestMetadataErrorTreePrecedence(t *testing.T) {
	first := New("first").WithType("first").WithTag("key", "first").WithUIError("first UI")
	second := New("second").WithType("second").WithTag("key", "second").WithUIError("second UI")
	for _, tt := range []struct {
		name string
		err  error
		want string
		ui   string
	}{
		{"left sibling first", stderr.Join(first, second), "first", "first UI"},
		{"reversed siblings", stderr.Join(second, first), "second", "second UI"},
		{"depth before next sibling", stderr.Join(fmt.Errorf("wrapped: %w", stderr.Join(first)), second), "first", "first UI"},
		{"outer metadata first", New("outer").WithType("outer").WithTag("key", "outer").WithUIError("outer UI").Wrap(stderr.Join(first, second)), "outer", "outer UI"},
		{"single chain keeps outer metadata", first.Wrap(second), "first", "first UI"},
		{"stop before later branch", stderr.Join(first, unvisitedError{}), "first", "first UI"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.True(t, HasType(tt.err, tt.want))
			require.Equal(t, tt.want, Tag(tt.err, "key"))
			require.Equal(t, tt.ui, UIError(tt.err))
		})
	}

	t.Run("type search continues past a different type", func(t *testing.T) {
		err := first.Wrap(stderr.Join(stderr.New("plain"), second))
		require.True(t, HasType(err, "second"))
	})
	t.Run("tag search continues past an absent or nil tag", func(t *testing.T) {
		err := New("outer").WithTag("key", nil).Wrap(stderr.Join(New("untagged"), second))
		require.Equal(t, "second", Tag(err, "key"))
	})
	for _, value := range []any{0, false, "", (*int)(nil)} {
		t.Run(fmt.Sprintf("non-nil tag value %T", value), func(t *testing.T) {
			err := stderr.Join(New("first").WithTag("key", value), second)
			require.Equal(t, value, Tag(err, "key"))
		})
	}
	t.Run("empty UI message is still a match", func(t *testing.T) {
		require.Equal(t, "", UIError(stderr.Join(uiErr(""), second)))
	})
	t.Run("first enriched error keeps its own UI fallback", func(t *testing.T) {
		fallback := New("internal error")
		require.Equal(t, fallback.Error(), UIError(stderr.Join(fallback, second)))
		err := fallback.Wrap(stderr.Join(first, second))
		require.Equal(t, err.Error(), UIError(err))
	})
}

func TestMetadataErrorTreeNoMatch(t *testing.T) {
	for _, tt := range []struct {
		name string
		err  error
	}{
		{"nil", nil},
		{"empty join", stderr.Join(nil, nil)},
		{"empty children", multiError{}},
		{"nil children", multiError{nil, nil}},
		{"plain join", stderr.Join(stderr.New("first"), stderr.New("second"))},
		{"wrapped join", fmt.Errorf("outer: %w", stderr.Join(stderr.New("first"), stderr.New("second")))},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.False(t, HasType(tt.err, "missing"))
			require.Nil(t, Tag(tt.err, "missing"))
			if tt.err == nil {
				require.True(t, HasType(tt.err, ""))
				require.Empty(t, UIError(tt.err))
			} else {
				// Nil children must not match the special top-level nil type.
				require.False(t, HasType(tt.err, ""))
				require.True(t, HasType(tt.err, Type(tt.err)))
				require.Equal(t, tt.err.Error(), UIError(tt.err))
			}
		})
	}
}

func TestUnwrapErrorTrees(t *testing.T) {
	joined := stderr.Join(stderr.New("first"), stderr.New("second"))
	// The public helper keeps the standard library's single-error semantics.
	require.Nil(t, Unwrap(joined))
	require.Nil(t, Unwrap(fmt.Errorf("%w and %w", stderr.New("first"), stderr.New("second"))))
	require.Nil(t, Unwrap(nil))
	require.Same(t, joined, Unwrap(fmt.Errorf("outer: %w", joined)))
}

func BenchmarkIs(b *testing.B) {
	b.Run("plain enriched error", func(b *testing.B) {
		err := New("err")
		target := err

		for b.Loop() {
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

		for b.Loop() {
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

		for b.Loop() {
			if !Is(err, target) {
				b.Fatal("expected match")
			}
		}
	})

	b.Run("wrapped enriched error", func(b *testing.B) {
		wrapped := New("wrapped").WithType("wrapped-code")
		err := New("err").Wrap(wrapped)
		target := wrapped

		for b.Loop() {
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

		for b.Loop() {
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

		for b.Loop() {
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

		for b.Loop() {
			if !err.Is(target) {
				b.Fatal("expected match")
			}
		}
	})

	b.Run("wrapped enriched error", func(b *testing.B) {
		wrapped := io.EOF
		err := New("err").Wrap(wrapped)
		target := err

		for b.Loop() {
			if !err.Is(target) {
				b.Fatal("expected match")
			}
		}
	})
}
