package app

import (
	"bytes"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

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
		})
	}
}

func TestPrintHTML(t *testing.T) {
	var b bytes.Buffer
	PrintHTML(&b, Div())
	require.Equal(t, "<div></div>", b.String())
}

func TestNodeManagerMount(t *testing.T) {
	ctx := makeTestContext()

	t.Run("mounting a text succeeds", func(t *testing.T) {
		var m nodeManager

		hello, err := m.Mount(ctx, 1, Text("hello"))
		require.NoError(t, err)
		require.NotZero(t, hello)
		require.True(t, hello.Mounted())
		require.Equal(t, "hello", hello.(*text).value)
		require.NotNil(t, hello.JSValue())
	})

	t.Run("mounting an already mounted text returns an error", func(t *testing.T) {
		var m nodeManager

		text, err := m.Mount(ctx, 1, Text("hello"))
		require.NoError(t, err)

		text, err = m.Mount(ctx, 1, text)
		require.Error(t, err)
		require.Zero(t, text)
		t.Log(err)
	})

	t.Run("mounting html succeeds", func(t *testing.T) {
		var m nodeManager

		div, err := m.Mount(ctx, 1, Div())
		require.NoError(t, err)
		require.NotZero(t, div)
		require.True(t, div.Mounted())
		require.NotNil(t, div.JSValue())
		require.Equal(t, uint(1), div.(HTML).depth())
	})

	t.Run("mounting html body succeeds", func(t *testing.T) {
		var m nodeManager

		body, err := m.Mount(ctx, 1, Body())
		require.NoError(t, err)
		require.NotZero(t, body)
		require.True(t, body.Mounted())
		require.NotNil(t, body.JSValue())
		require.Equal(t, uint(1), body.(HTML).depth())
	})

	t.Run("mounting html with attributes succeeds", func(t *testing.T) {
		var m nodeManager

		div, err := m.Mount(ctx, 1, Img().
			Class("test").
			Src("/web/test.webp"))
		require.NoError(t, err)
		require.True(t, div.Mounted())
	})

	t.Run("mounting html with event handlers succeeds", func(t *testing.T) {
		var m nodeManager
		var wg sync.WaitGroup

		elem, err := m.Mount(ctx, 1, Body().
			On("testgoapp", func(ctx Context, e Event) {
				wg.Done()
			}))
		require.NoError(t, err)
		require.True(t, elem.Mounted())

		if IsServer {
			return
		}

		wg.Add(1)
		customEvent := Window().Get("CustomEvent").New("testgoapp", map[string]any{
			"detail": "a js custom event",
		})
		elem.JSValue().Call("dispatchEvent", customEvent)
		wg.Wait()
	})

	t.Run("mounting html with children succeeds", func(t *testing.T) {
		var m nodeManager

		div, err := m.Mount(ctx, 1, Div().Body(
			Span(),
		))
		require.NoError(t, err)
		require.NotZero(t, div)

		body := div.(HTML).body()
		require.NotEmpty(t, body)

		span := body[0]
		require.True(t, span.Mounted())
		require.Equal(t, uint(2), span.(HTML).depth())
	})

	t.Run("mounting an already mounted html returns an error", func(t *testing.T) {
		var m nodeManager

		div, err := m.Mount(ctx, 1, Div())
		require.NoError(t, err)

		div, err = m.Mount(ctx, 1, div)
		require.Error(t, err)
		require.Zero(t, div)
		t.Log(err)
	})

	t.Run("mounting html with non mountable child returns an error", func(t *testing.T) {
		var m nodeManager

		div, err := m.Mount(ctx, 1, Div().Body(
			&compoWithNilRendering{},
		))
		require.Error(t, err)
		require.Nil(t, div)
		t.Log(err)
	})

	t.Run("mounting a component succeeds", func(t *testing.T) {
		var m nodeManager

		compo, err := m.Mount(ctx, 1, &hello{})
		require.NoError(t, err)
		require.NotNil(t, compo)
		require.True(t, compo.Mounted())
		require.Equal(t, uint(1), compo.(Composer).depth())
		if IsServer {
			require.True(t, compo.(*hello).preRendered)
		} else {
			require.True(t, compo.(*hello).mounted)
		}

		root := compo.(Composer).root()
		require.NotNil(t, root)
		require.IsType(t, Div(), root)
		require.True(t, root.Mounted())
		require.NotNil(t, root.(HTML).parent())
	})

	t.Run("mounting a component which renders nil returns an error", func(t *testing.T) {
		var m nodeManager

		compo, err := m.Mount(ctx, 1, &compoWithNilRendering{})
		require.Error(t, err)
		require.Nil(t, compo)
		t.Log(err)
	})

	t.Run("mounting an already mounted component returns an error", func(t *testing.T) {
		var m nodeManager

		compo, err := m.Mount(ctx, 1, &hello{})
		require.NoError(t, err)

		compo, err = m.Mount(ctx, 1, compo)
		require.Error(t, err)
		require.Nil(t, compo)
		t.Log(err)
	})

	t.Run("mounting a component with a non mountable root returns an error", func(t *testing.T) {
		var m nodeManager

		compo, err := m.Mount(ctx, 1, &compoWithNonMountableRoot{})
		require.Error(t, err)
		require.Nil(t, compo)
		t.Log(err)
	})

	t.Run("mounting raw html succeeds", func(t *testing.T) {
		var m nodeManager

		span, err := m.Mount(ctx, 1, Raw(`<span>hello</span>`))
		require.NoError(t, err)
		require.NotZero(t, span)
		require.NotNil(t, span.(*raw).jsElement)
		require.True(t, span.Mounted())
		require.Equal(t, "<span>hello</span>", span.(*raw).value)
	})

	t.Run("mounting an already mounted raw html returns an error", func(t *testing.T) {
		var m nodeManager

		span, err := m.Mount(ctx, 1, Raw(`<span>hello</span>`))
		require.NoError(t, err)

		span, err = m.Mount(ctx, 1, span)
		require.Error(t, err)
		require.Zero(t, span)
	})

	t.Run("mounting not supported element returns an error", func(t *testing.T) {
		var m nodeManager

		condition, err := m.Mount(ctx, 1, condition{})
		require.Error(t, err)
		require.Zero(t, condition)
	})
}

