package app

import (
	"fmt"
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

function Call(msg) {
	msg = JSON.stringify(msg);
	%v
}
    `
)

// MurlokJS returns the javascript code allowing bidirectional communication
// between a context and it's webview.
// Should be used only in drivers implementations.
func MurlokJS() string {
	return fmt.Sprintf(jsFmt, driver.JavascriptBridge())
}
