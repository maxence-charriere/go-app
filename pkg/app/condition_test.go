package app

import "testing"

func TestCondition(t *testing.T) {
	testUpdate(t, []updateTest{
		{
			scenario: "if is interpreted",
			a: Div().Body(
				If(false, func() UI {
					return H1()
				}),
			),
			b: Div().Body(
				If(true, func() UI {
					return H1()
				}),
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},
				{
					Path:     TestPath(0),
					Expected: H1(),
				},
			},
		},
		{
			scenario: "if is not interpreted",
			a: Div().Body(
				If(true, func() UI {
					return H1()
				}),
			),
			b: Div().Body(
				If(false, func() UI {
					return H1()
				}),
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},
				{
					Path:     TestPath(0),
					Expected: nil,
				},
			},
		},
		{
			scenario: "else if is interpreted",
			a: Div().Body(
				If(true, func() UI {
					return H1()
				}).ElseIf(false, func() UI {
					return H2()
				}),
			),
			b: Div().Body(
				If(false, func() UI {
					return H1()
				}).ElseIf(true, func() UI {
					return H2()
				}),
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},

				{
					Path:     TestPath(0),
					Expected: H2(),
				},
			},
		},
		{
			scenario: "else if is not interpreted",
			a: Div().Body(
				If(false, func() UI {
					return H1()
				}).ElseIf(true, func() UI {
					return H2()
				}),
			),
			b: Div().Body(
				If(false, func() UI {
					return H1()
				}).ElseIf(false, func() UI {
					return H2()
				}),
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},

				{
					Path:     TestPath(0),
					Expected: nil,
				},
			},
		},
		{
			scenario: "else is interpreted",
			a: Div().Body(
				If(false, func() UI {
					return H1()
				}).ElseIf(true, func() UI {
					return H2()
				}).Else(func() UI {
					return H3()
				}),
			),
			b: Div().Body(
				If(false, func() UI {
					return H1()
				}).ElseIf(false, func() UI {
					return H2()
				}).Else(func() UI {
					return H3()
				}),
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},

				{
					Path:     TestPath(0),
					Expected: H3(),
				},
			},
		},
		{
			scenario: "else is not interpreted",
			a: Div().Body(
				If(false, func() UI {
					return H1()
				}).ElseIf(true, func() UI {
					return H2()
				}).Else(func() UI {
					return H3()
				}),
			),
			b: Div().Body(
				If(true, func() UI {
					return H1()
				}).ElseIf(false, func() UI {
					return H2()
				}).Else(func() UI {
					return H3()
				}),
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},

				{
					Path:     TestPath(0),
					Expected: H1(),
				},
			},
		},
	})
}

func BenchmarkCondition(b *testing.B) {
	for n := 0; n < b.N; n++ {
		If(true, func() UI {
			return Div()
		})
	}
}
