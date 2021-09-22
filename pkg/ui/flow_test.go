package ui

import (
	"testing"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func TestFlowPreRender(t *testing.T) {
	utests := []struct {
		scenario string
		flow     app.UI
	}{
		{
			scenario: "empty flow",
			flow:     Flow(),
		},
		{
			scenario: "flow with content",
			flow:     Flow().Content(app.Div()),
		},
		{
			scenario: "flow with multiple content",
			flow: Flow().Content(
				app.Div(),
				app.Div(),
				app.Div(),
				app.Div(),
			),
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			d := app.NewServerTester(u.flow)
			defer d.Close()
			d.PreRender()
		})
	}
}
