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

func (n jsNode) updateText(s string) error {
	n.Set("nodeValue", s)
	return nil
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

func (n jsNode) upsertAttr() error {
	panic("not implemented")
}

func (n jsNode) deleteAttr() error {
	panic("not implemented")
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
