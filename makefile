bootstrap:
	@echo "\033[94m• Setting up go test for wasm to run in the browser\033[00m"
	go install github.com/agnivade/wasmbrowsertest@latest
	mv `go env GOPATH`/bin/wasmbrowsertest `go env GOPATH`/bin/go_js_wasm_exec
	go install golang.org/x/tools/cmd/godoc@latest

.PHONY: test
test:
	@echo "\033[94m• Running Go vet\033[00m"
	go vet ./...
	@echo "\033[94m\n• Running Go tests\033[00m"
	go test -race ./...
	@echo "\033[94m\n• Running go wasm tests\033[00m"
	GOARCH=wasm GOOS=js go test ./pkg/app

release: test
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

build:
	@echo "\033[94m• Building go-app documentation PWA\033[00m"
	@godoc -url /pkg/github.com/maxence-charriere/go-app/v9/pkg/app > ./docs/web/documents/reference.html
	# @GOARCH=wasm GOOS=js go build -v -o docs/web/app.wasm ./docs/src
	tinygo build -o docs/web/app.wasm -target wasm ./docs/src
	@echo "\033[94m• Building go-app documentation\033[00m"
	@go build -o docs/documentation ./docs/src

run: build
	@echo "\033[94m• Running go-app documentation server\033[00m"
	@cd docs && ./documentation local

github: build
	@echo "\033[94m• Generating GitHub Pages\033[00m"
	@cd docs && ./documentation github

clean:
	@go clean -v ./...
	-@rm docs/documentation
	@go mod tidy
