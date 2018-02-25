# conf [![CircleCI](https://circleci.com/gh/segmentio/conf.svg?style=shield)](https://circleci.com/gh/segmentio/conf) [![Go Report Card](https://goreportcard.com/badge/github.com/segmentio/conf)](https://goreportcard.com/report/github.com/segmentio/conf) [![GoDoc](https://godoc.org/github.com/segmentio/conf?status.svg)](https://godoc.org/github.com/segmentio/conf)
Go package for loading program configuration from multiple sources.

Motivations
-----------

Loading program configurations is usually done by parsing the arguments passed
to the command line, and in this case the standard library offers a good support
with the `flag` package.  
However, there are times where the standard is just too limiting, for example
when the program needs to load configuration from other sources (like a file, or
the environment variables).  
The `conf` package was built to address these issues, here were the goals:

- **Loading the configuration has to be type-safe**, there were other packages
available that were covering the same use-cases but they often required doing
type assertions on the configuration values which is always an opportunity to
get the program to panic.

- **Keeping the API minimal**, while the `flag` package offered the type safety
we needed it is also very verbose to setup. With `conf`, only a single function
call is needed to setup and load the entire program configuration.

- **Supporting richer syntaxes**, because program configurations are often
generated dynamically, the `conf` package accepts YAML values as input to all
configuration values. It also has support for sub-commands on the command line,
which is a common approach used by CLI tools.

- **Supporting multiple sources**, because passing values through the command
line is not always the best appraoch, programs may need to receive their
configuration from files, environment variables, secret stores, or other network
locations.

Basic Usage
-----------

A program using the `conf` package needs to declare a struct which is passed to
`conf.Load` to populate the fields with the configuration that was made
available at runtime through a configuration file, environment variables or the
program arguments.

Each field of the structure may declare a `conf` tag which sets the name of the
property, and a `help` tag to provide a help message for the configuration.

The `conf` package will automatically understand the structure of the program
configuration based on the struct it receives, as well as generating the program
usage and help messages if the `-h` or `-help` options are passed (or an error
is detected).

The `conf.Load` function adds support for a `-config-file` option on the program
arguments which accepts the path to a file that the configuration may be loaded
from as well.

Here's an example of how a program would typically use the package:
```go
package main

import (
    "fmt"

    "github.com/segmentio/conf"
)

func main() {
    var config struct {
        Message string `conf:"m" help:"A message to print."`
    }

    // Load the configuration, either from a config file, the environment or the program arguments.
    conf.Load(&config)

    fmt.Println(config.Message)
}
```
```
$ go run ./example.go -m 'Hello World!'
Hello World!
```

Advanced Usage
--------------

While the `conf.Load` function is good enough for common use cases, programs
sometimes need to customize the default behavior.  
A program may then use the `conf.LoadWith` function, which accepts a
`conf.Loader` as second argument to gain more control over how the configuration
is loaded.

Here's the `conf.Loader` definition:
```go
package conf

type Loader struct {
     Name     string    // program name
     Usage    string    // program usage
     Args     []string  // list of arguments
     Commands []Command // list of commands
     Sources  []Source  // list of sources to load configuration from.
}
```

The `conf.Load` function is actually just a wrapper around `conf.LoadWith` that
passes a default loader. The default loader gets the program name from the first
program argument, supports no sub-commands, and has two custom sources setup to
potentially load its configuration from a configuration file or the environment
variables.

Here's an example showing how to configure a CLI tool that supports a couple of
sub-commands:
```go
package main

import (
    "fmt"

    "github.com/segmentio/conf"
)

func main() {
    // If nil is passed instead of a configuration struct no arguments are
    // parsed, only the command is extracted.
    cmd, args := conf.LoadWith(nil, conf.Loader{
        Name:     "example",
        Args:     os.Args[1:],
        Commands: []conf.Command{
            {"print", "Print the message passed to -m"},
            {"version", "Show the program version"},
        },
    })

    switch cmd {
    case "print":
        var config struct{
            Message string `conf:"m" help:"A message to print."`
        }

        conf.LoadWith(&config, conf.Loader{
            Name: "example print",
            Args: args,
        })

        fmt.Println(config.Message)

    case "version":
        fmt.Println("1.2.3")
    }
}
```
```
$ go run ./example.go version
1.2.3
$ go run ./example.go print -m 'Hello World!'
Hello World!
```

Custom Sources
--------------

We mentionned the `conf.Loader` type supported setting custom sources that the
program configuration can be loaded from. Here's the the `conf.Source` interface
definition:
```go
package conf

type Source interface {
    Load(dst Map)
}
```

The source has a single method which receives a `conf.Map` value which is an
itermediate representation of the configuration struct that was received by the
loader.  
The package uses this type internally as well for loading configuration values
from the program arguments, it can be seen as a reflective representiong of the
original value which exposes an API that is more convenient to use that having
a raw `reflect.Value`.

One of the advantages of the `conf.Map` type is that it implements the
[objconv.ValueDecoder](https://godoc.org/github.com/segmentio/objconv#ValueDecoder)
interface and therefore can be used directly to load configurations from a
serialized format (like JSON for example).

Validation
----------

Last but not least, the `conf` package also supports automatic validation of the
fields in the configuration struct. This happens after the values were loaded
and is based on [gopkg.in/validator.v2](https://godoc.org/gopkg.in/validator.v2).

This step could have been done outside the package however it is both convenient
and useful to have all configuration errors treated the same way (getting the
usage and help message shown when something is wrong).
