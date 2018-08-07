package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type historyAction struct {
	name string
	url  string
}

func TestHistory(t *testing.T) {
	const (
		current     = "Current"
		new         = "NewEntry"
		canPrevious = "CanPrevious"
		previous    = "Previous"
		canNext     = "CanNext"
		next        = "Next"
	)

	tests := []struct {
		scenario    string
		actions     []historyAction
		expectedLen int
		expectedURL string
	}{
		{
			scenario: "fetching current url in empty history",
			actions: []historyAction{
				{current, ""},
			},
		},
		{
			scenario: "set a new entry",
			actions: []historyAction{
				{new, "hello"},
			},
			expectedURL: "hello",
			expectedLen: 1,
		},
		{
			scenario: "set a new empty entry",
			actions: []historyAction{
				{new, ""},
			},
			expectedURL: "",
			expectedLen: 0,
		},
		{
			scenario: "set a multiple new entries",
			actions: []historyAction{
				{new, "hello"},
				{new, "big"},
				{new, "world"},
			},
			expectedURL: "world",
			expectedLen: 3,
		},
		{
			scenario: "get previous entry",
			actions: []historyAction{
				{new, "hello"},
				{new, "big"},
				{new, "world"},
				{canPrevious, ""},
				{previous, ""},
			},
			expectedURL: "big",
			expectedLen: 3,
		},
		{
			scenario: "get previous entry from first entry",
			actions: []historyAction{
				{new, "hello"},
				{new, "big"},
				{new, "world"},
				{previous, ""},
				{previous, ""},
				{previous, ""},
			},
			expectedURL: "",
			expectedLen: 3,
		},
		{
			scenario: "get previous entry from empty history",
			actions: []historyAction{
				{previous, ""},
			},
		},
		{
			scenario: "get next entry",
			actions: []historyAction{
				{new, "hello"},
				{new, "big"},
				{new, "world"},
				{previous, ""},
				{previous, ""},
				{canNext, ""},
				{next, ""},
			},
			expectedURL: "big",
			expectedLen: 3,
		},
		{
			scenario: "get next entry from the last entry",
			actions: []historyAction{
				{new, "hello"},
				{new, "big"},
				{new, "world"},
				{next, ""},
			},
			expectedURL: "",
			expectedLen: 3,
		},
		{
			scenario: "get next entry from empty history",
			actions: []historyAction{
				{next, ""},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			var url string

			h := NewHistory()

			for _, action := range test.actions {
				switch action.name {
				case current:
					url = h.Current()

				case new:
					url = action.url
					h.NewEntry(action.url)
					assert.Equal(t, action.url, h.Current())

				case canPrevious:
					h.CanPrevious()

				case previous:
					url = h.Previous()

				case canNext:
					h.CanNext()

				case next:
					url = h.Next()
				}
			}

			assert.Equal(t, test.expectedURL, url)
			assert.Equal(t, test.expectedLen, h.Len())
		})
	}
}
