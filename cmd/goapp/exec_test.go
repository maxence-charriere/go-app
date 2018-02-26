package main

import "testing"

func TestExec(t *testing.T) {
	if err := execute("ls", "-la"); err != nil {
		t.Fatal(err)
	}
}
