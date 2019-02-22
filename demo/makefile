build:
	@goapp build -v 

deploy: build
	@cp -r web cmd/demo-server/web
	gcloud app deploy cmd/demo-server 
	@rm -r cmd/demo-server/web

clean:
	@goapp clean -v 

run:
	@goapp run -v -b chrome

install:
	@go get -u -v github.com/maxence-charriere/app/cmd/goapp
