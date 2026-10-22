package app

import "testing"

func BenchmarkCondition(b *testing.B) {
	for b.Loop() {
		If(true, func() UI {
			return Div()
		})
	}
}