func BenchmarkNodeManagerMount(b *testing.B) {
	var m nodeManager
	ctx := makeTestContext()

	for n := 0; n < b.N; n++ {
		m.Mount(ctx, 1, Div())
	}
}

func TestNodeManagerDismount(t *testing.T) {
	ctx := makeTestContext()

	t.Run("html is dismounted", func(t *testing.T) {
		var m nodeManager

		div, err := m.Mount(ctx, 1, Div())
		require.NoError(t, err)

		m.Dismount(div)
		require.False(t, div.Mounted())
		require.Nil(t, div.JSValue())
	})

	t.Run("html child is dismounted", func(t *testing.T) {
		var m nodeManager

		div, err := m.Mount(ctx, 1, Div().Body(
			Span(),
		))
		require.NoError(t, err)
		span := div.(HTML).body()[0]

		m.Dismount(div)
		require.False(t, span.Mounted())
	})

	t.Run("html event handler is dismounted", func(t *testing.T) {
		var m nodeManager

		div, err := m.Mount(ctx, 1, Div().
			On("", func(ctx Context, e Event) {}))
		require.NoError(t, err)

		m.Dismount(div)
	})

	t.Run("component is dismounted", func(t *testing.T) {
		var m nodeManager

		compo, err := m.Mount(ctx, 1, &hello{})
		require.NoError(t, err)

		m.Dismount(compo)
		require.False(t, compo.Mounted())
		require.False(t, compo.(*hello).mounted)
		require.False(t, compo.(Composer).root().Mounted())
	})

	t.Run("raw html is dismounted", func(t *testing.T) {
		var m nodeManager

		span, err := m.Mount(ctx, 1, Raw("<span>hi</span>"))
		require.NoError(t, err)

		m.Dismount(span)
		require.Nil(t, span.(*raw).jsElement)
		require.False(t, span.Mounted())
	})
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

	t.Run("replacer component can be updated", func(t *testing.T) {
		var m nodeManager
		require.True(t, m.CanUpdate(&replacerComponent{}, &replacerComponent{}))
	})

	t.Run("replacer component cannot be updated", func(t *testing.T) {
		var m nodeManager
		require.False(t, m.CanUpdate(&replacerComponent{replace: true}, &replacerComponent{replace: true}))
	})
}

