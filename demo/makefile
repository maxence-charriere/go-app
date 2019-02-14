build:
	@go build -o demo ./server 
	@GOOS=js GOARCH=wasm go build -o web/demo.wasm ./client
	@cp client/*.css ./web/


clean:
	-@rm demo
	-@rm -r web

run: build
	./demo
