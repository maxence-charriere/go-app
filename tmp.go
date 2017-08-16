package app

// FilePicker is a struct that describes a file picker.
// It will be used by a driver to create a native file picker that allow to
// select files and directories filenames.
type FilePicker struct {
	MultipleSelection bool
	NoDir             bool
	NoFile            bool
	OnPick            func(filenames []string)
}

// Share is a struct that describes a share.
// It will be used by a driver to create a native share panel.
type Share struct {
	Value interface{}
}

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
