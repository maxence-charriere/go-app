package app

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

type dom struct {
	compoBuilder        compoBuilder
	msgs                *messenger
	callOnUI            func(func())
	trackCursorPosition func(js.Value)
	contextMenu         Compo

	once           sync.Once
	components     map[Compo]*node
	converters     map[string]interface{}
	root           *node
	ctxMenuRoot    *node
	attrTransforms []attrTransform
}

func (d *dom) init() {
	d.components = make(map[Compo]*node)
	d.converters = map[string]interface{}{
		"compo": urlToHTMLTag,
		"json":  jsonFormat,
		"raw":   rawHTML,
		"time":  timeFormat,
		"emit":  emitter,
	}
	d.attrTransforms = []attrTransform{eventTransform}
}

func (d *dom) newBody(c Compo) error {
	d.once.Do(d.init)

	if err := d.render(c); err != nil {
		return err
	}
	root, ok := d.components[c]
	if !ok {
		return errors.New("root not found")
	}
	d.root = root

	if err := d.render(d.contextMenu); err != nil {
		return err
	}
	ctxMenu, ok := d.components[d.contextMenu]
	if !ok {
		return errors.New("context menu not found")
	}
	d.ctxMenuRoot = ctxMenu

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

func (d *dom) render(c Compo) error {
	n, ok := d.components[c]
	if !ok {
		n = &node{}
	}
	return d.renderCompo(c, n)
}

func (d *dom) renderCompo(c Compo, n *node) error {
	requireMount := c != n.compo

	rendering, err := d.compoToHTML(c)
	if err != nil {
		return err
	}

	if err := d.renderNode(renderContext{
		tokenizer: html.NewTokenizer(bytes.NewBufferString(rendering)),
		compo:     c,
		dom:       d,
	}, n); err != nil {
		return err
	}

	if requireMount {
		n.CompoName = compoName(c)
		d.mount(n)
	}

	return nil
}

func (d *dom) compoToHTML(c Compo) (string, error) {
	var extendedFuncs map[string]interface{}
	if extended, ok := c.(CompoWithExtendedRender); ok {
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

func (d *dom) renderNode(ctx renderContext, n *node) error {
	switch typ := ctx.tokenizer.Next(); typ {
	case html.TextToken:
		return d.renderText(ctx, n)

	case html.SelfClosingTagToken, html.StartTagToken:
		return d.renderTag(ctx, n, typ)

	case html.EndTagToken:
		return d.renderEndTag(ctx, n)

	case html.ErrorToken:
		err := ctx.tokenizer.Err()
		if err == io.EOF {
			return d.renderEndTag(ctx, n)
		}
		return err

	default:
		return d.renderNode(ctx, n)
	}
}

func (d *dom) renderText(ctx renderContext, n *node) error {
	text := string(ctx.tokenizer.Text())
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

	n.compo = ctx.compo

	if n.Text != text {
		n.Text = text
		n.updateText(text)
	}

	return nil
}

func (d *dom) renderTag(ctx renderContext, n *node, typ html.TokenType) error {
	tagName, hasAttr := ctx.tokenizer.TagName()
	name := string(tagName)

	if ctx.namespace == "" {
		ctx.namespace = namespaces[name]
	}

	if isVoidElem(name) {
		return d.renderSelfClosingTag(ctx, n, name, hasAttr)
	}

	if isCompoNode(name, ctx.namespace) {
		return d.renderCompoNode(ctx, n, name, hasAttr)
	}

	switch typ {
	case html.SelfClosingTagToken:
		return d.renderSelfClosingTag(ctx, n, name, hasAttr)

	default:
		return d.renderStartTag(ctx, n, name, hasAttr)
	}
}

func (d *dom) renderSelfClosingTag(ctx renderContext, n *node, name string, hasAttr bool) error {
	if n.isZero() {
		n.Name = name
		n.new(name, ctx.namespace)
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
		n.change(name, ctx.namespace)
	}

	n.compo = ctx.compo
	d.renderTagAttrs(ctx, n, hasAttr)
	return nil
}

func (d *dom) renderStartTag(ctx renderContext, n *node, name string, hasAttr bool) error {
	if n.isZero() {
		n.Name = name
		n.Text = ""
		n.new(name, ctx.namespace)
	}

	if n.Name != name {
		d.dismount(n)
		n.Name = name
		n.Text = ""
		n.Attrs = nil
		n.compo = ctx.compo
		n.change(name, ctx.namespace)
	}

	n.compo = ctx.compo
	d.renderTagAttrs(ctx, n, hasAttr)

	var childrenToDelete []*node
	for i, c := range n.Children {
		if err := d.renderNode(ctx, c); err != nil {
			return err
		}

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
		var c node
		if err := d.renderNode(ctx, &c); err != nil {
			return err
		}

		if c.isEnd {
			return nil
		}

		n.Children = append(n.Children, &c)
		n.appendChild(&c)
	}
}

func (d *dom) renderTagAttrs(ctx renderContext, n *node, hasAttr bool) {
	var attrs map[string]string
	if hasAttr {
		attrs = make(map[string]string)
	}

	for hasAttr {
		var tmpK []byte
		var tmpV []byte

		tmpK, tmpV, hasAttr = ctx.tokenizer.TagAttr()
		k := string(tmpK)
		v := string(tmpV)

		switch ctx.namespace {
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

func (d *dom) setEventHandler(ctx renderContext, n *node, k, v string) {
	k = strings.TrimPrefix(k, "on")

	if n.eventCloses == nil {
		n.eventCloses = make(map[string]func())
	}

	v = strings.TrimPrefix(v, "//go:")
	v = strings.ReplaceAll(v, " ", "")
	n.eventCloses[k] = n.addEventListener(ctx, k, v)
}

func (d *dom) closeEventHandler(ctx renderContext, n *node, k string) {
	close, ok := n.eventCloses[k]
	if !ok {
		return
	}

	close()
	delete(n.eventCloses, k)
}

func (d *dom) renderEndTag(ctx renderContext, n *node) error {
	n.isEnd = true
	return nil
}

func (d *dom) renderCompoNode(ctx renderContext, n *node, name string, hasAttr bool) error {
	var compo Compo
	var err error

	if n.isZero() {
		compo, err = d.compoBuilder.new(name)
	} else if name != compoName(n.compo) {
		d.dismount(n)
		compo, err = d.compoBuilder.new(name)
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
	return d.renderCompo(compo, n)
}

func (d *dom) getCompoAttrs(ctx renderContext, hasAttr bool) map[string]string {
	var attrs map[string]string
	if hasAttr {
		attrs = make(map[string]string)
	}

	for hasAttr {
		var k []byte
		var v []byte
		k, v, hasAttr = ctx.tokenizer.TagAttr()
		attrs[string(k)] = string(v)
	}
	return attrs
}

func (d *dom) mount(n *node) {
	d.components[n.compo] = n

	if m, ok := n.compo.(Mounter); ok {
		m.OnMount()
	}
}

func (d *dom) dismount(n *node) {
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
	delete(d.components, n.compo)

	if d, ok := n.compo.(Dismounter); ok {
		d.OnDismount()
	}
}

func (d *dom) setBindingClose(c Compo, close func()) error {
	n, ok := d.components[c]
	if !ok {
		return errors.New("root not found")
	}

	n.bindingCloses = append(n.bindingCloses, close)
	return nil
}

func (d *dom) clean() {
	for _, n := range d.components {
		d.dismount(n)
	}
	d.root = nil
	d.ctxMenuRoot = nil

}

type renderContext struct {
	tokenizer *html.Tokenizer
	compo     Compo
	namespace string
	dom       *dom
}
