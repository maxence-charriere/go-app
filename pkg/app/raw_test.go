package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRawRootTagName(t *testing.T) {
	tests := []struct {
		scenario string
		raw      string
		expected string
	}{
		{
			scenario: "tag set",
			raw: `
			<div>
				<span></span>
			</div>`,
			expected: "div",
		},
		{
			scenario: "tag is empty",
		},
		{
			scenario: "opening tag missing",
			raw:      "</div>",
		},
		{
			scenario: "tag is not set",
			raw:      "div",
		},
		{
			scenario: "tag is not closing",
			raw:      "<div",
		},
		{
			scenario: "tag is not closing",
			raw:      "<div",
		},
		{
			scenario: "tag without value",
			raw:      "<>",
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			tag := rawRootTagName(test.raw)
			require.Equal(t, test.expected, tag)
		})
	}
}

func TestRawMountDismount(t *testing.T) {
	testMountDismount(t, []mountTest{
		{
			scenario: "raw html element",
			node:     Raw(`<h1>Hello</h1>`),
		},
		{
			scenario: "raw svg element",
			node:     Raw(`<svg></svg>`),
		},
	})
}

func TestRawUpdate(t *testing.T) {
	testUpdate(t, []updateTest{
		{
			scenario: "raw html element is replace by another raw html element",
			a: Div().Body(
				Raw("<div></div>"),
			),
			b: Div().Body(
				Raw("<svg></svg>"),
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},
				{
					Path:     TestPath(0),
					Expected: Raw("<svg></svg>"),
				},
			},
		},
		{
			scenario: "raw html element is replace by non-raw html element",
			a: Div().Body(
				Raw("<div></div>"),
			),
			b: Div().Body(
				Text("hello"),
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},
				{
					Path:     TestPath(0),
					Expected: Text("hello"),
				},
			},
		},
	})
}
