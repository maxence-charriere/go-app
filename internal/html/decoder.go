package html

import (
	"bytes"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func decodeNodes(s string, rootID string) (node, error) {
	d := &decoder{
		tokenizer: html.NewTokenizer(bytes.NewBufferString(s)),
	}
	return d.decode(rootID)
}

type decoder struct {
	tokenizer   *html.Tokenizer
	decodingSVG bool
}

func (d *decoder) decode(id string) (node, error) {
	switch d.tokenizer.Next() {
	case html.TextToken:
		return d.decodeText(id)

	case html.SelfClosingTagToken:
		return d.decodeSelfClosingElem(id)

	case html.StartTagToken:
		return d.decodeElem(id)

	case html.EndTagToken:
		return d.closeElem()

	case html.ErrorToken:
		return nil, d.tokenizer.Err()

	default:
		// Nothing we care about, decode the next node.
		return d.decode(id)
	}
}

func (d *decoder) decodeText(id string) (node, error) {
	text := string(d.tokenizer.Text())
	text = strings.TrimSpace(text)

	if len(text) == 0 {
		// Text is empty, decode the next node.
		return d.decode(id)
	}

	t := newTextNode(id)
	t.SetText(text)
	return t, nil
}

func (d *decoder) decodeSelfClosingElem(id string) (node, error) {
	name, hasAttr := d.tokenizer.TagName()
	tagName := string(name)

	if isCompoTagName(tagName, d.decodingSVG) {
		return newCompoNode(id, tagName, d.decodeAttrs(hasAttr)), nil
	}

	elem := newElemNode(id, tagName)
	elem.SetAttrs(d.decodeAttrs(hasAttr))
	return elem, nil
}

func (d *decoder) decodeAttrs(hasAttr bool) map[string]string {
	if !hasAttr {
		return nil
	}

	attrs := make(map[string]string)
	for {
		name, val, moreAttr := d.tokenizer.TagAttr()
		attrs[string(name)] = string(val)

		if !moreAttr {
			break
		}
	}
	return attrs
}

func (d *decoder) decodeElem(id string) (node, error) {
	name, hasAttr := d.tokenizer.TagName()
	tagName := string(name)

	if isCompoTagName(tagName, d.decodingSVG) {
		return newCompoNode(id, tagName, d.decodeAttrs(hasAttr)), nil
	}

	if tagName == "svg" {
		d.decodingSVG = true
	}

	elem := newElemNode(id, tagName)
	elem.SetAttrs(d.decodeAttrs(hasAttr))

	if isVoidElem(tagName) {
		return elem, nil
	}

	for {
		child, err := d.decode(uuid.New().String())
		if err != nil {
			return nil, err
		}
		if child == nil {
			break
		}
		elem.appendChild(child)
	}
	return elem, nil
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
