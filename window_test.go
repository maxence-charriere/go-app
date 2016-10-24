package app

import "testing"

func TestNewWindow(t *testing.T) {
	w := Window{}
	t.Log(NewWindow(w))
}
