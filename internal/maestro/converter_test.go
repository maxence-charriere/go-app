package maestro

import (
	"html/template"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBind(t *testing.T) {
	s := bind("Hello", 42)
	assert.Equal(t, template.JS("Hello.42"), s)
}

func TestCompoHTMLTag(t *testing.T) {
	tests := []struct {
		url         string
		expectedTag template.HTML
	}{
		{
			url:         "test.hello",
			expectedTag: "<test.hello>",
		},
		{
			url:         "/test.hello",
			expectedTag: "<test.hello>",
		},
		{
			url:         "/test.hello?value=42",
			expectedTag: `<test.hello value="42">`,
		},
		{
			url:         "/test.hello?Value=Hello+World",
			expectedTag: `<test.hello value="Hello World">`,
		},
	}

	for _, test := range tests {
		res := urlToHTMLTag(test.url)
		require.Equal(t, test.expectedTag, res)
	}
}

func TestJSONFormat(t *testing.T) {
	s, _ := jsonFormat(42)
	assert.Equal(t, "42", s)
}

func TestRawHTML(t *testing.T) {
	s := rawHTML("hello")
	assert.Equal(t, template.HTML("hello"), s)
}

func TestTimeFormat(t *testing.T) {
	s := timeFormat(time.Date(1986, 2, 14, 0, 0, 0, 0, time.UTC), "2006")
	assert.Equal(t, "1986", s)
}
