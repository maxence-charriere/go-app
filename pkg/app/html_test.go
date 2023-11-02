package app

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkMountHTMLElement(b *testing.B) {
	var m nodeManager

	for n := 0; n < b.N; n++ {
		m.Mount(makeTestContext(), 1, Div().
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

// func BenchmarkHTMLElementHTML(b *testing.B) {
// 	div := Div().
// 		Class("shell").
// 		Body(
// 			H1().Class("title").
// 				Text("Hello"),
// 			Input().
// 				Type("text").
// 				Class("in").
// 				Value("World").
// 				Placeholder("Type a name.").
// 				OnChange(func(ctx Context, e Event) {
// 					fmt.Println("Yo!")
// 				}),
// 		)

// 	for n := 0; n < b.N; n++ {
// 		var bytes bytes.Buffer
// 		div.html(&bytes)
// 	}
// }

// func BenchmarkHTMLElementHTMLIndent(b *testing.B) {
// 	div := Div().
// 		Class("shell").
// 		Body(
// 			H1().Class("title").
// 				Text("Hello"),
// 			Input().
// 				Type("text").
// 				Class("in").
// 				Value("World").
// 				Placeholder("Type a name.").
// 				OnChange(func(ctx Context, e Event) {
// 					fmt.Println("Yo!")
// 				}),
// 		)

// 	for n := 0; n < b.N; n++ {
// 		var bytes bytes.Buffer
// 		div.htmlWithIndent(&bytes, 0)
// 	}
// }

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
