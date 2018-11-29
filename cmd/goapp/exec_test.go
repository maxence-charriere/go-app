package main

import (
	"context"
	"runtime"
	"testing"
)

func TestExec(t *testing.T) {
	cmd := []string{"ls", "-l"}
	if runtime.GOOS == "windows" {
		cmd = []string{"powershell", "ls"}
	}

	if err := execute(context.Background(), cmd[0], cmd[1:]...); err != nil {
		t.Fatal(err)
	}
}
