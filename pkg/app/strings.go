package app

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
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

	case bool:
		return strconv.FormatBool(v)

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

func stringTo(s string, v interface{}) error {
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Ptr {
		return errors.New("receiver in not a pointer").Tag("receiver-type", val.Type())
	}
	val = val.Elem()

	switch val.Kind() {
	case reflect.String:
		val.SetString(s)

	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		i, _ := strconv.ParseInt(s, 10, 0)
		val.SetInt(i)

	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		i, _ := strconv.ParseUint(s, 10, 0)
		val.SetUint(i)

	case reflect.Float64:
		f, _ := strconv.ParseFloat(s, 64)
		val.SetFloat(f)

	case reflect.Float32:
		f, _ := strconv.ParseFloat(s, 32)
		val.SetFloat(f)

	default:
		return errors.New("string cannot be converted to receiver type").
			Tag("string", s).
			Tag("receiver-type", val.Type())
	}

	return nil
}

// AppendClass adds c to the given class string.
func AppendClass(class, c string) string {
	if c == "" {
		return class
	}
	if class != "" {
		class += " "
	}
	class += c
	return class
}
