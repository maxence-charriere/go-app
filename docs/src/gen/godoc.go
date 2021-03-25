package main

import (
	"bytes"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
	"github.com/maxence-charriere/go-app/v8/pkg/errors"
	"golang.org/x/net/html"
)

func main() {
	cmd := exec.Command("godoc",
		"-url", "/pkg/github.com/maxence-charriere/go-app/v8/pkg/app")
	out, err := cmd.Output()
	if err != nil {
		app.Log(errors.New("reading godoc failed").Wrap(err))
		return
	}

	document, err := html.Parse(bytes.NewReader(out))
	if err != nil {
		app.Log(errors.New("parsing html failed").Wrap(err))
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
			app.Log(errors.New("finding html node failed").
				Tag("id", d.targetID).
				Wrap(err))
			return
		}
		normalizeNode(n)

		f, err := os.Create(d.file)
		if err != nil {
			app.Log(errors.New("creating godoc file failed").
				Tag("file", d.file).
				Wrap(err))
			return
		}
		defer f.Close()

		if err = html.Render(f, n); err != nil {
			app.Log(errors.New("writting godoc file failed").
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

func normalizeNode(n *html.Node) {
	if n.Type == html.ElementNode {
		id := ""

		for i, a := range n.Attr {
			if a.Key != "href" {
				continue
			}

			u, err := url.Parse(a.Val)
			if err != nil {
				continue
			}

			switch {
			case strings.HasPrefix(u.Path, "/src/github.com/maxence-charriere/go-app/v8"):
				u.RawQuery = ""
				u.Path = strings.TrimPrefix(u.Path, "/src/github.com/maxence-charriere/go-app/v8")
				u.Path = "/maxence-charriere/go-app/blob/master" + u.Path
				u.Scheme = "https"
				u.Host = "github.com"

			case strings.HasPrefix(u.Path, "/pkg/"):
				u.Scheme = "https"
				u.Host = "golang.org"

			case u.Scheme == "" && u.Fragment != "":
				id = linkID(u.Fragment)
			}

			a.Val = u.String()
			n.Attr[i] = a
			break
		}

		if id != "" {
			n.Attr = append(n.Attr, html.Attribute{
				Key: "id",
				Val: id,
			})
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		normalizeNode(c)
	}
}

func linkID(fragment string) string {
	return "src-" + fragment
}
