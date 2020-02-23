# hello-local

hello-local is a demo that shows how to run a progressive web app created with the [app package](https://github.com/maxence-charriere/app) on your local machine.

## Build and run

Go to the hello-local directory:

```sh
cd $GOPATH/src/github.com/maxence-charriere/app/demo/hello-local
```

Build the hello app:

```sh
GOARCH=wasm GOOS=js go build -o app.wasm ../hello
```

Build the server:

```sh
go build
```

Your directory should look like the following:

```sh
# github.com/maxence-charriere/app/demo/hello-local
.
├── README.md
├── app.wasm
├── hello-local
└── main.go
```

Run the server:

```sh
./hello-local
```
