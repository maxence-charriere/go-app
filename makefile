bootstrap:
	@echo "\033[94m• Setting up go test for wasm to run in the browser\033[00m"
	go get -u github.com/agnivade/wasmbrowsertest
	mv ${GOPATH}/bin/wasmbrowsertest ${GOPATH}/bin/go_js_wasm_exec

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
	@git checkout v6
	@git tag ${VERSION}
	@git push origin ${VERSION}

else
	@echo "\033[94m\n• Releasing version\033[00m"
	@echo "\033[91mVERSION is not defided\033[00m"
	@echo "~> make VERSION=\033[90mv6.0.0\033[00m release"
endif
	

build:
	@echo "\033[94m• Building go-app client\033[00m"
	@GOARCH=wasm GOOS=js go build -o app.wasm ./bin/client
	@echo "\033[94m\n• Building go-app server\033[00m"
	@go build ./bin/server

run: build
	@echo "\033[94m\n• Running go-app server\033[00m"
	@./server

clean:
	@go clean -v ./...
	-@rm server
	-@rm app.wasm
