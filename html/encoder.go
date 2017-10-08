package html

import (
	"bufio"
	"io"
	"strings"

	"github.com/murlokswarm/app"
	"github.com/pkg/errors"
)

type Encoder struct {
	writer *bufio.Writer
	markup app.Markup
}

func NewTagEncoder(w io.Writer, markup app.Markup) *Encoder {
	return &Encoder{
		writer: bufio.NewWriter(w),
		markup: markup,
	}
}

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
		e.writer.WriteByte(' ')
		e.writer.WriteString(name)

		if len(val) == 0 {
			continue
		}

		e.writer.WriteString(`="`)

		if strings.HasPrefix(name, "on") && !strings.HasPrefix(val, "js:") {
			e.writer.WriteString(`CallGoHandler('`)
			e.writer.WriteString(tag.CompoID.String())
			e.writer.WriteString(val)
			e.writer.WriteString(`', this, event)"`)
			continue
		}

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
		return errors.Wrap(err, "encoding component failed")
	}

	root, _ := e.markup.Root(compo)
	return e.encode(root, indent)
}

func (e *Encoder) encodeIndent(indent int) {
	for i := 0; i < indent; i++ {
		e.writer.WriteString("  ")
	}
}
