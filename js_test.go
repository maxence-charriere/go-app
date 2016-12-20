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

func TestCallComponentMethod(t *testing.T) {
	bar := &Bar{}

	ctx := NewZeroContext("test for call")
	defer ctx.Close()

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

	CallComponentMethod(string(msg))
}

func TestCallComponentJSONError(t *testing.T) {
	bar := &Bar{}

	ctx := NewZeroContext("test for call")
	defer ctx.Close()

	ctx.Mount(bar)
	msg := "}{}"

	CallComponentMethod(string(msg))
}

func TestCallComponentArgError(t *testing.T) {
	bar := &Bar{}

	ctx := NewZeroContext("test for call")
	defer ctx.Close()

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

	CallComponentMethod(string(msg))
}

func TestMurlokJS(t *testing.T) {
	t.Log(MurlokJS())
}
