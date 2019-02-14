build:
	@go build -o demo ./server 
	@GOOS=js GOARCH=wasm go build -o demo.wasm ./client
	@cp client/*.css .

run: build
	./demo

clean:
	-@rm demo*
	-@rm *.css