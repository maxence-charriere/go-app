package app

import "testing"

func TestCondition(t *testing.T) {
	testUpdate(t, []updateTest{
		{
			scenario: "if is interpreted",
			a: Div().Body(
				If(false,
					H1(),
				),
			),
			b: Div().Body(
				If(true,
					H1(),
				),
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
				If(true,
					H1(),
				),
			),
			b: Div().Body(
				If(false,
					H1(),
				),
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
				If(true,
					H1(),
				).ElseIf(false,
					H2(),
				),
			),
			b: Div().Body(
				If(false,
					H1(),
				).ElseIf(true,
					H2(),
				),
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
				If(false,
					H1(),
				).ElseIf(true,
					H2(),
				),
			),
			b: Div().Body(
				If(false,
					H1(),
				).ElseIf(false,
					H2(),
				),
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
				If(false,
					H1(),
				).ElseIf(true,
					H2(),
				).Else(
					H3(),
				),
			),
			b: Div().Body(
				If(false,
					H1(),
				).ElseIf(false,
					H2(),
				).Else(
					H3(),
				),
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
				If(false,
					H1(),
				).ElseIf(true,
					H2(),
				).Else(
					H3(),
				),
			),
			b: Div().Body(
				If(true,
					H1(),
				).ElseIf(false,
					H2(),
				).Else(
					H3(),
				),
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
