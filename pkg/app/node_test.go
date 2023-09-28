package app

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterUIElems(t *testing.T) {
	t.Run("filter empty elements returns nil", func(t *testing.T) {
		require.Nil(t, FilterUIElems())
	})

	t.Run("nil pointer is removed", func(t *testing.T) {
		require.Empty(t, FilterUIElems(nil))
	})

	t.Run("nil element is removed", func(t *testing.T) {
		var foo *foo
		require.Empty(t, FilterUIElems(foo))
	})

	t.Run("condition is inserted", func(t *testing.T) {
		elems := FilterUIElems(
			Div(),
			If(true, func() UI {
				return Span()
			}),
		)
		require.Len(t, elems, 2)
		require.IsType(t, Div(), elems[0])
		require.IsType(t, Span(), elems[1])
	})

	t.Run("condition is removed", func(t *testing.T) {
		elems := FilterUIElems(
			Div(),
			If(false, func() UI {
				return Span()
			}),
		)
		require.Len(t, elems, 1)
		require.IsType(t, Div(), elems[0])
	})

	t.Run("range is inserted", func(t *testing.T) {
		slice := []UI{Span()}

		elems := FilterUIElems(
			Div(),
			Range(slice).Slice(func(i int) UI {
				return slice[i]
			}),
			Div(),
		)
		require.Len(t, elems, 3)
		require.IsType(t, Div(), elems[0])
		require.IsType(t, Span(), elems[1])
		require.IsType(t, Div(), elems[2])
	})

	t.Run("range is removed", func(t *testing.T) {
		var slice []UI

		elems := FilterUIElems(
			Div(),
			Range(slice).Slice(func(i int) UI {
				return slice[i]
			}),
			Div(),
		)
		require.Len(t, elems, 2)
		require.IsType(t, Div(), elems[0])
		require.IsType(t, Div(), elems[1])
	})

	t.Run("no elements are removed", func(t *testing.T) {
		foo := &foo{}
		div := Div()
		text := Text("hello")
		raw := Raw("<br>")

		elems := FilterUIElems(
			foo,
			div,
			text,
			raw,
		)
		require.Len(t, elems, 4)
		require.Equal(t, foo, elems[0])
		require.Equal(t, div, elems[1])
		require.Equal(t, text, elems[2])
		require.Equal(t, raw, elems[3])
	})
}

func BenchmarkFilterUIElems(b *testing.B) {
	for n := 0; n < b.N; n++ {
		FilterUIElems(Div().
			Class("shell").
			Body(
				H1().Class("title").
					Text("Hello"),
				Input().
					Type("text").
					Class("in").
					Value("World").
					Placeholder("Type a name.").
					OnChange(func(ctx Context, e Event) {
						fmt.Println("Yo!")
					}),
			))
	}
}

type mountTest struct {
	scenario string
	node     UI
}

func testMountDismount(t *testing.T, utests []mountTest) {
	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			n := u.node

			d := NewClientTester(n)

			d.Consume()
			testMounted(t, n)

			d.Close()
			testDismounted(t, n)
		})
	}
}

func testMounted(t *testing.T, n UI) {
	require.NotNil(t, n.JSValue())
	require.NotNil(t, n.getDispatcher())
	require.True(t, n.Mounted())

	switch n.(type) {
	case *text, *raw:

	default:
		require.NoError(t, n.getContext().Err())
		require.NotNil(t, n.self())
	}

	for _, c := range n.getChildren() {
		require.Equal(t, n, c.getParent())
		testMounted(t, c)
	}
}

func testDismounted(t *testing.T, n UI) {
	require.NotNil(t, n.getDispatcher())
	require.False(t, n.Mounted())

	switch n.(type) {
	case *text, *raw:

	default:
		require.Error(t, n.getContext().Err())
		require.Nil(t, n.self())
	}

	for _, c := range n.getChildren() {
		testDismounted(t, c)
	}
}

type updateTest struct {
	scenario string
	a        UI
	b        UI
	matches  []TestUIDescriptor
}

func testUpdate(t *testing.T, utests []updateTest) {
	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			d := NewClientTester(u.a)
			defer d.Close()
			d.Consume()

			err := update(u.a, u.b)
			require.NoError(t, err)

			d.Consume()

			for _, d := range u.matches {
				require.NoError(t, TestMatch(u.a, d))
			}
		})
	}
}

func TestHTMLString(t *testing.T) {
	utests := []struct {
		scenario string
		root     UI
	}{
		{
			scenario: "hmtl element",
			root:     Div().ID("test"),
		},
		{
			scenario: "text",
			root:     Text("hello"),
		},
		{
			scenario: "component",
			root:     &hello{},
		},
		{
			scenario: "nested component",
			root:     Div().Body(&hello{}),
		},
		{
			scenario: "nested nested component",
			root:     Div().Body(&foo{Bar: "bar"}),
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			t.Log(HTMLString(u.root))
			t.Log(HTMLStringWithIndent(u.root))
		})
	}
}
