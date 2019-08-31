// +build js

package maestro

import (
	"reflect"
	"strconv"
	"strings"
	"syscall/js"

	"github.com/maxence-charriere/app/pkg/log"
	"golang.org/x/net/html/atom"
)

// Node represents a dom node.
type Node struct {
	js.Value

	Name      string
	Text      string
	Attrs     map[string]string
	Children  []*Node
	CompoName string

	compo         Compo
	eventCloses   map[string]func()
	bindingCloses []func()
	isEnd         bool
}

func (n *Node) isZero() bool {
	return n.Name == "" &&
		n.Text == "" &&
		n.Attrs == nil &&
		n.Children == nil &&
		n.CompoName == "" &&
		n.compo == nil &&
		n.eventCloses == nil &&
		n.bindingCloses == nil &&
		!n.isEnd
}

func (n *Node) isCompoRoot() bool {
	return n.CompoName != ""
}

func (n *Node) new(tag, namespace string) {
	if namespace != "" {
		n.Value = js.Global().Get("document").Call("createElementNS", namespace, tag)
	} else {
		n.Value = js.Global().Get("document").Call("createElement", tag)
	}
}

func (n *Node) newText() {
	n.Value = js.Global().Get("document").Call("createTextNode", "")
}

func (n *Node) change(tag, namespace string) {
	parent := n.Get("parentNode")
	if t := parent.Type(); t == js.TypeUndefined || t == js.TypeNull {
		panic("parentNode is not set")
	}

	old := n.Value

	if tag == "" {
		n.Value = js.Global().Get("document").Call("createTextNode", "")
	} else if namespace != "" {
		n.Value = js.Global().Get("document").Call("createElementNS", namespace, tag)
	} else {
		n.Value = js.Global().Get("document").Call("createElement", tag)
	}

	parent.Call("replaceChild", n, old)
}

func (n *Node) updateText(s string) {
	n.Set("nodeValue", s)
}

func (n *Node) appendChild(c *Node) {
	n.Call("appendChild", c)
}

func (n *Node) removeChild(c *Node) {
	n.Call("removeChild", c)
}

func (n *Node) upsertAttr(k, v string) {
	n.Call("setAttribute", k, v)
}

func (n *Node) deleteAttr(k string) {
	n.Call("removeAttribute", k)
}

func (n *Node) addEventListener(ctx renderContext, event string, target string) func() {
	preventDefault := event == "contextmenu"

	execBinding := func(this js.Value, args []js.Value) interface{} {
		ctx.Dom.CallOnUI(func() {
			var event js.Value
			if len(args) >= 1 {
				event = args[0]
			}

			if preventDefault {
				event.Call("preventDefault")
			}

			recv, err := getReceiver(ctx.Compo, target)
			if err != nil {
				log.Error("adding event listener failed").
					T("reason", err).
					T("component", reflect.TypeOf(ctx.Compo)).
					T("target", target)
				return
			}

			switch recv.Kind() {
			case reflect.Func:
				if reflect.TypeOf(func(s, e js.Value) {}) == recv.Type() {
					recv.Call([]reflect.Value{
						reflect.ValueOf(this),
						reflect.ValueOf(event),
					})
				}

			case reflect.String:
				value := this.Get("value")
				recv.SetString(value.String())
				ctx.Dom.Render(ctx.Compo)

			case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
				value := this.Get("value").String()
				i, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					log.Error("adding event listener failed").
						T("reason", err).
						T("component", reflect.TypeOf(ctx.Compo)).
						T("target", target)
					return
				}
				recv.SetInt(i)
				ctx.Dom.Render(ctx.Compo)

			case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
				value := this.Get("value").String()
				u, err := strconv.ParseUint(value, 10, 64)
				if err != nil {
					log.Error("adding event listener failed").
						T("reason", err).
						T("component", reflect.TypeOf(ctx.Compo)).
						T("target", target)
					return
				}
				recv.SetUint(u)
				ctx.Dom.Render(ctx.Compo)

			case reflect.Float64, reflect.Float32:
				value := this.Get("value").String()
				f, err := strconv.ParseFloat(value, 64)
				if err != nil {
					log.Error("adding event listener failed").
						T("reason", err).
						T("component", reflect.TypeOf(ctx.Compo)).
						T("target", target)
					return
				}
				recv.SetFloat(f)
				ctx.Dom.Render(ctx.Compo)

			case reflect.Bool:
				value := this.Get("value").String()
				b, err := strconv.ParseBool(value)
				if err != nil {
					log.Error("adding event listener failed").
						T("reason", err).
						T("component", reflect.TypeOf(ctx.Compo)).
						T("target", target)
					return
				}
				recv.SetBool(b)
				ctx.Dom.Render(ctx.Compo)

			default:
				log.Error("adding event listener failed").
					T("reason", "unsupported target kind").
					T("component", reflect.TypeOf(ctx.Compo)).
					T("target", target).
					T("target type", recv.Type())
			}
		})

		return nil
	}

	cb := js.FuncOf(execBinding)
	n.Call("addEventListener", event, cb)

	return func() {
		n.Call("removeEventListener", event, cb)
		cb.Release()
	}
}

type eventHandler func(js.Value)

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

func isGoEventAttr(k, v string) bool {
	return strings.HasPrefix(k, "on") && strings.HasPrefix(v, "//go: ")
}