func BenchmarkNodeManagerCanUpdate(b *testing.B) {
	var m nodeManager
	for n := 0; n < b.N; n++ {
		m.CanUpdate(Div(), Div())
	}
}

func TestNodeManagerUpdate(t *testing.T) {
	ctx := makeTestContext()

	t.Run("updating a non mounted element returns an error", func(t *testing.T) {
		var m nodeManager

		_, err := m.Update(ctx, Div(), Div())
		require.Error(t, err)
		t.Log(err)
	})

	t.Run("updating text succeeds", func(t *testing.T) {
		var m nodeManager

		greeting, err := m.Mount(ctx, 1, Text("hello"))
		require.NoError(t, err)

		greeting, err = m.Update(ctx, greeting, Text("bye"))
		require.NoError(t, err)
		require.Equal(t, "bye", greeting.(*text).value)
	})

	t.Run("updating same text succeeds", func(t *testing.T) {
		var m nodeManager

		greeting, err := m.Mount(ctx, 1, Text("hello"))
		require.NoError(t, err)

		greeting, err = m.Update(ctx, greeting, Text("hello"))
		require.NoError(t, err)
		require.Equal(t, "hello", greeting.(*text).value)
	})

	t.Run("update html adds an attribute", func(t *testing.T) {
		var m nodeManager

		div, err := m.Mount(ctx, 1, Div())
		require.NoError(t, err)
		require.Empty(t, div.(HTML).attrs())

		div, err = m.Update(ctx, div, Div().Class("test"))
		require.NoError(t, err)
		require.Len(t, div.(HTML).attrs(), 1)
		require.Equal(t, "test", div.(HTML).attrs()["class"])

		div, err = m.Update(ctx, div, Div().
			Class("test").
			ID("test"))
		require.NoError(t, err)
		require.Len(t, div.(HTML).attrs(), 2)
		require.Equal(t, "test", div.(HTML).attrs()["id"])
	})

	t.Run("update html updates an attribute", func(t *testing.T) {
		var m nodeManager

		div, err := m.Mount(ctx, 1, Div().Class("hello"))
		require.NoError(t, err)
		require.Equal(t, "hello", div.(HTML).attrs()["class"])

		div, err = m.Update(ctx, div, Div().Class("bye"))
		require.NoError(t, err)
		require.Len(t, div.(HTML).attrs(), 1)
		require.Equal(t, "bye", div.(HTML).attrs()["class"])
	})

	t.Run("update html removes an attribute", func(t *testing.T) {
		var m nodeManager

		div, err := m.Mount(ctx, 1, Div().Class("hello"))
		require.NoError(t, err)
		require.Len(t, div.(HTML).attrs(), 1)
		require.Equal(t, "hello", div.(HTML).attrs()["class"])

		div, err = m.Update(ctx, div, Div())
		require.NoError(t, err)
		require.Empty(t, div.(HTML).attrs()["class"])
	})

	t.Run("update html adds an event handler", func(t *testing.T) {
		var m nodeManager

		div, err := m.Mount(ctx, 1, Div())
		require.NoError(t, err)
		require.Empty(t, div.(HTML).events())

		handler1 := func(ctx Context, e Event) {}
		div, err = m.Update(ctx, div, Div().OnClick(handler1))
		require.NoError(t, err)
		require.Len(t, div.(HTML).events(), 1)
		require.True(t, div.(HTML).events()["click"].Equal(eventHandler{
			event:     "click",
			goHandler: handler1,
		}))

		handler2 := func(ctx Context, e Event) {}
		div, err = m.Update(ctx, div, Div().
			OnClick(handler1).
			OnChange(handler2))
		require.NoError(t, err)
		require.Len(t, div.(HTML).events(), 2)
		require.True(t, div.(HTML).events()["change"].Equal(eventHandler{
			event:     "change",
			goHandler: handler2,
		}))
	})

	t.Run("update html updates an event handler", func(t *testing.T) {
		var m nodeManager

		handler1 := func(ctx Context, e Event) {}
		div, err := m.Mount(ctx, 1, Div().OnClick(handler1))
		require.NoError(t, err)

		handler2 := func(ctx Context, e Event) {}
		div, err = m.Update(ctx, div, Div().OnClick(handler2))
		require.NoError(t, err)
		require.Len(t, div.(HTML).events(), 1)
		require.False(t, div.(HTML).events()["click"].Equal(eventHandler{
			event:     "click",
			goHandler: handler1,
		}))
		require.True(t, div.(HTML).events()["click"].Equal(eventHandler{
			event:     "click",
			goHandler: handler2,
		}))
	})

	t.Run("udpate html removes an event handler", func(t *testing.T) {
		var m nodeManager

		div, err := m.Mount(ctx, 1, Div().OnClick(func(ctx Context, e Event) {}))
		require.NoError(t, err)
		require.Len(t, div.(HTML).events(), 1)

		div, err = m.Update(ctx, div, Div())
		require.NoError(t, err)
		require.Empty(t, div.(HTML).events())
	})

	t.Run("update html adds a child", func(t *testing.T) {
		var m nodeManager

		div, err := m.Mount(ctx, 1, Div())
		require.NoError(t, err)
		require.Empty(t, div.(HTML).body())

		div, err = m.Update(ctx, div, Div().Body(
			Span(),
		))
		require.NoError(t, err)
		require.Len(t, div.(HTML).body(), 1)
		require.IsType(t, Span(), div.(HTML).body()[0])
	})

	t.Run("update html by adding non mountable child returns an error", func(t *testing.T) {
		var m nodeManager

		div, err := m.Mount(ctx, 1, Div())
		require.NoError(t, err)

		div, err = m.Update(ctx, div, Div().Body(
			&compoWithNilRendering{},
		))
		require.Error(t, err)
		require.Zero(t, div)
		t.Log(err)
	})

	t.Run("update html updates a child", func(t *testing.T) {
		var m nodeManager

		div, err := m.Mount(ctx, 1, Div().Body(
			Span(),
		))
		require.NoError(t, err)
		require.Len(t, div.(HTML).body(), 1)
		child := div.(HTML).body()[0]
		require.IsType(t, Span(), child)
		require.Empty(t, child.(HTML).attrs())

		div, err = m.Update(ctx, div, Div().Body(
			Span().Class("test"),
		))
		require.NoError(t, err)
		require.Len(t, div.(HTML).body(), 1)
		child = div.(HTML).body()[0]
		require.IsType(t, Span(), child)
		require.Equal(t, "test", child.(HTML).attrs()["class"])
	})

	t.Run("update html by updating child with a non mountable element returns an error", func(t *testing.T) {
		var m nodeManager

		div, err := m.Mount(ctx, 1, Div().Body(
			Span().Text("test"),
		))
		require.NoError(t, err)

		div, err = m.Update(ctx, div, Div().Body(
			Span().Body(
				&compoWithNilRendering{},
			),
		))
		require.Error(t, err)
		require.Zero(t, div)
		t.Log(err)
	})

	t.Run("update html replaces a child", func(t *testing.T) {
		var m nodeManager

		div, err := m.Mount(ctx, 1, Div().Body(
			Span(),
		))
		require.NoError(t, err)
		require.Len(t, div.(HTML).body(), 1)
		child := div.(HTML).body()[0]
		require.IsType(t, Span(), child)

		div, err = m.Update(ctx, div, Div().Body(
			Div(),
		))
		require.NoError(t, err)
		require.Len(t, div.(HTML).body(), 1)
		child = div.(HTML).body()[0]
		require.IsType(t, Div(), child)
	})

	t.Run("update html by replacing child with a non mountable element returns an error", func(t *testing.T) {
		var m nodeManager

		div, err := m.Mount(ctx, 1, Div().Body(
			Span(),
		))
		require.NoError(t, err)

		div, err = m.Update(ctx, div, Div().Body(
			&compoWithNilRendering{},
		))
		require.Error(t, err)
		require.Zero(t, div)
		t.Log(err)
	})

	t.Run("update html removes a child", func(t *testing.T) {
		var m nodeManager

		div, err := m.Mount(ctx, 1, Div().Body(
			Span(),
		))
		require.NoError(t, err)
		require.Len(t, div.(HTML).body(), 1)

		div, err = m.Update(ctx, div, Div())
		require.NoError(t, err)
		require.Empty(t, div.(HTML).body())
	})

	t.Run("update component updates a field", func(t *testing.T) {
		var m nodeManager

		compo, err := m.Mount(ctx, 1, &bar{})
		require.NoError(t, err)
		require.NotNil(t, compo)
		require.Equal(t, "", compo.(Composer).root().(*text).value)

		updatedCompo, err := m.Update(ctx, compo, &bar{
			Value: "bar",
		})
		require.NoError(t, err)
		require.Equal(t, compo, updatedCompo)
		require.Equal(t, "bar", compo.(Composer).root().(*text).value)
	})

	t.Run("update component skips update", func(t *testing.T) {
		var m nodeManager

		compo, err := m.Mount(ctx, 1, &bar{
			Value: "bar",
		})
		require.NoError(t, err)
		require.Equal(t, "bar", compo.(Composer).root().(*text).value)

		updatedCompo, err := m.Update(ctx, compo, &bar{
			Value: "bar",
		})
		require.NoError(t, err)
		require.Equal(t, compo, updatedCompo)
		require.Equal(t, "bar", compo.(Composer).root().(*text).value)
	})

	t.Run("update component with non renderable component returns an error", func(t *testing.T) {
		var m nodeManager

		compo, err := m.Mount(ctx, 1, &compoWithCustomRoot{Root: Text("hi")})
		require.NoError(t, err)

		newCompo, err := m.Update(ctx, compo, &compoWithCustomRoot{Root: nil})
		require.Error(t, err)
		require.Nil(t, newCompo)
	})

	t.Run("update component with non mountable root returns an error", func(t *testing.T) {
		var m nodeManager

		compo, err := m.Mount(ctx, 1, &compoWithCustomRoot{Root: Div()})
		require.NoError(t, err)

		newCompo, err := m.Update(ctx, compo, &compoWithCustomRoot{Root: Div().Body(&compoWithNilRendering{})})
		require.Error(t, err)
		require.Nil(t, newCompo)
	})

	t.Run("update component replaces its root", func(t *testing.T) {
		var m nodeManager

		body, err := m.Mount(ctx, 1, Body().privateBody(
			&compoWithCustomRoot{Root: Div()}),
		)
		require.NoError(t, err)
		require.NotZero(t, body)
		require.NotEmpty(t, body.(HTML).body())

		compo := body.(HTML).body()[0]
		require.NotNil(t, compo)
		require.IsType(t, &compoWithCustomRoot{}, compo)
		require.IsType(t, Div(), compo.(*compoWithCustomRoot).Root)

		newCompo, err := m.Update(ctx, compo, &compoWithCustomRoot{Root: Span()})
		require.NoError(t, err)
		require.Equal(t, compo, newCompo)
		require.IsType(t, Span(), compo.(*compoWithCustomRoot).Root)
	})

	t.Run("update component with non mountable replaced component returns an error", func(t *testing.T) {
		var m nodeManager

		compo, err := m.Mount(ctx, 1, &compoWithCustomRoot{
			Root: &compoWithCustomRoot{Root: Div()},
		})
		require.NoError(t, err)
		require.NotNil(t, compo)

		newCompo, err := m.Update(ctx, compo, &compoWithCustomRoot{
			Root: &compoWithCustomRoot{Root: &compoWithNilRendering{}},
		})
		require.Error(t, err)
		require.Nil(t, newCompo)
	})

	t.Run("update raw html replaces its value", func(t *testing.T) {
		var m nodeManager

		body, err := m.Mount(ctx, 1, Body().privateBody(
			Raw("<span>hi</span>"),
		))
		require.NoError(t, err)
		require.NotZero(t, body)
		require.NotEmpty(t, body.(HTML).body())

		span := body.(HTML).body()[0]
		require.NotZero(t, span)
		require.IsType(t, Raw(""), span)

		newSpan, err := m.Update(ctx, span, Raw("<span>bye</span>"))
		require.NoError(t, err)
		require.NotEqual(t, span, newSpan)
		require.IsType(t, Raw(""), newSpan)
		require.Equal(t, "<span>bye</span>", newSpan.(*raw).value)
		require.True(t, newSpan.Mounted())
		require.False(t, span.Mounted())
	})

	t.Run("update raw html skips update", func(t *testing.T) {
		var m nodeManager

		body, err := m.Mount(ctx, 1, Body().privateBody(
			Raw("<span>hi</span>"),
		))
		require.NoError(t, err)
		require.NotZero(t, body)
		require.NotEmpty(t, body.(HTML).body())

		span := body.(HTML).body()[0]
		require.NotZero(t, span)
		require.IsType(t, Raw(""), span)

		newSpan, err := m.Update(ctx, span, Raw("<span>hi</span>"))
		require.NoError(t, err)
		require.Equal(t, span, newSpan)
		require.IsType(t, Raw(""), newSpan)
		require.Equal(t, "<span>hi</span>", newSpan.(*raw).value)
		require.True(t, newSpan.Mounted())
	})
}

