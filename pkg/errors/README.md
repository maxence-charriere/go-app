# errors

Package errors implements functions to manipulate errors.

Errors created are taggable and wrappable.

```go
errWithTags := errors.New("an error with tags").
    Tag("a", 42).
    Tag("b", 21)

errWithWrap := errors.New("error").
    Tag("a", 42).
    Wrap(errors.New("wrapped error"))
```

The package mirrors https:golang.org/pkg/errors package.
