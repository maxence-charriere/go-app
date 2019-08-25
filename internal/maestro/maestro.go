// +build js

package maestro

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"strings"
	"sync"
	"syscall/js"

	"golang.org/x/net/html"
)

// Maestro is a document object model that manage components and html nodes.
type Maestro struct {
	CompoBuilder   CompoBuilder
	components     map[Compo]*Node
	converters     map[string]interface{}
	attrTransforms []attrTransform
	mutex          sync.Mutex
}

// NewMaestro creates a maestro.
func NewMaestro(compoBuilder CompoBuilder) *Maestro {
	return &Maestro{
		CompoBuilder: compoBuilder,
		components:   make(map[Compo]*Node),
		converters: map[string]interface{}{
			"compo": urlToHTMLTag,
			"json":  jsonFormat,
			"raw":   rawHTML,
			"time":  timeFormat,
		},
		attrTransforms: []attrTransform{eventTransform},
	}
}

// NewBody insert the given component into the document body.
func (m *Maestro) NewBody(c Compo) error {
	if err := m.Render(c); err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	root, ok := m.components[c]
	if !ok {
		return errors.New("root not found")
	}

	body := js.Global().Get("document").Get("body")

	for {
		firstChild := body.Get("firstChild")
		if !firstChild.Truthy() {
			break
		}
		body.Call("removeChild", firstChild)
	}

	body.Call("appendChild", root)
	return nil
}

// Render renders the given component.
func (m *Maestro) Render(c Compo) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

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
		Maestro:   m,
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
		n.newText()
	}

	if n.Name != "" {
		m.dismount(n)
		n.Name = ""
		n.Attrs = nil
		n.change("", "")
	}

	n.compo = ctx.Compo

	if n.Text != text {
		n.Text = text
		n.updateText(text)
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
		n.new(name, ctx.Namespace)
	}

	for _, c := range n.Children {
		m.dismount(c)
		c.removeChild(c)
	}
	n.Children = nil

	if n.Name != name {
		m.dismount(n)
		n.Name = name
		n.Text = ""
		n.Attrs = nil
		n.change(name, ctx.Namespace)
	}

	n.compo = ctx.Compo
	m.renderTagAttrs(ctx, n, hasAttr)
	return nil
}

func (m *Maestro) renderStartTag(ctx renderContext, n *Node, name string, hasAttr bool) error {
	if n.isZero() {
		n.Name = name
		n.Text = ""
		n.new(name, ctx.Namespace)
	}

	if n.Name != name {
		m.dismount(n)
		n.Name = name
		n.Text = ""
		n.Attrs = nil
		n.compo = ctx.Compo
		n.change(name, ctx.Namespace)
	}

	n.compo = ctx.Compo
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
			n.removeChild(c)
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
		n.appendChild(&c)
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
		}

		for _, transform := range m.attrTransforms {
			k, v = transform(k, v)
		}

		attrs[k] = v
	}

	for k, v := range n.Attrs {
		if _, ok := attrs[k]; !ok {
			n.deleteAttr(k)

			if isGoEventAttr(k, v) {
				m.closeEventHandler(ctx, n, k)
			}

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

		if isGoEventAttr(k, v) {
			m.closeEventHandler(ctx, n, k)
			m.setEventHandler(ctx, n, k, v)
		}

		n.upsertAttr(k, v)
		n.Attrs[k] = v
	}
}

func (m *Maestro) setEventHandler(ctx renderContext, n *Node, k, v string) {
	k = strings.TrimPrefix(k, "on")

	if n.eventCloses == nil {
		n.eventCloses = make(map[string]func())
	}

	n.eventCloses[k] = n.addEventListener(
		ctx,
		k,
		strings.TrimPrefix(v, "//go: "),
	)
}

func (m *Maestro) closeEventHandler(ctx renderContext, n *Node, k string) {
	close, ok := n.eventCloses[k]
	if !ok {
		return
	}

	close()
	delete(n.eventCloses, k)
}

func (m *Maestro) renderEndTag(ctx renderContext, n *Node) error {
	n.isEnd = true
	return nil
}

func (m *Maestro) renderCompoNode(ctx renderContext, n *Node, name string, hasAttr bool) error {
	var compo Compo
	var err error

	if n.isZero() {
		compo, err = m.CompoBuilder.New(name)
	} else if name != compoName(n.compo) {
		m.dismount(n)
		compo, err = m.CompoBuilder.New(name)
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

	for k, close := range n.eventCloses {
		close()
		delete(n.eventCloses, k)
	}
	n.Attrs = nil

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
	Maestro   *Maestro
}
