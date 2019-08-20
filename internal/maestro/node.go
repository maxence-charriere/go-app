package maestro

import (
	"strings"

	"golang.org/x/net/html/atom"
)

// JSNode is the interface that describes a javascript node.
type JSNode interface {
	new(string, string) error
	newText(s string) error
	updateText(s string)
	appendChild(c JSNode)
	removeChild(c JSNode)
	changeType(string, string) error
	upsertAttr(string, string)
	deleteAttr(string)
	delete() error
}

// Node represents a dom node.
type Node struct {
	Type     string
	Text     string
	IsCompo  bool
	Attrs    map[string]string
	Compo    Compo
	Parent   *Node
	Children []*Node

	jsNode JSNode
}

// JSNode returns the underlying javascript node.
func (n *Node) JSNode() JSNode {
	if !n.IsCompo {
		return n.jsNode
	}

	panic("not implemented")
}

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
