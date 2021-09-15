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
