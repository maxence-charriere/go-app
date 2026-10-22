package app

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkMountHTMLElement(b *testing.B) {
	var m nodeManager
	ctx := makeTestContext()
	makeElement := func() HTML {
		return Div().
			Class("shell").
			Body(
				H1().Class("title").
					Text("Hello"),
				Input().
					Type("text").
					Class("in").
					Value("World").
					Placeholder("Type a name.").
					OnChange(func(Context, Event) {}),
			)
	}
	check, err := m.Mount(ctx, 1, makeElement())
	require.NoError(b, err)
	require.True(b, check.Mounted())
	require.Len(b, check.(HTML).body(), 2)
	m.Dismount(check)
	require.False(b, check.Mounted())
	b.ReportAllocs()
	// Include construction, mounting, and dismounting. Releasing event
	// callbacks on every iteration keeps browser benchmarks bounded too.
	for b.Loop() {
		elem, err := m.Mount(ctx, 1, makeElement())
		if err != nil {
			b.Fatal(err)
		}
		m.Dismount(elem)
	}
}

func BenchmarkHTMLElementHTML(b *testing.B) {
	var m nodeManager
	ctx := makeTestContext()

	div := Div().
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
		)

	var check bytes.Buffer
	m.Encode(ctx, &check, div)
	require.Contains(b, check.String(), `<h1 class="title">Hello</h1>`)
	require.Contains(b, check.String(), `value="World"`)
	b.ReportAllocs()
	b.SetBytes(int64(check.Len()))
	for b.Loop() {
		var out bytes.Buffer
		m.Encode(ctx, &out, div)
	}
}

func TestKeyCondition(t *testing.T) {
	utests := []struct {
		key      string
		value    string
		expected bool
	}{
		{
			key:      "class",
			value:    "hi",
			expected: false,
		},
		{
			key:      "id",
			value:    "hi",
			expected: false,
		},
		{
			key:      "style",
			value:    "",
			expected: false,
		},
		{
			key:      "class",
			value:    "",
			expected: true,
		},
		{
			key:      "id",
			value:    "",
			expected: true,
		},
	}

	for i, u := range utests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			require.Equal(t, u.expected, (u.key == "id" || u.key == "class") && u.value == "")
		})
	}
}
