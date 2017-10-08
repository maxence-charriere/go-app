package html

import "golang.org/x/net/html/atom"

func isVoidElement(name string, decodingSvg bool) bool {
	if decodingSvg {
		return false
	}
	_, ok := voidElems[name]
	return ok
}

var (
	voidElems = map[string]struct{}{
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

func isComponent(name string, decodingSvg bool) bool {
	if len(name) == 0 {
		return false
	}
	if decodingSvg {
		return false
	}

	// Any non standard html tag name describes a component name.
	return atom.Lookup([]byte(name)) == 0
}
