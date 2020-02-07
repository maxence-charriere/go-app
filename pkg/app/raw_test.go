package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRawOpenTag(t *testing.T) {
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
			tag := rawOpenTag(test.raw)
			require.Equal(t, test.expected, tag)
		})
	}
}

func TestRaw(t *testing.T) {
	tests := []struct {
		scenario string
		raw      string
		panic    bool
	}{
		{
			scenario: "missing opening tag",
			raw:      "</div>",
			panic:    true,
		},
		{
			scenario: "missing closing tag",
			raw:      "<div>",
			panic:    true,
		},
		{
			scenario: "different closing tag",
			raw:      "<div></span>",
			panic:    true,
		},
		{
			scenario: "well formed value",
			raw:      "<div></div>",
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			if test.panic {
				require.Panics(t, func() {
					Raw(test.raw)
				})
				return
			}

			require.NotNil(t, Raw(test.raw))
		})
	}
}
