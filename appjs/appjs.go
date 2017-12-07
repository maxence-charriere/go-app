package appjs

import "html/template"

//go:generate npx babel src/app.js --out-file src/appjs
//go:generate npx babel src/test.js --out-file test/test.js
//go:generate go run gen.go
//go:generate go fmt

// AppJS returns a string that contains the app.js script.
func AppJS() template.JS {
	return appjs
}
