## What is an Action?

**An [action](/reference#Action) is a custom event propagated across the app**. It can be asynchronously handled in a separate goroutine or by any component.

![Actions diagram](/web/images/actions.svg)

## Create

An action is created from a [Context](/reference#Context) by calling `NewAction(name string, tags ...Tagger)`:

```go
func (h *hello) onInputChange(ctx app.Context, e app.Event) {
	ctx.NewAction("greet")
}
```

A payload can also be attached with `NewActionWithValue(name string, v interface{}, tags ...Tagger):`

```go
func (h *hello) onInputChange(ctx app.Context, e app.Event) {
	name := ctx.JSSrc().Get("value").String()

	ctx.NewActionWithValue("greet", name)
}
```

A bit like an HTTP header, additional info can be attached to actions by setting tags:

```go
func (h *hello) onInputChange(ctx app.Context, e app.Event) {
	name := ctx.JSSrc().Get("value").String()

	ctx.NewActionWithValue("greet", name,
		app.T("source", "input"),
		app.T("event", "change"),
	)
}
```

## Handling

Once an [action](/reference#Action) is created, it is propagated across the app. It can then be handled at global and/or component levels with an [ActionHandler](/reference#ActionHandler):

```go
type ActionHandler func(Context, Action)
```

### Global Level

Dealing with actions at a global level is done by registering an [ActionHandler](/reference#ActionHandler) with the [Handle](/reference#Handle) function:

```go
func main() {
	app.Handle("greet", handleGreet) // Registering action handler.

	app.Route("/", &hello{})
	app.RunWhenOnBrowser()

	// ...
}

// Action handler that is called on a separate goroutine when a "greet" action
// is created.
func handleGreet(ctx app.Context, a app.Action) {
	name, ok := a.Value.(string) // Checks if a name was given.
	if !ok {
		fmt.Println("Hello, World")
		return
	}

	fmt.Println("Hello,", name)
}
```

**Executed asynchronously on a separate goroutine**, handling an action globally is **used to centralize and separate the business logic from the UI**.

### Component Level

Actions can also be handled at component level by registering an [ActionHandler](/reference#ActionHandler) from a [Context](/reference#Context):

```go
func (h *hello) OnMount(ctx app.Context) {
	ctx.Handle("greet") // Registering action handler.
}

func (h *hello) onChange(ctx app.Context, e app.EventHandler) {
	name := ctx.JSSrc().Get("value").String()
	ctx.NewActionWithValue("greet", name) // Creating "greet" action.
}

// Action handler that is called on the UI goroutine when a "greet" action is
// created.
func (h *hello) handleGreet(ctx app.Context, a app.Action) {
	name, ok := a.Value.(string) // Checks if a name was given.
	if !ok {
		return
	}
	h.name = name
}
```

**Executed on the UI goroutine**, handling actions from components can help **to send data from a component to another**.

## Next

- [State Management](/states)
- [Reference](/reference)
