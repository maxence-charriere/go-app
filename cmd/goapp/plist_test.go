// +build darwin,amd64

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneratePlist(t *testing.T) {
	tests := []struct {
		filename string
		template string
	}{
		{
			filename: "Info.plist",
			template: plist,
		},
		{
			filename: ".entitlements",
			template: entitlements,
		},
	}

	for _, test := range tests {
		t.Run(test.filename, func(t *testing.T) {
			dir := "fixtures"
			err := os.Mkdir(dir, os.ModeDir|0755)
			require.NoError(t, err)
			defer os.RemoveAll(dir)

			filename := filepath.Join(dir, test.filename)

			err = generatePlist(filename, test.template, bundle{})
			assert.NoError(t, err)

			var b []byte
			b, err = ioutil.ReadFile(filename)
			assert.NoError(t, err)

			t.Log(string(b))
		})
	}
}
