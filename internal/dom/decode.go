package dom

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func decodeNodes(s string, hrefFmt bool) (node, error) {
	d := &decoder{
		tokenizer: html.NewTokenizer(bytes.NewBufferString(s)),
		hrefFmt:   hrefFmt,
	}

	return d.decode()
}

type decoder struct {
	tokenizer   *html.Tokenizer
	decodingSVG bool
	hrefFmt     bool
}

func (d *decoder) decode() (node, error) {
	switch d.tokenizer.Next() {
	case html.TextToken:
		return d.decodeText()

	case html.SelfClosingTagToken:
		return d.decodeSelfClosingElem()

	case html.StartTagToken:
		return d.decodeElem()

	case html.EndTagToken:
		return d.closeElem()

	case html.ErrorToken:
		return nil, d.tokenizer.Err()

	default:
		// Nothing we care about, decode the next node.
		return d.decode()
	}
}

func (d *decoder) decodeText() (node, error) {
	text := string(d.tokenizer.Text())
	text = strings.TrimSpace(text)

	if len(text) == 0 {
		// Text is empty, decode the next node.
		return d.decode()
	}

	t := newText()
	t.SetText(text)
	return t, nil
}

func (d *decoder) decodeSelfClosingElem() (node, error) {
	name, hasAttr := d.tokenizer.TagName()
	tagName := string(name)

	if isCompoTagName(tagName, d.decodingSVG) {
		return newCompo(tagName, d.decodeAttrs(hasAttr)), nil
	}

	e := newElem(tagName)
	e.SetAttrs(d.decodeAttrs(hasAttr))
	return e, nil
}

func (d *decoder) decodeElem() (node, error) {
	name, hasAttr := d.tokenizer.TagName()
	tagName := string(name)

	if isCompoTagName(tagName, d.decodingSVG) {
		return newCompo(tagName, d.decodeAttrs(hasAttr)), nil
	}

	if tagName == "svg" {
		d.decodingSVG = true
	}

	e := newElem(tagName)
	e.SetAttrs(d.decodeAttrs(hasAttr))

	if isVoidElem(tagName) {
		return e, nil
	}

	for {
		child, err := d.decode()
		if err != nil {
			return nil, err
		}
		if child == nil {
			break
		}
		e.appendChild(child)
	}

	return e, nil
}

func (d *decoder) decodeAttrs(hasAttr bool) map[string]string {
	if !hasAttr {
		return nil
	}

	attrs := make(map[string]string)
	for {
		name, val, moreAttr := d.tokenizer.TagAttr()
		n := string(name)
		v := string(val)

		switch {
		case strings.HasPrefix(n, "on") && !strings.HasPrefix(v, "js:"):
			v = fmt.Sprintf(`callCompoHandler(this, event, '%s')`, v)

		case n == "href" && d.hrefFmt:
			if u, err := url.Parse(v); err == nil && len(u.Scheme) == 0 {
				if len(u.Path) != 0 && u.Path[0] != '/' {
					u.Path = "/" + u.Path
				}

				u.Scheme = "compo"
				v = u.String()
			}
		}

		attrs[n] = v

		if !moreAttr {
			break
		}
	}

	return attrs
}

func (d *decoder) closeElem() (node, error) {
	tagName, _ := d.tokenizer.TagName()
	if string(tagName) == "svg" {
		d.decodingSVG = false
	}

	return nil, nil
}

func isHTMLTagName(tagName string) bool {
	return atom.Lookup([]byte(tagName)) != 0
}

func isCompoTagName(tagName string, decodingSVG bool) bool {
	if decodingSVG {
		return false
	}
	return !isHTMLTagName(tagName)
}

var voidElems = map[string]struct{}{
	"area":   {},
	"base":   {},
	"br":     {},
	"col":    {},
	"embed":  {},
	"hr":     {},
	"img":    {},
	"input":  {},
	"keygen": {},
	"link":   {},
	"meta":   {},
	"param":  {},
	"source": {},
	"track":  {},
	"wbr":    {},
}

func isVoidElem(tagName string) bool {
	_, ok := voidElems[tagName]
	return ok
}
