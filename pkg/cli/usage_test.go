package cli

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

func TestCommandUsage(t *testing.T) {
	cmd := &command{
		name: "foo bar",
		help: `
			Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do
			eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim
			ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut
			aliquip ex ea commodo consequat. Duis aute irure dolor in
			reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla
			pariatur. Excepteur sint occaecat cupidatat non proident, sunt in
			culpa qui officia deserunt mollit anim id est laborum.
			`,
	}

	opts := []option{
		{
			name: "foo",
			help: `
				Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod
				tempor incididunt ut labore et dolore magna aliqua.
				`,
			envKey: "FOO",
			value:  reflect.ValueOf(42),
		},
		{
			name:   "bar",
			help:   "Bar option description.",
			envKey: "-",
			value:  reflect.ValueOf("bar"),
		},
		{
			name:   "alakazam",
			help:   "Alakazam option description.",
			envKey: "BAR",
			value:  reflect.ValueOf(0),
		},
	}

	w := bytes.NewBufferString("\n")
	usage := commandUsage(w, cmd, opts)
	usage()

	t.Log(w.String())
}

func TestCommandUsageIndex(t *testing.T) {
	m := commandManager{}

	m.register("foo", "bar").Help(`
	Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do
	eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim
	ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut
	aliquip ex ea commodo consequat. Duis aute irure dolor in
	reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla
	pariatur. Excepteur sint occaecat cupidatat non proident, sunt in
	culpa qui officia deserunt mollit anim id est laborum.
	`)
	m.register("foo", "foo").Help("Foo lolilol.")
	m.register("foo", "buu").Help("A more simple help.")

	w := bytes.NewBufferString("\n")
	usage := commandUsageIndex(w, m.commands)
	usage()

	t.Log(w.String())
}

func TestPrintError(t *testing.T) {
	w := bytes.NewBufferString("\n")
	printError(w, errors.New("an error for testing printing"))
	t.Log(w.String())
}
