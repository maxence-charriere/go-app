# hello

hello is a demo that shows how to use the [app package](https://github.com/maxence-charriere/app) to build a GUI.

## Build

```sh
GOARCH=wasm GOOS=js go build -o app.wasm
```

Note that `app.wasm` binary requires to be moved at the server location that will serve it. See the other hello examples:

- [hello-docker](https://github.com/maxence-charriere/app/tree/master/demo/hello-docker)
- [hello-local](https://github.com/maxence-charriere/app/tree/master/demo/hello-local)
