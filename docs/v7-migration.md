# How to migrate from v6 to v7

**go-app** v7 mainly introduced internal improvements. Despite trying to keep the API the same as v6, those changes resulted in API minor modifications that break compatibility.

Here is a guide on how to adapt v6 code to make it work with v7.

## Imports

Go `import` instructions using v6:

```go
import (
    "github.com/maxence-charriere/go-app/v6/pkg/app"
)
```

must be changed to v7:

```go
import (
    "github.com/maxence-charriere/go-app/v7/pkg/app"
)

```

## http.Handler

- `RootDir` field has been removed.
  Handler now uses the `Resources` field to define where app resources and static resources are located. This has a default value that is retro compatible with local static resources. If static resources are located on a remote bucket, use the following:

  ```go
  app.Handler{
      Resources: app.RemoteBucket("BUCKET_URL"),
  }
  ```

- `UseMinimalDefaultStyles` field has been removed.
  **go-app** now use CSS styles that only apply to the loading screen and context menus. The former default styles have been removed since they could conflict with some CSS frameworks.

## Component interfaces

Some interfaces to deals with the lifecycle of components have been modified to take a context as the first argument. Component with methods that satisfies the following interfaces must be updated:

- **Mounter**: `OnMount()` become `OnMount(ctx app.Context)`
- **Navigator**: `OnNav(u *url.URL)` become `OnNav(ctx app.Context u *url.URL)`

## Event handlers

EventHandler function signature has been also modified to take a context as the first argument:

- `func(src app.Value, e app.Event)` become `func(ctx app.Context, e app.Event)`
- Source is now accessed from the context:

  ```go
  // Event handler that retrieve input value when onchange is fired:
  func (c *myCompo) OnChange(ctx app.Context, e app.Event) {
      v := ctx.JSSrc().Get("value")
  }
  ```
