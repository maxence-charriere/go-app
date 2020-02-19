.PHONY: demo
demo:
	GOOS=js GOARCH=wasm go build -o ./demo/hello-local/web/app.wasm ./demo/hello
	go build  -o ./demo/hello-local/hello-local ./demo/hello-local
	cd ./demo/hello-local && ./hello-local

clean:
	@go clean -v ./...
