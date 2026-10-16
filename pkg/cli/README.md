# cli

Go package for building cli tool that can load program configuration from multiple sources:

- command flags
- Environment

Environment values override struct defaults, and command flags override valid
environment values. Invalid environment values stop loading before flags are
parsed, even when a flag would override them. The error identifies the option,
environment key, and expected type without printing the environment value.

## Usage

Define a Go struct that uses field tags to describe a command configuration:

```go
type config struct {
	String string `cli:"string" env:"CONF_STRING" help:"A string argmument for demo purpose."`
	Int    int    `cli:"int"    env:"CONF_INT"    help:"A int argmument for demo purpose."`
	Help   bool   `cli:"h"      env:"-"           help:"Show help."`
}
```

| Field Tag | Description                                                           |
| --------- | --------------------------------------------------------------------- |
| cli       | Maps a cli flag for the given field.                                  |
| env       | Maps environment variable for the given field.                        |
| help      | Setup a description for the given field when using the help flag `-h` |

Load the config:

```go
func main() {
	cfg := config{
		Int: 42, // Set 42 as default value for Int field.
	}

	// Registers the config.
	cli.Register().
		Help("A demo cli program").
		Options(&cfg)

	cli.Load() // Stores corresponding cli and env into cfg.

	fmt.Println(cfg)
}
```

command output with `-h` flag:

```
▶ ./my-program -h
Usage:

    hds [options]

Description:

    A demo cli program

Options:

    -string  string    A string argmument for demo purpose.
                       Env:     CONF_STRING

    -int     int       A int argmument for demo purpose.
                       Env:     CONF_INT
                       Default: 42

    -h       bool      Show help.
```
