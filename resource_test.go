package app

import "testing"
import "os"
import "path/filepath"

func TestResourcePathJoin(t *testing.T) {
	l := ResourcePath("resources")

	if j := l.Join("css"); j != "resources/css" {
		t.Error("j should be resources/css:", j)
	}
}

func TestResourcePathCSS(t *testing.T) {
	r := ResourcePath("resources")
	cssPath := r.Join("css")

	if err := os.MkdirAll(cssPath, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(Resources().Path())

	os.Create(filepath.Join(cssPath, "foo.css"))
	os.Create(filepath.Join(cssPath, "hello"))
	os.Create(filepath.Join(cssPath, "world.txt"))
	os.Mkdir(filepath.Join(cssPath, "bar"), os.ModePerm)

	cssFilenames := r.CSS()

	if l := len(cssFilenames); l != 1 {
		t.Error("cssFilenames should have 1 element:", l)
	}

	if cssFilenames[0] != "foo.css" {
		t.Error("cssFilenames[0] should be foo.css:", cssFilenames[0])
	}
}

func TestResourcePathCSSError(t *testing.T) {
	// No css directory.
	r := ResourcePath("resources")
	r.CSS()

	// css is not a directory.
	cssPath := r.Join("css")
	os.Mkdir(r.Path(), os.ModePerm)
	os.Create(cssPath)
	defer os.RemoveAll(Resources().Path())
	r.CSS()
}

func TestResourcePathJS(t *testing.T) {
	r := ResourcePath("resources")
	jsPath := r.Join("js")

	if err := os.MkdirAll(jsPath, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(Resources().Path())

	os.Create(filepath.Join(jsPath, "foo.js"))
	os.Create(filepath.Join(jsPath, "hello"))
	os.Create(filepath.Join(jsPath, "world.txt"))
	os.Mkdir(filepath.Join(jsPath, "bar"), os.ModePerm)

	jsFilenames := r.JS()

	if l := len(jsFilenames); l != 1 {
		t.Error("jsFilenames should have 1 element:", l)
	}

	if jsFilenames[0] != "foo.js" {
		t.Error("jsFilenames[0] should be foo.js:", jsFilenames[0])
	}
}

func TestResourcePathJSError(t *testing.T) {
	// No js directory.
	r := ResourcePath("resources")
	r.JS()

	// js is not a directory.
	jsPath := r.Join("js")
	os.Mkdir(r.Path(), os.ModePerm)
	os.Create(jsPath)
	defer os.RemoveAll(Resources().Path())
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
