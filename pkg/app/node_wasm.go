package app

import (
	"reflect"
	"strconv"
	"strings"
	"syscall/js"

	"github.com/maxence-charriere/go-app/v5/pkg/log"
	"golang.org/x/net/html/atom"
)

// CompoHandler represents a function that can be associated to an html node
// event.
type CompoHandler func(src, event js.Value)

type node struct {
	js.Value

	Name      string
	Text      string
	CompoName string
	Attrs     map[string]string
	Children  []*node

	compo         Compo
	eventCloses   map[string]func()
	bindingCloses []func()
	isEnd         bool
}

func (n *node) isZero() bool {
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

func (n *node) isCompoRoot() bool {
	return n.CompoName != ""
}

func (n *node) new(tag, namespace string) {
	if namespace != "" {
		n.Value = js.Global().Get("document").Call("createElementNS", namespace, tag)
	} else {
		n.Value = js.Global().Get("document").Call("createElement", tag)
	}
}

func (n *node) newText() {
	n.Value = js.Global().Get("document").Call("createTextNode", "")
}

func (n *node) change(tag, namespace string) {
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

func (n *node) updateText(s string) {
	n.Set("nodeValue", s)
}

func (n *node) appendChild(c *node) {
	n.Call("appendChild", c)
}

func (n *node) removeChild(c *node) {
	n.Call("removeChild", c)
}

func (n *node) upsertAttr(k, v string) {
	n.Call("setAttribute", k, v)
}

func (n *node) deleteAttr(k string) {
	n.Call("removeAttribute", k)
}

func (n *node) addEventListener(ctx renderContext, eventname string, target string) func() {
	cb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var event js.Value
		if len(args) >= 1 {
			event = args[0]
		}
		ctx.dom.trackCursorPosition(event)

		if eventname == "contextmenu" {
			event.Call("preventDefault")
		}

		if strings.HasPrefix(target, "emit:") {
			msg := strings.TrimPrefix(target, "emit:")
			go ctx.dom.msgs.emit(msg, this, event)
			return nil
		}

		ctx.dom.callOnUI(func() {
			n.executeGoCallback(ctx, target, this, event)
		})
		return nil
	})
	n.Call("addEventListener", eventname, cb)

	return func() {
		n.Call("removeEventListener", eventname, cb)
		cb.Release()
	}
}

func (n *node) executeGoCallback(ctx renderContext, target string, this, event js.Value) {
	recv, err := getReceiver(ctx.compo, target)
	if err != nil {
		log.Error("calling event listener failed").
			T("reason", err).
			T("component", reflect.TypeOf(ctx.compo)).
			T("target", target)
		return
	}

	switch recv.Kind() {
	case reflect.Func:
		if jsHandlerType != recv.Type() {
			log.Error("calling event listener failed").
				T("reason", "bad receiver function type").
				T("component", reflect.TypeOf(ctx.compo)).
				T("target", target).
				T("expected type", jsHandlerType).
				T("receiver type", recv.Type())
			return
		}
		recv.Call([]reflect.Value{
			reflect.ValueOf(this),
			reflect.ValueOf(event),
		})

	case reflect.String:
		value := this.Get("value")
		recv.SetString(value.String())
		ctx.dom.render(ctx.compo)

	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		value := this.Get("value").String()
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			log.Error("adding event listener failed").
				T("reason", err).
				T("component", reflect.TypeOf(ctx.compo)).
				T("target", target)
			return
		}
		recv.SetInt(i)
		ctx.dom.render(ctx.compo)

	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		value := this.Get("value").String()
		u, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			log.Error("adding event listener failed").
				T("reason", err).
				T("component", reflect.TypeOf(ctx.compo)).
				T("target", target)
			return
		}
		recv.SetUint(u)
		ctx.dom.render(ctx.compo)

	case reflect.Float64, reflect.Float32:
		value := this.Get("value").String()
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Error("adding event listener failed").
				T("reason", err).
				T("component", reflect.TypeOf(ctx.compo)).
				T("target", target)
			return
		}
		recv.SetFloat(f)
		ctx.dom.render(ctx.compo)

	case reflect.Bool:
		value := this.Get("value").String()
		b, err := strconv.ParseBool(value)
		if err != nil {
			log.Error("adding event listener failed").
				T("reason", err).
				T("component", reflect.TypeOf(ctx.compo)).
				T("target", target)
			return
		}
		recv.SetBool(b)
		ctx.dom.render(ctx.compo)

	default:
		log.Error("adding event listener failed").
			T("reason", "unsupported target kind").
			T("component", reflect.TypeOf(ctx.compo)).
			T("target", target).
			T("target type", recv.Type())
	}
}

var (
	jsHandlerType = reflect.TypeOf(func(s, e js.Value) {})
)

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
	return strings.HasPrefix(k, "on") && strings.HasPrefix(v, "//go:")
}
