package conf

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/segmentio/objconv"
)

// PrintError outputs the error message for err to stderr.
func (ld Loader) PrintError(err error) {
	w := bufio.NewWriter(os.Stderr)
	ld.fprintError(w, err, stderr())
	w.Flush()
}

// FprintError outputs the error message for err to w.
func (ld Loader) FprintError(w io.Writer, err error) {
	ld.fprintError(w, err, monochrome())
}

// PrintHelp outputs the help message for cfg to stderr.
func (ld Loader) PrintHelp(cfg interface{}) {
	w := bufio.NewWriter(os.Stderr)
	ld.fprintHelp(w, cfg, stderr())
	w.Flush()
}

// FprintHelp outputs the help message for cfg to w.
func (ld Loader) FprintHelp(w io.Writer, cfg interface{}) {
	ld.fprintHelp(w, cfg, monochrome())
}

func (ld Loader) fprintError(w io.Writer, err error, col colors) {
	var errors errorList

	if e, ok := err.(errorList); ok {
		errors = e
	} else {
		errors = errorList{err}
	}

	fmt.Fprintf(w, "%s\n", col.titles("Error:"))

	for _, e := range errors {
		fmt.Fprintf(w, "  %s\n", col.errors(e.Error()))
	}

	fmt.Fprintln(w)
}

func (ld Loader) fprintHelp(w io.Writer, cfg interface{}, col colors) {
	var m Map

	if cfg != nil {
		v := reflect.ValueOf(cfg)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		m = makeNodeStruct(v, v.Type())
	}

	fmt.Fprintf(w, "%s\n", col.titles("Usage:"))
	switch {
	case len(ld.Usage) != 0:
		fmt.Fprintf(w, "  %s %s\n\n", ld.Name, ld.Usage)
	case len(ld.Commands) != 0:
		fmt.Fprintf(w, "  %s [command] [options...]\n\n", ld.Name)
	default:
		fmt.Fprintf(w, "  %s [-h] [-help] [options...]\n\n", ld.Name)
	}

	if len(ld.Commands) != 0 {
		fmt.Fprintf(w, "%s\n", col.titles("Commands:"))
		width := 0

		for _, c := range ld.Commands {
			if n := len(col.cmds(c.Name)); n > width {
				width = n
			}
		}

		cmdfmt := fmt.Sprintf("  %%-%ds  %%s\n", width)

		for _, c := range ld.Commands {
			fmt.Fprintf(w, cmdfmt, col.cmds(c.Name), c.Help)
		}

		fmt.Fprintln(w)
	}

	set := newFlagSet(m, ld.Name, ld.Sources...)
	if m.Len() != 0 {
		fmt.Fprintf(w, "%s\n", col.titles("Options:"))
	}

	// Outputs the flags following the same format than the standard flag
	// package. The main difference is in the type names which are set to
	// values returned by prettyType.
	set.VisitAll(func(f *flag.Flag) {
		var t string
		var h []string
		var empty bool
		var boolean bool
		var object bool
		var list bool

		switch v := f.Value.(type) {
		case Node:
			x := reflect.ValueOf(v.Value())
			t = prettyType(x.Type())
			empty = isEmptyValue(x)

			switch v.(type) {
			case Map:
				object = true
			case Array:
				list = true
			default:
				boolean = isBoolFlag(x)
			}

		case FlagSource:
			t = "source"
		default:
			t = "value"
		}

		fmt.Fprintf(w, "  %s", col.keys("-"+f.Name))

		switch {
		case !boolean:
			fmt.Fprintf(w, " %s\n", col.types(t))
		case len(f.Name) >= 4: // put help message inline for boolean flags
			fmt.Fprint(w, "\n")
		}

		if s := f.Usage; len(s) != 0 {
			h = append(h, s)
		}

		if s := f.DefValue; len(s) != 0 && !empty && !(boolean || object || list) {
			h = append(h, col.defvals("(default "+s+")"))
		}

		if len(h) != 0 {
			if !boolean || len(f.Name) >= 4 {
				fmt.Fprint(w, "    ")
			}
			fmt.Fprintf(w, "\t%s\n", strings.Join(h, " "))
		}

		fmt.Fprint(w, "\n")
	})
}

func prettyType(t reflect.Type) string {
	if t == nil {
		return "unknown"
	}

	if _, ok := objconv.AdapterOf(t); ok {
		return "value"
	}

	switch {
	case t.Implements(objconvValueDecoderInterface):
		return "value"
	case t.Implements(textUnmarshalerInterface):
		return "string"
	}

	switch t {
	case timeDurationType:
		return "duration"
	case timeTimeType:
		return "time"
	}

	switch t.Kind() {
	case reflect.Struct, reflect.Map:
		return "object"
	case reflect.Slice, reflect.Array:
		if t.Elem().Kind() == reflect.Uint8 {
			return "base64"
		}
		return "list"
	case reflect.Ptr:
		return prettyType(t.Elem())
	default:
		return t.String()
	}
}

type colors struct {
	titles  func(string) string
	cmds    func(string) string
	keys    func(string) string
	types   func(string) string
	defvals func(string) string
	errors  func(string) string
}

func stderr() colors {
	if isTerminal(2) {
		return colorized()
	}
	return monochrome()
}

func colorized() colors {
	return colors{
		titles:  bold,
		cmds:    magenta,
		keys:    blue,
		types:   green,
		defvals: grey,
		errors:  red,
	}
}

func monochrome() colors {
	return colors{
		titles:  normal,
		cmds:    normal,
		keys:    normal,
		types:   normal,
		defvals: normal,
		errors:  normal,
	}
}

func bold(s string) string {
	return "\033[1m" + s + "\033[0m"
}

func blue(s string) string {
	return "\033[1;34m" + s + "\033[0m"
}

func green(s string) string {
	return "\033[1;32m" + s + "\033[0m"
}

func red(s string) string {
	return "\033[1;31m" + s + "\033[0m"
}

func magenta(s string) string {
	return "\033[1;35m" + s + "\033[0m"
}

func grey(s string) string {
	return "\033[1;30m" + s + "\033[0m"
}

func normal(s string) string {
	return s
}

func isEmptyValue(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}

	switch v.Kind() {
	case reflect.Slice, reflect.Map:
		return v.Len() == 0

	case reflect.Struct:
		return v.NumField() == 0
	}

	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

func isBoolFlag(v reflect.Value) bool {
	type iface interface {
		IsBoolFlag() bool
	}

	if !v.IsValid() {
		return false
	}

	if x, ok := v.Interface().(iface); ok {
		return x.IsBoolFlag()
	}

	return v.Kind() == reflect.Bool
}
