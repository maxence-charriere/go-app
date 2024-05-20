package app

import "testing"

func BenchmarkCondition(b *testing.B) {
	for n := 0; n < b.N; n++ {
		If(true, func() UI {
			return Div()
		})
	}
}
