build: build-app build-server

build-app:
	@GOOS=js GOARCH=wasm go build -o web/demo.wasm ./cmd/demo

build-server:
	@go build ./cmd/demo-server

deploy:
	@cp -r web cmd/demo-server/web
	gcloud app deploy cmd/demo-server 
	@rm -r cmd/demo-server/web

clean:
	-@rm demo-server
	-@rm web/demo.wasm

run: build
	./demo-server
