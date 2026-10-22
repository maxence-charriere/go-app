install:
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" pkg/app/gen/wasm_exec.js
	@echo "\033[94m• Setting up go test for wasm to run in the browser\033[00m"
	go install github.com/agnivade/wasmbrowsertest@latest
	mv `go env GOPATH`/bin/wasmbrowsertest `go env GOPATH`/bin/go_js_wasm_exec
	go install golang.org/x/perf/cmd/benchstat@latest

.PHONY: test
test:
	@echo "\033[94m• Running Go vet\033[00m"
	go vet ./...
	@echo "\033[94m\n• Running Go tests\033[00m"
	go test -race ./...
	@echo "\033[94m\n• Running go wasm tests\033[00m"
	GOARCH=wasm GOOS=js go test ./pkg/app

.PHONY: build
build:
	@echo "\033[94m• Generating HTML Syntax\033[00m"
	@go generate ./pkg/app
	@echo "\033[94m• Building Go packages\033[00m"
	@go build ./...

clean:
	@go clean -v ./...
	@go mod tidy
