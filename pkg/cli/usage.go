package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
)

const (
	defaultColor = "\033[0m"
	errorColor   = "\033[91m"
	successColor = "\033[92m"
	accentColor  = "\033[94m"
	focusColor   = "\033[1m"
	subColor     = "\033[2m"
)

func commandUsage(w io.Writer, cmd *command, opts []option) func() {
	return func() {
		// Usage:
		fmt.Fprintf(w, "%sUsage:%s\n\n", accentColor, defaultColor)

		indent(w, 4)
		fmt.Fprint(w, focusColor, filepath.Base(os.Args[0]))
		if cmd.name != "" {
			fmt.Fprint(w, " ", cmd.name)
		}
		fmt.Fprint(w, defaultColor, accentColor, " [options]", defaultColor)
		fmt.Fprintln(w)
		fmt.Fprintln(w)

		// Help:
		if cmd.help != "" {
			fmt.Fprintf(w, "%sDescription:%s\n\n", accentColor, defaultColor)
			indent(w, 4)
			writeText(w, cmd.help, 4, 80)
			fmt.Fprintln(w)
		}

		// Options:
		if len(opts) == 0 {
			return
		}

		fmt.Fprintf(w, "%sOptions:%s\n\n", accentColor, defaultColor)

		optsInfo := optionsInfo(opts)
		for _, o := range opts {
			indent(w, 4)
			fmt.Fprintf(w, "%s-%s%s", focusColor, o.name, defaultColor)
			indent(w, optsInfo.nameLen-len(o.name)+2)

			typeName := o.value.Type().String()
			typeName = strings.TrimPrefix(typeName, "main.")
			fmt.Fprintf(w, "%s%s%s", accentColor, typeName, defaultColor)
			indent(w, optsInfo.typeLen-len(typeName)+4)

			lastColIndent := 4 + 1 + optsInfo.nameLen + 2 + optsInfo.typeLen + 4
			if o.help != "" {
				writeText(w, o.help, lastColIndent, 80)
				indent(w, lastColIndent)
			}

			if o.envKey != "-" {
				fmt.Fprintf(w, "%sEnv:%s     %s%s%s\n", subColor, defaultColor, accentColor, o.envKey, defaultColor)
				indent(w, lastColIndent)
			}

			if !o.value.IsZero() {
				switch o.value.Kind() {
				case reflect.String,
					reflect.Struct,
					reflect.Map,
					reflect.Array,
					reflect.Slice:
					fmt.Fprintf(w, "%sDefault:%s %q\n", subColor, defaultColor, o)

				default:
					fmt.Fprintf(w, "%sDefault:%s %s\n", subColor, defaultColor, o)

				}
			}

			fmt.Fprintln(w)
		}
	}
}

func commandUsageIndex(w io.Writer, cmds map[string]Command) func() {
	return func() {
		// Usage:
		fmt.Fprintf(w, "%sUsage:%s\n\n", accentColor, defaultColor)

		indent(w, 4)
		fmt.Fprint(w, focusColor, filepath.Base(os.Args[0]))
		fmt.Fprint(w, defaultColor, accentColor, " <command>", defaultColor)
		fmt.Fprintln(w)
		fmt.Fprintln(w)

		// Commands:
		fmt.Fprintf(w, "%sCommands:%s\n\n", accentColor, defaultColor)

		cmdlist := make([]*command, 0, len(cmds))
		maxLenName := 0
		for _, c := range cmds {
			cmd := c.(*command)

			if l := len(cmd.name); l > maxLenName {
				maxLenName = l
			}

			cmdlist = append(cmdlist, cmd)
		}

		sort.Slice(cmdlist, func(i, j int) bool {
			return strings.Compare(cmdlist[i].name, cmdlist[j].name) < 0
		})

		for _, c := range cmdlist {
			indent(w, 4)
			fmt.Fprintf(w, "%s%s%s", focusColor, c.name, defaultColor)
			indent(w, maxLenName-len(c.name)+4)
			writeText(w, c.help, 4+maxLenName+4, 80)
			fmt.Fprintln(w)
		}
	}
}

func printError(w io.Writer, err error) {
	fmt.Fprintf(w, "%sError:%s\n\n", errorColor, defaultColor)
	indent(w, 4)
	fmt.Fprintln(w, err)
	fmt.Fprintln(w)
}

func indent(w io.Writer, level int) int {
	count := 0
	for i := 0; i < level; i++ {
		n, _ := w.Write([]byte(" "))
		count += n
	}
	return count
}

func writeText(w io.Writer, text string, level, maxLen int) {
	buff := []byte(text)
	var trailingWord []byte
	count := level

	defer fmt.Fprintln(w)

	for {
		if len(trailingWord) != 0 {
			count += indent(w, level)
			written, _ := w.Write(trailingWord)
			count += written
			trailingWord = nil
			continue
		}

		advance, word, err := bufio.ScanWords(buff, true)
		if err != nil || advance == 0 {
			return
		}
		buff = buff[advance:]

		if count+len(word) > maxLen {
			count = 0
			trailingWord = word
			fmt.Fprintln(w)
			continue
		}

		if count > level {
			count += indent(w, 1)
		}

		written, _ := w.Write(word)
		count += written
	}
}

type optionFormatInfo struct {
	nameLen int
	typeLen int
}

func optionsInfo(opts []option) optionFormatInfo {
	var info optionFormatInfo
	for _, o := range opts {
		if l := len(o.name); l > info.nameLen {
			info.nameLen = l
		}

		if l := len(o.value.Type().String()); l > info.typeLen {
			info.typeLen = l
		}
	}
	return info
}
