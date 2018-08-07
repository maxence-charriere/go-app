package html

import "testing"

func TestPage(t *testing.T) {
	tests := []struct {
		scenario string
		config   PageConfig
	}{
		{
			scenario: "returns a page from default config",
		},
		{
			scenario: "returns a page from filled config",
			config: PageConfig{
				Title: "page test",
				Metas: []Meta{
					{
						Name:    DescriptionMeta,
						Content: "A test page.",
					},
					{
						HTTPEquiv: RefreshMeta,
						Content:   "42",
					},
				},
				DefaultCompo:    "<div></div>",
				DisableAppStyle: true,
				CSS: []string{
					"hello.css",
					"world.css",
				},
				AppJS: "alert('some javascript code!')",
				Javascripts: []string{
					"hello.js",
					"world.js",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			page := NewPage(test.config)
			if len(page) == 0 {
				t.Fatal("page is empty")
			}
			t.Log(page)
		})
	}
}
