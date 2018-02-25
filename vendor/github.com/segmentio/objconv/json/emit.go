package json

import (
	"encoding/base64"
	"errors"
	"io"
	"math"
	"strconv"
	"time"

	"github.com/segmentio/objconv"
	"github.com/segmentio/objconv/objutil"
)

var (
	nullBytes  = [...]byte{'n', 'u', 'l', 'l'}
	trueBytes  = [...]byte{'t', 'r', 'u', 'e'}
	falseBytes = [...]byte{'f', 'a', 'l', 's', 'e'}

	arrayOpen  = [...]byte{'['}
	arrayClose = [...]byte{']'}

	mapOpen  = [...]byte{'{'}
	mapClose = [...]byte{'}'}

	comma  = [...]byte{','}
	column = [...]byte{':'}

	newline = [...]byte{'\n'}
	spaces  = [...]byte{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '}
)

// Emitter implements a JSON emitter that satisfies the objconv.Emitter
// interface.
type Emitter struct {
	w io.Writer
	s []byte
	a [128]byte
}

func NewEmitter(w io.Writer) *Emitter {
	e := &Emitter{w: w}
	e.s = e.a[:0]
	return e
}

func (e *Emitter) Reset(w io.Writer) {
	e.w = w
}

func (e *Emitter) EmitNil() (err error) {
	_, err = e.w.Write(nullBytes[:])
	return
}

func (e *Emitter) EmitBool(v bool) (err error) {
	if v {
		_, err = e.w.Write(trueBytes[:])
	} else {
		_, err = e.w.Write(falseBytes[:])
	}
	return
}

func (e *Emitter) EmitInt(v int64, _ int) (err error) {
	_, err = e.w.Write(strconv.AppendInt(e.s[:0], v, 10))
	return
}

func (e *Emitter) EmitUint(v uint64, _ int) (err error) {
	_, err = e.w.Write(strconv.AppendUint(e.s[:0], v, 10))
	return
}

func (e *Emitter) EmitFloat(v float64, bitSize int) (err error) {
	switch {
	case math.IsNaN(v):
		err = errors.New("NaN has no json representation")

	case math.IsInf(v, +1):
		err = errors.New("+Inf has no json representation")

	case math.IsInf(v, -1):
		err = errors.New("-Inf has no json representation")

	default:
		_, err = e.w.Write(strconv.AppendFloat(e.s[:0], v, 'g', -1, bitSize))
	}
	return
}

func (e *Emitter) EmitString(v string) (err error) {
	i := 0
	j := 0
	n := len(v)
	s := append(e.s[:0], '"')

	for j != n {
		b := v[j]
		j++

		switch b {
		case '"', '\\':
			// b = b

		case '\b':
			b = 'b'

		case '\f':
			b = 'f'

		case '\n':
			b = 'n'

		case '\r':
			b = 'r'

		case '\t':
			b = 't'

		default:
			continue
		}

		s = append(s, v[i:j-1]...)
		s = append(s, '\\', b)
		i = j
	}

	s = append(s, v[i:j]...)
	s = append(s, '"')
	e.s = s[:0] // in case the buffer was reallocated

	_, err = e.w.Write(s)
	return
}

func (e *Emitter) EmitBytes(v []byte) (err error) {
	s := e.s[:0]
	n := base64.StdEncoding.EncodedLen(len(v)) + 2

	if cap(s) < n {
		s = make([]byte, 0, align(n, 1024))
		e.s = s
	}

	s = s[:n]
	s[0] = '"'
	base64.StdEncoding.Encode(s[1:], v)
	s[n-1] = '"'

	_, err = e.w.Write(s)
	return
}

func (e *Emitter) EmitTime(v time.Time) (err error) {
	s := e.s[:0]

	s = append(s, '"')
	s = v.AppendFormat(s, time.RFC3339Nano)
	s = append(s, '"')

	e.s = s[:0]
	_, err = e.w.Write(s)
	return
}

func (e *Emitter) EmitDuration(v time.Duration) (err error) {
	s := e.s[:0]

	s = append(s, '"')
	s = objutil.AppendDuration(s, v)
	s = append(s, '"')

	e.s = s[:0]
	_, err = e.w.Write(s)
	return
}

func (e *Emitter) EmitError(v error) (err error) {
	return e.EmitString(v.Error())
}

func (e *Emitter) EmitArrayBegin(_ int) (err error) {
	_, err = e.w.Write(arrayOpen[:])
	return
}

func (e *Emitter) EmitArrayEnd() (err error) {
	_, err = e.w.Write(arrayClose[:])
	return
}

func (e *Emitter) EmitArrayNext() (err error) {
	_, err = e.w.Write(comma[:])
	return
}

func (e *Emitter) EmitMapBegin(_ int) (err error) {
	_, err = e.w.Write(mapOpen[:])
	return
}

func (e *Emitter) EmitMapEnd() (err error) {
	_, err = e.w.Write(mapClose[:])
	return
}

func (e *Emitter) EmitMapValue() (err error) {
	_, err = e.w.Write(column[:])
	return
}

func (e *Emitter) EmitMapNext() (err error) {
	_, err = e.w.Write(comma[:])
	return
}

func (e *Emitter) TextEmitter() bool {
	return true
}

func (e *Emitter) PrettyEmitter() objconv.Emitter {
	return NewPrettyEmitter(e.w)
}

func align(n int, a int) int {
	if (n % a) == 0 {
		return n
	}
	return ((n / a) + 1) * a
}

type PrettyEmitter struct {
	Emitter
	i int
	s []int
	a [8]int
}

func NewPrettyEmitter(w io.Writer) *PrettyEmitter {
	e := &PrettyEmitter{
		Emitter: *NewEmitter(w),
	}
	e.s = e.a[:0]
	return e
}

func (e *PrettyEmitter) Reset(w io.Writer) {
	e.Emitter.Reset(w)
	e.i = 0
	e.s = e.s[:0]
}

func (e *PrettyEmitter) EmitArrayBegin(n int) (err error) {
	if err = e.Emitter.EmitArrayBegin(n); err != nil {
		return
	}
	if e.push(n) != 0 {
		err = e.indent()
	}
	return
}

func (e *PrettyEmitter) EmitArrayEnd() (err error) {
	if e.pop() != 0 {
		if err = e.indent(); err != nil {
			return
		}
	}
	return e.Emitter.EmitArrayEnd()
}

func (e *PrettyEmitter) EmitArrayNext() (err error) {
	if err = e.Emitter.EmitArrayNext(); err != nil {
		return
	}
	return e.indent()
}

func (e *PrettyEmitter) EmitMapBegin(n int) (err error) {
	if err = e.Emitter.EmitMapBegin(n); err != nil {
		return
	}
	if e.push(n) != 0 {
		err = e.indent()
	}
	return
}

func (e *PrettyEmitter) EmitMapEnd() (err error) {
	if e.pop() != 0 {
		if err = e.indent(); err != nil {
			return
		}
	}
	return e.Emitter.EmitMapEnd()
}

func (e *PrettyEmitter) EmitMapValue() (err error) {
	if err = e.Emitter.EmitMapValue(); err != nil {
		return
	}
	_, err = e.w.Write(spaces[:1])
	return
}

func (e *PrettyEmitter) EmitMapNext() (err error) {
	if err = e.Emitter.EmitMapNext(); err != nil {
		return
	}
	return e.indent()
}

func (e *PrettyEmitter) TextEmitter() bool {
	return true
}

func (e *PrettyEmitter) indent() (err error) {
	if _, err = e.w.Write(newline[:]); err != nil {
		return
	}

	for n := 2 * e.i; n != 0; {
		n1 := n
		n2 := len(spaces)

		if n1 > n2 {
			n1 = n2
		}

		if _, err = e.w.Write(spaces[:n1]); err != nil {
			return
		}

		n -= n1
	}

	return
}

func (e *PrettyEmitter) push(n int) int {
	if n != 0 {
		e.i++
	}
	e.s = append(e.s, n)
	return n
}

func (e *PrettyEmitter) pop() int {
	i := len(e.s) - 1
	n := e.s[i]
	e.s = e.s[:i]
	if n != 0 {
		e.i--
	}
	return n
}
