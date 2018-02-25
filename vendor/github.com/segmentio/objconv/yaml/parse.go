package yaml

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/segmentio/objconv"
)

type Parser struct {
	r io.Reader // reader to load bytes from
	s []byte    // string buffer
	// This stack is used to iterate over the arrays and maps that get loaded in
	// the value field.
	stack []parser
}

func NewParser(r io.Reader) *Parser {
	return &Parser{r: r}
}

func (p *Parser) Reset(r io.Reader) {
	p.r = r
	p.s = nil
	p.stack = nil
}

func (p *Parser) Buffered() io.Reader {
	return bytes.NewReader(nil)
}

func (p *Parser) ParseType() (typ objconv.Type, err error) {
	if p.stack == nil {
		var b []byte
		var v interface{}

		if b, err = ioutil.ReadAll(p.r); err != nil {
			return
		}
		if err = yaml.Unmarshal(b, &v); err != nil {
			return
		}
		p.push(newParser(v))
	}

	switch v := p.value(); v.(type) {
	case nil:
		typ = objconv.Nil

	case bool:
		typ = objconv.Bool

	case int, int64:
		typ = objconv.Int

	case uint64:
		typ = objconv.Uint

	case float64:
		typ = objconv.Float

	case string:
		typ = objconv.String

	case yaml.MapSlice:
		typ = objconv.Map

	case []interface{}:
		typ = objconv.Array

	case eof:
		err = io.EOF

	default:
		err = fmt.Errorf("objconv/yaml: gopkg.in/yaml.v2 generated an unsupported value of type %T", v)
	}

	return
}

func (p *Parser) ParseNil() (err error) {
	p.pop()
	return
}

func (p *Parser) ParseBool() (v bool, err error) {
	v = p.pop().value().(bool)
	return
}

func (p *Parser) ParseInt() (v int64, err error) {
	switch x := p.pop().value().(type) {
	case int:
		v = int64(x)
	default:
		v = x.(int64)
	}
	return
}

func (p *Parser) ParseUint() (v uint64, err error) {
	v = p.pop().value().(uint64)
	return
}

func (p *Parser) ParseFloat() (v float64, err error) {
	v = p.pop().value().(float64)
	return
}

func (p *Parser) ParseString() (v []byte, err error) {
	s := p.pop().value().(string)
	n := len(s)

	if cap(p.s) < n {
		p.s = make([]byte, 0, ((n/1024)+1)*1024)
	}

	v = p.s[:n]
	copy(v, s)
	return
}

func (p *Parser) ParseBytes() (v []byte, err error) {
	panic("objconv/yaml: ParseBytes should never be called because YAML has no bytes type, this is likely a bug in the decoder code")
}

func (p *Parser) ParseTime() (v time.Time, err error) {
	panic("objconv/yaml: ParseBytes should never be called because YAML has no time type, this is likely a bug in the decoder code")
}

func (p *Parser) ParseDuration() (v time.Duration, err error) {
	panic("objconv/yaml: ParseDuration should never be called because YAML has no duration type, this is likely a bug in the decoder code")
}

func (p *Parser) ParseError() (v error, err error) {
	panic("objconv/yaml: ParseError should never be called because YAML has no error type, this is likely a bug in the decoder code")
}

func (p *Parser) ParseArrayBegin() (n int, err error) {
	if n = p.top().len(); n != 0 {
		p.push(newParser(p.top().next()))
	}
	return
}

func (p *Parser) ParseArrayEnd(n int) (err error) {
	p.pop()
	return
}

func (p *Parser) ParseArrayNext(n int) (err error) {
	p.push(newParser(p.top().next()))
	return
}

func (p *Parser) ParseMapBegin() (n int, err error) {
	if n = p.top().len(); n != 0 {
		p.push(newParser(p.top().next()))
	}
	return
}

func (p *Parser) ParseMapEnd(n int) (err error) {
	p.pop()
	return
}

func (p *Parser) ParseMapValue(n int) (err error) {
	p.push(newParser(p.top().next()))
	return
}

func (p *Parser) ParseMapNext(n int) (err error) {
	p.push(newParser(p.top().next()))
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

func (p *Parser) push(v parser) {
	p.stack = append(p.stack, v)
}

func (p *Parser) pop() parser {
	i := len(p.stack) - 1
	v := p.stack[i]
	p.stack = p.stack[:i]
	return v
}

func (p *Parser) top() parser {
	return p.stack[len(p.stack)-1]
}

func (p *Parser) value() interface{} {
	n := len(p.stack)
	if n == 0 {
		return eof{}
	}
	return p.stack[n-1].value()
}

type parser interface {
	value() interface{}
	next() interface{}
	len() int
}

type valueParser struct {
	self interface{}
}

func (p *valueParser) value() interface{} {
	return p.self
}

func (p *valueParser) next() interface{} {
	panic("objconv/yaml: invalid call of next method on simple value parser")
}

func (p *valueParser) len() int {
	panic("objconv/yaml: invalid call of len method on simple value parser")
}

type arrayParser struct {
	self []interface{}
	off  int
}

func (p *arrayParser) value() interface{} {
	return p.self
}

func (p *arrayParser) next() interface{} {
	v := p.self[p.off]
	p.off++
	return v
}

func (p *arrayParser) len() int {
	return len(p.self)
}

type mapParser struct {
	self yaml.MapSlice
	off  int
	val  bool
}

func (p *mapParser) value() interface{} {
	return p.self
}

func (p *mapParser) next() (v interface{}) {
	if p.val {
		v = p.self[p.off].Value
		p.val = false
		p.off++
	} else {
		v = p.self[p.off].Key
		p.val = true
	}
	return
}

func (p *mapParser) len() int {
	return len(p.self)
}

func newParser(v interface{}) parser {
	switch x := v.(type) {
	case map[interface{}]interface{}:
		return &mapParser{self: makeMapSlice(x)}

	case []interface{}:
		return &arrayParser{self: x}

	default:
		return &valueParser{self: x}
	}
}

func makeMapSlice(m map[interface{}]interface{}) yaml.MapSlice {
	s := make(yaml.MapSlice, 0, len(m))

	for k, v := range m {
		s = append(s, yaml.MapItem{
			Key:   k,
			Value: v,
		})
	}

	return s
}

// eof values are returned by the top method to indicate that all values have
// already been consumed.
type eof struct{}
