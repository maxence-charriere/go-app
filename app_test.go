package app

import (
	"testing"

	"github.com/murlokswarm/log"
)

type Hello struct {
	Greeting  string
	BadMarkup bool
}

func (h *Hello) OnInputChange(a ChangeArg) {
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
    <input _onchange="OnInputChange" />

	{{if .BadMarkup}}<div></span>{{end}}
</div>
    `
}

func init() {
	RegisterComponent(&Hello{})
}

func TestRun(t *testing.T) {
	OnFinalize = func() {
		log.Info("OnFinalize called")
	}

	Run()
}

func TestRender(t *testing.T) {
	hello := &Hello{}

	ctx := NewZeroContext("rendering")
	defer ctx.Close()

	ctx.Mount(hello)
	hello.Greeting = "Maxence"
	Render(hello)
}

func TestRenderPanicCompoCtxError(t *testing.T) {
	defer func() { recover() }()

	hello := &Hello{}
	Render(hello)
	t.Error("should panic")
}

func TestMenu(t *testing.T) {
	t.Log(Menu())
}

func TestDock(t *testing.T) {
	t.Log(Dock())
}
