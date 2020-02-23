# hello-local

hello-local is a demo that shows how to run a progressive web app created with the [app package](https://github.com/maxence-charriere/app) on your local machine.

## Build and run

```sh
# Go to the hello-local directory:
cd $GOPATH/src/github.com/maxence-charriere/app/demo/hello-local

# Build the hello app:
GOARCH=wasm GOOS=js go build -o app.wasm ../hello

# Build the server:
go build

# Run the server:
./hello-local
```
