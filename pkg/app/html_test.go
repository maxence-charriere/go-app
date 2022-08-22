package app

import (
	"bytes"
	"fmt"
	"testing"
)

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
