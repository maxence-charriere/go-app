package dom

import (
	"testing"
)

func TestPageString(t *testing.T) {
	p := Page{
		GoRequest:     "alert",
		RootCompoName: "hello.hello",
	}

	t.Log(p)
}
