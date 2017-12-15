package appjs

import "html/template"
import "fmt"

//go:generate npx babel src/app.js --out-file src/appjs
//go:generate npx babel src/test.js --out-file test/test.js
//go:generate go run gen.go
//go:generate go fmt

// AppJS returns a string that contains the app.js script and use the given
// function to perform a go bridge request.
func AppJS(golangRequest string) template.JS {
	return template.JS(fmt.Sprintf(`
function golangRequest(payload) {
	%s(payload);
}
%s`, golangRequest, appjs))
}
