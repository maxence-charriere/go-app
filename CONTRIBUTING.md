# Contributing

## Code

A PR would be awesome!
Here is some guidelines when submitting code:

- use gofmt
- use govet
- use golint
- test coverage for ```./internal/dom``` stays at 100%
- try to keep consistent coding style
- avoid naked returns (if you deal with a part of the code that have some, please refactor).
- run [goreportcard](https://goreportcard.com/report/) with your branch, everything must be 100%.
- when resolving an issue, try write a simple example that show how to use a feature.
- please, **make a PR on the `stage` branch**.
