# Contributing

Contributions to go-app are welcome: bug fixes, tests, improvements, and documentation. Please follow the [Code of Conduct](CODE_OF_CONDUCT.md).

## Reporting bugs and proposing changes

Before opening an issue, check whether it has already been reported. For a bug, include:

- A small example that reproduces the problem.
- The expected behavior and what actually happened.
- Your go-app version, Go version, operating system, and browser version when relevant.
- Relevant errors or logs.

For substantial features or API changes, open an issue to discuss the approach before starting implementation. Follow the [security policy](SECURITY.md) when reporting a vulnerability.

## Development setup

You will need:

- Go matching the version in [go.mod](go.mod).
- Git, Make, and a POSIX shell for the commands below.
- Chrome or Chromium for browser tests.
- Node.js for the JavaScript and service-worker tests.
- A C compiler and CGO enabled for the native race-detector tests.

Fork and clone the repository, then create a branch for your change. From the repository root, run:

```sh
export PATH="$(go env GOPATH)/bin:$PATH"
make install
make build
```

The makefile expects Go tools to be installed in `$(go env GOPATH)/bin`; leave `GOBIN` unset when using `make install`.

`make install` installs the browser test runner as `go_js_wasm_exec` and installs `benchstat`. It also copies your Go toolchain's `wasm_exec.js` into `pkg/app/gen/wasm_exec.js`. Review that file's diff if you use a different Go version.

`make build` regenerates the HTML API and its tests, then compiles all Go packages.

## Making changes

- Use expressive names and format Go code with `gofmt`.
- Use explicit return values instead of naked returns.
- Follow the surrounding code's conventions and keep changes focused.
- Add regression tests for bug fixes and tests for new behavior.
- Consider both server-side rendering and browser behavior when changing `pkg/app`.
- Explain any public API or compatibility changes in the pull request.

The HTML API and its tests are generated from `pkg/app/gen/html.go`. Edit the generator, then run `make build` and include the resulting changes to `pkg/app/html_gen.go` and `pkg/app/html_gen_test.go`.

## Running tests

Run the full local checks before submitting code changes:

```sh
make test
```

This runs `go vet ./...`, native tests with the race detector, and the `pkg/app` WebAssembly tests in a browser. Browser tests require the setup above; native tests skip browser-only behavior. Some JavaScript tests skip when Node.js is unavailable, so ensure it is installed.

For focused checks while developing:

```sh
# Native tests for one package.
go test ./pkg/app

# Browser tests for the app package.
GOOS=js GOARCH=wasm go test ./pkg/app

# WebAssembly tests for errors and logs, using Node.js as CI does.
GOOS=js GOARCH=wasm go test \
  -exec "$(go env GOROOT)/lib/wasm/go_js_wasm_exec" \
  ./pkg/errors ./pkg/logs
```

The build workflow runs vet, native race tests, and Node.js-backed WebAssembly tests for `pkg/errors` and `pkg/logs`. Run the browser tests locally for changes that affect browser behavior.

## Documentation

The documentation website lives in [go-app-docs](https://github.com/maxence-charriere/go-app-docs). Submit website content and example changes there; API comments and this repository's Markdown files can be updated here.

To preview documentation against your local library changes, clone `go-app-docs` alongside this repository and follow its [local development instructions](https://github.com/maxence-charriere/go-app-docs#local-development). Its build uses the neighboring `go-app` checkout when available.

## Pull requests

Describe the problem, the resulting behavior, and how you tested the change. Link related issues and include a small example or screenshots when they help explain a UI change.

Keep unrelated refactoring and dependency updates separate. Include generated files when their source changes, and mention any checks you could not run. Ensure the build workflow passes and address review feedback before merging.
