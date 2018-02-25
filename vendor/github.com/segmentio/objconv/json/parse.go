package json

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"time"
	"unicode/utf16"
	"unicode/utf8"
	"unsafe"

	"github.com/segmentio/objconv"
	"github.com/segmentio/objconv/objutil"
)

type Parser struct {
	r io.Reader // reader to load bytes from
	s []byte    // buffer used for building strings
	i int       // offset of the first byte in b
	j int       // offset of the last byte in b
	b [128]byte // buffer where bytes are loaded from the reader
	c [128]byte // initial backend array for s
}

func NewParser(r io.Reader) *Parser {
	p := &Parser{r: r}
	p.s = p.c[:0]
	return p
}

func (p *Parser) Reset(r io.Reader) {
	p.r = r
	p.i = 0
	p.j = 0
}

func (p *Parser) Buffered() io.Reader {
	return bytes.NewReader(p.b[p.i:p.j])
}

func (p *Parser) ParseType() (t objconv.Type, err error) {
	var b byte

	if err = p.skipSpaces(); err != nil {
		return
	}

	if b, err = p.peekByteAt(0); err != nil {
		return
	}

	switch {
	case b == '"':
		t = objconv.String

	case b == '{':
		t = objconv.Map

	case b == '[':
		t = objconv.Array

	case b == 'n':
		t = objconv.Nil

	case b == 't':
		t = objconv.Bool

	case b == 'f':
		t = objconv.Bool

	case b == '-' || (b >= '0' && b <= '9'):
		t = objconv.Int

		chunk, _ := p.peekNumber()

		for _, c := range chunk {
			if c == '.' || c == 'e' || c == 'E' {
				t = objconv.Float
				break
			}
		}

		// Cache the result of peekNumber for the following call to ParseInt or
		// ParseFloat.
		p.s = append(p.s[:0], chunk...)

	default:
		err = fmt.Errorf("objconv/json: expected token but found '%c'", b)
	}

	return
}

func (p *Parser) ParseNil() (err error) {
	return p.readToken(nullBytes[:])
}

func (p *Parser) ParseBool() (v bool, err error) {
	var b byte

	if b, err = p.peekByteAt(0); err != nil {
		return
	}

	switch b {
	case 'f':
		v, err = false, p.readToken(falseBytes[:])

	case 't':
		v, err = true, p.readToken(trueBytes[:])

	default:
		err = fmt.Errorf("objconv/json: expected boolean but found '%c'", b)
	}

	return
}

func (p *Parser) ParseInt() (v int64, err error) {
	if v, err = objutil.ParseInt(p.s); err != nil {
		return
	}
	p.i += len(p.s)
	return
}

func (p *Parser) ParseUint() (v uint64, err error) {
	panic("objconv/json: ParseUint should never be called because JSON has no unsigned integer type, this is likely a bug in the decoder code")
}

func (p *Parser) ParseFloat() (v float64, err error) {
	if v, err = strconv.ParseFloat(stringNoCopy(p.s), 64); err != nil {
		return
	}
	p.i += len(p.s)
	return
}

func (p *Parser) ParseString() (v []byte, err error) {
	if p.i == p.j {
		if err = p.fill(); err != nil {
			return
		}
	}

	// fast path: look for an unescaped string in the read buffer.
	if p.i != p.j && p.b[p.i] == '"' {
		chunk := p.b[p.i+1 : p.j]
		off1 := bytes.IndexByte(chunk, '"')
		off2 := bytes.IndexByte(chunk, '\\')

		if off1 >= 0 && off2 < 0 {
			v = p.b[p.i+1 : p.i+1+off1]
			p.i += off1 + 2
			return
		}
	}

	// there are escape characters or the string didn't fit in the read buffer.
	if err = p.readByte('"'); err != nil {
		return
	}

	escaped := false
	v = p.s[:0]

	for {
		var b byte

		if b, err = p.peekByteAt(0); err != nil {
			return
		}
		p.i++

		if escaped {
			escaped = false
			switch b {
			case '"', '\\', '/':
				// simple escaped character
			case 'n':
				b = '\n'

			case 'r':
				b = '\r'

			case 't':
				b = '\t'

			case 'b':
				b = '\b'

			case 'f':
				b = '\f'

			case 'u':
				var r1 rune
				var r2 rune
				if r1, err = p.readUnicode(); err != nil {
					return
				}
				if utf16.IsSurrogate(r1) {
					if r2, err = p.readUnicode(); err != nil {
						return
					}
					r1 = utf16.DecodeRune(r1, r2)
				}
				v = append(v, 0, 0, 0, 0) // make room for 4 bytes
				i := len(v) - 4
				n := utf8.EncodeRune(v[i:], r1)
				v = v[:i+n]
				continue

			default: // not sure what this escape sequence is
				v = append(v, '\\')
			}
		} else if b == '\\' {
			escaped = true
			continue
		} else if b == '"' {
			break
		}

		v = append(v, b)
	}

	p.s = v[:0]
	return
}

func (p *Parser) ParseBytes() (v []byte, err error) {
	panic("objconv/json: ParseBytes should never be called because JOSN has no bytes, this is likely a bug in the decoder code")
}

func (p *Parser) ParseTime() (v time.Time, err error) {
	panic("objconv/json: ParseBytes should never be called because JSON has no time type, this is likely a bug in the decoder code")
}

func (p *Parser) ParseDuration() (v time.Duration, err error) {
	panic("objconv/json: ParseDuration should never be called because JSON has no duration type, this is likely a bug in the decoder code")
}

func (p *Parser) ParseError() (v error, err error) {
	panic("objconv/json: ParseError should never be called because JSON has no error type, this is likely a bug in the decoder code")
}

