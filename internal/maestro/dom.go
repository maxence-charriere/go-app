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

// Dom is a document object model that manage components and html nodes.
type Dom struct {
	CompoBuilder        CompoBuilder
	CallOnUI            func(func())
	TrackCursorPosition func(js.Value)
	ContextMenu         Compo

	once           sync.Once
	components     map[Compo]*Node
	converters     map[string]interface{}
	attrTransforms []attrTransform
}

func (d *Dom) init() {
	d.components = make(map[Compo]*Node)
	d.converters = map[string]interface{}{
		"compo": urlToHTMLTag,
		"json":  jsonFormat,
		"raw":   rawHTML,
		"time":  timeFormat,
	}
	d.attrTransforms = []attrTransform{eventTransform}
}

// NewBody insert the given component into the document body.
func (d *Dom) NewBody(c Compo) error {
	d.once.Do(d.init)

	if err := d.Render(c); err != nil {
		return err
	}
	root, ok := d.components[c]
	if !ok {
		return errors.New("root not found")
	}

	if err := d.Render(d.ContextMenu); err != nil {
		return err
	}
	ctxMenu, ok := d.components[d.ContextMenu]
	if !ok {
		return errors.New("context menu not found")
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
	body.Call("appendChild", ctxMenu)
	return nil
}

// Render renders the given component.
func (d *Dom) Render(c Compo) error {
	n, ok := d.components[c]
	if !ok {
		n = &Node{}
	}

	return d.render(c, n)
}

func (d *Dom) render(c Compo, n *Node) error {
	requireMount := c != n.compo

	rendering, err := d.compoToHTML(c)
	if err != nil {
		return err
	}

	if err := d.renderNode(renderContext{
		Tokenizer: html.NewTokenizer(bytes.NewBufferString(rendering)),
		Compo:     c,
		Dom:       d,
	}, n); err != nil {
		return nil
	}

	if requireMount {
		n.CompoName = CompoName(c)
		d.mount(n)
	}

	return nil
}

func (d *Dom) compoToHTML(c Compo) (string, error) {
	var extendedFuncs map[string]interface{}
	if extended, ok := c.(compoWithExtendedRender); ok {
		extendedFuncs = extended.Funcs()
	}

	// The number of template functions. It contains the
	// component extended functions, the converters and
	// the resources accessor.
	funcsCount := len(d.converters) + len(extendedFuncs) + 1

	funcs := make(template.FuncMap, funcsCount)

	for k, v := range d.converters {
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

func (d *Dom) renderNode(ctx renderContext, n *Node) error {
	switch typ := ctx.Tokenizer.Next(); typ {
	case html.TextToken:
		return d.renderText(ctx, n)

	case html.SelfClosingTagToken, html.StartTagToken:
		return d.renderTag(ctx, n, typ)

	case html.EndTagToken:
		return d.renderEndTag(ctx, n)

	case html.ErrorToken:
		err := ctx.Tokenizer.Err()
		if err == io.EOF {
			return d.renderEndTag(ctx, n)
		}
		return err

	default:
		return d.renderNode(ctx, n)
	}
}

func (d *Dom) renderText(ctx renderContext, n *Node) error {
	text := string(ctx.Tokenizer.Text())
	text = strings.TrimSpace(text)

	// Skip empty text.
	if text == "" {
		return d.renderNode(ctx, n)
	}

	if n.isZero() {
		n.newText()
	}

	if n.Name != "" {
		d.dismount(n)
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

func (d *Dom) renderTag(ctx renderContext, n *Node, typ html.TokenType) error {
	tagName, hasAttr := ctx.Tokenizer.TagName()
	name := string(tagName)

	if ctx.Namespace == "" {
		ctx.Namespace = namespaces[name]
	}

	if isVoidElem(name) {
		return d.renderSelfClosingTag(ctx, n, name, hasAttr)
	}

	if isCompoNode(name, ctx.Namespace) {
		return d.renderCompoNode(ctx, n, name, hasAttr)
	}

	switch typ {
	case html.SelfClosingTagToken:
		return d.renderSelfClosingTag(ctx, n, name, hasAttr)

	default:
		return d.renderStartTag(ctx, n, name, hasAttr)
	}
}

func (d *Dom) renderSelfClosingTag(ctx renderContext, n *Node, name string, hasAttr bool) error {
	if n.isZero() {
		n.Name = name
		n.new(name, ctx.Namespace)
	}

	for _, c := range n.Children {
		d.dismount(c)
		c.removeChild(c)
	}
	n.Children = nil

	if n.Name != name {
		d.dismount(n)
		n.Name = name
		n.Text = ""
		n.Attrs = nil
		n.change(name, ctx.Namespace)
	}

	n.compo = ctx.Compo
	d.renderTagAttrs(ctx, n, hasAttr)
	return nil
}

func (d *Dom) renderStartTag(ctx renderContext, n *Node, name string, hasAttr bool) error {
	if n.isZero() {
		n.Name = name
		n.Text = ""
		n.new(name, ctx.Namespace)
	}

	if n.Name != name {
		d.dismount(n)
		n.Name = name
		n.Text = ""
		n.Attrs = nil
		n.compo = ctx.Compo
		n.change(name, ctx.Namespace)
	}

	n.compo = ctx.Compo
	d.renderTagAttrs(ctx, n, hasAttr)

	var childrenToDelete []*Node
	for i, c := range n.Children {
		d.renderNode(ctx, c)

		if c.isEnd {
			childrenToDelete = n.Children[i:]
			n.Children = n.Children[:i]
			break
		}
	}

	if childrenToDelete != nil {
		for i, c := range childrenToDelete {
			d.dismount(c)
			n.removeChild(c)
			childrenToDelete[i] = nil
		}
		return nil
	}

	for {
		var c Node
		d.renderNode(ctx, &c)

		if c.isEnd {
			return nil
		}

		n.Children = append(n.Children, &c)
		n.appendChild(&c)
	}
}

func (d *Dom) renderTagAttrs(ctx renderContext, n *Node, hasAttr bool) {
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

		for _, transform := range d.attrTransforms {
			k, v = transform(k, v)
		}

		attrs[k] = v
	}

	for k, v := range n.Attrs {
		if _, ok := attrs[k]; !ok {
			n.deleteAttr(k)

			if isGoEventAttr(k, v) {
				d.closeEventHandler(ctx, n, k)
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
			d.closeEventHandler(ctx, n, k)
			d.setEventHandler(ctx, n, k, v)
		}

		n.upsertAttr(k, v)
		n.Attrs[k] = v
	}
}

func (d *Dom) setEventHandler(ctx renderContext, n *Node, k, v string) {
	k = strings.TrimPrefix(k, "on")

	if n.eventCloses == nil {
		n.eventCloses = make(map[string]func())
	}

	v = strings.TrimPrefix(v, "//go: ")
	v = strings.ReplaceAll(v, " ", "")
	n.eventCloses[k] = n.addEventListener(ctx, k, v)
}

func (d *Dom) closeEventHandler(ctx renderContext, n *Node, k string) {
	close, ok := n.eventCloses[k]
	if !ok {
		return
	}

	close()
	delete(n.eventCloses, k)
}

func (d *Dom) renderEndTag(ctx renderContext, n *Node) error {
	n.isEnd = true
	return nil
}

func (d *Dom) renderCompoNode(ctx renderContext, n *Node, name string, hasAttr bool) error {
	var compo Compo
	var err error

	if n.isZero() {
		compo, err = d.CompoBuilder.New(name)
	} else if name != CompoName(n.compo) {
		d.dismount(n)
		compo, err = d.CompoBuilder.New(name)
	} else {
		compo = n.compo
	}
	if err != nil {
		return err
	}

	attrs := d.getCompoAttrs(ctx, hasAttr)

	if err = mapCompoFields(compo, attrs); err != nil {
		return err
	}

	return d.render(compo, n)
}

func (d *Dom) getCompoAttrs(ctx renderContext, hasAttr bool) map[string]string {
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

func (d *Dom) mount(n *Node) {
	d.components[n.compo] = n

	if m, ok := n.compo.(mounter); ok {
		m.OnMount()
	}
}

func (d *Dom) dismount(n *Node) {
	for _, c := range n.Children {
		d.dismount(c)
	}

	for k, close := range n.eventCloses {
		close()
		delete(n.eventCloses, k)
	}

	for _, close := range n.bindingCloses {
		close()
	}
	n.bindingCloses = nil

	n.Attrs = nil

	if !n.isCompoRoot() {
		return
	}

	n.CompoName = ""
	d.components[n.compo] = nil

	if d, ok := n.compo.(dismounter); ok {
		d.OnDismount()
	}
}

func (d *Dom) SetBindingClose(c Compo, close func()) error {
	n, ok := d.components[c]
	if !ok {
		return errors.New("root not found")
	}

	n.bindingCloses = append(n.bindingCloses, close)
	return nil
}

type renderContext struct {
	Tokenizer *html.Tokenizer
	Compo     Compo
	Namespace string
	Dom       *Dom
}
