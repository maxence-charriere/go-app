package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/murlokswarm/app/internal/file"
)

func main() {
	release := filepath.Join(
		os.Getenv("GOPATH"),
		"src",
		"github.com",
		"murlokswarm",
		"uwp",
		"uwp",
		"bin",
		"x64",
		"Release",
	)

	files := []string{
		"App.xbf",
		"WindowPage.xbf",
		filepath.Join("AppX", "clrcompression.dll"),
		filepath.Join("AppX", "uwp.dll"),
		filepath.Join("AppX", "uwp.exe"),
	}

	for _, f := range files {
		src := filepath.Join(release, f)
		dst := filepath.Join("uwp", filepath.Base(f))

		if err := file.Copy(dst, src); err != nil {
			fmt.Println(err)
		}
	}

	goappDLL := filepath.Join(
		os.Getenv("GOPATH"),
		"src",
		"github.com",
		"murlokswarm",
		"uwp",
		"x64",
		"Release",
		"goapp.dll",
	)

	if err := file.Copy(
		filepath.Join("uwp", "goapp.dll"),
		filepath.Join(goappDLL),
	); err != nil {
		fmt.Println(err)
	}

}
