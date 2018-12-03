package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestValidateWinFileTypes(t *testing.T) {
	tests := []struct {
		scenario string
		fileType winFileType
		err      bool
	}{
		{
			scenario: "valid file type",
			fileType: winFileType{
				Name: "test",
				Extensions: []winFileExtension{
					{Ext: ".test"},
				},
			},
		},
		{
			scenario: "no name returns an error",
			fileType: winFileType{},
			err:      true,
		},
		{
			scenario: "no extensions returns an error",
			fileType: winFileType{Name: "test"},
			err:      true,
		},
		{
			scenario: "extension without '.' prefix returns an error",
			fileType: winFileType{
				Name: "test",
				Extensions: []winFileExtension{
					{Ext: "test"},
				},
			},
			err: true,
		},
	}

	for _, test := range tests {
		t.Run("test.scenario", func(t *testing.T) {
			err := validateWinFileTypes([]winFileType{test.fileType})

			if test.err {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
