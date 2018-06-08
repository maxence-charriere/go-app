package app

// func TestComponentNameFromURL(t *testing.T) {
// 	tests := []struct {
// 		rawurl       string
// 		expectedName string
// 	}{
// 		{
// 			rawurl:       "/hello",
// 			expectedName: "hello",
// 		},
// 		{
// 			rawurl:       "/hello?int=42",
// 			expectedName: "hello",
// 		},
// 		{
// 			rawurl:       "/hello/world",
// 			expectedName: "hello",
// 		},
// 		{
// 			rawurl:       "hello",
// 			expectedName: "hello",
// 		},
// 		{
// 			rawurl: "test://hello",
// 		},
// 		{
// 			rawurl: "compo://",
// 		},
// 		{
// 			rawurl: "http://www.github.com",
// 		},
// 	}

// 	for _, test := range tests {
// 		if name := ComponentNameFromURLString(test.rawurl); name != test.expectedName {
// 			t.Errorf(`name is not "%s": "%s"`, test.expectedName, name)
// 		}
// 	}
// }

// func TestNormalizeComponentName(t *testing.T) {
// 	if name := "lib.FooBar"; normalizeComponentName(name) != "lib.foobar" {
// 		t.Errorf("name is not lib.foobar: %s", name)
// 	}

// 	if name := "main.FooBar"; normalizeComponentName(name) != "foobar" {
// 		t.Errorf("name is not foobar: %s", name)
// 	}
// }
