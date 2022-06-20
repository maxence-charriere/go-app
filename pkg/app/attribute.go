package app

import (
	"strconv"
	"strings"
)

type attributes map[string]string

func (a attributes) Set(name string, value any) {
	switch name {
	case "style", "allow":
		a[name] += toAttributeValue(value) + ";"

	case "class":
		s := a[name]
		if s != "" {
			s += " "
		}
		s += toAttributeValue(value)
		a[name] = s

	default:
		a[name] = toAttributeValue(value)
	}
}

func (a attributes) Mount(jsElement Value, resolveURL attributeURLResolver) {
	for name, value := range a {
		setJSAttribute(jsElement, name, resolveAttributeURLValue(
			name,
			value,
			resolveURL,
		))
	}
}

func (a attributes) Update(jsElement Value, b attributes, resolveURL attributeURLResolver) {
	for name := range a {
		if _, ok := b[name]; !ok {
			deleteJSAttribute(jsElement, name)
			delete(a, name)
		}
	}

	for name, value := range b {
		if a[name] == value {
			continue
		}

		a[name] = value
		setJSAttribute(jsElement, name, resolveAttributeURLValue(
			name,
			value,
			resolveURL,
		))
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
		"src",
		"srcset":
		return resolve(value)

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

	case "class":
		jsElement.Set("className", value)

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
