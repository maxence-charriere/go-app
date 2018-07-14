package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		expectsErr  bool
	}{
		{
			scenario: "fetching current url in empty history returns an error",
			actions: []historyAction{
				{current, ""},
			},
			expectsErr: true,
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
			scenario: "get previous entry from first entry returns an error",
			actions: []historyAction{
				{new, "hello"},
				{new, "big"},
				{new, "world"},
				{previous, ""},
				{previous, ""},
				{previous, ""},
			},
			expectsErr: true,
		},
		{
			scenario: "get previous entry from empty history returns an error",
			actions: []historyAction{
				{previous, ""},
			},
			expectsErr: true,
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
			scenario: "get next entry from the last entry returns an error",
			actions: []historyAction{
				{new, "hello"},
				{new, "big"},
				{new, "world"},
				{next, ""},
			},
			expectsErr: true,
		},
		{
			scenario: "get next entry from empty history returns an error",
			actions: []historyAction{
				{next, ""},
			},
			expectsErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			var url string
			var currentURL string
			var err error

			h := NewHistory()

			for _, action := range test.actions {
				switch action.name {
				case current:
					url, err = h.Current()

				case new:
					url = action.url
					h.NewEntry(action.url)

				case canPrevious:
					h.CanPrevious()

				case previous:
					url, err = h.Previous()

				case canNext:
					h.CanNext()

				case next:
					url, err = h.Next()
				}

				if err != nil {
					break
				}
			}

			if test.expectsErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			currentURL, _ = h.Current()
			assert.Equal(t, url, currentURL)
			assert.Equal(t, test.expectedURL, url)
			assert.Equal(t, test.expectedLen, h.Len())
		})
	}
}
