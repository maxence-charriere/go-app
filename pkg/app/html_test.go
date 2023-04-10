package app

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTMLElementHTML(t *testing.T) {
	t.Run("srcset", func(t *testing.T) {
		e := Img().Src("/web/test.jpg").
			SrcSet("/web/test.webp").
			SrcSet("/web/test/jpg")

		var b strings.Builder
		e.html(&b)

		html := b.String()

		require.Equal(t, 1, strings.Count(html, "/web/test.webp"))
		t.Log(html)
	})
}

func BenchmarkMountHTMLElement(b *testing.B) {
	for n := 0; n < b.N; n++ {
		client := NewClientTester(Div().
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
		client.Close()
	}
}

func BenchmarkHTMLElementHTML(b *testing.B) {
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

	for n := 0; n < b.N; n++ {
		var bytes bytes.Buffer
		div.html(&bytes)
	}
}

func BenchmarkHTMLElementHTMLIndent(b *testing.B) {
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

	for n := 0; n < b.N; n++ {
		var bytes bytes.Buffer
		div.htmlWithIndent(&bytes, 0)
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
