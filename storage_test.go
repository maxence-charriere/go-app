package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileHaveExtension(t *testing.T) {
	if name, ext := "foo.jpg", ".jpg"; !FileHaveExtension(name, ext) {
		t.Errorf("%v not found in %v", ext, name)
	}
	if name, ext := "foo.jpg", "jpg"; FileHaveExtension(name, ext) {
		t.Errorf("%v should not be a valid extension name", ext)
	}
	if name, ext := "foo.jpg", ".png"; FileHaveExtension(name, ext) {
		t.Errorf("%v should have be found in %v", ext, name)
	}
}

func TestFileIsSupportedIcon(t *testing.T) {
	if name := "foo.jpg"; !FileIsSupportedIcon(name) {
		t.Errorf("%v should be supported", name)
	}

	if name := "foo.pnh"; FileIsSupportedIcon(name) {
		t.Errorf("%v should not be supported", name)
	}
}

func TestGetFilenamesFromDir(t *testing.T) {
	dirname := "testdir"
	os.Mkdir(dirname, os.ModePerm)
	os.Mkdir(filepath.Join(dirname, "foo"), os.ModePerm)
	defer os.RemoveAll(dirname)

	filenames := []string{
		filepath.Join(dirname, "hello.css"),
		filepath.Join(dirname, "hello.js"),
		filepath.Join(dirname, "hello.png"),
		filepath.Join(dirname, "foo/hello.png"),
		filepath.Join(dirname, "foo/hello.css"),
	}
	for _, name := range filenames {
		f, _ := os.Create(name)
		defer f.Close()
	}

	names := GetFilenamesFromDir(dirname, ".png", ".css")
	if l := len(names); l != 4 {
		t.Error("l should be 4:", l)
	}
	t.Log(names)

	// No files with ext.
	if names = GetFilenamesFromDir("hello", ".jpg"); len(names) != 0 {
		t.Fatal("name should be empty")
	}

	// No directory
	name := "foofile"
	f, _ := os.Create(name)
	defer os.Remove(name)
	defer f.Close()

	if names = GetFilenamesFromDir(name, ".jpg"); len(names) != 0 {
		t.Fatal("name should be empty")
	}
}
