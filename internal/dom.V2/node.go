package dom

import (
	"strings"

	"golang.org/x/net/html/atom"
)

type node struct {
	ID        string
	ParentID  string
	CompoID   string
	Type      string
	Namespace string
	Text      string
	Attrs     map[string]string
	ChildIDs  []string
	Dom       *Engine
}

type change struct {
	Action     changeAction
	NodeID     string
	Type       string `json:",omitempty"`
	Namespace  string `json:",omitempty"`
	Key        string `json:",omitempty"`
	Value      string `json:",omitempty"`
	ChildID    string `json:",omitempty"`
	NewChildID string `json:",omitempty"`
}

type changeAction int

const (
	newNode changeAction = iota
	delNode
	setAttr
	delAttr
	setText
	appendChild
	removeChild
	replaceChild
)

var (
	svgNamespace    = "http://www.w3.org/2000/svg"
	specialTagNames map[string]string
	voidElems       = map[string]struct{}{
		"area":   {},
		"base":   {},
		"br":     {},
		"col":    {},
		"embed":  {},
		"hr":     {},
		"img":    {},
		"input":  {},
		"keygen": {},
		"link":   {},
		"meta":   {},
		"param":  {},
		"source": {},
		"track":  {},
		"wbr":    {},
	}
)

func init() {
	svgSpecialTagNames := []string{
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

	specialTagNames = make(map[string]string, len(svgSpecialTagNames))
	for _, n := range svgSpecialTagNames {
		specialTagNames[strings.ToLower(n)] = n
	}
}

func nodeType(n string) string {
	if sn, ok := specialTagNames[n]; ok {
		return sn
	}
	return n
}

func isHTMLNode(tagName string) bool {
	return atom.Lookup([]byte(tagName)) != 0
}

func isVoidElem(tagName string) bool {
	_, ok := voidElems[tagName]
	return ok
}

func isCompoNode(tagName, namespace string) bool {
	if len(namespace) != 0 {
		return false
	}
	return !isHTMLNode(tagName)
}
