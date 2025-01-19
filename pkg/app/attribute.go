package app

import (
	"strconv"
	"strings"
)

type attributes map[string]string

func (a attributes) Set(name string, value any) {
	switch name {
	case "value":
		a[name] = toString(value)
	case "style", "allow":
		a[name] += toAttributeValue(value) + ";"

	case "class":
		if v := strings.TrimSpace(toAttributeValue(value)); v != "" {
			s := a[name]
			if s != "" {
				s += " "
			}
			a[name] = s + v
		}

	case "srcset":
		s := a[name]
		if s != "" {
			s += ", "
		}
		a[name] = s + toAttributeValue(value)

	default:
		a[name] = toAttributeValue(value)
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