func TestNodeManagerContext(t *testing.T) {
	var m nodeManager

	ctx := makeTestContext()
	div, err := m.Mount(ctx, 1, Div())
	require.NoError(t, err)

	ctx = m.context(ctx, div)
	require.NotZero(t, ctx)
	require.NotNil(t, ctx.Context)
	require.NotNil(t, ctx.page)
	require.NotNil(t, ctx.resolveURL)
	require.NotNil(t, ctx.localStorage)
	require.NotNil(t, ctx.sessionStorage)
	require.NotNil(t, ctx.dispatch)
	require.NotNil(t, ctx.defere)

	require.NotNil(t, ctx.sourceElement)
	require.NotNil(t, ctx.notifyComponentEvent)
}

func TestNodeManagerNotifyComponentEvent(t *testing.T) {
	ctx := makeTestContext()

	t.Run("nav event is notified", func(t *testing.T) {
		updates := make(map[UI]struct{})
		ctx.addComponentUpdate = func(c Composer, v int) {
			updates[c] = struct{}{}
		}

		var m nodeManager
		compo := &hello{}
		div, err := m.Mount(ctx, 1, Div().Body(compo))
		require.NoError(t, err)

		m.NotifyComponentEvent(ctx, div, nav{})
		require.NotEmpty(t, compo.onNavURL)
		require.Contains(t, updates, compo)
	})

	t.Run("app update event is notified", func(t *testing.T) {
		updates := make(map[UI]struct{})
		ctx.addComponentUpdate = func(c Composer, v int) {
			updates[c] = struct{}{}
		}

		var m nodeManager
		compo := &hello{}
		div, err := m.Mount(ctx, 1, Div().Body(compo))
		require.NoError(t, err)

		m.NotifyComponentEvent(ctx, div, appUpdate{})
		require.True(t, compo.appUpdated)
		require.Contains(t, updates, compo)
	})

	t.Run("app install change event is notified", func(t *testing.T) {
		updates := make(map[UI]struct{})
		ctx.addComponentUpdate = func(c Composer, v int) {
			updates[c] = struct{}{}
		}

		var m nodeManager
		compo := &hello{}
		div, err := m.Mount(ctx, 1, Div().Body(compo))
		require.NoError(t, err)

		m.NotifyComponentEvent(ctx, div, appInstallChange{})
		require.True(t, compo.appInstalled)
		require.Contains(t, updates, compo)
	})

	t.Run("resize change event is notified", func(t *testing.T) {
		updates := make(map[UI]struct{})
		ctx.addComponentUpdate = func(c Composer, v int) {
			updates[c] = struct{}{}
		}

		var m nodeManager
		compo := &hello{}
		div, err := m.Mount(ctx, 1, Div().Body(compo))
		require.NoError(t, err)

		m.NotifyComponentEvent(ctx, div, resize{})
		require.True(t, compo.appResized)
		require.Contains(t, updates, compo)
	})
}

