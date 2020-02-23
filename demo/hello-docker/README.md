# hello-docker

hello-docker is a demo that shows how to deploy a progressive web app created with the [app package](https://github.com/maxence-charriere/app) in a Docker container.

## Prerequisites

- [Docker](https://www.docker.com) installed on your machine

## Build and run Docker contrainer

```sh
# Go to the hello-local directory:
cd $GOPATH/src/github.com/maxence-charriere/app/demo/hello-local

# Set dependencies:
go mod init
go mod tidy

# Build the hello app:
GOARCH=wasm GOOS=js go build -o app.wasm ../hello

# Build Docker image:
docker build -t hello-docker .

# Run the Docker container:
docker run -d -p 7000:7000 hello-docker
```
