# Contributing
Here is some basic guidelines:

- use gofmt
- use govet
- use golint
- test coverage for ```.``` stay at 100%
- test coverage for ```./appjs``` stays at 100%
- test coverage for ```./bridge``` stays at 100%
- test coverage for ```./html``` stays at 100%
- try to keep consistent coding style
- avoid naked returns (if you deal with a part of the code that have some, please refactor).
- run [goreportcard](https://goreportcard.com/report/) with your branch, everything that is not gocyclo must be 100%.
- when resolving an issue, try write a simple example that show how to use a feature.
