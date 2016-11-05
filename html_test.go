package app

import (
	"testing"

	"github.com/murlokswarm/uid"
)

func TestHTMLContextHTML(t *testing.T) {
	c := HTMLContext{
		ID:       uid.Context(),
		Title:    "Test",
		Lang:     "fr",
		MurlokJS: MurlokJS(),
		JS:       []string{"test.js"},
		CSS:      []string{"test.css", "test2.css"},
	}

	t.Log(c.HTML())
}
