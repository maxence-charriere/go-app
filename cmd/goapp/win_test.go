package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWinSDKDirectory(t *testing.T) {
	defer os.RemoveAll("winSDKTest")

	expected := filepath.Join("winSDKTest", "10.0.8484")

	dirs := []string{
		filepath.Join("winSDKTest", "10.0.4242"),
		expected,
		filepath.Join("winSDKTest", "10.0.2121"),
	}

	for _, d := range dirs {
		os.MkdirAll(d, 0755)
	}

	sdkDir := winSDKDirectory("winSDKTest")
	require.Equal(t, expected, sdkDir)
}
