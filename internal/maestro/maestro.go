package maestro

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/template"

	"golang.org/x/net/html"
)

// Maestro is a document object model that manage components and html nodes.
type Maestro struct {
	compoBuilders map[string]reflect.Type
	components    map[Compo]*Node
	converters    map[string]interface{}
}

// NewMaestro creates a maestro.
func NewMaestro() *Maestro {
	return &Maestro{
		compoBuilders: make(map[string]reflect.Type),
		components:    make(map[Compo]*Node),
		converters: map[string]interface{}{
			"bind":  bind,
			"compo": urlToHTMLTag,
			"json":  jsonFormat,
			"raw":   rawHTML,
			"time":  timeFormat,
		},
	}
}

// Import creates a builder that can build the given component.
func (m *Maestro) Import(c Compo) error {
	v := reflect.ValueOf(c)
	if v.Kind() != reflect.Ptr {
		return errors.New("component is not a pointer")
	}
	if v = v.Elem(); v.Kind() != reflect.Struct {
		return errors.New("component is not implemented on a struct")
	}
	if v.NumField() == 0 {
		return errors.New("component does not have fields")
	}

	m.compoBuilders[compoName(c)] = v.Type()
	return nil
}

// New creates the named component.
func (m *Maestro) New(name string) (Compo, error) {
	t, ok := m.compoBuilders[name]
	if !ok {
		return nil, errors.New("component " + name + " is not imported")
	}
	return reflect.New(t).Interface().(Compo), nil
}

// NewBody insert the given component into the document body.
func (m *Maestro) NewBody(c Compo) error {
	if err := m.Render(c); err != nil {
		return err
	}

	root, ok := m.components[c]
	if !ok {
		return errors.New("root not found")
	}

	root.jsNode.addToBody()
	return nil
}

// Render renders the given component.
func (m *Maestro) Render(c Compo) error {
	n, ok := m.components[c]
	if !ok {
		n = &Node{}
	}

	return m.render(c, n)
}

func (m *Maestro) render(c Compo, n *Node) error {
	requireMount := c != n.compo

	rendering, err := m.compoToHTML(c)
	if err != nil {
		return err
	}

	if err := m.renderNode(renderContext{
		Tokenizer: html.NewTokenizer(bytes.NewBufferString(rendering)),
		Compo:     c,
	}, n); err != nil {
		return nil
	}

	if requireMount {
		n.CompoName = compoName(c)
		m.mount(n)
	}

	return nil
}

func (m *Maestro) compoToHTML(c Compo) (string, error) {
	var extendedFuncs map[string]interface{}
	if extended, ok := c.(compoWithExtendedRender); ok {
		extendedFuncs = extended.Funcs()
	}

	// The number of template functions. It contains the
	// component extended functions, the converters and
	// the resources accessor.
	funcsCount := len(m.converters) + len(extendedFuncs) + 1

	funcs := make(template.FuncMap, funcsCount)

	for k, v := range m.converters {
		funcs[k] = v
	}

	for k, v := range extendedFuncs {
		if _, ok := funcs[k]; ok {
			return "", errors.New("template extension can't be named: " + k)
		}
		funcs[k] = v
	}

	tmpl, err := template.
		New(fmt.Sprintf("%T", c)).
		Funcs(funcs).
		Parse(c.Render())
	if err != nil {
		return "", err
	}

	var w bytes.Buffer
	if err = tmpl.Execute(&w, c); err != nil {
		return "", err
	}

	html := strings.TrimSpace(w.String())
	if len(html) == 0 {
		return "", errors.New("component does not render anything")
	}

	return html, nil
}

