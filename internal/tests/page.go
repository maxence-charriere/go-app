package tests

// // PageTester is the interface that wraps the NewTestPage inteface.
// type PageTester interface {
// 	// NewTestPage creates a page for test.
// 	NewTestPage(c app.PageConfig) (app.Page, error)
// }

// func testPage(t *testing.T, d app.Driver) {
// 	tester, ok := d.(PageTester)
// 	if !ok {
// 		return
// 	}

// 	tests := []struct {
// 		scenario string
// 		config   app.PageConfig
// 		function func(t *testing.T, w app.Page)
// 	}{
// 		{
// 			scenario: "create",
// 		},
// 		{
// 			scenario: "create with a default component",
// 			config: app.PageConfig{
// 				DefaultURL: "tests.hello",
// 			},
// 		},
// 		{
// 			scenario: "url",
// 			function: testPageURL,
// 		},
// 		{
// 			scenario: "referer",
// 			function: testPageReferer,
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.scenario, func(t *testing.T) {
// 			p, err := tester.NewTestPage(test.config)
// 			if err == app.ErrNotSupported {
// 				return
// 			}
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			if test.function == nil {
// 				return
// 			}
// 			test.function(t, p)
// 		})
// 	}

// 	testElemWithCompo(t, func() app.ElemWithCompo {
// 		p, _ := tester.NewTestPage(app.PageConfig{})
// 		return p
// 	})

// 	testElementWithNavigation(t, func() app.Navigator {
// 		p, _ := tester.NewTestPage(app.PageConfig{})
// 		return p
// 	})
// }

// func testPageURL(t *testing.T, p app.Page) {
// 	t.Log(p.URL())
// }

// func testPageReferer(t *testing.T, p app.Page) {
// 	t.Log(p.Referer())
// }
