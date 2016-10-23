package app

import (
	"testing"

	"github.com/murlokswarm/log"
	"github.com/murlokswarm/markup"
	"github.com/murlokswarm/uid"
)

type Hello struct {
	Greeting string
}

func (h *Hello) OnInputChange(a OnChangeArg) {
	h.Greeting = a.Value
	Render(h)
}

func (h *Hello) Render() string {
	return `
<div>
    Hello, 
    <span>
        {{if .Greeting}}
            {{html .Greeting}}
        {{else}}
            World
        {{end}}
    </span>
    <input onchange="@OnInputChange" />
</div>
    `
}

func TestRun(t *testing.T) {
	OnFinalize = func() {
		log.Info("OnFinalize called")
	}

	Run()
	Run()
	Finalize()
	Finalize()
}

func TestRender(t *testing.T) {
	Run()
	defer Finalize()

	hello := &Hello{}
	markup.Mount(hello, uid.Context())

	hello.Greeting = "Maxence"
	Render(hello)
}
