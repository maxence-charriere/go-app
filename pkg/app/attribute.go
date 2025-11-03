package app

import (
	"strconv"
	"strings"
)

type attributes map[string]string

func (a attributes) Set(name string, value any) {
	var v string
	switch name {
	case "value":
		v = toString(value)

	case "style", "allow":
		v = strings.TrimLeft(a[name]+";"+toAttributeValue(value), ";")

	case "class":
		v = strings.TrimSpace(a[name] + " " + toAttributeValue(value))

	case "srcset":
		v = strings.TrimLeft(a[name]+", "+toAttributeValue(value), ", ")

	default:
		v = toAttributeValue(value)
	}

	switch v {
	case "cite",
		"data",
		"download",
		"href",
		"src",
		"value":
		a[name] = v

	default:
		if v != "" {
			a[name] = v
		}
	}
}

type attributeURLResolver func(string) string

func toAttributeValue(v any) string {
	return strings.TrimSpace(toString(v))
}

func resolveAttributeURLValue(name, value string, resolve attributeURLResolver) string {
	switch name {
	case "cite",
		"data",
		"href",
		"src":
		return resolve(value)

	case "srcset":
		srcs := strings.Split(value, ", ")
		for i, src := range srcs {
			srcs[i] = resolve(src)
		}
		return strings.Join(srcs, ", ")

	default:
		return value
	}
}

func setJSAttribute(jsElement Value, name, value string) {
	toBool := func(v string) bool {
		b, _ := strconv.ParseBool(v)
		return b
	}

	switch name {
	case "value":
		jsElement.Set(name, value)

	case "contenteditable":
		jsElement.Set("contentEditable", value)

	case "ismap":
		jsElement.Set("isMap", toBool(value))

	case "readonly":
		jsElement.Set("readOnly", toBool(value))

	case "async",
		"autofocus",
		"autoplay",
		"checked",
		"default",
		"defer",
		"disabled",
		"hidden",
		"loop",
		"multiple",
		"muted",
		"open",
		"required",
		"reversed",
		"selected":
		jsElement.Set(name, toBool(value))

	default:
		jsElement.Call("setAttribute", name, value)
	}
}

func deleteJSAttribute(jsElement Value, name string) {
	jsElement.Call("removeAttribute", name)
}
