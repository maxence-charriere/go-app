package cli

import (
	"encoding/json"
	"flag"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

type optionParser struct {
	flags   *flag.FlagSet
	options []option
}

func (p *optionParser) parse(v interface{}) ([]option, error) {
	p.options = nil

	if v == nil {
		return nil, nil
	}

	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Ptr {
		return nil, errors.New("receiver is not a pointer").
			WithTag("type", val.Type()).
			WithTag("kind", val.Kind())
	}

	if val = val.Elem(); val.Kind() != reflect.Struct {
		return nil, errors.New("receiver does not point to a struct").
			WithTag("type", val.Type()).
			WithTag("kind", val.Kind())
	}

	p.parseStruct("", val)

	for _, o := range p.options {
		if o.name != "h" && o.name != "help" {
			p.flags.Var(o, o.name, o.help)
		}
	}

	return p.options, nil
}

func (p *optionParser) parseStruct(prefix string, v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		fval := v.Field(i)
		if !fval.CanSet() {
			continue
		}

		finfo := v.Type().Field(i)
		fname := finfo.Tag.Get("cli")

		if fname == "" {
			fname = finfo.Name
		}
		fname = normalizeCLIOptionName(fname)

		if prefix != "" {
			fname = prefix + "." + fname
		}

		envKey := finfo.Tag.Get("env")
		if envKey == "" {
			envKey = normalizeEnvOptionName(fname)
		}

		o := option{
			name:   fname,
			help:   finfo.Tag.Get("help"),
			envKey: envKey,
			value:  fval,
		}

		if envVal, ok := os.LookupEnv(envKey); ok && envKey != "-" {
			o.Set(envVal)
		}

		p.options = append(p.options, o)

		if fval.Kind() == reflect.Struct {
			p.parseStruct(fname, fval)
		}
	}
}

type option struct {
	name   string
	help   string
	envKey string
	value  reflect.Value
}

func (o option) IsBoolFlag() bool {
	return o.value.Kind() == reflect.Bool
}

func (o option) String() string {
	switch o.value.Kind() {
	case reflect.String:
		return o.value.String()
	}

	switch value := o.value.Interface().(type) {
	case time.Duration:
		return value.String()
	}

	b, _ := json.Marshal(o.value.Interface())
	return string(b)
}

func (o option) Set(s string) error {
	switch o.value.Kind() {
	case reflect.String:
		o.value.SetString(s)
		return nil
	}

	switch o.value.Interface().(type) {
	case time.Duration:
		return setDuration(o.value, s)

	case time.Time, *time.Time:
		s = strconv.Quote(s)
	}

	return json.Unmarshal([]byte(s), o.value.Addr().Interface())
}

func setDuration(v reflect.Value, s string) error {
	d, err := time.ParseDuration(s)
	if err != nil {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}

		v.SetInt(i)
		return nil
	}

	v.SetInt(int64(d))
	return nil
}

func normalizeOptionName(name string, sep string) string {
	var b strings.Builder

	write := func(s string) {
		if s != "" {
			if b.Len() != 0 {
				b.WriteString(sep)
			}
			b.WriteString(strings.ToLower(s))
		}
	}

	start := 0
	end := 0

	for end < len(name) {
		switch {
		case isUpperCase(name[end]) && end > 0 && !isUpperCase(name[end-1]):
			write(name[start:end])
			start = end

		case name[end] == '-',
			name[end] == '_',
			name[end] == ' ',
			name[end] == '\t',
			name[end] == '.':
			write(name[start:end])
			start = end + 1
		}

		end++
	}

	write(name[start:end])
	return b.String()
}

func normalizeCLIOptionName(name string) string {
	return normalizeOptionName(name, "-")
}

func normalizeEnvOptionName(name string) string {
	name = normalizeOptionName(name, "_")
	return strings.ToUpper(name)
}

func isUpperCase(b byte) bool {
	return b >= 'A' && b <= 'Z'
}
