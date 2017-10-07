package html

import "github.com/murlokswarm/app"

func VoidElement(tag app.Tag) bool {
	if tag.Svg {
		return false
	}
	_, ok := voidElems[tag.Name]
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