func TestNodeManagerEncode(t *testing.T) {
	t.Run("encode indent with no depth", func(t *testing.T) {
		var m nodeManager
		var b bytes.Buffer

		m.encodeIndent(&b, 0)
		require.Empty(t, b.Bytes())
	})

	t.Run("encode indent with depth", func(t *testing.T) {
		var m nodeManager
		var b bytes.Buffer

		m.encodeIndent(&b, 2)
		require.Equal(t, "    ", b.String())
	})

	t.Run("encode empty text", func(t *testing.T) {
		var m nodeManager
		var b bytes.Buffer

		m.Encode(makeTestContext(), &b, Text(""))
		require.Empty(t, b.String())
	})

	t.Run("encode text", func(t *testing.T) {
		var m nodeManager
		var b bytes.Buffer

		m.Encode(makeTestContext(), &b, Text("hello"))
		require.Equal(t, "hello", b.String())
		t.Log(b.String())
	})

	t.Run("encode self closing html", func(t *testing.T) {
		var m nodeManager
		var b bytes.Buffer

		m.Encode(makeTestContext(), &b, Img())
		require.Equal(t, "<img>", b.String())
		t.Log(b.String())
	})

	t.Run("encode html", func(t *testing.T) {
		var m nodeManager
		var b bytes.Buffer

		m.Encode(makeTestContext(), &b, Div())
		require.Equal(t, "<div></div>", b.String())
		t.Log(b.String())
	})

	t.Run("encode html with attribute", func(t *testing.T) {
		var m nodeManager
		var b bytes.Buffer

		m.Encode(makeTestContext(), &b, Div().Class("test"))
		require.Equal(t, `<div class="test"></div>`, b.String())
		t.Log(b.String())
	})

	t.Run("encode html with empty class and id attributes", func(t *testing.T) {
		var m nodeManager
		var b bytes.Buffer

		m.Encode(makeTestContext(), &b, Div().
			Class("").
			ID(""))
		require.Equal(t, `<div></div>`, b.String())
		t.Log(b.String())
	})

	t.Run("encode html with a true boolean attribute", func(t *testing.T) {
		var m nodeManager
		var b bytes.Buffer

		m.Encode(makeTestContext(), &b, Script().Async(true))
		require.Equal(t, `<script async></script>`, b.String())
		t.Log(b.String())
	})

	t.Run("encode html with single text child", func(t *testing.T) {
		var m nodeManager
		var b bytes.Buffer

		m.Encode(makeTestContext(), &b, Div().Text("hello"))
		require.Equal(t, `<div>hello</div>`, b.String())
		t.Log(b.String())
	})

	t.Run("encode html with single non text child", func(t *testing.T) {
		var m nodeManager
		var b bytes.Buffer

		m.Encode(makeTestContext(), &b, Div().Body(
			Span(),
		))
		require.Equal(t, "<div>\n  <span></span>\n</div>", b.String())
		t.Log(b.String())
	})

	t.Run("encode html with multiple children", func(t *testing.T) {
		var m nodeManager
		var b bytes.Buffer

		m.Encode(makeTestContext(), &b, Div().Body(
			Text("hello"),
			Span().Body(
				Img(),
			),
		))
		require.Equal(t, "<div>\n  hello\n  <span>\n    <img>\n  </span>\n</div>", b.String())
		t.Log(b.String())
	})

	t.Run("encode component", func(t *testing.T) {
		var m nodeManager
		var b bytes.Buffer

		m.Encode(makeTestContext(), &b, &compoWithCustomRoot{Root: Div()})
		require.Equal(t, "<div></div>", b.String())
		t.Log(b.String())
	})

	t.Run("encode nested component", func(t *testing.T) {
		var m nodeManager
		var b bytes.Buffer

		m.Encode(makeTestContext(), &b, Span().Body(
			&compoWithCustomRoot{Root: Div()},
		))
		require.Equal(t, "<span>\n  <div></div>\n</span>", b.String())
		t.Log(b.String())
	})

	t.Run("encode empty raw", func(t *testing.T) {
		var m nodeManager
		var b bytes.Buffer

		m.Encode(makeTestContext(), &b, Raw(""))
		require.Equal(t, "<div></div>", b.String())
		t.Log(b.String())
	})

	t.Run("encode raw", func(t *testing.T) {
		var m nodeManager
		var b bytes.Buffer

		m.Encode(makeTestContext(), &b, Raw("<img>"))
		require.Equal(t, "<img>", b.String())
		t.Log(b.String())
	})

	t.Run("encode nested raw", func(t *testing.T) {
		var m nodeManager
		var b bytes.Buffer

		m.Encode(makeTestContext(), &b, Div().Body(
			Raw("<img>"),
		))
		require.Equal(t, "<div>\n  <img>\n</div>", b.String())
		t.Log(b.String())
	})
}