func (p *Parser) ParseArrayBegin() (n int, err error) {
	return -1, p.readByte('[')
}

func (p *Parser) ParseArrayEnd(n int) (err error) {
	if err = p.skipSpaces(); err != nil {
		return
	}
	return p.readByte(']')
}

func (p *Parser) ParseArrayNext(n int) (err error) {
	var b byte

	if err = p.skipSpaces(); err != nil {
		return
	}

	if b, err = p.peekByteAt(0); err != nil {
		return
	}

	switch {
	case b == ',' && n != 0:
		p.i++
	case b == ']':
		err = objconv.End
	default:
		if n != 0 { // we likely are not in an empty array, there's a value to parse
			err = fmt.Errorf("objconv/json: expected ',' or ']' but found '%c'", b)
		}
	}

	return
}

func (p *Parser) ParseMapBegin() (n int, err error) {
	return -1, p.readByte('{')
}

func (p *Parser) ParseMapEnd(n int) (err error) {
	if err = p.skipSpaces(); err != nil {
		return
	}
	return p.readByte('}')
}

func (p *Parser) ParseMapValue(n int) (err error) {
	if err = p.skipSpaces(); err != nil {
		return
	}
	return p.readByte(':')
}

func (p *Parser) ParseMapNext(n int) (err error) {
	var b byte

	if err = p.skipSpaces(); err != nil {
		return
	}

	if b, err = p.peekByteAt(0); err != nil {
		return
	}

	switch b {
	case ',':
		p.i++
	case '}':
		err = objconv.End
	default:
		if n != 0 { // the map is not empty, likely there's a value to parse
			err = fmt.Errorf("objconv/json: expected ',' or '}' but found '%c'", b)
		}
	}

	return
}

func (p *Parser) TextParser() bool {
	return true
}

func (p *Parser) DecodeBytes(b []byte) (v []byte, err error) {
	var n int
	if n, err = base64.StdEncoding.Decode(b, b); err != nil {
		return
	}
	v = b[:n]
	return
}

func (p *Parser) peek(n int) (b []byte, err error) {
	for (p.i + n) > p.j {
		if err = p.fill(); err != nil {
			return
		}
	}
	b = p.b[p.i : p.i+n]
	return
}

func (p *Parser) peekByteAt(i int) (b byte, err error) {
	for (p.i + i + 1) > p.j {
		if err = p.fill(); err != nil {
			return
		}
	}
	b = p.b[p.i+i]
	return
}

func isNumberByte(b byte) bool {
	return (b >= '0' && b <= '9') || (b == '.') || (b == '+') || (b == '-') || (b == 'e') || (b == 'E')
}

func (p *Parser) peekNumber() (b []byte, err error) {
	// fast path: if the number is loaded in the read buffer we avoid the costly
	// calls to peekByteAt.
	for i, c := range p.b[p.i:p.j] {
		if !isNumberByte(c) {
			b = p.b[p.i : p.i+i]
			return
		}
	}

	// slow path: the number was likely at the end of the read buffer, so there
	// may be some missing digits, loading the read buffer and peeking bytes
	// a non-numeric character is found.
	var i int
	for i = 0; true; i++ {
		var c byte

		if c, err = p.peekByteAt(i); err != nil {
			break
		}

		if !isNumberByte(c) {
			break
		}
	}
	b = p.b[p.i : p.i+i]
	return
}

func (p *Parser) readByte(b byte) (err error) {
	var c byte

	if c, err = p.peekByteAt(0); err == nil {
		if b == c {
			p.i++
		} else {
			err = fmt.Errorf("objconv/json: expected '%c' but found '%c'", b, c)
		}
	}

	return
}

func (p *Parser) readToken(token []byte) (err error) {
	var chunk []byte
	var n = len(token)

	if chunk, err = p.peek(n); err == nil {
		if bytes.Equal(chunk, token) {
			p.i += n
		} else {
			err = fmt.Errorf("objconv/json: expected %#v but found %#v", string(token), string(chunk))
		}
	}

	return
}

func (p *Parser) readUnicode() (r rune, err error) {
	var chunk []byte
	var code uint64

	if chunk, err = p.peek(4); err != nil {
		return
	}

	if code, err = objutil.ParseUintHex(chunk); err != nil {
		err = fmt.Errorf("objconv/json: expected an hexadecimal unicode code point but found %#v", string(chunk))
		return
	}

	if code > objutil.Uint16Max {
		err = fmt.Errorf("objconv/json: expected an hexadecimal unicode code points but found an overflowing value %X", code)
		return
	}

	p.i += 4
	r = rune(code)
	return
}

func (p *Parser) skipSpaces() (err error) {
	for {
		if p.i == p.j {
			if err = p.fill(); err != nil {
				return
			}
		}

		// seek the first byte in the read buffer that isn't a space character.
		for _, b := range p.b[p.i:p.j] {
			switch b {
			case ' ', '\n', '\t', '\r', '\b', '\f':
				p.i++
			default:
				return
			}
		}

		// all trailing bytes in the read buffer were spaces, clear and refill.
		p.i = 0
		p.j = 0
	}
}

func (p *Parser) fill() (err error) {
	n := p.j - p.i
	copy(p.b[:n], p.b[p.i:p.j])
	p.i = 0
	p.j = n

	if n, err = p.r.Read(p.b[p.j:]); n > 0 {
		err = nil
		p.j += n
	} else if err != nil {
		return
	} else {
		err = io.ErrNoProgress
		return
	}

	return
}

func stringNoCopy(b []byte) string {
	n := len(b)
	if n == 0 {
		return ""
	}
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(&b[0])),
		Len:  n,
	}))
}
