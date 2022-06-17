package app

import (
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
