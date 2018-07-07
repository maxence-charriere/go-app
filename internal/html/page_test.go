package html

import (
	"testing"

	"github.com/murlokswarm/app"
)

func TestPage(t *testing.T) {
	p := Page(app.HTMLConfig{})
	t.Log(p)
}
