package tests

import (
	"testing"

	"github.com/murlokswarm/app"
)

func testPage(t *testing.T, d app.Driver) {
	tests := []struct {
		scenario string
		config   app.PageConfig
		function func(t *testing.T, w app.Page)
	}{
		{
			scenario: "create",
		},
		{
			scenario: "create with a default component",
			config: app.PageConfig{
				DefaultURL: "tests.hello",
			},
		},
		{
			scenario: "page is decorated with logs",
			function: testPageIsDecorated,
		},
		{
			scenario: "url",
			function: testPageURL,
		},
		{
			scenario: "referer",
			function: testPageReferer,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			p, err := d.NewPage(test.config)
			if app.NotSupported(err) {
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if test.function == nil {
				return
			}
			test.function(t, p)
		})
	}

	testElementWithComponent(t, func() (app.ElementWithComponent, error) {
		return d.NewPage(app.PageConfig{})
	})

	testElementWithNavigation(t, func() (app.ElementWithNavigation, error) {
		return d.NewPage(app.PageConfig{})
	})
}

func testPageIsDecorated(t *testing.T, p app.Page) {
	if base := p.Base(); base == p {
		t.Error("page is not decorated")
	}
}

func testPageURL(t *testing.T, p app.Page) {
	t.Log(p.URL())
}

func testPageReferer(t *testing.T, p app.Page) {
	t.Log(p.Referer())
}
