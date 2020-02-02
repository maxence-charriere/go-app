build:
	@go build
	@GOOS=js GOARCH=wasm go build -o ./web/app.wasm ./app

run: build
	./demo

deploy: build
	rm demo
	gcloud app deploy . --project=murlok

clean:
	@go clean -v
	@rm ./web/app.wasm
