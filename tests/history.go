package tests

import (
	"testing"

	"github.com/murlokswarm/app"
)

type historyAction struct {
	name string
	url  string
}

// TestHistory is a test suite used to ensure that all history implementations
// behave the same.
func TestHistory(t *testing.T, newHistory func() app.History) {
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

			history := newHistory()

			for _, action := range test.actions {
				switch action.name {
				case current:
					url, err = history.Current()

				case new:
					url = action.url
					history.NewEntry(action.url)

				case canPrevious:
					history.CanPrevious()

				case previous:
					url, err = history.Previous()

				case canNext:
					history.CanNext()

				case next:
					url, err = history.Next()
				}

				if err != nil {
					break
				}
			}

			if err == nil && test.expectsErr {
				t.Fatal("error is nil")
			}
			if err != nil && test.expectsErr {
				return
			}

			currentURL, _ = history.Current()
			if currentURL != url {
				t.Errorf("current url is not %v: %v", url, currentURL)
			}

			if url != test.expectedURL {
				t.Errorf("expected url is not %v: %v", test.expectedURL, url)
			}
			if l := history.Len(); l != test.expectedLen {
				t.Errorf("expected len is not %v: %v", test.expectedLen, l)
			}
		})
	}
}
