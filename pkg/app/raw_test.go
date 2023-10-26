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
