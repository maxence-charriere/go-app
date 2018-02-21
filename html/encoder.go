package html

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/pkg/errors"
)

// Encoder is a tag encoder based on HTML5.
// It implements the app.TagEncoder interface.
type Encoder struct {
	writer *bufio.Writer
	markup app.Markup
}

// NewEncoder create a tag encoder that writes on the given writer.
func NewEncoder(w io.Writer, markup app.Markup) *Encoder {
	return &Encoder{
		writer: bufio.NewWriter(w),
		markup: markup,
	}
}

// Encode encodes the given tag to HTML5.
// It satisfies the app.TagEncoder interface.
func (e *Encoder) Encode(tag app.Tag) error {
	if err := e.encode(tag, 0); err != nil {
		return err
	}
	return e.writer.Flush()
}

func (e *Encoder) encode(tag app.Tag, indent int) error {
	switch tag.Type {
	case app.SimpleTag:
		return e.encodeSimple(tag, indent)

	case app.TextTag:
		return e.encodeText(tag, indent)

	case app.CompoTag:
		return e.encodeComponent(tag, indent)

	default:
		return errors.Errorf("encoding tag %s of type %v is not supported", tag.Name, tag.Type)
	}
}

func (e *Encoder) encodeSimple(tag app.Tag, indent int) error {
	e.encodeIndent(indent)
	e.writer.WriteByte('<')
	e.writer.WriteString(tag.Name)
	e.encodeAttributes(tag)
	e.writer.WriteByte('>')

	if isVoidElement(tag.Name, tag.Svg) {
		return nil
	}

	if len(tag.Children) == 0 {
		e.writer.WriteString("</")
		e.writer.WriteString(tag.Name)
		e.writer.WriteByte('>')
		return nil
	}

	for _, child := range tag.Children {
		e.writer.WriteByte('\n')
		if err := e.encode(child, indent+1); err != nil {
			return err
		}
	}

	e.writer.WriteByte('\n')
	e.encodeIndent(indent)
	e.writer.WriteString("</")
	e.writer.WriteString(tag.Name)
	e.writer.WriteByte('>')
	return nil
}

func (e *Encoder) encodeAttributes(tag app.Tag) {
	for name, val := range tag.Attributes {
		val = AppJSAttributeValue(name, val, e.markup.Factory(), tag.CompoID)

		e.writer.WriteByte(' ')
		e.writer.WriteString(name)

		if len(val) == 0 {
			continue
		}

		e.writer.WriteString(`="`)
		e.writer.WriteString(val)
		e.writer.WriteByte('"')
	}

	e.writer.WriteString(` data-goapp-id="`)
	e.writer.WriteString(tag.ID.String())
	e.writer.WriteByte('"')
}

func (e *Encoder) encodeText(tag app.Tag, indent int) error {
	e.encodeIndent(indent)
	e.writer.WriteString(tag.Text)
	return nil
}

func (e *Encoder) encodeComponent(tag app.Tag, indent int) error {
	compo, err := e.markup.Component(tag.ID)
	if err != nil {
		return err
	}

	root, _ := e.markup.Root(compo)
	return e.encode(root, indent)
}

func (e *Encoder) encodeIndent(indent int) {
	for i := 0; i < indent; i++ {
		e.writer.WriteString("  ")
	}
}

// AppJSAttributeValue is an helper method that make the given attribute value
// compatible with appjs.
func AppJSAttributeValue(name, val string, f app.Factory, compoID uuid.UUID) string {
	if name == "href" {
		u, _ := url.Parse(val)
		compoName := app.ComponentNameFromURL(u)

		if f.Registered(compoName) {
			u.Scheme = "compo"
			u.Path = "/" + compoName
		}
		return u.String()
	}

	if !strings.HasPrefix(name, "on") {
		return val
	}

	if strings.HasPrefix(val, "js:") {
		return strings.TrimPrefix(val, "js:")
	}

	return fmt.Sprintf(`callGoEventHandler('%s', '%s', this, event)`,
		compoID,
		val,
	)
}
