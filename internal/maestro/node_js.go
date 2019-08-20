package maestro

import (
	"errors"
	"syscall/js"
)

type jsNode struct {
	js.Value
}

func (n jsNode) new(typ, namespace string) error {
	var v js.Value
	if namespace != "" {
		v = js.Global().Call("createElementNS", namespace, typ)
	} else {
		v = js.Global().Call("createElement", typ)
	}
	if v.Type() == js.TypeUndefined {
		return errors.New("createElement returned an undefined value")
	}
	return nil
}

func (n jsNode) newText(s string) error {
	v := js.Global().Call("createTextNode", s)
	if v.Type() == js.TypeUndefined {
		return errors.New("createTextNode returned an undefined value")
	}
	n.Value = v
	return nil
}

func (n jsNode) updateText(s string) {
	n.Set("nodeValue", s)
}

func (n jsNode) removeChild(c JSNode) {
	n.Call("removeChild", c)
}

func (n jsNode) changeType(typ, namespace string) error {
	parent := n.Get("parentNode")
	if t := parent.Type(); t == js.TypeUndefined || t == js.TypeNull {
		return errors.New("parentNode is not set")
	}

	var v js.Value
	if typ == "text" {
		v = js.Global().Call("createTextNode", "")
	} else if namespace != "" {
		v = js.Global().Call("createElementNS", namespace, typ)
	} else {
		v = js.Global().Call("createElement", typ)
	}
	if v.Type() == js.TypeUndefined {
		return errors.New("changing element type returned an undefined value")
	}

	parent.Call("replaceChild", v, n.Value)
	n.Value = v
	return nil
}

func (n jsNode) upsertAttr(k, v string) {
	n.Call("setAttribute", k, v)
}

func (n jsNode) deleteAttr(k string) {
	n.Call("removeAttribute", k)
}

func (n jsNode) delete() error {
	parent := n.Get("parentNode")
	if t := parent.Type(); t == js.TypeUndefined || t == js.TypeNull {
		return errors.New("parentNode is not set")
	}
	parent.Call("removeChild", n.Value)
	n.Value = js.Undefined()
	return nil
}
