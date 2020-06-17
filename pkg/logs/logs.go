// Package logs implements functions to manipulate logs.
//
// Logs created are taggable.
//
//   logWithTags := logs.New("a log with tags").
//       Tag("a", 42).
// 	     Tag("b", 21)
package logs

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unsafe"
)

// New returns a log with the given description that can be tagged.
func New(v string) Log {
	return Log{
		description: v,
	}
}

// Newf returns a log with the given formatted description that can be tagged.
func Newf(format string, v ...interface{}) Log {
	return New(fmt.Sprintf(format, v...))
}

// Log is a implementation that supports tagging.
type Log struct {
	description string
	tags        []tag
	maxKeyLen   int
}

// Tag sets the named tag with the given value.
func (l Log) Tag(k string, v interface{}) Log {
	if l.tags == nil {
		l.tags = make([]tag, 0, 8)
	}

	if length := len(k); length > l.maxKeyLen {
		l.maxKeyLen = length
	}

	switch v := v.(type) {
	case string:
		l.tags = append(l.tags, tag{key: k, value: v})
	case fmt.Stringer:
		l.tags = append(l.tags, tag{key: k, value: v.String()})
	case fmt.GoStringer:
		l.tags = append(l.tags, tag{key: k, value: v.GoString()})

	case int:
		l.tags = append(l.tags, tag{key: k, value: strconv.Itoa(v)})
	case int8:
		l.tags = append(l.tags, tag{key: k, value: strconv.FormatInt(int64(v), 10)})
	case int16:
		l.tags = append(l.tags, tag{key: k, value: strconv.FormatInt(int64(v), 10)})
	case int32:
		l.tags = append(l.tags, tag{key: k, value: strconv.FormatInt(int64(v), 10)})
	case int64:
		l.tags = append(l.tags, tag{key: k, value: strconv.FormatInt(int64(v), 10)})

	case uint:
		l.tags = append(l.tags, tag{key: k, value: strconv.FormatUint(uint64(v), 10)})
	case uint8:
		l.tags = append(l.tags, tag{key: k, value: strconv.FormatUint(uint64(v), 10)})
	case uint16:
		l.tags = append(l.tags, tag{key: k, value: strconv.FormatUint(uint64(v), 10)})
	case uint32:
		l.tags = append(l.tags, tag{key: k, value: strconv.FormatUint(uint64(v), 10)})
	case uint64:
		l.tags = append(l.tags, tag{key: k, value: strconv.FormatUint(uint64(v), 10)})

	case float32:
		l.tags = append(l.tags, tag{key: k, value: strconv.FormatFloat(float64(v), 'f', 4, 32)})
	case float64:
		l.tags = append(l.tags, tag{key: k, value: strconv.FormatFloat(v, 'f', 4, 64)})

	default:
		l.tags = append(l.tags, tag{key: k, value: fmt.Sprintf("%+v", v)})
	}

	return l
}

func (l Log) String() string {
	w := bytes.NewBuffer(make([]byte, 0, len(l.description)+len(l.tags)*(l.maxKeyLen+11)))
	l.format(w, 0)
	return bytesToString(w.Bytes())
}

func (l Log) format(w *bytes.Buffer, indent int) {
	w.WriteString(l.description)
	if len(l.tags) != 0 {
		w.WriteByte(':')
	}

	tags := l.tags
	sort.Slice(tags, func(a, b int) bool {
		return strings.Compare(tags[a].key, tags[b].key) < 0
	})

	for _, t := range l.tags {
		k := t.key
		v := t.value

		w.WriteByte('\n')
		l.indent(w, indent+4)
		w.WriteString(k)
		w.WriteByte(':')
		l.indent(w, l.maxKeyLen-len(k)+1)
		w.WriteString(v)
	}
}

func (l Log) indent(w *bytes.Buffer, n int) {
	for i := 0; i < n; i++ {
		w.WriteByte(' ')
	}
}

type tag struct {
	key   string
	value string
}

func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
