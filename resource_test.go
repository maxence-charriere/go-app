package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResourcePathJoin(t *testing.T) {
	l := ResourcePath("resources")
	if j, exp := l.Join("css"), filepath.Join("resources", "css"); j != exp {
		t.Errorf("j should be %v: %v", exp, j)
	}
}

func TestResourcePathCSS(t *testing.T) {
	r := ResourcePath("resources")
	cssPath := r.Join("css")
	if err := os.MkdirAll(cssPath, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(Resources().Path())

	os.Mkdir(filepath.Join(cssPath, "bar"), os.ModePerm)
	ffoo, _ := os.Create(filepath.Join(cssPath, "foo.css"))
	fhello, _ := os.Create(filepath.Join(cssPath, "hello"))
	fworld, _ := os.Create(filepath.Join(cssPath, "world.txt"))
	defer ffoo.Close()
	defer fhello.Close()
	defer fworld.Close()

	cssFilenames := r.CSS()
	if l := len(cssFilenames); l != 1 {
		t.Error("cssFilenames should have 1 element:", l)
	}
	if exp := filepath.Join("css", "foo.css"); cssFilenames[0] != exp {
		t.Errorf("cssFilenames[0] should be %v: %v", exp, cssFilenames[0])
	}
}

func TestResourcePathCSSError(t *testing.T) {
	// No css directory.
	r := ResourcePath("resources")
	r.CSS()

	// css is not a directory.
	cssPath := r.Join("css")
	os.Mkdir(r.Path(), os.ModePerm)
	defer os.RemoveAll(Resources().Path())
	f, _ := os.Create(cssPath)
	defer f.Close()

	r.CSS()
}

func TestResourcePathJS(t *testing.T) {
	r := ResourcePath("resources")
	jsPath := r.Join("js")
	if err := os.MkdirAll(jsPath, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(Resources().Path())

	os.Mkdir(filepath.Join(jsPath, "bar"), os.ModePerm)
	ffoo, _ := os.Create(filepath.Join(jsPath, "foo.js"))
	fhello, _ := os.Create(filepath.Join(jsPath, "hello"))
	fworld, _ := os.Create(filepath.Join(jsPath, "world.txt"))
	defer ffoo.Close()
	defer fhello.Close()
	defer fworld.Close()

	jsFilenames := r.JS()
	if l := len(jsFilenames); l != 1 {
		t.Error("jsFilenames should have 1 element:", l)
	}
	if exp := filepath.Join("js", "foo.js"); jsFilenames[0] != exp {
		t.Errorf("jsFilenames[0] should be %v: %v", exp, jsFilenames[0])
	}
}

func TestResourcePathJSError(t *testing.T) {
	// No js directory.
	r := ResourcePath("resources")
	r.JS()

	// js is not a directory.
	jsPath := r.Join("js")
	os.Mkdir(r.Path(), os.ModePerm)
	defer os.RemoveAll(Resources().Path())
	f, _ := os.Create(jsPath)
	defer f.Close()

	r.JS()
}

func TestResources(t *testing.T) {
	t.Log(Resources())
}

func TestIsSupportedImageExtension(t *testing.T) {
	if !IsSupportedImageExtension("logo.png") {
		t.Error("logo.png should be a supported image")
	}

	if IsSupportedImageExtension("logo.txt") {
		t.Error("logo.txt should not be a supported image")
	}
}
