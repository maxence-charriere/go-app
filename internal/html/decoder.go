package html

// func decodeNodes(s string) (node, error) {
// 	d := &decoder{
// 		tokenizer: html.NewTokenizer(bytes.NewBufferString(s)),
// 	}
// 	return d.decode()
// }

// type decoder struct {
// 	tokenizer   *html.Tokenizer
// 	decodingSVG bool
// }

// func (d *decoder) decode() (node, error) {
// 	switch d.tokenizer.Next() {
// 	case html.TextToken:
// 		return d.decodeText()

// 	case html.SelfClosingTagToken:
// 		return d.decodeSelfClosingElem()

// 	case html.StartTagToken:
// 		return d.decodeElem()

// 	case html.EndTagToken:
// 		return d.closeElem()

// 	case html.ErrorToken:
// 		return nil, d.tokenizer.Err()

// 	default:
// 		// Nothing we care about, decode the next node.
// 		return d.decode()
// 	}
// }

// func (d *decoder) decodeText() (node, error) {
// 	text := string(d.tokenizer.Text())
// 	text = strings.TrimSpace(text)

// 	if len(text) == 0 {
// 		// Text is empty, decode the next node.
// 		return d.decode()
// 	}

// 	return &textNode{
// 		id:   uuid.New().String(),
// 		text: text,
// 	}, nil
// }

// func (d *decoder) decodeSelfClosingElem() (node, error) {
// 	name, hasAttr := d.tokenizer.TagName()
// 	tagName := string(name)

// 	if isCompoTagName(tagName, d.decodingSVG) {
// 		return d.decodeCompo(tagName, hasAttr)
// 	}

// 	return &elemNode{
// 		id:      uuid.New().String(),
// 		tagName: tagName,
// 		attrs:   d.decodeAttrs(hasAttr),
// 	}, nil
// }

// func (d *decoder) decodeAttrs(hasAttr bool) map[string]string {
// 	if !hasAttr {
// 		return nil
// 	}

// 	attrs := make(map[string]string)
// 	for {
// 		name, val, moreAttr := d.tokenizer.TagAttr()
// 		attrs[string(name)] = string(val)

// 		if !moreAttr {
// 			break
// 		}
// 	}
// 	return attrs
// }

// func (d *decoder) decodeElem() (node, error) {
// 	name, hasAttr := d.tokenizer.TagName()
// 	tagName := string(name)

// 	if isCompoTagName(tagName, d.decodingSVG) {
// 		return d.decodeCompo(tagName, hasAttr)
// 	}

// 	if tagName == "svg" {
// 		d.decodingSVG = true
// 	}

// 	elem := &elemNode{
// 		id:      uuid.New().String(),
// 		tagName: tagName,
// 		attrs:   d.decodeAttrs(hasAttr),
// 	}

// 	if isVoidElem(tagName) {
// 		return elem, nil
// 	}

// 	for {
// 		child, err := d.decode()
// 		if err != nil {
// 			return nil, err
// 		}
// 		if child == nil {
// 			break
// 		}
// 		elem.appendChild(child)
// 	}
// 	return elem, nil
// }

// func (d *decoder) closeElem() (node, error) {
// 	tagName, _ := d.tokenizer.TagName()
// 	if string(tagName) == "svg" {
// 		d.decodingSVG = false
// 	}
// 	return nil, nil
// }

// func (d *decoder) decodeCompo(tagName string, hasAttr bool) (node, error) {
// 	return &compoNode{
// 		id:     uuid.New().String(),
// 		name:   tagName,
// 		fields: d.decodeAttrs(hasAttr),
// 	}, nil
// }

// func isHTMLTagName(tagName string) bool {
// 	return atom.Lookup([]byte(tagName)) != 0
// }

// func isCompoTagName(tagName string, decodingSVG bool) bool {
// 	if decodingSVG {
// 		return false
// 	}
// 	return !isHTMLTagName(tagName)
// }

// var voidElems = map[string]struct{}{
// 	"area":   {},
// 	"base":   {},
// 	"br":     {},
// 	"col":    {},
// 	"embed":  {},
// 	"hr":     {},
// 	"img":    {},
// 	"input":  {},
// 	"keygen": {},
// 	"link":   {},
// 	"meta":   {},
// 	"param":  {},
// 	"source": {},
// 	"track":  {},
// 	"wbr":    {},
// }

// func isVoidElem(tagName string) bool {
// 	_, ok := voidElems[tagName]
// 	return ok
// }
