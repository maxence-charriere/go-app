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

release:
ifdef VERSION
	@echo "\033[94m\n• Releasing ${VERSION}\033[00m"
	@git tag ${VERSION}
	@git push origin ${VERSION}

else
	@echo "\033[94m\n• Releasing version\033[00m"
	@echo "\033[91mVERSION is not defided\033[00m"
	@echo "~> make VERSION=\033[90mv6.0.0\033[00m release"
endif
	
gen:
	@echo "\033[94m• Generating HTML Syntax\033[00m"
	@go generate ./pkg/app

clean:
	@go clean -v ./...
	@go mod tidy
