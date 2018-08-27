package dom

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/stretchr/testify/assert"
)

func TestPage(t *testing.T) {
	p := Page(app.HTMLConfig{}, "alert", "")
	t.Log(p)
}

func TestCleanWindowsPath(t *testing.T) {
	assert.Equal(t,
		[]string{"a/b", "c/d"},
		cleanWindowsPath([]string{"a/b", `c\d`}),
	)
}
