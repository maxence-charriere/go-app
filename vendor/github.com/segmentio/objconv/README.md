objconv [![CircleCI](https://circleci.com/gh/segmentio/objconv.svg?style=shield)](https://circleci.com/gh/segmentio/objconv) [![Go Report Card](https://goreportcard.com/badge/github.com/segmentio/objconv)](https://goreportcard.com/report/github.com/segmentio/objconv) [![GoDoc](https://godoc.org/github.com/segmentio/objconv?status.svg)](https://godoc.org/github.com/segmentio/objconv)
=======

This Go package provides the implementation of high performance encoder and
decoders for JSON-like object representations.

The top-level package exposes the generic types and algorithms for encoding and
decoding values, while each sub-package implements the parser and emitters for
specific types.

### Breaking changes introduced in [#18](https://github.com/segmentio/objconv/pull/18)

The `Encoder` type used to have methods exposed to encode specific types for
optimization purposes. The generic `Encode` method has been optimized to make
those other methods obsolete and they were therefore removed.

Compatibility with the standard library
---------------------------------------

The sub-packages providing implementation for specific formats also expose APIs
that mirror those of the standard library to make it easy to integrate with the
objconv package. However there are a couple of differences that need to be taken
in consideration:

- Encoder and Decoder types are not exposed in the objconv sub-packages, instead
the types from the top-level package are used. For example, variables declared
with the `json.Encoder` type would have to be replaced with `objconv.Encoder`.

- Interfaces like `json.Marshaler` or `json.Unmarshaler` are not supported.
However the `encoding.TextMarshaler` and `encoding.TextUnmarshaler` interfaces
are.

Encoder
-------

The package exposes a generic encoder API that let's the program serialize
native values into various formats.

Here's an example of how to serialize a structure to JSON:
```go
package main

import (
    "os"

    "github.com/segmentio/objconv/json"
)

func main() {
    e := json.NewEncoder(os.Stdout)
    e.Encode(struct{ Hello string }{"World"})
}
```
```
$ go run ./example.go
{"Hello":"World"}
```

Note that this code is fully compatible with the standard `encoding/json`
package.

Decoder
-------

Here's an example of how to use a JSON decoder:
```go
package main

import (
    "fmt"
    "os"

    "github.com/segmentio/objconv/json"
)

func main() {
    v := struct{ Message string }{}

    d := json.NewDecoder(os.Stdin)
    d.Decode(&v)

    fmt.Println(v.Message)
}
```
```
$ echo '{ "Message": "Hello World!" }' | go run ./example.go
Hello World!
```

Streaming
---------

One of the interesting features of the `objconv` package is the ability to read
and write streams of data. This has several advantages in terms of memory usage
and latency when passing data from service to service.  
The package exposes the `StreamEncoder` and `StreamDecoder` types for this
purpose.

For example the JSON stream encoder and decoder can produce a JSON array as a
stream where data are produced and consumed on the fly as they become available,
here's an example:
```go
package main

import (
    "io"

    "github.com/segmentio/objconv/json"
)

func main() {
     r, w := io.Pipe()

    go func() {
        defer w.Close()

        e := json.NewStreamEncoder(w)
        defer e.Close()

        // Produce values to the JSON stream.
        for i := 0; i != 1000; i++ {
            e.Encode(i)
        }
    }()

    d := json.NewStreamDecoder(r)

    // Consume values from the JSON stream.
    var v int

    for d.Decode(&v) == nil {
        // v => {0..999}
        // ...
    }
}
```

Stream decoders are capable of reading values from either arrays or single
values, this is very convenient when an program cannot predict the structure of
the stream. If the actual data representation is not an array the stream decoder
will simply behave like a normal decoder and produce a single value.

Encoding and decoding custom types
----------------------------------

To override the default encoder and decoder behaviors a type may implement the
`ValueEncoder` or `ValueDecoder` interface. The method on these interfaces are
called to customize the default behavior.

This can prove very useful to represent slice of pairs as maps for example:
```go
type KV struct {
    K string
    V interface{}
}

type M []KV

// Implement the ValueEncoder interface to provide a custom encoding.
func (m M) EncodeValue(e objconv.Encoder) error {
    i := 0
    return e.EncodeMap(len(m), func(k objconv.Encoder, v objconv.Encoder) (err error) {
        if err = k.Encode(m[i].K); err != nil {
            return
        }
        if err = v.Encode(m[i].V); err != nil {
            return
        }
        i++
        return
    })
}
```

Mime Types
----------

The `objconv` package exposes APIs for registering codecs for specific mime
types. When an objconv package for a specific format is imported
it registers itself on the global registry to be later referred by name.

```go
import (
    "bytes"

    "github.com/segmentio/objconv"
    _ "github.com/segmentio/objconv/json" // registers the JSON codec
)

func main() {
    // Lookup the JSON codec.
    jsonCodec, ok := objconv.Lookup("application/json")

    if !ok {
        panic("unreachable")
    }

    // Create a new encoder from the codec.
    b := &bytes.Buffer{}
    e := jsonCodec.NewEncoder(b)

    // ...
}
```
