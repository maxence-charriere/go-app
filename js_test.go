package app

import (
	"encoding/json"
	"testing"

	"github.com/murlokswarm/log"
	"github.com/murlokswarm/markup"
)

type CallArg struct {
	Nb int
}

type Bar struct {
	placeholder string
}

func (b *Bar) Render() string {
	return `<div>Bar</div>`
}

func (b *Bar) OnCall(arg CallArg) {
	log.Infof("bar.OnCall(%+v) --> success", arg)
}

func init() {
	RegisterComponent(&Bar{})
}

func TestHandleEvent(t *testing.T) {
	ctx := newTestContext("test for call")
	defer ctx.Close()

	bar := &Bar{}
	ctx.Mount(bar)
	root := markup.Root(bar)

	jsMsg := jsMsg{
		ID:     root.ID,
		Method: "OnCall",
		Arg:    `{"Nb": 42}`,
	}
	msg, err := json.Marshal(jsMsg)
	if err != nil {
		t.Fatal(err)
	}
	HandleEvent(string(msg))
}

func TestCallComponentJSONError(t *testing.T) {
	ctx := newTestContext("test for call")
	defer ctx.Close()

	bar := &Bar{}
	ctx.Mount(bar)
	msg := "}{}"
	HandleEvent(string(msg))
}

func TestCallComponentArgError(t *testing.T) {
	ctx := newTestContext("test for call")
	defer ctx.Close()

	bar := &Bar{}
	ctx.Mount(bar)
	root := markup.Root(bar)

	jsMsg := jsMsg{
		ID:     root.ID,
		Method: "OnCall",
		Arg:    "42",
	}
	msg, err := json.Marshal(jsMsg)
	if err != nil {
		t.Fatal(err)
	}
	HandleEvent(string(msg))
}

func TestMurlokJS(t *testing.T) {
	t.Log(MurlokJS())
}
