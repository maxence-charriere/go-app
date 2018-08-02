package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompoNameFromURL(t *testing.T) {
	tests := []struct {
		rawurl       string
		expectedName string
	}{
		{
			rawurl:       "/hello",
			expectedName: "hello",
		},
		{
			rawurl:       "/Hello",
			expectedName: "hello",
		},
		{
			rawurl:       "/hello?int=42",
			expectedName: "hello",
		},
		{
			rawurl:       "/hello/world",
			expectedName: "hello",
		},
		{
			rawurl:       "hello",
			expectedName: "hello",
		},
		{
			rawurl:       "main.hello",
			expectedName: "hello",
		},
		{
			rawurl:       "main.hello?foo=bar",
			expectedName: "hello",
		},
		{
			rawurl:       "hello?foo=bar",
			expectedName: "hello",
		},
		{
			rawurl: "test://hello",
		},
		{
			rawurl: "compo://",
		},
		{
			rawurl: "http://www.github.com",
		},
	}

	for _, test := range tests {
		name := CompoNameFromURLString(test.rawurl)
		assert.Equal(t, test.expectedName, name)
	}
}
