package markup

import (
	"bufio"
	"io"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Tag represents an HTML tag.
type Tag struct {
	ID       uuid.UUID
	CompoID  uuid.UUID
	Name     string
	Text     string
	Svg      bool
	Attrs    AttrMap
	Children []Tag
}

// IsEmpty reports whether its argument t is nil.
// Empty tags have empty name and empty text.
func (t *Tag) IsEmpty() bool {
	return len(t.Name) == 0 && len(t.Text) == 0
}

// IsText reports whether its argument t represents a text.
// Text tags have empty name and non empty text.
func (t *Tag) IsText() bool {
	if t.IsEmpty() {
		return false
	}
	return len(t.Name) == 0 && len(t.Text) != 0
}

// IsComponent reports whether its argument t represents a component.
// Component tags have non standard HTML5 tag name.
func (t *Tag) IsComponent() bool {
	if len(t.Name) == 0 {
		return false
	}

	if t.Svg {
		return false
	}

	a := atom.Lookup([]byte(t.Name))
	return a == 0
}

// IsVoidElem reports whether its argument t represents a void element.
// Void elements are tags listed at
// https://www.w3.org/TR/html5/syntax.html#void-elements.
func (t *Tag) IsVoidElem() bool {
	if t.Svg {
		return false
	}
	_, ok := voidElems[t.Name]
	return ok
}

var (
	voidElems = map[string]struct{}{
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
)

// AttrMap represents a map of attributes.
type AttrMap map[string]string

// AttrEquals reports wheter its arguments l and r are equals.
func AttrEquals(l, r AttrMap) bool {
	if len(l) != len(r) {
		return false
	}

	for k, v := range l {
		otherVal, ok := r[k]
		if !ok {
			return false
		}
		if v != otherVal {
			return false
		}
	}
	return true
}

// TagEncoder is the interface that describes a encoder that convert a Tag to
// HTML5.
type TagEncoder interface {
	Encode(t Tag) error
}

// NewTagEncoder creates a new tag encoder.
func NewTagEncoder(w io.Writer, env Env) TagEncoder {
	return &tagEncoder{
		w:   bufio.NewWriter(w),
		env: env,
	}
}

type tagEncoder struct {
	w   *bufio.Writer
	env Env
	svg bool
}

func (e *tagEncoder) Encode(t Tag) error {
	if err := e.encode(t, 0); err != nil {
		return err
	}
	return e.w.Flush()
}

func (e *tagEncoder) encode(t Tag, indent int) error {
	if t.IsText() {
		e.encodeIndent(indent)
		e.w.WriteString(t.Text)
		return nil
	}

	if t.IsComponent() {
		return e.encodeComponent(t, indent)
	}

	e.encodeIndent(indent)
	e.w.WriteString("<")
	e.w.WriteString(t.Name)
	e.encodeAttributes(t)
	e.w.WriteRune('>')

	if t.IsVoidElem() {
		return nil
	}

	if len(t.Children) == 0 {
		e.w.WriteString("</")
		e.w.WriteString(t.Name)
		e.w.WriteRune('>')
		return nil
	}

	for _, child := range t.Children {
		e.w.WriteRune('\n')
		e.encode(child, indent+1)
	}

	e.w.WriteRune('\n')
	e.encodeIndent(indent)
	e.w.WriteString("</")
	e.w.WriteString(t.Name)
	e.w.WriteRune('>')
	return nil
}

func (e *tagEncoder) encodeComponent(t Tag, indent int) error {
	c, err := e.env.Component(t.ID)
	if err != nil {
		return errors.Wrap(err, "can't encode component")
	}

	root, _ := e.env.Root(c)
	return e.encode(root, indent)
}

func (e *tagEncoder) encodeAttributes(t Tag) {
	for k, v := range t.Attrs {
		if len(v) == 0 {
			e.w.WriteRune(' ')
			e.w.WriteString(k)
			continue
		}

		if strings.HasPrefix(k, "on") {
			e.w.WriteRune(' ')
			e.w.WriteString(k)
			e.w.WriteString(`="CallGoHandler('`)
			e.w.WriteString(t.CompoID.String())
			e.w.WriteString(`', '`)
			e.w.WriteString(v)
			e.w.WriteString(`', this, event)"`)
			continue
		}

		e.w.WriteRune(' ')
		e.w.WriteString(k)
		e.w.WriteString(`="`)
		e.w.WriteString(v)
		e.w.WriteString(`"`)
	}

	e.w.WriteString(` data-go-id="`)
	e.w.WriteString(t.ID.String())
	e.w.WriteString(`"`)
}

func (e *tagEncoder) encodeIndent(indent int) {
	for i := 0; i < indent; i++ {
		e.w.WriteString("  ")
	}
}

// TagDecoder is the interface that describes a decoder that can read HTML5 code
// and translate it to a Tag tree.
// Additionally, HTML5 can embed custom component tags.
type TagDecoder interface {
	Decode(t *Tag) error
}

// NewTagDecoder creates a new tag decoder.
func NewTagDecoder(r io.Reader) TagDecoder {
	return &tagDecoder{
		tokenizer: html.NewTokenizer(r),
	}
}

type tagDecoder struct {
	tokenizer *html.Tokenizer
	svg       bool
	err       error
}

func (d *tagDecoder) Decode(t *Tag) error {
	d.decode(t)

	if t.IsEmpty() {
		return errors.New("can't decode an empty html")
	}

	return d.err
}

func (d *tagDecoder) decode(t *Tag) bool {
	switch d.tokenizer.Next() {
	case html.StartTagToken:
		return d.decodeTag(t)

	case html.EndTagToken:
		return d.decodeEndTag(t)

	case html.TextToken:
		return d.decodeText(t)

	case html.SelfClosingTagToken:
		return d.decodeSelfClosingTag(t)

	case html.ErrorToken:
		return false
	}
	return d.decode(t)
}

func (d *tagDecoder) decodeTag(t *Tag) bool {
	bname, hasAttr := d.tokenizer.TagName()
	name := string(bname)
	t.Name = name

	if name == "svg" {
		d.svg = true
	}
	t.Svg = d.svg

	if hasAttr {
		d.decodeAttrs(t)
	}

	if t.IsComponent() || t.IsVoidElem() {
		return true
	}

	for {
		c := Tag{}
		if !d.decode(&c) {
			return false
		}

		// An empty child can result only if the decoder encountered an end tag.
		// It means the current tag doesn't have other children.
		if c.IsEmpty() {
			return true
		}

		t.Children = append(t.Children, c)
	}
}

func (d *tagDecoder) decodeAttrs(t *Tag) {
	attrs := make(AttrMap)
	for {
		key, val, more := d.tokenizer.TagAttr()
		attrs[string(key)] = string(val)
		if !more {
			break
		}
	}
	t.Attrs = attrs
}

func (d *tagDecoder) decodeEndTag(t *Tag) bool {
	bname, _ := d.tokenizer.TagName()
	name := string(bname)

	if name == "svg" {
		d.svg = false
	}
	return true
}

func (d *tagDecoder) decodeSelfClosingTag(t *Tag) bool {
	bname, hasAttr := d.tokenizer.TagName()
	name := string(bname)

	if !d.svg || name == "svg" {
		d.err = errors.Errorf("%s should not be a self closing tag", name)
		return false
	}

	t.Name = name
	t.Svg = true

	if hasAttr {
		d.decodeAttrs(t)
	}
	return true
}

func (d *tagDecoder) decodeText(t *Tag) bool {
	text := string(d.tokenizer.Text())
	text = strings.TrimSpace(text)

	// There is no need to have empty text tag. If it is the case we try to
	// decode the next tag.
	if len(text) == 0 {
		return d.decode(t)
	}

	t.Text = text
	return true
}
