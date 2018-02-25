package conf

import (
	"flag"
	"io/ioutil"
	"strings"
)

func newFlagSet(cfg Map, name string, sources ...Source) *flag.FlagSet {
	set := flag.NewFlagSet(name, flag.ContinueOnError)
	set.SetOutput(ioutil.Discard)

	cfg.Scan(func(path []string, item MapItem) {
		set.Var(item.Value, strings.Join(append(path, item.Name), "."), item.Help)
	})

	for _, source := range sources {
		if f, ok := source.(FlagSource); ok {
			set.Var(f, f.Flag(), f.Help())
		}
	}

	return set
}
