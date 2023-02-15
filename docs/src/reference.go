package main

import (
	"bytes"
	"net/url"
	"strings"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/maxence-charriere/go-app/v9/pkg/errors"
	"golang.org/x/net/html"
)

const (
	getReference   = "/reference/get"
	referenceState = "/reference"
)

func handleGetReference(ctx app.Context, a app.Action) {
	state := referenceState

	var ref htmlContent
	ctx.GetState(state, &ref)
	switch ref.Status {
	case loading, loaded:
		return
	}

	ref.Status = loading
	ref.Err = nil
	ctx.SetState(state, ref)

	res, err := get(ctx, "/web/documents/reference.html")
	if err != nil {
		ref.Status = loadingErr
		ref.Err = errors.New("getting reference failed").Wrap(err)
		ctx.SetState(state, ref)
		return
	}

	doc, err := html.Parse(bytes.NewReader(res))
	if err != nil {
		ref.Status = loadingErr
		ref.Err = errors.New("parsing reference failed").Wrap(err)
		ctx.SetState(state, ref)
		return
	}

	content, err := getHTML(doc, "#page")
	if err != nil {
		ref.Status = loadingErr
		ref.Err = errors.New("generating reference content failed").Wrap(err)
		ctx.SetState(state, ref)
		return
	}
	content = strings.ReplaceAll(content, "▾", "")
	content = strings.ReplaceAll(content, `title="Click to hide Overview section"`, "")
	content = strings.ReplaceAll(content, `title="Click to hide Index section"`, "")
	content = strings.ReplaceAll(content, "/src/github.com/maxence-charriere/go-app/v9/", "https://github.com/maxence-charriere/go-app/blob/master/")

	index, err := getHTML(doc, "#manual-nav")
	if err != nil {
		ref.Status = loadingErr
		ref.Err = errors.New("generating reference content failed").Wrap(err)
		ctx.SetState(state, ref)
		return
	}

	ref.Content = content
	ref.Index = index
	ref.Status = loaded
	ctx.SetState(state, ref)
}

type htmlContent struct {
	Status  status
	Err     error
	Index   string
	Content string
}

func getHTML(n *html.Node, class string) (string, error) {
	section, err := findHTMLNode(n, class)
	if err != nil {
		return "", errors.New("finding html node failed").
			WithTag("target", class).
			Wrap(err)
	}

	normalizeHTMLNode(section)

	var b bytes.Buffer
	if err := html.Render(&b, section); err != nil {
		return "", errors.New("rendering html failed").
			WithTag("target", class).
			Wrap(err)
	}
	return b.String(), nil
}

func findHTMLNode(n *html.Node, sel string) (*html.Node, error) {
	if n.Type == html.ElementNode {
		for _, a := range n.Attr {
			switch {
			case a.Key == "class" && a.Val == strings.TrimPrefix(sel, "."):
				return n, nil

			case a.Key == "id" && a.Val == strings.TrimPrefix(sel, "#"):
				return n, nil
			}

		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if child, err := findHTMLNode(c, sel); err == nil {
			return child, nil
		}
	}

	return nil, errors.New("node not found")
}

func normalizeHTMLNode(n *html.Node) {
	if n.Type == html.ElementNode {
		var id string

		for i, a := range n.Attr {
			if a.Key != "href" {
				continue
			}

			u, err := url.Parse(a.Val)
			if err != nil {
				continue
			}

			switch {

			case strings.HasPrefix(u.Path, "/builtin"):
				u.Scheme = "https"
				u.Host = "pkg.go.dev"

			case u.Scheme == "" && u.Fragment != "":
				id = refLinkID(u.Fragment)
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
		normalizeHTMLNode(c)
	}
}

func refLinkID(v string) string {
	return "ref-link-" + v
}
