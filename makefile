.PHONY: demo

bootstrap:
	@echo "\033[94m• Setting up go test for wasm to run in the browser\033[00m"
	go get -u github.com/agnivade/wasmbrowsertest
	mv ${GOPATH}/bin/wasmbrowsertest ${GOPATH}/bin/go_js_wasm_exec

test:
	@echo "\033[94m• Running Go vet\033[00m"
	go vet ./...
	@echo "\033[94m\n• Running Go tests\033[00m"
	go test -race ./...
	@echo "\033[94m\n• Running go wasm tests\033[00m"
	GOARCH=wasm GOOS=js go test ./pkg/app

demo:
	GOARCH=wasm GOOS=js go build -o ./demo/hello-local/app.wasm ./demo/hello
	go build  -o ./demo/hello-local/hello-local ./demo/hello-local
	cd ./demo/hello-local && ./hello-local

clean:
	@go clean -v ./...
