package main

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/maxence-charriere/go-app/v7/pkg/app"
	"github.com/maxence-charriere/go-app/v7/pkg/errors"
	"golang.org/x/net/html"
)

func main() {
	cmd := exec.Command("godoc",
		"-url", "/pkg/github.com/maxence-charriere/go-app/v7/pkg/app")
	out, err := cmd.Output()
	if err != nil {
		app.Log("%s", errors.New("reading godoc failed").Wrap(err))
		return
	}

	document, err := html.Parse(bytes.NewReader(out))
	if err != nil {
		app.Log("%s", errors.New("parsing html failed").Wrap(err))
		return
	}

	docs := []struct {
		file     string
		targetID string
	}{
		{
			file:     "../web/godoc-index.html",
			targetID: "manual-nav",
		},
		{
			file:     "../web/godoc.html",
			targetID: "page",
		},
	}

	for _, d := range docs {
		n, err := findHTMLNode(document, d.targetID)
		if err != nil {
			app.Log("%s", errors.New("finding html node failed").
				Tag("id", d.targetID).
				Wrap(err))
			return
		}

		f, err := os.Create(d.file)
		if err != nil {
			app.Log("%s", errors.New("creating godoc file failed").
				Tag("file", d.file).
				Wrap(err))
			return
		}
		defer f.Close()

		if err = html.Render(f, n); err != nil {
			app.Log("%s", errors.New("writting godoc file failed").
				Tag("file", d.file).
				Wrap(err))
			return
		}
	}
}

func findHTMLNode(n *html.Node, id string) (*html.Node, error) {
	if n.Type == html.ElementNode {
		for _, a := range n.Attr {
			if a.Key == "id" && a.Val == id {
				return n, nil
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if child, err := findHTMLNode(c, id); err == nil {
			return child, nil
		}
	}

	return nil, errors.New("not found")
}
