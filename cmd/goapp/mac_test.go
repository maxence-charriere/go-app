package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateMacFileTypes(t *testing.T) {
	tests := []struct {
		scenario string
		fileType macFileType
		err      bool
	}{
		{
			scenario: "valid file type",
			fileType: macFileType{
				Name: "test",
				Icon: "test.png",
				UTIs: []string{"public.png"},
			},
		},
		{
			scenario: "no name returns an error",
			fileType: macFileType{},
			err:      true,
		},
		{
			scenario: "no png icon returns an error",
			fileType: macFileType{
				Name: "test",
				Icon: "test.jpg",
			},
			err: true,
		},
		{
			scenario: "no uti returns an error",
			fileType: macFileType{
				Name: "test",
				Icon: "test.png",
			},
			err: true,
		},
		{
			scenario: "empty uti returns an error",
			fileType: macFileType{
				Name: "test",
				Icon: "test.png",
				UTIs: []string{""},
			},
			err: true,
		},
	}

	for _, test := range tests {
		t.Run("test.scenario", func(t *testing.T) {
			err := validateMacFileTypes(test.fileType)

			if test.err {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
