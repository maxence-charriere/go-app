package maestro

import (
	"bytes"
	"errors"
	"strings"

	"golang.org/x/net/html"
)

type Maestro struct {
	components map[Compo]*Node
	root       *Node
}

func NewMaestro() *Maestro {
	return &Maestro{
		components: make(map[Compo]*Node),
	}
}

func (m *Maestro) Render(c Compo) error {
	// root, ok := m.components[c]
	panic("not implemented")
}

func (m *Maestro) render(c Compo) error {
	root, ok := m.components[c]
	if !ok {
		return errors.New("component not found")
	}

	return m.renderNode(renderContext{
		Tokenizer: html.NewTokenizer(bytes.NewBufferString(c.Render())),
		Compo:     c,
	}, root)

}

func (m *Maestro) renderNode(ctx renderContext, n *Node) error {
	switch ctx.Tokenizer.Next() {
	case html.TextToken:
		return m.renderText(ctx, n)

	case html.SelfClosingTagToken:
		panic("not implemented")

	case html.StartTagToken:
		panic("not implemented")

	case html.EndTagToken:
		panic("not implemented")

	case html.ErrorToken:
		return ctx.Tokenizer.Err()

	default:
		return m.renderNode(ctx, n)
	}
}

func (m *Maestro) renderText(ctx renderContext, n *Node) error {
	text := string(ctx.Tokenizer.Text())
	text = strings.TrimSpace(text)

	// Skip empty text.
	if text == "" {
		return m.renderNode(ctx, n)
	}

	if n.JSNode == nil {
		n.Type = "text"
		n.Text = text
		return n.newText(text)
	}

	if n.Type != "text" {
		n.Type = "text"
		n.Text = text
		n.IsCompo = false
		if err := n.changeType("text", ""); err != nil {
			return err
		}
		n.updateText(text)
		return nil
	}

	if n.Text != text {
		n.Text = text
		n.updateText(text)
	}
	return nil
}

func (m *Maestro) renderSelfClosingTag(ctx renderContext, n *Node) error {
	tagName, hasAttr := ctx.Tokenizer.TagName()
	typ := string(tagName)

	if isCompoNode(typ, "") {
		return m.renderCompoNode(ctx, n)
	}

	if n.JSNode == nil {
		n.Type = typ
		return n.new(typ, "")
	}

	for _, c := range n.Children {
		var err error
		if c.IsCompo {
			err = m.dismount(c.Compo)
		} else {
			c.delete()
		}
		if err != nil {
			return err
		}
	}
	n.Children = nil

	if n.Type != typ {
		n.Type = typ
		n.Text = ""
		n.IsCompo = false
		n.Attrs = nil
		if err := n.changeType(typ, ""); err != nil {
			return err
		}
	}

	return m.renderTagAttrs(ctx, n, hasAttr)
}

func (m *Maestro) renderTagAttrs(ctx renderContext, n *Node, hasAttr bool) error {
	var attrs map[string]string
	if hasAttr {
		attrs = make(map[string]string)
	}

	for hasAttr {
		var tmpK []byte
		var tmpV []byte

		tmpK, tmpV, hasAttr = ctx.Tokenizer.TagAttr()
		k := string(tmpK)
		v := string(tmpV)

		if ctx.Namespace != "" {
			switch k {
			case svgNamespace:
				k = svgAttr(k)
			}
		}

		// TODO: attr transform.
		attrs[k] = v
	}

	for k := range n.Attrs {
		if _, ok := attrs[k]; !ok {
			n.deleteAttr(k)
			delete(n.Attrs, k)
		}
	}

	if n.Attrs == nil {
		n.Attrs = make(map[string]string, len(attrs))
	}

	for k, v := range attrs {
		if oldv, ok := n.Attrs[k]; ok && oldv == v {
			continue
		}
		n.upsertAttr(k, v)
		n.Attrs[k] = v
	}

	return nil
}

func (m *Maestro) renderCompoNode(ctx renderContext, n *Node) error {
	panic("not implemented")
}

func (m *Maestro) dismount(c Compo) error {
	panic("not implemented")
}

type renderContext struct {
	Tokenizer  *html.Tokenizer
	Compo      Compo
	ParentNode *Node
	Namespace  string
}
