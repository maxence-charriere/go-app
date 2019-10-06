build:
	@goapp build -v 

deploy: build
	gcloud app deploy cmd/demo-server --project=murlok

clean:
	@goapp clean -v 

run:
	@goapp run -v -b chrome

install:
	@go get -u -v github.com/maxence-charriere/app/cmd/goapp
