# Static resources

Static resources represent resources that are not dynamically generated such as:

- Styles (\*.css)
- Scripts (\*.js)
- Images
- Sounds
- Documents

As it is presented in the [architecture](/architecture#static-resources) section, static resources are located by default in a directory named `web` that is relative to the server executable.

In a real scenario, this is not always the case. Depending on how and where the server is deployed, static resources could also be deployed in another directory or a remote bucket such as [S3](https://aws.amazon.com/s3) or [Google Cloud Storage](https://cloud.google.com/storage).

## Access static resources

Whether there are in a local or remote location, **static resources are always located in a single directory referred to as the `web` directory**:

```go
/web/RESOURCE_NAME
```

### In Handler

Static resources used in a [Handler](/reference#Handler) are usually icons, CSS, and Javascript files.

```go
http.Handle("/", &app.Handler{
	Name:        "Hello",
	Description: "An Hello World! example",
	Icon: app.Icon{
		Default:    "/web/logo.png",       // Specify default favicon.
		AppleTouch: "/web/logo-apple.png", // Specify icon on IOS devices.
	},
	Styles: []string{
		"/web/hello.css", // Loads hello.css file.
	},
	Scripts: []string{
		"/web/hello.js", // Loads hello.js file.
	},
})
```

### In components

```go
func (f *foo) Render() app.UI {
	return app.Img().
		Alt("An image").
		Src("/web/foo.png") // Specify image source to foo.png.
}
```

### In CSS files

Static resources can also be referred to in a CSS file. Here is an example that specifies a background image:

```css
.bg {
  background-image: url("/web/bg.jpg");
}
```

## Setup local web directory

By default, the web directory is located next to the server binary.

```bash
.
├── ...     # Other source files.
├── hello   # Server binary.
└── web     # Web directory.
    └── ... # Static resources.
```

The location of the web directory is changed by setting the [Handler](/reference#Handler) with a [LocalDir](/reference#LocalDir) resource provider:

```go
http.Handle("/", &app.Handler{
	Name:        "Hello",
	Description: "An Hello World! example",
	Resources:   app.LocalDir("/tmp/web"),
})
```

In the example above, static resources must be located in `/tmp/web`.

Note that within the [Handler](/reference#Handler) and the **app**, they will still be accessed by using the pattern `/web/RESOURCE_NAME`.

## Setup remote web directory

When deployed on a cloud provider, it is a common practice to put static resources in a storage service such as [S3](https://aws.amazon.com/s3) or [Google Cloud Storage](https://cloud.google.com/storage). In this scenario, changing the web directory to a remote bucket is done by using the [RemoteBucket](/reference#RemoteBucket) resource provider.

```go
http.Handle("/", &app.Handler{
	Name:        "Hello",
	Description: "An Hello World! example",
	Resources:   app.RemoteBucket("https://storage.googleapis.com/myapp.appspot.com"),
})
```

In the example above, static resources must be located in the [Google Cloud Storage](https://cloud.google.com/storage) bucket with `https://storage.googleapis.com/myapp.appspot.com` URL.

Note that within the [Handler](/reference#Handler) and the **app**, they will still be accessed by using the pattern `/web/RESOURCE_NAME`. You may also have to configure the remote bucket to avoid [CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS) issues.

## Fully static app

Apps built with this package can be generated to be deployed as a fully static website. It is useful in order to be deployed on platforms such as [GitHub Pages](https://pages.github.com). Static website files are generated with the [GenerateStaticWebsite()](/reference#GenerateStaticWebsite) function:

```go
func main() {
	err := app.GenerateStaticWebsite("/test-app", &app.Handler{
		Name:        "Hello",
		Description: "An Hello World! example",
    })

    if err != nil {
        log.Fatal(err)
    }
}
```

When built and launched, the example above generates a full static website in the `/test-app` directory. The generated website will have the following structure:

```bash
.                        # /test-app
├── app-worker.js        # service-worker file (Generated).
├── app.js               # Js support file (Generated).
├── index.html           # Index page (Generated).
├── manifest.webmanifest # PWA manifest (Generated).
├── wasm_exec.js         # Wasm support file (Generated).
└── web                  # Web directory.
    └── app.wasm         # Wasm app (Manually built).
```

Note that `app.wasm` is still built separately with:

```bash
GOARCH=wasm GOOS=js go build -o /test-app/web/app.wasm`
```

### GitHub Pages

The generated static website can be dropped directly in a GitHub repository, either in the root or in the `docs` directory depending on how [GitHub Pages](https://pages.github.com) is configured.

By default, there is no domain associated to the repository GitHub Pages. In that case the [Handler](/reference#Handler) resource provider must be set with [GitHubPages](/reference#GitHubPages):

```go
func main() {
	app.GenerateStaticWebsite("/test-app", &app.Handler{
		Name:        "Hello",
		Description: "An Hello World! example",
		Resources:   app.GitHubPages("REPOSITORY_NAME"),
	})
}
```

This will fix static resources issues cause by repository name that GitHub adds as a prefix in the URL path.

## Next

- [Understand go-app architecture](/architecture)
- [How to create a component](/components)
- [API reference](/reference)