func (m *Maestro) renderNode(ctx renderContext, n *Node) error {
	switch typ := ctx.Tokenizer.Next(); typ {
	case html.TextToken:
		return m.renderText(ctx, n)

	case html.SelfClosingTagToken, html.StartTagToken:
		return m.renderTag(ctx, n, typ)

	case html.EndTagToken:
		return m.renderEndTag(ctx, n)

	case html.ErrorToken:
		err := ctx.Tokenizer.Err()
		if err == io.EOF {
			return m.renderEndTag(ctx, n)
		}
		return err

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

	if n.isZero() {
		n.compo = ctx.Compo
		n.jsNode.newText()
	}

	if n.Name != "" {
		m.dismount(n)
		n.Name = ""
		n.Attrs = nil
		n.compo = ctx.Compo
		n.jsNode.change("", "")
	}

	if n.Text != text {
		n.Text = text
		n.jsNode.updateText(text)
	}

	return nil
}

func (m *Maestro) renderTag(ctx renderContext, n *Node, typ html.TokenType) error {
	tagName, hasAttr := ctx.Tokenizer.TagName()
	name := string(tagName)

	if ctx.Namespace == "" {
		ctx.Namespace = namespaces[name]
	}

	if isVoidElem(name) {
		return m.renderSelfClosingTag(ctx, n, name, hasAttr)
	}

	if isCompoNode(name, ctx.Namespace) {
		return m.renderCompoNode(ctx, n, name, hasAttr)
	}

	switch typ {
	case html.SelfClosingTagToken:
		return m.renderSelfClosingTag(ctx, n, name, hasAttr)

	default:
		return m.renderStartTag(ctx, n, name, hasAttr)
	}
}

func (m *Maestro) renderSelfClosingTag(ctx renderContext, n *Node, name string, hasAttr bool) error {
	if n.isZero() {
		n.Name = name
		n.compo = ctx.Compo
		n.jsNode.new(name, ctx.Namespace)
	}

	for _, c := range n.Children {
		m.dismount(c)
		c.jsNode.removeChild(c.jsNode)
	}
	n.Children = nil

	if n.Name != name {
		m.dismount(n)
		n.Name = name
		n.Text = ""
		n.Attrs = nil
		n.compo = ctx.Compo
		n.jsNode.change(name, ctx.Namespace)
	}

	m.renderTagAttrs(ctx, n, hasAttr)
	return nil
}

func (m *Maestro) renderStartTag(ctx renderContext, n *Node, name string, hasAttr bool) error {
	if n.isZero() {
		n.Name = name
		n.Text = ""
		n.compo = ctx.Compo
		n.jsNode.new(name, ctx.Namespace)
	}

	if n.Name != name {
		m.dismount(n)
		n.Name = name
		n.Text = ""
		n.Attrs = nil
		n.compo = ctx.Compo
		n.jsNode.change(name, ctx.Namespace)
	}

	m.renderTagAttrs(ctx, n, hasAttr)

	var childrenToDelete []*Node
	for i, c := range n.Children {
		m.renderNode(ctx, c)

		if c.isEnd {
			childrenToDelete = n.Children[i:]
			n.Children = n.Children[:i]
			break
		}
	}

	if childrenToDelete != nil {
		for i, c := range childrenToDelete {
			m.dismount(c)
			n.jsNode.removeChild(c.jsNode)
			childrenToDelete[i] = nil
		}
		return nil
	}

	for {
		var c Node
		m.renderNode(ctx, &c)

		if c.isEnd {
			return nil
		}

		n.Children = append(n.Children, &c)
		n.jsNode.appendChild(c.jsNode)
	}
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

		switch ctx.Namespace {
		case namespaces["svg"]:
			k = svgAttr(k)

		default:
			// TODO: attr transform.
		}

		attrs[k] = v
	}

	for k := range n.Attrs {
		if _, ok := attrs[k]; !ok {
			n.jsNode.deleteAttr(k)
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
		n.jsNode.upsertAttr(k, v)
		n.Attrs[k] = v
	}
}

func (m *Maestro) renderEndTag(ctx renderContext, n *Node) error {
	n.isEnd = true
	return nil
}

func (m *Maestro) renderCompoNode(ctx renderContext, n *Node, name string, hasAttr bool) error {
	var compo Compo
	var err error

	if n.isZero() {
		compo, err = m.New(name)
	} else if name != compoName(compo) {
		m.dismount(n)
		compo, err = m.New(name)
	} else {
		compo = n.compo
	}

	if err != nil {
		return err
	}

	attrs := m.getCompoAttrs(ctx, hasAttr)
	if err = mapCompoFields(compo, attrs); err != nil {
		return err
	}

	return m.render(compo, n)
}

func (m *Maestro) getCompoAttrs(ctx renderContext, hasAttr bool) map[string]string {
	var attrs map[string]string
	if hasAttr {
		attrs = make(map[string]string)
	}

	for hasAttr {
		var k []byte
		var v []byte
		k, v, hasAttr = ctx.Tokenizer.TagAttr()
		attrs[string(k)] = string(v)
	}
	return attrs
}

func (m *Maestro) mount(n *Node) {
	m.components[n.compo] = n

	if m, ok := n.compo.(mounter); ok {
		m.OnMount()
	}
}

func (m *Maestro) dismount(n *Node) {
	for _, c := range n.Children {
		m.dismount(c)
	}

	if !n.isCompoRoot() {
		return
	}

	n.CompoName = ""
	m.components[n.compo] = nil

	if d, ok := n.compo.(dismounter); ok {
		d.OnDismount()
	}
}

type renderContext struct {
	Tokenizer *html.Tokenizer
	Compo     Compo
	Namespace string
}
