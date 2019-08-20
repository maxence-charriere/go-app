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
		return m.renderSelfClosingTag(ctx, n)

	case html.StartTagToken:
		return m.renderStartTag(ctx, n)

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
		return n.JSNode().newText(text)
	}

	if n.Type != "text" {
		n.Type = "text"
		n.Text = text
		n.IsCompo = false
		if err := n.JSNode().changeType("text", ""); err != nil {
			return err
		}
		n.JSNode().updateText(text)
		return nil
	}

	if n.Text != text {
		n.Text = text
		n.JSNode().updateText(text)
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
		return n.JSNode().new(typ, "")
	}

	m.removeChildren(ctx, n, 0)

	if n.Type != typ {
		n.Type = typ
		n.Text = ""
		n.IsCompo = false
		n.Attrs = nil

		if err := n.JSNode().changeType(typ, ""); err != nil {
			return err
		}
	}

	m.renderTagAttrs(ctx, n, hasAttr)
	return nil
}

func (m *Maestro) renderStartTag(ctx renderContext, n *Node) error {
	tagName, hasAttr := ctx.Tokenizer.TagName()
	typ := string(tagName)

	switch typ {
	case "svg":
		ctx.Namespace = svgNamespace
	}

	if isCompoNode(typ, ctx.Namespace) {
		return errors.New("component is not a self closing tag: " + typ)
	}

	if n.JSNode() == nil {
		n.Type = typ
		if err := n.JSNode().new(typ, ctx.Namespace); err != nil {
			return err
		}
	}

	if n.Type != typ {
		n.Type = typ
		n.Text = ""
		n.IsCompo = false
		n.Attrs = nil
		m.removeChildren(ctx, n, 0)

		if err := n.JSNode().changeType(typ, ctx.Namespace); err != nil {
			return err
		}
	}

	m.renderTagAttrs(ctx, n, hasAttr)

	for i, c := range n.Children {
		if err := m.renderNode(ctx, c); err != nil {
			return err
		}
		if c.Type == "" {
			m.removeChildren(ctx, n, i)
			return nil
		}
	}

	for {
		var c Node
		if err := m.renderNode(ctx, &c); err != nil {
			return err
		}
		if c.Type == "" {
			return nil
		}

		n.Children = append(n.Children, &c)
		n.JSNode().appendChild(c.JSNode())
	}
}

func (m *Maestro) removeChildren(ctx renderContext, n *Node, start int) {
	children := n.Children[start:]
	for len(children) != 0 {
		c := children[0]
		if c.IsCompo {
			m.dismount(c)
		}

		n.JSNode().removeChild(c.JSNode())
		c.Parent = nil
		children[0] = nil
		children = children[1:]
	}
	n.Children = n.Children[:start]
}

func (m *Maestro) renderTagAttrs(ctx renderContext, n *Node, hasAttr bool) {
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
			n.JSNode().deleteAttr(k)
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
		n.JSNode().upsertAttr(k, v)
		n.Attrs[k] = v
	}
}

func (m *Maestro) renderCompoNode(ctx renderContext, n *Node) error {
	panic("not implemented")
}

func (m *Maestro) dismount(n *Node) {
	panic("not implemented")
}

type renderContext struct {
	Tokenizer  *html.Tokenizer
	Compo      Compo
	ParentNode *Node
	Namespace  string
}
