## Lifecycle Overview

Apps created with go-app are WebAssembly binaries that are served through HTTP requests.

Because they are Progressive Web apps, they have to be available for offline mode. Under the hood, this is done by using [service workers](https://developer.mozilla.org/en-US/docs/Web/API/Service_Worker_API) and caching mechanisms, which results in different app loading scenarios.

![app load flow](/web/images/app-lifecycle.svg)

### First loading

App first loading is the app install. It occurs the first time a user goes on the app with a given web browser.

1. Page is downloaded
2. Service worker is downloaded
3. App resources (app.wasm, CSS files and JS files) are downloaded
4. Page, Service worker, and app resources are cached into the web browser

### Recurrent loadings

App recurrent loadings are occurring when a user is coming back on the app.

1. Page is loaded from the web browser cache
2. Service worker is downloaded
3. **Service worker is compared with its cached version => They are identical**
4. App resources are loaded from the web browser cache

### Loading after an app update

App loading after an update is occurring when a user comes back on the app but the app has been modified since his/her last visit.

1. Page is loaded from the web browser cache
2. Service worker is downloaded
3. **Service worker is compared with its cached version => They are different**
4. Page and app resources are downloaded
5. Page, Service worker, and app resources are cached into the web browser

The trigger to update the app is a diff between the cached service worker and the live one. **Once the app is updated, the page still has to be reloaded** to be able to see the modifications since the current version displayed is the cached one.

## Listen for App Updates

Apps are automatically updated in the background when there are modifications.

Since the displayed version of the app is the cached one, a common scenario is to visually notify the user that an update has been downloaded and that it is available by reloading the page.

This is done by implementing the [AppUpdater](/reference#AppUpdater) interface into a component:

```go
// A component that describes a UI.
type littleApp struct {
	app.Compo

	// Field that reports whether an app update is available. False by default.
	updateAvailable bool
}

// OnAppUpdate satisfies the app.AppUpdater interface. It is called when the app
// is updated in background.
func (a *littleApp) OnAppUpdate(ctx app.Context) {
	a.updateAvailable = ctx.AppUpdateAvailable() // Reports that an app update is available.
}

func (a *littleApp) Render() app.UI {
	return app.Main().Body(
		app.H1().Text("A little app"),
		app.P().Text("That only display a text."),

		// Displays an Update button when an update is available.
		app.If(a.updateAvailable,
			app.Button().
				Text("Update!").
				OnClick(a.onUpdateClick),
		),
	)
}

func (a *littleApp) onUpdateClick(ctx app.Context, e app.Event) {
	// Reloads the page to display the modifications.
	ctx.Reload()
}
```

## Next

- [Handling Install](/install)
- [Reference](/reference)
