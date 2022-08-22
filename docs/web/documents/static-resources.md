## Intro

Images and other resources are often used to enhance a user interface. They are referred to as static resources. Here is a list of common static resources:

- Images
- Documents
- Sounds
- CSS files
- JavaScript files

## Access static resources

To work with go-app, **static resources are to be put into a directory called `web`**, by default located next to the [server](/architecture#server) executable. They are then accessible by referring to the resource from the `/web/` prefix:

```go
/web/RESOURCE_NAME
```

You will find below a couple of examples about how static resources are referred to within go-app.

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

## Setup Custom Web directory

By default relative to the working directory, the `web` directory can be configured to be located in other locations such as a different local directory or a remote bucket.

Keep in mind that **wherever the directory is located, static resources will always be accessible from `/web/RESOURCE_NAME` in code**.

### Setup local web directory

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

In the example above, static resources will be located in `/tmp/web/`, but still accessed from `/web/` when referred elsewhere in the [Handler](/reference#Handler) and within components.

### Setup remote web directory

When deployed on a cloud provider, it is a common practice to put static resources in a storage service such as [S3](https://aws.amazon.com/s3) or [Google Cloud Storage](https://cloud.google.com/storage). In this scenario, changing the web directory to a remote bucket is done by using the [RemoteBucket](/reference#RemoteBucket) resource provider.

```go
http.Handle("/", &app.Handler{
	Name:        "Hello",
	Description: "An Hello World! example",
	Resources:   app.RemoteBucket("https://storage.googleapis.com/myapp.appspot.com"),
})
```

In the example above, static resources are located in the [Google Cloud Storage](https://cloud.google.com/storage) bucket, at the `https://storage.googleapis.com/myapp.appspot.com` URL. Static resources will still referred from `/web/` elsewhere in the [Handler](/reference#Handler) and within components.

You may also have to configure the remote bucket to avoid [CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS) issues.

## Next

- [JavaScript Interoperability](/js)
- [Reference](/reference)
