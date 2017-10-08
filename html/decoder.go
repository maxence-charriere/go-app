package html

import (
	"io"

	"github.com/murlokswarm/app"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

type Decoder struct {
	tokenizer   *html.Tokenizer
	decodingSvg bool
	err         error
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		tokenizer: html.NewTokenizer(r),
	}
}

func (d *Decoder) Decode(tag *app.Tag) error {
	d.decode(tag)
	if d.err == io.EOF {
		d.err = nil
	}
	if d.err != nil {
		return d.err
	}

	if tag.Is(app.ZeroTag) {
		return errors.Errorf("no html to decode")
	}
	return nil
}

func (d *Decoder) decode(tag *app.Tag) bool {
	switch d.tokenizer.Next() {
	case html.StartTagToken:
		return d.decodeTag(tag)
	case html.EndTagToken:
		return d.decodeClosingTag(tag)
	case html.SelfClosingTagToken:
		return d.decodeSelfClosingTag(tag)
	case html.ErrorToken:
		d.err = d.tokenizer.Err()
		return false
	default:
		return d.decode(tag)
	}
}

func (d *Decoder) decodeTag(tag *app.Tag) bool {
	name, hasAttr := d.tokenizer.TagName()
	tag.Name = string(name)
	tag.Type = app.SimpleTag

	if hasAttr {
		d.decodeAttributes(tag)
	}

	switch {
	case tag.Name == "svg":
		d.decodingSvg = true
		tag.Svg = true
	case isVoidElement(tag.Name, d.decodingSvg):
		return true
	case isComponent(tag.Name, d.decodingSvg):
		tag.Type = app.CompoTag
		return true
	}

	for {
		var child app.Tag
		if !d.decode(&child) {
			return false
		}

		// A zero tag results from decoding a closing tag. It means there is no
		// more child to decode for this tag.
		if child.Is(app.ZeroTag) {
			return true
		}

		tag.Children = append(tag.Children, child)
	}
}

func (d *Decoder) decodeAttributes(tag *app.Tag) {
	attrs := make(app.AttributeMap)
	for {
		name, val, moreAttr := d.tokenizer.TagAttr()
		attrs[string(name)] = string(val)

		if !moreAttr {
			break
		}
	}
	tag.Attributes = attrs
}

func (d *Decoder) decodeClosingTag(tag *app.Tag) bool {
	name, _ := d.tokenizer.TagName()
	if string(name) == "svg" {
		d.decodingSvg = false
	}
	return true
}

func (d *Decoder) decodeSelfClosingTag(tag *app.Tag) bool {
	if !d.decodingSvg {
		d.err = errors.Errorf("decoding a self closing tag is not allowed outside a svg context")
		return false
	}

	name, hasAttr := d.tokenizer.TagName()
	tag.Name = string(name)
	tag.Type = app.SimpleTag
	tag.Svg = true

	if hasAttr {
		d.decodeAttributes(tag)
	}
	return true
}

func (d *Decoder) decodeText(tag *app.Tag) bool {
	text := string(d.tokenizer.Text())
	if len(text) == 0 {
		return d.decode(tag)
	}

	tag.Text = string(text)
	tag.Type = app.TextTag
	return true
}
