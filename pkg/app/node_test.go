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

func TestNodeManagerMount(t *testing.T) {
	t.Run("mounting a text succeeds", func(t *testing.T) {
		var m nodeManager

		hello, err := m.Mount(1, Text("hello"))
		require.NoError(t, err)
		require.NotZero(t, hello)
		require.True(t, hello.Mounted())
		require.Equal(t, "hello", hello.(*text).value)
		require.NotNil(t, hello.JSValue())
	})

	t.Run("mounting an already mounted text returns an error", func(t *testing.T) {
		var m nodeManager

		text, err := m.Mount(1, Text("hello"))
		require.NoError(t, err)

		text, err = m.Mount(1, text)
		require.Error(t, err)
		require.Zero(t, text)
	})

	t.Run("mounting an html element succeeds", func(t *testing.T) {
		var m nodeManager

		div, err := m.Mount(1, Div())
		require.NoError(t, err)
		require.NotZero(t, div)
		require.True(t, div.Mounted())
		require.NotNil(t, div.JSValue())
		require.Equal(t, uint(1), div.(HTML).Depth())
	})

	t.Run("mounting an html element with attributes succeeds", func(t *testing.T) {
		var m nodeManager

		div, err := m.Mount(1, Img().
			Class("test").
			Src("/web/test.webp"))
		require.NoError(t, err)
		require.True(t, div.Mounted())
	})

	t.Run("mounting an html element with children succeeds", func(t *testing.T) {
		var m nodeManager

		div, err := m.Mount(1, Div().Body(
			Span(),
		))
		require.NoError(t, err)
		require.NotZero(t, div)

		body := div.(HTML).body()
		require.NotEmpty(t, body)

		span := body[0]
		require.True(t, span.Mounted())
		require.Equal(t, uint(2), span.(HTML).Depth())
	})

	t.Run("mounting an already mounted html element returns an error", func(t *testing.T) {
		var m nodeManager

		div, err := m.Mount(1, Div())
		require.NoError(t, err)

		div, err = m.Mount(1, div)
		require.Error(t, err)
		require.Zero(t, div)
	})
}

func BenchmarkNodeManagerMount(b *testing.B) {
	var m nodeManager

	for n := 0; n < b.N; n++ {
		m.Mount(1, Div())
	}
}

func TestNodeManagerCanUpdate(t *testing.T) {
	t.Run("elements with same type can be updated", func(t *testing.T) {
		var m nodeManager
		require.True(t, m.CanUpdate(Div(), Div()))
	})

	t.Run("elements with different types cannot be updated", func(t *testing.T) {
		var m nodeManager
		require.False(t, m.CanUpdate(Div(), Span()))
	})

	t.Run("generic html elements with same tag can be updated", func(t *testing.T) {
		var m nodeManager
		require.True(t, m.CanUpdate(Elem("div"), Elem("div")))
	})

	t.Run("generic html elements with different tag cannot be updated", func(t *testing.T) {
		var m nodeManager
		require.False(t, m.CanUpdate(Elem("div"), Elem("span")))
	})

	t.Run("generic self closing html elements with same tag can be updated", func(t *testing.T) {
		var m nodeManager
		require.True(t, m.CanUpdate(ElemSelfClosing("input"), ElemSelfClosing("input")))
	})

	t.Run("generic self closing html elements with different tag cannot be updated", func(t *testing.T) {
		var m nodeManager
		require.False(t, m.CanUpdate(ElemSelfClosing("input"), ElemSelfClosing("br")))
	})
}

func BenchmarkNodeManagerCanUpdate(b *testing.B) {
	var m nodeManager
	for n := 0; n < b.N; n++ {
		m.CanUpdate(Div(), Div())
	}
}

func TestNodeManagerUpdate(t *testing.T) {
	t.Run("updating text succeeds", func(t *testing.T) {
		var m nodeManager

		greeting, err := m.Mount(1, Text("hello"))
		require.NoError(t, err)

		greeting, err = m.Update(greeting, Text("bye"))
		require.NoError(t, err)
		require.Equal(t, "bye", greeting.(*text).value)
	})

	t.Run("updating same text succeeds", func(t *testing.T) {
		var m nodeManager

		greeting, err := m.Mount(1, Text("hello"))
		require.NoError(t, err)

		greeting, err = m.Update(greeting, Text("hello"))
		require.NoError(t, err)
		require.Equal(t, "hello", greeting.(*text).value)
	})
}
