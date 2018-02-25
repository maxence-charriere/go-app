package conf

import (
	"bytes"
	"flag"
	"strings"
	"text/template"

	"github.com/segmentio/objconv/json"
)

// Source is the interface that allow new types to be plugged into a loader to
// make it possible to load configuration from new places.
//
// When the configuration is loaded the Load method of each source that was set
// on a loader is called with an Node representating the configuration struct.
// The typical implementation of a source is to load the serialized version of
// the configuration and use an objconv decoder to build the node.
type Source interface {
	Load(dst Map) error
}

// FlagSource is a special case of a source that receives a configuration value
// from the arguments of a loader. It makes it possible to provide runtime
// configuration to the source from the command line arguments of a program.
type FlagSource interface {
	Source

	// Flag is the name of the flag that sets the source's configuration value.
	Flag() string

	// Help is called to get the help message to display for the source's flag.
	Help() string

	// flag.Value must be implemented by a FlagSource to receive their value
	// when the loader's arguments are parsed.
	flag.Value
}

// SourceFunc makes it possible to use basic function types as configuration
// sources.
type SourceFunc func(dst Map) error

// Load calls f.
func (f SourceFunc) Load(dst Map) error {
	return f(dst)
}

// NewEnvSource creates a new source which loads values from the environment
// variables given in env.
//
// A prefix may be set to namespace the environment variables that the source
// will be looking at.
func NewEnvSource(prefix string, env ...string) Source {
	vars := makeEnvVars(env)
	base := make([]string, 0, 10)

	if prefix != "" {
		base = append(base, prefix)
	}

	return SourceFunc(func(dst Map) (err error) {
		dst.Scan(func(path []string, item MapItem) {
			path = append(base, path...)
			path = append(path, item.Name)

			k := snakecaseUpper(strings.Join(path, "_"))

			if v, ok := vars[k]; ok {
				if e := item.Value.Set(v); e != nil {
					err = e
				}
			}
		})
		return
	})
}

// NewFileSource creates a new source which loads a configuration from a file
// identified by a path (or URL).
//
// The returned source satisfies the FlagSource interface because it loads the
// file location from the given flag.
//
// The vars argument may be set to render the configuration file if it's a
// template.
//
// The readFile function loads the file content in-memory from a file location
// given as argument, usually this is ioutil.ReadFile.
//
// The unmarshal function decodes the content of the configuration file into a
// configuration object.
func NewFileSource(flag string, vars interface{}, readFile func(string) ([]byte, error), unmarshal func([]byte, interface{}) error) FlagSource {
	return &fileSource{
		flag:      flag,
		vars:      vars,
		readFile:  readFile,
		unmarshal: unmarshal,
	}
}

type fileSource struct {
	flag      string
	path      string
	vars      interface{}
	readFile  func(string) ([]byte, error)
	unmarshal func([]byte, interface{}) error
}

func (f *fileSource) Load(dst Map) (err error) {
	var b []byte

	if len(f.path) == 0 {
		return
	}

	if b, err = f.readFile(f.path); err != nil {
		return
	}

	tpl := template.New(f.flag)
	buf := &bytes.Buffer{}
	buf.Grow(len(b))

	tpl = tpl.Funcs(template.FuncMap{
		"json": func(v interface{}) (string, error) {
			b, err := json.Marshal(v)
			return string(b), err
		},
	})

	if _, err = tpl.Parse(string(b)); err != nil {
		return
	}

	if err = tpl.Execute(buf, f.vars); err != nil {
		return
	}

	err = f.unmarshal(buf.Bytes(), dst)
	return
}

func (f *fileSource) Flag() string {
	return f.flag
}

func (f *fileSource) Help() string {
	return "Location to load the configuration file from."
}

func (f *fileSource) Set(s string) error {
	f.path = s
	return nil
}

func (f *fileSource) String() string {
	return f.path
}
