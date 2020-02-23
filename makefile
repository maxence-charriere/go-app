.PHONY: demo
demo:
	GOARCH=wasm GOOS=js go build -o ./demo/hello-local/app.wasm ./demo/hello
	go build  -o ./demo/hello-local/hello-local ./demo/hello-local
	cd ./demo/hello-local && ./hello-local

clean:
	@go clean -v ./...
