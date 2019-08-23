// +build !js

package maestro

type jsNode struct {
}

func (n jsNode) new(tag, namespace string)    {}
func (n jsNode) newText()                     {}
func (n jsNode) change(tag, namespace string) {}
func (n jsNode) updateText(s string)          {}
func (n jsNode) appendChild(c jsNode)         {}
func (n jsNode) removeChild(c jsNode)         {}
func (n jsNode) upsertAttr(k, v string)       {}
func (n jsNode) deleteAttr(k string)          {}
func (n jsNode) addToBody()                   {}
