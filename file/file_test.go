package file

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	tests := []struct {
		scenario  string
		src       string
		dst       string
		noSrc     bool
		createDst bool
		err       bool
	}{
		{
			scenario: "copy",
			src:      "src/test_copy",
			dst:      "dst/test_copy",
		},
		{
			scenario:  "copy dst already exists",
			src:       "src/test_copy_dst_exists",
			dst:       "dst/test_copy_dst_exists",
			createDst: true,
		},
		{
			scenario: "copy no src",
			dst:      "dst/test_copy_no_src",
			noSrc:    true,
			err:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			setup(t)
			defer teardown()

			if !test.noSrc {
				err := createFile(test.src)
				require.NoError(t, err)
			}

			if test.createDst {
				err := createFile(test.dst)
				require.NoError(t, err)
			}

			err := Copy(test.dst, test.src)
			if test.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.True(t, Matches(test.src, test.dst))
		})
	}
}

func TestSync(t *testing.T) {
	tests := []struct {
		scenario    string
		createDirs  []string
		createFiles []string
		src         string
		dst         string
	}{
		{
			scenario:    "sync files",
			createFiles: []string{"src/test"},
			src:         "src/test",
			dst:         "dst/test",
		},
		{
			scenario: "sync dirs",
			createFiles: []string{
				"src/test1",
				"src/test2",
				"src/test3",
			},
			src: "src",
			dst: "dst",
		},
		{
			scenario: "sync mixed",
			createDirs: []string{
				"src/foo",
				"src/bar",
				"src/boo",

				"dst/test2",
				"dst/bar",
				"dst/no",
			},
			createFiles: []string{
				"src/test",
				"src/test2",
				"src/foo/test",
				"src/bar/test",
				"src/boo/test",

				"dst/boo",
				"dst/test3",
			},
			src: "src",
			dst: "dst",
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			setup(t)
			defer teardown()

			for _, d := range test.createDirs {
				err := os.MkdirAll(d, os.ModeDir|os.ModePerm)
				require.NoError(t, err)
			}

			for _, f := range test.createFiles {
				err := createFile(f)
				require.NoError(t, err)
			}

			err := Sync(test.dst, test.src)
			assert.NoError(t, err)
		})
	}
}

func setup(t *testing.T) {
	_, err := os.Stat("src")
	if os.IsNotExist(err) {
		err = os.Mkdir("src", 0777)
	}
	require.NoError(t, err)

	if _, err = os.Stat("dst"); os.IsNotExist(err) {
		err = os.Mkdir("dst", 0777)
	}
	require.NoError(t, err)
}

func teardown() {
	os.RemoveAll("src")
	os.RemoveAll("dst")

}

func createFile(name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(uuid.New().String())
	return err
}
