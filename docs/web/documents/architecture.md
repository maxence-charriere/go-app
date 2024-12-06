## Overview

Like traditional websites, [progressive web apps](https://developers.google.com/web/progressive-web-apps) are provided by a server and are displayed in a web browser. This document provide a description of the different elements which interact each other to operate a PWA.

![architecture diagram](/web/images/architecture.svg)

## Web browser

A web browser is where the PWA is displayed. Here is a list of well-known web browser:

- [Chrome](https://www.google.com/chrome)
- [Safari](https://www.apple.com/safari)
- [Firefox](https://www.mozilla.org/firefox)
- [Electron (Chromium embedded)](https://www.electronjs.org/)

When a user wants to use an app, the web browser requests an [HTML pages](#html-pages) and their associated resources to the server.

Once the required resources are gathered, it displays the app to the user.

## Server

The server is what serves the files to make a go-app progressive web app work into the web browser:

- [HTML pages](#html-pages)
- [Package resources](#package-resources)
- [app.wasm](#app-wasm)
- [Static resources](#static-resources)

It is implemented with the standard [Go HTTP package](https://golang.org/pkg/net/http) and the [Handler](/reference#Handler).

## HTML pages

HTML pages are pages that indicate to [web browsers](#web-browser) what other resources are required to run the progressive web app:

- [Package resources](#package-resources)
- [app.wasm](#app-wasm)
- [Static resources](#static-resources)

They also contain the markup that provides a pre-rendered version of the requested page and that will be replaced by the app dynamic content once [app.wasm](#app-wasm) is loaded.

## Package resources

Package resources are the mandatory resources to run a go-app progressive web app into web browsers. Those resources are:

| Package resource         | Description                                             |
| ------------------------ | ------------------------------------------------------- |
| **wasm_exec.js**         | Script to interop Go and Javascrip APIs.                |
| **app.js**               | Script that loads app.wasm and go-app service workers.  |
| **app-worker.js**        | Script that implements go-app required service workers. |
| **manifest.webmanifest** | Manifest that describes the progressive web app.        |
| **app.css**              | go-app widgets styles.                                  |

They are served by the [server](#server)'s go-app [Handler](/reference#Handler) and are accessible from the root of the app domain. Eg: `/app.js`.

## app.wasm

app.wasm is the binary that contains the UI logic of the progressive web app. It is the app code, built to run on `wasm` architecture.

```bash
GOARCH=wasm GOOS=js go build -o web/app.wasm
```

It is a [static resource](#static-resources) that is **always located at `/web/app.wasm`**, and can be served by the [server](#server) or is available from a remote bucket such as [S3](https://aws.amazon.com/s3) or [Google Cloud Storage](https://cloud.google.com/storage).

Once loaded in a [web browser](#web-browser), it displays the app content and handles user interactions.

## Static resources

Static resource are files such as:

- CSS files
- JS files
- Images
- Videos
- Sounds
- Documents

They are always located within a directory named `web`, and can be served by the [server](#server) or are available from a remote bucket such as [S3](https://aws.amazon.com/s3) or [Google Cloud Storage](https://cloud.google.com/storage).

Static resources are located in a single directory referred to as the `web` directory:

```sh
/web/RESOURCE_NAME
```

By default served by the server, the `web` directory can also be located on a remote bucket such as [S3](https://aws.amazon.com/s3) or [Google Cloud Storage](https://cloud.google.com/storage).

## Next

- [Components](/components)
