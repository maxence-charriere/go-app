package app

import "testing"

func BenchmarkCondition(b *testing.B) {
	for n := 0; b.Loop(); n++ {
		If(true, func() UI {
			return Div()
		})
	}
}
