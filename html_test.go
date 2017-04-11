package app

import (
	"testing"

	"github.com/satori/go.uuid"
)

func TestHTMLContextHTML(t *testing.T) {
	c := HTMLContext{
		ID:       uuid.NewV1(),
		Title:    "Test",
		Lang:     "fr",
		MurlokJS: MurlokJS(),
		JS:       []string{"test.js"},
		CSS:      []string{"test.css", "test2.css"},
	}
	t.Log(c.HTML())
}
