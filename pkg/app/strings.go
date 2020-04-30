package app

import (
	"fmt"
	"io"
	"unsafe"
)

func writeIndent(w io.Writer, indent int) {
	for i := 0; i < indent*4; i++ {
		w.Write(stob(" "))
	}
}

func stob(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func btos(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func ln() []byte {
	return stob("\n")
}

func toString(v interface{}) string {
	switch v := v.(type) {
	case string:
		return v

	case []byte:
		return btos(v)

	default:
		return fmt.Sprint(v)
	}
}
