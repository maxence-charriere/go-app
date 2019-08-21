package maestro

import (
	"errors"
	"strings"

	"golang.org/x/net/html/atom"
)

// JSNode is the interface that describes a javascript node.
type JSNode interface {
	new(tag string, namespace string)
	newText(s string)
	updateText(s string)
	appendChild(c JSNode)
	removeChild(c JSNode)
	replaceChild(old, new JSNode)
	replace(new JSNode)
	upsertAttr(k string, v string)
	deleteAttr(k string)
}

// Node represents a dom node.
type Node struct {
	Kind     NodeKind
	Name     string
	Text     string
	Attrs    map[string]string
	Children []*Node

	jsNode JSNode
	compo  Compo
}

func (n *Node) isZero() bool {
	return n.Kind == StdNode &&
		n.Name == "" &&
		n.Text == "" &&
		n.Attrs == nil &&
		n.Children == nil &&
		n.jsNode == nil &&
		n.compo == nil
}

func (n *Node) jsRoot() JSNode {
	switch n.Kind {
	case CompoNode:
		return n.Children[0].jsRoot()

	default:
		return n.jsNode
	}
}

func (n *Node) appendChild(c *Node) error {
	if n.Kind == CompoNode && len(n.Children) >= 1 {
		return errors.New("component can't have more than a child")
	}

	switch n.Kind {
	case StdNode:
		n.Children = append(n.Children, c)
		n.jsNode.appendChild(c.jsRoot())

	case TextNode:
		return errors.New("text can't have children")

	case CompoNode:
		n.Children = append(n.Children, c)
	}

	return nil
}

func (n *Node) removeChild(c *Node) error {
	if n.Kind == TextNode || n.Kind == CompoNode {
		return errors.New("child can't be removed from node")
	}

	children := c.Children
	for i, child := range children {
		if child == c {
			copy(children[i:], children[i+1:])
			children[len(children)-1] = nil
			children = children[:len(children)-1]

			n.Children = children
			n.jsNode.removeChild(c.jsRoot())
			return nil
		}
	}

	return errors.New("child to remove not found in node")
}

func (n *Node) replaceChild(old, new *Node) error {
	if n.Kind == TextNode {
		return errors.New("text node child can't be replaced")
	}

	for i, c := range n.Children {
		if c == old {
			n.Children[i] = new

			switch n.Kind {
			case StdNode:
				n.jsNode.replaceChild(old.jsRoot(), new.jsRoot())

			case CompoNode:
				old.jsRoot().replace(new.jsRoot())
			}

			return nil
		}
	}

	return errors.New("child to replace not found in node")
}

// NodeKind represents a node kind.
type NodeKind byte

// Constants that enumerate nodes kind.
const (
	StdNode NodeKind = iota
	TextNode
	CompoNode
)

var (
	svgNamespace    = "http://www.w3.org/2000/svg"
	svgSpecialAttrs map[string]string
	voidElems       = map[string]struct{}{
		"area":     {},
		"base":     {},
		"br":       {},
		"col":      {},
		"embed":    {},
		"hr":       {},
		"img":      {},
		"input":    {},
		"keygen":   {},
		"link":     {},
		"meta":     {},
		"param":    {},
		"source":   {},
		"track":    {},
		"wbr":      {},
		"menuitem": {},
	}
)

func init() {
	svgSpecialAttrNames := []string{
		"allowReorder",
		"attributeName",
		"attributeType",
		"autoReverse",

		"baseFrequency",
		"baseProfile",

		"calcMode",
		"clipPathUnits",
		"contentScriptType",
		"contentStyleType",

		"diffuseConstant",

		"externalResourcesRequired",

		"filterRes",
		"filterUnits",

		"glyphRef",
		"gradientTransform",
		"gradientUnits",

		"kernelMatrix",
		"kernelUnitLength",
		"keyPoints",
		"keySplines",
		"keyTimes",

		"lengthAdjust",
		"limitingConeAngle",

		"markerHeight",
		"markerUnits",
		"markerWidth",
		"maskContentUnits",
		"maskUnits",

		"numOctaves",

		"pathLength",
		"patternContentUnits",
		"patternTransform",
		"patternUnits",
		"pointsAtX",
		"pointsAtY",
		"pointsAtZ",
		"preserveAlpha",
		"preserveAspectRatio",
		"primitiveUnits",

		"referrerPolicy",
		"refX",
		"refY",
		"repeatCount",
		"repeatDur",
		"requiredExtensions",
		"requiredFeatures",

		"specularConstant",
		"specularExponent",
		"spreadMethod",
		"startOffset",
		"stdDeviation",
		"stitchTiles",
		"surfaceScale",
		"systemLanguage",

		"tableValues",
		"targetX",
		"targetY",
		"textLength",

		"viewBox",
		"viewTarget",

		"xChannelSelector",

		"yChannelSelector",

		"zoomAndPan",
	}

	svgSpecialAttrs = make(map[string]string, len(svgSpecialAttrNames))
	for _, n := range svgSpecialAttrNames {
		svgSpecialAttrs[strings.ToLower(n)] = n
	}
}

func svgAttr(k string) string {
	if sk, ok := svgSpecialAttrs[k]; ok {
		return sk
	}
	return k
}

func isCompoNode(tagName, namespace string) bool {
	if len(namespace) != 0 {
		return false
	}
	return !isHTMLNode(tagName)
}

func isHTMLNode(tagName string) bool {
	return atom.Lookup([]byte(tagName)) != 0
}

func isVoidElem(tagName string) bool {
	_, ok := voidElems[tagName]
	return ok
}
