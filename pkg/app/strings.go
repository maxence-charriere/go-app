package app

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/maxence-charriere/go-app/v10/pkg/errors"
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

func stringTo(s string, v any) error {
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Ptr {
		return errors.New("receiver in not a pointer").WithTag("receiver-type", val.Type())
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
			WithTag("string", s).
			WithTag("receiver-type", val.Type())
	}

	return nil
}

// AppendClass adds values to the given class string.
func AppendClass(class string, v ...string) string {
	var b strings.Builder
	b.WriteString(strings.TrimSpace(class))

	for _, c := range v {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}

		if b.Len() != 0 {
			b.WriteByte(' ')
		}
		b.WriteString(c)
	}

	return b.String()
}

func jsonString(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(errors.New("converting value to json string failed").Wrap(err))
	}
	return string(b)
}

// Formats a string with the given format and values.
// It uses fmt.Sprintf when len(v) != 0.
//
// TODO: Write a faster Sprintf to use with values.
func FormatString(format string, v ...any) string {
	if len(v) == 0 {
		return format
	}
	return fmt.Sprintf(format, v...)
}

func previewText(v string) string {
	if len(v) <= 80 {
		return v
	}
	return v[:77] + "..."
}
