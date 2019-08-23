package maestro

import (
	"strings"

	"golang.org/x/net/html/atom"
)

// Node represents a dom node.
type Node struct {
	Name      string
	Text      string
	Attrs     map[string]string
	Children  []*Node
	CompoName string

	compo  Compo
	jsNode jsNode
	isEnd  bool
}

func (n *Node) isZero() bool {
	return n.Name == "" &&
		n.Text == "" &&
		n.Attrs == nil &&
		n.Children == nil &&
		n.CompoName == "" &&
		n.compo == nil &&
		!n.isEnd

}

func (n *Node) isCompoRoot() bool {
	return n.CompoName != ""
}

var (
	namespaces = map[string]string{
		"svg": "http://www.w3.org/2000/svg",
	}
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

func isCompoNode(name, namespace string) bool {
	if len(namespace) != 0 {
		return false
	}
	return !isHTMLNode(name)
}

func isHTMLNode(name string) bool {
	return atom.Lookup([]byte(name)) != 0
}

func isVoidElem(name string) bool {
	_, ok := voidElems[name]
	return ok
}
