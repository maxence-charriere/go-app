package app

import (
	"bytes"
	"testing"
)

func TestCustom(t *testing.T) {
	svg := "http://www.w3.org/2000/svg"
	elem := Elem("svg")
	elem.XMLNS(svg)
	elem.Body(
		Elem("line").XMLNS(svg).Attr("x1", 0),
		ElemSC("selfclosing"))
	var buf bytes.Buffer
	elem.htmlWithIndent(&buf, 0)
	expect := `<svg xmlns="http://www.w3.org/2000/svg">
  <line xmlns="http://www.w3.org/2000/svg" x1="0"></line>
  <selfclosing>
</svg>`
	if buf.String() != expect {
		t.Log("Expected: '" + expect + "'")
		t.Log("Got: '" + buf.String() + "'")
		t.Fail()
	}
}
