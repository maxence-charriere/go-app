package app

import (
	"encoding/json"
	"fmt"

	"github.com/murlokswarm/log"
	"github.com/murlokswarm/markup"
	"github.com/murlokswarm/uid"
)

const (
	jsFmt = `
function Mount(id, markup) {
	const sel = '[data-murlok-root="' + id + '"]';
    const elem = document.querySelector(sel);
    elem.innerHTML = markup;
}

function Render(id, markup) {
	const sel = '[data-murlok-id="' + id + '"]';
    const elem = document.querySelector(sel);
    elem.outerHTML = markup;
}

function CallEvent(id, method, src, e) {
}

function Call(id, method, arg) {
	const msg = {
		ID: id,
		Method: method,
		Arg: arg
	};
	
	msg = JSON.stringify(msg);
	%v
}
    `
)

type jsMsg struct {
	ID     uid.ID
	Method string
	Arg    string
}

// CallComponentMethod calls component method described by msg.
// Should be used only in a driver.
func CallComponentMethod(msg string) {
	var jsMsg jsMsg

	if err := json.Unmarshal([]byte(msg), &jsMsg); err != nil {
		log.Error(err)
		return
	}

	if err := markup.Call(jsMsg.ID, jsMsg.Method, jsMsg.Arg); err != nil {
		log.Error(err)
	}
}

// MurlokJS returns the javascript code allowing bidirectional communication
// between a context and it's webview.
// Should be used only in drivers implementations.
func MurlokJS() string {
	return fmt.Sprintf(jsFmt, driver.JavascriptBridge())
}
