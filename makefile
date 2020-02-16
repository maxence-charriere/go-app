.PHONY: demo
demo:
	GOOS=js GOARCH=wasm go build -o ./demo/local/web/app.wasm ./demo/hello
	go build  -o ./demo/local/local ./demo/local
	cd ./demo/local && ./local

clean:
	@go clean -v ./...