func TestCanUpdateValue(t *testing.T) {
	utests := []struct {
		a         any
		b         any
		canUpdate bool
	}{
		{
			a:         "hello",
			b:         "hello",
			canUpdate: false,
		},
		{
			a:         "hello",
			b:         "bye",
			canUpdate: true,
		},
		{
			a:         true,
			b:         true,
			canUpdate: false,
		},
		{
			a:         true,
			b:         false,
			canUpdate: true,
		},
		{
			a:         true,
			b:         true,
			canUpdate: false,
		},
		{
			a:         true,
			b:         false,
			canUpdate: true,
		},
		{
			a:         42,
			b:         42,
			canUpdate: false,
		},
		{
			a:         42,
			b:         21,
			canUpdate: true,
		},
		{
			a:         42.1,
			b:         42.1,
			canUpdate: false,
		},
		{
			a:         42.42,
			b:         42.0,
			canUpdate: true,
		},
		{
			a:         func() time.Time { return time.Date(2023, 2, 14, 0, 0, 0, 0, time.UTC) }(),
			b:         func() time.Time { return time.Date(2023, 2, 14, 0, 0, 0, 0, time.UTC) }(),
			canUpdate: false,
		},
		{
			a:         func() time.Time { return time.Date(2023, 2, 14, 0, 0, 0, 0, time.UTC) }(),
			b:         func() time.Time { return time.Date(1986, 2, 14, 0, 0, 0, 0, time.UTC) }(),
			canUpdate: true,
		},
		{
			a:         []int{1, 2, 3},
			b:         []int{1, 2, 3},
			canUpdate: false,
		},
		{
			a:         []int{1, 2, 3},
			b:         []int{3, 2, 1},
			canUpdate: true,
		},
	}

	for _, u := range utests {
		t.Run(fmt.Sprintf("%s can update %v", reflect.TypeOf(u.a), u.canUpdate), func(t *testing.T) {
			res := canUpdateValue(reflect.ValueOf(u.a), reflect.ValueOf(u.b))
			require.Equal(t, u.canUpdate, res)
		})
	}
}

func TestComponent(t *testing.T) {
	t.Run("parent component is returned", func(t *testing.T) {
		compo := &compoWithCustomRoot{Root: Div()}

		var m nodeManager
		_, err := m.Mount(makeTestContext(), 1, compo)
		require.NoError(t, err)

		c, ok := component(compo)
		require.True(t, ok)
		require.Equal(t, compo, c)
	})

	t.Run("parent component is returned", func(t *testing.T) {
		div := Div()
		compo := &compoWithCustomRoot{Root: div}

		var m nodeManager
		_, err := m.Mount(makeTestContext(), 1, compo)
		require.NoError(t, err)

		c, ok := component(div)
		require.True(t, ok)
		require.Equal(t, compo, c)
	})

	t.Run("no component returned", func(t *testing.T) {
		div := Div()

		var m nodeManager
		_, err := m.Mount(makeTestContext(), 1, div)
		require.NoError(t, err)

		c, ok := component(div)
		require.False(t, ok)
		require.Nil(t, c)
	})
}
