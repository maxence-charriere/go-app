package main

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

func main() {
	docs := []struct {
		url      string
		file     string
		targetID string
	}{
		{
			url:      "https://pkg.go.dev/github.com/maxence-charriere/go-app/v9/pkg/app",
			file:     "../web/documents/reference.html",
			targetID: "page",
		},
	}

	for _, d := range docs {
		f, err := os.Create(d.file)
		if err != nil {
			app.Log(errors.New("creating godoc file failed").
				Tag("file", d.file).
				Wrap(err))
			return
		}
		defer f.Close()

		res, err := http.Get(d.url)
		if err != nil {
			app.Log(errors.New("getting document failed").
				Tag("url", d.url).
				Wrap(err))
			return
		}
		defer res.Body.Close()

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			app.Log(errors.New("reading document failed").
				Tag("url", d.url).
				Wrap(err))
			return
		}

		f.Write(b)
	}
}
