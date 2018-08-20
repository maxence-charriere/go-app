package dom

import (
	"testing"

	"github.com/murlokswarm/app"
)

func TestPage(t *testing.T) {
	p := Page(app.HTMLConfig{}, "alert")
	t.Log(p)
}
