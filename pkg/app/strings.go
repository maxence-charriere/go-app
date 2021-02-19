package app

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"unsafe"
)

func toString(v interface{}) string {
	switch v := v.(type) {
	case string:
		return v

	case []byte:
		return btos(v)

	case int:
		return strconv.Itoa(v)

	case float64:
		return strconv.FormatFloat(v, 'f', 4, 64)

	case nil:
		return ""

	default:
		return fmt.Sprint(v)
	}
}

func toPath(v ...interface{}) string {
	var b strings.Builder

	for _, o := range v {
		s := toString(o)
		if s == "" {
			continue
		}
		b.WriteByte('/')
		b.WriteString(s)
	}

	return b.String()
}

func writeIndent(w io.Writer, indent int) {
	for i := 0; i < indent*2; i++ {
		w.Write(stob(" "))
	}
}

func ln() []byte {
	return stob("\n")
}

func btos(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func stob(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func pxToString(px int) string {
	return strconv.Itoa(px) + "px"
}
