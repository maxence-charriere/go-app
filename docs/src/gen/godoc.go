package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

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

	f, err := os.Create("godoc.go")
	if err != nil {
		app.Log("%s", errors.New("creating godoc.go failed").Wrap(err))
		return
	}
	defer f.Close()

	fmt.Fprintln(f, "package main")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "// Code generated. DO NOT EDIT")
	fmt.Fprintln(f)
	fmt.Fprintln(f, `import "github.com/maxence-charriere/go-app/v7/pkg/app"`)

	if err = generateFunction(htmlContext{
		Context:   context.Background(),
		File:      f,
		HTML:      html.NewTokenizer(bytes.NewBuffer(out)),
		RootID:    "manual-nav",
		RootClass: "godoc-menu",
	}, "GodocMenu"); err != nil {
		app.Log("%s", errors.New("creating godoc.go failed").Wrap(err))
		return
	}

	if err = generateFunction(htmlContext{
		Context:   context.Background(),
		File:      f,
		HTML:      html.NewTokenizer(bytes.NewBuffer(out)),
		RootID:    "page",
		RootClass: "godoc",
	}, "Godoc"); err != nil {
		app.Log("%s", errors.New("creating godoc.go failed").Wrap(err))
		return
	}
}

type htmlContext struct {
	context.Context
	*os.File

	HTML      *html.Tokenizer
	RootID    string
	RootClass string
	Depth     int
}

func generateFunction(ctx htmlContext, name string) error {
	for {
		switch t := ctx.HTML.Next(); t {
		case html.ErrorToken:
			return ctx.HTML.Err()

		case html.StartTagToken:
			tag, hasAttrs := ctx.HTML.TagName()
			for hasAttrs {
				var k []byte
				var v []byte

				k, v, hasAttrs = ctx.HTML.TagAttr()
				key := string(k)
				value := string(v)

				if key == "id" && value == ctx.RootID {
					ctx.Depth = 1

					fmt.Fprintf(ctx, `
					// %s returns a go-app representation of the named documentation.
					func %s() app.UI {
						return app.%s().
							Class(%q).
							ID(%q).
							Body(
					`,
						name,
						name,
						tagToGoappElem(string(tag)),
						ctx.RootClass,
						ctx.RootID,
					)

					err := generate(ctx)
					fmt.Fprintln(ctx, ")\n}")
					return err
				}
			}
		}
	}
}

func generate(ctx htmlContext) error {
	if ctx.Depth == 0 {
		return nil
	}

	switch t := ctx.HTML.Next(); t {
	case html.ErrorToken:
		return ctx.HTML.Err()

	case html.StartTagToken:
		return generateTag(ctx)

	case html.EndTagToken:
		return generateTagEnd(ctx)

	case html.SelfClosingTagToken:
		return generate(ctx)

	case html.TextToken:
		return generateText(ctx)

	default:
		return generate(ctx)
	}
}

func generateTag(ctx htmlContext) error {
	t, hasAttrs := ctx.HTML.TagName()
	tag := string(t)

	fmt.Fprintf(ctx, "app.%s()", tagToGoappElem(tag))

	if hasAttrs {
		fmt.Fprintln(ctx, ".")
		generateAttrs(ctx)
	}

	if isVoidElement(tag) {
		fmt.Fprintln(ctx, ",")
		return generate(ctx)
	}

	ctx.Depth++

	if hasAttrs {
		fmt.Fprintln(ctx, ".")
	} else {
		fmt.Fprint(ctx, ".")
	}

	fmt.Fprintln(ctx, "Body(")
	return generate(ctx)
}

func generateAttrs(ctx htmlContext) {
	k, v, moreAttrs := ctx.HTML.TagAttr()
	key := string(k)
	value := string(v)

	if strings.HasPrefix(key, "on") {
		if moreAttrs {
			generateAttrs(ctx)
		}
		return
	}

	switch key {
	case "style":
		styles := styleToMap(value)
		i := 0
		for k, v := range styles {
			fmt.Fprintf(ctx, "Style(%q, %q)", k, v)

			if i < len(styles)-1 {
				fmt.Fprint(ctx, ".")
			}
			i++
		}

	default:
		fmt.Fprintf(ctx, "%s(%q)", attrToGoappMethod(key), value)
	}

	if moreAttrs {
		fmt.Fprintln(ctx, ".")
		generateAttrs(ctx)
	}
}

func generateTagEnd(ctx htmlContext) error {
	ctx.Depth--

	if ctx.Depth != 0 {
		fmt.Fprintln(ctx, "),")
	}

	return generate(ctx)
}

func generateText(ctx htmlContext) error {
	text := string(ctx.HTML.Text())
	text = strings.TrimSpace(text)

	if text != "" {
		fmt.Fprintf(ctx, "app.Text(%q),\n", text)
	}

	return generate(ctx)
}

func tagToGoappElem(tag string) string {
	return strings.Title(tag)
}

func attrToGoappMethod(attr string) string {
	switch attr {
	case "id":
		return "ID"

	default:
		return strings.Title(attr)
	}
}

func isVoidElement(tag string) bool {
	switch tag {
	case "area",
		"base",
		"br",
		"col",
		"command",
		"embed",
		"hr",
		"img",
		"input",
		"keygen",
		"link",
		"meta",
		"param",
		"source",
		"track",
		"wbr":
		return true

	default:
		return false
	}
}

func styleToMap(s string) map[string]string {
	styles := strings.Split(s, ";")
	m := make(map[string]string)

	for _, s := range styles {
		s := strings.TrimSpace(s)
		if s == "" {
			continue
		}

		values := strings.Split(s, ":")
		if len(values) < 2 {
			continue
		}

		m[values[0]] = values[1]
	}

	return m
}
