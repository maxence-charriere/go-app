package app

import (
	"bytes"
	"testing"
)

func BenchmarkHTML(b *testing.B) {
	var buffer bytes.Buffer

	root := Div().
		Class("Menu").
		Body(
			Text("☰"),
			Main().
				Class("Content").
				Body(
					H1().
						Body(
							Text("Hello,"),
							Text("world"),
						),
					Input().
						Placeholder("What is your name?").
						AutoFocus(true),
				),
		)

	for i := 0; i < b.N; i++ {
		buffer.Reset()
		root.html(&buffer)
	}
}

func BenchmarkHTMLWithRoot(b *testing.B) {
	var buffer bytes.Buffer

	for i := 0; i < b.N; i++ {
		buffer.Reset()

		root := Div().
			Class("Menu").
			Body(
				Text("☰"),
				Main().
					Class("Content").
					Body(
						H1().
							Body(
								Text("Hello,"),
								Text("world"),
							),
						Input().
							Placeholder("What is your name?").
							AutoFocus(true),
					),
			)

		root.html(&buffer)
	}
}
