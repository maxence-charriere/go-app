build:
	@go build
	@GOOS=js GOARCH=wasm go build -o web/demo.wasm ./app

clean:
	-@rm demo
	-@rm web/demo.wasm

run: build
	./demo
