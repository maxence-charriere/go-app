package app

import "testing"

func TestNewShare(t *testing.T) {
	NewShare(Share{
		Value: "Hello Murlok",
	})
}
