package app

import "strings"

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

	switch name {
	case "cite",
		"data",
		"download",
		"href",
		"src",
		"value":
		a[name] = v

	default:
		if v != "" || isBooleanAttribute(name) {
			a[name] = v
		}
	}
}

type attributeURLResolver func(string) string

// isBooleanAttribute identifies attributes whose true/false values control
// presence. Hidden also supports this API, but other values (such as
// "until-found") must be preserved because it is an enumerated attribute.
func isBooleanAttribute(name string) bool {
	switch strings.ToLower(name) {
	case "allowfullscreen",
		"allowpaymentrequest",
		"async",
		"autofocus",
		"autoplay",
		"checked",
		"controls",
		"default",
		"defer",
		"disabled",
		"disablepictureinpicture",
		"disableremoteplayback",
		"formnovalidate",
		"hidden",
		"inert",
		"ismap",
		"itemscope",
		"loop",
		"multiple",
		"muted",
		"nomodule",
		"novalidate",
		"open",
		"playsinline",
		"readonly",
		"required",
		"reversed",
		"selected":
		return true

	default:
		return false
	}
}

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
	// Update the current value, not only the default value attribute.
	switch name {
	case "value":
		jsElement.Set(name, value)

	// Attributes specify defaults; properties update the current state.
	case "checked",
		"muted",
		"selected":
		jsElement.Set(name, value != "false")

	// Removing the attribute does not clear a new script's force-async flag.
	case "async":
		jsElement.Set(name, value != "false")

	default:
		if isBooleanAttribute(name) {
			switch value {
			case "false":
				deleteJSAttribute(jsElement, name)
				return

			case "true":
				value = ""
			}
		}
		jsElement.Call("setAttribute", name, value)
	}
}

func deleteJSAttribute(jsElement Value, name string) {
	jsElement.Call("removeAttribute", name)
}
