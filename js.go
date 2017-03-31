package app

import (
	"encoding/json"
	"fmt"

	"github.com/murlokswarm/log"
	"github.com/murlokswarm/markup"
	"github.com/satori/go.uuid"
)

const (
	jsFmt = `
function Mount(id, markup) {
    const sel = '[data-murlok-root="' + id + '"]';
    const elem = document.querySelector(sel);
    elem.innerHTML = markup;
}

function RenderFull(id, markup) {
    const sel = '[data-murlok-id="' + id + '"]';
    const elem = document.querySelector(sel);
    elem.outerHTML = markup;
}

function RenderAttributes(id, attrs) {
    const sel = '[data-murlok-id="' + id + '"]';
    const elem = document.querySelector(sel);

    for (var name in attrs) {
        if (elem.hasAttribute(name) && attrs[name].length == 0) {
            elem.removeAttribute(name);
            continue;
        }
        elem.setAttribute(name, attrs[name]);
    }
}

function GetAttributeValue(elem, name) {
    if (!elem.hasAttribute(name)) {
        return null;
    }
    return elem.getAttribute(name);
}

function CallEvent(id, method, self, event) {
    var arg;

    var value = null;
    if (self.value) {
        value = self.value;
    }

    switch (event.type) {
        case "click":
        case "contextmenu":
        case "dblclick":
        case "mousedown":
        case "mouseenter":
        case "mouseleave":
        case "mousemove":
        case "mouseover":
        case "mouseout":
        case "mouseup":
        case "drag":
        case "dragend":
        case "dragenter":
        case "dragleave":
        case "dragover":
        case "dragstart":
        case "drop":
            arg = MakeMouseArg(event);
            break;

        case "mousewheel":
            arg = MakeWheelArg(event);
            break;

        case "keydown":
        case "keypress":
        case "keyup":
            arg = MakeKeyboardArg(event);
            break;

        case "change":
            arg = MakeChangeArg(value);
            break;

        default:
            arg = {};
            break;
    }

    arg.Target = {
        ID: GetAttributeValue(self, "id"),
        Class: GetAttributeValue(self, "class"),
        Index: GetAttributeValue(self, "data-murlok-index"),
        Value: value,
        Tag: self.tagName.toLowerCase()
    };

    Call(id, method, arg);
}

function MakeMouseArg(event) {
    return {
        AltKey: event.altKey,
        Button: event.button,
        ClientX: event.clientX,
        ClientY: event.clientY,
        CtrlKey: event.ctrlKey,
        Detail: event.detail,
        MetaKey: event.metaKey,
        PageX: event.pageX,
        PageY: event.pageY,
        ScreenX: event.screenX,
        ScreenY: event.screenY,
        ShiftKey: event.shiftKey
    };
}

function MakeWheelArg(event) {
    return {
        DeltaX: event.deltaX,
        DeltaY: event.deltaY,
        DeltaZ: event.deltaZ,
        DeltaMode: event.deltaMode
    };
}

function MakeKeyboardArg(event) {
    return {
        AltKey: event.altKey,
        CtrlKey: event.ctrlKey,
        CharCode: event.charCode,
        KeyCode: event.keyCode,
        Location: event.location,
        MetaKey: event.metaKey,
        ShiftKey: event.shiftKey
    };
}

function MakeChangeArg(value) {
    return {
        Value: value
    };
}

function Call(id, method, arg) {
    let msg = {
        ID: id,
        Method: method,
        Arg: JSON.stringify(arg)
    };

    msg = JSON.stringify(msg);
    %v
}
    `
)

// DOMElement represents a DOM element.
type DOMElement struct {
	Tag   string // The tag of the element. e.g. div.
	ID    string // The id attribute.
	Class string // the class attribute.
	Value string // The value attribute.
	Index string // The data-murlok-index attribute.
}

type jsMsg struct {
	ID     uuid.UUID
	Method string
	Arg    string
}

// HandleEvent allows to call the component method or map the component field
// described in msg.
// Should be used only in a driver.
func HandleEvent(msg string) {
	var jsMsg jsMsg
	if err := json.Unmarshal([]byte(msg), &jsMsg); err != nil {
		log.Error(err)
		return
	}
	markup.HandleEvent(jsMsg.ID, jsMsg.Method, jsMsg.Arg)
}

// MurlokJS returns the javascript code allowing bidirectional communication
// between a context and it's webview.
// Should be used only in drivers implementations.
func MurlokJS() string {
	return fmt.Sprintf(jsFmt, driver.JavascriptBridge())
}
