package app

import (
	"testing"
)

type testCompo struct {
	Compo

	Num int    // export number
	num int    // unexport number
	Str string // export string
	str string // unexport string
}

func (c *testCompo) Render() UI {
	return Div().Body(
		Text(c.Str),
	)
}

func TestCompoUpdate(t *testing.T) {
	a := &testCompo{
		Num: 1,
		num: 2,
		Str: "a",
		str: "b",
	}
	b := &testCompo{
		Num: 3,
		num: 4,
		Str: "c",
		str: "d",
	}

	mount(a)
	a.update(b)

	if a.Num != 3 {
		t.Error("export number is not updated")
	}
	if a.num != 2 {
		t.Error("unexport number is updated")
	}
	if a.Str != "c" {
		t.Error("export string is not updated")
	}
	if a.str != "b" {
		t.Error("unexport string is updated")
	}
}
