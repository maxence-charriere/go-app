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

function RenderFull(id, markup) {
	const sel = '[data-murlok-id="' + id + '"]';
    const elem = document.querySelector(sel);
    elem.outerHTML = markup;
}

function RenderAttributes(id, attrs) {
	const sel = '[data-murlok-id="' + id + '"]';
    const elem = document.querySelector(sel);
    
    for (var name in attrs) {
        if (attrs.hasOwnProperty(name)) {
            if (attrs[name].length == 0) {
                elem.removeAttribute(name);
                continue;
            }
            elem.setAttribute(name, attrs[name]);
        }
    }
}

function CallEvent(id, method, self, event) {
	var arg;
	const eventType = event.type;

	switch (eventType) {
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
			arg = MakeChangeArg(self.value);
			break;

        default:
			alert("not supported event: " + eventType);
            return;
	}
	
	Call(id, method, arg);
}

function MakeMouseArg(event) {
	return {
        "AltKey": event.altKey,
        "Button": event.button,
        "ClientX": event.clientX,
        "ClientY": event.clientY,
        "CtrlKey": event.ctrlKey,
        "Detail": event.detail,
        "MetaKey": event.metaKey,
        "PageX": event.pageX,
        "PageY": event.pageY,
        "ScreenX": event.screenX,
        "ScreenY": event.screenY,
        "ShiftKey": event.shiftKey
    };
}

function MakeWheelArg(event) {
	return {
        "DeltaX": event.deltaX,
        "DeltaY": event.deltaY,
        "DeltaZ": event.deltaZ,
        "DeltaMode": event.deltaMode
    };
}

function MakeKeyboardArg(event) {
	return {
        "AltKey": event.altKey,
        "CtrlKey": event.ctrlKey,
        "CharCode": event.charCode,
        "KeyCode": event.keyCode,
        "Location": event.location,
        "MetaKey": event.metaKey,
        "ShiftKey": event.shiftKey
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

	markup.Call(jsMsg.ID, jsMsg.Method, jsMsg.Arg)
}

// MurlokJS returns the javascript code allowing bidirectional communication
// between a context and it's webview.
// Should be used only in drivers implementations.
func MurlokJS() string {
	return fmt.Sprintf(jsFmt, driver.JavascriptBridge())
}
