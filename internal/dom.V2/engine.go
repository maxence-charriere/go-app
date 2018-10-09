package dom

import (
	"bytes"
	"fmt"
	"strings"
	"sync"

	"github.com/alecthomas/template"
	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

// Engine represents a dom (document object model) engine.
// It manages components an nodes lifecycle and keep track of node changes.
// The engine can be synchronized with a remote dom like a web browser document.
type Engine struct {
	// The factory to decode component from html.
	Factory *app.Factory

	// AttrTransforms describes a set of transformation to apply to parsed node
	// attributes.
	AttrTransforms []func(k, v string) (string, string)

	// AllowedNodeTypes restricts allowed nodes by the given types.
	// No restrictions are enforced if the array is empty.
	AllowedNodeTypes []string

	// Sync is the function used to synchronize node changes with a remote dom.
	// No synchronisations are performed if the func in nil.
	Sync func(v interface{}) error

	once     sync.Once
	compos   map[app.Compo]compo
	compoIDs map[string]compo
	nodes    map[string]node
	creates  []change
	changes  []change
	deletes  []change
}

func (e *Engine) init() {
	if e.Sync == nil {
		e.Sync = func(v interface{}) error {
			return nil
		}
	}

	e.compos = make(map[app.Compo]compo)
	e.compoIDs = make(map[string]compo)
	e.nodes = make(map[string]node)
	e.creates = make([]change, 64)
	e.changes = make([]change, 64)
	e.deletes = make([]change, 64)
}

// Render renders the given component by updating the state described within
// c.Render().
func (e *Engine) Render(c app.Compo) error {
	e.once.Do(e.init)

	// ic, mounted := e.compos[c]
	// if !mounted {
	// 	ic = compo{
	// 		ID: app.CompoName(c) + ":" + uuid.New().String(),
	// 	}

	// 	e.compos[c] = ic
	// 	e.compoIDs[ic.ID] = ic
	// }

	// markup, err := e.compoToHTML(c)
	// if err != nil {
	// 	return errors.Wrap(err, "reading html failed")
	// }

	// oldRoot := e.nodes[ic.rootID]
	// newRoot := node{}

	// newRoot, err = e.render(iterator{
	// 	tokenizer: html.NewTokenizer(bytes.NewBufferString(markup)),
	// 	compoID:   ic.ID,
	// 	nodeID:    ic.rootID,
	// })
	// if err == io.EOF {
	// 	err = nil
	// }
	// if err != nil {
	// 	return errors.Wrap(err, markup)
	// }

	// ic.rootID = newRoot.ID
	// e.compos[c] = ic
	// e.compoIDs[ic.ID] = ic

	// if len(oldRoot.ParentID) != 0 {
	// 	e.replaceNode(oldRoot.ID, newRoot.ID)
	// } else {
	// 	e.deleteNode(oldRoot.ID)
	// }

	// if mounted {
	// 	return nil
	// }

	// if mounter, ok := c.(app.Mounter); ok {
	// 	mounter.OnMount()
	// }

	return nil
}

func (e *Engine) compoToHTML(c app.Compo) (string, error) {
	var extendedFuncs map[string]interface{}
	if extended, ok := c.(app.CompoWithExtendedRender); ok {
		extendedFuncs = extended.Funcs()
	}

	funcs := make(template.FuncMap, len(converters)+len(extendedFuncs))
	for k, v := range converters {
		funcs[k] = v
	}

	for k, v := range extendedFuncs {
		if _, ok := funcs[k]; ok {
			return "", errors.Errorf("template extension can't be named %s", k)
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

	return w.String(), nil
}

func (e *Engine) render(z *html.Tokenizer, n node, namespace string) (node, bool) {
	switch z.Next() {
	case html.TextToken:
		return e.renderText(z, n, namespace)

	case html.SelfClosingTagToken:
		return e.renderSelfClosingTag(z, n, namespace)

	case html.StartTagToken:
		return e.renderStartTag(z, n, namespace)

	case html.EndTagToken:
		return node{}, false

	default:
		return e.render(z, n, namespace)
	}
}

func (e *Engine) renderText(z *html.Tokenizer, n node, namespace string) (node, bool) {
	text := string(z.Text())
	text = strings.TrimSpace(text)

	if len(text) == 0 || len(namespace) != 0 {
		// Invalid text, iterator next node.
		return e.render(z, n, namespace)
	}

	if len(n.ID) == 0 || n.Type != "text" {
		n = node{
			ID:      "text:" + uuid.New().String(),
			CompoID: n.CompoID,
			Type:    "text",
			Dom:     e,
		}
		e.newNode(n)
	}

	if text != n.Text {
		n.Text = text
		e.changes = append(e.changes, change{
			Action: setText,
			NodeID: n.ID,
			Value:  text,
		})
		e.nodes[n.ID] = n
	}

	return n, true
}

func (e *Engine) renderSelfClosingTag(z *html.Tokenizer, n node, namespace string) (node, bool) {
	tagName, hasAttr := z.TagName()
	typ := string(tagName)

	if typ == "svg" {
		namespace = svg
	}

	if isCompoNode(typ, namespace) {
		return e.renderCompo(z, typ, hasAttr)
	}

	if len(n.ID) == 0 || n.Type != typ {
		n = node{
			ID:      typ + ":" + uuid.New().String(),
			CompoID: n.CompoID,
			Type:    typ,
			Dom:     e,
		}
		e.newNode(n)
	}

	n = e.renderTagAttrs(z, n, hasAttr)

	for _, childID := range n.ChildIDs {
		e.deleteNode(childID)
	}

	n.ChildIDs = n.ChildIDs[:0]
	e.nodes[n.ID] = n
	return n, true
}

func (e *Engine) renderStartTag(z *html.Tokenizer, n node, namespace string) (node, bool) {
	tagName, hasAttr := z.TagName()
	typ := string(tagName)

	if typ == "svg" {
		namespace = svg
	}

	if isCompoNode(typ, namespace) {
		return e.renderCompo(z, typ, hasAttr)
	}

	if len(n.ID) == 0 || n.Type != typ {
		n = node{
			ID:      typ + ":" + uuid.New().String(),
			CompoID: n.CompoID,
			Type:    typ,
			Dom:     e,
		}
		e.newNode(n)
	}

	n = e.renderTagAttrs(z, n, hasAttr)

	childIDs := n.ChildIDs
	moreChild := true
	count := 0

	// Replace children:
	for len(childIDs) != 0 {
		old := e.nodes[childIDs[0]]
		new := node{}

		if new, moreChild = e.render(z, old, namespace); !moreChild {
			break
		}

		if new.ID != old.ID {
			e.changes = append(e.changes, change{
				Action:     replaceChild,
				NodeID:     n.ID,
				ChildID:    old.ID,
				NewChildID: new.ID,
			})

			childIDs[0] = new.ID
			e.deleteNode(old.ID)
		}

		count++
		childIDs = childIDs[1:]
	}

	// Remove children:
	for i := count; i < len(childIDs); i++ {
		childID := childIDs[i]

		e.deleteNode(childID)
		e.changes = append(e.changes, change{
			Action:  removeChild,
			NodeID:  n.ID,
			ChildID: childID,
		})

	}
	childIDs = childIDs[:count]

	// Add children
	for moreChild {
		c := node{CompoID: n.CompoID}
		c, moreChild = e.render(z, c, namespace)

		if !moreChild {
			break
		}

		childIDs = append(childIDs, c.ID)
		e.changes = append(e.changes, change{
			Action:  appendChild,
			NodeID:  n.ID,
			ChildID: c.ID,
		})
	}

	n.ChildIDs = childIDs
	e.nodes[n.ID] = n
	return n, true
}

func (e *Engine) renderTagAttrs(z *html.Tokenizer, n node, moreAttr bool) node {
	attrs := make(map[string]string)

	for moreAttr {
		var rk []byte
		var rv []byte

		rk, rv, moreAttr = z.TagAttr()
		k := string(rk)
		v := string(rv)

		for _, t := range e.AttrTransforms {
			k, v = t(k, v)
		}

		attrs[k] = v
		e.changes = append(e.changes, change{
			Action: setAttr,
			NodeID: n.ID,
			Key:    k,
			Value:  v,
		})
	}

	for k := range n.Attrs {
		if _, ok := attrs[k]; !ok {
			continue
		}

		e.changes = append(e.changes, change{
			Action: delAttr,
			NodeID: n.ID,
			Key:    k,
		})
	}

	n.Attrs = attrs
	e.nodes[n.ID] = n
	return n
}

func (e *Engine) renderCompo(z *html.Tokenizer, typ string, hasAttr bool) (node, bool) {
	panic("not implemented")
}

func (e *Engine) nodesByIDs(ids ...string) []node {
	nodes := make([]node, 0, len(ids))
	for _, id := range ids {
		nodes = append(nodes, e.nodes[id])
	}

	return nodes
}

func (e *Engine) newNode(n node) {
	e.nodes[n.ID] = n

	e.creates = append(e.creates, change{
		Action:    newNode,
		NodeID:    n.ID,
		Type:      n.Type,
		Namespace: n.Namespace,
	})
}

func (e *Engine) deleteNode(id string) {
	n, ok := e.nodes[id]
	if !ok {
		return
	}

	for _, childID := range n.ChildIDs {
		e.deleteNode(childID)
	}

	c, mounted := e.compoIDs[n.CompoID]
	if !mounted {
		return
	}

	if c.rootID != n.ID {
		return
	}

	if dismounter, ok := c.compo.(app.Dismounter); ok {
		dismounter.OnDismount()
	}

	e.deletes = append(e.deletes, change{
		Action: delNode,
		NodeID: n.ID,
	})

	delete(e.nodes, n.ID)
	delete(e.compos, c.compo)
	delete(e.compoIDs, c.ID)
}
