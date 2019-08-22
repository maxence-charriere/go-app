package maestro

import (
	"syscall/js"
)

type jsNode struct {
	js.Value
}

func (n jsNode) new(tag, namespace string) {
	var v js.Value
	if namespace != "" {
		v = js.Global().Call("createElementNS", namespace, tag)
	} else {
		v = js.Global().Call("createElement", tag)
	}
	if v.Type() == js.TypeUndefined {
		panic("createElement returned an undefined value")
	}
	n.Value = v
}

func (n jsNode) newText() {
	v := js.Global().Call("createTextNode", "")
	if v.Type() == js.TypeUndefined {
		panic("createTextNode returned an undefined value")
	}
	n.Value = v
}

func (n jsNode) change(tag, namespace string) {
	parent := n.Get("parentNode")
	if t := parent.Type(); t == js.TypeUndefined || t == js.TypeNull {
		panic("parentNode is not set")
	}

	var v js.Value
	if tag == "" {
		v = js.Global().Call("createTextNode", "")
	} else if namespace != "" {
		v = js.Global().Call("createElementNS", namespace, tag)
	} else {
		v = js.Global().Call("createElement", tag)
	}
	if v.Type() == js.TypeUndefined {
		panic("replace returned an undefined value")
	}

	parent.Call("replaceChild", v, n.Value)
	n.Value = v
}

func (n jsNode) updateText(s string) {
	n.Set("nodeValue", s)
}

func (n jsNode) appendChild(c JSNode) {
	n.Call("appendChild", c.(jsNode).Value)
}

func (n jsNode) removeChild(c JSNode) {
	n.Call("removeChild", c.(jsNode).Value)
}

func (n jsNode) upsertAttr(k, v string) {
	n.Call("setAttribute", k, v)
}

func (n jsNode) deleteAttr(k string) {
	n.Call("removeAttribute", k)
}

func (n jsNode) addToBody() {
	body := js.Global().Get("body")
	if body.Type() == js.TypeUndefined {
		panic("body is undefined")
	}

	body.Call("appendChild", n.Value)
}
