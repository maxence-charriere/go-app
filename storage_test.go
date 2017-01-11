package app

import (
	"os"
	"path/filepath"
	"testing"
)

type StorageTest string

func (s StorageTest) Resources() string {
	return "resources"
}

func (s StorageTest) CSS() string {
	return filepath.Join(s.Resources(), "css")
}

func (s StorageTest) JS() string {
	return filepath.Join(s.Resources(), "js")
}

func (s StorageTest) Storage() string {
	return filepath.Join(s.Resources(), "storage")
}

func (s StorageTest) Preferences() string {
	return filepath.Join(s.Resources(), "preferences")
}

func TestIsSupportedExtension(t *testing.T) {
	if name, ext := "foo.jpg", ".jpg"; !IsSupportedExtension(name, ext) {
		t.Errorf("%v not found in %v", ext, name)
	}
	if name, ext := "foo.jpg", "jpg"; IsSupportedExtension(name, ext) {
		t.Errorf("%v should not be a valid extension name", ext)
	}
	if name, ext := "foo.jpg", ".png"; IsSupportedExtension(name, ext) {
		t.Errorf("%v should have be found in %v", ext, name)
	}
}

func TestIsSupportedImageExtension(t *testing.T) {
	if name := "foo.jpg"; !IsSupportedImageExtension(name) {
		t.Errorf("%v should be supported", name)
	}

	if name := "foo.pnh"; IsSupportedImageExtension(name) {
		t.Errorf("%v should not be supported", name)
	}
}

func TestGetFilenamesWithExtensionsFromDir(t *testing.T) {
	dirname := "testdir"
	os.Mkdir(dirname, os.ModePerm)
	os.Mkdir(filepath.Join(dirname, "foo"), os.ModePerm)
	defer os.RemoveAll(dirname)

	filenames := []string{
		filepath.Join(dirname, "hello.css"),
		filepath.Join(dirname, "hello.js"),
		filepath.Join(dirname, "hello.png"),
	}
	for _, name := range filenames {
		f, _ := os.Create(name)
		defer f.Close()
	}

	names, err := GetFilenamesWithExtensionsFromDir(dirname, ".png", ".css")
	if err != nil {
		t.Fatal(err)
	}
	if l := len(names); l != 2 {
		t.Error("l should be 2:", l)
	}
	if name := names[0]; name != "hello.css" {
		t.Error("name should be hello.css:", name)
	}
	if name := names[1]; name != "hello.png" {
		t.Error("name should be hello.png:", name)
	}

	if _, err = GetFilenamesWithExtensionsFromDir("hello", ".jpg"); err == nil {
		t.Error("err should not be nil")
	}

	name := "foofile"
	f, _ := os.Create(name)
	defer os.Remove(name)
	defer f.Close()
	if _, err = GetFilenamesWithExtensionsFromDir(name, ".jpg"); err == nil {
		t.Error("err should not be nil")
	}
}
