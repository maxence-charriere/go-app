package dom

import (
	"fmt"
	"net/url"
	"strings"
)

// Transform perform transformation for a given attribute.
type Transform func(name, value string) (string, string)

// JsToGoHandler convert a javascript handler to a go component handler.
func JsToGoHandler(name, value string) (string, string) {
	if !strings.HasPrefix(name, "on") {
		return name, value
	}

	if strings.HasPrefix(value, "js:") {
		return name, strings.TrimPrefix(value, "js:")
	}

	return name, fmt.Sprintf("callCompoHandler(this, event, '%s')", value)
}

// HrefCompoFmt format href attribute without scheme to target a component.
func HrefCompoFmt(name, value string) (string, string) {
	if name != "href" {
		return name, value
	}

	u, err := url.Parse(value)
	if err != nil {
		return name, value
	}

	if len(u.Scheme) != 0 {
		return name, value
	}

	if !strings.HasPrefix(u.Path, "/") {
		u.Path = "/" + u.Path
	}

	u.Scheme = "compo"
	return name, u.String()
}
