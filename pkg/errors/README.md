# errors

A package that provides enriched errors. It uses the same conventions and extends the Go standard [errors](https://pkg.go.dev/errors) package.

## Install

```sh
go get -u github.com/aukilabs/go-tooling/pkg/errors
```

## Usage

### Create

```go
err := errors.New("error message")                      // With a message.
err := errors.Newf("error message with format: %v", 42) // With a formatted message.
```

### Enrich With A Custom Type

```go
err := errors.New("handling http request failed").WithType("httpError")
```

### Enrich With Tags

```go
err := errors.New("handing http request failed").
    WithTag("method", "GET").
    WithTag("path", "/cookies").
    WithTag("code", 401)
```

### Wrap Another Error

```go
err := errors.New("handling http request failed").Wrap(fmt.Errorf("a fake simple http error"))
```

### Compose Multiple Enrichments

```go
err := errors.New("handing http request failed").
    WithType("httpError").
    WithTag("method", "GET").
    WithTag("path", "/cookies").
    WithTag("code", 401).
    Wrap(fmt.Errorf("a fake simple http error"))
```

### Get Error Type

```go
var err error

t := errors.Type(err)
```

### Get Error Tag

```go
var err error

foo := errors.Tag(err, "foo")
```
