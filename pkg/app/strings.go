package app

import (
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

func ln() []byte {
	return stob("\n")
}
