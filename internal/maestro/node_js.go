package maestro

import (
	"syscall/js"
)

type jsNode struct {
	js.Value
}

func (n *jsNode) new(tag, namespace string) {
	var v js.Value
	if namespace != "" {
		v = js.Global().Get("document").Call("createElementNS", namespace, tag)
	} else {
		v = js.Global().Get("document").Call("createElement", tag)
	}
	n.Value = v
}

func (n *jsNode) newText() {
	v := js.Global().Get("document").Call("createTextNode", "")
	n.Value = v
}

func (n *jsNode) change(tag, namespace string) {
	parent := n.Get("parentNode")
	if t := parent.Type(); t == js.TypeUndefined || t == js.TypeNull {
		panic("parentNode is not set")
	}

	var v js.Value
	if tag == "" {
		v = js.Global().Get("document").Call("createTextNode", "")
	} else if namespace != "" {
		v = js.Global().Get("document").Call("createElementNS", namespace, tag)
	} else {
		v = js.Global().Get("document").Call("createElement", tag)
	}

	parent.Call("replaceChild", v, n.Value)
	n.Value = v
}

func (n *jsNode) updateText(s string) {
	n.Set("nodeValue", s)
}

func (n *jsNode) appendChild(c jsNode) {
	n.Call("appendChild", c.Value)
}

func (n *jsNode) removeChild(c jsNode) {
	n.Call("removeChild", c.Value)
}

func (n *jsNode) upsertAttr(k, v string) {
	n.Call("setAttribute", k, v)
}

func (n *jsNode) deleteAttr(k string) {
	n.Call("removeAttribute", k)
}

func (n *jsNode) addToBody() {
	body := js.Global().Get("document").Get("body")

	for {
		firstChild := body.Get("firstChild")
		if !firstChild.Truthy() {
			break
		}
		body.Call("removeChild", firstChild)
	}

	body.Call("appendChild", n.Value)
}
