package app

import (
	"strconv"
	"strings"
)

type attributes map[string]string

func (a attributes) Set(name string, value any) {
	switch name {
	case "style", "allow":
		var b strings.Builder
		b.WriteString(a[name])
		b.WriteString(toAttributeValue(value))
		b.WriteByte(';')
		a[name] = b.String()

	case "class":
		var b strings.Builder
		b.WriteString(a[name])
		if b.Len() != 0 {
			b.WriteByte(' ')
		}
		b.WriteString(toAttributeValue(value))
		a[name] = b.String()

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
		jsElement.Set("value", value)

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
		jsElement.setAttr(name, value)
	}
}
