package app

import "testing"

func TestFlowPreRender(t *testing.T) {
	testSkipWasm(t)

	utests := []struct {
		scenario string
		flow     UI
	}{
		{
			scenario: "empty flow",
			flow:     Flow(),
		},
		{
			scenario: "flow with content",
			flow:     Flow().Content(Div()),
		},
		{
			scenario: "flow with multiple content",
			flow: Flow().Content(
				Div(),
				Div(),
				Div(),
				Div(),
			),
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			d := NewServerTester(u.flow)
			defer d.Close()
			d.PreRender()
		})
	}
}
