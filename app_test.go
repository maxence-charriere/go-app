package app

import (
	"testing"

	"github.com/murlokswarm/log"
)

type Hello struct {
	Greeting      string
	BadMarkup     bool
	BadMarkupSync bool
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
    <input onchange="OnInputChange" />

	{{if .BadMarkup}}<div></span>{{end}}

	{{if .BadMarkupSync}}
		<div></span>
	{{else}}
		<div>Foo</div>
	{{end}}
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
	ctx := newTestContext("rendering")
	defer ctx.Close()

	hello := &Hello{}
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

func TestRenderSyncError(t *testing.T) {
	ctx := newTestContext("rendering")
	defer ctx.Close()

	hello := &Hello{}
	ctx.Mount(hello)
	hello.BadMarkupSync = true
	Render(hello)
}

func TestMenuBar(t *testing.T) {
	t.Log(MenuBar())
}

func TestDock(t *testing.T) {
	t.Log(Dock())
}

func TestStorage(t *testing.T) {
	t.Log(Storage())
	t.Log(Resources())
}
