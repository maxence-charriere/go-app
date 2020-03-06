# hello-appengine

hello-appengine is a demo that shows how to run a progressive web app created with the [app package](https://github.com/maxence-charriere/go-app) on [Google Cloud App Engine](https://cloud.google.com/appengine).

## Prerequisites

- [Google Cloud SDK](https://cloud.google.com/sdk) installed on your machine
- A [Google Cloud project](https://console.cloud.google.com/cloud-resource-manager)
- An [app.yaml file](https://github.com/maxence-charriere/go-app/tree/master/demo/hello-appengine/app.yaml)

## Build and deploy

Go to the hello-appengine directory:

```sh
cd $GOPATH/src/github.com/maxence-charriere/go-app/demo/hello-appengine
```

Build the hello app:

```sh
GOARCH=wasm GOOS=js go build -o app.wasm ../hello
```

The current directory should look like the following:

```sh
# github.com/maxence-charriere/go-app/demo/hello-appengine
.
├── README.md
├── app.wasm
├── app.yaml
└── main.go

```

Deploy on Google Cloud App Engine:

```sh
gcloud app deploy . --project YOUR_PROJECT_ID
```

See the [live demo](https://demo.murlok.io).

## Contribute

Help to develop the [app](https://github.com/maxence-charriere/go-app) package by becoming a sponsor.
<br>[Become a sponsor](https://opencollective.com/go-app).
