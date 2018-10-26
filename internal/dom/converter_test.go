package dom

import (
	"html/template"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRawHTML(t *testing.T) {
	s := rawHTML("hello")
	assert.Equal(t, template.HTML("hello"), s)
}

func TestCompoHTMLTag(t *testing.T) {
	s := compoHTMLTag("test.hello")
	assert.Equal(t, template.HTML("<test.hello>"), s)
}

func TestTimeFormat(t *testing.T) {
	s := timeFormat(time.Date(1986, 2, 14, 0, 0, 0, 0, time.UTC), "2006")
	assert.Equal(t, "1986", s)
}

func TestJSONFormat(t *testing.T) {
	s, _ := jsonFormat(42)
	assert.Equal(t, "42", s)
}
