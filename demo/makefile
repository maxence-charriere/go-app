build: build-app build-server

build-app:
	@GOOS=js GOARCH=wasm go build -o web/demo.wasm ./bin/demo

build-server:
	@go build ./bin/demo-server

deploy:
	@cp -r web bin/demo-server/web
	gcloud app deploy bin/demo-server 
	@rm -r bin/demo-server/web

clean:
	-@rm demo-server
	-@rm web/demo.wasm

run: build
	./demo-server
