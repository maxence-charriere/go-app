package dom

import (
	"bytes"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	svgNamespace = "http://www.w3.org/2000/svg"
)

func decodeNodes(s string, t ...Transform) (node, error) {
	d := &decoder{
		tokenizer:  html.NewTokenizer(bytes.NewBufferString(s)),
		transforms: t,
	}

	return d.decode()
}

type decoder struct {
	tokenizer   *html.Tokenizer
	decodingSVG bool
	transforms  []Transform
}

func (d *decoder) decode() (node, error) {
	switch d.tokenizer.Next() {
	case html.TextToken:
		return d.decodeText()

	case html.SelfClosingTagToken:
		return d.decodeSelfClosingElem()

	case html.StartTagToken:
		return d.decodeElem()

	case html.EndTagToken:
		return d.closeElem()

	case html.ErrorToken:
		return nil, d.tokenizer.Err()

	default:
		// Nothing we care about, decode the next node.
		return d.decode()
	}
}

func (d *decoder) decodeText() (node, error) {
	text := string(d.tokenizer.Text())
	text = strings.TrimSpace(text)

	if len(text) == 0 {
		// Text is empty, decode the next node.
		return d.decode()
	}

	t := newText()
	t.SetText(text)
	return t, nil
}

func (d *decoder) decodeSelfClosingElem() (node, error) {
	name, hasAttr := d.tokenizer.TagName()
	tagName := string(name)
	namespace := ""

	if isCompoTagName(tagName, d.decodingSVG) {
		return newCompo(tagName, d.decodeAttrs(hasAttr)), nil
	}

	if d.decodingSVG {
		namespace = svgNamespace
	}

	e := newElem(tagName, namespace)
	e.SetAttrs(d.decodeAttrs(hasAttr))
	return e, nil
}

func (d *decoder) decodeElem() (node, error) {
	name, hasAttr := d.tokenizer.TagName()
	tagName := string(name)
	namespace := ""

	if isCompoTagName(tagName, d.decodingSVG) {
		return newCompo(tagName, d.decodeAttrs(hasAttr)), nil
	}

	if tagName == "svg" {
		d.decodingSVG = true
	}

	if d.decodingSVG {
		namespace = svgNamespace
	}

	e := newElem(tagName, namespace)
	e.SetAttrs(d.decodeAttrs(hasAttr))

	if isVoidElem(tagName) {
		return e, nil
	}

	for {
		child, err := d.decode()
		if err != nil {
			return nil, err
		}
		if child == nil {
			break
		}
		e.appendChild(child)
	}

	return e, nil
}

func (d *decoder) decodeAttrs(hasAttr bool) map[string]string {
	if !hasAttr {
		return nil
	}

	attrs := make(map[string]string)
	for {
		name, val, moreAttr := d.tokenizer.TagAttr()
		n := tagName(string(name))
		v := string(val)

		for _, t := range d.transforms {
			n, v = t(n, v)
		}

		attrs[n] = v

		if !moreAttr {
			break
		}
	}

	return attrs
}

func (d *decoder) closeElem() (node, error) {
	tagName, _ := d.tokenizer.TagName()
	if string(tagName) == "svg" {
		d.decodingSVG = false
	}

	return nil, nil
}

func isHTMLTagName(tagName string) bool {
	return atom.Lookup([]byte(tagName)) != 0
}

func isCompoTagName(tagName string, decodingSVG bool) bool {
	if decodingSVG {
		return false
	}
	return !isHTMLTagName(tagName)
}

var voidElems = map[string]struct{}{
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

func isVoidElem(tagName string) bool {
	_, ok := voidElems[tagName]
	return ok
}

var specialTagNames map[string]string

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

func tagName(n string) string {
	if sn, ok := specialTagNames[n]; ok {
		return sn
	}

	return n
}
