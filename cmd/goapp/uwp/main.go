package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/murlokswarm/app/internal/file"
)

func main() {
	appx := filepath.Join(
		os.Getenv("GOPATH"),
		"src",
		"github.com",
		"murlokswarm",
		"uwp",
		"uwp",
		"bin",
		"x64",
		"Release",
		"AppX",
	)

	files := []string{
		"clrcompression.dll",
		"uwp.dll",
		"uwp.exe",
	}

	for _, f := range files {
		src := filepath.Join(appx, f)
		dst := filepath.Join("uwp", f)

		if err := file.Copy(dst, src); err != nil {
			fmt.Println(err)
		}
	}
}
