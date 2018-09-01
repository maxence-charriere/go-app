package main

import (
	"context"
	"testing"
)

func TestExec(t *testing.T) {
	if err := execute(context.Background(), "ls", "-la"); err != nil {
		t.Fatal(err)
	}
}
