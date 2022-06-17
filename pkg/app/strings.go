package app

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

func toString(v any) string {
	switch v := v.(type) {
	case string:
		return v

	case []byte:
		return string(v)

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

func toPath(v ...any) string {
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
		io.WriteString(w, " ")
	}
}

func ln() []byte {
	return []byte("\n")
}

func pxToString(px int) string {
	return strconv.Itoa(px) + "px"
}

func stringTo(s string, v any) error {
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

func jsonString(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(errors.New("converting value to json string failed").Wrap(err))
	}
	return string(b)
}
